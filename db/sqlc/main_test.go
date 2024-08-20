package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/tuda4/mb-backend/configs"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load file config:: %v", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect db:: %v", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
