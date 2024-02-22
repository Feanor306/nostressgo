package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	conn *pgxpool.Pool
}

func NewDatabase(context context.Context, connString string) (*DB, error) {
	dbpool, err := pgxpool.New(context, connString)

	if err != nil {
		return nil, err
	}

	return &DB{
		conn: dbpool,
	}, nil
}

func (db *DB) InitDatabase(context context.Context) (string, error) {
	a, err := db.conn.Exec(context, `
	CREATE TABLE IF NOT EXISTS events (
		id text NOT NULL,
		pubkey text NOT NULL,
		created_at integer NOT NULL,
		kind integer NOT NULL,
		etags text[],
		ptags text[],
		dtag text,
		expiration integer,
		gtags text[],
		raw json
	  );
	`)

	if err != nil {
		return "", err
	}
	return a.String(), nil
}
