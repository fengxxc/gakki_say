package policy

type ReplyType int

const (
	Failed ReplyType = iota
	Text
	Image
)

type Reply struct {
	Type ReplyType
	Body []byte
}
