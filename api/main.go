package main

import (
	"fmt"
	"log"
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

	var app App
	app.Initialize(jwtKey)

	fmt.Println("hello world")

	addr := fmt.Sprintf("%s:%s", host, port)
	app.Run(addr)
}
