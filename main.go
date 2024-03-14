package main

import (
	"envme/cmd"
	"log"
)

var (
	version = "v0.0.1"
)

func main() {
	err := cmd.Execute(version)
	if err != nil {
		log.Fatal(err)
	}
}
