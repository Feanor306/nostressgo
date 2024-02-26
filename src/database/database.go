package database

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/feanor306/nostressgo/src/helpers"
	"github.com/nbd-wtf/go-nostr"
)

type DB struct {
	conn *sql.DB
	sq   squirrel.StatementBuilderType
}

func NewDatabase(connString string) (*DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn: db,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}, nil
}

func (db *DB) Close() {
	db.conn.Close()
}

func (db *DB) InitDatabase() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id text NOT NULL,
			pubkey text NOT NULL,
			content text NOT NULL,
			created_at integer NOT NULL,
			kind integer NOT NULL,
			etags text[],
			ptags text[],
			gtags text[],
			dtag text,
			expiration integer,
			raw json
		);

		CREATE UNIQUE INDEX IF NOT EXISTS events_id_index ON events USING btree (id text_pattern_ops);
		CREATE INDEX IF NOT EXISTS events_pubkey_index ON events USING btree (pubkey text_pattern_ops);
		CREATE INDEX IF NOT EXISTS events_created_at_index ON events (created_at DESC);
		CREATE INDEX IF NOT EXISTS events_kind_index ON events (kind);
		CREATE INDEX IF NOT EXISTS events_ptags_index ON events USING gin (etags);
		CREATE INDEX IF NOT EXISTS events_etags_index ON events USING gin (ptags);
		CREATE INDEX IF NOT EXISTS events_gtags_index ON events USING gin (gtags);
		CREATE INDEX IF NOT EXISTS events_expiration_index ON events (expiration DESC);
	`)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateEvent(event *nostr.Event) error {
	etags := helpers.GetEtags(event)
	ptags := helpers.GetPtags(event)
	gtags := helpers.GetGtags(event)
	dtag := event.Tags.GetD()

	expiration, err := helpers.GetExpiration(event)
	if err != nil {
		return err
	}

	json, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	var id string
	err = db.sq.Insert("events").
		Columns("id", "pubkey", "content", "created_at", "kind", "etags", "ptags", "gtags", "dtag", "expiration", "raw").
		Values(event.ID, event.PubKey, event.Content, event.CreatedAt.Time().Unix(), event.Kind, etags, ptags, gtags, dtag, expiration, string(json)).
		Suffix("RETURNING \"id\"").
		QueryRow().Scan(&id)

	if len(id) == 0 {
		return fmt.Errorf("create event failed")
	}

	return err
}

func (db *DB) EventZeroExists(event *nostr.Event) (string, error) {
	var id string
	err := db.sq.Select("id").
		From("events").
		Where("pubkey", event.PubKey).
		Where("kind", "0").
		QueryRow().Scan(&id)

	return id, err
}

func (db *DB) UpdateEventZero(id string, event *nostr.Event) error {
	json, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = db.sq.Update("events").
		Set("id", event.ID).
		Set("content", event.Content).
		Set("created_at", event.CreatedAt).
		Set("raw", string(json)).
		Where("id", id).
		Exec()

	return err
}
