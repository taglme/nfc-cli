package service

import "github.com/urfave/cli/v2"

type Flag = string

const (
	FlagHost     Flag = "host"
	FlagAdapters      = "adapters"
	FlagRepeat        = "repeat"
	FlagOutput        = "output"
	FlagAppend        = "append"
	FlagTimeout       = "timeout"
	FlagInput         = "input"
	FlagAuth          = "auth"
)

func (s *appService) getFlagsMap() map[string]cli.Flag {
	return map[string]cli.Flag{
		FlagHost: &cli.StringFlag{
			Name:        FlagHost,
			Value:       "127.0.0.1:3011",
			Usage:       "Target host and port",
			Destination: &s.host,
		},
		FlagAdapters: &cli.IntFlag{
			Name:        FlagAdapters,
			Value:       1,
			Usage:       "Adapter",
			Destination: &s.adapter,
		},
		FlagRepeat: &cli.IntFlag{
			Name:        FlagRepeat,
			Value:       1,
			Usage:       "Number of required repetitions of the task. Optional. If missing, the task is run once",
			Destination: &s.repeat,
		},
		FlagOutput: &cli.StringFlag{
			Name:        FlagOutput,
			Usage:       "File name for recording the results of the task. Optional. If there is no record of the results is not performed.",
			Destination: &s.output,
		},
		FlagAppend: &cli.BoolFlag{
			Name:        FlagAppend,
			Value:       false,
			Usage:       "Mode of writing the results to a file. Optional. If append = true, the results are added to the file. If absent or append = false after opening the file, its contents are cleared",
			Destination: &s.append,
		},
		FlagTimeout: &cli.IntFlag{
			Name:        FlagTimeout,
			Value:       60,
			Usage:       "Job timeout time in seconds. Optional. If absent equals 60",
			Destination: &s.timeout,
		},
		FlagInput: &cli.StringFlag{
			Name:        FlagInput,
			Usage:       "File name for loading data to form a command. Optional. If absent, data is formed from the arguments of the command. If present, then the command arguments are ignored, data is taken from the file.",
			Destination: &s.input,
		},
		FlagAuth: &cli.StringFlag{
			Name:        FlagAuth,
			Usage:       "an indication of the need for authorization before starting operations. When this argument is specified, the authorization operation (API CommandAuthPassword API command) is included in the task as the first operation. The authorization operation must be the first in the list of operations. The value of the argument is indicated as an array of bytes in hex format. Example \"03 AD F3 41\"",
			Destination: &s.auth,
		},
	}
}
