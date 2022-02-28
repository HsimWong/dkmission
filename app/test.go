package main

import (
	"dkmission/processor"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	//var dog C.struct_Dog;
	//for k, v := range C.dogs

	log.SetLevel(log.DebugLevel)
	JobIdentifier := make(chan string)
	p := processor.Processor{JobIdentifier: JobIdentifier}
	go p.Run()
	log.Debugf("Try pushing \"oiltank_317.jpg\" ")
	JobIdentifier <- "oiltank_317.jpg"
	log.Debugf("Sleep for 60 seconds")
	time.Sleep(60 * time.Second)
	log.Debugf("Try pushing \"oiltank_162.JPG\" ")
	JobIdentifier <- "oiltank_162.JPG"
	//print(p.Run())
	utils.ThreadBlock()
}
