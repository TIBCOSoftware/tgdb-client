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
 * File Name: attrimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: attrimpl.go 4266 2020-08-19 22:05:56Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"tgdb"
	"time"
)

// ======= Attribute Types =======
const (
	AttributeTypeInvalid = iota
	AttributeTypeBoolean
	AttributeTypeByte
	AttributeTypeChar
	AttributeTypeShort
	AttributeTypeInteger
	AttributeTypeLong
	AttributeTypeFloat
	AttributeTypeDouble
	AttributeTypeNumber
	AttributeTypeString
	AttributeTypeDate
	AttributeTypeTime
	AttributeTypeTimeStamp
	AttributeTypeBlob
	AttributeTypeClob
)

type AttributeType struct {
	typeId      int
	typeName    string
	implementor string
}

var PreDefinedAttributeTypes = map[int]AttributeType{
	AttributeTypeInvalid:   {typeId: AttributeTypeInvalid, typeName: "", implementor: ""},
	AttributeTypeBoolean:   {typeId: AttributeTypeBoolean, typeName: "bool", implementor: "BooleanAttribute"},
	AttributeTypeByte:      {typeId: AttributeTypeByte, typeName: "uint8", implementor: "ByteAttribute"},
	AttributeTypeChar:      {typeId: AttributeTypeChar, typeName: "byte", implementor: "CharAttribute"},
	AttributeTypeShort:     {typeId: AttributeTypeShort, typeName: "int16", implementor: "ShortAttribute"},
	AttributeTypeInteger:   {typeId: AttributeTypeInteger, typeName: "int", implementor: "IntegerAttribute"},
	AttributeTypeLong:      {typeId: AttributeTypeLong, typeName: "int64", implementor: "LongAttribute"},
	AttributeTypeFloat:     {typeId: AttributeTypeFloat, typeName: "float32", implementor: "FloatAttribute"},
	AttributeTypeDouble:    {typeId: AttributeTypeDouble, typeName: "float64", implementor: "DoubleAttribute"},
	AttributeTypeNumber:    {typeId: AttributeTypeNumber, typeName: "big.Int", implementor: "NumberAttribute"},
	AttributeTypeString:    {typeId: AttributeTypeString, typeName: "string", implementor: "StringAttribute"},
	AttributeTypeDate:      {typeId: AttributeTypeDate, typeName: "date", implementor: "TimestampAttribute"},
	AttributeTypeTime:      {typeId: AttributeTypeTime, typeName: "time", implementor: "TimestampAttribute"},
	AttributeTypeTimeStamp: {typeId: AttributeTypeTimeStamp, typeName: "time.Time", implementor: "TimestampAttribute"},
	AttributeTypeBlob:      {typeId: AttributeTypeBlob, typeName: "[]uint8", implementor: "BlobAttribute"},
	AttributeTypeClob:      {typeId: AttributeTypeClob, typeName: "[]rune", implementor: "ClobAttribute"},
}



////////////////////////////////////////
// Helper functions for AttributeType //
////////////////////////////////////////

func (obj *AttributeType) GetTypeId() int {
	return obj.typeId
}

func (obj *AttributeType) GetTypeName() string {
	return obj.typeName
}

func (obj *AttributeType) GetImplementor() string {
	return obj.implementor
}

func (obj *AttributeType) SetTypeId(id int) {
	obj.typeId = id
}

func (obj *AttributeType) SetTypeName(name string) {
	obj.typeName = name
}

func (obj *AttributeType) SetImplementor(impl string) {
	obj.implementor = impl
}

// GetAttributeTypeFromId returns the TGAttributeType given its id
func GetAttributeTypeFromId(aType int) *AttributeType {
	attrObj, ok := PreDefinedAttributeTypes[aType]
	if ok {
		return &attrObj
	} else {
		invalid := PreDefinedAttributeTypes[AttributeTypeInvalid]
		return &invalid
	}
}

// GetAttributeTypeFromName returns the TGAttributeType given its Name
func GetAttributeTypeFromName(aName string) *AttributeType {
	for _, aType := range PreDefinedAttributeTypes {
		if strings.ToLower(aType.typeName) == strings.ToLower(aName) {
			return &aType
		}
	}
	invalid := PreDefinedAttributeTypes[AttributeTypeInvalid]
	return &invalid
}


var LocalAttributeId int64

type AttributeDescriptor struct {
	SysType     tgdb.TGSystemType
	attributeId int64
	Name        string
	AttrType    int
	IsArray     bool
	IsEncrypted bool
	Precision   int16
	Scale       int16
}

func DefaultAttributeDescriptor() *AttributeDescriptor {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AttributeDescriptor{})

	newAttributeDescriptor := AttributeDescriptor{
		SysType:     tgdb.SystemTypeAttributeDescriptor,
		Name:        "",
		AttrType:    AttributeTypeInvalid,
		IsArray:     false,
		IsEncrypted: false,
		Precision:   0,
		Scale:       0,
	}
	newAttributeDescriptor.attributeId = atomic.AddInt64(&LocalAttributeId, 1)
	return &newAttributeDescriptor
}

func NewAttributeDescriptor(id int64) *AttributeDescriptor {
	newAttributeDescriptor := DefaultAttributeDescriptor()
	newAttributeDescriptor.attributeId = id
	return newAttributeDescriptor
}

func NewAttributeDescriptorWithType(name string, attrType int) *AttributeDescriptor {
	newAttributeDescriptor := DefaultAttributeDescriptor()
	newAttributeDescriptor.Name = name
	newAttributeDescriptor.AttrType = attrType
	if attrType == AttributeTypeNumber {
		newAttributeDescriptor.Precision = 20
		newAttributeDescriptor.Scale = 5
	}
	return newAttributeDescriptor
}

func NewAttributeDescriptorAsArray(name string, attrType int, isArray bool) *AttributeDescriptor {
	newAttributeDescriptor := NewAttributeDescriptorWithType(name, attrType)
	newAttributeDescriptor.IsArray = isArray
	return newAttributeDescriptor
}

// TODO: To be used when created from server side data
func NewAttributeDescriptorOnServer(name string, attrType int, isArray bool, attributeId int64) *AttributeDescriptor {
	newAttributeDescriptor := NewAttributeDescriptorAsArray(name, attrType, isArray)
	newAttributeDescriptor.attributeId = attributeId
	return newAttributeDescriptor
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGAbstractEntity
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) SetAttributeId(attrId int64) {
	obj.attributeId = attrId
}

func (obj *AttributeDescriptor) SetAttrType(attrType int) {
	obj.AttrType = attrType
}

func (obj *AttributeDescriptor) SetIsArray(arrayFlag bool) {
	obj.IsArray = arrayFlag
}

func (obj *AttributeDescriptor) SetIsEncrypted(encryptedFlag bool) {
	obj.IsEncrypted = encryptedFlag
}

func (obj *AttributeDescriptor) SetName(attrName string) {
	obj.Name = attrName
}

func (obj *AttributeDescriptor) SetSystemType(sysType tgdb.TGSystemType) {
	obj.SysType = sysType
}

