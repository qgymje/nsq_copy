package nsqd

const (
	MsgIDLength = 16
)

type MessageID [MsgIDLength]byte
