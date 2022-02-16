package main

import (
	"dkmission/utils"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func main() {
	go utils.FileServer("/home/ryan/codes/dkmission/", utils.FileServerPort)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("%s %s\n", event.Name, event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("public")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done


	utils.ThreadBlock()
}
