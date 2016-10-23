package mbayes

type Score struct {
	Ident string
	Value float64
}

type Scores []Score

func (ss Scores) Len() int {
	return len(ss)
}
func (ss Scores) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
func (ss Scores) Less(i, j int) bool {
	return ss[i].Value > ss[j].Value
}
