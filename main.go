package main

import (
	"github.com/margostino/gover/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
