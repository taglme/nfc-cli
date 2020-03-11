package models

type ErrorResponse struct {
	Message string `json:"error_message"`
	Info    string `json:"error_info"`
}
