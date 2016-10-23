package mbayes

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TA_TRAIN = iota
	TA_UNTRAIN
	TA_COMMIT
	TA_ROLLBACK
)

type trainingSample struct {
	action  int
	feature []byte
	label   string
}

type Classifier struct {
	db   *sql.DB
	err  error
	tx   *sql.Tx
	cats int
	que  chan trainingSample
	wg   sync.WaitGroup
}

type SessionError struct {
	msg string
}

func (se SessionError) Error() string {
	return se.msg
}

func Open(dsn string) (*Classifier, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (tok TEXT NOT NULL, cat
	    TEXT NOT NULL, cnt INTEGER DEFAULT (0), PRIMARY KEY (tok,cat))`)
	if err != nil {
		return nil, err
	}
	c := &Classifier{db: db}
	err = c.refreshCats()
	return c, err
}

func (cf *Classifier) Close() (err error) {
	return cf.db.Close()
}

func (cf *Classifier) Begin() (err error) {
	if cf.tx != nil {
		return SessionError{msg: "already in transaction"}
	}
	cf.tx, err = cf.db.Begin()
	if err != nil {
		return
	}
	cf.err = nil
	cf.que = make(chan trainingSample, 128)
	cf.wg.Add(1)
	go func() {
		for {
			if cf.err != nil {
				break
			}
			ts := <-cf.que
			switch ts.action {
			case TA_TRAIN:
				cf.err = cf.add(cf.tx, ts.label, ts.feature)
			case TA_UNTRAIN:
				cf.err = cf.del(cf.tx, ts.label, ts.feature)
			case TA_COMMIT:
				cf.err = cf.tx.Commit()
				cf.tx = nil
				cf.wg.Done()
				break
			case TA_ROLLBACK:
				cf.err = cf.tx.Rollback()
				cf.tx = nil
				cf.wg.Done()
				break
			}
		}
	}()
	return
}

func (cf *Classifier) Commit() (err error) {
	if cf.tx == nil {
		return SessionError{msg: "not in transaction"}
	}
	cf.que <- trainingSample{action: TA_COMMIT}
	cf.wg.Wait()
	close(cf.que)
	return cf.err
}

func (cf *Classifier) Rollback() (err error) {
	if cf.tx == nil {
		return SessionError{msg: "not in transaction"}
	}
	cf.que <- trainingSample{action: TA_ROLLBACK}
	cf.wg.Wait()
	close(cf.que)
	return cf.err
}
