package utils

import (
	log "github.com/sirupsen/logrus"
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

func Check(err error, msg string) {
	if err != nil {
		log.Warnf("ErrorOccurs: %v", msg)
	}
}