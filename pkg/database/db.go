package database

import (
	"log"

	"github.com/NineKanokpol/Nine-shop-test/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnection(cfg config.IDbConfig) *sqlx.DB {
	//Connect DB
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
