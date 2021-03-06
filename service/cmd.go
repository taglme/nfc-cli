package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/utils"
	"github.com/urfave/cli/v2"
)

func (s *appService) cmdVersion(*cli.Context) error {
	s.cliStartedCb(s.host)
	fmt.Printf("CLI version: %s\n", s.cliApp.Version)
	if len(s.cliApp.Version) > 0 {
		fmt.Printf("   Version: %s\n", s.cliApp.Version)
	}
	if len(s.config.Commit) > 0 {
		fmt.Printf("   Commit: %s\n", s.config.Commit)
	}
	if len(s.config.SDK) > 0 {
		fmt.Printf("   SDK: %s\n", s.config.SDK)
	}
	if len(s.config.Platform) > 0 {
		fmt.Printf("   Platform: %s\n", s.config.Platform)
	}
	if len(s.config.BuildTime) > 0 {
		fmt.Printf("   Build time: %s\n", s.config.BuildTime)
	}

	_, err := s.repository.GetVersion()

	return err
}

func (s *appService) cmdAdapters(*cli.Context) error {
	s.cliStartedCb(s.host)
	_, err := s.repository.GetAdapters(true)

	return err
}

func (s *appService) cmdRead(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdDump(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdLock(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdFormat(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}

	return s.exportData(export, nj)
}

func (s *appService) cmdRmPwd(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdSetPwd(ctx *cli.Context) error {
	password, err := utils.ParseHexString(ctx.String(models.FlagPwd))
	if err != nil {
		return errors.Wrap(err, "Can't parse password arg")
	}

	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdTransmit(ctx *cli.Context) error {
	target := ctx.String(models.FlagTarget)
	if target != "tag" && target != "adapter" {
		return errors.New("Wrong target flag value. Can be either \"tag\" or \"adapter\".")
	}

	txBytes, err := utils.ParseHexString(ctx.String(models.FlagTxBytes))
	if err != nil {
		return errors.Wrap(err, "Can't parse tx bytes string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdWrite(ctx *cli.Context) error {
	auth, err := utils.ParseHexString(s.auth)
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
	if err != nil {
		return err
	}
	s.ongoingJobs.published = s.repeat
	s.ongoingJobs.left = s.repeat

	return s.exportData(export, nj)
}

func (s *appService) cmdRun(ctx *cli.Context) error {
	file := ctx.String(models.FlagFile)
	jobsPublished, err := s.repository.AddJobFromFile(s.adapterId, file, models.GenericJobParams{Expire: s.timeout, JobName: s.jobName})

	s.ongoingJobs.published = jobsPublished
	s.ongoingJobs.left = jobsPublished

	return err
}

func (s *appService) exportData(export bool, data interface{}) error {
	if export && data != nil {
		err := s.writeToFile(s.output, data)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file")
		}
	}

	return nil
}
