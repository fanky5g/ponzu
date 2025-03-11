package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	s *storage.Client
}

// Open serves directly a file. In the future, we will support generating signed links and redirecting,
// so we don't have to stream the file (at a performance and data streaming cost). The format for the name of the file is
// bucket/path/to/file.ext
func (c *Client) Open(name string) (http.File, error) {
	bucket, key, err := parsePath(name)
	if err != nil {
		return nil, err
	}

	object := c.s.Bucket(bucket).Object(key)
	o, err := object.NewReader(context.Background())
	if err != nil {
		return nil, err
	}

	return &gcsFile{
		obj:    object,
		Reader: o,
	}, nil
}

// Save saves a file to storage. Argument name must be of the syntax bucket/path/to/file.ext
func (c *Client) Save(name string, file io.ReadCloser) (string, int64, error) {
	bucket, key, err := parsePath(name)
	if err != nil {
		return "", 0, err
	}

	object := c.s.Bucket(bucket).Object(key)
	w := object.NewWriter(context.Background())
	defer func() {
		if err = w.Close(); err != nil {
			log.WithField("Error", err).Error("Error closing writer")
		}
	}()

	written, err := io.Copy(w, file)
	if err != nil {
		return "", written, err
	}

	return strings.Join([]string{bucket, key}, "/"), written, nil
}

func (c *Client) Delete(name string) error {
	bucket, key, err := parsePath(name)
	if err != nil {
		return err
	}

	return c.s.Bucket(bucket).Object(key).Delete(context.Background())
}
