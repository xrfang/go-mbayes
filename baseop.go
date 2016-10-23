package mbayes

import (
	"database/sql"
)

func (cf *Classifier) add(tr *sql.Tx, category string, token []byte) (err error) {
	_, err = tr.Exec(`INSERT OR IGNORE INTO tokens (tok,cat) VALUES (?,?)`,
		string(token), category)
	if err != nil {
		return
	}
	_, err = tr.Exec(`UPDATE tokens SET cnt=cnt+1 WHERE tok=? AND cat=?`,
		string(token), category)
	return
}

func (cf *Classifier) del(tr *sql.Tx, category string, token []byte) (err error) {
	_, err = tr.Exec(`UPDATE tokens SET cnt=cnt-1 WHERE tok=? AND cat=?`,
		string(token), category)
	if err != nil {
		return
	}
	_, err = tr.Exec(`DELETE FROM tokens WHERE tok=? AND cat=? AND cnt<=0`,
		string(token), category)
	return
}

func (cf *Classifier) refreshCats() error {
	row := cf.db.QueryRow(`SELECT COUNT(DISTINCT cat) FROM tokens`)
	return row.Scan(&cf.cats)
}
