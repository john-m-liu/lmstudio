package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	queue := JobQueue{
		queue: make(chan Workload),
	}

	go func() {
		startHttpServer(&queue)
	}()

	queue.doWork()
}

func startHttpServer(queue *JobQueue) {
	workHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			request, err := io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			workload := deserializeRequest(request)
			queue.enqueueWork(workload)

			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	http.HandleFunc("/addWork", workHandler)

	fmt.Println("starting http server")
	http.ListenAndServe(":5000", nil)
	fmt.Println("http server exiting")
}

func deserializeRequest(payload []byte) Workload {
	myPayload := Workload{}
	err := json.Unmarshal(payload, &myPayload)
	if err != nil {
		panic(err)
	}
	return myPayload
}
