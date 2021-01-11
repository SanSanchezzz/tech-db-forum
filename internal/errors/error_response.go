package errors

type JsonError struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error error
	JsonError *JsonError
	StatusCode int
}

func NewErrorResponse(errorConst int, error error) *ErrorResponse {
	return &ErrorResponse{
		Error: error,
		JsonError: &JsonError{
			Message: Errors[errorConst].Error(),
		},
		StatusCode: StatusCode[errorConst],
	}
}
