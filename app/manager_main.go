package main

import (
	"dkmission/dkmanager"
	"dkmission/utils"
)

func main() {
	manager := dkmanager.NewDKManager()
	go manager.Run()
	utils.ThreadBlock()
}