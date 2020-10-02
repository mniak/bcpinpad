package transport

type DataSender interface {
}

type DataReceiver interface {
	Receive() (string, error)
	ReadACKorNAK() (bool, error)
}
