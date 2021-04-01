package main

import (
	"log"
	"runtime"

	"dinit/pkg/dinit/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
