package main

import (
	"dkmission/dkworker"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	worker := dkworker.NewWorker()
	log.Info("dkworker binary")
	go worker.Run()
	utils.ThreadBlock()
}