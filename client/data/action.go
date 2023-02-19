package data

type Action string

const (
	Message Action = "Message"
	Join    Action = "Join"
)

func (a Action) String() string {
	return string(a)
}
