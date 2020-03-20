package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/client"
	"log"
	"testing"
)

func TestRepositoryService_readFromFile(t *testing.T) {
	nfc := client.New("url", "en")
	rep := New(&nfc)

	data, err := rep.readFromFile("reader_test_file.json")
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "Write tag", data[0].JobName)
	assert.Equal(t, "Second job", data[1].JobName)
}
