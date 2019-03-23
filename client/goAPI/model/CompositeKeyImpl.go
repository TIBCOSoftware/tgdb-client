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
 * File name: TGKey.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type CompositeKey struct {
	graphMetadata *GraphMetadata
	keyName       string
	attributes    map[string]types.TGAttribute
}

func NewCompositeKey(graphMetadata *GraphMetadata, typeName string) *CompositeKey {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CompositeKey{})

	newCompositeKey := CompositeKey{
		graphMetadata: graphMetadata,
		keyName:       typeName,
		attributes:    make(map[string]types.TGAttribute, 0),
	}
	return &newCompositeKey
}

/////////////////////////////////////////////////////////////////
// Helper functions for CompositeKey
/////////////////////////////////////////////////////////////////

func (obj *CompositeKey) GetAttributes() map[string]types.TGAttribute {
	return obj.attributes
}

func (obj *CompositeKey) GetKeyName() string {
	return obj.keyName
}

func (obj *CompositeKey) SetAttributes(attrs map[string]types.TGAttribute) {
	obj.attributes = attrs
}

func (obj *CompositeKey) SetKeyName(name string) {
	obj.keyName = name
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGKey
/////////////////////////////////////////////////////////////////

// Dynamically set the attribute to this entity. If the AttributeDescriptor doesn't exist in the database, create a new one.
func (obj *CompositeKey) SetOrCreateAttribute(name string, value interface{}) types.TGError {
	logger.Log(fmt.Sprintf("Entering CompositeKey:SetOrCreateAttribute w/ N-V Pair as '%+v'='%+v'", name, value))
	if name == "" || value == nil {
		logger.Log(fmt.Sprint("ERROR: Returning CompositeKey:SetOrCreateAttribute as either name or value is EMPTY"))
		errMsg := "name or Value is null"
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// If attribute is not present in the set, create a new one
	attr := obj.attributes[name]
	if attr == nil {
		logger.Log(fmt.Sprintf("CompositeKey:SetOrCreateAttribute attribute '%+v' not found - trying to get descriptor from GraphMetadata", name))
		attrDesc, err := obj.graphMetadata.GetAttributeDescriptor(name)
		if err != nil {
			logger.Log(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to get descriptor for attribute '%s' w/ error '%+v'", name, err.Error()))
			return err
		}
		logger.Log(fmt.Sprintf("CompositeKey:SetOrCreateAttribute attribute descriptor for '%+v' found in GraphMetadata", name))
		if attrDesc == nil {
			aType := reflect.TypeOf(value).String()
			logger.Log(fmt.Sprintf("=======> CompositeKey SetOrCreateAttribute creating new attribute descriptor '%+v':'%+v'(%+v) <=======", name, value, aType))
			attrDesc = obj.graphMetadata.CreateAttributeDescriptorForDataType(name, aType)
		}
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc.(*AttributeDescriptor), value)
		if aErr != nil {
			logger.Log(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to create attribute '%s' w/ descriptor and value '%+v'", name, value))
			return aErr
		}
		//newAttr.SetOwner(obj)
		attr = newAttr
	}
	logger.Log(fmt.Sprintf("CompositeKey:SetOrCreateAttribute trying to set attribute '%+v' value as '%+v'", attr, value))
	// Set the attribute value
	err := attr.SetValue(value)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.attributes[name] = attr
	logger.Log(fmt.Sprintf("Returning CompositeKey:SetOrCreateAttribute w/ Key as '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *CompositeKey) ReadExternal(is types.TGInputStream) types.TGError {
	errMsg := "Not Supported operation"
	return exception.GetErrorByType(types.TGErrorIOException, "TGErrorIOException", errMsg, "")
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *CompositeKey) WriteExternal(os types.TGOutputStream) types.TGError {
	if obj.keyName != "" {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(true) //TypeName exists
		err := os.(*iostream.ProtocolDataOutputStream).WriteUTF(obj.keyName)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:WriteExternal - unable to write obj.KeyName w/ Error: '%+v'", err.Error()))
			return err
		}
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteBoolean(false)
	}
	os.(*iostream.ProtocolDataOutputStream).WriteShort(len(obj.attributes))
	for _, attr := range obj.attributes {
		// Null value is not allowed and therefore no need to include isNull flag
		err := attr.WriteExternal(os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:WriteExternal - unable to write attr w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	logger.Log(fmt.Sprintf("Returning CompositeKey:WriteExternal w/ NO error, for key: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *CompositeKey) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.graphMetadata, obj.keyName, obj.attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *CompositeKey) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.graphMetadata, &obj.keyName, &obj.attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
