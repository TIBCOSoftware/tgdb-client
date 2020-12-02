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
 * File Name: entityimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: entityimpl.go 4543 2020-10-22 18:36:16Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"tgdb"
)

type AbstractEntity struct {
	EntityId           int64
	EntityKind         tgdb.TGEntityKind
	EntityType         tgdb.TGEntityType
	isDeleted          bool
	isInitialized      bool
	isNew              bool
	Version            int
	VirtualId          int64
	graphMetadata      *GraphMetadata
	Attributes         map[string]tgdb.TGAttribute
	ModifiedAttributes []tgdb.TGAttribute
}

var EntitySequencer int64



func DefaultAbstractEntity() *AbstractEntity {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractEntity{})

	newAbstractEntity := AbstractEntity{
		EntityId: -1,
		//EntityKind:         types.EntityKindInvalid,
		isDeleted:          false,
		isNew:              true,
		isInitialized:      true,
		Version:            0,
		Attributes:         make(map[string]tgdb.TGAttribute, 0),
		ModifiedAttributes: make([]tgdb.TGAttribute, 0),
	}
	newAbstractEntity.VirtualId = atomic.AddInt64(&EntitySequencer, -1)
	newAbstractEntity.EntityType = DefaultEntityType()
	return &newAbstractEntity
}

func NewAbstractEntity(gmd *GraphMetadata) *AbstractEntity {
	newAbstractEntity := DefaultAbstractEntity()
	newAbstractEntity.graphMetadata = gmd
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
	buffer.WriteString(fmt.Sprintf(", IsDeleted: %+v", obj.isDeleted))
	buffer.WriteString(fmt.Sprintf(", IsInitialized: %+v", obj.isInitialized))
	buffer.WriteString(fmt.Sprintf(", IsNew: %+v", obj.isNew))
	buffer.WriteString(fmt.Sprintf(", Version: %d", obj.Version))
	buffer.WriteString(fmt.Sprintf(", virtualId: %d", obj.VirtualId))
	//buffer.WriteString(fmt.Sprintf(", GraphMetadata: %+v", obj.GraphMetadata))
	buffer.WriteString(fmt.Sprintf(", Attributes: %+v", obj.Attributes))
	buffer.WriteString(fmt.Sprintf(", ModifiedAttributes: %+v", obj.ModifiedAttributes))
	buffer.WriteString("}")
	return buffer.String()
}

func (obj *AbstractEntity) getAttribute(name string) tgdb.TGAttribute {
	attr := obj.Attributes[name]
	return attr
}

func (obj *AbstractEntity) getAttributes() ([]tgdb.TGAttribute, tgdb.TGError) {
	if obj.Attributes == nil {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("ERROR: Returning AbstractEntity:getAttributes as there are NO Attributes associated"))
		}
		errMsg := "This entity does not have any Attributes associated"
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	attrList := make([]tgdb.TGAttribute, 0)
	for attrName, attr := range obj.Attributes {
		if len(attrName) != 0 {
			if attr.GetAttributeDescriptor().GetName() != "" && strings.ToLower(attr.GetAttributeDescriptor().GetName()) != "@Name" {
				attrList = append(attrList, attr)
			}
		}
	}
	return attrList, nil
}

func (obj *AbstractEntity) getEntityKind() tgdb.TGEntityKind {
	return obj.EntityKind
}

func (obj *AbstractEntity) getEntityType() tgdb.TGEntityType {
	return obj.EntityType
}

func (obj *AbstractEntity) getGraphMetadata() *GraphMetadata {
	return obj.graphMetadata
}

func (obj *AbstractEntity) getModifiedAttributes() []tgdb.TGAttribute {
	return obj.ModifiedAttributes
}

func (obj *AbstractEntity) getVersion() int {
	return obj.Version
}

func (obj *AbstractEntity) getVirtualId() int64 {
	if obj.getIsNew() {
		return obj.VirtualId
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

func (obj *AbstractEntity) getIsDeleted() bool {
	return obj.isDeleted
}

func (obj *AbstractEntity) getIsInitialized() bool {
	return obj.isInitialized
}

func (obj *AbstractEntity) getIsNew() bool {
	return obj.isNew
}

func (obj *AbstractEntity) resetModifiedAttributes() {
	for _, attr := range obj.ModifiedAttributes {
		attr.ResetIsModified()
	}
	// Reset array of modified Attributes
	obj.ModifiedAttributes = make([]tgdb.TGAttribute, 0)
}

func (obj *AbstractEntity) setAttributes(attrs map[string]tgdb.TGAttribute) {
	obj.Attributes = attrs
}

func (obj *AbstractEntity) setEntityId(id int64) {
	if obj.getIsNew() {
		obj.VirtualId = id
	} else {
		obj.VirtualId = 0
	}
	obj.EntityId = id
}

func (obj *AbstractEntity) setEntityKind(kind tgdb.TGEntityKind) {
	obj.EntityKind = kind
}

func (obj *AbstractEntity) setEntityType(eType tgdb.TGEntityType) {
	obj.EntityType = eType
}

func (obj *AbstractEntity) setIsDeleted(flag bool) {
	obj.isDeleted = flag
}

func (obj *AbstractEntity) setIsInitialized(flag bool) {
	obj.isInitialized = flag
}

func (obj *AbstractEntity) setIsNew(flag bool) {
	obj.isNew = flag
}

func (obj *AbstractEntity) setModifiedAttributes(mAttrs []tgdb.TGAttribute) {
	obj.ModifiedAttributes = mAttrs
}

func (obj *AbstractEntity) setVersion(version int) {
	obj.Version = version
}

func (obj *AbstractEntity) setAttribute(attr tgdb.TGAttribute) tgdb.TGError {
	//attrName := attr.GetName()
	//if attrName == "" {
	//	logger.Log(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as attrName is EMPTY"))
	//	errMsg := fmt.Sprintf("Name of the attribute cannot be null")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	//}
	attrDesc := attr.GetAttributeDescriptor()
	if attrDesc == nil {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as AttrDesc is EMPTY"))
		}
		errMsg := fmt.Sprintf("Attribute Descriptor cannot be null")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	attrDescName := attrDesc.GetName()
	if attrDescName == "" {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("ERROR: Returning AbstractEntity:setAttribute as attrDescName is EMPTY"))
		}
		errMsg := fmt.Sprintf("Name of the attribute cannot be null")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.Attributes[attrDescName] = attr
	// Value can be null here
	if !attr.GetIsModified() {
		obj.ModifiedAttributes = append(obj.ModifiedAttributes, attr)
	}
	//logger.Log(fmt.Sprintf("=======> Abstract Entity has Attributes '%+v' <=======", obj.Attributes))
	return nil
}


func (obj *AbstractEntity) setOrCreateAttribute(name string, value interface{}) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering AbstractEntity:SetOrCreateAttribute received N-V Pair as '%+v'='%+v'", name, value))
	}
	if name == "" || value == nil {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprint("ERROR: Returning AbstractEntity:SetOrCreateAttribute as either Name or value is EMPTY"))
		}
		errMsg := fmt.Sprintf("Name of the attribute cannot be null")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// If attribute is not present in the set, create a new one
	attr := obj.GetAttribute(name)
	//logger.Debug(fmt.Sprintf("Inside AbstractEntity:SetOrCreateAttribute Abstract Entity has attribute '%+v' <=======", obj.Attributes))
	if attr == nil {
		gmd := obj.GetGraphMetadata()
		attrDesc, err := gmd.GetAttributeDescriptor(name)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to get descriptor for attribute '%s' w/ error '%+v'", name, err.Error()))
			return err
		}
		if attrDesc == nil {
			if value == nil {
				errMsg := fmt.Sprintf("Attribute value is required")
				return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			aType := reflect.TypeOf(value).String()
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside AbstractEntity:SetOrCreateAttribute Abstract Entity creating new attribute '%+v':'%+v'(%+v) <=======", name, value, aType))
			}
			// TODO: Do we need to validate if this descriptor exists as part of Graph Meta Data???
			attrDesc = gmd.CreateAttributeDescriptorForDataType(name, aType)
		}
		// TODO: Revisit later - For some reason, it goes in infinite loop, alternative is to create attribute and assign owner and value later
		//newAttr, aErr := CreateAttribute(obj, AttrDesc, value)
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc.(*AttributeDescriptor), nil)
		if aErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to create attribute '%s' w/ descriptor and value '%+v'", name, value))
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
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:SetOrCreateAttribute unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.Attributes[name] = attr
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractEntity:SetOrCreateAttribute created/set attribute for '%+v'='%+v'", name, value))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGAbstractEntity
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) GetIsInitialized() bool {
	return obj.getIsInitialized()
}

