package main

import (
	"log"

	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"github.com/taglme/nfc-cli/service"
	"github.com/taglme/nfc-goclient/pkg/client"
)

var Version string
var Commit string
var SDKInfo string
var Platform string
var BuildTime string
var AppID string
var AppSecret string
var AppCert string

func main() {
	var nfc *client.Client
	var rep *repository.RepositoryService
	var app service.AppService

	cbCliStarted := func(url string) {
		if AppID != "" && AppSecret != "" && AppCert != "" {
			privateRSAKey, err := client.PrivateRSAKeyFromB64String(AppSecret)
			if err != nil {
				log.Fatal(err)
			}
			auth := client.NewSigner(AppID, privateRSAKey, AppCert)
			nfc = client.New(url, auth)
		} else {
			nfc = client.New(url)
		}

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
