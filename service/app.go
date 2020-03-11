package service

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

type AppService interface {
	Start() error
	SetRepository(ApiService)
}

type appService struct {
	repository ApiService
	cliApp     cli.App

	exitCh chan struct{}

	flagsMap map[string]cli.Flag
	//  below is arguments controlled by ./flags.go
	host    string
	adapter int
	repeat  int
	output  string
	append  bool
	timeout int
	input   string
	auth    string

	cliStartedCb CbCliStarted
}

type CbCliStarted = func(url string)

func New(repository ApiService, cb CbCliStarted) *appService {
	return &appService{
		cliStartedCb: cb,
		repository:   repository,
		cliApp: cli.App{
			Name:        "nfc-cli",
			Version:     "v0.0.1",
			Description: "Cross-platform CLI for reading NFC tags ",
		},
	}
}

func (s *appService) Start() error {
	s.flagsMap = s.getFlagsMap()
	s.cliApp.Commands = s.getCommands()

	sort.Sort(cli.FlagsByName(s.cliApp.Flags))
	sort.Sort(cli.CommandsByName(s.cliApp.Commands))

	return s.cliApp.Run(os.Args)
}

func (s *appService) SetRepository(r ApiService) {
	s.repository = r
}

func (s *appService) getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  models.CommandVersion,
			Usage: "Application version",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
			},
			Action: s.cmdVersion,
		},
		{
			Name:  models.CommandAdapters,
			Usage: "Get adapters list",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
			},
			Action: s.cmdAdapters,
		},
		{
			Name:  models.CommandRead,
			Usage: "Read tag data with NDEF message",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdRead)
			},
		},
		{
			Name:  models.CommandDump,
			Usage: "Dump tag memory",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdDump)
			},
		},
		{
			Name:  models.CommandLock,
			Usage: "Lock tag memory",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdLock)
			},
		},
		{
			Name:  models.CommandFormat,
			Usage: "Lock tag memory",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdFormat)
			},
		},
		{
			Name:  models.CommandRmpwd,
			Usage: "Remove password for tag write acccess",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdRmPwd)
			},
		},
		{
			Name:  models.CommandSetpwd,
			Usage: "Remove password for tag write acccess",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
				s.flagsMap[models.FlagPwd],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdSetPwd)
			},
		},
		{
			Name:  models.CommandTransmit,
			Usage: "Transmit bytes to adapter or tag",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagInput],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
				s.flagsMap[models.FlagTarget],
				s.flagsMap[models.FlagTxBytes],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdTransmit)
			},
		},
		{
			Name:  models.CommandWrite,
			Usage: "Write NDEF message to the tag",
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdWrite)
			},
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagInput],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],

				s.flagsMap[models.FlagNdefType],
				s.flagsMap[models.FlagProtect],

				s.flagsMap[models.FlagNdefTypeRawId],
				s.flagsMap[models.FlagNdefTypeRawTnf],
				s.flagsMap[models.FlagNdefTypeType],
				s.flagsMap[models.FlagNdefTypeRawPayload],
				s.flagsMap[models.FlagNdefTypeUrl],
				s.flagsMap[models.FlagNdefTypeText],
				s.flagsMap[models.FlagNdefTypeLang],
				s.flagsMap[models.FlagNdefUri],
				s.flagsMap[models.FlagNdefTypeAarPackage],
				s.flagsMap[models.FlagNdefTypePhone],
				s.flagsMap[models.FlagNdefTypeVcardAddressCity],
				s.flagsMap[models.FlagNdefTypeVcardAddressCountry],
				s.flagsMap[models.FlagNdefTypeVcardAddressPostalCode],
				s.flagsMap[models.FlagNdefTypeVcardAddressRegion],
				s.flagsMap[models.FlagNdefTypeVcardAddressStreet],
				s.flagsMap[models.FlagNdefTypeVcardEmail],
				s.flagsMap[models.FlagNdefTypeVcardFirstName],
				s.flagsMap[models.FlagNdefTypeVcardLastName],
				s.flagsMap[models.FlagNdefTypeVcardOrganization],
				s.flagsMap[models.FlagNdefTypeVcardPhoneCell],
				s.flagsMap[models.FlagNdefTypeVcardPhoneHome],
				s.flagsMap[models.FlagNdefTypeVcardPhoneWork],
				s.flagsMap[models.FlagNdefTypeTitle],
				s.flagsMap[models.FlagNdefTypeVcardSite],
				s.flagsMap[models.FlagNdefTypeMimeFormat],
				s.flagsMap[models.FlagNdefTypeMimeContent],
				s.flagsMap[models.FlagNdefTypeGeoLat],
				s.flagsMap[models.FlagNdefTypeGeoLon],
			},
		},
	}
}
