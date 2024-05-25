package utils

type HttpError struct {
	Message    string
	StatusCode int
}

func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) HttpStatusCode() int {
	return e.StatusCode
}
