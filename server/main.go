package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Fully disable Garbage Collector
	debug.SetGCPercent(-1)

	go monitorRuntime()

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("response")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "response")
	})

	app.Listen(":4040")
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
			runtime.GC()
			fmt.Println("Triggered Garbage Collection")
		}
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
