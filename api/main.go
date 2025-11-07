package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host, err := resolveEnvVar("HOST")
	if err != nil {
		log.Panic(err)
	}

	port, err := resolveEnvVar("PORT")
	if err != nil {
		log.Panic(err)
	}

	jwtKey, err := resolveEnvVar("JWT_KEY")
	if err != nil {
		log.Panic(err)
	}

	go func() {
		var app App
		app.Initialize(jwtKey)

		addr := fmt.Sprintf("%s:%s", host, port)
		app.Run(addr)
	}()

	// await termination signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down...")
}
