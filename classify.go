package mbayes

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func (cf *Classifier) entropy(cats []string, tokens [][]byte) (ss Scores,
	err error) {
	var args []interface{}
	for _, t := range tokens {
		args = append(args, string(t))
	}
	sql := `SELECT tok,MAX(cnt),SUM(cnt)-MAX(cnt) FROM tokens WHERE tok IN (?` +
		strings.Repeat(",?", len(tokens)-1) + `)`
	cc := len(cats)
	if cc == 0 {
		row := cf.db.QueryRow(`SELECT COUNT(DISTINCT cat) FROM tokens`)
		err = row.Scan(&cc)
		if err != nil {
			return
		}
	} else {
		sql += ` AND cat IN (?` + strings.Repeat(",?", len(cats)-1) + `)`
		for _, c := range cats {
			args = append(args, c)
		}
	}
	sql += ` GROUP BY tok`
	rows, err := cf.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			tok        string
			max, delta int
		)
		err = rows.Scan(&tok, &max, &delta)
		if err != nil {
			return
		}
		if max == 0 {
			continue
		}
		if delta == 0 {
			ss = append(ss, Score{
				Ident: tok,
				Value: 65536,
			})
		} else {
			ss = append(ss, Score{
				Ident: tok,
				Value: float64(max*cc-max) / float64(delta),
			})
		}
	}
	err = rows.Err()
	if err == nil {
		sort.Sort(ss)
	}
	return
}

func (cf *Classifier) prob(tokens map[string]int) (ss Scores, err error) {
	var args []interface{}
	for t := range tokens {
		args = append(args, t)
	}
	sql := `SELECT t2.tok,t1.cat,1.0*(t1.cnt+1)/(SUM(t2.cnt)+%d) FROM tokens t1,
	    tokens t2 WHERE t1.tok=t2.tok AND t2.tok IN (?` + strings.Repeat(",?",
		len(tokens)-1) + `)`
	sql = fmt.Sprintf(sql, cf.cats) + ` GROUP BY t1.cat,t2.tok`
	rows, err := cf.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	ms := make(map[string]Scores)
	for rows.Next() {
		var (
			tok, cat string
			ratio    float64
		)
		err = rows.Scan(&tok, &cat, &ratio)
		if err != nil {
			return
		}
		ms[cat] = append(ms[cat], Score{
			Ident: tok,
			Value: ratio,
		})
	}
	err = rows.Err()
	if err != nil {
		return
	}
	for c, ts := range ms {
		ss = append(ss, func() Score {
			p := 1.0
			q := 1.0
			f := make(map[string]int)
			for _, s := range ts {
				p *= (1 - s.Value)
				q *= s.Value
				f[s.Ident] = 1
			}
			pr := 1 / float64(cf.cats)
			for t := range tokens {
				if f[t] != 0 {
					continue
				}
				p *= (1 - pr)
				q *= pr
			}
			n := 1 / float64(len(tokens))
			p = 1 - math.Pow(p, n)
			q = 1 - math.Pow(q, n)
			return Score{
				Ident: c,
				Value: ((p-q)/(p+q) + 1) / 2,
			}
		}())
	}
	return
}

func (cf *Classifier) Classify(tokens ...[]byte) (ss Scores, err error) {
	tks, err := cf.entropy([]string{}, tokens)
	if err != nil {
		return
	}
	selected := make(map[string]int)
	for _, tk := range tks {
		if len(selected) >= 30 {
			break
		}
		selected[tk.Ident] = 1
	}
	ss, err = cf.prob(selected)
	if err != nil {
		return
	}
	sort.Sort(ss)
	return
}

func (cf *Classifier) ClassifyWithin(cats []string, tokens ...[]byte) (
	ss Scores, err error) {
	tks, err := cf.entropy(cats, tokens)
	if err != nil {
		return
	}
	_ = tks
	return
}
