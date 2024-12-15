package integrations

import (
	"context"
	"flag"
	"fmt"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../configs/config-testing.toml", "Path to configuration file")
}

func SetupSuite() (conf *config.Config) {
	flag.Parse()
	conf = config.New(configFile)
	if conf == nil {
		panic("config file is invalid")
	}
	return
}

func PoolConnect(conf *config.Config) (pool *pgxpool.Pool) {
	db := conf.Database
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.Username, db.Password, db.DBName)
	pool, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		panic(err)
	}
	return
}

func StorageConnect(conf *config.Config) (store storage.Storage) {
	store = storage.New(conf)
	err := store.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	return
}
