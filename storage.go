package mbayes

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
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

func (cf *classifier) Add(token []byte, category string) (err error) {
	sum := sha1.Sum(token)
	tok := hex.EncodeToString(sum[:])
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("addtok1"), tok, category)
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("addtok2"), tok, category)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (cf *classifier) Delete(token []byte, category string) (err error) {
	sum := sha1.Sum(token)
	tok := hex.EncodeToString(sum[:])
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("deltok1"), tok, category)
	if err != nil {
		return
	}
	_, err = tx.Exec(SQL("deltok2"), tok, category)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
