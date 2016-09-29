package mbayes

func (cf *Classifier) Train(category string, tokens ...[]byte) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	for _, tk := range tokens {
		_, err = tx.Exec(SQL("addtok1"), tk, category)
		if err != nil {
			return
		}
		_, err = tx.Exec(SQL("addtok2"), tk, category)
		if err != nil {
			return
		}
	}
	return
}

func (cf *Classifier) Untrain(category string, tokens ...[]byte) (err error) {
	tx, err := cf.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	for _, tk := range tokens {
		_, err = tx.Exec(SQL("deltok1"), tk, category)
		if err != nil {
			return
		}
		_, err = tx.Exec(SQL("deltok2"), tk, category)
		if err != nil {
			return
		}
	}
	return
}
