package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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
		memStore, err := store.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "store error: %v\n", err)
			os.Exit(1)
		}
		if err := daemon.Run(backend, memStore); err != nil {
			fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		memStore, err := store.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "store error: %v\n", err)
			os.Exit(1)
		}
		entries := memStore.List()
		for i, entry := range entries {
			text := strings.ReplaceAll(entry.Text, "\n", "\\n")
			text = strings.ReplaceAll(text, "\t", "\\t")
			fmt.Printf("%d\t%s\t%s\n", i+1, entry.AddedAt.Format(time.RFC3339), text)
		}
	case "pick":
		fmt.Println("pick: not implemented")
	case "clear":
		memStore, err := store.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "store error: %v\n", err)
			os.Exit(1)
		}
		memStore.Clear()
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Printf("unknown command: %s\n\n", os.Args[1])
		usage()
	}
}
