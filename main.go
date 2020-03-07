package main

import (
	"github.com/taglme/nfc-cli/repository"
	"github.com/taglme/nfc-cli/service"
	"github.com/taglme/nfc-client/pkg/client"
	"log"
)

func main() {
	var nfc *client.Client
	var rep *repository.ApiService
	var app service.AppService

	cbCliStarted := func(url string) {
		nfc = client.New(url, "en")
		rep = repository.New(&nfc)
		app.SetRepository(rep)
	}

	app = service.New(rep, cbCliStarted)

	err := app.Start()

	if err != nil {
		log.Fatal(err)
	}
}
