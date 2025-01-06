package storage

import (
	"fmt"

	"github.com/ElladanTasartir/ledger-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

type Storage interface {
}

type storage struct {
	db *sqlx.DB
}

func NewStorageModule() fx.Option {
	return fx.Provide(
		NewStorage,
	)
}

func NewStorage(config *config.Config) (Storage, error) {
	db, err := sqlx.Connect("postgres", getDSN(config.DB))
	if err != nil {
		return nil, err
	}

	return storage{
		db: db,
	}, nil
}

func getDSN(config *config.DBConfig) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", config.User, config.Password, config.Database, config.Host, config.Port)
}
