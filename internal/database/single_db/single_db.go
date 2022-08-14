package single_db

import (
	"strings"

	"github.com/pluming/aurora/datastruct/dict"
	"github.com/pluming/aurora/datastruct/lock"
	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/consts"
	"github.com/pluming/aurora/internal/database/protocol"
	"github.com/pluming/aurora/internal/database/router"
	"github.com/pluming/aurora/internal/database/transaction"
)

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
	lockerSize   = 1024
)

// DB stores data and execute user's commands
type DB struct {
	index int
	// key -> DataEntity
	data dict.Dict
	// key -> expireTime (time.Time)
	ttlMap dict.Dict
	// key -> version(uint32)
	versionMap dict.Dict

	// dict.Dict will ensure concurrent-safety of its method
	// use this mutex for complicated command only, eg. rpush, incr ...
	locker *lock.Locks
	//addAof func(CmdLine)
}

func (db *DB) SetIndex(index int) {
	db.index = index
}

// MakeDB create DB instance
func MakeDB(index int) *DB {
	db := &DB{
		data:       dict.MakeConcurrent(dataDictSize),
		ttlMap:     dict.MakeConcurrent(ttlDictSize),
		versionMap: dict.MakeConcurrent(dataDictSize),
		locker:     lock.Make(lockerSize),
		index:      index,
		//addAof:     func(line CmdLine) {},
	}
	return db
}

func (db *DB) Exec(client client.Connection, cmdLine [][]byte) client.Reply {
	cmdName := strings.ToUpper(string(cmdLine[0]))

	switch cmdName {
	case consts.CMDMulti:
		return nil
	case consts.CMDDiscard:
		return nil
	case consts.CMDExec:
		return nil
	case consts.CMDWatch:
		return nil
	}
	if client != nil && client.InMultiState() {
		transaction.EnqueueCmd(client, cmdLine)
		return protocol.MakeQueuedReply()
	}
	//normal command
	return db.ExecNormalCommand(cmdLine)
}

func (db *DB) ExecNormalCommand(cmdLine [][]byte) client.Reply {
	cmdName := strings.ToUpper(string(cmdLine[0]))
	cmd, ok := router.GetCmdCommand(cmdName)
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}
	if !router.ValidateArity(cmd.Arity, cmdLine) {
		return protocol.MakeArgNumErrReply(cmdName)
	}

	prepare := cmd.Prepare
	write, read := prepare(cmdLine[1:])
	db.addVersion(write...)
	db.RWLocks(write, read)
	defer db.RWUnLocks(write, read)
	fun := cmd.Executor
	return fun(db, cmdLine[1:])
}
