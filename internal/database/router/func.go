package router

import (
	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/IDB"
)

// ExecFunc is IDB for command Executor
// args don't include cmd line
type ExecFunc func(db IDB.DBInstance, args [][]byte) client.Reply

// PreFunc analyses command line when queued command to `multi`
// returns related write keys and read keys
type PreFunc func(args [][]byte) ([]string, []string)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// UndoFunc returns Undo logs for the given command line
// execute from head to tail when Undo
type UndoFunc func(db IDB.DBInstance, args [][]byte) []CmdLine
