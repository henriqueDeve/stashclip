package popup

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ErrCanceled indicates the user closed/canceled the popup.
var ErrCanceled = errors.New("popup canceled")

// Item is a selectable clipboard entry shown in the popup.
type Item struct {
	ID      int
	AddedAt time.Time
	Text    string
}

// Select opens a popup and returns the selected item ID (1-based).
func Select(items []Item) (int, error) {
	if len(items) == 0 {
		return 0, fmt.Errorf("popup error: no entries available")
	}

	name := preferredProvider()
	switch name {
	case "yad":
		return selectWithYad(items)
	case "zenity":
		return selectWithZenity(items)
	case "kdialog":
		return selectWithKdialog(items)
	default:
		return 0, fmt.Errorf("popup error: no supported popup backend found (install one of: yad, zenity, kdialog)")
	}
}

func preferredProvider() string {
	if forced := strings.ToLower(strings.TrimSpace(os.Getenv("STASHCLIP_POPUP_PROVIDER"))); forced != "" {
		if hasCommand(forced) {
			return forced
		}
		return ""
	}
	for _, p := range []string{"yad", "zenity", "kdialog"} {
		if hasCommand(p) {
			return p
		}
	}
	return ""
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func selectWithYad(items []Item) (int, error) {
	args := []string{
		"--list",
		"--title=Stashclip",
		"--text=Selecione um item para copiar",
		"--width=980",
		"--height=600",
		"--button=Copiar:0",
		"--button=Fechar:1",
		"--column=ID:NUM",
		"--column=Data:TEXT",
		"--column=Texto:TEXT",
		"--print-column=1",
		"--separator=\n",
	}
	for _, item := range items {
		args = append(args, strconv.Itoa(item.ID), item.AddedAt.Format(time.RFC3339), sanitize(item.Text))
	}

	cmd := exec.Command("yad", args...)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return 0, ErrCanceled
		}
		return 0, fmt.Errorf("popup error: %w", err)
	}
	return parseSelectedID(out)
}

func selectWithZenity(items []Item) (int, error) {
	args := []string{
		"--list",
		"--title=Stashclip",
		"--text=Selecione um item para copiar",
		"--width=980",
		"--height=600",
		"--ok-label=Copiar",
		"--cancel-label=Fechar",
		"--column=ID",
		"--column=Data",
		"--column=Texto",
		"--hide-column=1",
		"--print-column=1",
	}
	for _, item := range items {
		args = append(args, strconv.Itoa(item.ID), item.AddedAt.Format(time.RFC3339), sanitize(item.Text))
	}

	cmd := exec.Command("zenity", args...)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return 0, ErrCanceled
		}
		return 0, fmt.Errorf("popup error: %w", err)
	}
	return parseSelectedID(out)
}

func selectWithKdialog(items []Item) (int, error) {
	args := []string{
		"--title", "Stashclip",
		"--menu", "Selecione um item para copiar",
	}
	for _, item := range items {
		args = append(args, strconv.Itoa(item.ID), fmt.Sprintf("%s  %s", item.AddedAt.Format(time.RFC3339), sanitize(item.Text)))
	}

	cmd := exec.Command("kdialog", args...)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return 0, ErrCanceled
		}
		return 0, fmt.Errorf("popup error: %w", err)
	}
	return parseSelectedID(out)
}

func parseSelectedID(out []byte) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("popup error: invalid selection")
	}
	return id, nil
}

func sanitize(text string) string {
	s := strings.ReplaceAll(text, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
