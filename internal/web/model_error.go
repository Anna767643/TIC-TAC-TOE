package web

// ErrorDTO - единый формат ошибки для клиента.
type errorDTO struct {
	Error string `json:"error"`
}
