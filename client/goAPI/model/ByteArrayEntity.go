package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
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
 * File name: ByteArrayEntity.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ByteArrayEntity struct {
	entityId []byte
}

func DefaultByteArrayEntity() *ByteArrayEntity {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ByteArrayEntity{})

	newByteArrayEntity := ByteArrayEntity{}
	newByteArrayEntity.entityId = make([]byte, 16)
	return &newByteArrayEntity
}

func NewByteArrayEntity(value []byte) *ByteArrayEntity {
	newByteArrayEntity := DefaultByteArrayEntity()
	newByteArrayEntity.entityId = value
	return newByteArrayEntity
}

func NewByteArrayEntityAsInt(value int64) *ByteArrayEntity {
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, value) // Ignore Error Handling
	return NewByteArrayEntity(b.Bytes())
}

/////////////////////////////////////////////////////////////////
// Private functions from Interface ==> types.TGEntityId
/////////////////////////////////////////////////////////////////

func (obj *ByteArrayEntity) equals(cmpObj types.TGEntityId) bool {
	return reflect.DeepEqual(obj.entityId, cmpObj.(*ByteArrayEntity).entityId)
}

func (obj *ByteArrayEntity) writeLongAt(pos int, value int64) int {
	// Splitting in into two ints, is much faster than shifting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
	a := int(value >> 32)
	b := int(value)

	obj.entityId[pos] = byte(a >> 24)
	obj.entityId[pos+1] = byte(a >> 16)
	obj.entityId[pos+2] = byte(a >> 8)
	obj.entityId[pos+3] = byte(a)
	obj.entityId[pos+4] = byte(b >> 24)
	obj.entityId[pos+5] = byte(b >> 16)
	obj.entityId[pos+6] = byte(b >> 8)
	obj.entityId[pos+7] = byte(b)

	return pos + 8
}

func (obj *ByteArrayEntity) GetEntityId() []byte {
	return obj.entityId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGEntityId
/////////////////////////////////////////////////////////////////

// ToBytes converts the Entity Id in binary format
func (obj *ByteArrayEntity) ToBytes() ([]byte, types.TGError) {
	var encodedMsg bytes.Buffer
	encoder := gob.NewEncoder(&encodedMsg)
	err := encoder.Encode(obj)
	//err := interfaceEncode(encoder, obj)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteArrayEntity:ToBytes - unable to write entity in byte format w/ Error: '%+v'", err.Error()))
		errMsg := fmt.Sprintf("Unable to convert message '%+v' in byte format", obj)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Log(fmt.Sprintf("Message: %v bytes <--encoded-- %v", encodedMsg.Len(), m)
	return encodedMsg.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ByteArrayEntity) ReadExternal(is types.TGInputStream) types.TGError {
	entityId, err := is.(*iostream.ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteArrayEntity:ReadExternal - unable to read entityId w/ Error: '%+v'", err.Error()))
		return err
	}
	obj.entityId = entityId
	logger.Log(fmt.Sprintf("Returning ByteArrayEntity:ReadExternal w/ NO error, for entity: '%+v'", obj))
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ByteArrayEntity) WriteExternal(os types.TGOutputStream) types.TGError {
	//logger.Log(fmt.Sprintf("Exported ByteArrayEntity object as '%+v' from byte format", obj))
	return os.(*iostream.ProtocolDataOutputStream).WriteBytes(obj.entityId)
}
