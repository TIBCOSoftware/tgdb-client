package types

import "strings"

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
 * File name: TGAttributeType.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**	Test Program validated in GO Playground
import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)
func main() {
	fmt.Println("Hello, playground")

	tod := time.Now()
	fmt.Printf("\t%+v is of type %s\n", tod, reflect.TypeOf(tod).String())
	str1 := true
	fmt.Printf("\t%+v is of type %s\n", str1, reflect.TypeOf(str1).String())
	str2 := byte(1)
	fmt.Printf("\t%+v is of type %s\n", str2, reflect.TypeOf(str2).String())
	str3 := byte('A')
	fmt.Printf("\t%+v is of type %s\n", str3, reflect.TypeOf(str3).String())
	str4 := int16(21)
	fmt.Printf("\t%+v is of type %s\n", str4, reflect.TypeOf(str4).String())
	str5 := 123456
	fmt.Printf("\t%+v is of type %s\n", str5, reflect.TypeOf(str5).String())
	str6 := int64(123456789098765)
	fmt.Printf("\t%+v is of type %s\n", str6, reflect.TypeOf(str6).String())
	str7 := "tradeVolume"
	fmt.Printf("\t%+v is of type %s\n", str7, reflect.TypeOf(str7).String())
	date := "2014-11-17"
	fmt.Printf("\t%s is of type %s\n", date, reflect.TypeOf(date).String())
	opPrice := float32(42.75)
	fmt.Printf("\t%f is of type %s\n", opPrice, reflect.TypeOf(opPrice).String())
	clPrice := 42.990002
	fmt.Printf("\t%f is of type %s\n", clPrice, reflect.TypeOf(clPrice).String())
	hiPrice := 42.730000
	fmt.Printf("\t%f is of type %s\n", hiPrice, reflect.TypeOf(hiPrice).String())
	loPrice := 42.919998
	fmt.Printf("\t%f is of type %s\n", loPrice, reflect.TypeOf(loPrice).String())
	avgPrice := 38.971539
	fmt.Printf("\t%f is of type %s\n", avgPrice, reflect.TypeOf(avgPrice).String())
	volume := big.NewInt(106640001235)
	fmt.Printf("\t%d is of type %s\n", volume, reflect.TypeOf(*volume).String())

	fmt.Printf("\t%+v is mapped to %+v\n", tod, GetAttributeTypeFromName(reflect.TypeOf(tod).String()))
	fmt.Printf("\t%+v is mapped to %+v\n", date, GetAttributeTypeFromName(reflect.TypeOf(date).String()))
	fmt.Printf("\t%+v is mapped to %+v\n", opPrice, GetAttributeTypeFromName(reflect.TypeOf(opPrice).String()))
	fmt.Printf("\t%+v is mapped to %+v\n", volume, GetAttributeTypeFromName(reflect.TypeOf(*volume).String()))
}
----------- Output -----------
Hello, playground
	2009-11-10 23:00:00 +0000 UTC m=+0.000000001 is of type time.Time
	true is of type bool
	1 is of type uint8
	65 is of type uint8
	21 is of type int16
	123456 is of type int
	123456789098765 is of type int64
	tradevolume is of type string
	2014-11-17 is of type string
	42.750000 is of type float32
	42.990002 is of type float64
	42.730000 is of type float64
	42.919998 is of type float64
	38.971539 is of type float64
	106640001235 is of type big.Int
	2009-11-10 23:00:00 +0000 UTC m=+0.000000001 is mapped to &{TypeId:13 TypeName:time.Time Implementor:TimestampAttribute}
	2014-11-17 is mapped to &{TypeId:10 TypeName:string Implementor:StringAttribute}
	42.75 is mapped to &{TypeId:7 TypeName:float32 Implementor:FloatAttribute}
	+106640001235 is mapped to &{TypeId:9 TypeName:big.Int Implementor:NumberAttribute}
*/

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

// GetAttributeTypeFromName returns the TGAttributeType given its name
func GetAttributeTypeFromName(aName string) *AttributeType {
	for _, aType := range PreDefinedAttributeTypes {
		if strings.ToLower(aType.typeName) == strings.ToLower(aName) {
			return &aType
		}
	}
	invalid := PreDefinedAttributeTypes[AttributeTypeInvalid]
	return &invalid
}

