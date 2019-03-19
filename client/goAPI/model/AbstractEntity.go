package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strings"
	"sync/atomic"
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
 * File name: AbstractEntity.go
 * Created on: Oct 06, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
var EntitySequencer int64

type AbstractEntity struct {
	EntityId           int64
	EntityKind         types.TGEntityKind
	EntityType         types.TGEntityType
	IsDeleted          bool
	IsInitialized      bool
	IsNew              bool
	Version            int
	virtualId          int64
	GraphMetadata      *GraphMetadata
	Attributes         map[string]types.TGAttribute
	ModifiedAttributes []types.TGAttribute
}

func DefaultAbstractEntity() *AbstractEntity {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractEntity{})

	newAbstractEntity := AbstractEntity{
		EntityId:           -1,
		//EntityKind:         types.EntityKindInvalid,
		IsDeleted:          false,
		IsNew:              true,
		IsInitialized:      true,
		Version:            0,
		Attributes:         make(map[string]types.TGAttribute, 0),
		ModifiedAttributes: make([]types.TGAttribute, 0),
	}
	newAbstractEntity.virtualId = atomic.AddInt64(&EntitySequencer, -1)
	newAbstractEntity.EntityType = DefaultEntityType()
	return &newAbstractEntity
}

func NewAbstractEntity(gmd *GraphMetadata) *AbstractEntity {
	newAbstractEntity := DefaultAbstractEntity()
	newAbstractEntity.GraphMetadata = gmd
	return newAbstractEntity
}

/////////////////////////////////////////////////////////////////
// Private functions for AbstractEntity / TGEntity - used in all derived entities
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) entityToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractEntity:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %+v", obj.EntityId))
	buffer.WriteString(fmt.Sprintf(", EntityKind: %d", obj.EntityKind))
	buffer.WriteString(fmt.Sprintf(", EntityType: %+v", obj.EntityType))
	buffer.WriteString(fmt.Sprintf(", IsDeleted: %+v", obj.IsDeleted))
	buffer.WriteString(fmt.Sprintf(", IsInitialized: %+v", obj.IsInitialized))
	buffer.WriteString(fmt.Sprintf(", IsNew: %+v", obj.IsNew))
	buffer.WriteString(fmt.Sprintf(", Version: %d", obj.Version))
	buffer.WriteString(fmt.Sprintf(", virtualId: %d", obj.virtualId))
	//buffer.WriteString(fmt.Sprintf(", GraphMetadata: %+v", obj.GraphMetadata))
	buffer.WriteString(fmt.Sprintf(", Attributes: %+v", obj.Attributes))
	buffer.WriteString(fmt.Sprintf(", ModifiedAttributes: %+v", obj.ModifiedAttributes))
	buffer.WriteString("}")
	return buffer.String()
}

func (obj *AbstractEntity) getAttribute(name string) types.TGAttribute {
	attr := obj.Attributes[name]
	return attr
}

