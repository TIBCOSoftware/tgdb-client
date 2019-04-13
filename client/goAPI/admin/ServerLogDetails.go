package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
 * File name: ServerLogDetails.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
	/**
		TGLC_COMMON_COREMEMORY (1L << 0),
		TGLC_COMMON_CORECOLLECTIONS (1L << 1),
		TGLC_COMMON_COREPLATFORM (1L << 2),
		TGLC_COMMON_CORESTRING (1L << 3),
		TGLC_COMMON_UTILS (1L << 4),
		TGLC_COMMON_GRAPH (1L << 5),
		TGLC_COMMON_MODEL (1L << 6),
		TGLC_COMMON_NET (1L << 7),
		TGLC_COMMON_PDU (1L << 8),
		TGLC_COMMON_SEC (1L << 9),
		TGLC_COMMON_FILES (1L << 10),
		TGLC_COMMON_RESV2 (1L << 11),

		//Server Components
		TGLC_SERVER_CDMP (1L << 12),
		TGLC_SERVER_DB (1L << 13),
		TGLC_SERVER_EXPIMP (1L << 14),
		TGLC_SERVER_INDEX (1L << 15),
		TGLC_SERVER_INDEXBTREE (1L << 16),
		TGLC_SERVER_INDEXISAM (1L << 17),
		TGLC_SERVER_QUERY (1L << 18),
		TGLC_SERVER_QUERY_RESV1 (1L << 19),
		TGLC_SERVER_QUERY_RESV2 (1L << 20),
		TGLC_SERVER_TXN (1L << 21),
		TGLC_SERVER_TXNLOG (1L << 22),
		TGLC_SERVER_TXNWRITER (1L << 23),
		TGLC_SERVER_STORAGE (1L << 24),
		TGLC_SERVER_STORAGEPAGEMANAGER (1L << 25),
		TGLC_SERVER_GRAPH (1L << 26),
		TGLC_SERVER_MAIN (1L << 27),
		TGLC_SERVER_RESV2 (1L << 28),
		TGLC_SERVER_RESV3 (1L << 29),
		TGLC_SERVER_RESV4 (1L << 30),

		//Security Components
		TGLC_SECURITY_DATA (1L << 31),
		TGLC_SECURITY_NET (1L << 32),
		TGLC_SECURITY_RESV1 (1L << 33),
		TGLC_SECURITY_RESV2 (1L << 34),

		TGLC_ADMIN_LANG (1L << 35),
		TGLC_ADMIN_CMD (1L << 36),
		TGLC_ADMIN_MAIN (1L << 37),
		TGLC_ADMIN_AST (1L << 38),
		TGLC_ADMIN_GREMLIN (1L << 39),

		TGLC_CUDA_GRAPHMGR (1L << 40),
		TGLC_CUDA_KERNELEXECUTIVE (1L << 41),
		TGLC_CUDA_RESV1 (1L << 42),

		TGLC_LOG_GLOBAL (0xFFFFFFFFFFFFFFFFL),

		// User Defined Components
		TGLC_LOG_COREALL (TGLC_COMMON_COREMEMORY.getLogComponent() | TGLC_COMMON_CORECOLLECTIONS.getLogComponent() | TGLC_COMMON_COREPLATFORM.getLogComponent() | TGLC_COMMON_CORESTRING.getLogComponent()),

		TGLC_LOG_GRAPHALL (TGLC_COMMON_GRAPH.getLogComponent() | TGLC_SERVER_GRAPH.getLogComponent()),
		TGLC_LOG_MODEL (TGLC_COMMON_MODEL.getLogComponent()),
		TGLC_LOG_NET (TGLC_COMMON_NET.getLogComponent()),
		TGLC_LOG_PDUALL (TGLC_COMMON_PDU.getLogComponent() | TGLC_SERVER_CDMP.getLogComponent()),
		TGLC_LOG_SECALL (TGLC_COMMON_SEC.getLogComponent() | TGLC_SECURITY_DATA.getLogComponent() | TGLC_SECURITY_NET.getLogComponent()),
		TGLC_LOG_CUDAALL (TGLC_LOG_GRAPHALL.getLogComponent() | TGLC_CUDA_GRAPHMGR.getLogComponent() | TGLC_CUDA_KERNELEXECUTIVE.getLogComponent()),
		TGLC_LOG_TXNALL (TGLC_SERVER_TXN.getLogComponent() | TGLC_SERVER_TXNLOG.getLogComponent() | TGLC_SERVER_TXNWRITER.getLogComponent()),
		TGLC_LOG_STORAGEALL (TGLC_SERVER_STORAGE.getLogComponent() | TGLC_SERVER_STORAGEPAGEMANAGER.getLogComponent()),
		TGLC_LOG_PAGEMANAGER (TGLC_SERVER_STORAGEPAGEMANAGER.getLogComponent()),
		TGLC_LOG_ADMINALL (TGLC_ADMIN_LANG.getLogComponent() | TGLC_ADMIN_CMD.getLogComponent() | TGLC_ADMIN_MAIN.getLogComponent() | TGLC_ADMIN_AST.getLogComponent() | TGLC_ADMIN_GREMLIN.getLogComponent()),
		TGLC_LOG_MAIN (TGLC_SERVER_MAIN.getLogComponent() | TGLC_ADMIN_MAIN.getLogComponent());
	 */
	TGLC_COMMON_COREMEMORY TGLogComponent = -2
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
