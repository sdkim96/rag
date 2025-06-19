package models

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    DTO    `json:"data,omitempty"`
}
