package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"testing"
)

func TestNew(t *testing.T) {
	var rep *repository.RepositoryService
	config := opts.Config{}
	cbCliStarted := func(string) {}

	app := New(rep, cbCliStarted, config)
	assert.NotNil(t, app)
}

func TestAppService_getCommands(t *testing.T) {
	var rep *repository.RepositoryService
	config := opts.Config{}
	cbCliStarted := func(string) {}

	app := New(rep, cbCliStarted, config)
	commands := app.getCommands()

	for _, f := range commands {
		assert.NotEmpty(t, f.Usage)
	}
}
