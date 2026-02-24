package main

import (
	"os"
	"stashclip/internal/cli"
)

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	if err := cli.Run(args); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}
	return 0
}
