package tcp

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/client/tcp"
	"github.com/pluming/aurora/internal/database/IDB"
	"github.com/pluming/aurora/internal/database/parser"
	"github.com/pluming/aurora/internal/database/protocol"
	"github.com/pluming/aurora/internal/database/standalone"
	"github.com/pluming/aurora/lib/logger"
	"github.com/pluming/aurora/lib/sync/atomic"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// Handler implements tcp.Handler and serves as a redis server
type Handler struct {
	activeConn sync.Map // *client -> placeholder
	dbServer   IDB.DBServer
	closing    atomic.Boolean // refusing new client and new request
}

func (h *Handler) closeClient(client client.Connection) {
	_ = client.Close()
	h.dbServer.AfterClientClose(client)
	h.activeConn.Delete(client)
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handlerI refuse new connection
		_ = conn.Close()
		return
	}

	cli := tcp.NewConn(conn)
	h.activeConn.Store(cli, 1)

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(cli)
				logger.Info("connection closed: " + cli.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := protocol.MakeErrReply(payload.Err.Error())
			err := cli.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(cli)
				logger.Info("connection closed: " + cli.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*protocol.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}
		result := h.dbServer.Exec(cli, r.Args)
		if result != nil {
			_ = cli.Write(result.ToBytes())
		} else {
			_ = cli.Write(unknownErrReplyBytes)
		}
	}

}

// Close stops handler
func (h *Handler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		cli := key.(client.Connection)
		h.closeClient(cli)
		return true
	})
	h.dbServer.Close()
	return nil
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	var db IDB.DBServer
	db = standalone.NewServer()
	return &Handler{
		dbServer: db,
	}
}
