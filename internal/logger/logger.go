// Package logger provides functions for logging messages.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

const logFilePermission os.FileMode = 0o600

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

	currentLevel atomic.Int32
)

func init() {
	file = os.Stdout
	currentLevel.Store(int32(LevelDebug))
}

func checkLevel(l LogLevel) bool {
	return l >= LogLevel(currentLevel.Load())
}

func SetLogLevel(l LogLevel) {
	currentLevel.Store(int32(l))
}

func SetupFile(logFilePath string, disableColor ...bool) error {
	logFile, err := openLogFile(logFilePath)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	if len(disableColor) > 0 {
		noColor = disableColor[0]
	}
	file = logFile
	return nil
}

func openLogFile(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePermission)
	if err == nil {
		return logFile, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	Info("The filepath `%s` doesn't exist\nDo you want to create it? [y/n]", logFilePath)
	var response string
	if _, scanErr := fmt.Scanln(&response); scanErr != nil || (response != "y" && response != "Y") {
		return nil, fmt.Errorf("log file creation aborted")
	}
	if err := os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePermission)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	return logFile, nil
}

func Close() {
	mu.Lock()
	defer mu.Unlock()

	if file != nil && file != os.Stdout {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}
}

func SetNoColor(disable bool) {
	mu.Lock()
	defer mu.Unlock()
	noColor = disable
}

func writeToFile(message string) {
	if file != nil && file != os.Stdout {
		if _, err := fmt.Fprint(file, strings.TrimSuffix(message, "\n")+"\n"); err != nil {
			fmt.Printf("failed to write to log file: %v", err)
		}
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
	plain := fmt.Sprintf("INFO:     [%s] %s\n", timestamp(), msg)
	writeToFile(plain)
	if noColor {
		fmt.Print(plain)
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
	plain := fmt.Sprintf("SUCCESS:  [%s] %s\n", timestamp(), msg)
	writeToFile(plain)
	if noColor {
		fmt.Print(plain)
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

func Errorln(args ...any) {
	if !checkLevel(LevelError) {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	msg := fmt.Sprintln(args...)
	msg = strings.TrimSuffix(msg, "\n")
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

func ErrorWriteOnlyFile(format string, args ...any) {
	if !checkLevel(LevelError) {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	caller := getCaller(2)

	plain := fmt.Sprintf("ERROR:    [%s] [%s] %s\n", timestamp(), caller, msg)
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
