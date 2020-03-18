package main

import (
	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"github.com/taglme/nfc-cli/service"
	"github.com/taglme/nfc-goclient/pkg/client"
	"log"
)

var Version string
var Commit string
var SDKInfo string
var Platform string
var BuildTime string

func main() {
	var nfc *client.Client
	var rep *repository.ApiService
	var app service.AppService

	cbCliStarted := func(url string) {
		nfc = client.New(url, "en")
		rep = repository.New(&nfc)
		app.SetRepository(rep)
	}

	config := opts.Config{
		Version:   Version,
		Commit:    Commit,
		SDK:       SDKInfo,
		Platform:  Platform,
		BuildTime: BuildTime,
	}

	app = service.New(rep, cbCliStarted, config)

	err := app.Start()

	if err != nil {
		log.Fatal(err)
	}
}
