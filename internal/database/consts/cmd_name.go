package consts

//server CMD COMMAND
const (
	CMDAuth         = "AUTH"
	CMDSlaveOf      = "SLAVEOF"
	CMDSubscribe    = "SUBSCRIBE"
	CMDPublish      = "PUBLISH"
	CMDUnsubscribe  = "UNSUBSCRIBE"
	CMDBgRewriteAof = "BGREWRITEAOF"
	CMDRewriteAof   = "REWRITEAOF"
	CMDFlushAll     = "flushall"
	CMDFlushBb      = "FLUSHDB"
	CMDSave         = "SAVE"
	CMDSelect       = "SELECT"
	CMDBgSave       = "BGSAVE"
	CMDCopy         = "COPY"
)

//db 事务命令
const (
	CMDMulti   = "MULTI"
	CMDDiscard = "DISCARD"
	CMDExec    = "EXEC"
	CMDWatch   = "WATCH"
)

const (
	CMDPing = "PING"
)
