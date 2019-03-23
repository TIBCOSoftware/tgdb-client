package types

import (
	"bytes"
	"io"
	"log"
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
 * File name: TGLogger.go
 * Created on: Dec 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// Log4J:: A log request of level p in a logger with level q is enabled if p >= q.
// This rule is at the heart of log4j. It assumes that levels are ordered.
// For the standard levels, we need ALL < DEBUG < INFO < WARN < ERROR < FATAL < OFF

// ======= Log Levels Configured =======
type LogLevel int

const (
	TraceLog   LogLevel = 1 << iota
	DebugLog            // 2
	InfoLog             // 4
	WarningLog          // 8
	ErrorLog            // 16
	FatalLog            // 32
	NoLog               // 64
)

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

const (
	DefaultLogPrefix = "TGDB-GOAPI-Logger: "
	DefaultLogLevel = DebugLog
	//DefaultLogFlags = log.Ldate|log.Lmicroseconds|log.Lshortfile|log.LUTC
	DefaultLogFlags  = log.Ldate | log.Lmicroseconds | log.LUTC
	DefaultCallDepth = 2
)

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
	// Additional Method to help debugging
	String() string
}
