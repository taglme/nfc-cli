package main

import (
	"github.com/taglme/nfc-cli/service"
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
}
