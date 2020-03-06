package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/urfave/cli/v2"
	"os"
)

type AppService interface {
	Start() error
	GetHost() string
	SetActionService(actions.ActionService)
}

type appService struct {
	actions    actions.ActionService
	host       string
	adapterStr string
	cliApp     cli.App
	actionsConstr FNewActionService
	nfcClientConstr FNewNfcClient
}

type FNewActionService func(client client.Client) actions.ActionService
type FNewNfcClient func(string, string) client.Client

func New(a FNewActionService, nfc FNewNfcClient) AppService {
	return &appService{
		cliApp: cli.App{
			Name:        "nfc-cli",
			Version:     "v0.0.1",
			Description: "Cross-platform CLI for reading NFC tags ",
			ac
		},
	}
}

func (s *appService) Start() error {
	s.cliApp.Commands = s.getCommands()

	err := s.cliApp.Run(os.Args)

	if err != nil {
		return errors.Wrap(err, "Can't start thee cli application:\n")
	}
	return nil
}

func (s *appService) GetHost() string {
	return s.host
}

func (s *appService) SetActionService(a actions.ActionService) {
	s.actions = a
}

func (s *appService) getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "host",
			Value:       "127.0.0.1:3011",
			Usage:       "Target host and port",
			Destination: &s.host,
		},
		&cli.StringFlag{
			Name:        "adapter",
			Value:       "1",
			Usage:       "Adapter",
			Aliases:     []string{"a"},
			Destination: &s.adapterStr,
		},
	}
}

func (s *appService) getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "version",
			Usage: "Application version",
			Flags: s.getFlags(),
			Action: s.cmdVersion,
		},
		{
			Name:  "adapters",
			Usage: "Get adapters list",
			Flags: s.getFlags(),
			Action: func(c *cli.Context) error {
				fmt.Println(s.host)
				return nil
			},
		},
	}
}

func (s *appService) cmdVersion(ctx *cli.Context) error {
	if s.actions != nil {
		info, err := s.actions.GetVersion()

		fmt.Println(info, err, s.host)
	}
	return nil
}
