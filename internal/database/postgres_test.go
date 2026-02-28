package database

import (
	"os"
	"testing"
)

func TestNewPostgresConnection(t *testing.T) {
	os.Setenv("DATABASE_HOST", "localhost")
	_, _ = NewPostgresConnection()
}
