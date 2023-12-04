package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

// the job queue is not currently written in a testable way and just prints
// stuff to stdout so, these tests just have you look at stdout. actual tests
// would have assertions
func TestSerialization(t *testing.T) {
	payload := makeMockPayload()

	fmt.Println(deserializeRequest(payload))
}

func makeMockPayload() []byte {
	myWorkload := Workload{
		WorkloadType: AddJobType,
		Body: AddNumbersBody{
			Int1: 2,
			Int2: 3,
		},
	}
	bytes, err := json.Marshal(myWorkload)
	if err != nil {
		panic(err)
	}
	return bytes
}

func TestAddJob(t *testing.T) {
	queue := JobQueue{}

	addJob := Workload{
		WorkloadType: AddJobType,
		Body: AddNumbersBody{
			Int1: 3,
			Int2: 4,
		},
	}

	queue.performJob(context.Background(), addJob)
}

func TestSubtractJob(t *testing.T) {
	queue := JobQueue{}

	addJob := Workload{
		WorkloadType: SubtractJobType,
		Body: SubtractNumbersBody{
			Int1: 9,
			Int2: 4,
		},
	}

	queue.performJob(context.Background(), addJob)
}

func TestPrintStringJob(t *testing.T) {
	queue := JobQueue{}

	addJob := Workload{
		WorkloadType: PrintStringJobType,
		Body: PrintBody{
			ToPrint: "foo",
		},
	}

	queue.performJob(context.Background(), addJob)
}

func TestShellJobTimesOut(t *testing.T) {
	queue := JobQueue{}

	addJob := Workload{
		WorkloadType: ShellJobType,
		Body: ShellBody{
			ToExecute: "sleep 2",
		},
	}

	queue.performJob(context.Background(), addJob)
}

func TestErrorJob(t *testing.T) {
	queue := JobQueue{}

	addJob := Workload{
		WorkloadType: ErrorJobType,
	}

	queue.performJob(context.Background(), addJob)
}
