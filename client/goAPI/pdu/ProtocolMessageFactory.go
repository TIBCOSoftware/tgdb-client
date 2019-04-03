package pdu

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: MessageFactory.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**
 * The Server describes the pdu header as below.
 * struct _tg_pduheader_t_ {
	 tg_int32    length;         //length of the message including the header
	 tg_int32    magic;          //Magic to recognize this is our message
	 tg_int16    protVersion;    //protocol version
	 tg_pduverb  verbId;         //we write the verb as a short value
	 tg_uint64   sequenceNo;     //message SystemTypeSequence No from the client
	 tg_uint64   timestamp;      //Timestamp of the message sent.
	 tg_uint64   requestId;      //Unique _request Identifier from the client, which is returned
	 tg_int32    dataOffset;     //Offset from where the payload begins
 }
*/

var logger = logging.DefaultTGLogManager().GetLogger()

/////////////////////////////////////////////////////////////////
// Helper Public functions for TGMessageFactory
/////////////////////////////////////////////////////////////////

// Create new message instance based on the input type
func CreateMessageForVerb(verbId int) (types.TGMessage, types.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputVerbId := verbId

	// Use a switch case to switch between message types, if a type exist then error is nil (null)
	// Whenever new message type gets into the mix, just add a case below
	switch inputVerbId {
	case VerbPingMessage:
		return DefaultPingMessage(), nil
	case VerbHandShakeRequest:
		return DefaultHandShakeRequestMessage(), nil
	case VerbHandShakeResponse:
		return DefaultHandShakeResponseMessage(), nil
	case VerbAuthenticateRequest:
		return DefaultAuthenticateRequestMessage(), nil
	case VerbAuthenticateResponse:
		return DefaultAuthenticateResponseMessage(), nil
	case VerbBeginTransactionRequest:
		return DefaultBeginTransactionRequestMessage(), nil
	case VerbBeginTransactionResponse:
		return DefaultBeginTransactionResponseMessage(), nil
	case VerbCommitTransactionRequest:
		return DefaultCommitTransactionRequestMessage(), nil
	case VerbCommitTransactionResponse:
		return DefaultCommitTransactionResponseMessage(), nil
	case VerbRollbackTransactionRequest:
		return DefaultRollbackTransactionRequestMessage(), nil
	case VerbRollbackTransactionResponse:
		return DefaultRollbackTransactionResponseMessage(), nil
	case VerbQueryRequest:
		return DefaultQueryRequestMessage(), nil
	case VerbQueryResponse:
		return DefaultQueryResponseMessage(), nil
	case VerbTraverseRequest:
		return DefaultTraverseRequestMessage(), nil
	case VerbTraverseResponse:
		return DefaultTraverseResponseMessage(), nil
	case VerbAdminRequest:
		fallthrough
		//return DefaultAdminRequestMessage(), nil
	case VerbAdminResponse:
		fallthrough
		//return DefaultAdminResponseMessage(), nil
	case VerbMetadataRequest:
		return DefaultMetadataRequestMessage(), nil
	case VerbMetadataResponse:
		return DefaultMetadataResponseMessage(), nil
	case VerbGetEntityRequest:
		return DefaultGetEntityRequestMessage(), nil
	case VerbGetEntityResponse:
		return DefaultGetEntityResponseMessage(), nil
	case VerbGetLargeObjectRequest:
		return DefaultGetLargeObjectRequestMessage(), nil
	case VerbGetLargeObjectResponse:
		return DefaultGetLargeObjectResponseMessage(), nil
	case VerbDumpStacktraceRequest:
		fallthrough
		//return DefaultDumpStacktraceRequestMessage(), nil
	case VerbDisconnectChannelRequest:
		return DefaultDisconnectChannelRequestMessage(), nil
	case VerbSessionForcefullyTerminated:
		return DefaultSessionForcefullyTerminatedMessage(), nil
	case VerbDecryptBufferRequest:
		return DefaultDecryptBufferRequestMessage(), nil
	case VerbDecryptBufferResponse:
		return DefaultDecryptBufferResponseMessage(), nil
	case VerbExceptionMessage:
		return DefaultExceptionMessage(), nil
	case VerbInvalidMessage:
		return DefaultInvalidMessage(), nil
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Message Type '%s'", GetVerb(inputVerbId).name)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
}

func CreateMessageWithToken(verbId int, authToken, sessionId int64) (types.TGMessage, types.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputVerbId := verbId

	// Use a switch case to switch between message types, if a type exist then error is nil (null)
	// Whenever new message type gets into the mix, just add a case below
	switch inputVerbId {
	case VerbPingMessage:
		return NewPingMessage(authToken, sessionId), nil
	case VerbHandShakeRequest:
		return NewHandShakeRequestMessage(authToken, sessionId), nil
	case VerbHandShakeResponse:
		return NewHandShakeResponseMessage(authToken, sessionId), nil
	case VerbAuthenticateRequest:
		return NewAuthenticateRequestMessage(authToken, sessionId), nil
	case VerbAuthenticateResponse:
		return NewAuthenticateResponseMessage(authToken, sessionId), nil
	case VerbBeginTransactionRequest:
		return NewBeginTransactionRequestMessage(authToken, sessionId), nil
	case VerbBeginTransactionResponse:
		return NewBeginTransactionResponseMessage(authToken, sessionId), nil
	case VerbCommitTransactionRequest:
		return NewCommitTransactionRequestMessage(authToken, sessionId), nil
	case VerbCommitTransactionResponse:
		return NewCommitTransactionResponseMessage(authToken, sessionId), nil
	case VerbRollbackTransactionRequest:
		return NewRollbackTransactionRequestMessage(authToken, sessionId), nil
	case VerbRollbackTransactionResponse:
		return NewRollbackTransactionResponseMessage(authToken, sessionId), nil
	case VerbQueryRequest:
		return NewQueryRequestMessage(authToken, sessionId), nil
	case VerbQueryResponse:
		return NewQueryResponseMessage(authToken, sessionId), nil
	case VerbTraverseRequest:
		return NewTraverseRequestMessage(authToken, sessionId), nil
	case VerbTraverseResponse:
		return NewTraverseResponseMessage(authToken, sessionId), nil
	case VerbAdminRequest:
		fallthrough
		//return NewAdminRequestMessage(authToken, sessionId), nil
	case VerbAdminResponse:
		fallthrough
		//return NewAdminResponseMessage(authToken, sessionId), nil
	case VerbMetadataRequest:
		return NewMetadataRequestMessage(authToken, sessionId), nil
	case VerbMetadataResponse:
		return NewMetadataResponseMessage(authToken, sessionId), nil
	case VerbGetEntityRequest:
		return NewGetEntityRequestMessage(authToken, sessionId), nil
	case VerbGetEntityResponse:
		return NewGetEntityResponseMessage(authToken, sessionId), nil
	case VerbGetLargeObjectRequest:
		return NewGetLargeObjectRequestMessage(authToken, sessionId), nil
	case VerbGetLargeObjectResponse:
		return NewGetLargeObjectResponseMessage(authToken, sessionId), nil
	case VerbDumpStacktraceRequest:
		fallthrough
		//return NewDumpStacktraceRequestMessage(authToken, sessionId), nil
	case VerbDisconnectChannelRequest:
		return NewDisconnectChannelRequestMessage(authToken, sessionId), nil
	case VerbSessionForcefullyTerminated:
		return NewSessionForcefullyTerminatedMessage(authToken, sessionId), nil
	case VerbDecryptBufferRequest:
		return NewDecryptBufferRequestMessage(authToken, sessionId), nil
	case VerbDecryptBufferResponse:
		return NewDecryptBufferResponseMessage(authToken, sessionId), nil
	case VerbExceptionMessage:
		return NewExceptionMessage(authToken, sessionId), nil
	case VerbInvalidMessage:
		return NewInvalidMessage(authToken, sessionId), nil
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Message Type '%s'", GetVerb(inputVerbId).name)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
}

func CreateMessageFromBuffer(buffer []byte, offset int, length int) (types.TGMessage, types.TGError) {
	buf := make([]byte, 0)
	logger.Log(fmt.Sprintf("Entering MessageFactory::CreateMessageFromBuffer received buffer(BufLen: %d, Offset: %d, Len: %d)", len(buffer), offset, length))

	if len(buffer) == length {
		buf = buffer
	} else {
		buf = append(buffer[offset:length])
	}
	//logger.Log(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer buf is '%+v'", buf))

	commandVerb, err := VerbIdFromBytes(buf)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in extracting verbId from message buffer: %s", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer retrieved VerbId '%s'", commandVerb.name))
	msg, err := CreateMessageForVerb(commandVerb.id)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in creating message for verb('%s'): %s", commandVerb.name, err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside MessageFactory::CreateMessageFromBuffer created '%+v'", msg))
	msg1, err := msg.FromBytes(buffer)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MessageFactory::CreateMessageFromBuffer w/ error in updating message contents from buffer: %s", err.Error()))
		return nil, err
	}
	return msg1, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

