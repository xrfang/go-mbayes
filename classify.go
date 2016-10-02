package mbayes

type BayesianScore struct {
	Label string
	Score float64
}
type BayesianScores []BayesianScore

func (cf *Classifier) Classify(limit *[]string, tokens ...[]byte) (
	scores BayesianScores) {

}
