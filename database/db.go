package database

import (
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu/config"
	analyticsRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/analytics"
	configRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/config"
	contentRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/content"
	credentialRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/credential"
	recoveryKeyRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/recovery-key"
	uploadRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/uploads"
	usersRepositoryFactory "github.com/fanky5g/ponzu/driver/db/repository/users"
	"log"
	"path/filepath"
)

type database struct {
	store        *bolt.DB
	repositories *Repositories
}

func (db *database) GetRepositories() (*Repositories, error) {
	return db.repositories, nil
}

func (db *database) Close() error {
	return db.store.Close()
}

func New() (Database, error) {
	store, err := bolt.Open(filepath.Join(config.DataDir(), "system.db"), 0666, nil)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	configRepository, err := configRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize config repository %v", err)
	}

	analyticsRepository, err := analyticsRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize analytics repository: %v", err)
	}

	userRepository, err := usersRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	contentRepository, err := contentRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize content repository: %v", err)
	}

	credentialRepository, err := credentialRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize credential repository: %v", err)
	}

	recoveryKeyRepository, err := recoveryKeyRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize recovery key repository: %v", err)
	}

	uploadRepository, err := uploadRepositoryFactory.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize upload repository: %v", err)
	}
	// End initialize repositories

	repositories := &Repositories{
		Analytics:        analyticsRepository,
		Config:           configRepository,
		Users:            userRepository,
		Content:          contentRepository,
		CredentialHashes: credentialRepository,
		RecoveryKeys:     recoveryKeyRepository,
		Uploads:          uploadRepository,
	}

	return &database{store: store, repositories: repositories}, nil
}
