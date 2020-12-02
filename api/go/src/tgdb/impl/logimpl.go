/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File Name: logimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: logimpl.go 4264 2020-08-19 21:59:19Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"tgdb"
	"time"
)


type Logger struct {
	depth int
	level tgdb.LogLevel
	log   *log.Logger
	size int
	count int
	baseDir string
	fileNameBase string
	currentIndex int
	mu     sync.Mutex
}

func defaultLogger() *Logger {
	ml := Logger{
		depth: tgdb.DefaultCallDepth,
		level: tgdb.DefaultLogLevel,
		//level: types.InfoLog,	// TODO: For Debugging Purpose Only - Delete It
		log: log.New(os.Stdout, tgdb.DefaultLogPrefix, tgdb.DefaultLogFlags),
	}
	return &ml
}

/*
// Level, Prefix and LogWriter can be read and supplied from outside as ENV variables or Command Line Parameters
func NewLogger(level int, logPrefix string, logWriter io.Writer, logFlags int) *Logger {
	ml := defaultLogger()
	ml.level = LogLevel(level)
	ml.log = log.New(logWriter, logPrefix, logFlags)
	return ml
}
*/


/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> Logger
/////////////////////////////////////////////////////////////////

func (m *Logger) formatMessage(callDepth int, logMsg string) string {
	// Format log message according to configured msgFormat
	formattedMsg := logMsg
	if m.log.Flags()&(log.Lshortfile|log.Llongfile) == 0 {
		// NOTE: IF the following call to GetFileAndLine needs to move to other location, please adjust call depth
		fileName, lineNo := m.GetFileAndLine(callDepth)
		gRoutineId := syscall.Getgid()
		//formattedMsg = fmt.Sprintf("%s:%d [%d] %s", fileName, lineNo, gRoutineId, logMsg)
		formattedMsg = fmt.Sprintf("%d %s:%d - %s", gRoutineId, fileName, lineNo, logMsg)
	}
	return formattedMsg
}