func (obj *AbstractEntity) GetModifiedAttributes() []tgdb.TGAttribute {
	return obj.getModifiedAttributes()
}

func (obj *AbstractEntity) SetAttributes(modAttrs map[string]tgdb.TGAttribute) {
	obj.setAttributes(modAttrs)
}

// TODO: Revisit later - Once SetAttributeViaDescriptor is properly implemented after discussing with TGDB Engineering Team
func (obj *AbstractEntity) SetAttributeViaDescriptor(attrDesc *AttributeDescriptor, value interface{}) tgdb.TGError {
	//return setAttributeViaDescriptor(obj, AttrDesc, value)
	if attrDesc == nil {
		errMsg := "Attribute descriptor is required"
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return nil
}

func (obj *AbstractEntity) SetEntityKind(kind tgdb.TGEntityKind) {
	obj.setEntityKind(kind)
}

func (obj *AbstractEntity) SetEntityType(eType tgdb.TGEntityType) {
	obj.setEntityType(eType)
}

func (obj *AbstractEntity) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

func (obj *AbstractEntity) SetModifiedAttributes(modAttrs []tgdb.TGAttribute) {
	obj.setModifiedAttributes(modAttrs)
}

func (obj *AbstractEntity) AbstractEntityReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractEntity:AbstractEntityReadExternal"))
	}
	newEntityFlag, err := is.(*ProtocolDataInputStream).ReadBoolean() // Should always be False.
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read newEntityFlag w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read newEntityFlag as '%+v'", newEntityFlag))
	}
	if newEntityFlag {
		//TGDB-504
		//logger.Warning(fmt.Sprint("WARNING: AbstractEntity:AbstractEntityReadExternal - de-serializing a new entity is NOT expected"))
		obj.SetIsNew(false)
	}

	eKind, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read eKind w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read eKind as '%+v'", eKind))
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal - Object Kind is '%+v'", obj.EntityKind))
	}
	if obj.GetEntityKind() != tgdb.TGEntityKind(eKind) {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractEntity:AbstractEntityReadExternal as obj.GetEntityKind() != types.TGEntityKind(eKind)"))
		errMsg := "Invalid object for deserialization. Expecting..." // TODO: SS
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Overwrite the entityId as set by Server
	entityId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read entityId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read entityId as '%d'", entityId))
	}

	version, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read version w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read version as '%d'", version))
	}

	var eType tgdb.TGEntityType
	entityTypeId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read entityTypeId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read entityTypeId as '%d'", entityTypeId))
	}
	if entityTypeId != 0 {
		eType1, err := obj.graphMetadata.GetNodeTypeById(entityTypeId)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read nodeType w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal verified eType1 (nodeTypeById) as '%+v'", eType1))
		}
		if eType1 == nil {
			eType, err = obj.graphMetadata.GetEdgeTypeById(entityTypeId)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read edgeType w/ Error: '%+v'", err.Error()))
				return err
			}
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal verified eType (edgeTypeById) as '%+v'", eType))
			}
			if eType == nil {
				// TODO: Revisit later - Should we retrieve entity desc together with the entity?
				logger.Warning(fmt.Sprintf("WARNING: Cannot lookup entity desc '%d' from graph meta data cache", entityTypeId))
			}
		} else {
			eType = eType1
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal inferred from metadata, eType as '%+v'", eType))
	}

	count, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read count w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read attribute count as '%d'", count))
	}
	for i := 0; i < count; i++ {
		attr, err := ReadExternalForEntity(obj, is)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read attr w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityReadExternal read attr as '%+v'", attr))
		}
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
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractEntity:AbstractEntityReadExternal w/ NO error, for entity: '%+v'", obj))
	}
	return nil
}

func (obj *AbstractEntity) AbstractEntityWriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractEntity:AbstractEntityWriteExternal"))
	}
	os.(*ProtocolDataOutputStream).WriteBoolean(obj.GetIsNew())
	os.(*ProtocolDataOutputStream).WriteByte(int(obj.GetEntityKind())) //Write the EntityKind
	//virtual id can be local or actual id
	os.(*ProtocolDataOutputStream).WriteLong(obj.GetVirtualId())
	os.(*ProtocolDataOutputStream).WriteInt(obj.GetVersion())
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractEntity:AbstractEntityWriteExternal - obj.EntityType is '%+v'", obj.GetEntityType()))
	}
	if obj.GetEntityType() == nil {
		os.(*ProtocolDataOutputStream).WriteInt(0)
	} else {
		os.(*ProtocolDataOutputStream).WriteInt(obj.GetEntityType().GetEntityTypeId())
	}

	// The attribute id can be temporary which is a negative number
	// The actual attribute id is > 0
	modCount := 0
	for _, attr := range obj.Attributes {
		if attr.GetIsModified() {
			modCount++
		}
	}
	os.(*ProtocolDataOutputStream).WriteInt(modCount)
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
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractEntity:AbstractEntityWriteExternal w/ NO error, for entity: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the Name specified
func (obj *AbstractEntity) GetAttribute(attrName string) tgdb.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the Attributes set
func (obj *AbstractEntity) GetAttributes() ([]tgdb.TGAttribute, tgdb.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *AbstractEntity) GetEntityKind() tgdb.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *AbstractEntity) GetEntityType() tgdb.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *AbstractEntity) GetGraphMetadata() tgdb.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *AbstractEntity) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *AbstractEntity) GetIsNew() bool {
	return obj.getIsNew()
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

// ResetModifiedAttributes resets the dirty flag on Attributes
func (obj *AbstractEntity) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *AbstractEntity) SetAttribute(attr tgdb.TGAttribute) tgdb.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *AbstractEntity) SetOrCreateAttribute(name string, value interface{}) tgdb.TGError {
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
func (obj *AbstractEntity) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return obj.AbstractEntityReadExternal(is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AbstractEntity) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return obj.AbstractEntityWriteExternal(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractEntity) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.VirtualId, obj.Version, obj.EntityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.Attributes, obj.ModifiedAttributes)
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
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.VirtualId, &obj.Version, &obj.EntityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.Attributes, &obj.ModifiedAttributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

type EntityType struct {
	SysType    tgdb.TGSystemType
	id         int // Issued only for creation and not valid later
	Name       string
	parent     tgdb.TGEntityType
	Attributes map[string]*AttributeDescriptor
}

