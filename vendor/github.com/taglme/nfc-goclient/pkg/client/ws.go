package client

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//EventHandler event handler func specific type
type EventHandler func(e models.Event)

//ErrorHandler error handler func specific type
type ErrorHandler func(e error)

//WsService websocket service interface
type WsService interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	OnEvent(EventHandler)
	ConnString() string
	OnError(ErrorHandler)
	SetLocale(locale string) error
}

type wsService struct {
	url           string
	path          string
	conn          *websocket.Conn
	handlers      []EventHandler
	errorHandlers []ErrorHandler
}

func newWsService(url string) WsService {
	return &wsService{
		url:           url,
		path:          "/ws",
		handlers:      make([]EventHandler, 0),
		errorHandlers: make([]ErrorHandler, 0),
	}
}

// Creating a WS connection and start listening for messages
func (s *wsService) Connect() (err error) {
	s.conn, _, err = websocket.DefaultDialer.Dial(s.url+s.path, nil)
	if err != nil {
		return errors.Wrap(err, "Can't connect to the ws endpoint")
	}

	go s.read()
	return nil
}

// Checks if WS connection is established
func (s *wsService) IsConnected() bool {
	return s.conn != nil
}

// Closing the WS connection and stop listening for messages
func (s *wsService) Disconnect() (err error) {
	err = s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "WS connection closed"))
	if err != nil {
		return errors.Wrap(err, "Error on close WS connection")
	}
	s.conn = nil
	return nil
}

// Sending a message to set locale via WS connection
// locale – locale identifier string
func (s *wsService) SetLocale(locale string) (err error) {
	if s.conn == nil {
		return errors.New("Can't set locale. Connection were not initialized")
	}

	loc, ok := models.StringToLocale(locale)
	if !ok {
		loc = models.LocaleEn
	}
	jobStep := models.JobStep{
		Command: models.CommandSetLocale,
		Params: models.SetLocaleParams{
			Locale: loc,
		},
	}

	jobStepResource := jobStep.ToResource()

	body, err := json.Marshal(jobStepResource)
	if err != nil {
		return errors.Wrap(err, "Error on marshall set locale resource")
	}

	err = s.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return errors.Wrap(err, "Error on send set locale resource")
	}

	return nil
}

// Handle event
// handler – Function which is called once event is emitted
func (s *wsService) OnEvent(handler EventHandler) {
	s.handlers = append(s.handlers, handler)
}

// Handle error
// h – Function which is called once error is emitted
func (s *wsService) OnError(h ErrorHandler) {
	s.errorHandlers = append(s.errorHandlers, h)
}

// Pass event to the registered event handlers
// e – Event
func (s *wsService) eventListener(e models.Event) {
	for _, handler := range s.handlers {
		handler(e)
	}
}

// Pass error to the registered error handlers
// e – error
func (s *wsService) errListener(e error) {
	for _, handler := range s.errorHandlers {
		handler(e)
	}
}

// Waiting for and parsing the messages from WS connection
func (s *wsService) read() {
	defer func() {
		if s.conn != nil {
			s.conn.Close()
		}
		s.conn = nil
	}()

	for {
		if !s.IsConnected() {
			return
		}

		var eventResource models.EventResource
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't read WS message"))
			return
		}
		err = json.Unmarshal(message, &eventResource)
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't unmarshall event resource"))
			continue
		}

		event, err := eventResource.ToEvent()
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't convert event resource to the event model"))

			continue
		}
		s.eventListener(event)
	}
}

// Get WS target url
func (s *wsService) ConnString() string {
	return s.url + s.path
}
