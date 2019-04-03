package pdu

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
 * File name: VerbId.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//type VerbId int

const (
	AbstractMessage int = -100
)

const (
	// Ping Message - Heart beats
	VerbPingMessage int = 0
	// HandShake Request/Response protocol
	VerbHandShakeRequest  int = 1
	VerbHandShakeResponse int = 2
	// Authenticate Request/Response protocol
	VerbAuthenticateRequest  int = 3
	VerbAuthenticateResponse int = 4
	// Transaction - begin/commit/rollback protocol verbs
	VerbBeginTransactionRequest     int = 5
	VerbBeginTransactionResponse    int = 6
	VerbCommitTransactionRequest    int = 7
	VerbCommitTransactionResponse   int = 8
	VerbRollbackTransactionRequest  int = 9
	VerbRollbackTransactionResponse int = 10
	// Query Request/Response verbs
	VerbQueryRequest  int = 11
	VerbQueryResponse int = 12
	// Graph Traversal verbs
	VerbTraverseRequest  int = 13
	VerbTraverseResponse int = 14
	// Admin Request/Response verbs
	VerbAdminRequest  int = 15
	VerbAdminResponse int = 16
	// Retrieve meta data
	VerbMetadataRequest  int = 19
	VerbMetadataResponse int = 20
	// Get entities
	VerbGetEntityRequest  int = 21
	VerbGetEntityResponse int = 22
	// Get LargeObject
	VerbGetLargeObjectRequest  int = 23
	VerbGetLargeObjectResponse int = 24
	// Import/Export verbs - They are admin request, and not supported by Java
	VerbBeginExportRequest    int = 25
	VerbBeginExportResponse   int = 26
	VerbPartialExportRequest  int = 27
	VerbPartialExportResponse int = 28
	VerbCancelExportRequest   int = 29
	VerbBeginImportRequest    int = 31
	BeginImportResponse       int = 32
	VerbPartialImportRequest  int = 33
	VerbPartialImportResponse int = 34
	// Dump Stacktrace request verb
	VerbDumpStacktraceRequest int = 39
	// Disconnect Request verbs
	VerbDisconnectChannelRequest    int = 40
	VerbSessionForcefullyTerminated int = 41
	// Decryption Request verbs
	VerbDecryptBufferRequest  int = 44
	VerbDecryptBufferResponse int = 45
	// Unknown Exception Message on the server.
	VerbExceptionMessage int = 100
	VerbInvalidMessage   int = -1
)

type CommandVerbs struct {
	id          int
	name        string
	implementor string
}

