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
 * File name: log.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: log.go 4264 2020-08-19 21:59:19Z nimish $
 */

package tgdb

import (
	"bytes"
	"io"
	"log"
)

type LogLevel int

type TGLogger interface {
	// Trace logs Trace Statements
	Trace(formattedLogMsg string)
	// Debug logs Debug Statements
	Debug(formattedLogMsg string)
	// Info logs Info Statements
	Info(formattedLogMsg string)
	// Warning logs Warning Statements
	Warning(formattedLogMsg string)
	// Error logs Error Statements
	Error(formattedLogMsg string)
	// Fatal logs Fatal Statements
	Fatal(formattedLogMsg string)
	// Log is a generic function that introspects the log level configuration set for the current session, and
	// appropriately calls the level specific logging functions
	Log(logMsg string)

	// GetLogLevel gets Log Level

	GetLogLevel() LogLevel
	// SetLogLevel sets Log Level
	SetLogLevel(level LogLevel)
	// GetLogPrefix gets Log Prefix
	GetLogPrefix() string
	// SetLogPrefix sets Log Prefix
	SetLogPrefix(prefix string)
	// SetLogWriter sets Log Writer
	SetLogWriter(writer io.Writer)
	// SetLogFormatFlags sets Log Format
	SetLogFormatFlags(flags int)

	SetLogBaseDir(logBaseDr string) error
	GetLogBaseDir() string
	SetFileNameBase(baseFileName string) error
	GetFileNameBase() string
	SetFileCount (count int)
	GetFileCount () int
	SetFileSize (size int)
	GetFileSize () int
	// Additional Method to help debugging
	String() string

	IsDebug() bool
	IsInfo() bool
	IsTrace() bool
	IsWarning() bool

}


const (
	DefaultLogPrefix = "TGDB-GOAPI-Logger: "
	DefaultLogLevel = WarningLog
	//DefaultLogFlags = log.Ldate|log.Lmicroseconds|log.Lshortfile|log.LUTC
	DefaultLogFlags  = log.Ldate | log.Lmicroseconds | log.LUTC
	DefaultCallDepth = 2
)

// ======= Log Levels Configured =======

const (
	TraceLog   LogLevel = 1 << iota
	DebugLog            // 2
	InfoLog             // 4
	WarningLog          // 8
	ErrorLog            // 16
	FatalLog            // 32
	NoLog               // 64
)


// Log4J:: A log request of level p in a logger with level q is enabled if p >= q.
// This rule is at the heart of log4j. It assumes that levels are ordered.
// For the standard levels, we need ALL < DEBUG < INFO < WARN < ERROR < FATAL < OFF


func (logLevel LogLevel) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if logLevel&TraceLog == TraceLog {
		buffer.WriteString("LogLevel=TRACE/DEBUG-WIRE")
	} else if logLevel&DebugLog == DebugLog {
		buffer.WriteString("LogLevel=DEBUG")
	} else if logLevel&InfoLog == InfoLog {
		buffer.WriteString("LogLevel=INFO")
	} else if logLevel&WarningLog == WarningLog {
		buffer.WriteString("LogLevel=WARNING")
	} else if logLevel&ErrorLog == ErrorLog {
		buffer.WriteString("LogLevel=ERROR")
	} else if logLevel&FatalLog == FatalLog {
		buffer.WriteString("LogLevel=FATAL/PANIC")
	} else if logLevel&NoLog == NoLog {
		buffer.WriteString("LogLevel=OFF")
	}
	return buffer.String()
}

