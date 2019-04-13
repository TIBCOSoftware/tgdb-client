package types

/**
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
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
 * File name: TGMessage.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// This is the main interface for all message types
type TGMessage interface {
	// FromBytes constructs a message object from the input buffer in the byte format
	FromBytes(buffer []byte) (TGMessage, TGError)
	// ToBytes converts a message object into byte format to be sent over the network to TGDB server
	ToBytes() ([]byte, int, TGError)
	// GetAuthToken gets the authToken
	GetAuthToken() int64
	// GetIsUpdatable checks whether this message updatable or not
	GetIsUpdatable() bool
	// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
	GetMessageByteBufLength() int
	// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
	GetRequestId() int64
	// GetSequenceNo gets the sequenceNo of the message
	GetSequenceNo() int64
	// GetSessionId gets the session id
	GetSessionId() int64
	// GetTimestamp gets the Timestamp
	GetTimestamp() int64
	// GetVerbId gets verbId of the message
	GetVerbId() int

	// Additional Method to help debugging
	String() string

	// SetAuthToken sets the authToken
	SetAuthToken(authToken int64)
	// SetDataOffset sets the offset at which data starts in the payload
	SetDataOffset(dataOffset int16)
	// SetIsUpdatable sets the updatable flag
	SetIsUpdatable(updateFlag bool)
	// SetMessageByteBufLength sets the message buffer length
	SetMessageByteBufLength(bufLength int)
	// SetRequestId sets the request id
	SetRequestId(requestId int64)
	// SetSequenceNo sets the sequenceNo
	SetSequenceNo(sequenceNo int64)
	// SetSessionId sets the session id
	SetSessionId(sessionId int64)
	// SetTimestamp sets the timestamp
	SetTimestamp(timestamp int64) TGError
	// SetVerbId sets verbId of the message
	SetVerbId(verbId int)

	// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
	// @param timestamp
	// @return TGMessage on success, error on failure
	UpdateSequenceAndTimeStamp(timestamp int64) TGError

	// ReadHeader reads the bytes from input stream and constructs a common header of network packet
	ReadHeader(is TGInputStream) TGError
	// WriteHeader exports the values of the common message header attributes to output stream
	WriteHeader(os TGOutputStream) TGError
	// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
	ReadPayload(is TGInputStream) TGError
	// WritePayload exports the values of the message specific payload attributes to output stream
	WritePayload(os TGOutputStream) TGError
}
