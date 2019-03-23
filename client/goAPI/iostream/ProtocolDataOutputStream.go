package iostream

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"math"
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
 * File name: ProtocolDataOutputStream.go
 * Created on: Nov 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/** Test to ensure float/double gets converted to/from int/int64 correctly
package main

import (
   "fmt"
   "math"
)

func main() {
   fmt.Println("Hello, playground")
   input := 1.23456789
   tmp := math.Float64bits(input)	// For float, use Float32bits()
   output := int64(tmp)
   fmt.Printf("Input: '%+v' Output: '%+v'\n", input, output)
   input1 := tmp
   output1 := math.Float64frombits(input1)	// For float, use Float32frombits()
   fmt.Printf("Input1: '%+v' Output1: '%+v'\n", input1, output1)
}
*/

//var mark = 0
var writeBuf = ""

type ProtocolDataOutputStream struct {
	Buf              []byte
	oStreamBufLen    int
	oStreamByteCount int
}

func DefaultProtocolDataOutputStream() *ProtocolDataOutputStream {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ProtocolDataOutputStream{})

	newStream := ProtocolDataOutputStream{
		Buf:              make([]byte, 256),
		oStreamBufLen:    256,
		oStreamByteCount: 0,
	}
	return &newStream
}

// Create New Input Stream Instance
func NewProtocolDataOutputStream(len int) *ProtocolDataOutputStream {
	newStream := DefaultProtocolDataOutputStream()
	newStream.Buf = make([]byte, len)
	newStream.oStreamBufLen = len
	return newStream
}

/////////////////////////////////////////////////////////////////
// Private functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

func intToBytes(value int, bytes []byte, offset int) {
	for i:=0; i<4; i++ {
		bytes[offset+i] = byte(value >> uint(8*(3-i)))
	}
}

/////////////////////////////////////////////////////////////////
// Helper functions for ProtocolDataInputStream
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) Ensure(len int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::Ensure('%d') in contents", len))
	if (msg.oStreamByteCount + len) <= msg.oStreamBufLen {
		//logger.Log(fmt.Sprint("Returning ProtocolDataOutputStream::Ensure as (msg.Count + len) <= msg.BufLen"))
		return
	}
	newLen := 0
	if len > 100000 {
		newLen = msg.oStreamByteCount + len + 2048
	} else {
		newLen = (msg.oStreamByteCount + len) * 2
	}
	b := make([]byte, newLen)
	copy(b, msg.Buf[0:msg.oStreamByteCount])
	msg.Buf = b
	msg.oStreamBufLen = newLen
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::Ensure('%d') in contents '%+v'", len, msg.Buf))
}

