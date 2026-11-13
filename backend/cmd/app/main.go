package main

import (
	"aibuddy/cmd/server"
)

func main() {
	srv, err := server.New()
	if err != nil{
		return
	}
	if err := srv.Start(); err != nil{
		return 
	}
}

