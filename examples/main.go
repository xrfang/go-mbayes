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
	assert(mb.Add([]byte("hello world"), "cat1"))
	//assert(mb.Delete([]byte("hello world"), "cat1"))
	fmt.Printf("done.\n")
}
