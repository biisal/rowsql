package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota + 1
	LevelInfo
	LevelSuccess
	LevelWarning
	LevelError
)

var (
	noColor bool
	file    *os.File
	mu      sync.Mutex
	level   = LevelInfo
)

func init() {
	file = os.Stdout
}

func checkLevel(l LogLevel) bool {
	return l >= level
}

func SetLogLevel(l LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	level = l
}

func SetupFile(logFilePath string, disableColor ...bool) error {
	mu.Lock()
	defer mu.Unlock()

	if len(disableColor) > 0 {
		noColor = disableColor[0]
	}

	if strings.HasPrefix(logFilePath, "~") {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return fmt.Errorf("failed to get user home directory: %w", homeErr)
		}
		logFilePath = strings.Replace(logFilePath, "~", homeDir, 1)
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			Info("The filepath `%s` doesn't exist\nDo you want to create it? [y/n]", logFilePath)
			var response string
			_, scanErr := fmt.Scanln(&response)
			if scanErr != nil || (response != "y" && response != "Y") {
				return fmt.Errorf("log file creation aborted")
			}
			dir := filepath.Dir(logFilePath)
			if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
				return fmt.Errorf("failed to create log directory: %w", mkErr)
			}
			f, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf("failed to create log file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	file = f
	return nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if file != nil && file != os.Stdout {
		return file.Close()
	}
	return nil
}

func SetNoColor(disable bool) {
	mu.Lock()
	defer mu.Unlock()
	noColor = disable
}

func writeToFile(message string) {
	if file != nil && file != os.Stdout {
		fmt.Fprintln(file, message)
	}
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func Info(format string, args ...any) {
	if !checkLevel(LevelInfo) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Printf("INFO:     [%s] %s\n", timestamp(), msg)
		return
	}
	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	fmt.Printf("%s %s %s\n", level, ts, msg)
}

func Success(format string, args ...any) {
	if !checkLevel(LevelSuccess) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Printf("SUCCESS:  [%s] %s\n", timestamp(), msg)
		return
	}
	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgGreen, color.Bold).Sprint("SUCCESS: ")
	fmt.Printf("%s %s %s\n", level, ts, msg)
}

func Warning(format string, args ...any) {
	if !checkLevel(LevelWarning) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	caller := getCaller(2)

	if noColor {
		logMsg := fmt.Sprintf("WARNING:  [%s] [%s] %s\n", timestamp(), caller, msg)
		fmt.Print(logMsg)
		writeToFile(logMsg)
		return
	}

	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgYellow, color.Bold).Sprint("WARNING: ")
	callerInfo := color.New(color.FgCyan).Sprintf("[%s]", caller)
	colored := fmt.Sprintf("%s %s %s %s\n", level, ts, callerInfo, msg)
	plain := fmt.Sprintf("WARNING:  [%s] [%s] %s\n", timestamp(), caller, msg)

	fmt.Print(colored)
	writeToFile(plain)
}

func Error(format string, args ...any) {
	if !checkLevel(LevelError) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	caller := getCaller(2)

	if noColor {
		logMsg := fmt.Sprintf("ERROR:    [%s] [%s] %s\n", timestamp(), caller, msg)
		fmt.Print(logMsg)
		writeToFile(logMsg)
		return
	}

	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgRed, color.Bold).Sprint("ERROR:   ")
	callerInfo := color.New(color.FgYellow).Sprintf("[%s]", caller)

	colored := fmt.Sprintf("%s %s %s %s\n", level, ts, callerInfo, msg)
	plain := fmt.Sprintf("ERROR:    [%s] [%s] %s\n", timestamp(), caller, msg)

	fmt.Print(colored)
	writeToFile(plain)
}

func Debug(format string, args ...any) {
	if !checkLevel(LevelDebug) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Printf("DEBUG:    [%s] %s\n", timestamp(), msg)
		return
	}
	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgMagenta, color.Bold).Sprint("DEBUG:   ")
	fmt.Printf("%s %s %s\n", level, ts, msg)
}

func ServerRunning(url string) {
	mu.Lock()
	defer mu.Unlock()

	if noColor {
		fmt.Printf("INFO:     [%s] Server running on %s\n", timestamp(), url)
		return
	}
	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	message := "Server running on"
	urlColored := color.New(color.FgBlue, color.Bold).Sprint(url)
	fmt.Printf("%s %s %s %s\n", level, ts, message, urlColored)
}

func ConfigSet(key, value string) {
	mu.Lock()
	defer mu.Unlock()

	if noColor {
		fmt.Printf("INFO:     [%s] Configuration set: %s=%s\n", timestamp(), key, value)
		return
	}
	ts := color.New(color.FgHiBlack).Sprintf("[%s]", timestamp())
	level := color.New(color.FgCyan, color.Bold).Sprint("INFO:    ")
	keyColored := color.New(color.FgCyan).Sprint(key)
	valueColored := color.New(color.FgGreen).Sprint(value)
	fmt.Printf("%s %s Configuration set: %s=%s\n", level, ts, keyColored, valueColored)
}
