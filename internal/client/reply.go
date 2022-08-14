package client

// Reply is the IDB of client serialization protocol message
type Reply interface {
	ToBytes() []byte
}
