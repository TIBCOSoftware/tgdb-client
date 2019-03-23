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
 * File name: TGEdge_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestEdgeEntity() *Edge {
	gmd := CreateTestGraphMetadata()
	newEdgeEntity := NewEdge(gmd)
	newEdgeEntity.virtualId = atomic.AddInt64(&EntitySequencer, 1)
	bAttrDesc := CreateTestAttributeDescriptor("Bool", types.AttributeTypeBoolean)
	boolAttr, _ := CreateAttributeWithDesc(newEdgeEntity, bAttrDesc, true)
	iAttrDesc := CreateTestAttributeDescriptor("Integer", types.AttributeTypeInteger)
	intAttr, _ := CreateAttributeWithDesc(newEdgeEntity, iAttrDesc, 22222)
	sAttrDesc := CreateTestAttributeDescriptor("String", types.AttributeTypeString)
	strAttr, _ := CreateAttributeWithDesc(newEdgeEntity, sAttrDesc, "InsideEdeEntity")
	attrMap := make(map[string]types.TGAttribute, 0)
	attrMap["Bool"] = boolAttr
	attrMap["Integer"] = intAttr
	attrMap["String"] = strAttr
	newEdgeEntity.attributes = attrMap
	fromNode := CreateTestNodeEntity()
	newEdgeEntity.fromNode = fromNode
	toNode := CreateTestNodeEntity()
	newEdgeEntity.toNode = toNode
	return newEdgeEntity
}

func TestEdgeEntityGetDirectionType(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	directionType := testEntity.GetDirectionType()
	t.Logf("TestEdgeEntity '%+v' has the following direction type '%+v'", testEntity, directionType.String())
}

func TestEdgeEntityGetVertices(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	vertices := testEntity.GetVertices()
	t.Logf("TestEdgeEntity '%+v' has the following vertices '%+v'", testEntity, vertices)
}

func TestEdgeEntityGetAttribute(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	bAttr := testEntity.GetAttribute("Bool")
	t.Logf("TestEdgeEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, bAttr, bAttr.GetValue())
}

func TestEdgeEntityGetAttributes(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	attrList, err := testEntity.GetAttributes()
	if err != nil {
		errMsg := "Error retrieving attributes of TestEdgeEntity"
		t.Errorf("TestEdgeEntityGetAttributes returned error message %s", errMsg)
	}
	for _, attr := range attrList {
		t.Logf("TestEdgeEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	}
}

func TestEdgeEntitySetAttribute(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	iAttr := testEntity.GetAttribute("Integer")
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, iAttr.GetName(), iAttr.GetValue())
	aErr := iAttr.SetValue(12345)
	if aErr != nil {
		errMsg := "Error modifying attribute value of TestEdgeEntity"
		t.Errorf("TestEdgeEntitySetAttribute returned error message %s", errMsg)
	}
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	err := testEntity.SetAttribute(iAttr)
	if err != nil {
		errMsg := "Error setting attributes of TestEdgeEntity"
		t.Errorf("TestEdgeEntitySetAttribute returned error message %s", errMsg)
	}
	t.Logf("Modified TestEdgeEntity '%+v'", testEntity)
}

func TestEdgeEntitySetOrCreateAttribute(t *testing.T) {
	testEntity := CreateTestEdgeEntity()
	err := testEntity.SetOrCreateAttribute("NumberDesc", 123.456)
	if err != nil {
		errMsg := "Error setting attributes of TestEdgeEntity"
		t.Errorf("TestEdgeEntitySetOrCreateAttribute returned error message %s", errMsg)
	}
	t.Logf("TestEdgeEntity '%+v' has set new attribute", testEntity)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestEdgeEntityWriteExternal(t *testing.T) {
	ToBeExportedEntityType := CreateTestEdgeEntity()
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	ToBeExportedEntityType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedEntityType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	ToBeExportedEntityType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedEntityType)
}
