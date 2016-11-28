package main

import (
	"fmt"
	"sync"
	"time"
)

type task struct {
	threadNb, nextThread int
}

var mutexList [5]sync.Mutex

func main() {
	var wg sync.WaitGroup

	taskList := [5]task{{0, 0}, {1, 2}, {2, 3}, {3, 1}, {4, 4}}
	mutexList[2].Lock()
	mutexList[3].Lock()

	for i := 0; i < 5; i++ {
		go startThread(&wg, taskList[i])
		wg.Add(1)
	}
	wg.Wait()
}

func startThread(wg *sync.WaitGroup, t task) {

	defer wg.Done()

	for {
		mutexList[t.threadNb].Lock()

		fmt.Printf("begin task%d\n", t.threadNb)
		task0()
		fmt.Printf("end task%d\n", t.threadNb)

		mutexList[t.nextThread].Unlock()
	}

}

func task0() {
	time.Sleep(2 * time.Second)
}
