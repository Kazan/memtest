package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var serverURL = "http://localhost:4040"

func main() {
	go monitorRuntime()

	data := strings.Repeat("x", 1048576) // 1048576 -> 1MB

	for {
		sendHTTPRequest(buildHTTPRequest(data), getHTTPClient())
		time.Sleep(10 * time.Millisecond)
	}
}

func buildHTTPRequest(data string) *http.Request {
	req, err := http.NewRequest(http.MethodPost, serverURL, strings.NewReader(data))
	if err != nil {
		panic(fmt.Sprintf("failed to build HTTP request: %v", err))
	}

	req.Close = true

	return req
}

func sendHTTPRequest(req *http.Request, cli *http.Client) {
	resp, err := cli.Do(req)
	if err != nil {
		// write tcp [::1]:54091->[::1]:4040: write: protocol wrong type for socket
		// write tcp [::1]:54096->[::1]:4040: write: broken pipe
		return
	}

	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)
}

func getHTTPClient() *http.Client {
	return http.DefaultClient
}

func monitorRuntime() {
	var m runtime.MemStats

	reportTicker := time.NewTicker(5 * time.Second)
	gcTicker := time.NewTicker(20 * time.Second)

	for {
		select {
		case <-reportTicker.C:
			runtime.ReadMemStats(&m)

			fmt.Println(
				fmt.Sprintf(
					"memo-sys: %5v \t\t heap-inuse: %5v \t\t heap-objects: %5v\t\tRoutines: %5v",
					(m.Sys),
					(m.HeapInuse),
					m.HeapObjects,
					runtime.NumGoroutine(),
				),
			)

		case <-gcTicker.C:
			// runtime.GC()
			// fmt.Println("Triggered Garbage Collection")
		}
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