func (obj *AttributeDescriptor) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AttributeDescriptor:{")
	buffer.WriteString(fmt.Sprintf("SysType: %d", obj.SysType))
	buffer.WriteString(fmt.Sprintf(", AttributeId: %d", obj.attributeId))
	buffer.WriteString(fmt.Sprintf(", Name: %s", obj.Name))
	buffer.WriteString(fmt.Sprintf(", AttrType: %s", GetAttributeTypeFromId(obj.AttrType).GetTypeName()))
	buffer.WriteString(fmt.Sprintf(", IsArray: %+v", obj.IsArray))
	buffer.WriteString(fmt.Sprintf(", Precision: %+v", obj.Precision))
	buffer.WriteString(fmt.Sprintf(", Scale: %+v", obj.Scale))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGAttributeDescriptor
/////////////////////////////////////////////////////////////////

// GetAttributeId returns the attributeId
func (obj *AttributeDescriptor) GetAttributeId() int64 {
	return obj.attributeId
}

// GetAttrType returns the type of Attribute Descriptor
func (obj *AttributeDescriptor) GetAttrType() int {
	return obj.AttrType
}

// GetPrecision returns the Precision for Attribute Descriptor of type Number. The default Precision is 20
func (obj *AttributeDescriptor) GetPrecision() int16 {
	return obj.Precision
}

// GetScale returns the Scale for Attribute Descriptor of type Number. The default Scale is 5
func (obj *AttributeDescriptor) GetScale() int16 {
	return obj.Scale
}

// IsAttributeArray checks whether the AttributeType an array desc or not
func (obj *AttributeDescriptor) IsAttributeArray() bool {
	return obj.IsArray
}

// Is_Encrypted checks whether this attribute is Encrypted or not
func (obj *AttributeDescriptor) Is_Encrypted() bool {
	return obj.IsEncrypted
}

// SetPrecision sets the prevision for Attribute Descriptor of type Number
func (obj *AttributeDescriptor) SetPrecision(precision int16) {
	if obj.AttrType == AttributeTypeNumber {
		obj.Precision = precision
	}
}

// SetScale sets the Scale for Attribute Descriptor of type Number
func (obj *AttributeDescriptor) SetScale(scale int16) {
	if obj.AttrType == AttributeTypeNumber {
		obj.Scale = scale
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSystemObject
/////////////////////////////////////////////////////////////////

// GetName gets the system object's Name
func (obj *AttributeDescriptor) GetName() string {
	return obj.Name
}

// GetSystemType gets system object's type
func (obj *AttributeDescriptor) GetSystemType() tgdb.TGSystemType {
	return tgdb.SystemTypeAttributeDescriptor
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *AttributeDescriptor) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	sType, err := is.(*ProtocolDataInputStream).ReadByte() // Read the sysobject desc field which should be 0 for attribute descriptor
	if err != nil {
		return err
	}
	if tgdb.TGSystemType(sType) != tgdb.SystemTypeAttributeDescriptor {
		// TODO: Revisit later - Do we need to throw exception is needed
		//errMsg := fmt.Sprintf("Attribute descriptor has invalid input stream value: %d", sType)
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		logger.Warning(fmt.Sprint("WARNING: AttributeDescriptor:ReadExternal - types.TGSystemType(sType) != types.SystemTypeAttributeDescriptor"))
	}
	attributeId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read attributeId w/ Error: '%+v'", err.Error()))
		return err
	}
	attrName, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read attrName w/ Error: '%+v'", err.Error()))
		return err
	}
	attrType, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read AttrType w/ Error: '%+v'", err.Error()))
		return err
	}
	isArray, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read IsArray w/ Error: '%+v'", err.Error()))
		return err
	}
	isEncrypted, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read Is_Encrypted w/ Error: '%+v'", err.Error()))
		return err
	}
	var precision, scale int16
	if attrType == AttributeTypeNumber {
		precision, err = is.(*ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read Precision w/ Error: '%+v'", err.Error()))
			return err
		}
		scale, err = is.(*ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:ReadExternal - unable to read Scale w/ Error: '%+v'", err.Error()))
			return err
		}
	}
	obj.SetSystemType(tgdb.TGSystemType(sType))
	obj.SetAttributeId(int64(attributeId))
	obj.SetName(attrName)
	obj.SetAttrType(int(attrType))
	obj.SetIsArray(isArray)
	obj.SetIsEncrypted(isEncrypted)
	obj.SetPrecision(precision)
	obj.SetScale(scale)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AttributeDescriptor:ReadExternal w/ NO error, for AttrDesc: '%+v'", obj))
	}
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AttributeDescriptor) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	os.(*ProtocolDataOutputStream).WriteByte(int(tgdb.SystemTypeAttributeDescriptor)) // SysObject desc attribute descriptor
	os.(*ProtocolDataOutputStream).WriteInt(int(obj.GetAttributeId()))
	err := os.(*ProtocolDataOutputStream).WriteUTF(obj.GetName())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:WriteExternal - unable to write AttrDesc Name w/ Error: '%+v'", err.Error()))
		return err
	}
	os.(*ProtocolDataOutputStream).WriteByte(obj.GetAttrType())
	os.(*ProtocolDataOutputStream).WriteBoolean(obj.IsAttributeArray())
	os.(*ProtocolDataOutputStream).WriteBoolean(obj.Is_Encrypted())
	if obj.AttrType == AttributeTypeNumber {
		os.(*ProtocolDataOutputStream).WriteShort(int(obj.GetPrecision()))
		os.(*ProtocolDataOutputStream).WriteShort(int(obj.GetScale()))
	}
	//logger.Log(fmt.Sprintf("Exported Attribute Descriptor object as '%+v' from byte format", obj))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.SysType, obj.attributeId, obj.Name, obj.AttrType, obj.IsArray, obj.Precision, obj.Scale)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *AttributeDescriptor) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.SysType, &obj.attributeId, &obj.Name, &obj.AttrType, &obj.IsArray, &obj.Precision, &obj.Scale)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AttributeDescriptor:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

const (
	DATE_ONLY    = 0
	TIME_ONLY    = 1
	TIMESTAMP    = 2
	TGNoZone     = -1
	TGZoneOffset = 0
	TGZoneId     = 1
	TGZoneName   = 2
)


type AbstractAttribute struct {
	owner      tgdb.TGEntity
	AttrDesc   *AttributeDescriptor
	AttrValue  interface{}
	IsModified bool
}

// Create New Attribute Instance
func defaultNewAbstractAttribute() *AbstractAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractAttribute{})

	newAttribute := AbstractAttribute{
		AttrDesc:   DefaultAttributeDescriptor(),
		IsModified: false,
	}
	return &newAttribute
}

func NewAbstractAttributeWithOwner(ownerEntity tgdb.TGEntity) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewAbstractAttribute(attrDesc *AttributeDescriptor) *AbstractAttribute {
	newAttribute := defaultNewAbstractAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewAbstractAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *AbstractAttribute {
	newAttribute := NewAbstractAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Private functions for AbstractAttribute / TGAttribute
/////////////////////////////////////////////////////////////////

// interfaceEncode encodes the interface value into the encoder.
func interfaceEncode(enc *gob.Encoder, p tgdb.TGAttribute) tgdb.TGError {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.

	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type. If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
		errMsg := "Unable to encode interface"
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	return nil
}

// interfaceDecode decodes the next interface value from the stream and returns it.
func interfaceDecode(dec *gob.Decoder) tgdb.TGAttribute {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p tgdb.TGAttribute
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}

func (obj *AbstractAttribute) getAttributeDescriptor() *AttributeDescriptor {
	return obj.AttrDesc
}

func (obj *AbstractAttribute) getIsModified() bool {
	return obj.IsModified
}

func (obj *AbstractAttribute) getName() string {
	return obj.GetAttributeDescriptor().GetName()
}

func (obj *AbstractAttribute) getOwner() tgdb.TGEntity {
	return obj.owner
}

func (obj *AbstractAttribute) getValue() interface{} {
	return obj.AttrValue
}

func (obj *AbstractAttribute) isNull() bool {
	return obj.AttrValue == nil
}

func (obj *AbstractAttribute) resetIsModified() {
	obj.IsModified = false
}

func (obj *AbstractAttribute) setIsModified(flag bool) {
	obj.IsModified = flag
}

func (obj *AbstractAttribute) setNull() {
	obj.AttrValue = nil
}

func (obj *AbstractAttribute) setOwner(ownerEntity tgdb.TGEntity) {
	obj.owner = ownerEntity
}

func (obj *AbstractAttribute) setValue(value interface{}) tgdb.TGError {
	if value == nil {
		errMsg := fmt.Sprintf("Attribute value is required")
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if !obj.isNull() && obj.getValue() == value {
		return nil
	}

	//var Precision, Scale int16
	//if obj.AttrDesc.GetAttrType() == types.AttributeTypeNumber {
	//	Precision = obj.AttrDesc.GetPrecision()
	//	Scale = obj.AttrDesc.GetScale()
	//}
	//obj.setValueWithPrecisionAndScale(value, Precision, Scale)

	obj.AttrValue = value
	obj.IsModified = true
	return nil
}

func (obj *AbstractAttribute) attributeToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractAttribute:{")
	//buffer.WriteString(fmt.Sprintf("owner: %+v ", obj.owner))
	buffer.WriteString(fmt.Sprintf("AttrDesc: %s", obj.AttrDesc.String()))
	buffer.WriteString(fmt.Sprintf(", AttrValue: %+v", obj.AttrValue))
	buffer.WriteString(fmt.Sprintf(", IsModified: %+v", obj.IsModified))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Helper functions for AbstractAttribute
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) SetIsModified(flag bool) {
	obj.setIsModified(flag)
}

func ReadExternalForEntity(owner tgdb.TGEntity, is tgdb.TGInputStream) (tgdb.TGAttribute, tgdb.TGError) {

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractAttribute:ReadExternalForEntity"))
	}
	attrId, err := is.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractAttribute::ReadExternalForEntity - read attrId: '%+v'", attrId))
	}

	//edge, ok := collection[i].(*impl.Edge)
	var gmd tgdb.TGGraphMetadata
	node, ok := owner.(*Node)
	if ok {
		gmd = node.AbstractEntity.GetGraphMetadata()
	} else {
		edge, ok := owner.(*Edge)
		if ok {
			gmd = edge.AbstractEntity.GetGraphMetadata()
		} else {
			absEntity, ok := owner.(*AbstractEntity)
			if ok {
				gmd = absEntity.GetGraphMetadata()
			}
		}
	}

	//gmd := owner.(*AbstractEntity).GetGraphMetadata()
	if gmd == nil {
		errMsg := fmt.Sprintf("Invalid graph meta data associated with owner: '%+v'", owner.(*AbstractEntity))
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	//logger.Debug(fmt.Sprintf("Inside AbstractAttribute::ReadExternalForEntity - read gmd: '%+v'", gmd))
	attrDesc, err := gmd.GetAttributeDescriptorById(int64(attrId))
	if err != nil {
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractAttribute::ReadExternalForEntity - read AttrDesc: '%+v'", attrDesc))
	}
	if attrDesc == nil {
		errMsg := fmt.Sprintf("Invalid attributeId:'%d' encountered while deserialized", attrId)
		return nil, GetErrorByType(TGErrorIOException, TGDB_CLIENT_READEXTERNAL, errMsg, "")
	}
	newAttr, err := CreateAttributeWithDesc(owner, attrDesc.(*AttributeDescriptor), nil)
	if err != nil {
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractAttribute::ReadExternalForEntity - created new attribute newAttr: '%+v'", newAttr))
	}
	err = newAttr.ReadExternal(is)
	if err != nil {
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractAttribute::ReadExternalForEntity - updated newAttr from stream: '%+v'", newAttr))
	}
	return newAttr, nil
}

func AbstractAttributeReadDecrypted(obj tgdb.TGAttribute, is tgdb.TGInputStream) tgdb.TGError {
	conn := obj.GetOwner().GetGraphMetadata().GetConnection()
	decryptBuf, err := conn.DecryptBuffer(is)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in conn.DecryptBuffer(encryptedBuf): %s", err.Error()))
		return err
	}
	value, err := ObjectFromByteArray(decryptBuf, obj.GetAttributeDescriptor().GetAttrType())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:AbstractAttributeReadDecrypted w/ Error in ObjectFromByteArray(): %s", err.Error()))
		return err
	}
	return obj.SetValue(value)
}

