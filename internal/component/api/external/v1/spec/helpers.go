package spec

func NewErrorResponse(errs ...error) *ErrorResponse {
	errStrings := make([]string, 0, len(errs))
	for _, err := range errs {
		errStrings = append(errStrings, err.Error())
	}
	return &ErrorResponse{
		Errors: errStrings,
	}
}
