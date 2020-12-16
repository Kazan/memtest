package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	debug.SetGCPercent(10)

	go monitorRuntime()

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status": 200,
			"result": "ok",
		})
	})

	app.Listen(":4040")
}

func monitorRuntime() {
	var m runtime.MemStats

	reportTicker := time.NewTicker(5 * time.Second)
	// gcTicker := time.NewTicker(20 * time.Second)

	var prev uint64

	prev = 0

	for {
		select {
		case <-reportTicker.C:
			runtime.ReadMemStats(&m)

			fmt.Printf(
				"memo-sys: %5v \t\t heap-inuse: %5v \t\t heap-objects: %5v\t\tGrowth: %5v\t\tRoutines: %5v\n",
				(m.Sys),
				(m.HeapInuse),
				m.HeapObjects,
				int64(m.HeapObjects-prev),
				runtime.NumGoroutine(),
			)

			prev = m.HeapObjects

			// case <-gcTicker.C:
			// 	runtime.GC()
			// 	fmt.Println("Triggered Garbage Collection")
		}
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
