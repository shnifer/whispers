package main

import (
	"fmt"
	"io/ioutil"
)

func main() {

	b1, err := ioutil.ReadFile("galaxy_warp.json")
	if err != nil {
		panic(err)
	}
	b2, err := ioutil.ReadFile("galaxy_IV5.json")
	if err != nil {
		panic(err)
	}
	c := NewComposer()
	c.Add("part1", b1)
	c.Add("part2", b2)
	res := c.Encode()
	was := len(b1) + len(b2)
	got := len(res)
	fmt.Println("LEN: ", was, "->", got, " (", got*100/was, "% )")
}
