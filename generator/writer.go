package generator

type Writer interface {
	Write(filePath string, buf []byte) error
}
