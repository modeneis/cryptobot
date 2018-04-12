package main

import (
	"os"

	"github.com/modeneis/cryptobot/src/server"
)

//Host eg: localhost:3001 ; or website.com
var Port = os.Getenv("PORT")

func main() {
	if Port == "" {
		Port = "3030"
	}
	server.StartServer(Port)
}