func AbstractAttributeReadExternal(obj tgdb.TGAttribute, is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractAttribute:EntityTypeReadExternal"))
	}
	// We have already read the AttributeId, so no need to read it.
	isNull, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeReadExternal w/ Error in reading isNull from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractAttribute::AbstractAttributeReadExternal - read isnull: '%+v'", isNull))
	}
	if isNull {
		obj.SetValue(nil)
		return nil
	}
	if obj.GetAttributeDescriptor().Is_Encrypted() &&
		obj.GetAttributeDescriptor().GetAttrType() != AttributeTypeBlob &&
		obj.GetAttributeDescriptor().GetAttrType() != AttributeTypeClob {
		return AbstractAttributeReadDecrypted(obj, is)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractAttribute::AbstractAttributeReadExternal after reading the attribute value"))
	}
	return obj.ReadValue(is)
}

func AbstractAttributeWriteEncrypted(obj tgdb.TGAttribute, os tgdb.TGOutputStream) tgdb.TGError {
	buff, err := ObjectToByteArray(obj.GetValue(), obj.GetAttributeDescriptor().GetAttrType())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeWriteEncrypted w/ Error in ObjectToByteArray"))
		return err
	}
	conn := obj.GetOwner().GetGraphMetadata().GetConnection()
	encryptedBuf, err := conn.EncryptEntity(buff)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AbstractAttribute:AbstractAttributeWriteEncrypted w/ Error in conn.EncryptEntity(buff)"))
		return err
	}
	return os.(*ProtocolDataOutputStream).WriteBytes(encryptedBuf)
}

func AbstractAttributeWriteExternal(obj tgdb.TGAttribute, os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractAttribute:AbstractAttributeWriteExternal"))
	}
	attrId := obj.GetAttributeDescriptor().GetAttributeId()
	// Null attribute is not allowed during entity creation
	os.(*ProtocolDataOutputStream).WriteInt(int(attrId))
	os.(*ProtocolDataOutputStream).WriteBoolean(obj.IsNull())
	if obj.IsNull() {
		return nil
	}
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		return AbstractAttributeWriteEncrypted(obj, os)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractAttribute:AbstractAttributeWriteExternal after writing the attribute value"))
	}
	return obj.WriteValue(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *AbstractAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *AbstractAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *AbstractAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *AbstractAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *AbstractAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *AbstractAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *AbstractAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *AbstractAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *AbstractAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.setNull()
		return nil
	}
	return obj.setValue(value)
}

// ReadValue reads the value from input stream
func (obj *AbstractAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	return nil
}

// WriteValue writes the value to output stream
func (obj *AbstractAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	return nil
}

func (obj *AbstractAttribute) String() string {
	return obj.attributeToString()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *AbstractAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *AbstractAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *AbstractAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

type BooleanAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultBooleanAttribute() *BooleanAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BooleanAttribute{})

	newAttribute := BooleanAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	newAttribute.AttrValue = false
	return &newAttribute
}

func NewBooleanAttributeWithOwner(ownerEntity tgdb.TGEntity) *BooleanAttribute {
	newAttribute := DefaultBooleanAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewBooleanAttribute(attrDesc *AttributeDescriptor) *BooleanAttribute {
	newAttribute := DefaultBooleanAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewBooleanAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *BooleanAttribute {
	newAttribute := NewBooleanAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}



/////////////////////////////////////////////////////////////////
// Helper functions for BooleanAttribute
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) SetBoolean(b bool) {
	if !obj.IsNull() && obj.AttrValue.(bool) == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *BooleanAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *BooleanAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *BooleanAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *BooleanAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *BooleanAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *BooleanAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *BooleanAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *BooleanAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *BooleanAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}
	if reflect.TypeOf(value).Kind() != reflect.Bool &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning BooleanAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to BooleanAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := strconv.ParseBool(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BooleanAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to BooleanAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBoolean(v)
	} else {
		obj.SetBoolean(value.(bool))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *BooleanAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	value, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading value from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("BooleanAttribute::ReadValue - read value: '%+v'", value))
	}
	obj.AttrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *BooleanAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.AttrValue == nil {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		var network bytes.Buffer
		enc := gob.NewEncoder(&network)
		err := enc.Encode(obj.AttrValue)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:WriteValue - unable to encode attribute value w/ '%s'", err.Error()))
			errMsg := "AbstractAttribute::WriteExternal - Unable to encode attribute value"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		dec := gob.NewDecoder(&network)
		var v bool
		err = dec.Decode(&v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:WriteValue - unable to decode attribute value w/ '%s'", err.Error()))
			errMsg := "AbstractAttribute::WriteExternal - Unable to decode attribute value"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		os.(*ProtocolDataOutputStream).WriteBoolean(v)
	}
	return nil
}

func (obj *BooleanAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BooleanAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *BooleanAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *BooleanAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *BooleanAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BooleanAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttributeFactory
/////////////////////////////////////////////////////////////////

// CreateAttributeByType creates a new attribute based on the type specified
func CreateAttributeByType(attrTypeId int) (tgdb.TGAttribute, tgdb.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	inputAttrTypeId := attrTypeId

	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case AttributeTypeBoolean:
		return DefaultBooleanAttribute(), nil
	case AttributeTypeByte:
		return DefaultByteAttribute(), nil
	case AttributeTypeChar:
		return DefaultCharAttribute(), nil
	case AttributeTypeShort:
		return DefaultShortAttribute(), nil
	case AttributeTypeInteger:
		return DefaultIntegerAttribute(), nil
	case AttributeTypeLong:
		return DefaultLongAttribute(), nil
	case AttributeTypeFloat:
		return DefaultFloatAttribute(), nil
	case AttributeTypeDouble:
		return DefaultDoubleAttribute(), nil
	case AttributeTypeNumber:
		return DefaultNumberAttribute(), nil
	case AttributeTypeString:
		return DefaultStringAttribute(), nil
	case AttributeTypeDate:
		return DefaultTimestampAttribute(), nil
	case AttributeTypeTime:
		return DefaultTimestampAttribute(), nil
	case AttributeTypeTimeStamp:
		return DefaultTimestampAttribute(), nil
	case AttributeTypeBlob:
		return DefaultBlobAttribute(), nil
	case AttributeTypeClob:
		return DefaultClobAttribute(), nil
	case AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return nil, nil
}

// CreateAttribute creates a new attribute based on AttributeDescriptor
func CreateAttribute(attrDesc *AttributeDescriptor) (tgdb.TGAttribute, tgdb.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	attrTypeId := attrDesc.GetAttrType()
	inputAttrTypeId := attrTypeId

	var newAttribute tgdb.TGAttribute
	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case AttributeTypeBoolean:
		// Execute Individual Attribute's method
		newAttribute = NewBooleanAttribute(attrDesc)
	case AttributeTypeByte:
		// Execute Individual Attribute's method
		newAttribute = NewByteAttribute(attrDesc)
	case AttributeTypeChar:
		// Execute Individual Attribute's method
		newAttribute = NewCharAttribute(attrDesc)
	case AttributeTypeShort:
		// Execute Individual Attribute's method
		newAttribute = NewShortAttribute(attrDesc)
	case AttributeTypeInteger:
		// Execute Individual Attribute's method
		newAttribute = NewIntegerAttribute(attrDesc)
	case AttributeTypeLong:
		// Execute Individual Attribute's method
		newAttribute = NewLongAttribute(attrDesc)
	case AttributeTypeFloat:
		// Execute Individual Attribute's method
		newAttribute = NewFloatAttribute(attrDesc)
	case AttributeTypeDouble:
		// Execute Individual Attribute's method
		newAttribute = NewDoubleAttribute(attrDesc)
	case AttributeTypeNumber:
		// Execute Individual Attribute's method
		newAttribute = NewNumberAttribute(attrDesc)
	case AttributeTypeString:
		// Execute Individual Attribute's method
		newAttribute = NewStringAttribute(attrDesc)
	case AttributeTypeDate:
		newAttribute = NewTimestampAttribute(attrDesc)
	case AttributeTypeTime:
		newAttribute = NewTimestampAttribute(attrDesc)
	case AttributeTypeTimeStamp:
		// Execute Individual Attribute's method
		newAttribute = NewTimestampAttribute(attrDesc)
	case AttributeTypeBlob:
		// Execute Individual Attribute's method
		newAttribute = NewBlobAttribute(attrDesc)
	case AttributeTypeClob:
		// Execute Individual Attribute's method
		newAttribute = NewClobAttribute(attrDesc)
	case AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return newAttribute, nil
}

// CreateAttributeWithDesc creates new attribute based on the owner and AttributeDescriptor
func CreateAttributeWithDesc(attrOwner tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) (tgdb.TGAttribute, tgdb.TGError) {
	// Store incoming identifier, in case there is a need to find more dependency or massaging
	attrTypeId := attrDesc.GetAttrType()
	inputAttrTypeId := attrTypeId

	var newAttribute tgdb.TGAttribute
	// Use a switch case to switch between attribute types, if a type exist then error is nil (null)
	// Whenever new attribute type gets into the mix, just add a case below
	switch inputAttrTypeId {
	case AttributeTypeBoolean:
		// Execute Individual Attribute's method
		newAttribute = NewBooleanAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeByte:
		// Execute Individual Attribute's method
		newAttribute = NewByteAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeChar:
		// Execute Individual Attribute's method
		newAttribute = NewCharAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeShort:
		// Execute Individual Attribute's method
		newAttribute = NewShortAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeInteger:
		// Execute Individual Attribute's method
		newAttribute = NewIntegerAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeLong:
		// Execute Individual Attribute's method
		newAttribute = NewLongAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeFloat:
		// Execute Individual Attribute's method
		newAttribute = NewFloatAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeDouble:
		// Execute Individual Attribute's method
		newAttribute = NewDoubleAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeNumber:
		// Execute Individual Attribute's method
		newAttribute = NewNumberAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeString:
		// Execute Individual Attribute's method
		newAttribute = NewStringAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeDate:
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeTime:
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeTimeStamp:
		// Execute Individual Attribute's method
		newAttribute = NewTimestampAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeBlob:
		// Execute Individual Attribute's method
		newAttribute = NewBlobAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeClob:
		// Execute Individual Attribute's method
		newAttribute = NewClobAttributeWithDesc(attrOwner, attrDesc, value)
	case AttributeTypeInvalid:
		fallthrough
	default:
		//if type is invalid, return an error
		errMsg := fmt.Sprintf("AttributeTypeInvalid Attribute Type '%s'", GetAttributeTypeFromId(inputAttrTypeId).GetTypeName())
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	return newAttribute, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute type
func GetAttributeDescriptor(attrTypeId int) tgdb.TGAttributeDescriptor {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return nil
	}
	// Execute Individual Attribute's method
	return attr.GetAttributeDescriptor()
}

// GetIsModified checks whether the attribute of this type is modified or not
func GetIsModified(attrTypeId int) bool {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetIsModified()
}

// ResetIsModified resets the IsModified flag of this attribute type - recursively, if needed
func ResetIsModified(attrTypeId int) {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return
	}
	// Execute Individual Attribute's method
	attr.ResetIsModified()
}

// GetName gets the Name for this attribute type as the most generic form
func GetName(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetName()
}

// GetOwner gets owner Entity of this attribute type
func GetOwner(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return nil
	}
	// Execute Individual Attribute's method
	return attr.GetOwner()
}

