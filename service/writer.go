package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func (s *appService) writeToFile(filename string, data interface{}) (err error) {
	var file *os.File
	if s.append || (s.ongoingJobs.published > 1 && s.ongoingJobs.published-s.ongoingJobs.left > 1) {
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return errors.Wrap(err, "Can't open the file: ")
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return errors.Wrap(err, "Can't encode the data on writing to the file: ")
	}

	return nil
}
