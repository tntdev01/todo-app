package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {

	if os.Getenv("TEST_DATABASE_URL") == "" {
		_ = godotenv.Load("../../.env.test")
	}

	connString := os.Getenv("TEST_DATABASE_URL")

	var err error
	testPool, err = pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		panic(err)
	}

	err = testPool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	code := m.Run()

	testPool.Close()

	os.Exit(code)
}
