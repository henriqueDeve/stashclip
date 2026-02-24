package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"stashclip/internal/clipboard"
	"stashclip/internal/daemon"
	"stashclip/internal/popup"
	"stashclip/internal/store"
)

// Run executes the CLI command based on args.
func Run(args []string) error {
	if len(args) < 2 {
		return runPopup()
	}

	switch args[1] {
	case "menu":
		return runPopup()
	case "popup":
		return runPopup()
	case "__daemon-run":
		return runDaemonForeground()
	case "-h", "--help", "help":
		usage()
		return nil
	default:
		return fmt.Errorf("unsupported command: use Ctrl+Alt+A to open stashclip popup")
	}
}

func usage() {
	fmt.Println("Usage: stashclip [popup]")
	fmt.Println()
	fmt.Println("Open popup and choose an item to copy again.")
	fmt.Println("Recommended usage: global shortcut Ctrl+Alt+A -> stashclip-popup")
}

func runDaemonCommand(args []string) error {
	if len(args) == 0 || args[0] == "start" {
		return startDaemon()
	}
	if len(args) != 1 {
		return fmt.Errorf("daemon error: invalid arguments")
	}

	switch args[0] {
	case "run":
		return runDaemonForeground()
	case "stop":
		return stopDaemon()
	case "status":
		return daemonStatus()
	default:
		return fmt.Errorf("daemon error: unknown action: %s", args[0])
	}
}

func runDaemonForeground() error {
	clipboardProvider, err := clipboard.NewProvider()
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	memStore, err := newStore()
	if err != nil {
		return err
	}
	if err := daemon.Run(clipboardProvider, memStore); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	return nil
}

func startDaemon() error {
	pidPath := daemonPIDPath()
	pid, running, err := readDaemonPID(pidPath)
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if running {
		fmt.Printf("daemon already running (pid %d)\n", pid)
		return nil
	}
	if pid != 0 {
		_ = os.Remove(pidPath)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	logFile, err := openDaemonLog()
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	defer logFile.Close()

	cmd := exec.Command(exe, "daemon", "run")
	cmd.Stdin = nil
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if err := writeDaemonPID(pidPath, cmd.Process.Pid); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	time.Sleep(700 * time.Millisecond)
	if !processRunning(cmd.Process.Pid) {
		_ = os.Remove(pidPath)
		return fmt.Errorf("daemon error: failed to stay running (check %s)", daemonLogPath())
	}
	fmt.Printf("daemon started (pid %d)\n", cmd.Process.Pid)
	return nil
}

func stopDaemon() error {
	pidPath := daemonPIDPath()
	pid, running, err := readDaemonPID(pidPath)
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if !running {
		_ = os.Remove(pidPath)
		return fmt.Errorf("daemon error: not running")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if waitProcessStop(pid, 3*time.Second) {
		_ = os.Remove(pidPath)
		fmt.Printf("daemon stopped (pid %d)\n", pid)
		return nil
	}

	if err := proc.Signal(syscall.SIGKILL); err != nil {
		return fmt.Errorf("daemon error: failed to stop pid %d after timeout: %w", pid, err)
	}
	if waitProcessStop(pid, 1*time.Second) {
		_ = os.Remove(pidPath)
		fmt.Printf("daemon stopped (pid %d)\n", pid)
		return nil
	}

	return fmt.Errorf("daemon error: process %d did not stop", pid)
}

func daemonStatus() error {
	pid, running, err := readDaemonPID(daemonPIDPath())
	if err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}
	if running {
		fmt.Printf("daemon running (pid %d)\n", pid)
	} else {
		fmt.Println("daemon not running")
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

	selected := len(entries)
	if len(args) > 1 {
		return fmt.Errorf("pick error: too many arguments")
	}
	if len(args) == 1 {
		n, convErr := strconv.Atoi(args[0])
		if convErr != nil {
			return fmt.Errorf("pick error: invalid index: %s", args[0])
		}
		selected = n
	}
	return writePickByIndex(entries, selected)
}

func runPopup() error {
	for {
		memStore, err := newStore()
		if err != nil {
			return err
		}
		entries := memStore.List()
		if len(entries) == 0 {
			return fmt.Errorf("popup error: no entries available")
		}
		items := make([]popup.Item, 0, len(entries))
		for i, entry := range entries {
			items = append(items, popup.Item{
				ID:      i + 1,
				AddedAt: entry.AddedAt,
				Text:    entry.Text,
			})
		}
		selected, err := popup.Select(items)
		if err != nil {
			if errors.Is(err, popup.ErrCanceled) {
				return nil
			}
			return err
		}
		if err := writePickByIndex(entries, selected); err != nil {
			return err
		}
	}
}

func writePickByIndex(entries []store.Entry, oneBasedIndex int) error {
	if oneBasedIndex < 1 || oneBasedIndex > len(entries) {
		return fmt.Errorf("pick error: index out of range: %d", oneBasedIndex)
	}
	clipboardProvider, err := clipboard.NewProvider()
	if err != nil {
		return fmt.Errorf("pick error: %w", err)
	}
	if err := clipboard.MarkIgnored(entries[oneBasedIndex-1].Text); err != nil {
		return fmt.Errorf("pick error: %w", err)
	}
	if err := clipboardProvider.Write(entries[oneBasedIndex-1].Text); err != nil {
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

func daemonPIDPath() string {
	return filepath.Join(daemonStateDir(), "daemon.pid")
}

func daemonLogPath() string {
	return filepath.Join(daemonStateDir(), "daemon.log")
}

func openDaemonLog() (*os.File, error) {
	path := daemonLogPath()
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

func writeDaemonPID(path string, pid int) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
}

func readDaemonPID(path string) (int, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, false, nil
		}
		return 0, false, err
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		return 0, false, nil
	}
	pid, err := strconv.Atoi(s)
	if err != nil || pid <= 0 {
		return 0, false, fmt.Errorf("invalid pid file")
	}
	return pid, processRunning(pid), nil
}

func processRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if proc.Signal(syscall.Signal(0)) != nil {
		return false
	}
	return !isZombieProcess(pid)
}

func daemonStateDir() string {
	path := store.DefaultPath()
	if path != "" {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err == nil && isDirWritable(dir) {
			return dir
		}
	}
	return "/tmp"
}

func isDirWritable(dir string) bool {
	f, err := os.CreateTemp(dir, ".stashclip-writecheck-")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}

func isZombieProcess(pid int) bool {
	statPath := filepath.Join("/proc", strconv.Itoa(pid), "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return false
	}
	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return false
	}
	return parts[2] == "Z"
}

func waitProcessStop(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !processRunning(pid) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return !processRunning(pid)
}
