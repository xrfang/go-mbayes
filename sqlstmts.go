package mbayes

import _ "github.com/mattn/go-sqlite3"

const DBTYPE = "sqlite3"

var stmts map[string]map[string]string

func init() {
	stmts = make(map[string]map[string]string)
	stmts["sqlite3"] = map[string]string{
		"create": `CREATE TABLE IF NOT EXISTS tokens (tok TEXT NOT NULL, cat TEXT
		    NOT NULL, cnt INTEGER DEFAULT (0), PRIMARY KEY (tok,cat))`,
		"addtok1": `INSERT OR IGNORE INTO tokens (tok,cat) VALUES (?,?)`,
		"addtok2": `UPDATE tokens SET cnt=cnt+1 WHERE tok=? AND cat=?`,
		"deltok1": `UPDATE tokens SET cnt=cnt-1 WHERE tok=? AND cat=?`,
		"deltok2": `DELETE FROM tokens WHERE tok=? AND cat=? AND cnt<=0`,
	}
}

func SQL(tag string) string {
	return stmts[DBTYPE][tag]
}
