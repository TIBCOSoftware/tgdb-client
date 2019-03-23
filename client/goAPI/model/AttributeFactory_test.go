package model

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"math/rand"
	"strconv"
	"testing"
	"time"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGAttributeFactory_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func createTestAttributeDescriptor(aType int) *AttributeDescriptor {
	attrDescName := "AttributeDescriptorType-" + strconv.Itoa(aType)
	return NewAttributeDescriptorWithType(attrDescName, aType)
}

func createTestBooleanAttribute() *BooleanAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeBoolean)
	return NewBooleanAttribute(attrDesc)
}

func createTestByteAttribute() *ByteAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeByte)
	return NewByteAttribute(attrDesc)
}

func createTestCharAttribute() *CharAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeChar)
	return NewCharAttribute(attrDesc)
}

func createTestShortAttribute() *ShortAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeShort)
	return NewShortAttribute(attrDesc)
}

func createTestIntegerAttribute() *IntegerAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeInteger)
	return NewIntegerAttribute(attrDesc)
}

func createTestLongAttribute() *LongAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeLong)
	return NewLongAttribute(attrDesc)
}

func createTestFloatAttribute() *FloatAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeFloat)
	return NewFloatAttribute(attrDesc)
}

func createTestDoubleAttribute() *DoubleAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeDouble)
	return NewDoubleAttribute(attrDesc)
}

func createTestNumberAttribute() *NumberAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeNumber)
	attrDesc.SetPrecision(24)
	attrDesc.SetScale(8)
	return NewNumberAttribute(attrDesc)
}

func createTestStringAttribute() *StringAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeString)
	return NewStringAttribute(attrDesc)
}

func createTestDateAttribute() *TimestampAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeDate)
	return NewTimestampAttribute(attrDesc)
}

func createTestTimeAttribute() *TimestampAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeTime)
	return NewTimestampAttribute(attrDesc)
}

func createTestTimeStampAttribute() *TimestampAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeTimeStamp)
	return NewTimestampAttribute(attrDesc)
}

func createTestBlobAttribute() *BlobAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeBlob)
	return NewBlobAttribute(attrDesc)
}

func createTestClobAttribute() *ClobAttribute {
	attrDesc := createTestAttributeDescriptor(types.AttributeTypeClob)
	return NewClobAttribute(attrDesc)
}

func createTestAttributeByType(aType int) types.TGAttribute {
	switch aType {
	case types.AttributeTypeBoolean:
		return createTestBooleanAttribute()
	case types.AttributeTypeByte:
		return createTestByteAttribute()
	case types.AttributeTypeChar:
		return createTestCharAttribute()
	case types.AttributeTypeShort:
		return createTestShortAttribute()
	case types.AttributeTypeInteger:
		return createTestIntegerAttribute()
	case types.AttributeTypeLong:
		return createTestLongAttribute()
	case types.AttributeTypeFloat:
		return createTestFloatAttribute()
	case types.AttributeTypeDouble:
		return createTestDoubleAttribute()
	case types.AttributeTypeNumber:
		return createTestNumberAttribute()
	case types.AttributeTypeString:
		return createTestStringAttribute()
	case types.AttributeTypeDate:
		return createTestDateAttribute()
	case types.AttributeTypeTime:
		return createTestTimeAttribute()
	case types.AttributeTypeTimeStamp:
		return createTestTimeStampAttribute()
	case types.AttributeTypeBlob:
		return createTestBlobAttribute()
	case types.AttributeTypeClob:
		return createTestClobAttribute()
	default:
		return defaultNewAbstractAttribute()
	}
}

func TestCreateAttributeByType(t *testing.T) {
	for attrTypeId := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		newAttr := createTestAttributeByType(attrTypeId)
		t.Logf("AttributeFactory returned %s attribute for attrType: '%+v' as '%+v'", types.GetAttributeTypeFromId(attrTypeId).GetTypeName(), attrTypeId, newAttr)
	}
}

