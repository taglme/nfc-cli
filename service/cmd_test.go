package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/mock"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
)

func Test_exportData(t *testing.T) {
	var rep *repository.RepositoryService
	config := opts.Config{}
	cbCliStarted := func(string) {}

	app := New(rep, cbCliStarted, config)

	err := app.exportData(false, nil)
	assert.Nil(t, err)

	err = app.exportData(true, nil)
	assert.Nil(t, err)

	err = app.exportData(true, "any data")
	assert.Error(t, err)

	app.output = "cmd_test_file.json"
	err = app.exportData(true, "any data")
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdVersion(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)

	err := app.cmdVersion(&cli.Context{})
	assert.Nil(t, err)
}

func Test_cmdAdapters(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)

	err := app.cmdVersion(&cli.Context{})
	assert.Nil(t, err)
}

func Test_cmdRead(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)

	err := app.cmdVersion(&cli.Context{})
	assert.Nil(t, err)
}

func Test_cmdDump(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandDump, "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdLock(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandLock, "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdFormat(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandFormat, "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdRmPwd(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandRmpwd, "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdSetPwd(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandSetpwd, "--" + models.FlagPwd, "AA AA AA AA", "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdTransmit(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandTransmit, "--" + models.FlagTarget, "adapter", "--" + models.FlagTxBytes, "AA AA AA AA", "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}

func Test_cmdWrite(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	os.Args = []string{"nfc-cli", models.CommandWrite, "--" + models.FlagNdefType, "url", "--" + models.NdefTypeUrl, "http://url.ulr", "--" + models.FlagExport, "--" + models.FlagOutput, "cmd_test_file.json"}
	err := app.Start()
	assert.Nil(t, err)
	err = os.Remove(app.output)
	assert.Nil(t, err)
}
