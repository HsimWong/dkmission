package main

import (
	"dkmission/dkmanager"
	"dkmission/utils"
)

func main() {
	manager := dkmanager.NewDKManager()
	//subtasks := dkmanager.
	go manager.Run()

	utils.ThreadBlock()
}