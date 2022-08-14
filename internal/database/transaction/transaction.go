package transaction

import (
	"strings"

	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/protocol"
	"github.com/pluming/aurora/internal/database/router"
)

// Watch set watching keys
//func Watch(db *single_db.DB, conn client.Connection, args [][]byte) client.Reply {
//	watching := conn.GetWatching()
//	for _, bkey := range args {
//		key := string(bkey)
//		watching[key] = db.GetVersion(key)
//	}
//	return protocol.MakeOkReply()
//}
//
//func execGetVersion(db *single_db.DB, args [][]byte) client.Reply {
//	key := string(args[0])
//	ver := db.GetVersion(key)
//	return protocol.MakeIntReply(int64(ver))
//}
//
//func init() {
//	//IDB.RegisterCommand("GetVer", execGetVersion, readAllKeys, nil, 2, flagReadOnly)
//}
//
//// invoker should lock watching keys
//func isWatchingChanged(db *single_db.DB, watching map[string]uint32) bool {
//	for key, ver := range watching {
//		currentVersion := db.GetVersion(key)
//		if ver != currentVersion {
//			return true
//		}
//	}
//	return false
//}

// StartMulti starts multi-command-transaction
func StartMulti(conn client.Connection) client.Reply {
	if conn.InMultiState() {
		return protocol.MakeErrReply("ERR MULTI calls can not be nested")
	}
	conn.SetMultiState(true)
	return protocol.MakeOkReply()
}

// EnqueueCmd puts command line into `multi` pending queue
func EnqueueCmd(conn client.Connection, cmdLine [][]byte) client.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := router.GetCmdCommand(cmdName)
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}
	if cmd.Prepare == nil {
		return protocol.MakeErrReply("ERR command '" + cmdName + "' cannot be used in MULTI")
	}
	if !router.ValidateArity(cmd.Arity, cmdLine) {
		// difference with redis: we won't enqueue command line with wrong arity
		return protocol.MakeArgNumErrReply(cmdName)
	}
	conn.EnqueueCmd(cmdLine)
	return protocol.MakeQueuedReply()
}

//func execMulti(db *single_db.DB, conn client.Connection) client.Reply {
//	if !conn.InMultiState() {
//		return protocol.MakeErrReply("ERR EXEC without MULTI")
//	}
//	defer conn.SetMultiState(false)
//	cmdLines := conn.GetQueuedCmdLine()
//	return db.ExecMulti(conn, conn.GetWatching(), cmdLines)
//}
//
//// ExecMulti executes multi commands transaction Atomically and Isolated
//func (db *single_db.DB) ExecMulti(conn client.Connection, watching map[string]uint32, cmdLines []CmdLine) client.Reply {
//	// prepare
//	writeKeys := make([]string, 0) // may contains duplicate
//	readKeys := make([]string, 0)
//	for _, cmdLine := range cmdLines {
//		cmdName := strings.ToLower(string(cmdLine[0]))
//		cmd := cmdTable[cmdName]
//		prepare := cmd.prepare
//		write, read := prepare(cmdLine[1:])
//		writeKeys = append(writeKeys, write...)
//		readKeys = append(readKeys, read...)
//	}
//	// set watch
//	watchingKeys := make([]string, 0, len(watching))
//	for key := range watching {
//		watchingKeys = append(watchingKeys, key)
//	}
//	readKeys = append(readKeys, watchingKeys...)
//	db.RWLocks(writeKeys, readKeys)
//	defer db.RWUnLocks(writeKeys, readKeys)
//
//	if isWatchingChanged(db, watching) { // watching keys changed, abort
//		return protocol.MakeEmptyMultiBulkReply()
//	}
//	// execute
//	results := make([]client.Reply, 0, len(cmdLines))
//	aborted := false
//	undoCmdLines := make([][]CmdLine, 0, len(cmdLines))
//	for _, cmdLine := range cmdLines {
//		undoCmdLines = append(undoCmdLines, db.GetUndoLogs(cmdLine))
//		result := db.execWithLock(cmdLine)
//		if protocol.IsErrorReply(result) {
//			aborted = true
//			// don't rollback failed commands
//			undoCmdLines = undoCmdLines[:len(undoCmdLines)-1]
//			break
//		}
//		results = append(results, result)
//	}
//	if !aborted { //success
//		db.addVersion(writeKeys...)
//		return protocol.MakeMultiRawReply(results)
//	}
//	// undo if aborted
//	size := len(undoCmdLines)
//	for i := size - 1; i >= 0; i-- {
//		curCmdLines := undoCmdLines[i]
//		if len(curCmdLines) == 0 {
//			continue
//		}
//		for _, cmdLine := range curCmdLines {
//			db.execWithLock(cmdLine)
//		}
//	}
//	return protocol.MakeErrReply("EXECABORT Transaction discarded because of previous errors.")
//}
//
//// DiscardMulti drops MULTI pending commands
//func DiscardMulti(conn client.Connection) client.Reply {
//	if !conn.InMultiState() {
//		return protocol.MakeErrReply("ERR DISCARD without MULTI")
//	}
//	conn.ClearQueuedCmds()
//	conn.SetMultiState(false)
//	return protocol.MakeOkReply()
//}
//
//// GetUndoLogs return rollback commands
//func (db *single_db.DB) GetUndoLogs(cmdLine [][]byte) []CmdLine {
//	cmdName := strings.ToLower(string(cmdLine[0]))
//	cmd, ok := cmdTable[cmdName]
//	if !ok {
//		return nil
//	}
//	undo := cmd.undo
//	if undo == nil {
//		return nil
//	}
//	return undo(db, cmdLine[1:])
//}
//
//// GetRelatedKeys analysis related keys
//func GetRelatedKeys(cmdLine [][]byte) ([]string, []string) {
//	cmdName := strings.ToLower(string(cmdLine[0]))
//	cmd, ok := cmdTable[cmdName]
//	if !ok {
//		return nil, nil
//	}
//	prepare := cmd.prepare
//	if prepare == nil {
//		return nil, nil
//	}
//	return prepare(cmdLine[1:])
//}
