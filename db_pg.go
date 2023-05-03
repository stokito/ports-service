package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sort"
)

// language=PostgreSQL
var sqlUpsertPort = `
INSERT INTO ports (
unlocs, port, ids)
VALUES ($1, $2, $3)
`

// language=PostgreSQL
var sqlGetPort = `
SELECT unlocs, port, ids
FROM ports 
WHERE $1 = ANY(ids)
`

type PostgresDb struct {
	pool        *pgxpool.Pool
	databaseUrl string
}

func NewPostgresDb(databaseUrl string) *PostgresDb {
	return &PostgresDb{databaseUrl: databaseUrl}
}

func (db *PostgresDb) Connect(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(db.databaseUrl)
	if err != nil {
		return err
	}
	pool, dbErr := pgxpool.ConnectConfig(ctx, poolConfig)
	if dbErr != nil {
		return dbErr
	}
	db.pool = pool
	log.Printf("INFO Connected to database\n")
	return nil
}

func (db *PostgresDb) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.pool = nil
		log.Printf("INFO DB disconnected\n")
	}
}

func (db *PostgresDb) UpsertPort(ctx context.Context, newUnloc string, port *Port) {
	sort.Strings(port.Unlocs)
	_, sqlErr := db.pool.Exec(ctx, sqlUpsertPort,
		newUnloc, port, port.Unlocs)
	if sqlErr != nil {
		log.Printf("WARN Fail to upsert port %v\n", sqlErr)
	}
}

func (db *PostgresDb) FindPort(ctx context.Context, newUnloc string) *Port {
	rows, err := db.pool.Query(ctx, sqlGetPort, newUnloc)
	if err != nil {
		log.Printf("ERR GetAll error: %s\n", err)
		return nil
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	port := db.scanRow(rows)
	return port
}

func (db *PostgresDb) GetAll(ctx context.Context) []*Port {
	rows, err := db.pool.Query(ctx, "SELECT unlocs, port, ids FROM ports")
	if err != nil {
		log.Printf("ERR GetAll error: %s\n", err)
		return []*Port{}
	}
	defer rows.Close()
	allPorts := []*Port{}
	for rows.Next() {
		port := db.scanRow(rows)
		if port == nil {
			continue
		}
		allPorts = append(allPorts, port)
	}
	return allPorts
}

func (db *PostgresDb) scanRow(rows pgx.Rows) *Port {
	var unlocs string
	var portJson string
	var ids []string
	scanErr := rows.Scan(&unlocs, &portJson, &ids)
	if scanErr != nil {
		log.Printf("ERR scan error: %s\n", scanErr)
		return nil
	}
	port := &Port{}
	jsonErr := json.Unmarshal([]byte(portJson), port)
	if jsonErr != nil {
		log.Printf("ERR scan error: %s\n", jsonErr)
		return nil
	}
	return port
}

func (db *PostgresDb) RemoveAll() {
	ctx := context.Background()
	_, err := db.pool.Exec(ctx, "DELETE FROM ports")
	if err != nil {
		log.Printf("ERR GetAll error: %s\n", err)
		return
	}
}
