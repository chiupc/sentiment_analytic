package main

import (
	"io/ioutil"
)

func clean(s []byte) []byte {
	j := 0
	var quoteInd []int
	for _, b := range s {
		if b < 128 {
			s[j] = b
			if b == 34{
				quoteInd = append(quoteInd, j)
			}else if b == 10{
				if len(quoteInd) > 2{
					for _, ind := range(quoteInd[1 : len(quoteInd)-1]) {
						s[ind] = 39
					}
				}
				quoteInd = []int{}
			}
			j++
		}
	}
	return s[:j]
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}