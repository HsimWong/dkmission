package dkmanager

import (
	"context"
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type HostInfo struct {
	HostStatus Status
	HostName string
	HostAddress string
	HostRegisterTime string
}


type Registry struct{
	hosts map[string]*HostInfo
}



func NewRegistry() *Registry {
	return &Registry{hosts: make(map[string]*HostInfo)}
}

func (r Registry) Register(ctx context.Context, info *HostRegisterInfo) (*RegisterResult, error) {
	if value, err := r.hosts[info.HostName]; err {
		return &RegisterResult{
			Result:       "HostExists",
			RegisterTime: value.HostRegisterTime,
		}, nil
	}
	host := &HostInfo{
		HostStatus:  Status_NOT_READY,
		HostName:    info.GetHostName(),
		HostAddress: info.GetHostAddr(),
	}
	r.hosts[host.HostName] = host
	registerTime := time.StampNano
	return &RegisterResult{
		Result:       "Success",
		RegisterTime: registerTime,
	}, nil
}

func (r Registry) ReportNodeStatus(ctx context.Context, report *HostReport) (*ReportStatus, error) {
	panic("implement me")
}


func (r Registry) ScheduleTask(ctx context.Context, empty *Empty) (*ScheduleResult, error) {
	// A method called from dispatcher
	panic("implement me")
}

func (r Registry) Run() {
	lis, err := net.Listen("tcp", utils.RegistryServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterRegistryServer(s, r)

	log.Infoln("Registry start serving")
	s.Serve(lis)
}

