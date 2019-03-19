package types

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
 * File name: TGOutputStream.go
 * Created on: Nov 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGOutputStream interface {
	// GetBuffer gets the underlying Buffer
	GetBuffer() []byte
	// GetLength gets the total write length
	GetLength() int
	// GetPosition gets the current write position
	GetPosition() int
	// SkipNBytes skips n bytes. Allocate if necessary
	SkipNBytes(n int)
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
	// WriteLongAt writes Long at the position. Buffer should have sufficient space to write the content.
	WriteLongAt(pos int, value int64) (int, TGError)
	// WriteShortAt writes a Java Char at the position. Buffer should have sufficient space to write the content.
	WriteShortAt(pos int, value int) (int, TGError)
	// WriteUTFString writes UTFString
	WriteUTFString(str string) (int, TGError)
	// WriteVarLong writes a long value as varying length into the buffer.
	WriteVarLong(value int64) TGError
}
