package mbayes

func (cf *classifier) Train(tokens [][]byte, category string) (err error) {
	for _, tk := range tokens {
		err = cf.add(tk, category)
		if err != nil {
			return
		}
	}
	return
}

func (cf *classifier) Untrain(tokens [][]byte, category string) (err error) {
	for _, tk := range tokens {
		err = cf.Untrain(tk, category)
		if err != nil {
			return
		}
	}
	return
}
