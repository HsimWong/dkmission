package dkmanager

import (
	"github.com/gammazero/deque"
	"sync"
)

type subtask struct {
	mainTaskID string
	SubtaskID   string
	SubtaskDataPath string
	DeployTarget string
}
var subTaskQueue *deque.Deque

var singletonSTQ sync.Once

func getSubTaskQueue() *deque.Deque {
	singletonSTQ.Do(func() {
		subTaskQueue = &deque.Deque{}
	})
	return subTaskQueue
}


type subResult struct {
	mainTaskID string
	subtaskIDs []string
	anchors [][]int32
	types []string
}

var subResults map[string]*subResult // map[mainTaskID]*subResult
var singletonSubResults sync.Once
func getSubResults() *map[string]*subResult {
	singletonSubResults.Do(func() {
		subResults = make(map[string]*subResult)
	})
	return &subResults
}
