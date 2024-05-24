package main

import (
	"fmt"
	"log"
)

func main() {
	env, err := LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load and parse env variables: %s", err)
	}

	addr := fmt.Sprintf("%s:%d", env.Host(), env.Port())
	log.Printf("Starting server in %s", addr)

	server := createServer()
	server.Run(addr)
}
