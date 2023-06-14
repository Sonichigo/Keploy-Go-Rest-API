package main

import (
	"log"

	"github.com/keploy/go-sdk/integrations/kmux"
	"github.com/keploy/go-sdk/keploy"
)

func main() {
	a := &App{}
	err := a.Initialize(
		"postgres",
		"password",
		"postgres")

	if err != nil {
		log.Fatal("Failed to initialize app", err)
	}

	port := "3200"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "my-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})

	a.Router.Use(kmux.MuxMiddleware(k))
	log.Printf("😃 Connected to 3200 port !!")

	a.Run(":3200")
}
