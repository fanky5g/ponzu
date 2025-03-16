package connection

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
)

var (
	once sync.Once
	conn *pgxpool.Pool
)

func Get(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		var cfg *Config
		cfg, err = getConfig()
		if err != nil {
			err = fmt.Errorf("failed to get config: %v", err)
			return
		}

		dbURI := fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)

		config, err := pgxpool.ParseConfig(dbURI)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to parse database connection: %v\n", err)
			return
		}

		conn, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %v", err)
			return
		}
	})

	return conn, err
}
