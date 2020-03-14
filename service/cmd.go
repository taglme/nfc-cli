package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
)

func (s *appService) cmdVersion(*cli.Context) error {
	s.cliStartedCb(s.host)
	fmt.Printf("CLI version: %s\n", s.cliApp.Version)
	_, err := s.repository.GetVersion()

	return err
}

func (s *appService) cmdAdapters(*cli.Context) error {
	s.cliStartedCb(s.host)
	_, err := s.repository.GetAdapters()

	return err
}

func (s *appService) cmdRead(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandRead,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}

	return err
}

func (s *appService) cmdDump(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}
	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandDump,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdLock(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)
	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandLock,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			JobName:   s.jobName,
		},
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdFormat(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)
	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandFormat,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdRmPwd(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)
	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandRmpwd,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdSetPwd(ctx *cli.Context) error {
	password, err := s.parseHexString(ctx.String(models.FlagPwd))
	if err != nil {
		return errors.Wrap(err, "Can't parse password arg")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)
	var nj interface{}
	_, nj, err = s.repository.AddSetPwdJob(
		models.GenericJobParams{
			Cmd:       models.CommandSetpwd,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
		password,
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdTransmit(ctx *cli.Context) error {
	target := ctx.String(models.FlagTarget)
	if target != "tag" && target != "adapter" {
		return errors.New("Wrong target flag value. Can be either \"tag\" or \"adapter\".")
	}

	txBytes, err := s.parseHexString(ctx.String(models.FlagTxBytes))
	if err != nil {
		return errors.Wrap(err, "Can't parse tx bytes string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddTransmitJob(
		models.GenericJobParams{
			Cmd:       models.CommandTransmit,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
		txBytes,
		target,
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}
	return err
}

func (s *appService) cmdWrite(ctx *cli.Context) error {
	auth, err := s.parseHexString(s.auth)
	if err != nil {
		return errors.Wrap(err, "Can't parse auth string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	payload, err := s.parseNdefPayloadFlags(ctx)
	if err != nil {
		return err
	}

	protect := ctx.Bool(models.FlagProtect)
	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddWriteJob(
		models.GenericJobParams{
			Cmd:       models.CommandTransmit,
			AdapterId: s.adapterId,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName:   s.jobName,
		},
		payload,
		protect,
	)

	s.ongoingJobs.published = 1
	s.ongoingJobs.left = 1

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}

	return err
}

func (s *appService) cmdRun(ctx *cli.Context) error {
	file := ctx.String(models.FlagFile)
	jobsPublished, err := s.repository.AddJobFromFile(s.adapterId, file, models.GenericJobParams{Expire: s.timeout, JobName: s.jobName})

	s.ongoingJobs.published = jobsPublished
	s.ongoingJobs.left = jobsPublished

	return err
}
