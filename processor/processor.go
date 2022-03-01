package processor
import "C"
// #cgo LDFLAGS: -L.
// #cgo linux LDFLAGS: -ldarknet
// #include "processor.h"
import "C"
import (
	"dkmission/comm/dkmanager"
	"dkmission/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"unsafe"
)

type BoundingBox C.struct_BoundingBox

type Processor struct {
	JobProcessIdf *utils.SyncMessenger
}

type SubTaskResult struct {

	BoundingBoxes []BoundingBox

}


func (p *Processor) Run() {
	//C.Hello()
	// Using this method for not knowing the exact type of processServer
	var processServer = C.newProcessServer()

	//print()
	//log.Println(reflect.Type(processServer))
	C.loadModel(processServer)
	log.Debugf("MOdel loading finished")
	for {
		jobName := p.JobProcessIdf.Serve().(string)
		targetNum := C.runDetection(processServer, C.CString(jobName))
		var boundingboxSlices []C.struct_BoundingBox = nil
		if targetNum > 0 {
			boundingboxSlices = make([]C.struct_BoundingBox, int(targetNum))
			C.getBoundingBox(processServer, targetNum, unsafe.Pointer(&boundingboxSlices[0]))
			fmt.Println(boundingboxSlices)
		}
		
		result := make([]*dkmanager.ObjectResult, int(targetNum))
		for i, v := range boundingboxSlices {
			result[i] = &dkmanager.ObjectResult{
				ObjectType: uint32(v.typeID),
				Width:      uint32(v.width * utils.WidthBase),
				Height:     uint32(v.height * utils.HeightBase),
				TopLeftX:   uint32(v.posx * utils.WidthBase - 0.5 * v.width * utils.WidthBase),
				TopLeftY:   uint32(v.posy * utils.HeightBase - 0.5 * v.height * utils.HeightBase),
			}
		}

		p.JobProcessIdf.Respond(result)
	}
	//log.Println("Start processing 1st img")

	//
	//log.Println("Start processing 2nd img")
	//typeNum = C.runDetection(processServer, C.CString("oiltank_317.jpg"))
	//boundingboxSlices = make([]C.struct_BoundingBox, int(typeNum))
	//C.getBoundingBox(processServer, typeNum, unsafe.Pointer(&boundingboxSlices[0]))
	////typeNum := C.runDetection(processServer, C.CString("oiltank_162.JPG"))

	//fmt.Println(boundingboxSlices)

	//return ""
}
