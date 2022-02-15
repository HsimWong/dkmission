package utils


type SyncMessenger struct {
	Requester chan string
	Responder chan string
}

func NewSyncMessenger() *SyncMessenger {
	return &SyncMessenger{
		Requester: make(chan string),
		Responder: make(chan string),
	}
}

func (s *SyncMessenger) Request(reqMsg string) string {
	s.Requester <- reqMsg
	return <- s.Responder
}

func (s *SyncMessenger) Serve() string {
	return <- s.Requester
}

func (s *SyncMessenger) Respond(resMsg string) {
	s.Responder <- resMsg
}




