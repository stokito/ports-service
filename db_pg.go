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
INSERT INTO ports (unlocs, port)
VALUES ($1, $2)
`

// language=PostgreSQL
var sqlDeletePort = `
DELETE FROM ports WHERE $1 = ANY(unlocs)
`

// language=PostgreSQL
var sqlGetAll = "SELECT port FROM ports"

// language=PostgreSQL
var sqlGetPort = `
SELECT port
FROM ports 
WHERE $1 = ANY(unlocs)
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

// UpsertPort Insert or update Port.
// Note: it's not CTE atomic so don't use concurrently
func (db *PostgresDb) UpsertPort(ctx context.Context, portUnloc string, port *Port) {
	db.RemovePort(ctx, portUnloc)
	sort.Strings(port.Unlocs)
	_, sqlErr := db.pool.Exec(ctx, sqlUpsertPort,
		port.Unlocs, port)
	if sqlErr != nil {
		log.Printf("WARN Fail to upsert port %v\n", sqlErr)
	}
}

func (db *PostgresDb) FindPort(ctx context.Context, portUnloc string) *Port {
	rows, err := db.pool.Query(ctx, sqlGetPort, portUnloc)
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
	rows, err := db.pool.Query(ctx, sqlGetAll)
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
	var portJson string
	scanErr := rows.Scan(&portJson)
	if scanErr != nil {
		log.Printf("ERR scan error: %s\n", scanErr)
		return nil
	}
	// the port with all fields is stored as is into a column so we need to unmarshal it
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

func (db *PostgresDb) RemovePort(ctx context.Context, portUnloc string) {
	_, err := db.pool.Exec(ctx, sqlDeletePort, portUnloc)
	if err != nil {
		log.Printf("ERR RemovePort error: %s\n", err)
		return
	}
}
