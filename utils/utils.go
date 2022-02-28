package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

const (
	TaskImageDir = "data/tasks"
	SubTaskImgDir = "data/subtasks"

	WidthBase = 480
	HeightBase = 480

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

func ReadFromCmd(command string) string{
	output, err := exec.Command("/bin/bash", "-c", command).Output()
	Check(err, "Execution failed")
	return string(output)
}

func GetCWD() string {
	dir, err := os.Getwd()
	Check(err, "getting current working directory failed")
	//fmt.Println(dir)
	return dir
}