func DefaultEntityType() *EntityType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(EntityType{})

	newEntityType := EntityType{
		//SysType:    types.SystemTypeEntity,
		Attributes: make(map[string]*AttributeDescriptor, 0),
	}
	// TODO: Check with TGDB Engineering Team how and when Id will be set - It is supposed to be set at the time of creation
	return &newEntityType
}

func NewEntityType(name string, parent *EntityType) *EntityType {
	newEntityType := DefaultEntityType()
	newEntityType.Name = name
	newEntityType.parent = parent
	return newEntityType
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGEntityType
/////////////////////////////////////////////////////////////////

func (obj *EntityType) UpdateMetadata(gmd *GraphMetadata) tgdb.TGError {
	return EntityTypeUpdateMetadata(obj, gmd)
}

func (obj *EntityType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.Attributes = attrMap
}

func (obj *EntityType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.Attributes[attrName] = attrDesc
}

func (obj *EntityType) SetParent(parentEntity tgdb.TGEntityType) {
	obj.parent = parentEntity
}

func EntityTypeReadExternal(obj tgdb.TGEntityType, is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering EntityType:EntityTypeReadExternal"))
	}
	// TODO: Revisit later - Do we save the desc value?
	sType, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read sType w/ Error: '%+v'", err.Error()))
		return err
	}
	if tgdb.TGSystemType(sType) == tgdb.SystemTypeInvalid {
		logger.Warning(fmt.Sprint("WARNING: EntityType:EntityTypeReadExternal - types.TGSystemType(sType) == types.SystemTypeInvalid"))
		// TODO: Revisit later - Do we need to throw Exception?
		//errMsg := fmt.Sprintf("Entity desc input stream has invalid desc value: %d", sType)
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	eId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read eId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EntityType:EntityTypeReadExternal read eId as '%d'", eId))
	}

	eName, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read eName w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EntityType:EntityTypeReadExternal read eName as '%s'", eName))
	}

	_, err = is.(*ProtocolDataInputStream).ReadInt() // pagesize
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read pageSize w/ Error: '%+v'", err.Error()))
		return err
	}

	// TODO: Check with TGDB Engineering Team why parent is not being sent over

	attrCount, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read attrCount w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EntityType:EntityTypeReadExternal read attrCount as '%d'", attrCount))
	}

	for i := 0; i < int(attrCount); i++ {
		attrName, err := is.(*ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeReadExternal - unable to read aName w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside EntityType:EntityTypeReadExternal read attrName as '%s'", attrName))
		}
		// TODO: The stream only contains Name of the descriptor. Do we need to lookup the descriptor from GraphMetaData?
		attrDesc := NewAttributeDescriptorWithType(attrName, AttributeTypeString)
		obj.AddAttributeDescriptor(attrName, attrDesc)
	}
	obj.SetEntityTypeId(eId)
	obj.SetName(eName)
	obj.SetSystemType(tgdb.TGSystemType(sType))
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning EntityType:EntityTypeReadExternal w/ NO error, for entityType: '%+v'", obj))
	}
	return nil
}

func EntityTypeUpdateMetadata(obj tgdb.TGEntityType, gmd *GraphMetadata) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering EntityType:EntityTypeUpdateMetadata"))
	}
	for attrName, _ := range obj.(*EntityType).Attributes {
		attrDesc, err := gmd.GetAttributeDescriptor(attrName)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning EntityType:EntityTypeUpdateMetadata - unable to get AttrDesc w/ Error: '%+v'", err.Error()))
			return err
		}
		// TODO: Revisit later - Something not correct - should we continue or throw an error
		if attrDesc == nil {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop EntityType:EntityTypeUpdateMetadata - cannot find '%s' attribute descriptor", attrName))
			continue
		}
		if attrDesc.GetAttrType() == AttributeTypeInvalid {
			logger.Warning(fmt.Sprint("WARNING: Continuing loop EntityType:EntityTypeUpdateMetadata as AttrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			continue
		}
		obj.AddAttributeDescriptor(attrName, attrDesc)
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning EntityType:EntityTypeUpdateMetadata w/ NO error, for entityType: '%+v'", obj))
	}
	return nil
}

func (obj *EntityType) entityTypeToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("EntityType:{")
	buffer.WriteString(fmt.Sprintf("SysType: %d", obj.SysType))
	buffer.WriteString(fmt.Sprintf(", Id: %+v", obj.id))
	buffer.WriteString(fmt.Sprintf(", Name: %+v", obj.Name))
	//buffer.WriteString(fmt.Sprintf(", Parent: %+v", obj.parent))
	//buffer.WriteString(fmt.Sprintf(", EntityTypeAttributes: %+v", obj.Attributes))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *EntityType) AddAttributeDescriptor(attrName string, attrDesc tgdb.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.Attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *EntityType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *EntityType) DerivedFrom() tgdb.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified Name
func (obj *EntityType) GetAttributeDescriptor(attrName string) tgdb.TGAttributeDescriptor {
	attrDesc := obj.Attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *EntityType) GetAttributeDescriptors() []tgdb.TGAttributeDescriptor {
	attrDescriptors := make([]tgdb.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.Attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *EntityType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's Name
func (obj *EntityType) SetName(eTypeName string) {
	obj.Name = eTypeName
}

// SetSystemType sets system object's type
func (obj *EntityType) SetSystemType(eSysType tgdb.TGSystemType) {
	obj.SysType = eSysType
}

func (obj *EntityType) String() string {
	return obj.entityTypeToString()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the Name for this entity type as the most generic form
func (obj *EntityType) GetName() string {
	return obj.Name
}

func (obj *EntityType) GetSystemType() tgdb.TGSystemType {
	return obj.SysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *EntityType) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return EntityTypeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *EntityType) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning EntityType:WriteExternal is not implemented"))
	//errMsg := fmt.Sprint("EntityType WriteExternal message is not implemented")
	//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *EntityType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.SysType, obj.id, obj.Name, obj.parent, obj.Attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *EntityType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.SysType, &obj.id, &obj.Name, &obj.parent, &obj.Attributes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EntityType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

type NodeType struct {
	*EntityType
	pKeys      []*AttributeDescriptor
	idxIds     []int
	numEntries int64
}

func DefaultNodeType() *NodeType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NodeType{})

	newNodeType := NodeType{
		EntityType: DefaultEntityType(),
		pKeys:      make([]*AttributeDescriptor, 0),
		idxIds:     make([]int, 0),
		numEntries: 0,
	}
	newNodeType.SysType = tgdb.SystemTypeNode
	return &newNodeType
}

