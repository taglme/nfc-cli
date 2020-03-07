package main

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/taglme/nfc-cli/service"
	"log"
	"os"
)

func main() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	app := service.New(t)
	err := app.Start()

	if err != nil {
		log.Fatal(err)
	}
}
