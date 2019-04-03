package pdu

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"math/rand"
	"reflect"
	"testing"
	"time"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: MessageFactory_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func createTestAbstractMessage() *AbstractProtocolMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewAbstractProtocolMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestAuthenticatedMessage() *AuthenticatedMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	connectionId := rand.Int()
	clientId := "Test-Client"
	msg := NewAuthenticatedMessage(authToken, sessionId)
	msg.SetConnectionId(connectionId)
	msg.SetClientId(clientId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestAuthenticateRequestMessage() *AuthenticateRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	clientId := "Test-Client"
	inboxAddr := "Test-Host"
	userName := "Test-User"
	password := []byte("Test-Password")
	msg := NewAuthenticateRequestMessage(authToken, sessionId)
	msg.SetClientId(clientId)
	msg.SetInboxAddr(inboxAddr)
	msg.SetUserName(userName)
	msg.SetPassword(password)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	fmt.Printf("MessageFactory::createTestAuthenticateRequestMessage starting w/ Message BufLength: '%+v'", msg.BufLength)
	return msg
}

func createTestAuthenticateResponseMessage() *AuthenticateResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewAuthenticateResponseMessage(authToken, sessionId)
	msg.SetSuccess(true)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestBeginTransactionRequestMessage() *BeginTransactionRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewBeginTransactionRequestMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestBeginTransactionResponseMessage() *BeginTransactionResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	transactionId := rand.Int63()
	msg := NewBeginTransactionResponseMessage(authToken, sessionId)
	msg.SetTransactionId(transactionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestCommitTransactionRequestMessage() *CommitTransactionRequest {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewCommitTransactionRequestMessage(authToken, sessionId)
	// TODO: Revisit later to put some dummy values in 3 maps and desc array
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestCommitTransactionResponseMessage() *CommitTransactionResponse {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewCommitTransactionResponseMessage(authToken, sessionId)
	// TODO: Revisit later to put some dummy values in graphmetadata, exception and inputstream
	msg.SetAddEntityCount(1)
	msg.SetAddedIdList([]int64{rand.Int63()})
	msg.SetUpdatedEntityCount(1)
	msg.SetUpdatedIdList([]int64{rand.Int63()})
	msg.SetRemovedEntityCount(1)
	msg.SetRemovedIdList([]int64{rand.Int63()})
	msg.SetAttrDescCount(1)
	msg.SetAttrDescId([]int64{rand.Int63()})
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestDisconnectChannelRequestMessage() *DisconnectChannelRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	connectionId := rand.Int()
	clientId := "Test-Client"
	msg := NewDisconnectChannelRequestMessage(authToken, sessionId)
	msg.SetConnectionId(connectionId)
	msg.SetClientId(clientId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestExceptionMessage() *ExceptionMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	exceptionId := rand.Int()
	exceptionMsg := "Test-Exception"
	msg := NewExceptionMessage(authToken, sessionId)
	msg.SetExceptionType(exceptionId)
	msg.SetExceptionMsg(exceptionMsg)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestGetEntityRequestMessage() *GetEntityRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	commandType := rand.Int()
	fetchSize := rand.Int()
	batchSize := rand.Int()
	traversalDepth := rand.Int()
	edgeLimit := rand.Int()
	resultId := rand.Int()
	// TODO: Revisit later to put some dummy values for key
	msg := NewGetEntityRequestMessage(authToken, sessionId)
	msg.SetCommand(int16(commandType))
	msg.SetFetchSize(fetchSize)
	msg.SetBatchSize(batchSize)
	msg.SetTraversalDepth(traversalDepth)
	msg.SetEdgeLimit(edgeLimit)
	msg.SetResultId(resultId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestGetEntityResponseMessage() *GetEntityResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	resultId := rand.Int()
	totalCount := rand.Int()
	resultCount := rand.Int()
	// TODO: Revisit later to put some dummy values for inputstream
	msg := NewGetEntityResponseMessage(authToken, sessionId)
	msg.SetResultId(resultId)
	msg.SetTotalCount(totalCount)
	msg.SetResultCount(resultCount)
	msg.SetHasResult(true)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestGetLargeObjectRequestMessage() *GetLargeObjectRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	entityId := rand.Int63()
	msg := NewGetLargeObjectRequestMessage(authToken, sessionId)
	msg.SetEntityId(entityId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestGetLargeObjectResponseMessage() *GetLargeObjectResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	entityId := rand.Int63()
	// TODO: Revisit later to put some dummy values for bostream
	msg := NewGetLargeObjectResponseMessage(authToken, sessionId)
	msg.SetEntityId(entityId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestHandShakeRequestMessage() *HandShakeRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	challenge := rand.Int63()
	handshakeType := rand.Int()
	msg := NewHandShakeRequestMessage(authToken, sessionId)
	msg.SetChallenge(challenge)
	msg.SetRequestType(handshakeType)
	msg.SetSslMode(false)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestHandShakeResponseMessage() *HandShakeResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	challenge := rand.Int63()
	responseStatus := rand.Int()
	msg := NewHandShakeResponseMessage(authToken, sessionId)
	msg.SetChallenge(challenge)
	msg.SetResponseStatus(responseStatus)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestInvalidMessage() *InvalidMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewInvalidMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestMetadataRequestMessage() *MetadataRequest {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewMetadataRequestMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestMetadataResponseMessage() *MetadataResponse {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	// TODO: Revisit later to put some dummy values for 3 arrays
	msg := NewMetadataResponseMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestPingMessage() *PingMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewPingMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestQueryRequestMessage() *QueryRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	queryExpr := "Test-Query-Expression"
	edgeExpr := "Test-Edge-Expression"
	traverseExpr := "Test-Traversal-Expression"
	endExpr := "Test-End-Expression"
	queryHashId := rand.Int63()
	command := rand.Int()
	fetchSize := rand.Int()
	batchSize := rand.Int()
	traversalDepth := rand.Int()
	edgeLimit := rand.Int()
	sortAttName := "Test-Sort-Attribute"
	sortResultLimit := rand.Int()
	// TODO: Revisit later to put some dummy values for queryObject
	msg := NewQueryRequestMessage(authToken, sessionId)
	msg.SetQuery(queryExpr)
	msg.SetEdgeFilter(edgeExpr)
	msg.SetTraversalCondition(traverseExpr)
	msg.SetEndCondition(endExpr)
	msg.SetQueryHashId(queryHashId)
	msg.SetCommand(command)
	msg.SetFetchSize(fetchSize)
	msg.SetBatchSize(batchSize)
	msg.SetTraversalDepth(traversalDepth)
	msg.SetEdgeLimit(edgeLimit)
	msg.SetSortAttrName(sortAttName)
	msg.SetSortOrderDsc(false)
	msg.SetSortResultLimit(sortResultLimit)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestQueryResponseMessage() *QueryResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	result := rand.Int()
	totalCount := rand.Int()
	resultCount := rand.Int()
	queryHashId := rand.Int63()
	// TODO: Revisit later to put some dummy values for 3 arrays
	msg := NewQueryResponseMessage(authToken, sessionId)
	msg.SetResult(result)
	msg.SetTotalCount(totalCount)
	msg.SetResultCount(resultCount)
	msg.SetQueryHashId(queryHashId)
	msg.SetHasResult(true)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestRollbackTransactionRequestMessage() *RollbackTransactionRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewRollbackTransactionRequestMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestRollbackTransactionResponseMessage() *RollbackTransactionResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewRollbackTransactionResponseMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestSessionForcefullyTerminatedMessage() *SessionForcefullyTerminatedMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	exceptionId := rand.Int()
	exceptionMsg := "Test-Forcefully-Terminated"
	msg := NewSessionForcefullyTerminatedMessage(authToken, sessionId)
	msg.SetExceptionType(exceptionId)
	msg.SetExceptionMsg(exceptionMsg)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestTraverseRequestMessage() *TraverseRequestMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewTraverseRequestMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestTraverseResponseMessage() *TraverseResponseMessage {
	authToken := rand.Int63()
	sessionId := rand.Int63()
	msg := NewTraverseResponseMessage(authToken, sessionId)
	msg.BufLength = int(reflect.TypeOf(*msg).Size())
	return msg
}

func createTestMessageForVerb(verbId int) types.TGMessage {
	switch verbId {
	case VerbPingMessage:
		return createTestPingMessage()
	case VerbHandShakeRequest:
		return createTestHandShakeRequestMessage()
	case VerbHandShakeResponse:
		return createTestHandShakeResponseMessage()
	case VerbAuthenticateRequest:
		return createTestAuthenticateRequestMessage()
	case VerbAuthenticateResponse:
		return createTestAuthenticateResponseMessage()
	case VerbBeginTransactionRequest:
		return createTestBeginTransactionRequestMessage()
	case VerbBeginTransactionResponse:
		return createTestBeginTransactionResponseMessage()
	case VerbCommitTransactionRequest:
		return createTestCommitTransactionRequestMessage()
	case VerbCommitTransactionResponse:
		return createTestCommitTransactionResponseMessage()
	case VerbRollbackTransactionRequest:
		return createTestRollbackTransactionRequestMessage()
	case VerbRollbackTransactionResponse:
		return createTestRollbackTransactionResponseMessage()
	case VerbQueryRequest:
		return createTestQueryRequestMessage()
	case VerbQueryResponse:
		return createTestQueryResponseMessage()
	case VerbTraverseRequest:
		return createTestTraverseRequestMessage()
	case VerbTraverseResponse:
		return createTestTraverseResponseMessage()
	case VerbAdminRequest:
		fallthrough
		//return createTestAdminRequestMessage()
	case VerbAdminResponse:
		fallthrough
		//return createTestAdminResponseMessage()
	case VerbMetadataRequest:
		return createTestMetadataRequestMessage()
	case VerbMetadataResponse:
		return createTestMetadataResponseMessage()
	case VerbGetEntityRequest:
		return createTestGetEntityRequestMessage()
	case VerbGetEntityResponse:
		return createTestGetEntityResponseMessage()
	case VerbGetLargeObjectRequest:
		return createTestGetLargeObjectRequestMessage()
	case VerbGetLargeObjectResponse:
		return createTestGetLargeObjectResponseMessage()
	case VerbDumpStacktraceRequest:
		fallthrough
		//return createTestDumpStacktraceRequestMessage()
	case VerbDisconnectChannelRequest:
		return createTestDisconnectChannelRequestMessage()
	case VerbSessionForcefullyTerminated:
		return createTestSessionForcefullyTerminatedMessage()
	case VerbDecryptBufferRequest:
		fallthrough
		//return createTestDecryptBufferRequestMessage()
	case VerbDecryptBufferResponse:
		fallthrough
		//return createTestDecryptBufferResponseMessage()
	case VerbExceptionMessage:
		return createTestExceptionMessage()
	case VerbInvalidMessage:
		return createTestInvalidMessage()
	default:
		return createTestInvalidMessage()
	}
}

func TestCreateMessageForVerb(t *testing.T) {
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 39 {
			continue
		}
		newMsg, err := CreateMessageForVerb(inputVerbId)
		if err != nil {
			t.Errorf("MessageFactory could not instantiate %s message for verbId: '%+v'", GetVerb(inputVerbId).name, inputVerbId)
			break
		}
		t.Logf("MessageFactory returned %s message for verbId: '%+v' as '%+v'", GetVerb(inputVerbId).name, inputVerbId, newMsg.String())
	}
}

func TestCreateMessageWithToken(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}
		authToken := rand.Int63()
		sessionId := rand.Int63()
		newMsg, err := CreateMessageWithToken(inputVerbId, authToken, sessionId)
		if err != nil {
			t.Errorf("MessageFactory could not instantiate %s message w/ token %d and sessionId %d for verbId: '%+v'", GetVerb(inputVerbId).name, authToken, sessionId, inputVerbId)
			break
		}
		t.Logf("MessageFactory returned %s message w/ token %d and sessionId %d for verbId: '%+v' as '%+v'", GetVerb(inputVerbId).name, authToken, sessionId, inputVerbId, newMsg.String())
	}
}

func TestCreateMessageFromBuffer(t *testing.T) {
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 3 || inputVerbId == 4 || inputVerbId == 8 || inputVerbId == 12 || inputVerbId == 20  || inputVerbId == 22 || inputVerbId == 24 {
			continue
		}
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}

		// Execute individual message type's method
		newMsg := createTestMessageForVerb(inputVerbId)
		t.Logf("MessageFactory::TestCreateMessageFromBuffer starting w/ '%+v'", newMsg.String())
		buf, bufLen, err := newMsg.ToBytes()
		if err != nil {
			t.Logf("Error: %+v", err)
			t.Errorf("MessageFactory::TestCreateMessageFromBuffer could not generate buffer for verbId: '%+v'", inputVerbId)
		}
		t.Logf("MessageFactory::TestCreateMessageFromBuffer ToBytes resulted in '%+v'", buf)

		constructedMsg, err := CreateMessageFromBuffer(buf[0:bufLen], 0, bufLen)
		if err != nil {
			t.Errorf("MessageFactory::TestCreateMessageFromBuffer could not instantiate %s message from buffer for verbId: '%+v'", GetVerb(inputVerbId).name, inputVerbId)
			break
		}
		t.Logf("MessageFactory::TestCreateMessageFromBuffer execution resulted in '%+v'", constructedMsg.String())
	}
}

// This automatically will test both APIs - (a) ToBytes and (b) FromBytes
func TestFromBytes(t *testing.T) {
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 3 || inputVerbId == 4 || inputVerbId == 8 || inputVerbId == 12 || inputVerbId == 20  || inputVerbId == 22 || inputVerbId == 24 {
			continue
		}
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}

		// Execute individual message type's method
		newMsg := createTestMessageForVerb(inputVerbId)
		t.Logf("MessageFactory::TestFromBytes starting w/ '%+v'-'%+v'", newMsg, newMsg.String())
		buf, bufLen, err := newMsg.ToBytes()
		if err != nil {
			t.Logf("Error: %+v", err)
			t.Errorf("MessageFactory::TestFromBytes could not generate buffer for verbId: '%+v'", inputVerbId)
		}
		t.Logf("MessageFactory::TestFromBytes ToBytes resulted in '%+v'", buf)

		// Execute individual message type's method
		constructedMsg, err := newMsg.FromBytes(buf[0:bufLen])
		if err != nil {
			t.Logf("Error: %+v", err)
			t.Errorf("MessageFactory::TestFromBytes could not generate buffer for verbId: '%+v'", inputVerbId)
		}
		t.Logf("MessageFactory::TestFromBytes execution resulted in '%+v'", constructedMsg.String())
	}
}

func TestToBytes(t *testing.T) {
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}
		// Execute individual message type's method
		newMsg := createTestMessageForVerb(inputVerbId)
		t.Logf("MessageFactory::TestToBytes starting w/ '%+v'", newMsg.String())
		buf, _, err := newMsg.ToBytes()
		if err != nil {
			t.Logf("Error: %+v", err)
			t.Errorf("MessageFactory::TestToBytes could not generate buffer for verbId: '%+v'", inputVerbId)
		}
		t.Logf("MessageFactory::ToBytes resulted in '%+v'", buf)
	}
}

func TestGetTimestamp(t *testing.T) {
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}
		ts := GetTimestamp(inputVerbId)
		t.Logf("MessageFactory::GetTimestamp resulted in '%+v'", ts)
	}
}

func TestSetTimestamp(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}
		ts := rand.Int63()
		err := SetTimestamp(inputVerbId, ts)
		if err != nil {
			t.Logf("Error: %+v", err)
		}
		t.Logf("MessageFactory::SetTimestamp updated messages w/ '%+v'", ts)
	}
}

func TestUpdateSequenceAndTimeStamp(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	for inputVerbId := range PreDefinedVerbs {
		if inputVerbId == 15 || inputVerbId == 16 || inputVerbId == 39 {
			continue
		}
		ts := rand.Int63()
		err := UpdateSequenceAndTimeStamp(inputVerbId, ts)
		if err != nil {
			t.Logf("Error: %+v", err)
		}
		t.Logf("MessageFactory::UpdateSequenceAndTimeStamp updated messages w/ '%+v'", ts)
	}
}
