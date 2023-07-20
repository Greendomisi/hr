package main

import (
	"fmt"
	"time"

	"m/pkg/task"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

func main() {
	const (
		taskBufNumber = 10
		sleepTime     = time.Second * 3
	)

	taskBufferedCh := make(chan task.Task, taskBufNumber)

	go task.CreateToChan(taskBufferedCh)

	doneTasks := make(chan task.Task)
	undoneTasks := make(chan error)

	go func() {
		// обработка и сортировка тасков
		for t := range taskBufferedCh {
			t.Work()

			t.SortToChannels(doneTasks, undoneTasks)
		}

		close(taskBufferedCh)
	}()

	result := map[int]task.Task{}
	err := []error{}

	go func() {
	Cycle:
		for {
			select {
			case doneTask := <-doneTasks:
				result[doneTask.ID] = doneTask
			case taskError := <-undoneTasks:
				err = append(err, taskError)
			case <-time.After(sleepTime):
				break Cycle
			}
		}

		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(sleepTime)

	printResults(result, err)
}

func printResults(result map[int]task.Task, errs []error) {
	fmt.Println("Undone tasks:")

	for r, err := range errs {
		fmt.Println(r, err)
	}

	fmt.Println("Done tasks:")

	for r, task := range result {
		fmt.Println(r, task)
	}
}
