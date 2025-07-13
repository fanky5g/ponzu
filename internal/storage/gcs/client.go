package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	internalStorage "github.com/fanky5g/ponzu/internal/storage"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
)

type Client struct {
	s *storage.Client
}

// Open serves directly a file. In the future, we will support generating signed links and redirecting,
// so we don't have to stream the file (at a performance and data streaming cost). The format for the name of the file is
// bucket/path/to/file.ext
func (c *Client) Open(name string) (http.File, error) {
	object, err := c.objectHandle(name)
	if err != nil {
		return nil, err
	}

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
	w.ObjectAttrs.ContentType = mime.TypeByExtension(path.Ext(name))

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

func (c *Client) Attributes(name string) (*internalStorage.FileAttributes, error) {
	object, err := c.objectHandle(name)
	if err != nil {
		return nil, err
	}

	objectAttrs, err := object.Attrs(context.Background())
	if err != nil {
		return nil, err
	}

	return &internalStorage.FileAttributes{ContentType: objectAttrs.ContentType}, nil
}

func (c *Client) Delete(name string) error {
	bucket, key, err := parsePath(name)
	if err != nil {
		return err
	}

	return c.s.Bucket(bucket).Object(key).Delete(context.Background())
}

func (c *Client) objectHandle(name string) (*storage.ObjectHandle, error) {
	bucket, key, err := parsePath(name)
	if err != nil {
		return nil, err
	}

	return c.s.Bucket(bucket).Object(key), nil
}
