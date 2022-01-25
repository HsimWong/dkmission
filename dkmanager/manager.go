package dkmanager

import (
	"github.com/gammazero/deque"
)


type DKManager struct {
	Subtasks *deque.Deque
	SubResults *map[string]*subResult
	registry *Registry
	dispatcher *dispatcher
	monitor *monitor
	merger *resultMerger
}

func NewDKManager() *DKManager{

	return &DKManager{
		Subtasks:   getSubTaskQueue(),
		SubResults: getSubResults(),
		dispatcher: NewDispatcher(),
		registry:	NewRegistry(),
		monitor:    nil,
		merger:     nil,
	}
}

func (m *DKManager) Run() {
	go m.registry.Run()
	go m.dispatcher.Run()
}