package main

import (
	"log"
)

func fatal(err error, msg string) {
	if err != nil {
		log.Fatalf("msg: %s", err)
	}
}
