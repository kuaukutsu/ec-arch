package handler

type ErrorResponse struct {
	Error   string       `json:"error"`
	Context ErrorContext `json:"context"`
}

type ErrorContext struct {
	RequestID string `json:"request_id,omitempty"`
}

func NewError(err string, ctx *ErrorContext) *ErrorResponse {
	if ctx != nil {
		return &ErrorResponse{
			Error:   err,
			Context: *ctx,
		}
	}

	return &ErrorResponse{
		Error: err,
	}
}
