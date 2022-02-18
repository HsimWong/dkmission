package dkmanager

import (
	"context"
	dkworkermesg "dkmission/comm/dkworker"
	"dkmission/utils"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)



type dispatcher struct {
	deployments map[string][]*subtask
	subtaskChan chan *subtask
	subtasksSig chan bool
	messageWithRegistry *utils.SyncMessenger
}

func
NewDispatcher(
messageWithRegistry *utils.SyncMessenger,
) *dispatcher {
	//subtasks := getSubTaskQueue()
	log.Infof("Dispatcher established")
	return &dispatcher{
		deployments: make(map[string][]*subtask),
		subtaskChan: make(chan *subtask, 0xffff),
		subtasksSig: make(chan bool, 0xffff),
		messageWithRegistry: messageWithRegistry,
	}
}

func (dsp *dispatcher) Run() {
	//log.Infof("split start running")
	err := os.MkdirAll(utils.TaskImageDir, os.ModePerm)
	utils.Check(err, "making path failed")
	go utils.FileServer("data/subtasks", "0.0.0.0"+utils.FileServerPort)
	go dsp.subtaskGenerate()
	go dsp.dispatch()
}


func (dsp *dispatcher) taskSplit(src string) []string {
	var subtaskNames []string

	sizeStr := strings.Split(strings.TrimSpace(utils.ReadFromCmd(
		fmt.Sprintf(`identify %s | awk '{ print $3 }'`,
		"\"" + path.Join(utils.TaskImageDir, src)+"\""))), "x")

	fmt.Println(sizeStr)
	width, err := strconv.Atoi(sizeStr[0])
	utils.Check(err, "Converting with error")
	height, err := strconv.Atoi(sizeStr[1])
	utils.Check(err, "Converting height error")

	log.Debugf("Generating subtasks...")
	for widthAnchor := 0; widthAnchor < width; widthAnchor += utils.WidthBase / 2 {
		//if widthAnchor > (width - utils.WidthBase) {
		//	widthAnchor = width - utils.WidthBase
		//}
		for heightAnchor := 0; heightAnchor < height; heightAnchor += utils.HeightBase / 2 {
			subtaskName := uuid.New().String()
			cmd := fmt.Sprintf("convert %s -crop %dx%d+%d+%d %s",
				path.Join(utils.TaskImageDir, src), utils.WidthBase,
				utils.HeightBase, widthAnchor, heightAnchor, path.Join(utils.SubTaskImgDir, subtaskName+".png"))
			//log.Println(cmd)
			utils.ReadFromCmd(cmd)
			subtaskNames = append(subtaskNames, subtaskName)
			newSubTask := &subtask{
				mainTaskID:      src,
				SubtaskID:       subtaskName,
				SubtaskDataPath: "",
				DeployTarget:    "",
				HeightAnchor:    heightAnchor,
				WidthAnchor:     widthAnchor,
			}
			dsp.logSubTask(newSubTask)
			dsp.subtasksSig <- true
			dsp.subtaskChan <- newSubTask
			dsp.subtasksSig <- true
			dsp.subtaskChan <- newSubTask
		}
	}
	log.Debugf("Subtasks generation for %s finished", src)
	return subtaskNames
}

func (dsp *dispatcher) logSubTask(sub *subtask) {
	dbInstance := utils.NewDatabase()
	sqlCmd := "insert into deployment(subtask_ID, subtask_create_time, main_task_ID, width_anchor, height_anchor) values (?,?,?,?,?);"
	statement, err := dbInstance.DbObject.Prepare(sqlCmd)
	utils.Check(err, "database logSubTask not prepared")
	_, err = statement.Exec(sub.SubtaskID, time.Now().String(), sub.mainTaskID, sub.WidthAnchor, sub.HeightAnchor)
	utils.Check(err, "Database operation for logSubTask failed")
}

func (dsp *dispatcher) dealSingleMainTask(src string) {
	mainTaskName := src
	err := os.MkdirAll(path.Join(utils.SubTaskImgDir, mainTaskName), os.ModePerm)
	utils.Check(err, "Failed to make directory")
	dsp.taskSplit(mainTaskName)
}

