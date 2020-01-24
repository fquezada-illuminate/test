package response

type errorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error errorObj `json:"error"`
}

// NewErrorResponse Creates a new instance of ErrorResponse
func NewErrorResponse(statusCode int, message string) ErrorResponse {
	return ErrorResponse{
		Error: errorObj{
			Code:    statusCode,
			Message: message,
		},
	}
}
