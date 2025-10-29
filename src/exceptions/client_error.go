package exceptions

type ClientError struct {
	Message string
}

// TODO: Di Go, kalau method punya pointer receiver (*ClientError), maka: Hanya pointer yang bisa memanggil method itu.
func (e *ClientError) Error() string {
	return e.Message
}
