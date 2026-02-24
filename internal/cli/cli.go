package cli

import (
	"fmt"
	"strings"
	"time"

	"stashclip/internal/clipboard"
	"stashclip/internal/daemon"
	"stashclip/internal/store"
)

// Run executes the CLI command based on args.
func Run(args []string) error {
	if len(args) < 2 {
		usage()
		return nil
	}

	switch args[1] {
	case "daemon":
		return runDaemon()
	case "list":
		return runList()
	case "pick":
		return runPick()
	case "clear":
		return runClear()
	case "-h", "--help", "help":
		usage()
		return nil
	default:
		fmt.Printf("unknown command: %s\n\n", args[1])
		usage()
		return nil
	}
}

func usage() {
	fmt.Println("Usage: stashclip <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  daemon  Start clipboard monitoring daemon")
	fmt.Println("  list    List stored entries")
	fmt.Println("  pick    Pick an entry to paste")
	fmt.Println("  clear   Clear stored entries")
}

func runDaemon() error {
	clipboardProvider := clipboard.NewX11()
	memStore, err := newStore()
	if err != nil {
		return err
	}
	if err := daemon.Run(clipboardProvider, memStore); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	return nil
}

func runList() error {
	memStore, err := newStore()
	if err != nil {
		return err
	}
	entries := memStore.List()
	for i, entry := range entries {
		text := strings.ReplaceAll(entry.Text, "\n", "\\n")
		text = strings.ReplaceAll(text, "\t", "\\t")
		fmt.Printf("%d\t%s\t%s\n", i+1, entry.AddedAt.Format(time.RFC3339), text)
	}
	return nil
}

func runPick() error {
	fmt.Println("pick: not implemented")
	return nil
}

func runClear() error {
	memStore, err := newStore()
	if err != nil {
		return err
	}
	memStore.Clear()
	return nil
}

func newStore() (*store.Store, error) {
	memStore, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("store error: %w", err)
	}
	return memStore, nil
}
