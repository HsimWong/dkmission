package utils

type SyncMessenger struct {
	Requester chan interface{}
	Responder chan interface{}
}

func NewSyncMessenger() *SyncMessenger {

	return &SyncMessenger{
		Requester: make(chan interface{}),
		Responder: make(chan interface{}),
	}
}

func (s *SyncMessenger) Request(reqMsg interface{}) interface{} {
	s.Requester <- reqMsg
	return <- s.Responder
}

func (s *SyncMessenger) Serve() interface{} {
	return <- s.Requester
}

func (s *SyncMessenger) Respond(resMsg interface{}) {
	s.Responder <- resMsg
}




