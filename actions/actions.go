package actions

import (
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/taglme/nfc-client/pkg/models"
)

type ActionService interface {
	GetVersion() (models.AppInfo, error)
}

type actionService struct {
	client *client.Client
}

func New(c *client.Client) ActionService {
	return &actionService{
		client: c,
	}
}

func (s *actionService) GetVersion() (models.AppInfo, error) {
	return s.client.About.Get()
}
