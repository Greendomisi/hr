package task

import (
	"fmt"
	"time"
)

const timeFormat = time.RFC3339

type Result string

const (
	ResultSuccess Result = "task has been successed"
	ResultError   Result = "something went wrong"
)

// Task - a task to handle on whatever.
type Task struct {
	ID         int
	CreatedAt  string
	ExecutedAt string
	Result     Result
}

func NewTask(id int, createTime string) *Task {
	return &Task{
		CreatedAt: createTime,
		ID:        id,
	}
}

func CreateToChan(taskCh chan Task) {
	const (
		timeError = "some error occurred"
	)

	for {
		var currentTime string

		now := time.Now()
		currentTime = now.Format(timeFormat)

		if now.Nanosecond()%3 != 0 {
			currentTime = timeError
		}

		taskCh <- *NewTask(int(time.Now().Unix()), currentTime)
	}
}

func (task *Task) Work() {
	const (
		sleepTime = time.Millisecond * 150
		afterTime = time.Second * -20
	)

	task.Result = ResultError

	taskCreateTime, _ := time.Parse(timeFormat, task.CreatedAt)
	if taskCreateTime.After(
		time.Now().Add(afterTime)) {
		task.Result = ResultSuccess
	}

	task.ExecutedAt = time.Now().Format(time.RFC3339Nano)

	time.Sleep(sleepTime)
}

func (task Task) SortToChannels(
	doneChan chan Task,
	undoneChan chan error,
) {
	if task.Result == ResultSuccess {
		doneChan <- task

		return
	}

	undoneChan <- fmt.Errorf(`
        Task id %d, 
        create time %s, 
        execute time: %s, 
        error %s`,
		task.ID, task.CreatedAt,
		task.ExecutedAt, task.Result)
}
