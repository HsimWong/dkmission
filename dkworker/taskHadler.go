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
	"path"
	"time"
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


	imgUrl := "http://" + utils.RegistryServerIP +
		 utils.FileServerPort + "/" + task.GetSubTaskID() + ".png"
	//resp, err := http.Get(imgUrl)
	log.Debugf("Try downloading img %s", imgUrl)
	imgDir := path.Join(utils.TmpImgDir, task.GetSubTaskID() + ".png")
	utils.DownloadFile(imgUrl, imgDir)
	log.Debugf("Downloading Finished")
	result := messenger.Request(imgDir).([]*dkmanager.ObjectResult)
	if result != nil {
		fmt.Println(result)
	}


	log.Debugf("Finished executing subtask: %s", task.SubTaskID)

	//rsp := th.msgToWorker
	rsp := th.msgToWorker.Request(&dkmanager.SubTaskResult{
		Subtask_ID: task.SubTaskID,
		Objects:    result,
	})

	//rsp := th.msgToWorker.Request(&dkmanager.ReleaseRequest{Subtask_ID: task.SubTaskID}).(string)
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
		//var testIDFS = utils.NewSyncMessenger()

		var JobProcIdfs = make([]*utils.SyncMessenger, utils.HostAvailability)
		var jobProc [utils.HostAvailability]*processor.Processor

		for i := 0; i < utils.HostAvailability; i++	 {
			JobProcIdfs[i] = utils.NewSyncMessenger()
			jobProc[i] = &processor.Processor{JobProcessIdf: JobProcIdfs[i]}
			go jobProc[i].Run()
		}


		//for index, messenger := range JobProcIdfs {
		//	JobProcIdfs[index] = utils.NewSyncMessenger()
		//	jobProc[index] = &processor.Processor{JobProcessIdf: messenger}
		//	jobProc[index].Run()
		//}

		time.Sleep(10 * time.Second)
		index := 0
		//
		for {
			log.Debugf("Trying to execute the task")
			go th.taskProcess(<-th.tasks, JobProcIdfs[index])
			index ++
			index %= utils.HostAvailability
			//index = index %
			//index ^= 1
		}
	}()


}