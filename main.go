// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// Package main is the entry point for the Kilonova CLI application.
// It initializes and executes the command-line interface.
package main

import "kncli/cmd"

// main is the entry point of the Kilonova CLI application.
// It calls the Execute function from the cmd package to start the CLI.
func main() {
	cmd.Execute()
}
