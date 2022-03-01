package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

const (
	TaskImageDir = "data/tasks"
	SubTaskImgDir = "data/subtasks"


	TmpImgDir = "tmp"

	WidthBase = 480
	HeightBase = 480

	HostAvailability = 2

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

func DownloadFile(imgUrl string, imgTarget string) {
	//cmd := spew.Sprintf("wget %s -O %s", )
	cmd := fmt.Sprintf("wget %s -O %s", imgUrl, imgTarget)
	ReadFromCmd(cmd)




	//resp, err := http.Get(imgUrl)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//// 创建一个文件用于保存
	//out, err := os.Create(imgTarget)
	//if err != nil {
	//	panic(err)
	//}
	//defer out.Close()
	//
	//// 然后将响应流和文件流对接起来
	//_, err = io.Copy(out, resp.Body)
	//if err != nil {
	//	panic(err)
	//}
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

func TripleOp(condition bool, a, b interface{}) interface{}{
	if condition {
		return a
	} else {
		return b
	}
}
