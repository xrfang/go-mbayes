package mbayes

import (
	"database/sql"
)

func (cf *Classifier) Train(category string, tokens ...[]byte) (err error) {
	if cf.err != nil {
		return
	}
	autoCommit := cf.tx == nil
	var tx *sql.Tx
	if autoCommit {
		tx, err = cf.db.Begin()
		if err != nil {
			return
		}
	} else {
		tx = cf.tx
	}
	for _, tk := range tokens {
		_, err = tx.Exec(SQL("addtok1"), tk, category)
		if err != nil {
			break
		}
		_, err = tx.Exec(SQL("addtok2"), tk, category)
		if err != nil {
			break
		}
	}
	if autoCommit {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	} else {
		cf.err = err
	}
	return
}

func (cf *Classifier) Untrain(category string, tokens ...[]byte) (err error) {
	if cf.err != nil {
		return
	}
	autoCommit := cf.tx == nil
	var tx *sql.Tx
	if autoCommit {
		tx, err = cf.db.Begin()
		if err != nil {
			return
		}
	} else {
		tx = cf.tx
	}
	for _, tk := range tokens {
		_, err = tx.Exec(SQL("deltok1"), tk, category)
		if err != nil {
			break
		}
		_, err = tx.Exec(SQL("deltok2"), tk, category)
		if err != nil {
			break
		}
	}
	if autoCommit {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	} else {
		cf.err = err
	}
	return
}
