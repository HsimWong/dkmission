package dkmanager

import (
	"github.com/gammazero/deque"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)



type dispatcher struct {
	subtasks *deque.Deque
	subtasksSig chan bool
}

func NewDispatcher() *dispatcher {
	subtasks := getSubTaskQueue()
	return &dispatcher{subtasks: subtasks, subtasksSig: make(chan bool)}
}

func (dsp *dispatcher) Run() {
	go dsp.split()
	go dsp.dispatch()
}

func (dsp *dispatcher) split() {
	for {
		mainUuid := uuid.New().String()
		for i := 0; i < 5; i++ {
			newSubTask := &subtask{
				mainTaskID:  mainUuid,
				SubtaskID:   uuid.New().String(),
				SubtaskData: nil,
			}
			dsp.subtasks.PushBack(newSubTask)
			dsp.subtasksSig <- true
			log.Info("A new task has been pushed: ", dsp.subtasks.Back())
			time.Sleep(2 * time.Second)
		}
	}

}

func (dsp *dispatcher) dispatch()  {
	for {
		if <- dsp.subtasksSig {
			log.Infoln("Received from dispatcher: ", dsp.subtasks.Front())
		}
	}
}

func (dsp *dispatcher) getTarget() net.Addr {
	return nil
}

