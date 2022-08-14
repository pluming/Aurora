package client

// Connection represents a connection with client client
type Connection interface {
	Write([]byte) error
	SetPassword(string)
	GetPassword() string

	// Subscribe client should keep its subscribing channels
	Subscribe(channel string)
	UnSubscribe(channel string)
	SubsCount() int
	GetChannels() []string

	// InMultiState used for `Multi` command
	InMultiState() bool
	SetMultiState(bool)
	GetQueuedCmdLine() [][][]byte
	EnqueueCmd([][]byte)
	ClearQueuedCmds()
	GetWatching() map[string]uint32

	// GetDBIndex used for multi database
	GetDBIndex() int
	SelectDB(int)
	// GetRole returns role of conn, such as connection with client, connection with master node
	GetRole() int32
	SetRole(int32)

	Close() error
}
