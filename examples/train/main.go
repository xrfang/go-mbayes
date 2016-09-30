package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/xrfang/go-mbayes"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	r    *regexp.Regexp
	c    *mbayes.Classifier
	wg   sync.WaitGroup
	rate chan int
)

func init() {
	r = regexp.MustCompile("[a-zA-Z0-9]+")
}

func train(path string) {
	cat := filepath.Base(filepath.Dir(path))
	f, err := os.Open(path)
	assert(err)
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	fmt.Println(cat, path)
	dedup := make(map[string]bool)
	for _, tk := range r.FindAll(data, len(data)) {
		dedup[string(tk)] = true
	}
	var tokens [][]byte
	for tk := range dedup {
		tokens = append(tokens, []byte(tk))
	}
	assert(c.Train(cat, tokens...))
	wg.Done()
	<-rate
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("USAGE: %s <dbname> <training-data-dir>\n", filepath.Base(
			os.Args[0]))
		return
	}
	db := os.Args[1]
	root := os.Args[2]
	var err error
	c, err = mbayes.Open(db)
	assert(err)
	assert(c.Begin())
	rate = make(chan int, 1000)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rate <- 1
		wg.Add(1)
		go train(path)
		return nil
	})
	wg.Wait()
	assert(c.Commit())
	close(rate)
}