// GetValue gets the value for this attribute type as the most generic form
func GetValue(attrTypeId int) interface{} {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.GetValue()
}

// IsNull checks whether the value of this attribute type is null or not
func IsNull(attrTypeId int) bool {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return false
	}
	// Execute Individual Attribute's method
	return attr.IsNull()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func SetOwner(attrTypeId int, attrOwner tgdb.TGEntity) {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return
	}
	// Execute Individual Attribute's method
	attr.SetOwner(attrOwner)
}

// SetValue sets the value for this attribute type. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func SetValue(attrTypeId int, value interface{}) tgdb.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.SetValue(value)
}

// ReadValue reads the value of this attribute type from input stream
func ReadValue(attrTypeId int, is tgdb.TGInputStream) tgdb.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.ReadValue(is)
}

// WriteValue writes the value of this attribute type to output stream
func WriteValue(attrTypeId int, os tgdb.TGOutputStream) tgdb.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.WriteValue(os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func ReadExternal(attrTypeId int, is tgdb.TGInputStream) tgdb.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.ReadExternal(is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func WriteExternal(attrTypeId int, os tgdb.TGOutputStream) tgdb.TGError {
	attr, err := CreateAttributeByType(attrTypeId)
	if err != nil {
		return err
	}
	// Execute Individual Attribute's method
	return attr.WriteExternal(os)
}


type ByteAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultByteAttribute() *ByteAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ByteAttribute{})

	newAttribute := ByteAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewByteAttributeWithOwner(ownerEntity tgdb.TGEntity) *ByteAttribute {
	newAttribute := DefaultByteAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewByteAttribute(attrDesc *AttributeDescriptor) *ByteAttribute {
	newAttribute := DefaultByteAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewByteAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *ByteAttribute {
	newAttribute := NewByteAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for ByteAttribute
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) SetByte(b uint8) {
	if !obj.IsNull() {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *ByteAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *ByteAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *ByteAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *ByteAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *ByteAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *ByteAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *ByteAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *ByteAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *ByteAttribute) SetValue(value interface{}) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering ByteAttribute::SetValue trying to set attribute value '%+v' of type '%+v'", value, reflect.TypeOf(value).Kind()))
	}
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(value)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to encode attribute value"))
		errMsg := "ByteAttribute::SetValue - Unable to encode attribute value"
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)

	if reflect.TypeOf(value).Kind() != reflect.Bool &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.Uint &&
		reflect.TypeOf(value).Kind() != reflect.Uint8 &&
		reflect.TypeOf(value).Kind() != reflect.Uint16 &&
		reflect.TypeOf(value).Kind() != reflect.Uint32 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to ByteAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.Bool {
		var v bool
		err = dec.Decode(&v)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to decode attribute value"))
			errMsg := "ByteAttribute::SetValue - Unable to decode attribute value"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
		}
		if v {
			obj.SetByte(1)
		} else {
			obj.SetByte(0)
		}
	} else if reflect.TypeOf(value).Kind() == reflect.String ||
		reflect.TypeOf(value).Kind() == reflect.Float32 ||
		reflect.TypeOf(value).Kind() == reflect.Float64 {
		v := value.(string)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to ByteAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		v1, _ := strconv.Atoi(v)
		obj.SetByte(uint8(v1))
	} else {
		v := uint8(reflect.ValueOf(value).Uint())
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning CharAttribute::SetValue - trying to set attribute value '%+v' of type '%+v'", v, reflect.TypeOf(v).Kind()))
		}
		obj.SetByte(v)
	}

	return nil
}

// ReadValue reads the value from input stream
func (obj *ByteAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	value, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:ReadValue w/ Error in reading value from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ByteAttribute::ReadValue - read value: '%+v'", value))
	}
	obj.AttrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *ByteAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(obj.AttrValue)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:WriteValue - unable to encode attribute value"))
		errMsg := "AbstractAttribute::WriteExternal - Unable to encode attribute value"
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)
	var v byte
	err = dec.Decode(&v)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ByteAttribute:WriteValue - unable to decode attribute value"))
		errMsg := "AbstractAttribute::WriteExternal - Unable to decode attribute value"
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	os.(*ProtocolDataOutputStream).WriteByte(int(v))
	return nil
}

func (obj *ByteAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ByteAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ByteAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ByteAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ByteAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ByteAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type CharAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultCharAttribute() *CharAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CharAttribute{})

	newAttribute := CharAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewCharAttributeWithOwner(ownerEntity tgdb.TGEntity) *CharAttribute {
	newAttribute := DefaultCharAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewCharAttribute(attrDesc *AttributeDescriptor) *CharAttribute {
	newAttribute := DefaultCharAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewCharAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *CharAttribute {
	newAttribute := NewCharAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for CharAttribute
/////////////////////////////////////////////////////////////////

func (obj *CharAttribute) SetChar(v int32) {
	if !obj.IsNull() && obj.AttrValue == v {
		return
	}
	obj.AttrValue = v
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *CharAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *CharAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *CharAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *CharAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *CharAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *CharAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *CharAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *CharAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *CharAttribute) SetValue(value interface{}) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering CharAttribute::SetValue trying to set attribute value '%+v' of type '%+v'", value, reflect.TypeOf(value).Kind()))
	}
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning CharAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to DoubleAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToInteger(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning CharAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to DoubleAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetChar(int32(v))
	} else if reflect.TypeOf(value).Kind() == reflect.Int32 {
		v := reflect.ValueOf(value).Int()
		//logger.Log(fmt.Sprintf("Returning CharAttribute::SetValue trying to set attribute value '%+v' of type '%+v'", v, reflect.TypeOf(v).Kind()))
		obj.SetChar(int32(v))
	} else {
		//logger.Log(fmt.Sprintf("Returning CharAttribute::SetValue finally trying to set attribute value '%+v' of type '%+v'", value, reflect.TypeOf(value).Kind()))
		obj.SetChar(value.(int32))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *CharAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	value, err := is.(*ProtocolDataInputStream).ReadChar()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning CharAttribute:ReadValue w/ Error in reading value from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning CharAttribute::ReadValue - read value: '%+v'", value))
	}
	obj.AttrValue = value
	return nil
}

// WriteValue writes the value to output stream
func (obj *CharAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering CharAttribute::WriteValue trying to write attribute value '%+v'", obj.AttrValue))
	}
	iValue := reflect.ValueOf(obj.AttrValue).Int()
	os.(*ProtocolDataOutputStream).WriteChar(int(iValue))
	return nil
}

