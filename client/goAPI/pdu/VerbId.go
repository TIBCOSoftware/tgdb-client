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
	// Unknown Exception Message on the server.
	VerbExceptionMessage int = 100
	VerbInvalidMessage   int = -1
)

type CommandVerbs struct {
	Id          int
	Name        string
	Implementor string
}

var PreDefinedVerbs = map[int]CommandVerbs{
	VerbPingMessage:                 {Id: VerbPingMessage, Name: "VerbPingMessage", Implementor: "pdu.VerbPingMessage"},
	VerbHandShakeRequest:            {Id: VerbHandShakeRequest, Name: "VerbHandShakeRequest", Implementor: "pdu.HandshakeRequest"},
	VerbHandShakeResponse:           {Id: VerbHandShakeResponse, Name: "VerbHandShakeResponse", Implementor: "pdu.HandshakeResponse"},
	VerbAuthenticateRequest:         {Id: VerbAuthenticateRequest, Name: "VerbAuthenticateRequest", Implementor: "pdu.VerbAuthenticateRequest"},
	VerbAuthenticateResponse:        {Id: VerbAuthenticateResponse, Name: "VerbAuthenticateResponse", Implementor: "pdu.VerbAuthenticateResponse"},
	VerbBeginTransactionRequest:     {Id: VerbBeginTransactionRequest, Name: "VerbBeginTransactionRequest", Implementor: "pdu.VerbBeginTransactionRequest"},
	VerbBeginTransactionResponse:    {Id: VerbBeginTransactionResponse, Name: "VerbBeginTransactionResponse", Implementor: "pdu.VerbBeginTransactionResponse"},
	VerbCommitTransactionRequest:    {Id: VerbCommitTransactionRequest, Name: "VerbCommitTransactionRequest", Implementor: "pdu.VerbCommitTransactionRequest"},
	VerbCommitTransactionResponse:   {Id: VerbCommitTransactionResponse, Name: "VerbCommitTransactionResponse", Implementor: "pdu.VerbCommitTransactionResponse"},
	VerbRollbackTransactionRequest:  {Id: VerbRollbackTransactionRequest, Name: "VerbRollbackTransactionRequest", Implementor: "pdu.VerbRollbackTransactionRequest"},
	VerbRollbackTransactionResponse: {Id: VerbRollbackTransactionResponse, Name: "VerbRollbackTransactionResponse", Implementor: "pdu.VerbRollbackTransactionResponse"},
	VerbQueryRequest:                {Id: VerbQueryRequest, Name: "VerbQueryRequest", Implementor: "pdu.VerbQueryRequest"},
	VerbQueryResponse:               {Id: VerbQueryResponse, Name: "VerbQueryResponse", Implementor: "pdu.VerbQueryResponse"},
	VerbTraverseRequest:             {Id: VerbTraverseRequest, Name: "VerbTraverseRequest", Implementor: "pdu.VerbTraverseRequest"},
	VerbTraverseResponse:            {Id: VerbTraverseResponse, Name: "VerbTraverseResponse", Implementor: "pdu.VerbTraverseResponse"},
	VerbAdminRequest:                {Id: VerbAdminRequest, Name: "VerbAdminRequest", Implementor: "pdu.VerbAdminRequest"},
	VerbAdminResponse:               {Id: VerbAdminResponse, Name: "VerbAdminResponse", Implementor: "pdu.VerbAdminResponse"},
	VerbMetadataRequest:             {Id: VerbMetadataRequest, Name: "VerbMetadataRequest", Implementor: "pdu.VerbMetadataRequest"},
	VerbMetadataResponse:            {Id: VerbMetadataResponse, Name: "VerbMetadataResponse", Implementor: "pdu.VerbMetadataResponse"},
	VerbGetEntityRequest:            {Id: VerbGetEntityRequest, Name: "VerbGetEntityRequest", Implementor: "pdu.VerbGetEntityRequest"},    //0 = mean immediate, AttributeTypeInteger Max for indefinite
	VerbGetEntityResponse:           {Id: VerbGetEntityResponse, Name: "VerbGetEntityResponse", Implementor: "pdu.VerbGetEntityResponse"}, //Represented in ms. Default Value is 10sec
	VerbGetLargeObjectRequest:       {Id: VerbGetLargeObjectRequest, Name: "VerbGetLargeObjectRequest", Implementor: "pdu.VerbGetLargeObjectRequest"},
	VerbGetLargeObjectResponse:      {Id: VerbGetLargeObjectResponse, Name: "VerbGetLargeObjectResponse", Implementor: "pdu.VerbGetLargeObjectResponse"},
	VerbDisconnectChannelRequest:    {Id: VerbDisconnectChannelRequest, Name: "VerbDisconnectChannelRequest", Implementor: "pdu.VerbDisconnectChannelRequest"},
	VerbSessionForcefullyTerminated: {Id: VerbSessionForcefullyTerminated, Name: "VerbSessionForcefullyTerminated", Implementor: "pdu.VerbSessionForcefullyTerminated"},
	VerbExceptionMessage:            {Id: VerbExceptionMessage, Name: "VerbExceptionMessage", Implementor: "pdu.VerbExceptionMessage"},
	VerbInvalidMessage:              {Id: VerbInvalidMessage, Name: "VerbInvalidMessage", Implementor: "pdu.VerbInvalidMessage"},
}

func NewVerbId(id int, name, impl string) *CommandVerbs {
	newConfig := &CommandVerbs{Id: id, Name: name, Implementor: impl}
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
