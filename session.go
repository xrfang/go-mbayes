package mbayes

import (
	"database/sql"
)

type Classifier struct {
	db  *sql.DB
	err error
	tx  *sql.Tx
}

type SessionError struct {
	msg string
}

func (se SessionError) Error() string {
	return se.msg
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

func (cf *Classifier) Begin() (err error) {
	if cf.tx != nil {
		return SessionError{msg: "already in transaction"}
	}
	cf.tx, err = cf.db.Begin()
	return
}

func (cf *Classifier) Commit() (err error) {
	if cf.tx == nil {
		return SessionError{msg: "not in transaction"}
	}
	err = cf.tx.Commit()
	cf.tx = nil
	cf.err = nil
	return
}

func (cf *Classifier) Rollback() (err error) {
	if cf.tx == nil {
		return SessionError{msg: "not in transaction"}
	}
	err = cf.tx.Rollback()
	cf.tx = nil
	cf.err = nil
	return
}

func (cf *Classifier) Err() error {
	return cf.err
}
