package exceptions

type ClientError struct {
	Message string
}

// TODO: Di Go, kalau method punya pointer receiver (*ClientError), maka: Hanya pointer (alamat asli) yang bisa memanggil method itu.
//TODO: Bisa dibilang bahwa *ClientError itu artinya meminta alamat asli dari objek ClientError, bukan salinannya.
func (e *ClientError) Error() string {
	return e.Message
}
