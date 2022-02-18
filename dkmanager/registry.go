package dkmanager

import (
	"context"
	dkmanagermesg "dkmission/comm/dkmanager"
	dkworkermesg "dkmission/comm/dkworker"
	"dkmission/utils"
	"encoding/json"
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
	HostRegisterTime time.Time
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
				RegisterTime: value.HostRegisterTime.UnixMilli(),
			}, nil
		}
	} else {
		host := &HostInfo{
			HostStatus:       dkmanagermesg.Status_READY,
			HostName:         info.GetHostName(),
			HostAddress:      hostAddr,
			HostRegisterTime: time.Now(),
			HostAvailability:   3,
		}

		r.hosts[host.HostAddress] = host
		r.logRegister(host)
		//registerTime := time.Now().Unix()
	}
	return &dkmanagermesg.RegisterResult{
		Result:       "Success",
		RegisterTime: r.hosts[hostAddr].HostRegisterTime.UnixMilli(),
	}, nil

}

func (r *Registry) logRegister(host *HostInfo) {
	dbInstance := utils.NewDatabase()
	execSql := `insert into nodes(node_address,node_join_time,node_current_status)`
	execSql += ` values(?,?,?)`
	statement, err := dbInstance.DbObject.Prepare(execSql)
	utils.Check(err, "database preparing failed for logRegister")
	_, err = statement.Exec(host.HostAddress, host.HostRegisterTime.String(), "NOT_READY")
	utils.Check(err, "Database execution failed for logRegister")
	//dbInstance.DbObject.Close()
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

	needleValue := rand.Int63()
	respond, err := client.StatusTest(ctx, &dkworkermesg.Needle{NeedleValue: needleValue})
	//log.Infoln(respond.NegNeedleVal)
	prevStatus := r.hosts[workerAddr].HostStatus
	status := ""
	if err == nil {
		if respond.GetNegNeedleVal() == -needleValue {
			log.Infoln("Successfully inspected ")
			status = "READY"
			r.hosts[workerAddr].HostStatus = dkmanagermesg.Status_READY
		} else {
			log.Infof("Original: %d, received: %d", needleValue, respond.NegNeedleVal)
			log.Warnf("Fail to inspect")
			status = "ERROR"
			r.hosts[workerAddr].HostStatus = dkmanagermesg.Status_ERROR
		}
	} else {
		//r.hosts[workerAddr] = nil
		r.hosts[workerAddr].HostStatus = dkmanagermesg.Status_OFFLINE
		log.Warnf("node:%s detected failed", workerAddr)
		status = "OFFLINE"
		// error happens, node should be deleted
	}
	if prevStatus != r.hosts[workerAddr].HostStatus {
		r.logInspectHost(status, workerAddr)
	}


	cancel()
	err = conn.Close()
	log.Infof("connection closed")

	if err != nil {
		log.Warnf("Connection to %s close failed",
			workerAddr)
	}
}

func (r *Registry) inspectHosts() {
	for {
		//log.Infof("Monitoring")
		//log.Warnf("host length: %d", len(r.hosts))
		for _, host := range r.hosts {
			r.monitorNodeStatus(host.HostAddress)
		}
		time.Sleep(2 * time.Second)
	}
}

func (r *Registry) logInspectHost(status string, hostAddress string) {
	sqlCmd := `update nodes set node_current_status
		= ? where node_address = ?`
	db := utils.NewDatabase()
	statement, err := db.DbObject.Prepare(sqlCmd)
	utils.Check(err, "DB prepare failed for logInspectHost")
	_, err = statement.Exec(status, hostAddress)
	utils.Check(err, "DB exec failed for logInspectHost")
	//db.DbObject.Close()

}

type except struct {
	Exception string `json:"exception"`
}
func (r *Registry) respondHostRequest() {
	for {
		log.Infof("respondant start serving...")
		exceptJson := r.messageWithDispatcher.Serve()
		log.Infof(exceptJson)
		exceptionObject := except{}
		err := json.Unmarshal([]byte(exceptJson), &exceptionObject)
		utils.Check(err, "Unmarshalling exception failed")
		log.Infof("EXCEPTION: %s", exceptionObject.Exception)

		log.Infof("Received a request.")
		flagFoundHost := false
		for k, host := range r.hosts {
			if k == exceptionObject.Exception {
				continue
			}
			if host.HostAvailability >= 1 {
				log.Infof("Found an available node: %s", k)
				r.messageWithDispatcher.Respond(k)
				host.HostAvailability --
				flagFoundHost = true
				break
			}
		}
		if !flagFoundHost {
			log.Debugf("Unable to find any suitable node")
			r.messageWithDispatcher.Respond("None")
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
	log.Debugf("Inspector Start serving")
	go r.inspectHosts()
	log.Debugf(	"scheduler start serving")
	go r.respondHostRequest()

}

