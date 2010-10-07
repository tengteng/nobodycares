package main

import (
	"flag"
	"fmt"
)

var input = flag.String("i", "", "input to sha256 hash")

func main() {
	flag.Parse()
	if len(*input) > 0 {
		fmt.Printf("%s\n", Hash(*input))
	} else {
		fmt.Printf("usage: pwhash -i=<password>\n")
	}
}
