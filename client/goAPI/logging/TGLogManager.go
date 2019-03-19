package logging

import (
	"encoding/gob"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"io"
	"sync"
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

type TGLogManager struct {
	logger types.TGLogger
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

// GetLogger gets the logger instance handle
func (m *TGLogManager) GetLogger() types.TGLogger {
	return m.logger
}

// SetLogger sets the logger instance handle
func (m *TGLogManager) SetLogger(logger types.TGLogger) {
	m.logger = logger
}
