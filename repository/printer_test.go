package repository

import (
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"testing"
)

func TestApiService_printAdapters(t *testing.T) {
	adapters := []apiModels.Adapter{
		{
			Name:      "Adapter 1",
			AdapterID: "Adapter ID 1",
		},
		{
			Name:      "Adapter 2",
			AdapterID: "Adapter ID 2",
		},
	}

	nfc := client.New("url", "en")
	rep := New(&nfc)
	rep.printAdapters(adapters)
}

func TestApiService_printAppInfo(t *testing.T) {
	appInfo := apiModels.AppInfo{
		Version:   "1.1.2",
		Commit:    "asd23d",
		SDKInfo:   "sdk info",
		Platform:  "darwin",
		BuildTime: "build time",
	}

	nfc := client.New("url", "en")
	rep := New(&nfc)
	rep.printAppInfo(appInfo)
}
