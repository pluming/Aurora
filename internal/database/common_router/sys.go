package common_router

import (
	"github.com/pluming/aurora/config"
	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/IDB"
	"github.com/pluming/aurora/internal/database/consts"
	"github.com/pluming/aurora/internal/database/protocol"
	"github.com/pluming/aurora/internal/database/router"
)

// Ping the server
func Ping(db IDB.DBInstance, args [][]byte) client.Reply {
	if len(args) == 0 {
		return &protocol.PongReply{}
	} else if len(args) == 1 {
		return protocol.MakeStatusReply(string(args[0]))
	} else {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ping' command")
	}
}

// Auth validate client's password
func Auth(c client.Connection, args [][]byte) client.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'auth' command")
	}
	if config.Properties.RequirePass == "" {
		return protocol.MakeErrReply("ERR Client sent AUTH, but no password is set")
	}
	passwd := string(args[0])
	c.SetPassword(passwd)
	if config.Properties.RequirePass != passwd {
		return protocol.MakeErrReply("ERR invalid password")
	}
	return &protocol.OkReply{}
}

func IsAuthenticated(c client.Connection) bool {
	if config.Properties.RequirePass == "" {
		return true
	}
	return c.GetPassword() == config.Properties.RequirePass
}

func init() {
	router.RegisterCommand(consts.CMDPing, Ping, router.NoPrepare, nil, -1, router.FlagReadOnly)
}
