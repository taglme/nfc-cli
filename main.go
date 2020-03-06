package main

import (
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-cli/service"
	"github.com/taglme/nfc-client/pkg/client"
	"log"
)

func main() {
	//var adapterStr string
	//var host string
	//var nfcClient *client.Client
	//
	//sort.Sort(cli.FlagsByName(app.Flags))
	//sort.Sort(cli.CommandsByName(app.Commands))

	app := service.New()
	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}

	nfcClient := client.New(app.GetHost(), "en")
	actionService := actions.New(nfcClient)
	app.SetActionService(actionService)
}
