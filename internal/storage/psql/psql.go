package psql

import (
	"context"
	"fmt"

	"log"

	"github.com/bulgil/blog-rest-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewStorage(storageCfg config.PGStorage) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/blog",
		storageCfg.Login, storageCfg.Password, storageCfg.Address, storageCfg.Port)
	conpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id serial PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		tags TEXT[] NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
	`

	_, err = conpool.Exec(context.Background(), stmt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("storage initialized")
	return conpool
}
