package main

import (
	"log"

	"github.com/ryotarai/cmstore/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
