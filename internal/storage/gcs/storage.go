package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
)

func GetStorageClient(ctx context.Context) (*storage.Client, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	opts := make([]option.ClientOption, 0)
	if cfg.ServiceAccountFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.ServiceAccountFile))
	}

	return storage.NewClient(ctx, opts...)
}

func New() (*Client, error) {
	s, err := GetStorageClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Client{s: s}, nil
}
