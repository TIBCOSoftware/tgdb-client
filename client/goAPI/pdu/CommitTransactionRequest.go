package pdu

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
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
 * File name: VerbCommitTransactionRequest.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type CommitTransactionRequest struct {
	*AbstractProtocolMessage
	AddedList   map[int64]types.TGEntity
	//ChangedList map[int64]types.TGEntity
	UpdatedList map[int64]types.TGEntity
	RemovedList map[int64]types.TGEntity
	AttrDescSet []types.TGAttributeDescriptor
}

func DefaultCommitTransactionRequestMessage() *CommitTransactionRequest {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CommitTransactionRequest{})

	newMsg := CommitTransactionRequest{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.VerbId = VerbCommitTransactionRequest
	newMsg.AddedList = make(map[int64]types.TGEntity, 0)
	//newMsg.ChangedList = make(map[int64]types.TGEntity, 0)
	newMsg.UpdatedList = make(map[int64]types.TGEntity, 0)
	newMsg.RemovedList = make(map[int64]types.TGEntity, 0)
	newMsg.AttrDescSet = make([]types.TGAttributeDescriptor, 0)
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewCommitTransactionRequestMessage(authToken, sessionId int64) *CommitTransactionRequest {
	newMsg := DefaultCommitTransactionRequestMessage()
	newMsg.AuthToken = authToken
	newMsg.SessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for VerbCommitTransactionRequest message
/////////////////////////////////////////////////////////////////

func (msg *CommitTransactionRequest) AddCommitLists(addedList, updatedList, removedList map[int64]types.TGEntity, attrDescriptors []types.TGAttributeDescriptor) *CommitTransactionRequest {
	if len(addedList) > 0 {
		msg.AddedList = addedList
	}
	if len(updatedList) > 0 {
		msg.UpdatedList = updatedList
	}
	if len(removedList) > 0 {
		msg.RemovedList = removedList
	}
	if len(attrDescriptors) > 0 {
		msg.AttrDescSet = attrDescriptors
	}
	return msg
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *CommitTransactionRequest) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering CommitTransactionRequest:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}
	logger.Log(fmt.Sprint("Entering CommitTransactionRequest::FromBytes"))

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionRequest:FromBytes - about to readHeader"))
	err = msg.readHeader(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionRequest:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("Returning CommitTransactionRequest::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *CommitTransactionRequest) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering CommitTransactionRequest::ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside CommitTransactionRequest:ToBytes - about to writeHeader"))
	err := msg.writeHeader(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside CommitTransactionRequest:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("Returning CommitTransactionRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()[:os.GetLength()]))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *CommitTransactionRequest) GetAuthToken() int64 {
	return msg.getAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *CommitTransactionRequest) GetIsUpdatable() bool {
	return msg.getIsUpdatable()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *CommitTransactionRequest) GetMessageByteBufLength() int {
	return msg.getMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *CommitTransactionRequest) GetRequestId() int64 {
	return msg.getRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *CommitTransactionRequest) GetSequenceNo() int64 {
	return msg.getSequenceNo()
}

// GetSessionId gets the session id
func (msg *CommitTransactionRequest) GetSessionId() int64 {
	return msg.getSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *CommitTransactionRequest) GetTimestamp() int64 {
	return msg.getTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *CommitTransactionRequest) GetVerbId() int {
	return msg.getVerbId()
}

// SetAuthToken sets the authToken
func (msg *CommitTransactionRequest) SetAuthToken(authToken int64) {
	msg.setAuthToken(authToken)
}

// SetRequestId sets the request id
func (msg *CommitTransactionRequest) SetRequestId(requestId int64) {
	msg.setRequestId(requestId)
}

// SetSessionId sets the session id
func (msg *CommitTransactionRequest) SetSessionId(sessionId int64) {
	msg.setSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *CommitTransactionRequest) SetTimestamp(timestamp int64) types.TGError {
	return msg.setTimestamp(timestamp)
}

func (msg *CommitTransactionRequest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CommitTransactionRequest:{")
	buffer.WriteString(fmt.Sprint("AttrDescSet: "))
	buffer.WriteString("{")
	//for _, d := range msg.AttrDescSet {
	//	buffer.WriteString(fmt.Sprintf("Attribute Descriptor: %+v ", d))
	//}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", AddedList:{"))
	for k, v := range msg.AddedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
	}
	buffer.WriteString("}")
	//buffer.WriteString(fmt.Sprint("ChangedList: "))
	//for k, v := range msg.ChangedList {
	//	buffer.WriteString(fmt.Sprintf("EntityId: %d Entity: %+v ", k, v))
	//}
	buffer.WriteString(fmt.Sprint(", UpdatedList:{"))
	for k, v := range msg.UpdatedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", RemovedList:{"))
	for k, v := range msg.RemovedList {
		buffer.WriteString(fmt.Sprintf("EntityId: %d=Entity: %+v ", k, v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.messageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *CommitTransactionRequest) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.updateSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *CommitTransactionRequest) ReadHeader(is types.TGInputStream) types.TGError {
	return msg.readHeader(is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *CommitTransactionRequest) WriteHeader(os types.TGOutputStream) types.TGError {
	return msg.writeHeader(os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *CommitTransactionRequest) ReadPayload(is types.TGInputStream) types.TGError {
	//Commit response need to send back real id for all entities and descriptors.
	// No-op for Now
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *CommitTransactionRequest) WritePayload(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering CommitTransactionRequest:ReadPayload at output buffer position: '%d'", startPos))
	os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the commit buffer length
	os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the checksum for the commit buffer to be added later.  Currently not used
	////<A> for attribute descriptor, <N> for node desc definitions, <E> for edge desc definitions
	////meta should be sent before the instance data
	if len(msg.AttrDescSet) > 0 {
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0x1010) // For attribute descriptor
		// There should be nothing after the marker due to no new attribute descriptor
		// Need to check for new descriptor only with attribute id as negative number
		// Check for size overrun
		newAttrCount := 0
		for _, attrDesc := range msg.AttrDescSet {
			if attrDesc.GetAttributeId() < 0 {
				newAttrCount++
			}
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:ReadPayload - There are '%d' new attribute descriptors", newAttrCount))
		os.(*iostream.ProtocolDataOutputStream).WriteInt(newAttrCount)
		for _, attrDesc := range msg.AttrDescSet {
				err := attrDesc.WriteExternal(os)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing attrDesc '%+v' to message buffer", attrDesc))
					return err
				}
			//}
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:ReadPayload - '%d' attribute descriptors are written in byte format", len(msg.AttrDescSet)))
	}
	if len(msg.AddedList) > 0 {
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0x1011) // For entity creation
		os.(*iostream.ProtocolDataOutputStream).WriteInt(len(msg.AddedList))
		for _, entity := range msg.AddedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing addedEntity to message buffer"))
				return err
			}
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' new entities are written in byte format", len(msg.AddedList)))
	}
	//TODO: Ask TGDB Engineering Team - Need to write only the modified attributes
	if len(msg.UpdatedList) > 0 {
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0x1012) // For entity update
		os.(*iostream.ProtocolDataOutputStream).WriteInt(len(msg.UpdatedList))
		for _, entity := range msg.UpdatedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing updatedEntity to message buffer"))
				return err
			}
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' updateable entities are written in byte format", len(msg.UpdatedList)))
	}
	//TODO: Ask TGDB Engineering Team - Need to write the id only
	if len(msg.RemovedList) > 0 {
		os.(*iostream.ProtocolDataOutputStream).WriteShort(0x1013) // For deleted entities
		os.(*iostream.ProtocolDataOutputStream).WriteInt(len(msg.RemovedList))
		for _, entity := range msg.RemovedList {
			err := entity.WriteExternal(os)
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning CommitTransactionRequest:WritePayload w/ Error in writing removedEntity to message buffer"))
				return err
			}
		}
		logger.Log(fmt.Sprintf("Inside CommitTransactionRequest:WritePayload - '%d' removable entities are written in byte format", len(msg.RemovedList)))
	}
	currPos := os.GetPosition()
	length := currPos - startPos
	_, err := os.(*iostream.ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		return err
	}
	logger.Log(fmt.Sprintf("Returning CommitTransactionRequest::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *CommitTransactionRequest) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.VerbId, msg.SequenceNo, msg.Timestamp,
		msg.RequestId, msg.AuthToken, msg.SessionId, msg.DataOffset, msg.IsUpdatable, msg.AddedList, //msg.ChangedList,
		msg.UpdatedList, msg.RemovedList, msg.AttrDescSet)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *CommitTransactionRequest) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.VerbId, &msg.SequenceNo,
		&msg.Timestamp, &msg.RequestId, &msg.AuthToken, &msg.SessionId, &msg.DataOffset, &msg.IsUpdatable,
		&msg.AddedList, &msg.UpdatedList, &msg.RemovedList, &msg.AttrDescSet)
		//&msg.AddedList, &msg.ChangedList, &msg.UpdatedList, &msg.RemovedList, &msg.AttrDescSet)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CommitTransactionRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
