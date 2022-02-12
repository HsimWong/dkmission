package utils

import (
	"sync"
)

const (
	RegistryServerIP = "localhost"
	RegistryServerPort = ":60000"
	ResultServerPort = ":60001"
	WorkerPort = ":60002"
)


func ThreadBlock() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

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