var PreDefinedVerbs = map[int]CommandVerbs{
	VerbPingMessage:                 {id: VerbPingMessage, name: "VerbPingMessage", implementor: "pdu.VerbPingMessage"},
	VerbHandShakeRequest:            {id: VerbHandShakeRequest, name: "VerbHandShakeRequest", implementor: "pdu.HandshakeRequest"},
	VerbHandShakeResponse:           {id: VerbHandShakeResponse, name: "VerbHandShakeResponse", implementor: "pdu.HandshakeResponse"},
	VerbAuthenticateRequest:         {id: VerbAuthenticateRequest, name: "VerbAuthenticateRequest", implementor: "pdu.VerbAuthenticateRequest"},
	VerbAuthenticateResponse:        {id: VerbAuthenticateResponse, name: "VerbAuthenticateResponse", implementor: "pdu.VerbAuthenticateResponse"},
	VerbBeginTransactionRequest:     {id: VerbBeginTransactionRequest, name: "VerbBeginTransactionRequest", implementor: "pdu.VerbBeginTransactionRequest"},
	VerbBeginTransactionResponse:    {id: VerbBeginTransactionResponse, name: "VerbBeginTransactionResponse", implementor: "pdu.VerbBeginTransactionResponse"},
	VerbCommitTransactionRequest:    {id: VerbCommitTransactionRequest, name: "VerbCommitTransactionRequest", implementor: "pdu.VerbCommitTransactionRequest"},
	VerbCommitTransactionResponse:   {id: VerbCommitTransactionResponse, name: "VerbCommitTransactionResponse", implementor: "pdu.VerbCommitTransactionResponse"},
	VerbRollbackTransactionRequest:  {id: VerbRollbackTransactionRequest, name: "VerbRollbackTransactionRequest", implementor: "pdu.VerbRollbackTransactionRequest"},
	VerbRollbackTransactionResponse: {id: VerbRollbackTransactionResponse, name: "VerbRollbackTransactionResponse", implementor: "pdu.VerbRollbackTransactionResponse"},
	VerbQueryRequest:                {id: VerbQueryRequest, name: "VerbQueryRequest", implementor: "pdu.VerbQueryRequest"},
	VerbQueryResponse:               {id: VerbQueryResponse, name: "VerbQueryResponse", implementor: "pdu.VerbQueryResponse"},
	VerbTraverseRequest:             {id: VerbTraverseRequest, name: "VerbTraverseRequest", implementor: "pdu.VerbTraverseRequest"},
	VerbTraverseResponse:            {id: VerbTraverseResponse, name: "VerbTraverseResponse", implementor: "pdu.VerbTraverseResponse"},
	VerbAdminRequest:                {id: VerbAdminRequest, name: "VerbAdminRequest", implementor: "pdu.VerbAdminRequest"},
	VerbAdminResponse:               {id: VerbAdminResponse, name: "VerbAdminResponse", implementor: "pdu.VerbAdminResponse"},
	VerbMetadataRequest:             {id: VerbMetadataRequest, name: "VerbMetadataRequest", implementor: "pdu.VerbMetadataRequest"},
	VerbMetadataResponse:            {id: VerbMetadataResponse, name: "VerbMetadataResponse", implementor: "pdu.VerbMetadataResponse"},
	VerbGetEntityRequest:            {id: VerbGetEntityRequest, name: "VerbGetEntityRequest", implementor: "pdu.VerbGetEntityRequest"},    //0 = mean immediate, AttributeTypeInteger Max for indefinite
	VerbGetEntityResponse:           {id: VerbGetEntityResponse, name: "VerbGetEntityResponse", implementor: "pdu.VerbGetEntityResponse"}, //Represented in ms. Default Value is 10sec
	VerbGetLargeObjectRequest:       {id: VerbGetLargeObjectRequest, name: "VerbGetLargeObjectRequest", implementor: "pdu.VerbGetLargeObjectRequest"},
	VerbGetLargeObjectResponse:      {id: VerbGetLargeObjectResponse, name: "VerbGetLargeObjectResponse", implementor: "pdu.VerbGetLargeObjectResponse"},
	VerbDumpStacktraceRequest:       {id: VerbDumpStacktraceRequest, name: "VerbDumpStacktraceRequest", implementor: "pdu.VerbDumpStacktraceRequest"},
	VerbDisconnectChannelRequest:    {id: VerbDisconnectChannelRequest, name: "VerbDisconnectChannelRequest", implementor: "pdu.VerbDisconnectChannelRequest"},
	VerbSessionForcefullyTerminated: {id: VerbSessionForcefullyTerminated, name: "VerbSessionForcefullyTerminated", implementor: "pdu.VerbSessionForcefullyTerminated"},
	VerbDecryptBufferRequest:        {id: VerbDecryptBufferRequest, name: "VerbDecryptBufferRequest", implementor: "pdu.VerbDecryptBufferRequest"},
	VerbDecryptBufferResponse:       {id: VerbDecryptBufferResponse, name: "VerbDecryptBufferResponse", implementor: "pdu.VerbDecryptBufferResponse"},
	VerbExceptionMessage:            {id: VerbExceptionMessage, name: "VerbExceptionMessage", implementor: "pdu.VerbExceptionMessage"},
	VerbInvalidMessage:              {id: VerbInvalidMessage, name: "VerbInvalidMessage", implementor: "pdu.VerbInvalidMessage"},
}

func NewVerbId(id int, name, impl string) *CommandVerbs {
	newConfig := &CommandVerbs{id: id, name: name, implementor: impl}
	return newConfig
}

// Return the commandVerbs given its id
// @param id VerbId
// @return commandVerbs associated to the id
func GetVerb(id int) *CommandVerbs {
	verb, ok := PreDefinedVerbs[id]
	if ok {
		return &verb
	} else {
		invalid := PreDefinedVerbs[VerbInvalidMessage]
		return &invalid
	}
}

/////////////////////////////////////////////
// Helper Public functions for CommonVerbs //
/////////////////////////////////////////////

func (obj *CommandVerbs) GetID() int {
	return obj.id
}

func (obj *CommandVerbs) GetName() string {
	return obj.name
}

func (obj *CommandVerbs) GetImplementor() string {
	return obj.implementor
}
