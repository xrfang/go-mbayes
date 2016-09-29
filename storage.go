package mbayes

import (
	"database/sql"
)

type classifier struct {
	db *sql.DB
}

func Open(dsn string) (*classifier, error) {
	db, err := sql.Open(DBTYPE, dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(SQL("create"))
	if err != nil {
		return nil, err
	}
	return &classifier{db: db}, nil
}

func (cf *classifier) Close() (err error) {
	return cf.db.Close()
}

func (cf *classifier) add(category string, token []byte) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("addtok1"), token, category)
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("addtok2"), token, category)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (cf *classifier) delete(category string, token []byte) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("deltok1"), token, category)
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("deltok2"), token, category)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
