package mbayes

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
)

type classifier struct {
	digest bool
	db     *sql.DB
}

func Open(dsn string, digest bool) (*classifier, error) {
	db, err := sql.Open(DBTYPE, dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(SQL("create"))
	if err != nil {
		return nil, err
	}
	return &classifier{db: db, digest: digest}, nil
}

func (cf *classifier) Close() (err error) {
	return cf.db.Close()
}

func (cf *classifier) stringy(token []byte) string {
	if cf.digest {
		sum := sha1.Sum(token)
		return hex.EncodeToString(sum[:])
	}
	return hex.EncodeToString(token)
}

func (cf *classifier) add(token []byte, category string) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	tok := cf.stringy(token)
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

func (cf *classifier) delete(token []byte, category string) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	tok := cf.stringy(token)
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
