package main

import (
	"fmt"
	"os"
)

const help = "Kilonova CLI - ver 0.0.9\n\n-signin <USERNAME> <PASSWORD>\n-langs <ID>\n-search <PROBLEM ID or NAME>\n-submit <PROBLEM ID> <LANGUAGE> <solution>\n-statement <PROBLEM ID>\n-logout"

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
	} else {
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
}
