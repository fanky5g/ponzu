package exceptions

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

type ExceptionSeverity string

const (
	INFO    ExceptionSeverity = "INFO"
	DEBUG   ExceptionSeverity = "DEBUG"
	WARNING ExceptionSeverity = "WARNING"
	ERROR   ExceptionSeverity = "ERROR"
)

type Exception struct {
	Message  string
	Cause    error
	Severity ExceptionSeverity
	Context  map[string]interface{}
}

// By default all exceptions will not be propagated to external clients. Eg. over http.
// To explicitly propagate exceptions to client, use ClientException(s)
type ClientException struct {
	Exception
	Code string
}

func (exception Exception) Error() string {
	return exception.Message
}

func New(message string, severity ExceptionSeverity, cause error, context map[string]interface{}) *Exception {
	return &Exception{
		Message:  message,
		Severity: severity,
		Cause:    cause,
		Context:  context,
	}
}

func NewClientException(
	message,
	code string,
	severity ExceptionSeverity,
	cause error,
	context map[string]interface{},
) *ClientException {
	return &ClientException{
		Code:      code,
		Exception: *New(message, severity, cause, context),
	}
}

func logException(exception *Exception) {
	localLogger := log.New()
	localLogger.SetLevel(getLogLevel(exception.Severity))
	localLogger.SetFormatter(&log.JSONFormatter{})

	logEntry := localLogger.WithContext(context.TODO())
	if exception.Cause != nil {
		logEntry = logEntry.WithError(exception.Cause)
	}

	if exception.Context != nil {
		logEntry = logEntry.WithFields(log.Fields(exception.Context))
	}

	logEntry.Log(localLogger.Level, exception.Message)
}

func Log(err error) {
	if err == nil {
		return
	}

	var exception *Exception
	var clientException *ClientException
	if errors.As(err, &exception) {
		logException(err.(*Exception))
		return
	}

	if errors.As(err, &clientException) {
		logException(&err.(*ClientException).Exception)
		return
	}

	log.WithError(err).Error(err.Error())
}

func getLogLevel(severity ExceptionSeverity) log.Level {
	switch severity {
	case INFO:
		return log.InfoLevel
	case ERROR:
		return log.ErrorLevel
	case WARNING:
		return log.WarnLevel
	default:
		return log.DebugLevel
	}
}