func NewNodeType(name string, parent tgdb.TGEntityType) *NodeType {
	newNodeType := DefaultNodeType()
	newNodeType.Name = name
	newNodeType.parent = parent
	return newNodeType
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGNodeType
/////////////////////////////////////////////////////////////////

func (obj *NodeType) GetIndexIds() []int {
	return obj.idxIds
}

func (obj *NodeType) GetNumEntries() int64 {
	return obj.numEntries
}

func (obj *NodeType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.Attributes = attrMap
}

func (obj *NodeType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.Attributes[attrName] = attrDesc
}

func (obj *NodeType) SetParent(parentEntity tgdb.TGEntityType) {
	obj.parent = parentEntity
}

func (obj *NodeType) SetNumEntries(num int64) {
	obj.numEntries = num
}

func (obj *NodeType) UpdateMetadata(gmd *GraphMetadata) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering NodeType:UpdateMetadata"))
	}
	// Base Class EntityType::UpdateMetadata()

	//err := EntityTypeUpdateMetadata(obj, gmd)
	err := EntityTypeUpdateMetadata(obj.EntityType, gmd)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside NodeType:UpdateMetadata, updated base entity type's Attributes"))
	}
	for id, key := range obj.pKeys {
		attrDesc, err := gmd.GetAttributeDescriptor(key.GetName())
		if err == nil {
			//logger.Warning(fmt.Sprintf("WARNING: Continuing loop NodeType:UpdateMetadata - cannot find '%s' attribute descriptor", key.GetName()))
			continue
		}
		if attrDesc.GetAttrType() == AttributeTypeInvalid {
			//logger.Warning(fmt.Sprint("WARNING: Continuing loop NodeType:UpdateMetadata as AttrDesc.GetAttrType() == types.AttributeTypeInvalid"))
			continue
		}
		obj.pKeys[id] = attrDesc.(*AttributeDescriptor)
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning NodeType:UpdateMetadata w/ NO error, for entityType: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNodeType
/////////////////////////////////////////////////////////////////

// GetPKeyAttributeDescriptors returns a set of primary key descriptors
func (obj *NodeType) GetPKeyAttributeDescriptors() []tgdb.TGAttributeDescriptor {
	pkDesc := make([]tgdb.TGAttributeDescriptor, 0)
	for _, pk := range obj.pKeys {
		pkDesc = append(pkDesc, pk)
	}
	return pkDesc
}

// SetPKeyAttributeDescriptors sets primary key descriptors
func (obj *NodeType) SetPKeyAttributeDescriptors(keys []*AttributeDescriptor) {
	obj.pKeys = keys
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *NodeType) AddAttributeDescriptor(attrName string, attrDesc tgdb.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.Attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *NodeType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *NodeType) DerivedFrom() tgdb.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified Name
func (obj *NodeType) GetAttributeDescriptor(attrName string) tgdb.TGAttributeDescriptor {
	attrDesc := obj.Attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *NodeType) GetAttributeDescriptors() []tgdb.TGAttributeDescriptor {
	attrDescriptors := make([]tgdb.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.Attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *NodeType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's Name
func (obj *NodeType) SetName(eTypeName string) {
	obj.Name = eTypeName
}

// SetSystemType sets system object's type
func (obj *NodeType) SetSystemType(eSysType tgdb.TGSystemType) {
	obj.SysType = eSysType
}

func (obj *NodeType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NodeType:{")
	//buffer.WriteString(fmt.Sprintf("PKeys: %+v", obj.PKeys))
	buffer.WriteString(fmt.Sprintf(", IdxIds: %+v", obj.idxIds))
	buffer.WriteString(fmt.Sprintf(", NumEntries: %+v", obj.numEntries))
	strArray := []string{buffer.String(), obj.entityTypeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's Name
func (obj *NodeType) GetName() string {
	return obj.Name
}

// GetSystemType gets system object's type
func (obj *NodeType) GetSystemType() tgdb.TGSystemType {
	return obj.SysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *NodeType) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering NodeType:ReadExternal"))
	}
	// Base Class EntityType's ReadExternal()
	err := EntityTypeReadExternal(obj, is)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside NodeType:ReadExternal, read base entity type's Attributes"))
	}

	attrCount, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read attrCount w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside NodeType:ReadExternal read attrCount as '%d'", attrCount))
	}
	for i := 0; i < int(attrCount); i++ {
		attrName, err := is.(*ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read attrName w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside NodeType:ReadExternal read attrName as '%s'", attrName))
		}
		attrDesc := NewAttributeDescriptorWithType(attrName, AttributeTypeString)
		obj.pKeys = append(obj.pKeys, attrDesc)
	}

	idxCount, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read idxCount w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside NodeType:ReadExternal read idxCount as '%d'", idxCount))
	}
	for i := 0; i < int(idxCount); i++ {
		// TODO: Revisit later to get meta data needs to return index definitions
		indexId, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read indexId w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside NodeType:ReadExternal read indexId as '%d'", indexId))
		}
		obj.idxIds = append(obj.idxIds, indexId)
	}

	numEntries, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:ReadExternal - unable to read numEntries w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside NodeType:ReadExternal read numEntries as '%d'", numEntries))
	}

	obj.SetNumEntries(numEntries)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning NodeType:ReadExternal w/ NO error, for NodeType: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *NodeType) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning NodeType:WriteExternal is not implemented"))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NodeType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.SysType, obj.id, obj.Name, obj.parent, obj.Attributes, obj.pKeys, obj.idxIds, obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NodeType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.SysType, &obj.id, &obj.Name, &obj.parent, &obj.Attributes, &obj.pKeys,
		&obj.idxIds, &obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NodeType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type Node struct {
	*AbstractEntity
	Edges []tgdb.TGEdge
}

func DefaultNode() *Node {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Node{})

	newNode := Node{
		AbstractEntity: DefaultAbstractEntity(),
	}
	newNode.EntityKind = tgdb.EntityKindNode
	newNode.EntityType = DefaultNodeType()
	newNode.Edges = make([]tgdb.TGEdge, 0)
	return &newNode
}

func NewNode(gmd *GraphMetadata) *Node {
	newNode := DefaultNode()
	newNode.graphMetadata = gmd
	return newNode
}

func NewNodeWithType(gmd *GraphMetadata, nodeType tgdb.TGNodeType) *Node {
	newNode := NewNode(gmd)
	newNode.EntityType = nodeType
	return newNode
}

/////////////////////////////////////////////////////////////////
// Helper functions for Node
/////////////////////////////////////////////////////////////////

func (obj *Node) GetIsInitialized() bool {
	return obj.isInitialized
}

func (obj *Node) GetModifiedAttributes() []tgdb.TGAttribute {
	return obj.getModifiedAttributes()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNode
/////////////////////////////////////////////////////////////////

func (obj *Node) AddEdge(edge tgdb.TGEdge) {
	obj.Edges = append(obj.Edges, edge)
}

func (obj *Node) AddEdgeWithDirectionType(node tgdb.TGNode, edgeType tgdb.TGEdgeType, directionType tgdb.TGDirectionType) tgdb.TGEdge {
	newEdge := NewEdgeWithDirection(obj.graphMetadata, obj, node, directionType)
	obj.AddEdge(newEdge)
	return newEdge
}

func (obj *Node) GetEdges() []tgdb.TGEdge {
	return obj.Edges
}

func (obj *Node) GetEdgesForDirectionType(directionType tgdb.TGDirectionType) []tgdb.TGEdge {
	edgesWithDirections := make([]tgdb.TGEdge, 0)
	if len(obj.Edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Node:GetEdgesForDirectionType as there are NO edges"))
		return edgesWithDirections
	}

	for _, edge := range obj.Edges {
		if edge.(*Edge).directionType == directionType {
			edgesWithDirections = append(edgesWithDirections, edge)
		}
	}
	return edgesWithDirections
}

func (obj *Node) GetEdgesForEdgeType(edgeType tgdb.TGEdgeType, direction tgdb.TGDirection) []tgdb.TGEdge {
	edgesWithDirections := make([]tgdb.TGEdge, 0)
	if len(obj.Edges) == 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning Node:GetEdgesForEdgeType as there are NO edges"))
		return edgesWithDirections
	}

	if edgeType == nil && direction == tgdb.DirectionAny {
		for _, edge := range obj.Edges {
			if edge.(*Edge).GetIsInitialized() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		}
		return obj.Edges
	}

	for _, edge := range obj.Edges {
		if !edge.(*Edge).GetIsInitialized() {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop Node:GetEdgesForEdgeType - skipping uninitialized edge '%+v'", edge))
			continue
		}
		eType := edge.GetEntityType()
		if edgeType != nil && eType != nil && eType.GetName() != edgeType.GetName() {
			logger.Warning(fmt.Sprintf("WARNING: Continuing loop Node:GetEdgesForEdgeType - skipping (entity type NOT matching) edge '%+v'", edge))
			continue
		}
		if direction == tgdb.DirectionAny {
			edgesWithDirections = append(edgesWithDirections, edge)
		} else if direction == tgdb.DirectionOutbound {
			edgesForThisNode := edge.GetVertices()
			if obj.GetVirtualId() == edgesForThisNode[0].GetVirtualId() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		} else {
			edgesForThisNode := edge.GetVertices()
			if obj.GetVirtualId() == edgesForThisNode[1].GetVirtualId() {
				edgesWithDirections = append(edgesWithDirections, edge)
			}
		}
	}
	return edgesWithDirections
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the Name specified
func (obj *Node) GetAttribute(attrName string) tgdb.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the Attributes set
func (obj *Node) GetAttributes() ([]tgdb.TGAttribute, tgdb.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *Node) GetEntityKind() tgdb.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *Node) GetEntityType() tgdb.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *Node) GetGraphMetadata() tgdb.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *Node) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *Node) GetIsNew() bool {
	return obj.getIsNew()
}

