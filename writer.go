// Package: fileLogger
// File: writer.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage:
// DATE: 14-8-24 12:40
package fileLogger

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
	"runtime/debug"
)

const (
	DEFAULT_PRINT_INTERVAL = 300
)

// Receive logStr from f's logChan and print logstr to file
func (f *FileLogger) logWriter() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("FileLogger's LogWritter() catch panic: %v\n", err)
		}
	}()

	//TODO let printInterVal can be configure, without change sourceCode
	printInterval := DEFAULT_PRINT_INTERVAL

	seqTimer := time.NewTicker(time.Duration(printInterval) * time.Second)
	for {
		select {
		case str := <-f.logChan:
			f.p(str)
		case <-seqTimer.C:

		}
	}
}

// print log
func (f *FileLogger) p(str string) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	f.lg.Output(2, str)
	f.pc(str)
}

// print log in console, default log string wont be print in console
// NOTICE: when console is on, the process will really slowly
func (f *FileLogger) pc(str string) {
	if f.logConsole {
		if log.Prefix() != f.prefix {
			log.SetPrefix(f.prefix)
		}
		log.Println(str)
	}
}

func getPrefix(depth int) (string, string, string) {
	pc, file, line, _ := runtime.Caller(depth)
	function := runtime.FuncForPC(pc)
	funcName := function.Name()
	arr := strings.Split(funcName, ".")
	funcName = arr[len(arr)-1]
	return shortFileName(file), funcName, strconv.Itoa(line)
}

// Printf throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (f *FileLogger) Printf(depth int, format string, v ...interface{}) {
	fileName, funcName, line := getPrefix(depth)
	f.logChan <- fmt.Sprintf("[%v:%v:%v] ", fileName, funcName, line) + fmt.Sprintf(format, v...)
}

// Print throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (f *FileLogger) Print(depth int, v ...interface{}) {
	fileName, funcName, line := getPrefix(depth)
	f.logChan <- fmt.Sprintf("[%v:%v:%v] ", fileName, funcName, line) + fmt.Sprint(v...)
}

// Println throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (f *FileLogger) Println(depth int, v ...interface{}) {
	fileName, funcName, line := getPrefix(depth)
	f.logChan <- fmt.Sprintf("[%v:%v:%v]", fileName, funcName, line) + fmt.Sprintln(v...)
}

//======================================================================================================================
// Trace log
func (f *FileLogger) Trace(depth int, format string, v ...interface{}) {
	if f.logLevel <= TRACE {
		fileName, funcName, line := getPrefix(depth)
		f.logChan <- fmt.Sprintf("[%v:%v:%v]", fileName, funcName, line) + fmt.Sprintf("[TRACE]"+format, v...)
	}
}

// same with Trace()
func (f *FileLogger) T(depth int, format string, v ...interface{}) {
	f.Trace(depth, format, v...)
}

// info log
func (f *FileLogger) Info(depth int, format string, v ...interface{}) {
	if f.logLevel <= INFO {
		fileName, funcName, line := getPrefix(depth)
		f.logChan <- fmt.Sprintf("[%v:%v:%v]", fileName, funcName, line) + fmt.Sprintf("[INFO]"+format, v...)
	}
}

// same with Info()
func (f *FileLogger) I(depth int, format string, v ...interface{}) {
	f.Info(depth, format, v...)
}

// warning log
func (f *FileLogger) Warn(depth int, format string, v ...interface{}) {
	if f.logLevel <= WARN {
		fileName, funcName, line := getPrefix(depth)
		f.logChan <- fmt.Sprintf("[%v:%v:%v]", fileName, funcName, line) + fmt.Sprintf("\033[1;33m[WARN]\033[0m"+format, v...)
	}
}

// same with Warn()
func (f *FileLogger) W(depth int, format string, v ...interface{}) {
	f.Warn(depth, format, v...)
}

// error log
func (f *FileLogger) Error(depth int, format string, v ...interface{}) {
	if f.logLevel <= ERROR {
		fileName, funcName, line := getPrefix(depth)
		f.logChan <- fmt.Sprintf("[%v:%v:%v]", fileName, funcName, line) + fmt.Sprintf("\033[1;31m[ERROR]"+format+"\033[0m", v...)
		f.logChan <- fmt.Sprintf("%s", debug.Stack())
	}
}

// same with Error()
func (f *FileLogger) E(depth int, format string, v ...interface{}) {
	f.Error(depth, format, v...)
}
