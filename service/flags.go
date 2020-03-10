package service

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
	"strings"
)



func (s *appService) parseAuthString(pwd string) (res []byte, err error) {
	decoded, err := hex.DecodeString(strings.Replace(pwd, " ", "", -1))
	if err != nil {
		return res, errors.Wrap(err, "Can't decode password string")
	}

	return decoded, nil
}


func (s *appService) getFlagsMap() map[string]cli.Flag {
	return map[string]cli.Flag{
		models.FlagHost: &cli.StringFlag{
			Name:        models.FlagHost,
			Value:       "127.0.0.1:3011",
			Usage:       "Target host and port",
			Destination: &s.host,
		},
		models.FlagAdapter: &cli.IntFlag{
			Name:        models.FlagAdapter,
			Value:       1,
			Usage:       "Adapter",
			Destination: &s.adapter,
		},
		models.FlagRepeat: &cli.IntFlag{
			Name:        models.FlagRepeat,
			Value:       1,
			Usage:       "Number of required repetitions of the task. Optional. If missing, the task is run once",
			Destination: &s.repeat,
		},
		models.FlagOutput: &cli.StringFlag{
			Name:        models.FlagOutput,
			Usage:       "File name for recording the results of the task. Optional. If there is no record of the results is not performed.",
			Destination: &s.output,
		},
		models.FlagAppend: &cli.BoolFlag{
			Name:        models.FlagAppend,
			Value:       false,
			Usage:       "Mode of writing the results to a file. Optional. If append = true, the results are added to the file. If absent or append = false after opening the file, its contents are cleared",
			Destination: &s.append,
		},
		models.FlagTimeout: &cli.IntFlag{
			Name:        models.FlagTimeout,
			Value:       60,
			Usage:       "Job timeout time in seconds. Optional. If absent equals 60",
			Destination: &s.timeout,
		},
		models.FlagInput: &cli.StringFlag{
			Name:        models.FlagInput,
			Usage:       "File name for loading data to form a command. Optional. If absent, data is formed from the arguments of the command. If present, then the command arguments are ignored, data is taken from the file.",
			Destination: &s.input,
		},
		models.FlagAuth: &cli.StringFlag{
			Name:        models.FlagAuth,
			Usage:       "an indication of the need for authorization before starting operations. When this argument is specified, the authorization operation (API CommandAuthPassword API command) is included in the task as the first operation. The authorization operation must be the first in the list of operations. The value of the argument is indicated as an array of bytes in hex format. Example \"03 AD F3 41\"",
			Destination: &s.auth,
		},
	}
}
