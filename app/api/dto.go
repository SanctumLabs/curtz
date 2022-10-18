package api

// ErrorDto represents the error payload that's returned to the client
type ErrorDto struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
