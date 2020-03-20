package repository

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		resp, err := json.Marshal(apiModels.EventResource{
			EventID:     "123",
			Name:        apiModels.EventNameAdapterDiscovery.String(),
			AdapterID:   "123",
			AdapterName: "aname",
			Data:        nil,
			CreatedAt:   "2006-01-02T15:04:05Z",
		})

		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}

		err = c.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			break
		}
	}
}

func eventHandler(e models.Event, data interface{}) {
	fmt.Println("Event received")
}

func TestRepositoryService_RunWsConnection_StopWSConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	nfc := client.New(strings.Replace(s.URL, "http://", "", -1), "en")
	rep := New(&nfc)

	assert.Equal(t, false, rep.client.Ws.IsConnected())
	err := rep.RunWsConnection(eventHandler)
	assert.Nil(t, err)
	assert.Equal(t, true, rep.client.Ws.IsConnected())
	err = rep.StopWsConnection()
	assert.Nil(t, err)
	assert.Equal(t, false, rep.client.Ws.IsConnected())
}

func TestRepositoryService_eventHandler(t *testing.T) {
	nfc := client.New("url", "en")
	rep := New(&nfc)

	e := apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameJobSubmited,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name":     "Job Name",
			"adapter_Name": "Adapter Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e = apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameJobActivated,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name": "Job Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e = apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameRunStarted,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name": "Job Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e = apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameJobDeleted,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name": "Job Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e = apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameJobFinished,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name": "Job Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e = apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameRunError,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"job_name": "Job Name",
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)
}

func TestRepositoryService_eventHandler_run(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/adapters/Adapter%20ID/jobs/Job%20ID", req.URL.String())
		resp, err := json.Marshal(apiModels.JobResource{
			JobID:       "id",
			JobName:     "name",
			AdapterID:   "adid",
			AdapterName: "adname",
			CreatedAt:   "2006-01-02T15:04:05Z",
			Status:      apiModels.JobStatusActive.String(),
			Steps: []apiModels.JobStepResource{{
				Command: apiModels.CommandRemovePassword.String(),
				Params:  apiModels.RemovePasswordParamsResource{},
			}},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()
	nfc := client.New(strings.Replace(server.URL, "http://", "", -1), "en")
	rep := New(&nfc)

	getTags := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{
				"tag_id":       "id",
				"kind":         "kind",
				"href":         "link",
				"type":         apiModels.TagTypeBluetooth.String(),
				"adapter_id":   "id",
				"adapter_name": "name",
				"uid":          "qhIyag==",
				"atr":          "qhIyag==",
				"product":      "product",
				"vendor":       "vendor",
			},
		},
	}

	e := apiModels.Event{
		EventID:     "",
		Name:        apiModels.EventNameRunError,
		AdapterID:   "",
		AdapterName: "Adapter name",
		Data: map[string]interface{}{
			"run_id":       "Run ID",
			"job_id":       "Job ID",
			"job_name":     "Job Name",
			"status":       apiModels.JobRunStatusStarted.String(),
			"adapter_id":   "Adapter ID",
			"adapter_name": "Adapter Name",
			"created_at":   "2020-03-19T16:10:33.580Z",
			"tag": map[string]interface{}{
				"tag_id":       "id",
				"kind":         "kind",
				"href":         "link",
				"type":         apiModels.TagTypeBluetooth.String(),
				"adapter_id":   "id",
				"adapter_name": "name",
				"uid":          "qhIyag==",
				"atr":          "qhIyag==",
				"product":      "product",
				"vendor":       "vendor",
			},
			"results": []interface{}{
				map[string]interface{}{
					"message": "Msg it is",
					"status":  apiModels.CommandStatusSuccess.String(),
					"command": apiModels.CommandGetTags.String(),
					"output":  getTags,
					"params":  map[string]interface{}{},
				},
			},
		},
		CreatedAt: time.Time{},
	}
	rep.eventHandler(e)

	e.Name = apiModels.EventNameRunSuccess
	rep.eventHandler(e)
}
