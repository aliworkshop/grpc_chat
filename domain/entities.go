package domain

type Group struct {
	Name    string
	Id      string
	Members []string
}

type Channel struct {
	Name    string
	Id      string
	Admin   string
	Members []string
}
