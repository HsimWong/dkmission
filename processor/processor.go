package processor
import "C"
// #cgo LDFLAGS: -L.
// #cgo linux LDFLAGS: -ldarknet
// #include "processor.h"
import "C"
import (
	"dkmission/utils"
	"fmt"
	"unsafe"
)

type BoundingBox C.struct_BoundingBox

type Processor struct {
	JobProcessIdf *utils.SyncMessenger
}


func (p *Processor) Run() {
	//C.Hello()
	// Using this method for not knowing the exact type of processServer
	var processServer = C.newProcessServer()

	//print()
	//log.Println(reflect.Type(processServer))
	C.loadModel(processServer)
	for {
		jobName := p.JobProcessIdf.Serve().(string)
		typeNum := C.runDetection(processServer, C.CString(jobName))
		var boundingboxSlices []C.struct_BoundingBox = nil
		if typeNum > 0 {
			boundingboxSlices = make([]C.struct_BoundingBox, int(typeNum))
			C.getBoundingBox(processServer, typeNum, unsafe.Pointer(&boundingboxSlices[0]))
			fmt.Println(boundingboxSlices)
		}

		p.JobProcessIdf.Respond(boundingboxSlices)
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
