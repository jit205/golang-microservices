package main

import (
	"fmt"
	"log"
	"net/http"
)

const webport = "80"

type Config struct{}

func main() {

	app := Config{}

	log.Printf("stating the broker service on port : %s\n",webport)

	// define the http server 

	server := &http.Server{
		Addr: fmt.Sprintf(":%s",webport),
		Handler: app.routes(),
	}

	//start the server 
	err :=	server.ListenAndServe()

	if err!= nil{
		log.Panic(err)
	}

}