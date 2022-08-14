package standalone

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/pluming/aurora/config"
	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/IDB"
	"github.com/pluming/aurora/internal/database/common_router"
	"github.com/pluming/aurora/internal/database/consts"
	"github.com/pluming/aurora/internal/database/protocol"
	"github.com/pluming/aurora/internal/database/single_db"
	"github.com/pluming/aurora/lib/logger"
)

type MultiDB struct {
	dbSet []*atomic.Value // *DB
}

func NewServer() *MultiDB {
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	mdb := &MultiDB{
		dbSet: make([]*atomic.Value, config.Properties.Databases),
	}
	for i := range mdb.dbSet {
		singleDB := single_db.MakeDB(i)
		//singleDB.SetIndex(i)
		holder := &atomic.Value{}
		holder.Store(singleDB)
		mdb.dbSet[i] = holder
	}
	return mdb
}

// Exec executes command
// parameter `cmdLine` contains command and its arguments, for example: "set key value"
func (mdb *MultiDB) Exec(c client.Connection, cmdLine [][]byte) (result client.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToUpper(string(cmdLine[0]))

	switch cmdName {
	case consts.CMDAuth:
		return common_router.Auth(c, cmdLine[1:])
	default:
		common_router.IsAuthenticated(c)
	}
	switch cmdName {
	case consts.CMDSlaveOf:
		return nil
	default:
	}

	switch cmdName {
	case consts.CMDSubscribe:
		return nil
	case consts.CMDPublish:
		return nil
	case consts.CMDUnsubscribe:
		return nil
	case consts.CMDBgRewriteAof:
		return nil
	case consts.CMDRewriteAof:
		return nil
	case consts.CMDFlushAll:
		return nil
	case consts.CMDFlushBb:
		return nil
	case consts.CMDSave:
		return nil
	case consts.CMDBgSave:
		return nil
	case consts.CMDSelect:
		return nil
	case consts.CMDCopy:
		return nil
	default:
	}
	// todo: support multi database transaction

	// normal commands
	dbIndex := c.GetDBIndex()
	selectedDB, errReply := mdb.selectDB(dbIndex)
	if errReply != nil {
		return errReply
	}
	return selectedDB.Exec(c, cmdLine)
}

// AfterClientClose does some clean after client close connection
func (mdb *MultiDB) AfterClientClose(c client.Connection) {
	//pubsub.UnsubscribeAll(mdb.hub, c)
}

func (mdb *MultiDB) selectDB(dbIndex int) (IDB.DBInstance, *protocol.StandardErrReply) {
	if dbIndex >= len(mdb.dbSet) || dbIndex < 0 {
		return nil, protocol.MakeErrReply("ERR DB index is out of range")
	}
	return mdb.dbSet[dbIndex].Load().(IDB.DBInstance), nil
}

// Close graceful shutdown database
func (mdb *MultiDB) Close() {
	// stop replication first
	//mdb.replication.close()
	//if mdb.aofHandler != nil {
	//	mdb.aofHandler.Close()
	//}
}