func (msg *ProtocolDataOutputStream) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ProtocolDataOutputStream:{")
	buffer.WriteString(fmt.Sprintf("Buffer Length: %d", msg.oStreamBufLen))
	buffer.WriteString(fmt.Sprintf(", Count: %d", msg.oStreamByteCount))
	buffer.WriteString(fmt.Sprintf(", Buffer: "))
	strArray := []string{buffer.String(), bytes.NewBuffer(msg.Buf).String()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

func (msg *ProtocolDataOutputStream) WriteBoolean(value bool) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBoolean('%+v') in contents", value))
	if msg.oStreamByteCount >= msg.oStreamBufLen {
		msg.Ensure(1)
	}
	if value {
		msg.Buf[msg.oStreamByteCount] = byte(1)
	} else {
		msg.Buf[msg.oStreamByteCount] = byte(0)
	}
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBoolean('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteByte(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteByte('%+v') in contents", value))
	if msg.oStreamByteCount >= msg.oStreamBufLen {
		msg.Ensure(1)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteByte('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteBytesFromString(value string) types.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesFromString('%s') in contents", value))
	buf := []byte(value)
	return msg.WriteBytes(buf)
}

func (msg *ProtocolDataOutputStream) WriteChar(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteChar('%+v') in contents", value))
	if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
		msg.Ensure(2)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteChar('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteChars(value string) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteChars('%+v') in contents", value))
	sLen := len(value)
	msg.Ensure(sLen)
	for i := 0 ; i < sLen ; i++ {
		msg.WriteChar(int(value[i]))
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteChars('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteDouble(value float64) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteDouble('%+v') in contents", value))
	msg.WriteLong(int64(math.Float64bits(value)))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteDouble('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteFloat(value float32) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteFloat('%+v') in contents", value))
	msg.WriteInt(int(math.Float32bits(value)))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteFloat('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteInt(value int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteInt('%+v') in contents", value))
	if (msg.oStreamByteCount + 4) >= msg.oStreamBufLen {
		msg.Ensure(4)
	}
	//intToBytes(value, msg.Buf, msg.Count)
	msg.Buf[msg.oStreamByteCount] = byte(value >> 24)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value >> 16)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteInt('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteLong(value int64) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteLong('%+v') in contents", value))
	if (msg.oStreamByteCount + 8) >= msg.oStreamBufLen {
		msg.Ensure(8)
	}
	// Splitting in into two ints, is much faster than shifting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
	a := int(value >> 32)
	b := int(value)
	//intToBytes(a, msg.Buf, msg.Count)
	msg.Buf[msg.oStreamByteCount] = byte(a >> 24)
	msg.Buf[msg.oStreamByteCount+1] = byte(a >> 16)
	msg.Buf[msg.oStreamByteCount+2] = byte(a >> 8)
	msg.Buf[msg.oStreamByteCount+3] = byte(a)
	//intToBytes(b, msg.Buf, msg.Count+4)
	msg.Buf[msg.oStreamByteCount+4] = byte(b >> 24)
	msg.Buf[msg.oStreamByteCount+5] = byte(b >> 16)
	msg.Buf[msg.oStreamByteCount+6] = byte(b >> 8)
	msg.Buf[msg.oStreamByteCount+7] = byte(b)
	msg.oStreamByteCount += 8
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteLong('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteShort(value int)  {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteShort('%+v') in contents", value))
	if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
		msg.Ensure(2)
	}
	msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
	msg.oStreamByteCount++
	msg.Buf[msg.oStreamByteCount] = byte(value)
	msg.oStreamByteCount++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteShort('%+v') in contents '%+v'", value, msg.Buf))
}

func (msg *ProtocolDataOutputStream) WriteBytesFromPos(value []byte, writePos, writeLen int) types.TGError {
	logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesFromPos('%d') for len '%d' from '%+v' in contents", writePos, writeLen, value))
	if value == nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesFromPos - invalid byte value specified to be written at pos %d", writePos))
		errMsg := fmt.Sprintf("Invalid byte value specified to be written at pos %d", writePos)
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if writePos < 0 || writePos > len(value) || writeLen < 0 || (writePos + writeLen) > len(value) || (writePos + writeLen) < 0 {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesFromPos - either invalid pos %d or write length %d specified", writeLen, writePos))
		errMsg := fmt.Sprintf("Either invalid pos %d or write length %d specified", writeLen, writePos)
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if writeLen == 0 {
		return nil
	}
	msg.Ensure(writeLen)
	tempBuf := value[writePos:]
	copy(msg.Buf[msg.oStreamByteCount:], tempBuf[:writeLen])
	msg.oStreamByteCount += writeLen
	logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytesFromPos(''%d') for len '%d' from '%+v' in contents '%+v'", writePos, writeLen, value, msg.Buf))
	return nil
}

func (msg *ProtocolDataOutputStream) WriteUTF(str string) types.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteUTF('%+v') in contents", str))
	start := msg.oStreamByteCount
	sLen := 0

	if msg.oStreamByteCount+ 2 > msg.oStreamBufLen {
		msg.Ensure(2 + len(str) * 3)
	}
	msg.oStreamByteCount += 2

	sLen, err := msg.WriteUTFString(str)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteUTF - UTF Data Format Issue"))
		msg.oStreamByteCount = start
		errMsg := fmt.Sprint("UTF Data Format Issue")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.GetErrorDetails())
	}

	// Now write length
	msg.Buf[start]   = byte((sLen >> 8) & 0xFF)
	msg.Buf[start+1] = byte(sLen & 0xFF)
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteUTF('%+v') in contents '%+v'", str, msg.Buf))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGInputStream
/////////////////////////////////////////////////////////////////

// GetBuffer gets the underlying Buffer
func (msg *ProtocolDataOutputStream) GetBuffer() []byte {
	return msg.Buf
}

// GetLength gets the total write length
func (msg *ProtocolDataOutputStream) GetLength() int {
	return msg.oStreamByteCount
}

// GetPosition gets the current write position
func (msg *ProtocolDataOutputStream) GetPosition() int {
	return msg.oStreamByteCount
}

// SkipNBytes skips n bytes. Allocate if necessary
func (msg *ProtocolDataOutputStream) SkipNBytes(n int) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::SkipNBytes('%d') in contents", n))
	if msg.oStreamByteCount+ n > msg.oStreamBufLen {
		msg.Ensure(n)
	}
	msg.oStreamByteCount += n
	logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::SkipNBytes('%d') in contents are '%+v'", n, msg))
}

// WriteBooleanAt writes boolean at a given position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteBooleanAt(pos int, value bool) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBooleanAt(Pos '%d' - '%+v') in contents", pos, value))
	if pos >= msg.oStreamByteCount {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBooleanAt - Invalid position '%d' specified", pos))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value {
		msg.Buf[pos] = byte(1)
	} else {
		msg.Buf[pos] = byte(0)
	}
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBooleanAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteByteAt writes a byte at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteByteAt(pos int, value int) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteByteAt(Pos '%d' - '%+v') in contents", pos, value))
	if pos >= msg.oStreamByteCount {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteByteAt - Invalid position '%d' specified", pos))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteByteAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteBytes writes the len, and the byte array into the buffer
func (msg *ProtocolDataOutputStream) WriteBytes(buf []byte) types.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytes('%+v') in contents", buf))
	bLen := len(buf)
	msg.Ensure(bLen+4)
	msg.WriteInt(bLen)

	for i := 0; i < bLen ; i++ {
		msg.WriteByte(int(buf[i]))
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytes('%+v') in contents are '%+v'", buf, msg.Buf))
	return nil
}

// WriteBytesAt writes string at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteBytesAt(pos int, s string) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteBytesAt(Pos '%d' - '%+v') in contents", pos, s))
	sLen := len([]rune(s))
	if (pos + sLen) >= msg.oStreamByteCount {
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	for i := 0; i < sLen ; i++ {
		v, err := fmt.Scanf("%d", s[i])
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesAt - Unable to get character at position %d from string %s", i, s))
			errMsg := fmt.Sprintf("Unable to get character at position %d from string %s", i, s)
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		pos, err := msg.WriteByteAt(pos, v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteBytesAt - Invalid position '%d' specified", pos))
			errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteBytesAt(Pos '%d' - '%+v') in contents are '%+v'", pos, s, msg.Buf))
	return pos, nil
}

// WriteCharAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteCharAt(pos int, value int) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteCharAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 2) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteCharAt as (pos + 2) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteCharAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteCharsAt writes Chars at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteCharsAt(pos int, s string) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteCharsAt(Pos '%d' - '%+v') in contents", pos, s))
	sLen := len([]rune(s))
	if (pos + sLen) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt as (pos + sLen) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	for i := 0; i < sLen ; i++ {
		v, err := fmt.Scanf("%d", s[i])
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt - Unable to get character at position %d from string %s", i, s))
			errMsg := fmt.Sprintf("Unable to get character at position %d from string %s", i, s)
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
		pos, err := msg.WriteCharAt(pos, v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:WriteCharsAt - Invalid position '%d' specified", pos))
			errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
			return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteCharsAt(Pos '%d' - '%+v') in contents are '%+v'", pos, s, msg.Buf))
	return pos, nil
}

// WriteDoubleAt writes Double at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteDoubleAt(pos int, value float64) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteDoubleAt(Pos '%d' - '%+v') in contents", pos, value))
	v, err := msg.WriteLongAt(pos, int64(value))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteDoubleAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return v, err
}

// WriteFloatAt writes Float at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteFloatAt(pos int, value float32) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteFloatAt(Pos '%d' - '%+v') in contents", pos, value))
	v, err := msg.WriteIntAt(pos, int(value))
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteFloatAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return v, err
}

// WriteIntAt writes Integer at the position.Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteIntAt(pos int, value int) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteIntAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 4) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteIntAt as (pos + 4) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 24)
	pos++
	msg.Buf[pos] = byte(value >> 16)
	pos++
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteIntAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteLongAt writes Long at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteLongAt(pos int, value int64) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteLongAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 8) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteLongAt as (pos + 8) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// Splitting in into two ints, is much faster than shifting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
	a := int(value >> 32)
	b := int(value)

	msg.Buf[pos] = byte(a >> 24)
	msg.Buf[pos+1] = byte(a >> 16)
	msg.Buf[pos+2] = byte(a >> 8)
	msg.Buf[pos+3] = byte(a)
	msg.Buf[pos+4] = byte(b >> 24)
	msg.Buf[pos+5] = byte(b >> 16)
	msg.Buf[pos+6] = byte(b >> 8)
	msg.Buf[pos+7] = byte(b)

	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteLongAt(Pos '%d' - '%+v') in contents are '%+v'", pos+8, value, msg.Buf)
	return pos+8, nil
}

// WriteShortAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
func (msg *ProtocolDataOutputStream) WriteShortAt(pos int, value int) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteShortAt(Pos '%d' - '%+v') in contents", pos, value))
	if (pos + 2) >= msg.oStreamByteCount {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteShortAt as (pos + 2) >= msg.Count"))
		errMsg := fmt.Sprintf("Invalid position '%d' specified", pos)
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.Buf[pos] = byte(value >> 8)
	pos++
	msg.Buf[pos] = byte(value)
	pos++
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteShortAt(Pos '%d' - '%+v') in contents are '%+v'", pos, value, msg.Buf))
	return pos, nil
}

// WriteUTFString writes UTFString
func (msg *ProtocolDataOutputStream) WriteUTFString(str string) (int, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteUTFString('%+v') in contents", str))
	start := msg.oStreamByteCount
	sLen := len(str)
	i := 0
	c := 0

	if (msg.oStreamByteCount + 3 * sLen) > msg.oStreamBufLen {
		msg.Ensure(len(str) * 3)
	}
	writeBuf = str[0:sLen]

	for {
		if ! (i < sLen) {
			//logger.Log(fmt.Sprintf("ProtocolDataOutputStream::WriteUTFString - Breaking loop i='%d', c='%d', msg.Count='%d', msg.Buf[msg.Count]='%d'", i, c, msg.Count, msg.Buf[msg.Count]))
			break
		}
		c = int(writeBuf[i])
		i++
		if (c >= 0x0001) && (c <= 0x007F) {
			msg.Buf[msg.oStreamByteCount] = byte(c)
			msg.oStreamByteCount++
		} else if c > 0x07FF {
			msg.Buf[msg.oStreamByteCount] = byte(0xE0 | ((c >> 12) & 0x0F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 6) & 0x3F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 0) & 0x3F))
			msg.oStreamByteCount++
		} else {
			msg.Buf[msg.oStreamByteCount] = byte(0xC0 | ((c >> 6) & 0x1F))
			msg.oStreamByteCount++
			msg.Buf[msg.oStreamByteCount] = byte(0x80 | ((c >> 0) & 0x3F))
			msg.oStreamByteCount++
		}
		//logger.Log(fmt.Sprintf("ProtocolDataOutputStream::WriteUTFString - Inside loop i='%d', c='%d', msg.Count='%d', msg.Buf[msg.Count]='%d'", i, c, msg.Count, msg.Buf[msg.Count]))
	}

	writtenLen := msg.oStreamByteCount - start
	if writtenLen > 65535 {
		msg.oStreamByteCount = start
		errMsg := fmt.Sprint("Input String is too long")
		return -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteUTFString('%+v') in contents are '%+v'", str, msg.Buf))
	return writtenLen, nil
}

// WriteVarLong writes a long value as varying length into the buffer.
func (msg *ProtocolDataOutputStream) WriteVarLong(value int64) types.TGError {
	//logger.Log(fmt.Sprintf("Entering ProtocolDataOutputStream::WriteVarLong('%+v') in contents", value))
	if value < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning ProtocolDataOutputStream:WriteVarLong as value < 0"))
		errMsg := fmt.Sprint("Can not pack negative long value")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value == types.U64_NULL {
		if msg.oStreamByteCount >= msg.oStreamBufLen {
			msg.Ensure(1)
		}
		msg.Buf[msg.oStreamByteCount] = types.U64PACKED_NULL
		msg.oStreamByteCount++
		logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		return nil
	}

	if value <= 0x7f {
		if msg.oStreamByteCount >= msg.oStreamBufLen {
			msg.Ensure(1)
		}
		msg.Buf[msg.oStreamByteCount] = byte(value)
		msg.oStreamByteCount++
		logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		return nil
	}

	if value <= 0x3fff {
		if (msg.oStreamByteCount + 2) >= msg.oStreamBufLen {
			msg.Ensure(2)
		}
		value |= 0x00008000
		msg.Buf[msg.oStreamByteCount] = byte(value >> 8)
		msg.oStreamByteCount++
		msg.Buf[msg.oStreamByteCount] = byte(value)
		msg.oStreamByteCount++
		logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		return nil
	}

	if value <= 0x1fffffff {
		if (msg.oStreamByteCount + 4) >= msg.oStreamBufLen {
			msg.Ensure(4)
		}
		value |= 0xC0000000
		msg.oStreamByteCount += 4
		for i := 1; i < 4; i++ {
			msg.Buf[msg.oStreamByteCount- 1] = byte(value)
			value = value >> 8
		}
		logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
		return nil
	}

	// We may need up to 9 bytes
	if (msg.oStreamByteCount + 9) > msg.oStreamBufLen {
		msg.Ensure(9)
	}

	// Calculate the number of non-zero bytes
	mask := 0xff0000000000000
	cnt := 8
	for i := 0; i < 8; i++ {
		if (value & int64(mask)) != 0 {
			break
		}
		cnt--
		mask = mask >> 8
	}

	b := byte(cnt | 0xE0)
	msg.Buf[msg.oStreamByteCount] = b
	msg.oStreamByteCount++

	cnt += cnt
	for i:=1; i<=cnt; i++ {
		msg.Buf[cnt - i] = byte(value)
		value = value >> 8
	}
	logger.Log(fmt.Sprintf("Returning ProtocolDataOutputStream::WriteVarLong('%+v') in contents are '%+v'", value, msg.Buf))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.Buf, msg.oStreamBufLen, msg.oStreamByteCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (msg *ProtocolDataOutputStream) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.Buf, &msg.oStreamBufLen, &msg.oStreamByteCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ProtocolDataOutputStream:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
