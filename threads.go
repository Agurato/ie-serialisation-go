package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type task struct {
	taskName                   string
	threadNb, nextThread, line int
}

type line struct {
	firstTask, deadline int
	start, end          time.Time
}

var mutexList []sync.Mutex
var lineList []line

func main() {
	var wg sync.WaitGroup

	funcs := map[string]func(){"0": task0, "1": task1, "2": task2, "3": task3, "4": task4}

	taskList := []task{}

	taskFile, errFile := os.Open("taskList.txt")
	if errFile != nil {
		return
	}
	lineScanner := bufio.NewScanner(taskFile)
	lineScanner.Split(bufio.ScanLines)

	taskNb, lineNb, currTask, currLine := 0, 0, 0, 0

	for lineScanner.Scan() {
		infos := strings.Split(lineScanner.Text(), ":")
		switch infos[0] {
		case "TASK_NB":
			taskNb, _ = strconv.Atoi(infos[1])
		case "LINE_NB":
			lineNb, _ = strconv.Atoi(infos[1])
		default:
			//lineTasks := strings.Split(infos[1], "-")
			lineList = append(lineList, line{currTask, -1, time.Now(), time.Now()})
			lineInfos := strings.Split(infos[1], "-")
		lineInfosLoop:
			for index, taskName := range lineInfos {
				if taskName == "END" {
					taskList[currTask-1].nextThread = lineList[currLine].firstTask
					lineList[currLine].deadline, _ = strconv.Atoi(lineInfos[index+1])

					break lineInfosLoop
				} else {
					taskList = append(taskList, task{taskName, currTask, currTask + 1, currLine})
					mutexList = append(mutexList, sync.Mutex{})
					if currTask != lineList[currLine].firstTask {
						mutexList[currTask].Lock()
					}
					currTask++
				}
			}
			currLine++
		}
	}

	fmt.Printf("%d tasks, %d lines\n", taskNb, lineNb)
	taskFile.Close()

	for i := 0; i < 5; i++ {
		go startThread(&wg, funcs, taskList[i])
		wg.Add(1)
	}
	wg.Wait()
}

func startThread(wg *sync.WaitGroup, funcs map[string]func(), t task) {
	defer wg.Done()
	taskNameInt, _ := strconv.Atoi(t.taskName)
	for {
		mutexList[t.threadNb].Lock()

		if t.threadNb == lineList[t.line].firstTask {
			fmt.Printf("\x1b[%dmLine %d : task %d begin - timer starts\x1b[0m\n", 31+taskNameInt, t.line, t.threadNb)
			lineList[t.line].start = time.Now()
		} else {
			fmt.Printf("\x1b[%dmLine %d : task %d begin\x1b[0m\n", 31+taskNameInt, t.line, t.threadNb)
		}

		funcs[t.taskName]()

		diff := time.Since(lineList[t.line].start)
		milliDiff := diff.Nanoseconds() / 1000000
		deadline := int64(lineList[t.line].deadline)
		if milliDiff < deadline {
			if t.nextThread == lineList[t.line].firstTask {
				fmt.Printf("\x1b[%dmLine %d : task %d end - line %d ended before deadline (%d < %d)\x1b[0m\n", 31+taskNameInt, t.line, t.threadNb, t.line, milliDiff, deadline)
			} else {
				fmt.Printf("\x1b[%dmLine %d : task %d end\x1b[0m\n", 31+taskNameInt, t.line, t.threadNb)
			}
			mutexList[t.nextThread].Unlock()
		} else {
			fmt.Printf("\x1b[%dmLine %d : task %d end - deadline reached (%d > %d) : line %d stopped at task %d\x1b[0m\n", 31+taskNameInt, t.line, t.threadNb, milliDiff, deadline, t.line, t.threadNb)
			mutexList[lineList[t.line].firstTask].Unlock()
		}
	}
}
