package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyb88999/Go-SimpleBank/util"
	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
