package data

type Data struct {
	Action Action      `json:"action"`
	Type   Type        `json:"type"`
	Id     string      `json:"id"`
	Body   interface{} `json:"body,omitempty"`
}
