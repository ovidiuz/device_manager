package domain

const (
	InvalidRequestBody  = "invalid request body"
	InternalServerError = "internal server error"
)

type HandlerErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHandlerErr(code int, message string) *HandlerErr {
	return &HandlerErr{
		Code:    code,
		Message: message,
	}
}
