package dkmanager

type HostInfo struct {

}


type Registry struct{
	hosts HostInfo
}

func (r Registry) Register(c interface{}, info *HostRegisterInfo) (*RegisterResult, error) {
	panic("implement me")
}

func (r Registry) ReportNodeStatus(c interface{}, report *HostReport) (*ReportStatus, error) {
	panic("implement me")
}

//type Registry struct{}
//
//func (Registry) Register(c interface{}, info *HostInfo) (*RegisterResult, error) {
//	panic("implement me")
//}
//
//
//func (r *Registry) ReportNodeStatus()  {
//
//}
//
//func (r *Registry) GetNodesStatus() {
//
//}

