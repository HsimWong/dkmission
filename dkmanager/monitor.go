package dkmanager

type monitor struct {
	workerAddr string
}

func (m *monitor) Run() {
	//conn, err := grpc.Dial(m.workerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//client := dkworker.NewTaskHandleClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
	//defer cancel()
	//needleValue := rand.Int()
	//respond, err := client.StatusTest(ctx, &dkworker.Needle{NeedleValue: int32(needleValue)})
	//if respond.GetNegNeedleVal() == int32(-needleValue) {
	//	log.Infoln("Successfully inspected ")
	//}

}
