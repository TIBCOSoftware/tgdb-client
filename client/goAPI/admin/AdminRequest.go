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
 * File name: VerbAdminRequest.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type AdminRequestMessage struct {
	*pdu.AbstractProtocolMessage
	command    AdminCommand
	logDetails *ServerLogDetails
}

func DefaultAdminRequestMessage() *AdminRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AdminRequestMessage{})

	newMsg := AdminRequestMessage{
		AbstractProtocolMessage: pdu.DefaultAbstractProtocolMessage(),
		command:                 AdminCommandInvalid,
	}
	newMsg.HeaderSetVerbId(pdu.VerbAdminRequest)
	newMsg.HeaderSetMessageByteBufLength(int(reflect.TypeOf(newMsg).Size()))
	newMsg.SetUpdatableFlag(true)
	return &newMsg
}

// Create New Message Instance
func NewAdminRequestMessage(authToken, sessionId int64) *AdminRequestMessage {
	newMsg := DefaultAdminRequestMessage()
	newMsg.HeaderSetAuthToken(authToken)
	newMsg.HeaderSetSessionId(sessionId)
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *AdminRequestMessage) GetCommand() AdminCommand {
	return msg.command
}

func (msg *AdminRequestMessage) GetLogLevel() *ServerLogDetails {
	return msg.logDetails
}

func (msg *AdminRequestMessage) SetCommand(cmd AdminCommand) {
	msg.command = cmd
}

func (msg *AdminRequestMessage) SetLogLevel(connId *ServerLogDetails) {
	msg.logDetails = connId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AdminRequestMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminRequest:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminRequest:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminRequest:FromBytes - about to APMReadHeader"))
	err = pdu.APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminRequest:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("AdminRequest::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AdminRequestMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering AdminRequest:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside AdminRequest:ToBytes - about to APMWriteHeader"))
	err := pdu.APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside AdminRequest:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("AdminRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AdminRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AdminRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AdminRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AdminRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AdminRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AdminRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AdminRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AdminRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AdminRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AdminRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AdminRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AdminRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AdminRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AdminRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AdminRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AdminRequestMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", pdu.GetVerb(msg.GetVerbId()).GetName())
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AdminRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AdminRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AdminRequest:{")
	buffer.WriteString(fmt.Sprintf("Command: %d", msg.command))
	buffer.WriteString(fmt.Sprintf(", ServerLogDetails: %+v", msg.logDetails))
	strArray := []string{buffer.String(), msg.APMMessageToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AdminRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AdminRequestMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return pdu.APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *AdminRequestMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return pdu.APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *AdminRequestMessage) ReadPayload(is types.TGInputStream) types.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *AdminRequestMessage) WritePayload(os types.TGOutputStream) types.TGError {
	//startPos := os.GetPosition()
	//logger.Log(fmt.Sprintf("Entering AdminRequest:WritePayload at output buffer position: '%d'", startPos))
	//os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the commit buffer length
	//os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the checksum for the commit buffer to be added later.  Currently not used

	dataLen := 0
	checkSum := 0

	switch msg.command {
	case AdminCommandCreateUser:
	case AdminCommandCreateAttrDesc:
	case AdminCommandCreateIndex:
	case AdminCommandCreateNodeType:
	case AdminCommandCreateEdgeType:
	case AdminCommandShowUsers:
		fallthrough
	case AdminCommandShowAttrDescs:
		fallthrough
	case AdminCommandShowIndices:
		fallthrough
	case AdminCommandShowTypes:
		fallthrough
	case AdminCommandShowInfo:
		fallthrough
	case AdminCommandShowConnections:
		os.(*iostream.ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(int(msg.command))
	case AdminCommandDescribe:
	case AdminCommandSetLogLevel:
		os.(*iostream.ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(int(msg.command))
		os.(*iostream.ProtocolDataOutputStream).WriteShort(int(msg.logDetails.GetLogLevel()))
		os.(*iostream.ProtocolDataOutputStream).WriteLong(int64(msg.logDetails.GetLogComponent()))
	case AdminCommandStopServer:
		fallthrough
	case AdminCommandCheckpointServer:
		os.(*iostream.ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(int(msg.command))
	case AdminCommandDisconnectClient:
	case AdminCommandKillConnection:
		os.(*iostream.ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*iostream.ProtocolDataOutputStream).WriteInt(int(msg.command))
		os.(*iostream.ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true)
	default:
	}

	//currPos := os.GetPosition()
	//length := currPos - startPos
	//_, err := os.(*iostream.ProtocolDataOutputStream).WriteIntAt(startPos, length)
	//if err != nil {
	//	return err
	//}
	//logger.Log(fmt.Sprintf("Returning AdminRequest::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AdminRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.HeaderGetMessageByteBufLength(), msg.HeaderGetVerbId(),
		msg.HeaderGetSequenceNo(), msg.HeaderGetTimestamp(), msg.HeaderGetRequestId(),
		msg.HeaderGetDataOffset(), msg.HeaderGetAuthToken(), msg.HeaderGetSessionId(), msg.GetIsUpdatable(),
		msg.command, msg.logDetails)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AdminRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	var bLen, vId int
	var seq, tStamp, reqId, token, sId int64
	var offset int16
	var uFlag bool
	_, err := fmt.Fscanln(b, &bLen, &vId, &seq, &tStamp, &reqId, &offset, &token, &sId, &uFlag, &msg.command, &msg.logDetails)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
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
