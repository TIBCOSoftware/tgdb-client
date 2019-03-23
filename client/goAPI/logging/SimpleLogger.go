package logging

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"syscall"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: SimpleLogger.go
 * Created on: Dec 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// This is pretty basic implementation of Logger that uses bare-bone functionality provided in out-of-the-box GO log
// package. This can easily be extended and/or another implementation based on TGLogger interface can be developed.

/** The following sample was tested on GO Playground
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Primitive Log Implementation # 1

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	LogFlags int
)

func Init(traceHandle, infoHandle, warningHandle, errorHandle io.Writer) {
	LogFlags := log.Ldate | log.Lmicroseconds | log.Llongfile | log.LUTC

	//lvl, ok := os.LookupEnv("LOG_LEVEL")
	//// LOG_LEVEL not set, let's default to debug
	//if !ok {
	//    lvl = "debug"
	//}

	Trace = log.New(traceHandle, "TRACE: ", LogFlags)
	Info = log.New(infoHandle, "INFO: ", LogFlags)
	Warning = log.New(warningHandle, "WARNING: ", LogFlags)
	Error = log.New(errorHandle, "ERROR: ", LogFlags)
}

// Primitive Log Implementation # 2

// Log4J:: A log request of level p in a logger with level q is enabled if p >= q.
// This rule is at the heart of log4j. It assumes that levels are ordered.
// For the standard levels, we need ALL < DEBUG < INFO < WARN < ERROR < FATAL < OFF

type logLevel int

const (
	traceLog   logLevel = 1 << iota
	debugLog            // 2
	infoLog             // 4
	warningLog          // 8
	errorLog            // 16
	fatalLog            // 32
)

type Logger struct {
	level logLevel
	log   *log.Logger
}

// Level, Prefix and LogWriter can be read and supplied from outside as ENV variables or Command Line Parameters
func NewLogger(level int, logPrefix string, logWriter io.Writer) *Logger {
	ml := Logger{}
	ml.level = logLevel(level)
	ml.log = log.New(logWriter, logPrefix, log.Ldate|log.Lmicroseconds|log.Llongfile|log.LUTC)
	return &ml
}

func (m *Logger) Trace(args ...interface{}) {
	if m.level <= traceLog {
		m.log.Print(args...)
	}
}

func (m *Logger) Debug(args ...interface{}) {
	if m.level <= debugLog {
		m.log.Print(args...)
	}
}

func (m *Logger) Info(args ...interface{}) {
	if m.level <= infoLog {
		m.log.Print(args...)
	}
}

func (m *Logger) Warning(args ...interface{}) {
	if m.level <= warningLog {
		m.log.Print(args...)
	}
}

func (m *Logger) Error(args ...interface{}) {
	if m.level <= errorLog {
		m.log.Print(args...)
	}
}

func (m *Logger) Fatal(args ...interface{}) {
	if m.level <= fatalLog {
		m.log.Fatal(args...)
	}
}

func main() {
	fmt.Println("Hello, playground")

	// Primitive Log Implementation # 1
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Trace.Println("I have something standard to say")
	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Something has failed")

	// Primitive Log Implementation # 2
	log := NewLogger(8, "TGDB-GOAPI-Logger: ", os.Stdout)
	log.Trace("I have something standard to say")
	log.Debug("Something to debug")
	log.Info("Special Information")
	log.Warning("There is something you need to know about")
	log.Error("Something has failed")
	log.Fatal("Something is seriously out-of-whack!!!")
}
*/

type Logger struct {
	level types.LogLevel
	log   *log.Logger
}

func defaultLogger() *Logger {
	ml := Logger{
		level: types.DefaultLogLevel,
		//level: types.InfoLog,	// TODO: For Debugging Purpose Only - Delete It
		log: log.New(os.Stdout, types.DefaultLogPrefix, types.DefaultLogFlags),
	}
	return &ml
}

