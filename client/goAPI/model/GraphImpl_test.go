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
 * File name: TGGraph_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestGraphEntity() *Graph {
	gmd := CreateTestGraphMetadata()
	newGraphEntity := NewGraph(gmd)
	newGraphEntity.virtualId = atomic.AddInt64(&EntitySequencer, 1)
	bAttrDesc := CreateTestAttributeDescriptor("Bool", types.AttributeTypeBoolean)
	boolAttr, _ := CreateAttributeWithDesc(newGraphEntity, bAttrDesc, true)
	iAttrDesc := CreateTestAttributeDescriptor("Integer", types.AttributeTypeInteger)
	intAttr, _ := CreateAttributeWithDesc(newGraphEntity, iAttrDesc, 33333)
	sAttrDesc := CreateTestAttributeDescriptor("String", types.AttributeTypeString)
	strAttr, _ := CreateAttributeWithDesc(newGraphEntity, sAttrDesc, "InsideGraphEntity")
	attrMap := make(map[string]types.TGAttribute, 0)
	attrMap["Bool"] = boolAttr
	attrMap["Integer"] = intAttr
	attrMap["String"] = strAttr
	newGraphEntity.attributes = attrMap
	return newGraphEntity
}

// This will test 3 APIs - (a) AddEdge, (b) AddEdgeWithDirectionType, and (c) GetEdges
func TestGraphEntityGetEdges(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	edge1 := CreateTestEdgeEntity()
	edge1.directionType = types.DirectionTypeDirected
	testEntity.AddEdge(edge1)
	toNode := CreateTestNodeEntity()
	edge2 := CreateTestEdgeEntity()
	edge2.directionType = types.DirectionTypeUnDirected
	testEntity.AddEdgeWithDirectionType(toNode, DefaultEdgeType(), types.DirectionTypeUnDirected)
	edge3 := CreateTestEdgeEntity()
	edge3.directionType = types.DirectionTypeBiDirectional
	testEntity.AddEdge(edge3)
	edges := testEntity.GetEdges()
	if len(edges) == 0 {
		errMsg := "Error retrieving edges of TestGraphEntity"
		t.Errorf("TestGraphEntityGetEdges returned error message %s", errMsg)
	}
	t.Logf("TestGraphEntity has the following edges '%+v'", edges)
}

func TestGraphEntityGetEdgesForDirectionType(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	edge1 := CreateTestEdgeEntity()
	edge1.directionType = types.DirectionTypeDirected
	testEntity.AddEdge(edge1)
	toNode := CreateTestNodeEntity()
	edge2 := CreateTestEdgeEntity()
	edge2.directionType = types.DirectionTypeUnDirected
	testEntity.AddEdgeWithDirectionType(toNode, DefaultEdgeType(), types.DirectionTypeUnDirected)
	edge3 := CreateTestEdgeEntity()
	edge3.directionType = types.DirectionTypeBiDirectional
	testEntity.AddEdge(edge3)
	edges := testEntity.GetEdgesForDirectionType(types.DirectionTypeBiDirectional)
	if len(edges) == 0 {
		errMsg := "Error retrieving edges of TestGraphEntity"
		t.Errorf("TestGraphEntityGetEdgesForDirectionType returned error message %s", errMsg)
	}
	t.Logf("TestGraphEntity has the following edges '%+v'", edges)
}

func TestGraphEntityGetEdgesForEdgeType(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	edge1 := CreateTestEdgeEntity()
	edge1.directionType = types.DirectionTypeDirected
	testEntity.AddEdge(edge1)
	toNode := CreateTestNodeEntity()
	edge2 := CreateTestEdgeEntity()
	edge2.directionType = types.DirectionTypeUnDirected
	testEntity.AddEdgeWithDirectionType(toNode, DefaultEdgeType(), types.DirectionTypeUnDirected)
	edge3 := CreateTestEdgeEntity()
	edge3.directionType = types.DirectionTypeBiDirectional
	testEntity.AddEdge(edge3)
	edges := testEntity.GetEdgesForEdgeType(DefaultEdgeType(), types.DirectionAny)
	if len(edges) == 0 {
		errMsg := "Error retrieving edges of TestGraphEntity"
		t.Errorf("TestGraphEntityGetEdgesForDirectionType returned error message %s", errMsg)
	}
	t.Logf("TestGraphEntity has the following edges '%+v'", edges)
}

func TestGraphEntityGetAttribute(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	bAttr := testEntity.GetAttribute("Bool")
	t.Logf("TestGraphEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, bAttr, bAttr.GetValue())
}

func TestGraphEntityGetAttributes(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	attrList, err := testEntity.GetAttributes()
	if err != nil {
		errMsg := "Error retrieving attributes of TestGraphEntity"
		t.Errorf("TestGraphEntityGetAttributes returned error message %s", errMsg)
	}
	for _, attr := range attrList {
		t.Logf("TestGraphEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	}
}

func TestGraphEntitySetAttribute(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	iAttr := testEntity.GetAttribute("Integer")
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, iAttr.GetName(), iAttr.GetValue())
	aErr := iAttr.SetValue(12345)
	if aErr != nil {
		errMsg := "Error modifying attribute value of TestGraphEntity"
		t.Errorf("TestGraphEntitySetAttribute returned error message %s", errMsg)
	}
	//t.Logf("TestEntity '%+v' extracted attribute '%+v' that has value as '%+v'", testEntity, attr.GetName(), attr.GetValue())
	err := testEntity.SetAttribute(iAttr)
	if err != nil {
		errMsg := "Error setting attributes of TestGraphEntity"
		t.Errorf("TestGraphEntitySetAttribute returned error message %s", errMsg)
	}
	t.Logf("Modified TestGraphEntity '%+v'", testEntity)
}

func TestGraphEntitySetOrCreateAttribute(t *testing.T) {
	testEntity := CreateTestGraphEntity()
	err := testEntity.SetOrCreateAttribute("NumberDesc", 123.456)
	if err != nil {
		errMsg := "Error setting attributes of TestGraphEntity"
		t.Errorf("TestGraphEntitySetOrCreateAttribute returned error message %s", errMsg)
	}
	t.Logf("TestGraphEntity '%+v' has set new attribute", testEntity)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestGraphEntityWriteExternal(t *testing.T) {
	ToBeExportedEntityType := CreateTestGraphEntity()
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = ToBeExportedEntityType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedEntityType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	_ = ToBeExportedEntityType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedEntityType)
}
