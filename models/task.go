package models

import (
	"strconv"
	"time"
)

type Timer struct {
	Timeout int
}

type Date struct {
	Hour int
	Min  int
	Sec  int
	Date int
	Mon  int
	Year int
}

type Call struct {
	Caller string
	Called string
}

type Task struct {
	User string
	Id   string
	If   interface{}
	Then interface{}
	Loop string
}

var (
	TaskList   map[string]*Task
	runChannel chan Task
)

func init() {

	TaskList = make(map[string]*Task)
	runChannel = make(chan Task, 10000)
	go taskTimer()
}

func AddTask(t Task) string {

	t.Id = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	TaskList[t.Id] = &t

	runChannel <- t
	return t.Id
}

func timeoutDoThen(t int, do interface{}) {

	time.Sleep(time.Duration(t) * time.Second)

	switch do.(type) {
	case Call:
		{
			//call操作
			CallTel(do.(Call).Caller, do.(Call).Called)
		}
	}
}

func taskTimer() {

	for {

		select {
		case task := <-runChannel:
			{
				switch task.If.(type) {
				case Timer:
					timer := task.If.(Timer)
					go timeoutDoThen(timer.Timeout, task.Then)
				}
			}

		}
	}
}
