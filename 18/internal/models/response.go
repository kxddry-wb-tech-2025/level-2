package models

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}
