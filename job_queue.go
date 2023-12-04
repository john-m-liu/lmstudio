package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type JobQueue struct {
	queue chan Workload
}

func (j *JobQueue) enqueueWork(job Workload) {
	j.queue <- job
}

func (j *JobQueue) doWork() {
	const numWorkers = 3
	queueCtx, cancel := context.WithCancel(context.Background())

	// handle ctrl+C to gracefully shut down workers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	// starts 3 workers, each pulling from the same queue
	fmt.Println("starting job queue workers")
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go j.worker(queueCtx, &wg)
	}

	// block until all workers have exited
	wg.Wait()
	fmt.Println("all workers done, job queue exiting")
}

func (j *JobQueue) worker(queueCtx context.Context, wg *sync.WaitGroup) {
	workerCtx, cancel := context.WithCancel(queueCtx)
	defer cancel()
	for {
		select {
		case <-workerCtx.Done():
			fmt.Println("job worker done, exiting")
			wg.Done()
			return
		// in go, channels are threadsafe, so it's fine for 3 workers to access
		// this same queue without a mutex
		case job := <-j.queue:
			j.performJob(workerCtx, job)
		}
	}
}

func (j *JobQueue) performJob(queueCtx context.Context, job Workload) {
	// error handler that is called when a job panics, prevents panics from
	// crashing the job queue
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("runtime error! error is: %v\n", err)
		}
	}()

	switch job.WorkloadType {
	case AddJobType:
		body := job.Body.(AddNumbersBody)
		fmt.Printf("add numbers job: sum is %d\n", body.Int1+body.Int2)
	case SubtractJobType:
		body := job.Body.(SubtractNumbersBody)
		fmt.Printf("subtract numbers job: difference is %d\n", body.Int1-body.Int2)
	case PrintStringJobType:
		body := job.Body.(PrintBody)
		fmt.Printf("print string job: %s\n", body.ToPrint)
	case ShellJobType:
		// time out these jobs after 1 second
		jobCtx, cancel := context.WithTimeout(queueCtx, 1*time.Second)
		defer cancel()

		body := job.Body.(ShellBody)
		fmt.Printf("shell job: %s\n", body.ToExecute)

		cmdParts := strings.Split(body.ToExecute, " ")
		binary := cmdParts[0]
		args := []string{}
		if len(cmdParts) > 1 {
			args = cmdParts[1:]
		}

		// runs arbitrary shell cmd as the same user that ran this
		// binary - very dangerous
		cmd := exec.CommandContext(jobCtx, binary, args...)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error or timeout executing command: %v\n", err)
		}
	case ErrorJobType:
		panic("this is an expected panic from the error job")
	}
}