// GetVersion gets the version of the Entity
func (obj *Node) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *Node) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *Node) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on Attributes
func (obj *Node) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *Node) SetAttribute(attr tgdb.TGAttribute) tgdb.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *Node) SetOrCreateAttribute(name string, value interface{}) tgdb.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *Node) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *Node) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsInitialized set the initialized flag
func (obj *Node) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *Node) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *Node) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *Node) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Node:{")
	buffer.WriteString(fmt.Sprintf("Edges: %+v", obj.Edges))
	strArray := []string{buffer.String(), obj.entityToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *Node) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering Node:ReadExternal"))
	}
	// TODO: Revisit later - Do we need to validate length?
	nodeBufLen, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read length w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Node:ReadExternal read nodeBufLen as '%+v'", nodeBufLen))
	}

	err = obj.AbstractEntityReadExternal(is)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to obj.AbstractEntityReadExternal(is) w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside Node:ReadExternal read abstractEntity"))
	}

	edgeCount, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read edgeCount w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Node:ReadExternal read edgeCount as '%d'", edgeCount))
	}
	for i := 0; i < edgeCount; i++ {
		edgeId, err := is.(*ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning Node:ReadExternal - unable to read entId w/ Error: '%+v'", err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside Node:ReadExternal read edgeId as '%d'", edgeId))
		}
		var edge *Edge
		var entity tgdb.TGEntity
		refMap := is.(*ProtocolDataInputStream).GetReferenceMap()
		if refMap != nil {
			entity = refMap[edgeId]
		}
		if entity == nil {
			edge1 := NewEdge(obj.graphMetadata)
			edge1.SetEntityId(edgeId)
			edge1.SetIsInitialized(false)
			if refMap != nil {
				refMap[edgeId] = edge1
			}
			edge = edge1
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("Inside Node:ReadExternal created new edge: '%+v'", edge))
			}
		} else {
			edge = entity.(*Edge)
		}
		obj.Edges = append(obj.Edges, edge)
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside Node:ReadExternal Node has '%d' edges & StreamEntityCount is '%d'", len(obj.Edges), len(is.(*ProtocolDataInputStream).GetReferenceMap())))
		}
	}

	obj.SetIsInitialized(true)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning Node:ReadExternal w/ NO error, for node: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *Node) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering Node:WriteExternal"))
	}
	startPos := os.(*ProtocolDataOutputStream).GetPosition()
	os.(*ProtocolDataOutputStream).WriteInt(0)
	// Write Attributes from the base class
	err := obj.AbstractEntityWriteExternal(os)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside Node:WriteExternal - exported base entity Attributes"))
	}
	newCount := 0
	for _, edge := range obj.Edges {
		if edge.GetIsNew() {
			newCount++
		}
	}
	os.(*ProtocolDataOutputStream).WriteInt(newCount)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Node:WriteExternal - exported new edge count '%d'", newCount))
	}
	// Write the edges ids - ONLY include new edges
	for _, edge := range obj.Edges {
		if ! edge.GetIsNew() {
			continue
		}
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetVirtualId())
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside Node:WriteExternal - exported a new edge: '%+v'", edge))
		}
	}
	currPos := os.(*ProtocolDataOutputStream).GetPosition()
	length := currPos - startPos
	_, err = os.(*ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:WriteExternal - unable to update data length in the buffer w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning Node:WriteExternal w/ NO error, for node: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *Node) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.VirtualId, obj.Version, obj.EntityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.Attributes, obj.ModifiedAttributes, obj.Edges)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *Node) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.VirtualId, &obj.Version, &obj.EntityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.Attributes, &obj.ModifiedAttributes, &obj.Edges)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Node:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

type Edge struct {
	*AbstractEntity
	directionType tgdb.TGDirectionType
	FromNode      tgdb.TGNode
	ToNode        tgdb.TGNode
}

func DefaultEdge() *Edge {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Edge{})

	newEdge := Edge{
		AbstractEntity: DefaultAbstractEntity(),
	}
	newEdge.EntityKind = tgdb.EntityKindEdge
	newEdge.EntityType = DefaultEdgeType()
	return &newEdge
}

func NewEdge(gmd *GraphMetadata) *Edge {
	newEdge := DefaultEdge()
	newEdge.graphMetadata = gmd
	return newEdge
}

func NewEdgeWithDirection(gmd *GraphMetadata, fromNode tgdb.TGNode, toNode tgdb.TGNode, directionType tgdb.TGDirectionType) *Edge {
	newEdge := NewEdge(gmd)
	newEdge.directionType = directionType
	newEdge.FromNode = fromNode
	newEdge.ToNode = toNode
	return newEdge
}

func NewEdgeWithEdgeType(gmd *GraphMetadata, fromNode tgdb.TGNode, toNode tgdb.TGNode, edgeType tgdb.TGEdgeType) *Edge {
	newEdge := NewEdge(gmd)
	newEdge.FromNode = fromNode
	newEdge.ToNode = toNode
	newEdge.EntityType = edgeType
	newEdge.directionType = edgeType.GetDirectionType()
	return newEdge
}

/////////////////////////////////////////////////////////////////
// Helper functions for Edge
/////////////////////////////////////////////////////////////////

func (obj *Edge) GetFromNode() tgdb.TGNode {
	return obj.FromNode
}

func (obj *Edge) GetIsInitialized() bool {
	return obj.isInitialized
}

func (obj *Edge) GetModifiedAttributes() []tgdb.TGAttribute {
	return obj.getModifiedAttributes()
}

func (obj *Edge) GetToNode() tgdb.TGNode {
	return obj.ToNode
}

func (obj *Edge) SetDirectionType(dirType tgdb.TGDirectionType) {
	obj.directionType = dirType
}

func (obj *Edge) SetFromNode(node tgdb.TGNode) {
	obj.FromNode = node
}

