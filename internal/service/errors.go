package service

// AppError is a typed application error that carries an HTTP status code.
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}
