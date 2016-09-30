package mbayes

import (
	"database/sql"
)

func (cf *Classifier) add(tr *sql.Tx, category string, token []byte) (err error) {
	_, err = tr.Exec(SQL("addtok1"), token, category)
	if err != nil {
		return
	}
	_, err = tr.Exec(SQL("addtok2"), token, category)
	return
}

func (cf *Classifier) del(tr *sql.Tx, category string, token []byte) (err error) {
	_, err = tr.Exec(SQL("deltok1"), token, category)
	if err != nil {
		return
	}
	_, err = tr.Exec(SQL("deltok2"), token, category)
	return
}
