package actions

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-client/pkg/client"
)

type ActionService interface {
	GetVersion() (string, error)
}

type actionService struct {
	client *client.Client
}

func New(c *client.Client) ActionService {
	return &actionService{
		client: c,
	}
}

func (s *actionService) GetVersion() (res string, err error) {
	info, err := s.client.About.Get()
	if err != nil {
		return res, errors.Wrap(err, "Can't get application info from API: ")
	}

	tmpl := `
Name: %s
Version: %s
Commit: %s
SDK Info: %s
Platform: %s
Build time: %s
CheckSuccess: %t
Supported: %t
Have update: %t
Update version: %s
Update download: %s
Started at: %s
	`

	return fmt.Sprintf(
		tmpl,
		info.Name,
		info.Version,
		info.Commit,
		info.SDKInfo,
		info.Platform,
		info.BuildTime,
		info.CheckSuccess,
		info.Supported,
		info.HaveUpdate,
		info.UpdateVersion,
		info.UpdateDownload,
		info.StartedAt,
	), nil
}
