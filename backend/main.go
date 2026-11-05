package main

import (
	server "gateway/Server"
	"log"
	"os"

	"github.com/joho/godotenv"
)


func main(){
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    serverAddr := os.Getenv("SERVER_ADDR")

	Server:=server.Initialize(serverAddr)
	Server.StartServer()
}