func FromBytes(verbId int, buffer []byte) (types.TGMessage, types.TGError) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return nil, err
	}
	// Execute Individual Message's method
	return msg.FromBytes(buffer)
}

func ToBytes(verbId int) ([]byte, int, types.TGError) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return nil, -1, err
	}
	// Execute Individual Message's method
	return msg.ToBytes()
}

// GetAuthToken gets Authorization Token for specified message type
func GetAuthToken(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetAuthToken()
}

func GetIsUpdatable(verbId int) bool {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetIsUpdatable()
}

func GetMessageByteBufLength(verbId int) int {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetMessageByteBufLength()
}

func GetRequestId(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetRequestId()
}

func GetSequenceNo(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetSequenceNo()
}

// GetSessionId gets Session id for specified message type
func GetSessionId(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetSessionId()
}

func GetTimestamp(verbId int) int64 {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetTimestamp()
}

func GetVerbId(verbId int) int {
	msg, _ := CreateMessageForVerb(verbId)
	// Execute Derived / Dependent Message's method
	return msg.GetVerbId()
}

func SetAuthToken(verbId int, authToken int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetAuthToken(authToken)
}

//func SetIsUpdatable(verbId int, updateFlag bool) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetIsUpdatable(updateFlag)
//}

//func SetMessageByteBufLength(verbId int, bufLength int) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetMessageByteBufLength(bufLength)
//}