func (obj *CharAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CharAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *CharAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *CharAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *CharAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CharAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *CharAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CharAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type DoubleAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultDoubleAttribute() *DoubleAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DoubleAttribute{})

	newAttribute := DoubleAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewDoubleAttributeWithOwner(ownerEntity tgdb.TGEntity) *DoubleAttribute {
	newAttribute := DefaultDoubleAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewDoubleAttribute(attrDesc *AttributeDescriptor) *DoubleAttribute {
	newAttribute := DefaultDoubleAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewDoubleAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *DoubleAttribute {
	newAttribute := NewDoubleAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for DoubleAttribute
/////////////////////////////////////////////////////////////////

func (obj *DoubleAttribute) SetDouble(b float64) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *DoubleAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *DoubleAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *DoubleAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *DoubleAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *DoubleAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *DoubleAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *DoubleAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *DoubleAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *DoubleAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int &&
		reflect.TypeOf(value).Kind() != reflect.Int16 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning DoubleAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to DoubleAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToDouble(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DoubleAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to DoubleAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		obj.SetDouble(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v := reflect.ValueOf(value).Float()
		obj.SetDouble(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Int64 {
		v := reflect.ValueOf(value).Float()
		obj.SetDouble(float64(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int32 || reflect.TypeOf(value).Kind() != reflect.Int {
		v := reflect.ValueOf(value).Float()
		obj.SetDouble(float64(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int16 {
		v := reflect.ValueOf(value).Int()
		obj.SetDouble(float64(v))
	} else {
		obj.SetDouble(value.(float64))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *DoubleAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DoubleAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadDouble()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DoubleAttribute:ReadValue w/ Error in reading value from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("DoubleAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *DoubleAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning DoubleAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "DoubleAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		iValue := reflect.ValueOf(obj.AttrValue).Float()
		os.(*ProtocolDataOutputStream).WriteDouble(iValue)
	}
	return nil
}

func (obj *DoubleAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DoubleAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *DoubleAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *DoubleAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *DoubleAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DoubleAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *DoubleAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DoubleAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type FloatAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultFloatAttribute() *FloatAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(FloatAttribute{})

	newAttribute := FloatAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewFloatAttributeWithOwner(ownerEntity tgdb.TGEntity) *FloatAttribute {
	newAttribute := DefaultFloatAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewFloatAttribute(attrDesc *AttributeDescriptor) *FloatAttribute {
	newAttribute := DefaultFloatAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewFloatAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *FloatAttribute {
	newAttribute := NewFloatAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for FloatAttribute
/////////////////////////////////////////////////////////////////

func (obj *FloatAttribute) SetFloat(b float32) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *FloatAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *FloatAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *FloatAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *FloatAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *FloatAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *FloatAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *FloatAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *FloatAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *FloatAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int &&
		reflect.TypeOf(value).Kind() != reflect.Int16 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning FloatAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to DoubleAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToFloat(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning FloatAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to DoubleAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetFloat(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		v := reflect.ValueOf(value).Float()
		obj.SetFloat(float32(v))
	} else if reflect.TypeOf(value).Kind() == reflect.Int64 {
		v := reflect.ValueOf(value).Float()
		obj.SetFloat(float32(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int32 || reflect.TypeOf(value).Kind() != reflect.Int {
		v := reflect.ValueOf(value).Float()
		obj.SetFloat(float32(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int16 {
		v := reflect.ValueOf(value).Float()
		obj.SetFloat(float32(v))
	} else {
		obj.SetFloat(value.(float32))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *FloatAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning FloatAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadFloat()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning FloatAttribute:ReadValue w/ Error in reading value from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning FloatAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *FloatAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning FloatAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "FloatAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		iValue := reflect.ValueOf(obj.AttrValue).Float()
		os.(*ProtocolDataOutputStream).WriteFloat(float32(iValue))
	}
	return nil
}

func (obj *FloatAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("FloatAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *FloatAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *FloatAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (obj *FloatAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning FloatAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaler
/////////////////////////////////////////////////////////////////

func (obj *FloatAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning FloatAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

type IntegerAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultIntegerAttribute() *IntegerAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(IntegerAttribute{})

	newAttribute := IntegerAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewIntegerAttributeWithOwner(ownerEntity tgdb.TGEntity) *IntegerAttribute {
	newAttribute := DefaultIntegerAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewIntegerAttribute(attrDesc *AttributeDescriptor) *IntegerAttribute {
	newAttribute := DefaultIntegerAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewIntegerAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *IntegerAttribute {
	newAttribute := NewIntegerAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for IntegerAttribute
/////////////////////////////////////////////////////////////////

func round64(val float64) int {
	if val < 0 { return int(val-0.5) }
	return int(val+0.5)
}

func round32(val float32) int {
	if val < 0 { return int(val-0.5) }
	return int(val+0.5)
}

func (obj *IntegerAttribute) SetInteger(b int) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *IntegerAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *IntegerAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *IntegerAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *IntegerAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *IntegerAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *IntegerAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *IntegerAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *IntegerAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *IntegerAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int &&
		reflect.TypeOf(value).Kind() != reflect.Int16 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning IntegerAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to IntegerAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToInteger(reflect.ValueOf(value).String())
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning IntegerAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to IntegerAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetInteger(int(v))
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v := reflect.ValueOf(value).Float()
		obj.SetInteger(round32(float32(v)))
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		v := reflect.ValueOf(value).Float()
		obj.SetInteger(round64(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int32 {
		v := reflect.ValueOf(value).Int()
		obj.SetInteger(int(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int16 {
		v := reflect.ValueOf(value).Int()
		obj.SetInteger(int(v))
	} else {
		obj.SetInteger(value.(int))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *IntegerAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning IntegerAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning IntegerAttribute:SetValue - unable to extract attribute value in string format/type"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning IntegerAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *IntegerAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning IntegerAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "IntegerAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		iValue := reflect.ValueOf(obj.AttrValue).Int()
		os.(*ProtocolDataOutputStream).WriteInt(int(iValue))
	}
	return nil
}

func (obj *IntegerAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("IntegerAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *IntegerAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *IntegerAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *IntegerAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IntegerAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *IntegerAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IntegerAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return err
}



type LongAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultLongAttribute() *LongAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(LongAttribute{})

	newAttribute := LongAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewLongAttributeWithOwner(ownerEntity tgdb.TGEntity) *LongAttribute {
	newAttribute := DefaultLongAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewLongAttribute(attrDesc *AttributeDescriptor) *LongAttribute {
	newAttribute := DefaultLongAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewLongAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *LongAttribute {
	newAttribute := NewLongAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for LongAttribute
/////////////////////////////////////////////////////////////////

func lRound64(val float64) int64 {
	if val < 0 { return int64(val-0.5) }
	return int64(val+0.5)
}

func lRound32(val float32) int64 {
	if val < 0 { return int64(val-0.5) }
	return int64(val+0.5)
}

func (obj *LongAttribute) SetLong(b int64) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *LongAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *LongAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *LongAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *LongAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *LongAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *LongAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *LongAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *LongAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *LongAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int &&
		reflect.TypeOf(value).Kind() != reflect.Int16 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Int64 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning LongAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to LongAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToLong(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning LongAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to LongAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetLong(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v := reflect.ValueOf(value).Float()
		obj.SetLong(lRound32(float32(v)))
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		v := reflect.ValueOf(value).Float()
		obj.SetLong(lRound64(float64(v)))
	} else if reflect.TypeOf(value).Kind() != reflect.Int {
		v := reflect.ValueOf(value).Int()
		obj.SetLong(int64(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int32 {
		v := reflect.ValueOf(value).Int()
		obj.SetLong(int64(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int16 {
		v := reflect.ValueOf(value).Int()
		obj.SetLong(int64(v))
	} else {
		obj.SetLong(value.(int64))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *LongAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning LongAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning LongAttribute:ReadValue w/ Error in reading value from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning LongAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *LongAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning LongAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "LongAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		iValue := reflect.ValueOf(obj.AttrValue).Int()
		os.(*ProtocolDataOutputStream).WriteLong(iValue)
	}
	return nil
}

func (obj *LongAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("LongAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *LongAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *LongAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *LongAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning LongAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *LongAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning LongAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return err
}


type NumberAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultNumberAttribute() *NumberAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NumberAttribute{})

	newAttribute := NumberAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	newAttribute.GetAttributeDescriptor().SetPrecision(20)
	newAttribute.GetAttributeDescriptor().SetScale(5)
	return &newAttribute
}

func NewNumberAttributeWithOwner(ownerEntity tgdb.TGEntity) *NumberAttribute {
	newAttribute := DefaultNumberAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewNumberAttribute(attrDesc *AttributeDescriptor) *NumberAttribute {
	newAttribute := DefaultNumberAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewNumberAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *NumberAttribute {
	newAttribute := NewNumberAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for NumberAttribute
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) SetDecimal(b TGDecimal, precision, scale int) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::SetDecimal about to set value '%+v' w/ Precision:Scale %d/%d", b, precision, scale))
	}
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.GetAttributeDescriptor().SetPrecision(int16(precision))
	obj.GetAttributeDescriptor().SetScale(int16(scale))
	obj.AttrValue = b.String() //strings.Replace(b.String(), ".", "", -1)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::SetDecimal AttrValue is '%+v'", obj.AttrValue))
	}
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *NumberAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *NumberAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *NumberAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *NumberAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *NumberAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *NumberAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *NumberAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *NumberAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *NumberAttribute) SetValue(value interface{}) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::SetValue about to set value '%+v' of kind '%+v'", value, reflect.TypeOf(value).Kind()))
	}
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	precision := 20 	// Default Precision
	scale := 5			// Default Scale

	if 	reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Int64 &&
		reflect.TypeOf(value).Kind() != reflect.Struct &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to NumberAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v1 := value.(string)
		parts := strings.Split(v1, ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		//v, err := strconv.ParseFloat(v1, 64)
		v, err := NewTGDecimalFromString(v1)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to NumberAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetDecimal(v, precision, scale)
	} else if reflect.TypeOf(value).Kind() == reflect.Int32 ||
		reflect.TypeOf(value).Kind() == reflect.Int64 {
		v1 := value.(int64)
		scale = 0
		precision = len(strconv.Itoa(int(v1)))
		obj.SetDecimal(NewTGDecimal(value.(int64),0), precision, scale)
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 ||
		reflect.TypeOf(value).Kind() == reflect.Float64 {
		v1 := value.(float64)
		v2 := fmt.Sprintf("%f", v1)
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside NumberAttribute::SetValue - read v1: '%f' and v2: '%s'", v1, v2))
		}
		parts := strings.Split(v2, ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		obj.SetDecimal(NewTGDecimalFromFloatWithExponent(v1, 0), precision, scale)
	} else {
		v1 := value.(TGDecimal)
		parts := strings.Split(v1.String(), ".")
		if len(parts) == 1 {
			// There is no decimal point, we can just parse the original string as an int
			scale = 0
			precision = len(parts[0])
		} else if len(parts) == 2 {
			scale = len(parts[1])
			precision = len(parts[0]) + len(parts[1])
		}
		obj.SetDecimal(value.(TGDecimal), precision, scale)
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *NumberAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ShortAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
		return nil
	}

	precision, err := is.(*ProtocolDataInputStream).ReadShort() // Precision
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading Precision from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::ReadValue - read Precision: '%+v'", precision))
	}
	scale, err := is.(*ProtocolDataInputStream).ReadShort() // Scale
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading Scale from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::ReadValue - read Scale: '%+v'", scale))
	}

	bdStr, err := is.(*ProtocolDataInputStream).ReadUTF()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning NumberAttribute:ReadValue w/ Error in reading bdStr from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::ReadValue - read bdStr: '%+v'", bdStr))
	}
	obj.AttrValue, _ = NewTGDecimalFromString(bdStr)
	return nil
}

// WriteValue writes the value to output stream
func (obj *NumberAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "StringAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
		return nil
	}

	os.(*ProtocolDataOutputStream).WriteShort(int(obj.GetAttributeDescriptor().GetPrecision()))
	os.(*ProtocolDataOutputStream).WriteShort(int(obj.GetAttributeDescriptor().GetScale()))
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::WriteValue AttrValue is '%+v'", obj.AttrValue))
	}
	//dValue := reflect.ValueOf(obj.AttrValue).Float()
	//strValue := strconv.FormatFloat(dValue, 'f', int(obj.GetAttributeDescriptor().GetPrecision()), 64)
	//newStr := strings.Replace(strValue, ".", "", -1)
	dValue := obj.AttrValue.(string)
	newStr := dValue	//strings.Replace(dValue, ".", "", -1)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside NumberAttribute::WriteValue newStr is '%+v'", newStr))
	}
	return os.(*ProtocolDataOutputStream).WriteUTF(newStr)
}

func (obj *NumberAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NumberAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *NumberAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *NumberAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NumberAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NumberAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NumberAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ShortAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultShortAttribute() *ShortAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ShortAttribute{})

	newAttribute := ShortAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewShortAttributeWithOwner(ownerEntity tgdb.TGEntity) *ShortAttribute {
	newAttribute := DefaultShortAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewShortAttribute(attrDesc *AttributeDescriptor) *ShortAttribute {
	newAttribute := DefaultShortAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewShortAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *ShortAttribute {
	newAttribute := NewShortAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for ShortAttribute
/////////////////////////////////////////////////////////////////

func sRound64(val float64) int16 {
	if val < 0 { return int16(val-0.5) }
	return int16(val+0.5)
}

func sRound32(val float32) int16 {
	if val < 0 { return int16(val-0.5) }
	return int16(val+0.5)
}

func (obj *ShortAttribute) SetShort(b int16) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *ShortAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *ShortAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *ShortAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *ShortAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *ShortAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *ShortAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *ShortAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *ShortAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *ShortAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Int &&
		reflect.TypeOf(value).Kind() != reflect.Int16 &&
		reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning ShortAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to ShortAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToShort(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ShortAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprintf("Failure to covert string to ShortAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetShort(int16(v))
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v := reflect.ValueOf(value).Float()
		obj.SetShort(sRound32(float32(v)))
	} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v := reflect.ValueOf(value).Float()
		obj.SetShort(sRound64(float64(v)))
	} else if reflect.TypeOf(value).Kind() != reflect.Int32 {
		v := reflect.ValueOf(value).Int()
		obj.SetShort(int16(v))
	} else if reflect.TypeOf(value).Kind() != reflect.Int {
		v := reflect.ValueOf(value).Int()
		obj.SetShort(int16(v))
	} else {
		obj.SetShort(value.(int16))
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *ShortAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ShortAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning ShortAttribute:ReadValue w/ Error in reading value from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("ShortAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *ShortAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ShortAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "ShortAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		iValue := reflect.ValueOf(obj.AttrValue).Int()
		os.(*ProtocolDataOutputStream).WriteShort(int(iValue))
	}
	return nil
}

func (obj *ShortAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ShortAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ShortAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ShortAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ShortAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ShortAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ShortAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TimestampAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

const MaxStringAttrLength = 1000 - 2 - 1 // 2 ==> size of int16 or short

type StringAttribute struct {
	*AbstractAttribute
}

// Create New Attribute Instance
func DefaultStringAttribute() *StringAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(StringAttribute{})

	newAttribute := StringAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewStringAttributeWithOwner(ownerEntity tgdb.TGEntity) *StringAttribute {
	newAttribute := DefaultStringAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewStringAttribute(attrDesc *AttributeDescriptor) *StringAttribute {
	newAttribute := DefaultStringAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewStringAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *StringAttribute {
	newAttribute := NewStringAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for StringAttribute
/////////////////////////////////////////////////////////////////

func UtfLength(str string) int {
	strLen := len(str)
	utfLen := 0
	for i := 0; i < strLen; i++ {
		c := str[i]
		if c >= 0x0001 && c <= 0x007F {
			utfLen++
		} else if int(c) > 0x07FF {
			utfLen += 3
		} else {
			utfLen += 2
		}
	}
	return utfLen
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *StringAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *StringAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *StringAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *StringAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *StringAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *StringAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *StringAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *StringAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *StringAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	s := value.(string)
	strLen := UtfLength(s)
	if strLen > MaxStringAttrLength {
		logger.Error(fmt.Sprint("ERROR: Returning StringAttribute:SetValue as strLen > MaxStringAttrLength"))
		errMsg := fmt.Sprintf("UTF length of String exceed the maximum string length supported for a String Attribute (%d > %d", strLen, MaxStringAttrLength)
		return GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.AttrValue = value.(string)
	obj.setIsModified(true)
	return nil
}

// ReadValue reads the value from input stream
func (obj *StringAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning StringAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
	} else {
		value, err := is.(*ProtocolDataInputStream).ReadUTF()
		if err != nil {
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("StringAttribute::ReadValue - read value: '%+v'", value))
		}
		obj.AttrValue = value
	}
	return nil
}

// WriteValue writes the value to output stream
func (obj *StringAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "StringAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	} else {
		return os.(*ProtocolDataOutputStream).WriteUTF(obj.AttrValue.(string))
	}
	return nil
}

func (obj *StringAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("StringAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *StringAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *StringAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (obj *StringAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *StringAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type TimestampAttribute struct {
	*AbstractAttribute
}

// Create NewTGDecimal Attribute Instance
func DefaultTimestampAttribute() *TimestampAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TimestampAttribute{})

	newAttribute := TimestampAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
	}
	return &newAttribute
}

func NewTimestampAttributeWithOwner(ownerEntity tgdb.TGEntity) *TimestampAttribute {
	newAttribute := DefaultTimestampAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewTimestampAttribute(attrDesc *AttributeDescriptor) *TimestampAttribute {
	newAttribute := DefaultTimestampAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewTimestampAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *TimestampAttribute {
	newAttribute := NewTimestampAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for TimestampAttribute
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) SetCalendar(b time.Time) {
	if ! obj.IsNull() {
		return
	}
	obj.AttrValue = b
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *TimestampAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *TimestampAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *TimestampAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *TimestampAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *TimestampAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *TimestampAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *TimestampAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *TimestampAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *TimestampAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering TimestampAttribute:SetValue w/ Input value '%+v' is of type: '%+v'\n", value, reflect.TypeOf(value).Kind()))
	}
	if reflect.TypeOf(value).Kind() != reflect.Int32 &&
		reflect.TypeOf(value).Kind() != reflect.Int64 &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprint("Failure to cast the attribute value to TimestampAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.String {
		v, err := StringToCalendar(value.(string))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:SetValue - unable to extract attribute value in string format/type"))
			errMsg := fmt.Sprint("Failure to covert string to TimestampAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TimestampAttribute:SetValue - Transformed value '%+v' is of type: '%+v'\n", v, reflect.TypeOf(v).Kind()))
		}
		obj.SetCalendar(v)
		return nil
	} else if reflect.TypeOf(value).Kind() == reflect.Int32 ||
		reflect.TypeOf(value).Kind() == reflect.Int64 {
		v := LongToCalendar(value.(int64))
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TimestampAttribute:SetValue - Transformed value '%+v' is of type: '%+v'\n", v, reflect.TypeOf(v).Kind()))
		}
		obj.SetCalendar(v)
		return nil
	} else {
		obj.AttrValue = value
		obj.setIsModified(true)
	}
	return nil
}

// ReadValue reads the value from input stream
func (obj *TimestampAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TimestampAttribute::ReadValue"))
	}
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeReadDecrypted(obj, is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning StringAttribute:ReadValue w/ Error in AbstractAttributeReadDecrypted()"))
			return err
		}
		return nil
	}

	var v time.Time
	var year, mon, dom, hr, min, sec, ms, tzType int
	era, err := is.(*ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading era from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read era: '%+v'", era))
	}

	yr, err := is.(*ProtocolDataInputStream).ReadShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading yr from message buffer"))
		return err
	}
	year = int(yr)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read year: '%+v'", year))
	}

	mth, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mth from message buffer"))
		return err
	}
	mon = int(mth)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read mth: '%+v'", mth))
	}

	day, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading day from message buffer"))
		return err
	}
	dom = int(day)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read dom: '%+v'", dom))
	}

	hour, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading hour from message buffer"))
		return err
	}
	hr = int(hour)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read hr: '%+v'", hr))
	}

	mts, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mts from message buffer"))
		return err
	}
	min = int(mts)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read min: '%+v'", min))
	}

	secs, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading secs from message buffer"))
		return err
	}
	sec = int(secs)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read sec: '%+v'", sec))
	}

	mSec, err := is.(*ProtocolDataInputStream).ReadUnsignedShort()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading mSec from message buffer"))
		return err
	}
	ms = int(mSec)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read ms: '%+v'", ms))
	}

	tz, err := is.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading tz from message buffer"))
		return err
	}
	tzType = int(tz)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read tz: '%+v' and tzType: '%+v'", tz, tzType))
	}

	if tzType != -1 && tzType != 255 {
		tzId, err := is.(*ProtocolDataInputStream).ReadShort()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning TimestampAttribute:ReadValue w/ Error in reading tzId from message buffer"))
			return err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - read tzId: '%+v'", tzId))
		}
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::ReadValue - attribute type: '%+v'", obj.AttrDesc.GetAttrType()))
	}
	switch obj.AttrDesc.GetAttrType() {
	case AttributeTypeDate:
		v = time.Date(year, time.Month(mon), dom, 0, 0, 0, 0, time.Local)
		break
	case AttributeTypeTime:
		v = time.Date(1970, time.January, 01, hr, min, sec, ms*1000, time.Local)
		break
	case AttributeTypeTimeStamp:
		v = time.Date(year, time.Month(mon), dom, hr, min, sec, ms*1000, time.Local)
		break
	default:
		errMsg := fmt.Sprintf("Bad Descriptor: %s", string(obj.AttrDesc.GetAttrType()))
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning TimestampAttribute::ReadValue - read v: '%+v'", v))
	}

	//obj.AttrValue = v.In(loc)		// TODO: Revisit later to use this once location/zone information is available
	obj.AttrValue = v
	return nil
}

// WriteValue writes the value to output stream
func (obj *TimestampAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TimestampAttribute::WriteValue"))
	}
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		err := AbstractAttributeWriteEncrypted(obj, os)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning StringAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
			errMsg := "StringAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
		return nil
	}

	era := true // Corresponding to GregorianCalendar.AD = 1
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside TimestampAttribute::WriteValue - Object value '%+v' is of type: '%+v'\n", obj.GetValue(), reflect.TypeOf(obj.GetValue()).Kind()))
	}
	v := obj.GetValue().(time.Time)
	yr, mth, day := v.Date()
	hr := v.Hour()
	min := v.Minute()
	sec := v.Second()
	msec := v.Nanosecond() / 1000
	switch obj.AttrDesc.GetAttrType() {
	case AttributeTypeDate:
		os.(*ProtocolDataOutputStream).WriteBoolean(era)
		os.(*ProtocolDataOutputStream).WriteShort(yr)
		os.(*ProtocolDataOutputStream).WriteByte(int(mth))
		os.(*ProtocolDataOutputStream).WriteByte(day)
		os.(*ProtocolDataOutputStream).WriteByte(0)
		os.(*ProtocolDataOutputStream).WriteByte(0)
		os.(*ProtocolDataOutputStream).WriteByte(0)
		os.(*ProtocolDataOutputStream).WriteShort(0)
		os.(*ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	case AttributeTypeTime:
		os.(*ProtocolDataOutputStream).WriteBoolean(era)
		os.(*ProtocolDataOutputStream).WriteShort(0)
		os.(*ProtocolDataOutputStream).WriteByte(0)
		os.(*ProtocolDataOutputStream).WriteByte(0)
		os.(*ProtocolDataOutputStream).WriteByte(hr)
		os.(*ProtocolDataOutputStream).WriteByte(min)
		os.(*ProtocolDataOutputStream).WriteByte(sec)
		os.(*ProtocolDataOutputStream).WriteShort(msec)
		os.(*ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	case AttributeTypeTimeStamp:
		os.(*ProtocolDataOutputStream).WriteBoolean(era)
		os.(*ProtocolDataOutputStream).WriteShort(yr)
		os.(*ProtocolDataOutputStream).WriteByte(int(mth))
		os.(*ProtocolDataOutputStream).WriteByte(day)
		os.(*ProtocolDataOutputStream).WriteByte(hr)
		os.(*ProtocolDataOutputStream).WriteByte(min)
		os.(*ProtocolDataOutputStream).WriteByte(sec)
		os.(*ProtocolDataOutputStream).WriteShort(msec)
		os.(*ProtocolDataOutputStream).WriteByte(TGNoZone)
		break
	default:
		errMsg := fmt.Sprintf("Bad Descriptor: %s", string(obj.AttrDesc.GetAttrType()))
		return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering TimestampAttribute::WriteValue"))
	}
	return nil
}

func (obj *TimestampAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TimestampAttribute:{")
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *TimestampAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *TimestampAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TimestampAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TimestampAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TimestampAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


var UniqueId int64

type BlobAttribute struct {
	*AbstractAttribute
	entityId int64
	isCached bool
}

// Create NewTGDecimal Attribute Instance
func DefaultBlobAttribute() *BlobAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BlobAttribute{})

	newAttribute := BlobAttribute{
		AbstractAttribute: defaultNewAbstractAttribute(),
		isCached:          false,
	}
	newAttribute.entityId = atomic.AddInt64(&UniqueId, 1)
	newAttribute.AttrValue = []byte{}
	return &newAttribute
}

func NewBlobAttributeWithOwner(ownerEntity tgdb.TGEntity) *BlobAttribute {
	newAttribute := DefaultBlobAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewBlobAttribute(attrDesc *AttributeDescriptor) *BlobAttribute {
	newAttribute := DefaultBlobAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewBlobAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *BlobAttribute {
	newAttribute := NewBlobAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for BlobAttribute
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) SetBlob(b []byte) {
	obj.AttrValue = b
	obj.setIsModified(true)
}

func (obj *BlobAttribute) GetAsBytes() []byte {
	if obj.entityId < 0 || obj.isCached {
		return obj.AttrValue.([]byte)
	}
	conn := obj.GetOwner().GetGraphMetadata().GetConnection()
	if obj.GetAttributeDescriptor().Is_Encrypted() {
		v, err := conn.DecryptEntity(obj.entityId)
		if err != nil {
			obj.AttrValue = nil
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("BlobAttribute::GetAsBytes - Unable to conn.DecryptEntity()"))
			}
			return nil
		}
		obj.AttrValue = v
	} else {
		v, err := conn.GetLargeObjectAsBytes(obj.entityId, false)
		if err != nil {
			obj.AttrValue = nil
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("BlobAttribute::GetAsBytes - Unable to conn.GetLargeObjectAsBytes()"))
			}
			return nil
		}
		obj.AttrValue = v
	}
	obj.isCached = true
	return obj.AttrValue.([]byte)
}

func (obj *BlobAttribute) GetAsByteBuffer() *bytes.Buffer {
	buf := obj.GetAsBytes()
	return bytes.NewBuffer(buf)
}

func (obj *BlobAttribute) GetEntityId() int64 {
	return obj.entityId
}

func (obj *BlobAttribute) GetIsCached() bool {
	return obj.isCached
}

func (obj *BlobAttribute) SetEntityId(eId int64) {
	obj.entityId = eId
}

func (obj *BlobAttribute) SetIsCached(flag bool) {
	obj.isCached = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *BlobAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *BlobAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *BlobAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *BlobAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *BlobAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *BlobAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *BlobAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *BlobAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *BlobAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() {
		return nil
	}

	if 	reflect.TypeOf(value).Kind() != reflect.Float32 &&
		reflect.TypeOf(value).Kind() != reflect.Float64 &&
		reflect.TypeOf(value).Kind() != reflect.Array &&
		reflect.TypeOf(value).Kind() != reflect.Struct &&
		reflect.TypeOf(value).Kind() != reflect.String {
		logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - attribute value is NOT in expected format/type"))
		errMsg := fmt.Sprintf("Failure to cast the attribute value to BlobAttribute")
		return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if reflect.TypeOf(value).Kind() == reflect.Float32 {
		v, err := FloatToByteArray(value.(float32))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - unable to extract attribute value in float format/type"))
			errMsg := fmt.Sprintf("Failure to covert float to BlobAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBlob(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		v, err := DoubleToByteArray(value.(float64))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - unable to extract attribute value in double format/type"))
			errMsg := fmt.Sprintf("Failure to covert double to BlobAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBlob(v)
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		v := []byte(value.(string))
		obj.SetBlob(v)
	} else if reflect.TypeOf(value).Kind() == reflect.Struct {
		v, err := InputStreamToByteArray(NewProtocolDataInputStream(value.([]byte)))
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:SetValue - unable to InputStreamToByteArray(iostream.NewProtocolDataInputStream(value.([]byte)))"))
			errMsg := fmt.Sprintf("Failure to covert instream bytes to BlobAttribute")
			return GetErrorByType(TGErrorTypeCoercionNotSupported, INTERNAL_SERVER_ERROR, errMsg, err.Error())
		}
		obj.SetBlob(v)
	} else {
		obj.AttrValue = value
		obj.setIsModified(true)
	}

	return nil
}

// ReadValue reads the value from input stream
func (obj *BlobAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	entityId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning BlobAttribute:ReadValue w/ Error in reading entityId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("BlobAttribute::ReadValue - read entityId: '%+v'", entityId))
	}
	obj.entityId = entityId
	obj.isCached = false
	return nil
}

// WriteValue writes the value to output stream
func (obj *BlobAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	os.(*ProtocolDataOutputStream).WriteLong(obj.entityId)
	if obj.AttrValue == nil {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
		if obj.GetAttributeDescriptor().Is_Encrypted() {
			err := AbstractAttributeWriteEncrypted(obj, os)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:WriteValue - Unable to AbstractAttributeWriteEncrypted() w/ Error: '%s'", err.Error()))
				errMsg := "BlobAttribute::WriteValue - Unable to AbstractAttributeWriteEncrypted()"
				return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
			}
		} else {
			err := os.(*ProtocolDataOutputStream).WriteBytes(obj.AttrValue.([]byte))
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:WriteValue - Unable to WriteBytes() w/ Error: '%s'", err.Error()))
				errMsg := "BlobAttribute::WriteValue - Unable to WriteBytes()"
				return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
			}
		}
	}
	return nil
}

func (obj *BlobAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BlobAttribute:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %+v", obj.entityId))
	buffer.WriteString(fmt.Sprintf(", IsCached: %+v", obj.isCached))
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *BlobAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *BlobAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified, obj.entityId, obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *BlobAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified, &obj.entityId, &obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning BlobAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ClobAttribute struct {
	*BlobAttribute
}

// Create New Attribute Instance
func DefaultClobAttribute() *ClobAttribute {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ClobAttribute{})

	newAttribute := ClobAttribute{
		BlobAttribute: DefaultBlobAttribute(),
	}
	return &newAttribute
}

func NewClobAttributeWithOwner(ownerEntity tgdb.TGEntity) *ClobAttribute {
	newAttribute := DefaultClobAttribute()
	newAttribute.owner = ownerEntity
	return newAttribute
}

func NewClobAttribute(attrDesc *AttributeDescriptor) *ClobAttribute {
	newAttribute := DefaultClobAttribute()
	newAttribute.AttrDesc = attrDesc
	return newAttribute
}

func NewClobAttributeWithDesc(ownerEntity tgdb.TGEntity, attrDesc *AttributeDescriptor, value interface{}) *ClobAttribute {
	newAttribute := NewClobAttributeWithOwner(ownerEntity)
	newAttribute.AttrDesc = attrDesc
	newAttribute.AttrValue = value
	return newAttribute
}

/////////////////////////////////////////////////////////////////
// Helper functions for ClobAttribute
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) getValueAsBytes() ([]byte, tgdb.TGError) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(obj.AttrValue)
	if err != nil {
		errMsg := "ClobAttribute::getValueAsBytes - Unable to encode attribute value"
		return nil, GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	dec := gob.NewDecoder(&network)
	var v []byte
	err = dec.Decode(&v)
	if err != nil {
		errMsg := "ClobAttribute::getValueAsBytes - Unable to decode attribute value"
		return nil, GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.Error())
	}
	return v, nil
}

func (obj *ClobAttribute) SetCharBuffer(b string) {
	if !obj.IsNull() && obj.AttrValue == b {
		return
	}
	obj.AttrValue = []byte(b)
	obj.setIsModified(true)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGAttribute
/////////////////////////////////////////////////////////////////

// GetAttributeDescriptor returns the AttributeDescriptor for this attribute
func (obj *ClobAttribute) GetAttributeDescriptor() tgdb.TGAttributeDescriptor {
	return obj.getAttributeDescriptor()
}

// GetIsModified checks whether the attribute modified or not
func (obj *ClobAttribute) GetIsModified() bool {
	return obj.getIsModified()
}

// GetName gets the Name for this attribute as the most generic form
func (obj *ClobAttribute) GetName() string {
	return obj.getName()
}

// GetOwner gets owner Entity of this attribute
func (obj *ClobAttribute) GetOwner() tgdb.TGEntity {
	return obj.getOwner()
}

// GetValue gets the value for this attribute as the most generic form
func (obj *ClobAttribute) GetValue() interface{} {
	return obj.getValue()
}

// IsNull checks whether the attribute value is null or not
func (obj *ClobAttribute) IsNull() bool {
	return obj.isNull()
}

// ResetIsModified resets the IsModified flag - recursively, if needed
func (obj *ClobAttribute) ResetIsModified() {
	obj.resetIsModified()
}

// SetOwner sets the owner entity - Need this indirection to traverse the chain
func (obj *ClobAttribute) SetOwner(ownerEntity tgdb.TGEntity) {
	obj.setOwner(ownerEntity)
}

// SetValue sets the value for this attribute. Appropriate data conversion to its attribute desc will be performed
// If the object is Null, then the object is explicitly set, but no value is provided.
func (obj *ClobAttribute) SetValue(value interface{}) tgdb.TGError {
	if value == nil {
		//errMsg := fmt.Sprintf("Attribute value is required")
		//return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		obj.AttrValue = value
		obj.setIsModified(true)
		return nil
	}
	if !obj.IsNull() && obj.AttrValue == value {
		return nil
	}
	// TODO: Revisit later
	//if (value == null)
	//{
	//	this.value = value;
	//	setModified();
	//	return;
	//}
	//else if (value instanceof char[]) {
	//	setCharBuffer(CharBuffer.wrap((char[])value));
	//}
	//else if (value instanceof CharBuffer) {
	//	setCharBuffer((CharBuffer) value);
	//}
	//else if (value instanceof CharSequence) {
	//	setCharBuffer(CharBuffer.wrap((CharSequence)value));
	//}
	//else {
	//	super.setValue(value);
	//}

	obj.SetCharBuffer(value.(string))
	return nil
}

// ReadValue reads the value from input stream
func (obj *ClobAttribute) ReadValue(is tgdb.TGInputStream) tgdb.TGError {
	entityId, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning ClobAttribute:ReadValue w/ Error in reading entityId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning BlobAttribute::ReadValue - read entityId: '%+v'", entityId))
	}
	obj.entityId = entityId
	obj.isCached = false
	return nil
}

// WriteValue writes the value to output stream
func (obj *ClobAttribute) WriteValue(os tgdb.TGOutputStream) tgdb.TGError {
	if obj.AttrValue == nil {
		os.(*ProtocolDataOutputStream).WriteBoolean(false)
	} else {
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
		v, err := obj.getValueAsBytes()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:WriteValue - Unable to decode attribute value w/ Error: '%s'", err.Error()))
			errMsg := "ClobAttribute::WriteValue - Unable to decode attribute value"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
		err = os.(*ProtocolDataOutputStream).WriteBytes(v)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:WriteValue - Unable to write attribute value w/ Error: '%s'", err.Error()))
			errMsg := "ClobAttribute::WriteValue - Unable to write attribute value"
			return GetErrorByType(TGErrorIOException, "TGErrorIOException", errMsg, err.GetErrorDetails())
		}
	}
	return nil
}

func (obj *ClobAttribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ClobAttribute:{")
	buffer.WriteString(fmt.Sprintf("EntityId: %+v", obj.entityId))
	buffer.WriteString(fmt.Sprintf(", IsCached: %+v", obj.isCached))
	strArray := []string{buffer.String(), obj.attributeToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

//@Override
//public char[] getAsChars() throws TGException {
//	return getAsChars("UTF-8");
//}
//
//@Override
//public CharBuffer getAsCharBuffer() throws TGException {
//	return getAsCharBuffer("UTF-8");
//}
//
//public char[] getAsChars(String encoding) throws TGException
//{
//	CharBuffer cb = getAsCharBuffer(encoding);
//	return cb.array();
//}
//
//public CharBuffer getAsCharBuffer(String encoding) throws TGException {
//	ByteBuffer bb = getAsByteBuffer();
//	Charset cs = Charset.forName(encoding);
//	CharBuffer cb = cs.decode(bb);
//	return cb;
//}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> types.TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *ClobAttribute) ReadExternal(is tgdb.TGInputStream) tgdb.TGError {
	return AbstractAttributeReadExternal(obj, is)
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *ClobAttribute) WriteExternal(os tgdb.TGOutputStream) tgdb.TGError {
	return AbstractAttributeWriteExternal(obj, os)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.owner, obj.AttrDesc, obj.AttrValue, obj.IsModified, obj.entityId, obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ClobAttribute) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.owner, &obj.AttrDesc, &obj.AttrValue, &obj.IsModified, &obj.entityId, &obj.isCached)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ClobAttribute:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
