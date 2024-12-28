package repository

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("DATABASE_HOST", "localhost")
	_ = os.Setenv("DATABASE_USER", "postgres")
	_ = os.Setenv("DATABASE_PASSWORD", "password")
	_ = os.Setenv("DATABASE_NAME", "ponzu")
	_ = os.Setenv("DATABASE_PORT", "5432")

	m.Run()
}
