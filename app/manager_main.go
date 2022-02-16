package main

import (
	"dkmission/dkmanager"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	manager := dkmanager.NewDKManager()
	go manager.Run()
	utils.ThreadBlock()
}