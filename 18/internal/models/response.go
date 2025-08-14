package models

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  error       `json:"error,omitempty"`
}