func SetRequestId(verbId int, requestId int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetRequestId(requestId)
}

//func SetSequenceNo(verbId int, sequenceNo int64) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetSequenceNo(sequenceNo)
//}

func SetSessionId(verbId int, sessionId int64) {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return
	}
	// Execute Derived / Dependent Message's method
	msg.SetSessionId(sessionId)
}

func SetTimestamp(verbId int, timestamp int64) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Derived / Dependent Message's method
	return msg.SetTimestamp(timestamp)
}

//func SetVerbId(verbId int) {
//	msg, err := CreateMessageForVerb(verbId)
//	if err != nil {
//		return
//	}
//	// Execute Derived / Dependent Message's method
//	return msg.SetVerbId(verbId)
//}

func UpdateSequenceAndTimeStamp(verbId int, timestamp int64) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.UpdateSequenceAndTimeStamp(timestamp)
}

func ReadHeader(verbId int, is types.TGInputStream) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	//
	// Execute Individual Message's method
	return msg.ReadHeader(is)
}

func WriteHeader(verbId int, os types.TGOutputStream) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.WriteHeader(os)
}

func ReadPayload(verbId int, is types.TGInputStream) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.ReadPayload(is)
}

func WritePayload(verbId int, os types.TGOutputStream) types.TGError {
	msg, err := CreateMessageForVerb(verbId)
	if err != nil {
		return err
	}
	// Execute Individual Message's method
	return msg.WritePayload(os)
}
