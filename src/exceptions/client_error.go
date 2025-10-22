package exceptions

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}
