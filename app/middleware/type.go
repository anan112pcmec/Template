package middleware

type RequestData struct {
	Tujuan interface{} `json:"tujuan"`
	Nama   string      `json:"nama"`
	Pw     string      `json:"password"`
}