func TestCreateAttribute(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	for attrTypeId, attrType := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		attrDesc := createTestAttributeDescriptor(attrType.GetTypeId())
		if attrTypeId == 9 {
			attrDesc.SetPrecision(24)
			attrDesc.SetScale(8)
		}
		value := rand.Int63()
		newAttr, err := CreateAttributeWithDesc(nil, attrDesc, value)
		if err != nil {
			t.Errorf("AttributeFactory could not instantiate %s attribute for attrType: '%+v'", types.GetAttributeTypeFromId(attrTypeId).GetTypeName(), attrTypeId)
		}
		t.Logf("AttributeFactory returned %s attribute for attrType: '%+v' as '%+v'", types.GetAttributeTypeFromId(attrTypeId).GetTypeName(), attrTypeId, newAttr)
	}
}

func TestGetName(t *testing.T) {
	for attrTypeId := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		name := GetName(attrTypeId)
		t.Logf("AttributeFactory returned '%+v' attribute name for attrType: '%+v'", name, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
	}
}

func TestGetValue(t *testing.T) {
	for attrTypeId, attrType := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		attrDesc := createTestAttributeDescriptor(attrType.GetTypeId())
		if attrTypeId == 9 {
			attrDesc.SetPrecision(24)
			attrDesc.SetScale(8)
		}
		value1 := rand.Int63()
		newAttr, err := CreateAttributeWithDesc(nil, attrDesc, value1)
		if err != nil {
			t.Errorf("AttributeFactory could not instantiate %s attribute for attrType: '%+v'", types.GetAttributeTypeFromId(attrTypeId).GetTypeName(), attrTypeId)
		}
		value := newAttr.GetValue()
		t.Logf("AttributeFactory returned '%+v' attribute value for attrType: '%+v'", value, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
	}
}

func TestIsNull(t *testing.T) {
	for attrTypeId := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		value := IsNull(attrTypeId)
		t.Logf("AttributeFactory verified '%+v' attribute value whether it is null (TRUE) or not (FALSE): '%+v'", types.GetAttributeTypeFromId(attrTypeId).GetTypeName(), value)
	}
}

const (
	dateTimeFormat = time.RFC3339
	dateFormat     = "2006-01-02"
	timeFormat     = "15:04"
)

