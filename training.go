package mbayes

import (
	"database/sql"
)

func (cf *Classifier) Train(category string, tokens ...[]byte) (err error) {
	if cf.err != nil {
		return
	}
	if cf.tx == nil { //auto commit
		var tx *sql.Tx
		tx, err = cf.db.Begin()
		if err != nil {
			return
		}
		for _, tk := range tokens {
			err = cf.add(tx, category, tk)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		err = tx.Commit()
	} else {
		for _, tk := range tokens {
			cf.que <- trainingSample{
				action:  TA_TRAIN,
				feature: tk,
				label:   category,
			}
		}
	}
	return
}

func (cf *Classifier) Untrain(category string, tokens ...[]byte) (err error) {
	if cf.err != nil {
		return
	}
	if cf.tx == nil { //auto commit
		var tx *sql.Tx
		tx, err = cf.db.Begin()
		if err != nil {
			return
		}
		for _, tk := range tokens {
			err = cf.del(tx, category, tk)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		err = tx.Commit()
	} else {
		for _, tk := range tokens {
			cf.que <- trainingSample{
				action:  TA_UNTRAIN,
				feature: tk,
				label:   category,
			}
		}
	}
	return
}
