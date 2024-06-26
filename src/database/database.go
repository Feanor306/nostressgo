package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/feanor306/nostressgo/src/types"
	"github.com/feanor306/nostressgo/src/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nbd-wtf/go-nostr"
)

const EVENTS = "events"

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
			subject text,
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
	etags := utils.GetEtags(event)
	ptags := utils.GetPtags(event)
	gtags := utils.GetGtags(event)
	dtag := utils.GetDtag(event)
	subject := utils.GetSubject(event)

	expiration, err := utils.GetExpiration(event)
	if err != nil {
		return err
	}

	json, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	var id string
	err = db.sq.Insert(EVENTS).
		Columns("id", "pubkey", "content", "subject", "created_at", "kind", "etags", "ptags", "gtags", "dtag", "expiration", "raw").
		Values(event.ID, event.PubKey, event.Content, subject, event.CreatedAt.Time().Unix(), event.Kind, etags, ptags, gtags, dtag, expiration, string(json)).
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
		From(EVENTS).
		Where(squirrel.Eq{"pubkey": event.PubKey}).
		Where(squirrel.Eq{"kind": nostr.KindProfileMetadata}).
		QueryRow().Scan(&id)

	if err == sql.ErrNoRows {
		return id, nil
	}

	return id, err
}

func (db *DB) UpdateEventZero(id string, event *nostr.Event) error {
	_, err := db.sq.Update(EVENTS).
		Set("content", event.Content).
		Set("created_at", event.CreatedAt).
		Where(squirrel.Eq{"id": id}).
		Exec()

	return err
}

func (db *DB) ExpireEvents(etags []string) error {
	_, err := db.sq.Update(EVENTS).
		Set("expiration", time.Now().Unix()).
		Where(squirrel.Eq{"id": etags}).
		Exec()
	return err
}

func (db *DB) GetEventsByFilter(filter *nostr.Filter, chanGroup *types.ChanGroup) error {
	defer chanGroup.Done()
	query := db.sq.
		Select("id", "pubkey", "content", "subject", "created_at", "kind", "array_to_json(etags)", "array_to_json(ptags)", "array_to_json(gtags)", "dtag", "expiration", "raw").
		From(EVENTS)

	query = BuildFilterQuery(filter, query)

	rows, err := query.Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var event types.Event
		var etags, ptags, gtags string

		if err := rows.Scan(&event.ID, &event.PubKey, &event.Content, &event.Subject, &event.CreatedAt,
			&event.Kind, &etags, &ptags, &gtags, &event.Dtag, &event.Expiration, &event.Json); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(etags), &event.Etags); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(ptags), &event.Ptags); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(gtags), &event.Gtags); err != nil {
			return err
		}

		event.SetTags()
		chanGroup.Chan <- event.ToEnvelopeWrapper()
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
