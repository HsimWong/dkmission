package dkmanager

import (
	"context"
	dkmanagermesg "dkmission/comm/dkmanager"
	dkworkermesg "dkmission/comm/dkworker"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"math/rand"
	"net"

	"strings"
	"time"
)


type HostInfo struct {
	HostStatus dkmanagermesg.Status
	HostName string
	HostAddress string
	HostRegisterTime int64
	HostAvailability int8
}


type Registry struct{
	hosts map[string]*HostInfo
	messageWithDispatcher *utils.SyncMessenger
}



func NewRegistry(messenger *utils.SyncMessenger) *Registry {
	return &Registry{
		hosts: make(map[string]*HostInfo),
		messageWithDispatcher: messenger,
	}
}

func (r Registry) Register(ctx context.Context, info *dkmanagermesg.HostRegisterInfo) (*dkmanagermesg.RegisterResult, error) {

	p, _ := peer.FromContext(ctx)
	hostIP := strings.Split(p.Addr.String(), ":")[0]

	hostAddr := hostIP + info.GetHostPort()
	//hostAddr := hostIP +
	log.Infof("Received register request from %s", hostAddr)
	if value, err := r.hosts[hostAddr]; err {
		if value.HostStatus == dkmanagermesg.Status_READY {
			log.Warnf("HostExists: %s", value.HostRegisterTime)
			return &dkmanagermesg.RegisterResult{
				Result:       "HostExists",
				RegisterTime: value.HostRegisterTime,
			}, nil
		}
	}

	host := &HostInfo{
		HostStatus:       dkmanagermesg.Status_READY,
		HostName:         info.GetHostName(),
		HostAddress:      hostAddr,
		HostRegisterTime: time.Now().Unix(),
		HostAvailability:   2,
	}

	r.hosts[host.HostAddress] = host
	registerTime := time.Now().Unix()
	return &dkmanagermesg.RegisterResult{
		Result:       "Success",
		RegisterTime: registerTime,
	}, nil
}

func (r Registry) ReportNodeStatus(ctx context.Context, report *dkmanagermesg.HostReport) (*dkmanagermesg.ReportStatus, error) {
	panic("implement me")
}


func (r Registry) ScheduleTask(ctx context.Context, empty *dkmanagermesg.Empty) (*dkmanagermesg.ScheduleResult, error) {
	// A method called from dispatcher

	panic("implement me")
}

func (r *Registry) monitorNodeStatus(workerAddr string)  {
	log.Infof("trying to monitor worker: %s", workerAddr)
	conn, err := grpc.Dial(workerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := dkworkermesg.NewTaskHandleClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
	//defer cancel()
	//defer func() {
	//	err = conn.Close()
	//	if err != nil {
	//		log.Warnf("Connection to %s close failed",
	//			workerAddr)
	//	}
	//}()

	needleValue := rand.Int63()
	respond, err := client.StatusTest(ctx, &dkworkermesg.Needle{NeedleValue: needleValue})
	//log.Infoln(respond.NegNeedleVal)
	if err == nil {
		if respond.GetNegNeedleVal() == -needleValue {
			log.Infoln("Successfully inspected ")
		} else {
			log.Infof("Original: %d, received: %d", needleValue, respond.NegNeedleVal)
			log.Warnf("Fail to inspect")
		}
	} else {
		//r.hosts[workerAddr] = nil
		r.hosts[workerAddr].HostStatus = dkmanagermesg.Status_BROKEN
		log.Warnf("node:%s detected failed", workerAddr)
		// error happens, node should be deleted
	}
	cancel()
	err = conn.Close()
	if err != nil {
		log.Warnf("Connection to %s close failed",
			workerAddr)
	}
}

func (r *Registry) inspectHosts() {
	for {
		log.Infof("Monitoring")
		log.Warnf("host length: %d", len(r.hosts))
		for _, host := range r.hosts {
			r.monitorNodeStatus(host.HostAddress)
		}
		time.Sleep(3 * time.Second)
	}
}

func (r *Registry) respondHostRequest() {
	for {
		log.Infof("respondant start serving...")
		_ = r.messageWithDispatcher.Serve()
		log.Infof("Received a request.")
		for k, host := range r.hosts {
			if host.HostAvailability >= 1 {
				log.Infof("Found an available node: %s", k)
				r.messageWithDispatcher.Respond(k)
				host.HostAvailability --

			}
		}
	}
}

func (r *Registry) Run() {

	lis, err := net.Listen("tcp", utils.RegistryServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	dkmanagermesg.RegisterRegistryServer(s, r)

	log.Infoln("Registry start serving")
	go func() {
		err  = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	//utils.RunGRPCServer(r, utils.RegistryServerPort, RegisterRegistryServer)
	// grpc serving starting finished, going on monitoring each node.
	log.Infof("Inspector Start serving")
	go r.inspectHosts()
	log.Infof("scheduler start serving")
	go r.respondHostRequest()

}