func (obj *Edge) SetToNode(node tgdb.TGNode) {
	obj.ToNode = node
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEdge
/////////////////////////////////////////////////////////////////

// GetDirectionType gets direction type as one of the constants
func (obj *Edge) GetDirectionType() tgdb.TGDirectionType {
	if obj.EntityType != nil {
		return obj.EntityType.(*EdgeType).GetDirectionType()
	} else {
		return obj.directionType
	}
}

// GetVertices gets array of NODE (Entity) types for this EDGE (Entity) type
func (obj *Edge) GetVertices() []tgdb.TGNode {
	return []tgdb.TGNode{obj.FromNode, obj.ToNode}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntity
/////////////////////////////////////////////////////////////////

// GetAttribute gets the attribute for the Name specified
func (obj *Edge) GetAttribute(attrName string) tgdb.TGAttribute {
	return obj.getAttribute(attrName)
}

// GetAttributes lists of all the Attributes set
func (obj *Edge) GetAttributes() ([]tgdb.TGAttribute, tgdb.TGError) {
	return obj.getAttributes()
}

// GetEntityKind returns the EntityKind as a constant
func (obj *Edge) GetEntityKind() tgdb.TGEntityKind {
	return obj.getEntityKind()
}

// GetEntityType returns the EntityType
func (obj *Edge) GetEntityType() tgdb.TGEntityType {
	return obj.getEntityType()
}

// GetGraphMetadata returns the Graph Meta Data	- New in GO Lang
func (obj *Edge) GetGraphMetadata() tgdb.TGGraphMetadata {
	return obj.getGraphMetadata()
}

// GetIsDeleted checks whether this entity is already deleted in the system or not
func (obj *Edge) GetIsDeleted() bool {
	return obj.getIsDeleted()
}

// GetIsNew checks whether this entity that is currently being added to the system is new or not
func (obj *Edge) GetIsNew() bool {
	return obj.getIsNew()
}

// GetVersion gets the version of the Entity
func (obj *Edge) GetVersion() int {
	return obj.getVersion()
}

// GetVirtualId gets Entity identifier
// At the time of creation before reaching the server, it is the virtual id
// Upon successful creation, server returns a valid entity id that gets set in place of virtual id
func (obj *Edge) GetVirtualId() int64 {
	return obj.getVirtualId()
}

// IsAttributeSet checks whether this entity is an Attribute set or not
func (obj *Edge) IsAttributeSet(attrName string) bool {
	return obj.isAttributeSet(attrName)
}

// ResetModifiedAttributes resets the dirty flag on Attributes
func (obj *Edge) ResetModifiedAttributes() {
	obj.resetModifiedAttributes()
}

// SetAttribute associates the specified Attribute to this Entity
func (obj *Edge) SetAttribute(attr tgdb.TGAttribute) tgdb.TGError {
	return obj.setAttribute(attr)
}

// SetOrCreateAttribute dynamically associates the attribute to this entity
// If the AttributeDescriptor doesn't exist in the database, create a new one
func (obj *Edge) SetOrCreateAttribute(name string, value interface{}) tgdb.TGError {
	return obj.setOrCreateAttribute(name, value)
}

// SetEntityId sets Entity id and reset Virtual id after creation
func (obj *Edge) SetEntityId(id int64) {
	obj.setEntityId(id)
}

// SetIsDeleted set the deleted flag
func (obj *Edge) SetIsDeleted(flag bool) {
	obj.setIsDeleted(flag)
}

// SetIsInitialized set the initialized flag
func (obj *Edge) SetIsInitialized(flag bool) {
	obj.setIsInitialized(flag)
}

// SetIsNew sets the flag that this is a new entity
func (obj *Edge) SetIsNew(flag bool) {
	obj.setIsNew(flag)
}

// SetVersion sets the version of the Entity
func (obj *Edge) SetVersion(version int) {
	obj.setVersion(version)
}

func (obj *Edge) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Edge:{")
	buffer.WriteString(fmt.Sprintf("DirectionType: %+v", obj.directionType))
	if obj.FromNode != nil {
		buffer.WriteString(fmt.Sprintf(", FromNode: %+v", obj.FromNode.GetVirtualId()))
	}
	if obj.ToNode != nil {
		buffer.WriteString(fmt.Sprintf(", ToNode: %+v", obj.ToNode.GetVirtualId()))
	}
	strArray := []string{buffer.String(), obj.entityToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *Edge) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering Edge:ReadExternal"))
	}
	// TODO: Revisit later - Do we need to validate length?
	edgeBufLen, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal read edgeBufLen as '%+v'", edgeBufLen))
	}

	err = obj.AbstractEntityReadExternal(is)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside Edge:ReadExternal read abstractEntity"))
	}

	direction, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read direction w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal read direction as '%+v'", direction))
	}
	if direction == 0 {
		obj.SetDirectionType(tgdb.DirectionTypeUnDirected)
	} else if direction == 1 {
		obj.SetDirectionType(tgdb.DirectionTypeDirected)
	} else {
		obj.SetDirectionType(tgdb.DirectionTypeBiDirectional)
	}

	var fromEntity, toEntity tgdb.TGEntity
	var fromNode, toNode *Node

	fromNodeId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read fromNodeId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal read fromNodeId as '%d'", fromNodeId))
	}
	refMap := is.(*ProtocolDataInputStream).GetReferenceMap()
	if refMap != nil {
		fromEntity = refMap[fromNodeId]
	}
	if fromEntity == nil {
		fNode := NewNode(obj.graphMetadata)
		fNode.SetEntityId(fromNodeId)
		fNode.SetIsInitialized(false)
		if refMap != nil {
			refMap[fromNodeId] = fNode
		}
		fromNode = fNode
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal created new fromNode: '%+v'", fromNode))
		}
	} else {
		fromNode = fromEntity.(*Node)
	}
	obj.SetFromNode(fromNode)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal Edge has fromNode & StreamEntityCount is '%d'", len(is.(*ProtocolDataInputStream).GetReferenceMap())))
	}

	toNodeId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:ReadExternal - unable to read toNodeId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal read toNodeId as '%d'", toNodeId))
	}
	if refMap != nil {
		toEntity = refMap[toNodeId]
	}
	if toEntity == nil {
		tNode := NewNode(obj.graphMetadata)
		tNode.SetEntityId(toNodeId)
		tNode.SetIsInitialized(false)
		if refMap != nil {
			refMap[toNodeId] = tNode
		}
		toNode = tNode
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal created new toNode: '%+v'", toNode))
		}
	} else {
		toNode = toEntity.(*Node)
	}
	obj.SetToNode(toNode)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Edge:ReadExternal Edge has toNode & StreamEntityCount is '%d'", len(is.(*ProtocolDataInputStream).GetReferenceMap())))
	}

	obj.SetIsInitialized(true)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning Edge:ReadExternal w/ NO error, for edge: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *Edge) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering Edge:WriteExternal"))
	}
	startPos := os.(*ProtocolDataOutputStream).GetPosition()
	os.(*ProtocolDataOutputStream).WriteInt(0)
	// Write Attributes from the base class
	err := obj.AbstractEntityWriteExternal(os)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside Edge:WriteExternal - exported base entity Attributes"))
	}
	// TODO: Revisit later - Check w/ TGDB Engineering team as what the difference should be for if-n-else conditions
	// Write the edges ids
	if obj.GetIsNew() {
		os.(*ProtocolDataOutputStream).WriteByte(int(obj.GetDirectionType()))
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetFromNode().GetVirtualId())
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetToNode().GetVirtualId())
	} else {
		os.(*ProtocolDataOutputStream).WriteByte(int(obj.GetDirectionType())) // The Server expects it - so better send it.
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetFromNode().GetVirtualId())
		os.(*ProtocolDataOutputStream).WriteLong(obj.GetToNode().GetVirtualId())
	}
	currPos := os.(*ProtocolDataOutputStream).GetPosition()
	length := currPos - startPos
	_, err = os.(*ProtocolDataOutputStream).WriteIntAt(startPos, length)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:WriteExternal - unable to update data length in the buffer w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning Edge:WriteExternal w/ NO error, for edge: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *Edge) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.isNew, obj.EntityKind, obj.VirtualId, obj.Version, obj.EntityId, obj.EntityType,
		obj.isDeleted, obj.isInitialized, obj.graphMetadata, obj.Attributes, obj.ModifiedAttributes,
		obj.directionType, obj.FromNode, obj.ToNode)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *Edge) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.isNew, &obj.EntityKind, &obj.VirtualId, &obj.Version, &obj.EntityId, &obj.EntityType,
		&obj.isDeleted, &obj.isInitialized, &obj.graphMetadata, &obj.Attributes, &obj.ModifiedAttributes,
		&obj.directionType, &obj.FromNode, &obj.ToNode)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning Edge:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type EdgeType struct {
	*EntityType
	directionType tgdb.TGDirectionType
	fromTypeId    int
	fromNodeType  tgdb.TGNodeType
	toTypeId      int
	toNodeType    tgdb.TGNodeType
	numEntries    int64
}

