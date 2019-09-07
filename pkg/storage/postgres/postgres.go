package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/qwwqe/VocabKing/pkg/errors"
)

type Postgres struct {
	Source string
	conn   *sql.DB
}

func (db *Postgres) Open() error {
	const op errors.Op = "postgres.Open"

	conn, err := sql.Open("postgres", db.Source)
	if err != nil {
		return errors.New(op, errors.KindDatabaseFailure, err)
	}

	db.conn = conn
	return nil
}

func (db *Postgres) Close() error {
	const op errors.Op = "postgres.Close"

	return errors.NilOrNew(
		op,
		errors.KindDatabaseFailure,
		db.conn.Close(),
	)
}
