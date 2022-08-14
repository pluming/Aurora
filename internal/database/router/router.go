package router

import (
	"strings"

	"github.com/pluming/aurora/internal/client"
	"github.com/pluming/aurora/internal/database/single_db"
)

var cmdTable = make(map[string]*command)

// ExecFunc is IDB for command Executor
// args don't include cmd line
type ExecFunc func(db *single_db.DB, args [][]byte) client.Reply

// PreFunc analyses command line when queued command to `multi`
// returns related write keys and read keys
type PreFunc func(args [][]byte) ([]string, []string)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// UndoFunc returns Undo logs for the given command line
// execute from head to tail when Undo
type UndoFunc func(db *single_db.DB, args [][]byte) []CmdLine

type command struct {
	Executor ExecFunc
	Prepare  PreFunc // return related keys command
	Undo     UndoFunc
	Arity    int // allow number of args, Arity < 0 means len(args) >= -Arity
	Flags    int
}

const (
	FlagWrite    = 0
	FlagReadOnly = 1
)

// RegisterCommand registers a new command
// Arity means allowed number of cmdArgs, Arity < 0 means len(args) >= -Arity.
// for example: the Arity of `get` is 2, `mget` is -2
func RegisterCommand(name string, executor ExecFunc, prepare PreFunc, rollback UndoFunc, arity int, flags int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		Executor: executor,
		Prepare:  prepare,
		Undo:     rollback,
		Arity:    arity,
		Flags:    flags,
	}
}

func isReadOnlyCommand(name string) bool {
	name = strings.ToUpper(name)
	cmd := cmdTable[name]
	if cmd == nil {
		return false
	}
	return cmd.Flags&FlagReadOnly > 0
}

func GetCmdCommand(cmdName string) (*command, bool) {
	command, ok := cmdTable[cmdName]
	return command, ok
}

func ValidateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

func NoPrepare(args [][]byte) ([]string, []string) {
	return nil, nil
}
