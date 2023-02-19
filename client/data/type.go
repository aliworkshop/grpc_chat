package data

type Type string

const (
	User    Type = "user"
	Channel Type = "channel"
	Group   Type = "group"
)

func (t Type) String() string {
	return string(t)
}
