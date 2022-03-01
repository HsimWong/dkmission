package dkworker

import (
	"context"
	"dkmission/comm/dkmanager"
	"dkmission/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

//var

type Worker struct {
	localPort string
	hostname string
	taskHandler *TaskHandler
	wkThComm *utils.SyncMessenger
}

func NewWorker() *Worker {
	hostname := uuid.New().String()
	wkThComm := utils.NewSyncMessenger()

	err := os.MkdirAll(utils.TmpImgDir, os.ModePerm)
	utils.Check(err, "Unable to make dir tmp")

	return &Worker{
		localPort: utils.WorkerPort,
		hostname:  hostname,
		taskHandler: NewTaskHandler(wkThComm),
		wkThComm: wkThComm,
	}
}

func(w *Worker) callRegistryFunc(funcName string,
	commInfo interface{})interface{} {
	//
	conn, err := grpc.Dial(utils.RegistryServerIP +
		utils.RegistryServerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := dkmanager.NewRegistryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var respond interface{}
	switch funcName {
	case "Register":
		respond, err = client.Register(ctx,
			commInfo.(*dkmanager.HostRegisterInfo))
		break
	case "ReleaseRequest":
		respond, err = client.ReportResult(ctx,
			commInfo.(*dkmanager.SubTaskResult))

	default:
		break
	}
	return respond
}

func (w *Worker) register() error{
	//return callRegistryFunc()
	log.Infof("Register to %s", utils.RegistryServerIP + utils.RegistryServerPort)
	respond := w.callRegistryFunc("Register", &dkmanager.HostRegisterInfo{
			HostName: "",
			HostPort: utils.WorkerPort,
	}).(*dkmanager.RegisterResult)
	//respond := respondInterface.(*dkmanager.RegisterResult)

	//utils.Check(err, "Calling callRegistryFunc failed")
	if respond.Result == "Success"{
		log.Infof("Successfully registered as port: %s", utils.WorkerPort)
		return nil
	} else if respond.Result == "HostExists" {
		panic("HostAlreadyExists")
	} else {
		panic("RegisterFailed")
	}

}


func (w *Worker) Run() {
	// Start Client Service
	go w.taskHandler.Run()
	go w.register()

	// listening to taskHandler:
	go func() {
		for {
			releaseRequest := w.wkThComm.Serve().(*dkmanager.SubTaskResult)
			releaseRespond := w.callRegistryFunc("ReleaseRequest",
				releaseRequest).(*dkmanager.ReleaseResult)
			if releaseRespond.GetReleaseResult() == "Success" {
				w.wkThComm.Respond("Success")
			} else {
				w.wkThComm.Respond("Failed")
			}
		}

	}()

}
