package dkmanager

import (
	"context"
	dkworkermesg "dkmission/comm/dkworker"
	"dkmission/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	//"github.com/gammazero/deque"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)



type dispatcher struct {
	//subtasks *deque.Deque
	subtaskChan chan *subtask
	subtasksSig chan bool
	messageWithRegistry *utils.SyncMessenger
}

func NewDispatcher(messageWithRegistry *utils.SyncMessenger) *dispatcher {
	//subtasks := getSubTaskQueue()
	log.Infof("Dispatcher established")
	return &dispatcher{
		subtaskChan: make(chan *subtask),
		subtasksSig: make(chan bool),
		messageWithRegistry: messageWithRegistry,
	}
}

func (dsp *dispatcher) Run() {
	log.Infof("split start running")
	go dsp.split()


	go dsp.dispatch()
}

func (dsp *dispatcher) split() {
	time.Sleep(10 * time.Second)
	for {
		mainUuid := uuid.New().String()
		log.Infof("generating %s", mainUuid)
		for i := 0; i < 5; i++ {
			newSubTask := &subtask{
				mainTaskID:  mainUuid,
				SubtaskID:   uuid.New().String(),
				SubtaskData: nil,
			}
			//dsp.subtasks.PushBack(newSubTask)
			dsp.subtasksSig <- true
			dsp.subtaskChan <- newSubTask

			log.Infof("A new task has been pushed: %s", newSubTask.SubtaskID)
			time.Sleep(20 * time.Second)
		}
	}

}

func (dsp *dispatcher) dispatch()  {
	log.Infof("dispatcher start running")
	for {
		if <- dsp.subtasksSig {
			subtaskInstance := <-dsp.subtaskChan
			log.Infoln("Dispatcher: Received from dispatcher: ")
			log.Infoln("Dispatcher: trying to dispatch it", subtaskInstance.mainTaskID)
			// get target from registry
			log.Infof("Dispatcher: Requesting candidate host\n")
			candidateHostAddr := dsp.messageWithRegistry.Request("Request")
			log.Infof("Dispatcher: Received candidate host: %s", candidateHostAddr)

			conn, err := grpc.Dial(candidateHostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			client := dkworkermesg.NewTaskHandleClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
			//front := dsp.subtasks.Front()
			//subtaskInstance = <-dsp.subtaskChan

			_, err = client.PushTask(ctx, &dkworkermesg.Task{
				ImageData:  nil, // needs further info
				SubTaskID:  subtaskInstance.SubtaskID,
				MainTaskID: subtaskInstance.mainTaskID,
			})
			//
			//
			//
			//
			//
			cancel()
			err = conn.Close()
			if err != nil {
				log.Warnf("Connection to %s close failed",
					candidateHostAddr)
			}





		}
	}
}

//func (dsp *dispatcher) getTarget() string {
//
//}

