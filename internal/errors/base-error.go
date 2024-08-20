package errors

type BaseError struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

func (e *BaseError) Error() string {
	return e.Message
}
