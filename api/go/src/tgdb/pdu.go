/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name: pdu.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: pdu.go 4575 2020-10-27 00:21:18Z nimish $
 */

package tgdb


type TGSerializable interface {
	// ReadExternal reads the byte format from an external input stream and constructs a system object
	ReadExternal(iStream TGInputStream) TGError
	// WriteExternal writes a system object into an appropriate byte format onto an external output stream
	WriteExternal(oStream TGOutputStream) TGError
}


type TGInputStream interface {
	// Available checks whether there is any data available on the stream to read
	Available() (int, TGError)
	// GetPosition gets the current position of internal cursor
	GetPosition() int64
	// GetReferenceMap returns a user maintained reference map
	GetReferenceMap() map[int64]TGEntity
	// Mark marks the current position
	Mark(readlimit int)
	// MarkSupported checks whether the marking is supported or not
	MarkSupported() bool
	// Read reads the current byte
	Read() (int, TGError)
	// ReadIntoBuffer copies bytes in specified buffer
	// The buffer cannot be NIL
	ReadIntoBuffer(b []byte) (int, TGError)
	// ReadAtOffset is similar to readFully.
	ReadAtOffset(b []byte, off int, length int) (int, TGError)
	// ReadBytes reads an encoded byte array. writeBytes encodes the length, and the byte[].
	// This is equivalent to do a readInt, and read(byte[])
	ReadBytes() ([]byte, TGError)
	// ReadVarLong reads a Variable long field
	ReadVarLong() (int64, TGError)
	// Reset brings internal moving cursor back to the old position
	Reset()
	// SetPosition sets the position of reading.
	SetPosition(position int64) int64
	// SetReferenceMap sets a user maintained map as reference data
	SetReferenceMap(rMap map[int64]TGEntity)
	// Skip skips n bytes
	Skip(n int64) (int64, TGError)
}

type TGOutputStream interface {
	// GetBuffer gets the underlying Buffer
	GetBuffer() []byte
	// GetLength gets the total write length
	GetLength() int
	// GetPosition gets the current write position
	GetPosition() int
	// SkipNBytes skips n bytes. Allocate if necessary
	SkipNBytes(n int)
	// ToByteArray returns a new constructed byte array of the data that is being streamed.
	ToByteArray() ([]byte, TGError)
	// WriteBooleanAt writes boolean at a given position. Buffer should have sufficient space to write the content.
	WriteBooleanAt(pos int, value bool) (int, TGError)
	// WriteByteAt writes a byte at the position. Buffer should have sufficient space to write the content.
	WriteByteAt(pos int, value int) (int, TGError)
	// WriteBytes writes the len, and the byte array into the buffer
	WriteBytes(buf []byte) TGError
	// WriteBytesAt writes string at the position. Buffer should have sufficient space to write the content.
	WriteBytesAt(pos int, s string) (int, TGError)
	// WriteCharAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
	WriteCharAt(pos int, value int) (int, TGError)
	// WriteCharsAt writes Chars at the position. Buffer should have sufficient space to write the content.
	WriteCharsAt(pos int, s string) (int, TGError)
	// WriteDoubleAt writes Double at the position. Buffer should have sufficient space to write the content.
	WriteDoubleAt(pos int, value float64) (int, TGError)
	// WriteFloatAt writes Float at the position. Buffer should have sufficient space to write the content.
	WriteFloatAt(pos int, value float32) (int, TGError)
	// WriteIntAt writes Integer at the position.Buffer should have sufficient space to write the content.
	WriteIntAt(pos int, value int) (int, TGError)
	// WriteLongAsBytes writes Long in byte format
	WriteLongAsBytes(value int64) TGError
	// WriteLongAt writes Long at the position. Buffer should have sufficient space to write the content.
	WriteLongAt(pos int, value int64) (int, TGError)
	// WriteShortAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
	WriteShortAt(pos int, value int) (int, TGError)
	// WriteUTFString writes UTFString
	WriteUTFString(str string) (int, TGError)
	// WriteVarLong writes a long value as varying length into the buffer.
	WriteVarLong(value int64) TGError
}

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
	// GetTenantId gets the tenantId of the message
	GetTenantId() int
	// SetTenantId sets the tenantId
    SetTenantId(id int)
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