func DefaultEdgeType() *EdgeType {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(EdgeType{})

	newEdgeType := &EdgeType{
		EntityType: DefaultEntityType(),
		numEntries: 0,
	}
	newEdgeType.SysType = tgdb.SystemTypeEdge
	return newEdgeType
}

func NewEdgeType(name string, directionType tgdb.TGDirectionType, parent tgdb.TGEntityType) *EdgeType {
	newEdgeType := DefaultEdgeType()
	newEdgeType.Name = name
	newEdgeType.directionType = directionType
	newEdgeType.parent = parent
	return newEdgeType
}

/////////////////////////////////////////////////////////////////
// Helper functions for EdgeType
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) SetAttributeMap(attrMap map[string]*AttributeDescriptor) {
	obj.Attributes = attrMap
}

func (obj *EdgeType) SetAttributeDesc(attrName string, attrDesc *AttributeDescriptor) {
	obj.Attributes[attrName] = attrDesc
}

func (obj *EdgeType) SetParent(parentEntity tgdb.TGEntityType) {
	obj.parent = parentEntity
}

func (obj *EdgeType) SetDirectionType(dirType tgdb.TGDirectionType) {
	obj.directionType = dirType
}

func (obj *EdgeType) SetNumEntries(num int64) {
	obj.numEntries = num
}

func (obj *EdgeType) GetNumEntries() int64 {
	return obj.numEntries
}

func (obj *EdgeType) UpdateMetadata(gmd *GraphMetadata) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering EdgeType:UpdateMetadata"))
	}
	// Base Class EntityType::UpdateMetadata()
	//err := EntityTypeUpdateMetadata(obj, gmd)
	err := EntityTypeUpdateMetadata(obj.EntityType, gmd)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside EdgeType:UpdateMetadata, updated base entity type's Attributes"))
	}
	nType, nErr := gmd.GetNodeTypeById(obj.fromTypeId)
	if nErr == nil {
		if nType != nil {
			obj.fromNodeType = nType.(*NodeType)
		}
	}
	nType, nErr = gmd.GetNodeTypeById(obj.toTypeId)
	if nErr == nil {
		if nType != nil {
			obj.toNodeType = nType.(*NodeType)
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning EdgeType:UpdateMetadata w/ NO error, for EdgeType: '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNodeType
/////////////////////////////////////////////////////////////////

// GetDirectionType gets direction type as one of the constants
func (obj *EdgeType) GetDirectionType() tgdb.TGDirectionType {
	return obj.directionType
}

// GetFromNodeType gets From-Node Type
func (obj *EdgeType) GetFromNodeType() tgdb.TGNodeType {
	return obj.fromNodeType
}

// GetFromTypeId gets From-Node ID
func (obj *EdgeType) GetFromTypeId() int {
	return obj.fromTypeId
}

// GetToNodeType gets To-Node Type
func (obj *EdgeType) GetToNodeType() tgdb.TGNodeType {
	return obj.toNodeType
}

// GetToTypeId gets To-Node ID
func (obj *EdgeType) GetToTypeId() int {
	return obj.toTypeId
}

// SetFromNodeType sets From-Node Type
func (obj *EdgeType) SetFromNodeType(fromNode tgdb.TGNodeType) {
	obj.fromNodeType = fromNode
}

// SetFromTypeId sets From-Node ID
func (obj *EdgeType) SetFromTypeId(fromTypeId int) {
	obj.fromTypeId = fromTypeId
}

// SetToNodeType sets From-Node Type
func (obj *EdgeType) SetToNodeType(toNode tgdb.TGNodeType) {
	obj.toNodeType = toNode
}

// SetToTypeId sets To-Node ID
func (obj *EdgeType) SetToTypeId(toTypeId int) {
	obj.toTypeId = toTypeId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGEntityType
/////////////////////////////////////////////////////////////////

// AddAttributeDescriptor add an attribute descriptor to the map
func (obj *EdgeType) AddAttributeDescriptor(attrName string, attrDesc tgdb.TGAttributeDescriptor) {
	if attrName != "" && attrDesc != nil {
		obj.Attributes[attrName] = attrDesc.(*AttributeDescriptor)
	}
}

// GetEntityTypeId gets Entity Type id
func (obj *EdgeType) GetEntityTypeId() int {
	return obj.id
}

// DerivedFrom gets the parent Entity Type
func (obj *EdgeType) DerivedFrom() tgdb.TGEntityType {
	return obj.parent
}

// GetAttributeDescriptor gets the attribute descriptor for the specified Name
func (obj *EdgeType) GetAttributeDescriptor(attrName string) tgdb.TGAttributeDescriptor {
	attrDesc := obj.Attributes[attrName]
	return attrDesc
}

// GetAttributeDescriptors returns a collection of attribute descriptors associated with this Entity Type
func (obj *EdgeType) GetAttributeDescriptors() []tgdb.TGAttributeDescriptor {
	attrDescriptors := make([]tgdb.TGAttributeDescriptor, 0)
	for _, attrDesc := range obj.Attributes {
		attrDescriptors = append(attrDescriptors, attrDesc)
	}
	return attrDescriptors
}

// SetEntityTypeId sets Entity Type id
func (obj *EdgeType) SetEntityTypeId(eTypeId int) {
	obj.id = eTypeId
}

// SetName sets the system object's Name
func (obj *EdgeType) SetName(eTypeName string) {
	obj.Name = eTypeName
}

// SetSystemType sets system object's type
func (obj *EdgeType) SetSystemType(eSysType tgdb.TGSystemType) {
	obj.SysType = eSysType
}

func (obj *EdgeType) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("EdgeType:{")
	buffer.WriteString(fmt.Sprintf("DirectionType: %+v", obj.directionType))
	buffer.WriteString(fmt.Sprintf(", FromTypeId: %+v", obj.fromTypeId))
	buffer.WriteString(fmt.Sprintf(", FromNodeType: %+v", obj.fromNodeType))
	buffer.WriteString(fmt.Sprintf(", ToTypeId: %+v", obj.toTypeId))
	buffer.WriteString(fmt.Sprintf(", ToNodeType: %+v", obj.toNodeType))
	buffer.WriteString(fmt.Sprintf(", NumEntries: %+v", obj.numEntries))
	strArray := []string{buffer.String(), obj.entityTypeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's Name
func (obj *EdgeType) GetName() string {
	return obj.Name
}

// GetSystemType gets system object's type
func (obj *EdgeType) GetSystemType() tgdb.TGSystemType {
	return obj.SysType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *EdgeType) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering EdgeType:ReadExternal"))
	}
	// Base Class EntityType's ReadExternal()
	err := EntityTypeReadExternal(obj, is)
	if err != nil {
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside EdgeType:ReadExternal, read base entity type's Attributes"))
	}

	fromTypeId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read fromTypeId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EdgeType:ReadExternal read fromTypeId as '%d'", fromTypeId))
	}

	toTypeId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read toTypeId w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EdgeType:ReadExternal read toTypeId as '%d'", toTypeId))
	}

	direction, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read direction w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EdgeType:ReadExternal read direction as '%+v'", direction))
	}
	if direction == 0 {
		obj.SetDirectionType(tgdb.DirectionTypeUnDirected)
	} else if direction == 1 {
		obj.SetDirectionType(tgdb.DirectionTypeDirected)
	} else {
		obj.SetDirectionType(tgdb.DirectionTypeBiDirectional)
	}

	numEntries, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:ReadExternal - unable to read numEntries w/ Error: '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside EdgeType:ReadExternal read numEntries as '%d'", numEntries))
	}

	obj.SetFromTypeId(fromTypeId)
	obj.SetToTypeId(toTypeId)
	obj.SetNumEntries(numEntries)
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning EdgeType:ReadExternal w/ NO error, for entityType: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *EdgeType) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	logger.Warning(fmt.Sprint("WARNING: Returning EdgeType:WriteExternal is not implemented"))
	//errMsg := fmt.Sprint("EdgeType WriteExternal message is not implemented")
	//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.SysType, obj.id, obj.Name, obj.parent, obj.Attributes, obj.directionType,
		obj.fromTypeId, obj.fromNodeType, obj.toTypeId, obj.toNodeType, obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *EdgeType) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.SysType, &obj.id, &obj.Name, &obj.parent, &obj.Attributes, &obj.directionType,
		&obj.fromTypeId, &obj.fromNodeType, &obj.toTypeId, &obj.toNodeType, &obj.numEntries)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning EdgeType:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

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

