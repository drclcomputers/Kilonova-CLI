package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		return
	}

	cmd := os.Args[1]
	switch {
	case cmd == "-signin":
		login()
	case cmd == "-search":
		searchProblems()
	case cmd == "-submit":
		uploadCode()
	case cmd == "-submissions":
		printSubmissions()
	case cmd == "-statement":
		printStatement()
	case cmd == "-logout":
		logout()
	case cmd == "-langs":
		checklangs()
	default:
		fmt.Println(help)
	}
}
