package transport

type DataSender interface {
}

type DataReceiver interface {
	Receive() ([]byte, error)
	ReadACKorNAK() (bool, error)
}
