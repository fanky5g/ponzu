package localstorage

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/storage"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (c *Client) Save(fileName string, file io.ReadCloser) (string, int64, error) {
	defer func(f io.ReadCloser) {
		err := f.Close()
		if err != nil {
			log.Printf("Failed to close file: %v\n", err)
		}
	}(file)

	ts := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
	tm := time.Unix(ts/1000, ts%1000)

	uploadDir := filepath.Join(c.storageDir, fmt.Sprintf("%d", tm.Year()), fmt.Sprintf("%02d", tm.Month()))
	err := os.MkdirAll(uploadDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return "", 0, err
	}

	// check if file at path exists, if so, add timestamp to file
	absPath := filepath.Join(uploadDir, fileName)
	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		fileName = fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)
		absPath = filepath.Join(uploadDir, fileName)
	}

	// save to disk (TODO: or check if S3 credentials exist, & save to cloud)
	dst, err := os.Create(absPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create destination file for storage: %s", err)
	}

	// copy file from src to dst on disk
	var size int64
	if size, err = io.Copy(dst, file); err != nil {
		return "", 0, fmt.Errorf("failed to copy uploaded file to destination: %s", err)
	}

	// add name:urlPath to req.PostForm to be inserted into db
	urlPath := fmt.Sprintf("%d/%02d/%s", tm.Year(), tm.Month(), fileName)
	return urlPath, size, nil
}

func (c *Client) Open(name string) (http.File, error) {
	return c.fs.Open(name)
}

func (c *Client) Attributes(name string) (*storage.FileAttributes, error) {
	if _, err := c.fs.Open(name); err != nil {
		return nil, err
	}

	return &storage.FileAttributes{ContentType: mime.TypeByExtension(filepath.Ext(name))}, nil
}

func (c *Client) Delete(path string) error {
	// split and rebuild path in OS friendly way
	// use path to delete the physical file from disk
	pathSplit := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	return os.Remove(filepath.Join(pathSplit...))
}
