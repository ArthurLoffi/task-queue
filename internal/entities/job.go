package entities

type Job struct {
	Id string `json:"id"`
	Type string `json:"task"`
	Priority string `json:"priority"`
	Status string `jsom:"status"`
	Payload Payload `json:"payload"`
}

// Adicionar um payload mais completo, dependendo do job
// Tipo mandar email, criar usuário, etc...
type Payload struct {
	Title string `json:"title"`
	Text string `json:"text"`
}