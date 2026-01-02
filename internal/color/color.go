package color

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	noColor bool
}

var Default = &Logger{}

func (l *Logger) timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *Logger) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if l.noColor {
		fmt.Printf("INFO:     [%s] %s\n", l.timestamp(), msg)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	fmt.Printf("%s %s %s\n", level, timestamp, msg)
}

func (l *Logger) Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if l.noColor {
		fmt.Printf("SUCCESS:  [%s] %s\n", l.timestamp(), msg)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgGreen, color.Bold).Sprint("SUCCESS: ")
	fmt.Printf("%s %s %s\n", level, timestamp, msg)
}

func (l *Logger) Warning(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if l.noColor {
		fmt.Printf("WARNING:  [%s] %s\n", l.timestamp(), msg)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgYellow, color.Bold).Sprint("WARNING: ")
	fmt.Printf("%s %s %s\n", level, timestamp, msg)
}

func (l *Logger) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if l.noColor {
		fmt.Printf("ERROR:    [%s] %s\n", l.timestamp(), msg)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgRed, color.Bold).Sprint("ERROR:   ")
	fmt.Printf("%s %s %s\n", level, timestamp, msg)
}

func (l *Logger) Debug(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if l.noColor {
		fmt.Printf("DEBUG:    [%s] %s\n", l.timestamp(), msg)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgMagenta, color.Bold).Sprint("DEBUG:   ")
	fmt.Printf("%s %s %s\n", level, timestamp, msg)
}

func (l *Logger) ServerRunning(url string) {
	if l.noColor {
		fmt.Printf("INFO:     [%s] Server running on %s\n", l.timestamp(), url)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	message := "Server running on"
	urlColored := color.New(color.FgBlue, color.Bold).Sprint(url)
	fmt.Printf("%s %s %s %s\n", level, timestamp, message, urlColored)
}

func (l *Logger) ConfigSet(key, value string) {
	if l.noColor {
		fmt.Printf("INFO:     [%s] Configuration set: %s=%s\n", l.timestamp(), key, value)
		return
	}
	timestamp := color.New(color.FgHiBlack).Sprintf("[%s]", l.timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	keyColored := color.New(color.FgCyan).Sprint(key)
	valueColored := color.New(color.FgGreen).Sprint(value)
	fmt.Printf("%s %s Configuration set: %s=%s\n", level, timestamp, keyColored, valueColored)
}

func (l *Logger) SetNoColor(disable bool) {
	l.noColor = disable
}
