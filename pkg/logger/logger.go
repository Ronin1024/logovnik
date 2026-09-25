// Search service - micro service for gruzoperevozki-rf.com
//
// Copyright 2025 Ivanov Vladimir Vyacheslavovich aka Ronin1024. All rights reserved.

package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG int8 = 0
	INFO  int8 = 1
	ERROR int8 = 2
	FATAL int8 = 3
)

type Config struct {
	Level   int8
	LogDir  string
	LogFile string
}

const logBufferLength int = 100

var LogSettings Config

var redisctx context.Context

var l_stdout *log.Logger
var l_stderr *log.Logger

var logBufferChannel chan (string)

// Инициализация пакета
func init() {
	logBufferChannel = make(chan string, logBufferLength)
	l_stdout = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	l_stderr = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)

	LogSettings.Level = DEBUG
	LogSettings.LogDir = ""
	LogSettings.LogFile = "search_service " + GetStringTimeStampInFilename() + ".log"
	go saveLogsCoroutine()
	Debug("logger:init", "Init")
}

// Установка уровня логирования
func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		LogSettings.Level = DEBUG
	case "info":
		LogSettings.Level = INFO
	case "error":
		LogSettings.Level = ERROR
	case "fatal":
		LogSettings.Level = FATAL
	default:
		LogSettings.Level = DEBUG
	}
}

func getLogFileName() string {
	return getLogDirName() + LogSettings.LogFile
}

func getLogDirName() string {
	return fmt.Sprintf("%s%s%s%s", filepath.Dir(os.Args[0]), string(os.PathSeparator), "logs", string(os.PathSeparator))
}

func getLogFile() (*os.File, error) {
	if LogSettings.LogFile == "" {
		return nil, os.ErrNotExist
	}
	localLogDir := getLogDirName()
	os.MkdirAll(localLogDir, os.ModePerm)

	file, err := os.OpenFile(getLogFileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}
	return file, err
}

func SaveLogsToFile(file *os.File, localLogDir string) {
	if LogSettings.LogDir != "" {
		if localLogDir != LogSettings.LogDir {
			file.Close()
			os.MkdirAll(LogSettings.LogDir, os.ModePerm)
			os.Rename(localLogDir+LogSettings.LogFile, LogSettings.LogDir+LogSettings.LogFile)
			newfile, err := os.OpenFile(LogSettings.LogDir+LogSettings.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			file = newfile
			localLogDir = LogSettings.LogDir
			Debug("logger:saveLogsToFile", "Log file moved to "+localLogDir+LogSettings.LogFile)
		}
	}

	for {
		logText := <-logBufferChannel
		if logText != "" {
			var nowTime []byte
			TimeStampInFile(&nowTime)
			file.Write(nowTime)
			file.Write([]byte(logText + "\n"))
		}
	}
}

// Корутина, отвечающая за сохранение логов в файл.
func saveLogsCoroutine() {
	file, _ := getLogFile()
	localLogDir := getLogDirName()
	defer file.Close()
	Debug("logger:saveLogsToFile", "Log save to "+getLogFileName())
	for {
		time.Sleep(10 * time.Millisecond)
		//При обновлении настроек логгирования - переносим лог файл в новое место
		SaveLogsToFile(file, localLogDir)
	}
}

// Добавление строки в буффер лога
func saveLogsToBuffer(text *string) {
	logBufferChannel <- *text
}

// Получить Имя функции и пакета
func getShortFunctionName(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	// Получаем полное имя: "main.main" или "github.com/user/pkg.Function"
	fullName := fn.Name()

	// Обрезаем до последней точки
	// parts := strings.Split(fullName, ".")

	// return parts[len(parts)-1]
	return fullName + ":" + strconv.Itoa(line)
}

func outputLog(logger *log.Logger, ident string, logLevel string, v ...any) {
	logtext := fmt.Sprint(fmt.Sprintf("[%s:%s]   ", logLevel, ident), fmt.Sprint(v...))
	logger.Println(logtext)
	saveLogsToBuffer(&logtext)
}

// Функция вывода ошибки
func Error(ident string, v ...any) {
	if LogSettings.Level > ERROR {
		return
	}

	outputLog(l_stderr, ident+":"+getShortFunctionName(2), "ERROR", v...)
	// outputLog(l_stderr, ident+":"+getShortFunctionName(2), "STACK", string(debug.Stack()))
}

func Info(ident string, v ...any) {
	if LogSettings.Level > INFO {
		return
	}
	outputLog(l_stdout, ident, "INFO", v...)
}

func Warn(ident string, v ...any) {
	if LogSettings.Level > INFO {
		return
	}
	outputLog(l_stderr, ident+":"+getShortFunctionName(2), "WARN", v...)
}

func Debug(ident string, v ...any) {
	if LogSettings.Level > DEBUG {
		return
	}
	outputLog(l_stdout, ident+":"+getShortFunctionName(2), "DEBUG", v...)
}

// output and exit (1)
func Fatal(ident string, v ...any) {
	if LogSettings.Level > FATAL {
		return
	}
	outputLog(l_stderr, ident+":"+getShortFunctionName(2), "FATAL", v...)
	outputLog(l_stderr, ident+":"+getShortFunctionName(2), "STACK", string(debug.Stack()))
	file, _ := getLogFile()
	SaveLogsToFile(file, getLogDirName())
	file.Close()
	os.Exit(1)
}
