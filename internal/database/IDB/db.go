package IDB

import (
	"github.com/pluming/aurora/internal/client"
)

// DBServer is the IDB for client style storage engine
type DBServer interface {
	Exec(client client.Connection, cmdLine [][]byte) client.Reply
	AfterClientClose(c client.Connection)
	Close()
}

type DBInstance interface {
	Exec(client client.Connection, cmdLine [][]byte) client.Reply
	ExecNormalCommand(cmdLine [][]byte) client.Reply
}

// EmbedDB is the embedding storage engine exposing more methods for complex application
//type EmbedDB interface {
//	DBServer
//	ExecWithLock(conn client.Connection, cmdLine [][]byte) client.Reply
//	ExecMulti(conn client.Connection, watching map[string]uint32, cmdLines []command_func.CmdLine) client.Reply
//	GetUndoLogs(dbIndex int, cmdLine [][]byte) []command_func.CmdLine
//	ForEach(dbIndex int, cb func(key string, data *DataEntity, expiration *time.Time) bool)
//	RWLocks(dbIndex int, writeKeys []string, readKeys []string)
//	RWUnLocks(dbIndex int, writeKeys []string, readKeys []string)
//	GetDBSize(dbIndex int) (int, int)
//}

// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity struct {
	Data interface{}
}
