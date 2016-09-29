package mbayes

import (
	"database/sql"
)

type Classifier struct {
	db *sql.DB
}

func Open(dsn string) (*Classifier, error) {
	db, err := sql.Open(DBTYPE, dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(SQL("create"))
	if err != nil {
		return nil, err
	}
	return &Classifier{db: db}, nil
}

func (cf *Classifier) Close() (err error) {
	return cf.db.Close()
}
