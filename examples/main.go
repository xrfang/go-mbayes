package main

import (
	"fmt"

	"github.com/xrfang/go-mbayes"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mb, err := mbayes.Open("mbayes.db")
	assert(err)
	defer mb.Close()
	assert(mb.Train("cat1", []byte("hello"), []byte("world")))
	assert(mb.Untrain("cat1", []byte("hello world")))
	fmt.Printf("done.\n")
}
