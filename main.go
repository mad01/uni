package main

import (
	"github.com/mad01/uni/cmd"
)

var (
	gitHash = "dev"
	dirty   = "false"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}