// This automatically will test both APIs - (a) GetValue and (b) SetValue
func TestSetValue(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	for attrTypeId := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}
		newAttr := createTestAttributeByType(attrTypeId)
		oldValue := newAttr.GetValue()
		t.Logf("AttributeFactory verified existing attribute value '%+v' for attrType: '%+v'", oldValue, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())

		var value interface{}
		switch attrTypeId {
		case types.AttributeTypeBoolean:
			value = true
		case types.AttributeTypeByte:
			value = byte('G')
		case types.AttributeTypeChar:
			value = 'A'
		case types.AttributeTypeShort:
			value = int16(rand.Int())
		case types.AttributeTypeInteger:
			value = rand.Int()
		case types.AttributeTypeLong:
			value = rand.Int63()
		case types.AttributeTypeFloat:
			value = rand.Float32()
		case types.AttributeTypeDouble:
			value = rand.Float64()
		case types.AttributeTypeNumber:
			value, _ = utils.NewTGDecimalFromString("123456.0987654")
		case types.AttributeTypeString:
			value = "Winner"
		case types.AttributeTypeDate:
			env := utils.NewTGEnvironment()
			value = time.Date(2019, time.January, 01, 0, 0, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
		case types.AttributeTypeTime:
			env := utils.NewTGEnvironment()
			value = time.Date(01, 01, 01, 23, 45, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
		case types.AttributeTypeTimeStamp:
			env := utils.NewTGEnvironment()
			value = time.Date(2019, time.January, 01, 23, 45, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
			//value, _ = time.Parse(dateTimeFormat, "2019-01-01 23:45:10")
		case types.AttributeTypeBlob:
			value = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		case types.AttributeTypeClob:
			value = `"/**
				 * Copyright 2018 TIBCO Software Inc. All rights reserved.
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
				 **/"`
		case types.AttributeTypeInvalid:
		default:
		}

		t.Logf("AttributeFactory trying to set the attribute value '%+v' for attrType: '%+v'", value, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
		err := newAttr.SetValue(value)
		if err != nil {
			t.Errorf("AttributeFactory could not set attribute value '%+v' for attrType: '%+v'", value, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
		}
		newValue := newAttr.GetValue()
		t.Logf("AttributeFactory verified modified attribute value '%+v' for attrType: '%+v'", newValue, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
		t.Log("")
	}
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestWriteExternal(t *testing.T) {
	for attrTypeId := range types.PreDefinedAttributeTypes {
		if attrTypeId == 0 {
			continue
		}

		var value interface{}
		switch attrTypeId {
		case types.AttributeTypeBoolean:
			value = true
		case types.AttributeTypeByte:
			value = byte('G')
		case types.AttributeTypeChar:
			value = 'A'
		case types.AttributeTypeShort:
			value = int16(rand.Int())
		case types.AttributeTypeInteger:
			value = rand.Int()
		case types.AttributeTypeLong:
			value = rand.Int63()
		case types.AttributeTypeFloat:
			value = rand.Float32()
		case types.AttributeTypeDouble:
			value = rand.Float64()
		case types.AttributeTypeNumber:
			value, _ = utils.NewTGDecimalFromString("123456.0987654")
		case types.AttributeTypeString:
			value = "Winner"
		case types.AttributeTypeDate:
			env := utils.NewTGEnvironment()
			value = time.Date(2019, time.January, 01, 0, 0, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
		case types.AttributeTypeTime:
			env := utils.NewTGEnvironment()
			//value, _ = time.Parse("03:04", "23:45")
			value = time.Date(01, 01, 01, 23, 45, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
		case types.AttributeTypeTimeStamp:
			env := utils.NewTGEnvironment()
			value = time.Date(2019, time.January, 01, 23, 45, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
			//value, _ = time.Parse(dateTimeFormat, "2019-01-01 23:45:10")
		case types.AttributeTypeBlob:
			value = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		case types.AttributeTypeClob:
			value = `"/**
				 * Copyright 2018 TIBCO Software Inc. All rights reserved.
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
				 **/"`
		case types.AttributeTypeInvalid:
		default:
		}

		newAttr := createTestAttributeByType(attrTypeId)
		err := newAttr.SetValue(value)
		if err != nil {
			t.Errorf("AttributeFactory could not set attribute value '%+v' for attrType: '%+v'", value, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())
		}
		t.Logf("AttributeFactory returned modified attribute '%+v' for attrType: '%+v'", newAttr, types.GetAttributeTypeFromId(attrTypeId).GetTypeName())

		//var network bytes.Buffer
		oNetwork := iostream.DefaultProtocolDataOutputStream()
		err = newAttr.WriteExternal(oNetwork)
		if err != nil {
			t.Errorf("AttributeFactory could not newAttr.WriteExternal for attr: '%+v'", newAttr)
		}
		t.Logf("AttributeFactory WriteExternal exported '%+v' attribute value", types.GetAttributeTypeFromId(attrTypeId).GetTypeName())

		//iNetwork := iostream.DefaultProtocolDataInputStream()
		//err = newAttr.ReadExternal(iNetwork)
		//if err != nil {
		//	t.Errorf("AttributeFactory could not newAttr.ReadExternal for attr: '%+v'", newAttr)
		//}
		//t.Logf("AttributeFactory ReadExternal imported '%+v' attribute value", types.GetAttributeTypeFromId(attrTypeId).TypeName)
	}
}
