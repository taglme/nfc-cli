package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/f2prateek/train"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//Client represents main client structure. It is used to communicate with server API
type Client struct {
	Adapters AdapterService
	About    AboutService
	Events   EventService
	Snippets SnippetService
	Tags     TagService
	Runs     RunService
	Jobs     JobService
	Licenses LicenseService
	Ws       WsService
}

//New create new client to communicate with server API
func New(host string, interceptors ...train.Interceptor) *Client {
	transport := train.Transport(interceptors...)
	httpClient := &http.Client{
		Transport: transport,
	}
	urlHttp := "http://" + host
	urlWs := "ws://" + host

	return &Client{
		Adapters: newAdapterService(httpClient, urlHttp),
		About:    newAboutService(httpClient, urlHttp),
		Events:   newEventService(httpClient, urlHttp),
		Snippets: newSnippetService(httpClient, urlHttp),
		Tags:     newTagService(httpClient, urlHttp),
		Runs:     newRunService(httpClient, urlHttp),
		Jobs:     newJobService(httpClient, urlHttp),
		Licenses: newLicenseService(httpClient, urlHttp),
		Ws:       newWsService(urlWs),
	}
}

func handleHttpResponseCode(statusCode int, body []byte) (err error) {
	if statusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return
		}
		err = fmt.Errorf("Server responded with an error: %s (%s)", errorResponse.Message, errorResponse.Info)
		return err
	}

	return err
}