func (obj *AbstractEntity) getAttributes() ([]types.TGAttribute, types.TGError) {
	if obj.Attributes == nil {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:getAttributes as there are NO attributes associated"))
		errMsg := "This entity does not have any attributes associated"
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	attrList := make([]types.TGAttribute, 0)
	for attrName, attr := range obj.Attributes {
		if len(attrName) != 0 {
			if attr.GetAttributeDescriptor().GetName() != "" && strings.ToLower(attr.GetAttributeDescriptor().GetName()) != "@name" {
				attrList = append(attrList, attr)
			}
		}
	}
	return attrList, nil
}

func (obj *AbstractEntity) getEntityKind() types.TGEntityKind {
	return obj.EntityKind
}

func (obj *AbstractEntity) getEntityType() types.TGEntityType {
	return obj.EntityType
}

func (obj *AbstractEntity) getGraphMetadata() *GraphMetadata {
	return obj.GraphMetadata
}

func (obj *AbstractEntity) getModifiedAttributes() []types.TGAttribute {
	return obj.ModifiedAttributes
}

func (obj *AbstractEntity) getVersion() int {
	return obj.Version
}

func (obj *AbstractEntity) getVirtualId() int64 {
	if obj.isNew() {
		return obj.virtualId
	}
	return obj.EntityId
}

func (obj *AbstractEntity) isAttributeSet(name string) bool {
	attr := obj.Attributes[name]
	if !attr.IsNull() {
		return true
	}
	return false
}

func (obj *AbstractEntity) isDeleted() bool {
	return obj.IsDeleted
}

func (obj *AbstractEntity) isInitialized() bool {
	return obj.IsInitialized
}

func (obj *AbstractEntity) isNew() bool {
	return obj.IsNew
}

func (obj *AbstractEntity) resetModifiedAttributes() {
	for _, attr := range obj.ModifiedAttributes {
		attr.ResetIsModified()
	}
	// Reset array of modified attributes
	obj.ModifiedAttributes = make([]types.TGAttribute, 0)
}

func (obj *AbstractEntity) setEntityId(id int64) {
	obj.virtualId = 0
	obj.EntityId = id
}

func (obj *AbstractEntity) setEntityKind(kind types.TGEntityKind) {
	obj.EntityKind = kind
}

func (obj *AbstractEntity) setEntityType(eType types.TGEntityType) {
	obj.EntityType = eType
}

func (obj *AbstractEntity) setIsDeleted(flag bool) {
	obj.IsDeleted = flag
}

func (obj *AbstractEntity) setIsInitialized(flag bool) {
	obj.IsInitialized = flag
}

func (obj *AbstractEntity) setIsNew(flag bool) {
	obj.IsNew = flag
}

func (obj *AbstractEntity) setVersion(version int) {
	obj.Version = version
}

func (obj *AbstractEntity) setAttribute(attr types.TGAttribute) types.TGError {
	//attrName := attr.GetName()
	//if attrName == "" {
	//	logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as attrName is EMPTY"))
	//	errMsg := fmt.Sprintf("name of the attribute cannot be null")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	attrDesc := attr.GetAttributeDescriptor()
	if attrDesc == nil {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as attrDesc is EMPTY"))
		errMsg := fmt.Sprintf("Attribute Descriptor cannot be null")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	attrDescName := attrDesc.GetName()
	if attrDescName == "" {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as attrDescName is EMPTY"))
		errMsg := fmt.Sprintf("name of the attribute cannot be null")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.Attributes[attrDescName] = attr
	// Value can be null here
	if !attr.GetIsModified() {
		obj.ModifiedAttributes = append(obj.ModifiedAttributes, attr)
	}
	//logger.Log(fmt.Sprintf("=======> Abstract Entity has attributes '%+v' <=======", obj.attributes))
	return nil
}

// TODO: Revisit later - Once SetAttributeViaDescriptor is properly implemented after discussing with TGDB Engineering Team
func setAttributeViaDescriptor(obj types.TGEntity, attrDesc *AttributeDescriptor, value interface{}) types.TGError {
	if attrDesc == nil {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as attrDesc is EMPTY"))
		errMsg := fmt.Sprintf("Attribute Descriptor cannot be null")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if value == nil {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as value is EMPTY"))
		errMsg := fmt.Sprintf("Attribute value is required")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// TODO: Do we need to validate if this descriptor exists as part of Graph Meta Data???
	// If attribute is not present in the set, create a new one
	attrDescName := attrDesc.GetName()
	attr := obj.(*AbstractEntity).Attributes[attrDescName]
	if attr == nil {
		if attrDesc.GetAttrType() == types.AttributeTypeInvalid {
			logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttributeViaDescriptor as attrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			errMsg := fmt.Sprintf("Attribute descriptor is of incorrect type")
			return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
		// TODO: Revisit later - For some reason, it goes in infinite loop, alternative is to create attribute and assign owner and value later
		//newAttr, aErr := CreateAttribute(obj, attrDesc, value)
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc, nil)
		if aErr != nil {
			logger.Log(fmt.Sprintf("ERROR: Returning AbstractEntity:setAttributeViaDescriptor unable to create attribute '%s' w/ descriptor and value '%+v'", attrDesc, value))
			return aErr
		}
		newAttr.SetOwner(obj)
		attr = newAttr
	}
	// Value can be null here
	if !attr.GetIsModified() {
		obj.(*AbstractEntity).ModifiedAttributes = append(obj.(*AbstractEntity).ModifiedAttributes, attr)
	}
	// Set the attribute value
	err := attr.SetValue(value)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: Returning AbstractEntity:setAttributeViaDescriptor unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.(*AbstractEntity).Attributes[attrDesc.name] = attr
	return nil
}

func (obj *AbstractEntity) setOrCreateAttribute(name string, value interface{}) types.TGError {
	logger.Log(fmt.Sprintf("Entering AbstractEntity:SetOrCreateAttribute received N-V Pair as '%+v'='%+v'", name, value))
	if name == "" || value == nil {
		logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:SetOrCreateAttribute as either name or value is EMPTY"))
		errMsg := fmt.Sprintf("name of the attribute cannot be null")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// If attribute is not present in the set, create a new one
	attr := obj.GetAttribute(name)
	//logger.Log(fmt.Sprintf("Inside AbstractEntity:SetOrCreateAttribute Abstract Entity has attribute '%+v' <=======", obj.attributes))
	if attr == nil {
		gmd := obj.GetGraphMetadata()
		attrDesc, err := gmd.GetAttributeDescriptor(name)
		if err != nil {
			logger.Log(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to get descriptor for attribute '%s' w/ error '%+v'", name, err.Error()))
			return err
		}
		if attrDesc == nil {
			if value == nil {
				errMsg := fmt.Sprintf("Attribute value is required")
				return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
			}
			aType := reflect.TypeOf(value).String()
			logger.Log(fmt.Sprintf("Inside AbstractEntity:SetOrCreateAttribute Abstract Entity creating new attribute '%+v':'%+v'(%+v) <=======", name, value, aType))
			// TODO: Do we need to validate if this descriptor exists as part of Graph Meta Data???
			attrDesc = gmd.CreateAttributeDescriptorForDataType(name, aType)
		}
		// TODO: Revisit later - For some reason, it goes in infinite loop, alternative is to create attribute and assign owner and value later
		//newAttr, aErr := CreateAttribute(obj, attrDesc, value)
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc.(*AttributeDescriptor), nil)
		if aErr != nil {
			logger.Log(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to create attribute '%s' w/ descriptor and value '%+v'", name, value))
			return aErr
		}
		newAttr.SetOwner(obj)
		attr = newAttr
	}
	// Value can be null here
	if !attr.GetIsModified() {
		obj.ModifiedAttributes = append(obj.ModifiedAttributes, attr)
	}
	// Set the attribute value
	err := attr.SetValue(value)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.Attributes[name] = attr
	logger.Log(fmt.Sprintf("Returning AbstractEntity:SetOrCreateAttribute created/set attribute for '%+v'='%+v'", name, value))
	return nil
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGAbstractEntity
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) GetIsInitialized() bool {
	return obj.isInitialized()
}

func (obj *AbstractEntity) GetModifiedAttributes() []types.TGAttribute {
	return obj.getModifiedAttributes()
}

// TODO: Revisit later - Once SetAttributeViaDescriptor is properly implemented after discussing with TGDB Engineering Team
func (obj *AbstractEntity) SetAttributeViaDescriptor(attrDesc *AttributeDescriptor, value interface{}) types.TGError {
	//return setAttributeViaDescriptor(obj, attrDesc, value)
	if attrDesc == nil {
		errMsg := "Attribute descriptor is required"
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return nil
}

func (obj *AbstractEntity) SetEntityKind(kind types.TGEntityKind) {
	obj.setEntityKind(kind)
}

func (obj *AbstractEntity) SetEntityType(eType types.TGEntityType) {
	obj.setEntityType(eType)
}

func (obj *AbstractEntity) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

func (obj *AbstractEntity) AbstractEntityReadExternal(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AbstractEntity:AbstractEntityReadExternal"))
	newEntityFlag, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean() // Should always be False.
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read newEntityFlag w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read newEntityFlag as '%+v'", newEntityFlag))
	if newEntityFlag {
		//TGDB-504
		logger.Warning(fmt.Sprint("WARNING: AbstractEntity:AbstractEntityReadExternal - de-serializing a new entity is NOT expected"))
		obj.SetIsNew(false)
	}

	eKind, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read eKind w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read eKind as '%+v'", eKind))
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal - Object Kind is '%+v'", obj.EntityKind))
	if obj.GetEntityKind() != types.TGEntityKind(eKind) {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractEntity:AbstractEntityReadExternal as obj.GetEntityKind() != types.TGEntityKind(eKind)"))
		errMsg := "Invalid object for deserialization. Expecting..." // TODO: SS
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Overwrite the entityId as set by Server
	entityId, err := is.(*iostream.ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read entityId w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read entityId as '%d'", entityId))

	version, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read version w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read version as '%d'", version))

	var eType types.TGEntityType
	entityTypeId, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read entityTypeId w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read entityTypeId as '%d'", entityTypeId))
	if entityTypeId != 0 {
		eType1, err := obj.GraphMetadata.GetNodeTypeById(entityTypeId)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read nodeType w/ Error: '%+v'", err.Error()))
			return err
		}
		logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal verified eType1 (nodeTypeById) as '%+v'", eType1))
		if eType1 == nil {
			eType, err = obj.GraphMetadata.GetEdgeTypeById(entityTypeId)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read edgeType w/ Error: '%+v'", err.Error()))
				return err
			}
			logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal verified eType (edgeTypeById) as '%+v'", eType))
			if eType == nil {
				// TODO: Revisit later - Should we retrieve entity desc together with the entity?
				logger.Warning(fmt.Sprintf("WARNING: Cannot lookup entity desc '%d' from graph meta data cache", entityTypeId))
			}
		} else {
			eType = eType1
		}
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal inferred from metadata, eType as '%+v'", eType))

	count, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read count w/ Error: '%+v'", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read attribute count as '%d'", count))
	for i := 0; i < count; i++ {
		attr, err := ReadExternalForEntity(obj, is)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read attr w/ Error: '%+v'", err.Error()))
			return err
		}
		logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read attr as '%+v'", attr))
		err = obj.SetAttribute(attr)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to set attr w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	obj.EntityId = entityId
	obj.SetEntityType(eType)
	obj.SetIsNew(newEntityFlag)
	obj.SetVersion(version)
	logger.Log(fmt.Sprintf("Returning AbstractEntity:AbstractEntityReadExternal w/ NO error, for entity: '%+v'", obj))
	return nil
}

func (obj *AbstractEntity) AbstractEntityWriteExternal(os types.TGOutputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering AbstractEntity:AbstractEntityWriteExternal"))
	os.(*iostream.ProtocolDataOutputStream).WriteBoolean(obj.GetIsNew())
	os.(*iostream.ProtocolDataOutputStream).WriteByte(int(obj.GetEntityKind())) //Write the EntityKind
	//virtual id can be local or actual id
	os.(*iostream.ProtocolDataOutputStream).WriteLong(obj.GetVirtualId())
	os.(*iostream.ProtocolDataOutputStream).WriteInt(obj.GetVersion())
	logger.Log(fmt.Sprintf("Inside AbstractEntity:AbstractEntityWriteExternal - obj.EntityType is '%+v'", obj.GetEntityType()))
	if obj.GetEntityType() == nil {
		os.(*iostream.ProtocolDataOutputStream).WriteInt(0)
	} else {
		os.(*iostream.ProtocolDataOutputStream).WriteInt(obj.GetEntityType().GetEntityTypeId())
	}

	// The attribute id can be temporary which is a negative number
	// The actual attribute id is > 0
	modCount := 0
	for _, attr := range obj.Attributes {
		if attr.GetIsModified() {
			modCount++
		}
	}
	os.(*iostream.ProtocolDataOutputStream).WriteInt(modCount)
	for _, attr := range obj.Attributes {
		// If an attribute is not modified, do not include in the stream
		if !attr.GetIsModified() {
			logger.Warning(fmt.Sprint("WARNING: Continuing loop AbstractEntity:AbstractEntityWriteExternal as attr is NOT modified"))
			continue
		}
		err := attr.WriteExternal(os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityWriteExternal - unable to write attr w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	logger.Log(fmt.Sprintf("Returning AbstractEntity:AbstractEntityWriteExternal w/ NO error, for entity: '%+v'", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the name specified
func (obj *AbstractEntity) GetAttribute(attrName string) types.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the attributes set
func (obj *AbstractEntity) GetAttributes() ([]types.TGAttribute, types.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *AbstractEntity) GetEntityKind() types.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *AbstractEntity) GetEntityType() types.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *AbstractEntity) GetGraphMetadata() types.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *AbstractEntity) GetIsDeleted() bool {
	return obj.isDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *AbstractEntity) GetIsNew() bool {
	return obj.isNew()
}

// GetVersion gets the version of the Entity
func (obj *AbstractEntity) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *AbstractEntity) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *AbstractEntity) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on attributes
func (obj *AbstractEntity) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *AbstractEntity) SetAttribute(attr types.TGAttribute) types.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *AbstractEntity) SetOrCreateAttribute(name string, value interface{}) types.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *AbstractEntity) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *AbstractEntity) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *AbstractEntity) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *AbstractEntity) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *AbstractEntity) String() string {
	return obj.entityToString()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *AbstractEntity) ReadExternal(is types.TGInputStream) types.TGError {
	return obj.AbstractEntityReadExternal(is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AbstractEntity) WriteExternal(os types.TGOutputStream) types.TGError {
	return obj.AbstractEntityWriteExternal(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.IsNew, obj.EntityKind, obj.virtualId, obj.Version, obj.EntityId, obj.EntityType,
		obj.IsDeleted, obj.IsInitialized, obj.GraphMetadata, obj.Attributes, obj.ModifiedAttributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.IsNew, &obj.EntityKind, &obj.virtualId, &obj.Version, &obj.EntityId, &obj.EntityType,
		&obj.IsDeleted, &obj.IsInitialized, &obj.GraphMetadata, &obj.Attributes, &obj.ModifiedAttributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
