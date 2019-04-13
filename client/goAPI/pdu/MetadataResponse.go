package pdu

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
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
 * File name: VerbMetadataResponse.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type MetadataResponse struct {
	*AbstractProtocolMessage
	attrDescList []types.TGAttributeDescriptor
	nodeTypeList []types.TGNodeType
	edgeTypeList []types.TGEdgeType
}

func DefaultMetadataResponseMessage() *MetadataResponse {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MetadataResponse{})

	newMsg := MetadataResponse{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.verbId = VerbMetadataResponse
	newMsg.attrDescList = make([]types.TGAttributeDescriptor, 0)
	newMsg.nodeTypeList = make([]types.TGNodeType, 0)
	newMsg.edgeTypeList = make([]types.TGEdgeType, 0)
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewMetadataResponseMessage(authToken, sessionId int64) *MetadataResponse {
	newMsg := DefaultMetadataResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for VerbMetadataResponse
/////////////////////////////////////////////////////////////////

func (msg *MetadataResponse) GetAttrDescList() []types.TGAttributeDescriptor {
	return msg.attrDescList
}

func (msg *MetadataResponse) GetNodeTypeList() []types.TGNodeType {
	return msg.nodeTypeList
}

func (msg *MetadataResponse) GetEdgeTypeList() []types.TGEdgeType {
	return msg.edgeTypeList
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *MetadataResponse) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering MetadataResponse:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside MetadataResponse:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside MetadataResponse:FromBytes - about to APMReadHeader"))
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside MetadataResponse:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("MetadataResponse::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *MetadataResponse) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering MetadataResponse:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside MetadataResponse:ToBytes - about to APMWriteHeader"))
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside MetadataResponse:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("MetadataResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *MetadataResponse) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *MetadataResponse) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *MetadataResponse) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *MetadataResponse) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *MetadataResponse) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *MetadataResponse) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *MetadataResponse) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *MetadataResponse) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *MetadataResponse) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *MetadataResponse) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *MetadataResponse) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *MetadataResponse) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *MetadataResponse) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *MetadataResponse) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *MetadataResponse) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *MetadataResponse) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *MetadataResponse) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *MetadataResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MetadataResponse:{")
	buffer.WriteString(fmt.Sprint("AttrDescSet:{"))
	for _, d := range msg.attrDescList {
		buffer.WriteString(fmt.Sprintf("Attribute Descriptor: %+v ", d))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", NodeTypeList:{"))
	for _, v := range msg.nodeTypeList {
		buffer.WriteString(fmt.Sprintf("NodeType: %+v ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprint(", EdgeTypeList:{"))
	for _, v := range msg.edgeTypeList {
		buffer.WriteString(fmt.Sprintf("EdgeType: %+v ", v))
	}
	buffer.WriteString("}")
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *MetadataResponse) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *MetadataResponse) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *MetadataResponse) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *MetadataResponse) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering MetadataResponse:ReadPayload"))
	avail, err := is.(*iostream.ProtocolDataInputStream).Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading available bytes from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read avail as '%+v'", avail))
	if avail == 0 {
		logger.Log(fmt.Sprintf("Metadata response has no data"))
		errMsg := fmt.Sprint("Metadata Response has no data")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	count, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading count from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read count as '%+v'", count))
	for {
		if count <= 0 {
			break
		}

		sysType, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading sysType from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read sysType as '%+v'", sysType))

		typeCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading typeCount from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read typeCount as '%+v'", typeCount))

		if types.TGSystemType(sysType) == types.SystemTypeAttributeDescriptor {
			for i := 0; i < typeCount; i++ {
				attrDesc := model.NewAttributeDescriptorWithType("temp", types.AttributeTypeString)
				err := attrDesc.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading attrDesc from message buffer"))
					return err
				}
				logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read attrDesc as '%+v'", attrDesc))
				msg.attrDescList = append(msg.attrDescList, attrDesc)
			}
			logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' attrDesc and assigned as '%+v'", typeCount, msg.attrDescList))
		} else if types.TGSystemType(sysType) == types.SystemTypeNode {
			for i := 0; i < typeCount; i++ {
				nodeType := model.NewNodeType("temp", nil)
				err := nodeType.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading node from message buffer"))
					return err
				}
				name := nodeType.GetName()
				if strings.HasPrefix(name, "@") || strings.HasPrefix(name, "$") {
					continue
				}
				logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read nodeType as '%+v'", nodeType))
				msg.nodeTypeList = append(msg.nodeTypeList, nodeType)
			}
			logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' nodes and assigned as '%+v'", typeCount, msg.nodeTypeList))
		} else if types.TGSystemType(sysType) == types.SystemTypeEdge {
			for i := 0; i < typeCount; i++ {
				edgeType := model.NewEdgeType("temp", types.DirectionTypeBiDirectional, nil)
				err := edgeType.ReadExternal(is)
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning MetadataResponse:ReadPayload w/ Error in reading edge from message buffer"))
					return err
				}
				logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read edgeType as '%+v'", edgeType))
				msg.edgeTypeList = append(msg.edgeTypeList, edgeType)
			}
			logger.Log(fmt.Sprintf("Inside MetadataResponse:ReadPayload read '%d' edges and assigned as '%+v'", typeCount, msg.edgeTypeList))
		} else {
			logger.Warning(fmt.Sprintf("WARNING: MetadataResponse:ReadPayload - Invalid metadata desc '%d' received", sysType))
			//TODO: Revisit later - Do we need to throw exception?
		}
		count -= typeCount
	}
	logger.Log(fmt.Sprintf("Returning MetadataResponse:ReadPayload w/ MedataResponse as '%+v'", msg))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *MetadataResponse) WritePayload(os types.TGOutputStream) types.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *MetadataResponse) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.attrDescList,
		msg.nodeTypeList, msg.edgeTypeList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataResponse:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *MetadataResponse) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.attrDescList, &msg.nodeTypeList, &msg.edgeTypeList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MetadataResponse:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
