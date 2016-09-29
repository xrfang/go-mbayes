package mbayes

func (cf *classifier) Train(category string, tokens ...[]byte) (err error) {
	for _, tk := range tokens {
		err = cf.add(category, tk)
		if err != nil {
			return
		}
	}
	return
}

func (cf *classifier) Untrain(category string, tokens ...[]byte) (err error) {
	for _, tk := range tokens {
		err = cf.delete(category, tk)
		if err != nil {
			return
		}
	}
	return
}
