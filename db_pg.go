package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var db *pgxpool.Pool

// language=PostgreSQL
var sqlUpsertPort = `INSERT INTO ports (
	unlocks, port)
	VALUES ($1, $2)
`

func dbConnect(ctx context.Context) error {
	if conf.DatabaseUrl == "" {
		return errors.New("database is not configured")
	}
	poolConfig, err := pgxpool.ParseConfig(conf.DatabaseUrl)
	if err != nil {
		return err
	}
	var dbErr error
	db, dbErr = pgxpool.ConnectConfig(ctx, poolConfig)
	if dbErr != nil {
		return err
	}
	log.Printf("INFO Connected to database\n")
	return nil
}

func dbClose() {
	if db != nil {
		db.Close()
		db = nil
		log.Printf("INFO DB disconnected\n")
	}
}

func UpsertPort(ctx context.Context, unlock string, port *Port) {
	_, sqlErr := db.Exec(ctx, sqlUpsertPort,
		unlock, port)
	if sqlErr != nil {
		log.Printf("WARN Fail to upsert port %v\n", sqlErr)
	}
}
