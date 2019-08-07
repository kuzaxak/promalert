package main

import (
	"log"
	"os"
)

func fatal(err error, msg string) {
	if err != nil {
		log.Fatalf("msg: %s", err)
		os.Exit(1)
	}
}
