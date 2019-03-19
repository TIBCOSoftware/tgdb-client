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
 * File name: TGInputStream.go
 * Created on: Nov 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
