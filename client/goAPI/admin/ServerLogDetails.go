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
 * File name: ServerLogDetails.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// ======= Server Log Level Types =======
type TGLogLevel int

const (
	TGLL_Console TGLogLevel = -2
	TGLL_Invalid TGLogLevel = -1
	TGLL_Fatal   TGLogLevel = iota
	TGLL_Error
	TGLL_Warn
	TGLL_Info
	TGLL_User
	TGLL_Debug
	TGLL_DebugFine
	TGLL_DebugFiner
	TGLL_MaxLogLevel
)

// ======= Server Log Component Types =======
type TGLogComponent int64

const (
	TGLC_COMMON_COREMEMORY TGLogComponent = iota
	TGLC_COMMON_CORECOLLECTIONS
	TGLC_COMMON_COREPLATFORM
	TGLC_COMMON_CORESTRING
	TGLC_COMMON_UTILS
	TGLC_COMMON_GRAPH
	TGLC_COMMON_MODEL
	TGLC_COMMON_NET
	TGLC_COMMON_PDU
	TGLC_COMMON_SEC
	TGLC_COMMON_FILES
	TGLC_COMMON_RESV2

	//Server Components
	TGLC_SERVER_CDMP
	TGLC_SERVER_DB
	TGLC_SERVER_EXPIMP
	TGLC_SERVER_INDEX
	TGLC_SERVER_INDEXBTREE
	TGLC_SERVER_INDEXISAM
	TGLC_SERVER_QUERY
	TGLC_SERVER_QUERY_RESV1
	TGLC_SERVER_QUERY_RESV2
	TGLC_SERVER_TXN
	TGLC_SERVER_TXNLOG
	TGLC_SERVER_TXNWRITER
	TGLC_SERVER_STORAGE
	TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_SERVER_GRAPH
	TGLC_SERVER_MAIN
	TGLC_SERVER_RESV2
	TGLC_SERVER_RESV3
	TGLC_SERVER_RESV4

	//Security Components
	TGLC_SECURITY_DATA
	TGLC_SECURITY_NET
	TGLC_SECURITY_RESV1
	TGLC_SECURITY_RESV2

	TGLC_ADMIN_LANG
	TGLC_ADMIN_CMD
	TGLC_ADMIN_MAIN
	TGLC_ADMIN_AST
	TGLC_ADMIN_GREMLIN

	TGLC_CUDA_GRAPHMGR
	TGLC_CUDA_KERNELEXECUTIVE
	TGLC_CUDA_RESV1
)

const (
	// User Defined Components
	TGLC_LOG_COREALL     TGLogComponent = TGLC_COMMON_COREMEMORY | TGLC_COMMON_CORECOLLECTIONS | TGLC_COMMON_COREPLATFORM | TGLC_COMMON_CORESTRING
	TGLC_LOG_GRAPHALL    TGLogComponent = TGLC_COMMON_GRAPH | TGLC_SERVER_GRAPH
	TGLC_LOG_MODEL       TGLogComponent = TGLC_COMMON_MODEL
	TGLC_LOG_NET         TGLogComponent = TGLC_COMMON_NET
	TGLC_LOG_PDUALL      TGLogComponent = TGLC_COMMON_PDU | TGLC_SERVER_CDMP
	TGLC_LOG_SECALL      TGLogComponent = TGLC_COMMON_SEC | TGLC_SECURITY_DATA | TGLC_SECURITY_NET
	TGLC_LOG_CUDAALL     TGLogComponent = TGLC_LOG_GRAPHALL | TGLC_CUDA_GRAPHMGR | TGLC_CUDA_KERNELEXECUTIVE
	TGLC_LOG_TXNALL      TGLogComponent = TGLC_SERVER_TXN | TGLC_SERVER_TXNLOG | TGLC_SERVER_TXNWRITER
	TGLC_LOG_STORAGEALL  TGLogComponent = TGLC_SERVER_STORAGE | TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_LOG_PAGEMANAGER TGLogComponent = TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_LOG_ADMINALL    TGLogComponent = TGLC_ADMIN_LANG | TGLC_ADMIN_CMD | TGLC_ADMIN_MAIN | TGLC_ADMIN_AST | TGLC_ADMIN_GREMLIN
	TGLC_LOG_MAIN        TGLogComponent = TGLC_SERVER_MAIN | TGLC_ADMIN_MAIN
)

const (
	TGLC_LOG_GLOBAL TGLogComponent = 0xFFFFFFFFFFFFFFF
)

type ServerLogDetails struct {
	logLevel     TGLogLevel
	logComponent TGLogComponent
}

func DefaultServerLogDetails() *ServerLogDetails {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerLogDetails{})

	return &ServerLogDetails{}
}

func NewServerLogDetails(logLevel TGLogLevel, _logComponent TGLogComponent) *ServerLogDetails {
	newServerLogDetails := DefaultServerLogDetails()
	newServerLogDetails.logLevel = logLevel
	newServerLogDetails.logComponent = _logComponent
	return newServerLogDetails
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerLogDetails
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerLogDetails:{")
	buffer.WriteString(fmt.Sprintf("LogLevel: '%+v'", obj.logLevel))
	buffer.WriteString(fmt.Sprintf(", LogComponent: '%+v'", obj.logComponent))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> ServerLogDetails
/////////////////////////////////////////////////////////////////

// GetLogLevel returns the server log level
func (obj *ServerLogDetails) GetLogLevel() TGLogLevel {
	return obj.logLevel
}

// GetLogComponent returns the server log component
func (obj *ServerLogDetails) GetLogComponent() TGLogComponent {
	return obj.logComponent
}

// SetLogLevel sets the server log level
func (obj *ServerLogDetails) SetLogLevel(level TGLogLevel) {
	obj.logLevel = level
}

// SetLogComponent sets the server log component
func (obj *ServerLogDetails) SetLogComponent(comp TGLogComponent) {
	obj.logComponent = comp
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.logLevel, obj.logComponent)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerLogDetails:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.logLevel, &obj.logComponent)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerLogDetails:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
