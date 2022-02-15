package dkworker

import (
	"context"
	"dkmission/comm/dkmanager"
	"dkmission/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Worker struct {
	localPort string
	hostname string
	taskHandler *TaskHandler
}

func NewWorker() *Worker {
	hostname := uuid.New().String()
	return &Worker{
		localPort: utils.WorkerPort,
		hostname:  hostname,
		taskHandler: &TaskHandler{},
	}

}

func (w *Worker) register() error{
	log.Infof("Register to %s", utils.RegistryServerIP + utils.RegistryServerPort)
	conn, err := grpc.Dial(utils.RegistryServerIP + utils.RegistryServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := dkmanager.NewRegistryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	respond, err := client.Register(ctx, &dkmanager.HostRegisterInfo{
		HostName: "",
		HostPort: utils.WorkerPort,
	})
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
	w.taskHandler.Run()
	w.register()
}
