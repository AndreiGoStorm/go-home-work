package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Wrong number of parameters")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Error reading dir env: %v", err)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}
