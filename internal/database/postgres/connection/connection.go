package connection

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)

		conn, err = pgxpool.New(ctx, dsn)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %v", err)
			return
		}
	})

	return conn, err
}
