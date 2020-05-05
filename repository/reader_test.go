package repository

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/client"
)

func TestRepositoryService_readFromFile(t *testing.T) {
	nfc := client.New("url")
	rep := New(&nfc)

	data, err := rep.readFromFile("reader_test_file.json")
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "Write tag", data[0].JobName)
	assert.Equal(t, "Second job", data[1].JobName)
}
