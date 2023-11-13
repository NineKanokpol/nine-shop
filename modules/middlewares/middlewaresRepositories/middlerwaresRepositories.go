package middlewaresRepositories

import "github.com/jmoiron/sqlx"

//interface,struct,constructor

type IMiddlewaresRepository interface {
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}
