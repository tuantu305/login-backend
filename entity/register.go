package entity

type RegisterResponse struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
}

type RegisterRequestMsg struct {
	Id   string `json:"id"`
	User User   `json:"user"`
}

type RegisterResponseMsg struct {
	Id       string           `json:"id"`
	Response RegisterResponse `json:"register_response"`
}
