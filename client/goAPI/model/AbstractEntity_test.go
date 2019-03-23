package model

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sync/atomic"
	"testing"
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
 * File name: AbstractEntity_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestAttributeDescriptor(name string, attributeType int) *AttributeDescriptor {
	newAttributeDescriptor := NewAttributeDescriptorWithType(name, attributeType)
	return newAttributeDescriptor
}

func CreateTestAbstractEntity() *AbstractEntity {
	newAbstractEntity := AbstractEntity{
		isNew:              true,
		EntityKind:         0,
		version:            0,
		entityId:           -1,
		isDeleted:          false,
		isInitialized:      true,
		attributes:         make(map[string]types.TGAttribute, 0),
		modifiedAttributes: make([]types.TGAttribute, 0),
	}
	newAbstractEntity.virtualId = atomic.AddInt64(&EntitySequencer, 1)
	bAttrDesc := CreateTestAttributeDescriptor("Bool", types.AttributeTypeBoolean)
	boolAttr, _ := CreateAttributeWithDesc(&newAbstractEntity, bAttrDesc, true)
	iAttrDesc := CreateTestAttributeDescriptor("Integer", types.AttributeTypeInteger)
	intAttr, _ := CreateAttributeWithDesc(&newAbstractEntity, iAttrDesc, 98765)
	sAttrDesc := CreateTestAttributeDescriptor("String", types.AttributeTypeString)
	strAttr, _ := CreateAttributeWithDesc(&newAbstractEntity, sAttrDesc, "StringAttribute")
	attrMap := make(map[string]types.TGAttribute, 0)
	attrMap["Bool"] = boolAttr
	attrMap["Integer"] = intAttr
	attrMap["String"] = strAttr
	newAbstractEntity.attributes = attrMap
	gmd := CreateTestGraphMetadata()
	newAbstractEntity.graphMetadata = gmd
	return &newAbstractEntity
}

func TestAbstractEntityGetAttribute(t *testing.T) {
	testEntity := CreateTestAbstractEntity()
	bAttr := testEntity.GetAttribute("Bool")
	t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, bAttr.GetName(), bAttr.GetValue())
}

func TestAbstractEntityGetAttributes(t *testing.T) {
	testEntity := CreateTestAbstractEntity()
	attrList, err := testEntity.GetAttributes()
	if err != nil {
		errMsg := "Error retrieving attributes of TestAbstractEntity"
		t.Errorf("TestAbstractEntityGetAttributes returned error message %s", errMsg)
	}
	for _, attr := range attrList {
		t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	}
}

func TestAbstractEntitySetAttribute(t *testing.T) {
	testEntity := CreateTestAbstractEntity()
	iAttr := testEntity.GetAttribute("Integer")
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, iAttr.GetName(), iAttr.GetValue())
	aErr := iAttr.SetValue(12345)
	if aErr != nil {
		errMsg := "Error modifying attribute value of TestAbstractEntity"
		t.Errorf("TestAbstractEntitySetAttribute returned error message %s", errMsg)
	}
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	err := testEntity.SetAttribute(iAttr)
	if err != nil {
		errMsg := "Error setting attributes of TestAbstractEntity"
		t.Errorf("TestAbstractEntitySetAttribute returned error message %s", errMsg)
	}
	t.Logf("Modified TestEntity '%+v'", testEntity)
}

func TestAbstractEntitySetOrCreateAttribute(t *testing.T) {
	testEntity := CreateTestAbstractEntity()
	err := testEntity.SetOrCreateAttribute("NumberDesc", 123.456)
	if err != nil {
		errMsg := "Error setting attributes of TestAbstractEntity"
		t.Errorf("TestAbstractEntitySetOrCreateAttribute returned error message %s", errMsg)
	}
	t.Logf("TestEntity '%+v' has set new attribute", testEntity)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestAbstractEntityWriteExternal(t *testing.T) {
	ToBeExportedEntityType := CreateTestAbstractEntity()
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = ToBeExportedEntityType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedEntityType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	_ = ToBeExportedEntityType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedEntityType)
}