func (obj *ByteArrayEntity) equals(cmpObj tgdb.TGEntityId) bool {
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
func (obj *ByteArrayEntity) ToBytes() ([]byte, tgdb.TGError) {
	var encodedMsg bytes.Buffer
	encoder := gob.NewEncoder(&encodedMsg)
	err := encoder.Encode(obj)
	//err := interfaceEncode(encoder, obj)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteArrayEntity:ToBytes - unable to write entity in byte format w/ Error: '%+v'", err.Error()))
		errMsg := fmt.Sprintf("Unable to convert message '%+v' in byte format", obj)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Log(fmt.Sprintf("Message: %v bytes <--encoded-- %v", encodedMsg.Len(), m)
	return encodedMsg.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ByteArrayEntity) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	entityId, err := is.(*ProtocolDataInputStream).ReadBytes()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteArrayEntity:ReadExternal - unable to read entityId w/ Error: '%+v'", err.Error()))
		return err
	}
	obj.entityId = entityId
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning ByteArrayEntity:ReadExternal w/ NO error, for entity: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ByteArrayEntity) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Exported ByteArrayEntity object as '%+v' from byte format", obj))
	return os.(*ProtocolDataOutputStream).WriteBytes(obj.entityId)
}


type CompositeKey struct {
	graphMetadata *GraphMetadata
	keyName       string
	attributes    map[string]tgdb.TGAttribute
}

func NewCompositeKey(graphMetadata *GraphMetadata, typeName string) *CompositeKey {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CompositeKey{})

	newCompositeKey := CompositeKey{
		graphMetadata: graphMetadata,
		keyName:       typeName,
		attributes:    make(map[string]tgdb.TGAttribute, 0),
	}
	return &newCompositeKey
}

/////////////////////////////////////////////////////////////////
// Helper functions for CompositeKey
/////////////////////////////////////////////////////////////////

func (obj *CompositeKey) GetAttributes() map[string]tgdb.TGAttribute {
	return obj.attributes
}

func (obj *CompositeKey) GetKeyName() string {
	return obj.keyName
}

func (obj *CompositeKey) SetAttributes(attrs map[string]tgdb.TGAttribute) {
	obj.attributes = attrs
}

func (obj *CompositeKey) SetKeyName(name string) {
	obj.keyName = name
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGKey
/////////////////////////////////////////////////////////////////

// Dynamically set the attribute to this entity. If the AttributeDescriptor doesn't exist in the database, create a new one.
func (obj *CompositeKey) SetOrCreateAttribute(name string, value interface{}) tgdb.TGError {
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Entering CompositeKey:SetOrCreateAttribute w/ N-V Pair as '%+v'='%+v'", name, value))
	}
	if name == "" || value == nil {
		logger.Error(fmt.Sprint("ERROR: Returning CompositeKey:SetOrCreateAttribute as either Name or value is EMPTY"))
		errMsg := "Name or Value is null"
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	// If attribute is not present in the set, create a new one
	attr := obj.attributes[name]
	if attr == nil {
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CompositeKey:SetOrCreateAttribute attribute '%+v' not found - trying to get descriptor from GraphMetadata", name))
		}
		attrDesc, err := obj.graphMetadata.GetAttributeDescriptor(name)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to get descriptor for attribute '%s' w/ error '%+v'", name, err.Error()))
			return err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside CompositeKey:SetOrCreateAttribute attribute descriptor for '%+v' found in GraphMetadata", name))
		}
		if attrDesc == nil {
			aType := reflect.TypeOf(value).String()
			if logger.IsDebug() {
							logger.Debug(fmt.Sprintf("=======> Inside CompositeKey SetOrCreateAttribute creating new attribute descriptor '%+v':'%+v'(%+v) <=======", name, value, aType))
			}
			attrDesc = obj.graphMetadata.CreateAttributeDescriptorForDataType(name, aType)
		}
		newAttr, aErr := CreateAttributeWithDesc(nil, attrDesc.(*AttributeDescriptor), value)
		if aErr != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to create attribute '%s' w/ descriptor and value '%+v'", name, value))
			return aErr
		}
		//newAttr.SetOwner(obj)
		attr = newAttr
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside CompositeKey:SetOrCreateAttribute trying to set attribute '%+v' value as '%+v'", attr, value))
	}
	// Set the attribute value
	err := attr.SetValue(value)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:SetOrCreateAttribute unable to set attribute value w/ error '%+v'", err.Error()))
		return err
	}
	// Add it to the set
	obj.attributes[name] = attr
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CompositeKey:SetOrCreateAttribute w/ Key as '%+v'", obj))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *CompositeKey) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	errMsg := "Not Supported operation"
	return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, "")
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *CompositeKey) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.keyName != "" {
		os.(*ProtocolDataOutputStream).WriteBoolean(true) //TypeName exists
		err := os.(*ProtocolDataOutputStream).WriteUTF(obj.keyName)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:WriteExternal - unable to write obj.KeyName w/ Error: '%+v'", err.Error()))
			return err
		}
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
	}
	os.(*ProtocolDataOutputStream).WriteShort(len(obj.attributes))
	for _, attr := range obj.attributes {
		// Null value is not allowed and therefore no need to include isNull flag
		err := attr.WriteExternal(os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning CompositeKey:WriteExternal - unable to write attr w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CompositeKey:WriteExternal w/ NO error, for key: '%+v'", obj))
	}
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
