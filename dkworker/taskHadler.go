package dkworker

import (
	"context"
	"dkmission/comm/dkworker"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type TaskHandler struct {
	tasks chan *dkworker.Task
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{tasks: make(chan *dkworker.Task, 2)}
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

func (th *TaskHandler) taskProcess(task *dkworker.Task) {
	// Processing code here
	



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
		for {
			go th.taskProcess(<-th.tasks)
		}
	}()


}