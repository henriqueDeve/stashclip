package main

import (
	"fmt"
	"os"

	"stashclip/internal/clipboard"
	"stashclip/internal/daemon"
	"stashclip/internal/store"
)

func usage() {
	fmt.Println("Usage: stashclip <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  daemon  Start clipboard monitoring daemon")
	fmt.Println("  list    List stored entries")
	fmt.Println("  pick    Pick an entry to paste")
	fmt.Println("  clear   Clear stored entries")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "daemon":
		backend := clipboard.NewX11()
		memStore := store.New()
		if err := daemon.Run(backend, memStore); err != nil {
			fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		fmt.Println("list: not implemented")
	case "pick":
		fmt.Println("pick: not implemented")
	case "clear":
		fmt.Println("clear: not implemented")
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Printf("unknown command: %s\n\n", os.Args[1])
		usage()
	}
}
