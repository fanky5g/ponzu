//go:build integration

package gcs

import (
	"cloud.google.com/go/storage"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"testing"
)

type StorageTestSuite struct {
	suite.Suite
	c      *Client
	bucket string
}

func (s *StorageTestSuite) SetupSuite() {
	var err error
	s.c, err = New()
	if err != nil {
		s.T().Fatal(err)
	}

	s.bucket = "bucket-34ff28a4-b823-4c05-8367-8d760927938b"
}

func (s *StorageTestSuite) writeTestFile() (string, int64, error) {
	f, err := os.Open("testdata/test.txt")
	if err != nil {
		return "", 0, err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.WithField("Error", err).Error("Error closing opened file")
		}
	}()

	name := fmt.Sprintf("%s/test.txt", s.bucket)
	var written int64
	_, written, err = s.c.Save(name, f)
	if err != nil {
		return "", 0, err
	}

	return name, written, nil
}

func (s *StorageTestSuite) TestOpenFile() {
	name, size, err := s.writeTestFile()
	if err != nil {
		s.T().Fatal(err)
		return
	}

	outFileName := "testdata/out.txt"
	dest, err := os.Create(outFileName)
	if err != nil {
		s.T().Fatal(err)
		return
	}

	defer func() {
		if err = dest.Close(); err != nil {
			log.WithField("Error", err).Error("Error closing destination file")
		}
	}()

	f, err := s.c.Open(name)
	if err != nil {
		s.T().Fatal(err)
	}

	var written int64
	written, err = io.Copy(dest, f)
	if assert.NoError(s.T(), err) {
		assert.FileExists(s.T(), outFileName)
		assert.Equal(s.T(), written, size)
	}
}

func (s *StorageTestSuite) TestSave() {
	f, err := os.Open("testdata/test.txt")
	if err != nil {
		s.T().Fatal(err)
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.WithField("Error", err).Error("Error closing opened file")
		}
	}()

	name := fmt.Sprintf("%s/test.txt", s.bucket)
	fileName, written, err := s.c.Save(name, f)
	if assert.NoError(s.T(), err) {
		assert.True(s.T(), written > 0)
		assert.Equal(s.T(), fileName, name)
	}

	var existingFile http.File
	var fileStat os.FileInfo
	existingFile, err = s.c.Open(name)
	if assert.NoError(s.T(), err) {
		defer func() {
			if err = existingFile.Close(); err != nil {
				log.WithField("Error", err).Error("Error closing file")
			}
		}()

		fileStat, err = existingFile.Stat()
		if assert.NoError(s.T(), err) {
			assert.Equal(s.T(), written, fileStat.Size())
		}
	}
}

func (s *StorageTestSuite) TestDelete() {
	name, _, err := s.writeTestFile()
	if err != nil {
		s.T().Fatal(err)
		return
	}

	if assert.NoError(s.T(), s.c.Delete(name)) {
		_, err = s.c.Open(name)
		if assert.Error(s.T(), err) {
			assert.Equal(s.T(), err, storage.ErrObjectNotExist)
		}
	}
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
