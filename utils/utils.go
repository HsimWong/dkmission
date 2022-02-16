package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

const (
	TaskImageDir = "data/tasks"
	SubTaskImgDir = "data/subtasks"


	RegistryServerIP = "localhost"
	RegistryServerPort = ":60000"
	ResultServerPort = ":60001"
	WorkerPort = ":60002"
	FileServerPort = ":60003"
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

func FileServer(directory string, address string) {
	file := http.FileServer(http.Dir(directory))
	http.Handle("/", http.StripPrefix("/", file))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		//log.Println(err)
		panic(err)
	}
}