func (m *Logger) simpleLog(logMsg string) {
	m.depth += 2
	switch m.level {
	case tgdb.FatalLog:
		m.Fatal(logMsg)
	case tgdb.ErrorLog:
		m.Error(logMsg)
	case tgdb.WarningLog:
		//m.Warning(logMsg)
	case tgdb.InfoLog:
		//m.Info(logMsg)
	case tgdb.DebugLog:
		//m.Debug(logMsg)
	case tgdb.TraceLog:
		//m.Trace(logMsg)
	default:
		//m.Debug(logMsg)
	}
	m.depth -= 2
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

func (m *Logger) IsTrace() bool {
	return m.level <= tgdb.TraceLog
}

// Trace logs Trace (Down-to-the-wire) Statements
func (m *Logger) Trace(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.TraceLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Trace " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}

func (m *Logger) IsDebug() bool {
	return m.level <= tgdb.DebugLog
}

// Debug logs Debug Statements
func (m *Logger) Debug(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.DebugLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Debug " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}

func (m *Logger) IsInfo() bool {
	return m.level <= tgdb.InfoLog
}


// Info logs Informative Statements
func (m *Logger) Info(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.InfoLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Info " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}

func (m *Logger) IsWarning() bool {
	return m.level <= tgdb.WarningLog
}


// Warning logs Warning Statements
func (m *Logger) Warning(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.WarningLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Warning " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}


func (m *Logger) IsError() bool {
	return m.level <= tgdb.ErrorLog
}

// Error logs Error Statements
func (m *Logger) Error(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.ErrorLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Error " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}

// Fatal logs Fatal Statements
func (m *Logger) Fatal(logMsg string) {
	callDepth := m.depth+1
	if m.level <= tgdb.FatalLog {
		formattedLogMsg := m.formatMessage(callDepth, logMsg)
		formattedLogMsg = " Fatal " + formattedLogMsg
		m.formattedMsgOutput(formattedLogMsg)
	}
}

func (m *Logger) formattedMsgOutput (formattedLogMsg string) {
	formattedLogMsg = m.formatTimeString() + formattedLogMsg
	msgLength := len(formattedLogMsg)
	if  msgLength == 0 || formattedLogMsg[msgLength-1] != '\n' {
		formattedLogMsg += "\n"
	}
	msgBytes := []byte(formattedLogMsg)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.log.Writer().Write(msgBytes)
	file := m.log.Writer().(*os.File)
	if *file != *os.Stdout {
		m.CheckAndUpdateFileHandle()
	}
}

func (m *Logger) formatTimeString () string {
	current_time := time.Now()
	var buff string
	buff = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%03d-00:00",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second(), current_time.Nanosecond()/1000000)
	return buff
}

// Log is a generic function that introspects the log level configuration set for the current session, and
// appropriately calls the level specific logging functions
func (m *Logger) Log(logMsg string) {
	m.simpleLog(logMsg)
}

// GetLogLevel gets Log Level
func (m *Logger) GetLogLevel() tgdb.LogLevel {
	return m.level
}

// SetLogLevel sets Log Level
func (m *Logger) SetLogLevel(level tgdb.LogLevel) {
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

func (m *Logger) SetLogBaseDir(logBaseDr string) error {
	m.baseDir = logBaseDr

	if _, err := os.Stat(m.baseDir); os.IsNotExist(err) {
		err = os.Mkdir(logBaseDr, 0777)
		if err != nil {
			fmt.Println ("Error: " + err.Error())
			return err
		}
	}
	return nil
}

func (m *Logger) GetLogBaseDir() (string) {
	return m.baseDir
}

func (m *Logger) SetFileNameBase(baseFileName string) error {
	m.fileNameBase = baseFileName
	return m.UpdateFileHandle()
}

func (m *Logger) UpdateFileHandle() error {
	logFileHandle, e := os.OpenFile(m.GetAbsoluteFileName(), os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0777)
	if e != nil {
		fmt.Println("Error: " + e.Error())
		return e
	}
	m.SetLogWriter(logFileHandle)
	return nil
}

func (m *Logger) CheckAndUpdateFileHandle () error {
	absFileName := m.GetAbsoluteFileName()
	fileInfo, err := os.Stat(absFileName)
	if err != nil {
		return err
	}
	if fileInfo.Size() < int64 (m.size) {
		return nil
	} else {
		file := m.log.Writer().(*os.File)
		file.Close()
		m.currentIndex = m.currentIndex + 1
		if m.currentIndex >= m.count {
			m.currentIndex = 0
		}
		m.removeIfExists ()
		return m.UpdateFileHandle()
	}
}

func (m *Logger) removeIfExists() {
	_, err := os.Stat(m.GetAbsoluteFileName())
	if err == nil {
		os.Remove(m.GetAbsoluteFileName())
	}
}

func (m *Logger) GetAbsoluteFileName () string {
	return m.baseDir + string(os.PathSeparator) + m.fileNameBase + ".log." + strconv.Itoa(m.currentIndex)
}

func (m *Logger) GetFileNameBase() string {
	return m.fileNameBase
}

func (m *Logger) SetFileCount (count int) {
	m.count = count
}

func (m *Logger) GetFileCount () int {
	return m.count
}

func (m *Logger) SetFileSize (size int) {
	m.size = size
}

func (m *Logger) GetFileSize () int {
	return m.size
}


// Additional Method to help debugging
func (m *Logger) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("SimpleLogger:{")
	//buffer.WriteString(fmt.Sprintf("LogLevel: %+v", m.level.String()))
	buffer.WriteString(fmt.Sprintf(", Log: %+v", m.log))
	buffer.WriteString("}")
	return buffer.String()
}



//////////////////////////////////////////////
// TGLogManager
//////////////////////////////////////////////

type TGLogManager struct {
	logger tgdb.TGLogger
}

var globalLogManager *TGLogManager
var lOnce sync.Once

// DefaultTGLogManager will be used by all TGDB GO API implementations internally.
// It uses
//     DefaultLogPrefix = "TGDB-GOAPI-Logger: "
// and DefaultLogFlags = log.Ldate|log.Lmicroseconds|log.Lshortfile|log.LUTC ̰
func DefaultTGLogManager() *TGLogManager {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGLogManager{})

	lOnce.Do(func() {
		globalLogManager = &TGLogManager{
			logger: defaultLogger(),
		}
	})
	return globalLogManager
}

/*
// Not used for now


// NewTGLogManager facilitates a provision for a new Logger w/ new Destination (e.g. File, Pipe or Database) that
// may have customized log levels and prefix to differentiate itself.
func NewTGLogManager(level int, logPrefix string, logWriter io.Writer, logFlags int) *TGLogManager {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGLogManager{})

	lOnce.Do(func() {
		globalLogManager = &TGLogManager{
			logger: NewLogger(level, logPrefix, logWriter, logFlags),
		}
	})
	return globalLogManager
}
*/
// GetLogger gets the logger instance handle
func (m *TGLogManager) GetLogger() tgdb.TGLogger {
	return m.logger
}

// SetLogger sets the logger instance handle
func (m *TGLogManager) SetLogger(logger tgdb.TGLogger) {
	m.logger = logger
}



