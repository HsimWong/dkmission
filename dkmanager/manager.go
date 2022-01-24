package dkmanager

import (
	"github.com/gammazero/deque"
)


type DKManager struct {
	Subtasks *deque.Deque
	SubResults *map[string]*subResult
	dispatcher *dispatcher
	monitor *monitor
	merger *resultMerger
}

func NewDKManager() *DKManager{

	return &DKManager{
		Subtasks:   getSubTaskQueue(),
		SubResults: getSubResults(),
		dispatcher: NewDispatcher(),
		monitor:    nil,
		merger:     nil,
	}
}

func (m *DKManager) Run() {
	go m.dispatcher.Run()
}