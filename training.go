package mbayes

func (cf *Classifier) doSingle(act int, cat string, toks ...[]byte) error {
	tx, err := cf.db.Begin()
	if err != nil {
		return err
	}
	op := cf.add
	if act == TA_UNTRAIN {
		op = cf.del
	}
	for _, tk := range toks {
		err = op(tx, cat, tk)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (cf *Classifier) doBatch(act int, cat string, toks ...[]byte) error {
	for _, tk := range toks {
		cf.que <- trainingSample{
			action:  act,
			feature: tk,
			label:   cat,
		}
	}
	return nil
}

func (cf *Classifier) Train(category string, tokens ...[]byte) error {
	if cf.err != nil {
		return cf.err
	}
	if cf.tx == nil { //auto commit
		return cf.doSingle(TA_TRAIN, category, tokens...)
	} else {
		return cf.doBatch(TA_TRAIN, category, tokens...)
	}
}

func (cf *Classifier) Untrain(category string, tokens ...[]byte) error {
	if cf.err != nil {
		return cf.err
	}
	if cf.tx == nil { //auto commit
		return cf.doSingle(TA_UNTRAIN, category, tokens...)
	} else {
		return cf.doBatch(TA_UNTRAIN, category, tokens...)
	}
}