// Level, Prefix and LogWriter can be read and supplied from outside as ENV variables or Command Line Parameters
func NewLogger(level int, logPrefix string, logWriter io.Writer, logFlags int) *Logger {
	ml := defaultLogger()
	ml.level = types.LogLevel(level)
	ml.log = log.New(logWriter, logPrefix, logFlags)
	return ml
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> Logger
/////////////////////////////////////////////////////////////////

func (m *Logger) formatMessage(logMsg string) string {
	// Format log message according to configured msgFormat
	formattedMsg := logMsg
	if m.log.Flags()&(log.Lshortfile|log.Llongfile) == 0 {
		// NOTE: IF the following call to GetFileAndLine needs to move to other location, please adjust call depth
		fileName, lineNo := m.GetFileAndLine(types.DefaultCallDepth+3)
		gRoutineId := syscall.Getgid()
		formattedMsg = fmt.Sprintf("%s:%d [%d] %s", fileName, lineNo, gRoutineId, logMsg)
	}
	return formattedMsg
}

func (m *Logger) simpleLog(logMsg string) {
	switch m.level {
	case types.FatalLog:
		m.Fatal(logMsg)
	case types.ErrorLog:
		m.Error(logMsg)
	case types.WarningLog:
		m.Warning(logMsg)
	case types.InfoLog:
		m.Info(logMsg)
	case types.DebugLog:
		m.Debug(logMsg)
	case types.TraceLog:
		m.Trace(logMsg)
	default:
		m.Debug(logMsg)
	}
}

// GetFileAndLine returns the file and line from the stack at the given call depth
func (m *Logger) GetFileAndLine(callDepth int) (string, int) {
	//var ok bool
	_, file, line, _ := runtime.Caller(callDepth)
	file = path.Base(file)
	return file, line
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGLogger
/////////////////////////////////////////////////////////////////

// Trace logs Trace (Down-to-the-wire) Statements
func (m *Logger) Trace(logMsg string) {
	if m.level >= types.TraceLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Debug logs Debug Statements
func (m *Logger) Debug(logMsg string) {
	if m.level >= types.DebugLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Info logs Informative Statements
func (m *Logger) Info(logMsg string) {
	if m.level >= types.InfoLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Warning logs Warning Statements
func (m *Logger) Warning(logMsg string) {
	if m.level >= types.WarningLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Error logs Error Statements
func (m *Logger) Error(logMsg string) {
	if m.level >= types.ErrorLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Fatal logs Fatal Statements
func (m *Logger) Fatal(logMsg string) {
	if m.level >= types.FatalLog {
		// Format log message according to configured msgFormat
		formattedLogMsg := m.formatMessage(logMsg)
		// Ignore Error Handling
		_ = m.log.Output(types.DefaultCallDepth+3, formattedLogMsg)
	}
}

// Log is a generic function that introspects the log level configuration set for the current session, and
// appropriately calls the level specific logging functions
func (m *Logger) Log(logMsg string) {
	m.simpleLog(logMsg)
}

// GetLogLevel gets Log Level
func (m *Logger) GetLogLevel() types.LogLevel {
	return m.level
}

// SetLogLevel sets Log Level
func (m *Logger) SetLogLevel(level types.LogLevel) {
	m.level = level
}

// GetLogPrefix gets Log Prefix
func (m *Logger) GetLogPrefix() string {
	return m.log.Prefix()
}

// SetLogPrefix sets Log Prefix
func (m *Logger) SetLogPrefix(prefix string) {
	m.log.SetPrefix(prefix)
}

// SetLogWriter sets Log Writer
func (m *Logger) SetLogWriter(writer io.Writer) {
	m.log.SetOutput(writer)
}

// SetLogFormatFlags sets Log Format
func (m *Logger) SetLogFormatFlags(flags int) {
	m.log.SetFlags(flags)
}

// Additional Method to help debugging
func (m *Logger) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("SimpleLogger:{")
	buffer.WriteString(fmt.Sprintf("LogLevel: %+v", m.level.String()))
	buffer.WriteString(fmt.Sprintf(", Log: %+v", m.log))
	buffer.WriteString("}")
	return buffer.String()
}
