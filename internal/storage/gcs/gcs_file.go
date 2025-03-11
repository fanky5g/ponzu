package gcs

import (
	"cloud.google.com/go/storage"
	"io/fs"
	"time"
)

type gcsFile struct {
	obj *storage.ObjectHandle
	*storage.Reader
}

type gcsFileInfo struct {
	storage.ReaderObjectAttrs
	obj *storage.ObjectHandle
}

func (g *gcsFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (g *gcsFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

func (g *gcsFile) Stat() (fs.FileInfo, error) {
	return &gcsFileInfo{
		obj:               g.obj,
		ReaderObjectAttrs: g.Attrs,
	}, nil
}

func (g gcsFileInfo) Name() string {
	return g.obj.ObjectName()
}

func (g gcsFileInfo) Size() int64 {
	return g.ReaderObjectAttrs.Size
}

func (g gcsFileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

func (g gcsFileInfo) ModTime() time.Time {
	return g.ReaderObjectAttrs.LastModified
}

func (g gcsFileInfo) IsDir() bool {
	return false
}

func (g gcsFileInfo) Sys() any {
	return nil
}
