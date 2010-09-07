package main

import (
    "nobodycares"
    "flag"
    "fmt"
)

var input = flag.String("i", "", "input to sha256 hash")

func main() {
    flag.Parse()
    if len(*input) > 0 {
        fmt.Printf("%s\n", nobodycares.Hash(*input))
    } else {
        fmt.Printf("usage: pwhash -i=<password>\n")
    }
}
