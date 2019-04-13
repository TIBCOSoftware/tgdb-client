package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strings"
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
 * File name: VerbAdminResponse.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type AdminResponseMessage struct {
	*pdu.AbstractProtocolMessage
	attrDescriptors []types.TGAttributeDescriptor
	connections     []TGConnectionInfo
	indices         []TGIndexInfo
	serverInfo      *ServerInfoImpl
	users           []TGUserInfo
}

func DefaultAdminResponseMessage() *AdminResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AdminResponseMessage{})

	newMsg := AdminResponseMessage{
		AbstractProtocolMessage: pdu.DefaultAbstractProtocolMessage(),
	}
	newMsg.HeaderSetVerbId(pdu.VerbAdminResponse)
	newMsg.HeaderSetMessageByteBufLength(int(reflect.TypeOf(newMsg).Size()))
	return &newMsg
}

// Create New Message Instance
func NewAdminResponseMessage(authToken, sessionId int64) *AdminResponseMessage {
	newMsg := DefaultAdminResponseMessage()
	newMsg.HeaderSetAuthToken(authToken)
	newMsg.HeaderSetSessionId(sessionId)
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *AdminResponseMessage) GetDescriptorList() []types.TGAttributeDescriptor {
	return msg.attrDescriptors
}

func (msg *AdminResponseMessage) GetConnectionList() []TGConnectionInfo {
	return msg.connections
}

func (msg *AdminResponseMessage) GetIndexList() []TGIndexInfo {
	return msg.indices
}

func (msg *AdminResponseMessage) GetServerInfo() *ServerInfoImpl {
	return msg.serverInfo
}

func (msg *AdminResponseMessage) GetUserList() []TGUserInfo {
	return msg.users
}

func (msg *AdminResponseMessage) SetDescriptorList(list []types.TGAttributeDescriptor) {
	msg.attrDescriptors = list
}

func (msg *AdminResponseMessage) SetConnectionList(list []TGConnectionInfo) {
	msg.connections = list
}

func (msg *AdminResponseMessage) SetIndexList(list []TGIndexInfo) {
	msg.indices = list
}

func (msg *AdminResponseMessage) SetServerInfo(sInfo *ServerInfoImpl) {
	msg.serverInfo = sInfo
}

func (msg *AdminResponseMessage) SetUserList(list []TGUserInfo) {
	msg.users = list
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AdminResponseMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminResponse:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponse:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminResponse:FromBytes - about to APMReadHeader"))
	err = pdu.APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminResponse:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("AdminResponse::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AdminResponseMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminResponse:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside AdminResponse:ToBytes - about to APMWriteHeader"))
	err := pdu.APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminResponse:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("AdminResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AdminResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AdminResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AdminResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AdminResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AdminResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AdminResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AdminResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AdminResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AdminResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AdminResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AdminResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AdminResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AdminResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AdminResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AdminResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AdminResponseMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", pdu.GetVerb(msg.GetVerbId()).GetName())
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AdminResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AdminResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AdminResponse:{")
	buffer.WriteString(fmt.Sprintf("AttributeDescriptors: %+v", msg.attrDescriptors))
	buffer.WriteString(fmt.Sprintf(", Connections: %+v", msg.connections))
	buffer.WriteString(fmt.Sprintf(", Indices: %+v", msg.indices))
	buffer.WriteString(fmt.Sprintf(", ServerInfo: %+v", msg.serverInfo))
	buffer.WriteString(fmt.Sprintf(", Users: %+v", msg.users))
	strArray := []string{buffer.String(), msg.APMMessageToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AdminResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AdminResponseMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return pdu.APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AdminResponseMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return pdu.APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *AdminResponseMessage) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AdminResponseMessage:ReadPayload"))
	//bLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buf length
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bLen from message buffer"))
	//	return err
	//}
	//logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read buf length as '%+v'", bLen))
	//
	//checkSum, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // checksum
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading checkSum from message buffer"))
	//	return err
	//}
	//logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read checkSum as '%+v'", checkSum))

	resultId, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // result Id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading resultId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read resultId as '%+v'", resultId))

	command, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // command
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading command from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read command as '%+v'", command))

	switch AdminCommand(command) {
	case AdminCommandCreateUser:
	case AdminCommandCreateAttrDesc:
	case AdminCommandCreateIndex:
	case AdminCommandCreateNodeType:
	case AdminCommandCreateEdgeType:
	case AdminCommandShowUsers:
		userList, err := extractUserListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userList from message buffer"))
			return err
		}
		msg.SetUserList(userList)
	case AdminCommandShowAttrDescs:
		attrDescList, err := extractDescriptorListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrDescList from message buffer"))
			return err
		}
		msg.SetDescriptorList(attrDescList)
	case AdminCommandShowIndices:
		indexList, err := extractIndexListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexList from message buffer"))
			return err
		}
		msg.SetIndexList(indexList)
	case AdminCommandShowTypes:
	case AdminCommandShowInfo:
		serverInfo, err := extractServerInfoFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading serverInfo from message buffer"))
			return err
		}
		msg.SetServerInfo(serverInfo)
	case AdminCommandShowConnections:
		connList, err := extractConnectionListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading connList from message buffer"))
			return err
		}
		msg.SetConnectionList(connList)
	case AdminCommandDescribe:
	case AdminCommandSetLogLevel:
	case AdminCommandStopServer:
	case AdminCommandCheckpointServer:
	case AdminCommandDisconnectClient:
	case AdminCommandKillConnection:
	default:
	}
	logger.Log(fmt.Sprint("Returning AdminResponseMessage:ReadPayload"))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *AdminResponseMessage) WritePayload(os types.TGOutputStream) types.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AdminResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.HeaderGetMessageByteBufLength(), msg.HeaderGetVerbId(),
		msg.HeaderGetSequenceNo(), msg.HeaderGetTimestamp(), msg.HeaderGetRequestId(),
		msg.HeaderGetDataOffset(), msg.HeaderGetAuthToken(), msg.HeaderGetSessionId(),
		msg.GetIsUpdatable(), msg.attrDescriptors, msg.connections, msg.indices, msg.serverInfo, msg.users)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminResponse:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AdminResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	var bLen, vId int
	var seq, tStamp, reqId, token, sId int64
	var offset int16
	var uFlag bool
	_, err := fmt.Fscanln(b, &bLen, &vId, &seq, &tStamp, &reqId, &offset, &token, &sId, &uFlag,
		&msg.attrDescriptors, &msg.connections, &msg.indices, &msg.serverInfo, &msg.users)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminResponse:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	msg.HeaderSetMessageByteBufLength(bLen)
	msg.HeaderSetVerbId(vId)
	msg.HeaderSetSequenceNo(seq)
	msg.HeaderSetTimestamp(tStamp)
	msg.HeaderSetRequestId(reqId)
	msg.HeaderSetAuthToken(token)
	msg.HeaderSetSessionId(sId)
	msg.HeaderSetDataOffset(offset)
	msg.SetIsUpdatable(uFlag)
	return nil
}
