package transport

type Sender interface {
}

type Receiver interface {
	Receive() (string, error)
	ReadACKorNAK() (bool, error)
}
