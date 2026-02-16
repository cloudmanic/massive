//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package main

import (
	"github.com/cloudmanic/massive-cli/cmd"
)

// main is the entry point for the Massive CLI application. It delegates
// all command parsing and execution to the cobra command framework.
func main() {
	cmd.Execute()
}
