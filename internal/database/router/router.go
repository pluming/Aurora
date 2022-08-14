package router

import (
	"strings"
)

var cmdTable = make(map[string]*command)

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
	name = strings.ToUpper(name)
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
