package main

import (
	"flag"
	"fmt"
)

func main() {
	target := flag.String("target", "tutorial", "select target")
	flag.Parse()

	switch *target {
	case "tutorial":
		printTutorial()
	case "box":
		printBox()
	default:
		fmt.Println("invalid target")
	}
}
