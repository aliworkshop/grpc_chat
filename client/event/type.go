package event

type Type string

const (
	TypeClosed  Type = "CLOSED"
	TypeMessage Type = "MESSAGE"
	TypeJoin    Type = "JOIN"
)