func (dsp *dispatcher) subtaskGenerate() {
	//time.Sleep(5 * time.Second)
	files, err := ioutil.ReadDir(utils.TaskImageDir)
	utils.Check(err, "Unable to read taskImageDir")

	for _, file := range files {
		if !file.IsDir() {
			//dsp.dealSingleMainTask(path.Join(utils.TaskImageDir, file.Name()))
			dsp.dealSingleMainTask(file.Name())
		}
	}

	// Receive re-pulling expired subtask request

	go func() {
		// infinite loop
		// receive
		// push
		// update pushing time
	}()

	watcher, err := fsnotify.NewWatcher()
	utils.Check(err, "file watcher init failed")
	defer func() {
		err = watcher.Close()
		utils.Check(err, "watcher closing failed")
	}()

	done := make(chan bool) // Not knowing what it is for, but added as official tutorial has it.
	go func() {
		defer close(done)
		for {
			event, ok := <-watcher.Events
			if !ok {
				log.Warnf("watcher events popping failed")
				continue
			}
			if event.Op == fsnotify.Create {
				dsp.dealSingleMainTask(event.Name[len(utils.TaskImageDir):])
			}
		}
	}()
	err = watcher.Add(utils.TaskImageDir)
	utils.Check(err, "Added path to file watcher failed")
	<- done

}

func (dsp *dispatcher) logDispatch(sub *subtask, subtaskRole string)  {
	log.Infof(subtaskRole)
	//sqlCmd := `update deployment set ("?_deploy_target", "?_deploy_time")
	//= ("?","?") where subtask_ID = "?"
//`

	sqlCmd := fmt.Sprintf(`update deployment set (%s_deploy_target, %s_deploy_time) 
			= ("%s", "%s") where subtask_ID = "%s"`, subtaskRole, subtaskRole,
			sub.DeployTarget, time.Now().String(), sub.SubtaskID)
	fmt.Println(sqlCmd)
	dbInstance := utils.NewDatabase()
	statement, err := dbInstance.DbObject.Prepare(sqlCmd)
	utils.Check(err, "database preparing for logDispatch failed")
	//_, err = statement.Exec(subtaskRole, subtaskRole, sub.DeployTarget, time.Now().String(), sub.SubtaskID)
	Result, err := statement.Exec()
	utils.Check(err, "database executing for logDispatch failed")
	log.Println(Result)
}


func (dsp *dispatcher) dispatch()  {
	log.Infof("dispatcher start running")
	for {
		if <- dsp.subtasksSig {

			exception := ""
			role := "main"
			subtaskInstance := <-dsp.subtaskChan
			if len(dsp.deployments[subtaskInstance.SubtaskID]) >= 1 {
				exception = dsp.deployments[subtaskInstance.SubtaskID][0].DeployTarget
				role = "back"
			}

			req := "{\"exception\": \"" + exception + "\"}"
			candidateHostAddr := dsp.messageWithRegistry.Request(req)
			if candidateHostAddr == "None" {
				dsp.subtaskChan <- subtaskInstance
				dsp.subtasksSig <- true
				time.Sleep(3 * time.Second)
				continue
			}
			log.Infof("Dispatcher: Received candidate host: %s", candidateHostAddr)


			subtaskInstance.DeployTarget = candidateHostAddr
			//if candidateHostAddr ==
			log.Infoln("Dispatcher: trying to dispatch it", subtaskInstance.SubtaskID)

			conn, err := grpc.Dial(candidateHostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			client := dkworkermesg.NewTaskHandleClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
			//front := dsp.subtasks.Front()
			//subtaskInstance = <-dsp.subtaskChan

			_, err = client.PushTask(ctx, &dkworkermesg.Task{
				ImageData:  nil, // needs further info
				SubTaskID:  subtaskInstance.SubtaskID,
				MainTaskID: subtaskInstance.mainTaskID,
			})

			dsp.deployments[subtaskInstance.SubtaskID] = append(
				dsp.deployments[subtaskInstance.SubtaskID], subtaskInstance)

			dsp.logDispatch(subtaskInstance, role)

			cancel()
			err = conn.Close()
			if err != nil {
				log.Warnf("Connection to %s close failed",
					candidateHostAddr)
			}
			log.Infof("The task has been dispatched to %s", candidateHostAddr)
		}
	}
}

//func (dsp *dispatcher) getTarget() string {
//
//}

