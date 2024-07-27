package infrastructure

import (
	"fmt"
	"github.com/pkg/errors"

	gcsStorage "github.com/fanky5g/ponzu-driver-gcs"
	localStorage "github.com/fanky5g/ponzu-driver-local-storage"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/driver"
)

var (
	ErrStorageDriverMissing     = errors.New("Invalid configuration: missing upload storage_driver")
	ErrUnsupportedStorageDriver = errors.New("Unsupported upload storage driver")
)

func getUploadStorageClient() (driver.StorageClientInterface, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	if cfg.StorageDriver == "" {
		return nil, ErrStorageDriverMissing
	}

	switch cfg.StorageDriver {
	case "local":
		uploadStorageClient, err := localStorage.New("")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage client: %v", err)
		}

		return uploadStorageClient, nil
	case "gcs":
		gcsStorageClient, err := gcsStorage.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize gcs storage driver")
		}

		return gcsStorageClient, nil
	default:
		return nil, ErrUnsupportedStorageDriver
	}
}
