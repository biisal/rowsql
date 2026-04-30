package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	AnsiBackCursorToLineStart = "\x1b[1G"
	AnsiClearScreen           = "\033[2J\033[H"
	AnsiHideCursor            = "\x1b[?25l"
	AnsiShowCursor            = "\x1b[?25h"
	AnsiAlternateScreen       = "\x1b[?1049h"
	AnsiExitAlternateScreen   = "\x1b[?1049l"
)

type Prompter interface {
	AskConnection(connections []ConnectionConfig) (*ConnectionConfig, error)
}

type PrompterImpl struct{}

func promptForDefaultEnv(dir, fileName string) error {
	path := filepath.Join(dir, fileName)
	color.Cyan("No %s found in %s\nDo you want to create one with default values? (y/n): ", fileName, dir)
	var choice string
	if _, err := fmt.Scan(&choice); err != nil {
		return fmt.Errorf("reading choice: %w", err)
	}

	if strings.ToLower(choice) != "y" {
		return fmt.Errorf("no %s file found and user declined creation", fileName)
	}

	cfg := DefaultConfig()
	fields := ConfigFileFields{
		Connections: []ConnectionConfig{cfg.ConnectionConfig},
		AppConfig:   cfg.AppConfig,
	}

	data, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	logger.Info("Default %s file created in %s", fileName, path)
	return nil
}

func makeListString(configs []ConnectionConfig, selected int) string {
	var sb strings.Builder
	sb.WriteString(AnsiClearScreen)
	header := color.HiGreenString("Choose a port to run:")
	helpText := color.HiBlueString("use j/k to up down")
	fmt.Fprintf(&sb, "%s\n%s%s\n%s\n%s\n", header, AnsiBackCursorToLineStart, helpText, AnsiBackCursorToLineStart, AnsiBackCursorToLineStart)
	for i, cfg := range configs {
		if i == selected {
			fmt.Fprintf(&sb, "  %s Port: %s  DB: %s\n%s",
				color.HiGreenString("▶ %d.", i),
				color.HiGreenString(fmt.Sprintf("%d", cfg.Port)),
				color.HiGreenString(cfg.DBString),
				AnsiBackCursorToLineStart,
			)
		} else {
			fmt.Fprintf(&sb, "    %d. Port: %d  DB: %s\n%s",
				i, cfg.Port, cfg.DBString, AnsiBackCursorToLineStart,
			)
		}
	}
	return sb.String()
}

func (p *PrompterImpl) AskConnection(configs []ConnectionConfig) (*ConnectionConfig, error) {
	oldState, err := term.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to make terminal raw: %w", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdout.Fd()), oldState); err != nil {
			logger.Error("Failed to restore terminal: %s", err)
		}
		fmt.Print(AnsiShowCursor)
		fmt.Print(AnsiExitAlternateScreen)
	}()
	fmt.Print(AnsiAlternateScreen)
	fmt.Print(AnsiHideCursor)

	if len(configs) == 0 {
		return nil, apperr.ErrorNoConfigsFound
	}

	configsLastIdx := len(configs) - 1
	currentIndex := 0

	fmt.Println(makeListString(configs, currentIndex))

	input := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(input)
		if err != nil || n == 0 {
			continue
		}
		char := input[0]
		switch char {
		case 'q':
			return nil, fmt.Errorf("operation cancelled by user")
		case 'j':
			if currentIndex < configsLastIdx {
				currentIndex++
				fmt.Println(makeListString(configs, currentIndex))
			}
		case 'k':
			if currentIndex > 0 {
				currentIndex--
				fmt.Println(makeListString(configs, currentIndex))
			}
		case '\n', '\r':
			return &configs[currentIndex], nil
		default:
			// Ignore other inputs
		}
	}
}

