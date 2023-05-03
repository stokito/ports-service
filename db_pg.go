package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// language=PostgreSQL
var sqlUpsertPort = `INSERT INTO ports (
	unlocks, port)
	VALUES ($1, $2)
`

type PostgresDb struct {
	pool        *pgxpool.Pool
	databaseUrl string
}

func NewPostgresDb(databaseUrl string) *PostgresDb {
	return &PostgresDb{databaseUrl: databaseUrl}
}

//Connect(ctx context.Context) error
//Close(ctx context.Context)
//UpsertPort(ctx context.Context, portUnlock string, port *Port) error
//FindPort(ctx context.Context, portUnlock string) *Port

func (db *PostgresDb) Connect(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(db.databaseUrl)
	if err != nil {
		return err
	}
	pool, dbErr := pgxpool.ConnectConfig(ctx, poolConfig)
	if dbErr != nil {
		return err
	}
	db.pool = pool
	log.Printf("INFO Connected to database\n")
	return nil
}

func (db *PostgresDb) Close() {
	if db != nil {
		db.pool.Close()
		db.pool = nil
		log.Printf("INFO DB disconnected\n")
	}
}

func (db *PostgresDb) UpsertPort(ctx context.Context, portUnlock string, port *Port) {
	_, sqlErr := db.pool.Exec(ctx, sqlUpsertPort,
		portUnlock, port)
	if sqlErr != nil {
		log.Printf("WARN Fail to upsert port %v\n", sqlErr)
	}
}

func (db *PostgresDb) FindPort(ctx context.Context, portUnlock string) *Port {
	//TODO
	return nil
}

func (db *PostgresDb) GetAll(ctx context.Context) []*Port {
	//TODO
	return nil
}
