package cli

import (
	"fmt"
	"strconv"
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
		return runPick(args[2:])
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

func runPick(args []string) error {
	memStore, err := newStore()
	if err != nil {
		return err
	}
	entries := memStore.List()
	if len(entries) == 0 {
		return fmt.Errorf("pick error: no entries available")
	}

	selected := len(entries) - 1
	if len(args) > 1 {
		return fmt.Errorf("pick error: too many arguments")
	}
	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("pick error: invalid index: %s", args[0])
		}
		if n < 1 || n > len(entries) {
			return fmt.Errorf("pick error: index out of range: %d", n)
		}
		selected = n - 1
	}

	clipboardProvider := clipboard.NewX11()
	if err := clipboardProvider.Write(entries[selected].Text); err != nil {
		return fmt.Errorf("pick error: %w", err)
	}
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
