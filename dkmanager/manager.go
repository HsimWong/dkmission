package dkmanager

import (
	"dkmission/utils"
	"github.com/gammazero/deque"
	log "github.com/sirupsen/logrus"
)


type DKManager struct {
	Subtasks *deque.Deque
	SubResults *map[string]*subResult
	registry *Registry
	dispatcher *dispatcher
	monitor *monitor
	merger *resultMerger
	commDispRegi *utils.SyncMessenger
}

func NewDKManager() *DKManager{
	//log.Info("S")
	commDispRegi := utils.NewSyncMessenger()
	return &DKManager{
		Subtasks:   getSubTaskQueue(),
		SubResults: getSubResults(),
		dispatcher: NewDispatcher(commDispRegi),
		registry:	NewRegistry(commDispRegi),
		monitor:    nil,
		merger:     nil,
		commDispRegi: commDispRegi,
	}
}

func (m *DKManager) Run() {
	go m.registry.Run()
	log.Infoln("Trying to start dispatcher")
	go m.dispatcher.Run()
}