package dkworker

import (
	"context"
	"dkmission/comm/dkmanager"
	"dkmission/comm/dkworker"
	"dkmission/processor"
	"dkmission/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type TaskHandler struct {
	tasks chan *dkworker.Task
	msgToWorker *utils.SyncMessenger
}

func NewTaskHandler(msgToWorker *utils.SyncMessenger) *TaskHandler {
	return &TaskHandler{
		tasks: make(chan *dkworker.Task, 2),
		msgToWorker: msgToWorker,
	}
}

func (th TaskHandler) StatusTest(ctx context.Context, needle *dkworker.Needle) (*dkworker.NeedleReply, error) {
	log.Debugf("Received needle: %d", needle.NeedleValue)
	return &dkworker.NeedleReply{NegNeedleVal: -needle.GetNeedleValue()}, nil
}

func (th TaskHandler) PushTask(ctx context.Context, task *dkworker.Task) (*dkworker.TaskPushingReply, error) {
	log.Infof("Received Request: %s", task.GetSubTaskID())
	th.tasks <- task
	return &dkworker.TaskPushingReply{
		TaskPushingReplyContent: "TaskPushingSuccess",
	}, nil
}

func (th *TaskHandler) taskProcess(task *dkworker.Task, messenger *utils.SyncMessenger) {
	// Processing code here
	log.Debugf("Executing subtask: %s", task.SubTaskID)
	//time.Sleep(5 * time.Second)

	result := messenger.Request(task.GetSubTaskID())
	if result != nil {
		fmt.Println()
	}



	log.Debugf("Finished executing subtask: %s", task.SubTaskID)

	//rsp := th.msgToWorker

	rsp := th.msgToWorker.Request(&dkmanager.ReleaseRequest{Subtask_ID: task.SubTaskID}).(string)
	if rsp != "Success" {
		panic("Release Failed:")
	}
}

func (th *TaskHandler) Run() {
	lis, err := net.Listen("tcp", utils.WorkerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	dkworker.RegisterTaskHandleServer(s, *th)

	log.Infoln("Registry start serving")
	go func() {
		err = s.Serve(lis)
		utils.Check(err, "Task Handler serving failed")
	}()

	// Run services for processing images.
	go func() {
		var JobProcIdfs = make([]*utils.SyncMessenger, 2)
		var jobProc [2]*processor.Processor
		for index, messenger := range JobProcIdfs {
			jobProc[index] = &processor.Processor{JobProcessIdf: messenger}
			jobProc[index].Run()
		}
		index := 0

		for {
			go th.taskProcess(<-th.tasks, JobProcIdfs[index])
			index ^= 1
		}
	}()


}