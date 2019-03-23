package query

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGResultSet_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestGraphMetadata() *model.GraphMetadata {
	newGraphMetadata := model.NewGraphMetadata(nil)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	hAttrDesc := CreateTestAttributeDescriptor("ShortDesc", types.AttributeTypeShort)
	lAttrDesc := CreateTestAttributeDescriptor("LongDesc", types.AttributeTypeLong)
	nAttrDesc := CreateTestAttributeDescriptor("NumberDesc", types.AttributeTypeNumber)
	tAttrDesc := CreateTestAttributeDescriptor("TimeDesc", types.AttributeTypeTimeStamp)
	descMap := make(map[string]types.TGAttributeDescriptor, 7)
	descMap["BoolDesc"] = bAttrDesc
	descMap["IntegerDesc"] = iAttrDesc
	descMap["StringDesc"] = sAttrDesc
	descMap["ShortDesc"] = hAttrDesc
	descMap["LongDesc"] = lAttrDesc
	descMap["NumberDesc"] = nAttrDesc
	descMap["TimeDesc"] = tAttrDesc
	newGraphMetadata.SetAttributeDescriptors(descMap)
	descById := make(map[int64]types.TGAttributeDescriptor, 0)
	for _, attrDesc := range descMap {
		descById[attrDesc.GetAttributeId()] = attrDesc
	}
	newGraphMetadata.SetAttributeDescriptorsById(descById)

	nodeType1 := CreateTestNodeType("Node-1", types.SystemTypeNode, model.DefaultNodeType())
	nodeType2 := CreateTestNodeType("ChildNode-1", types.SystemTypeNode, nodeType1)
	nodeType3 := CreateTestNodeType("ChildNode-2", types.SystemTypeNode, nodeType1)
	nodeMap := make(map[string]types.TGNodeType, 3)
	nodeMap["Node-1"] = nodeType1
	nodeMap["ChildNode-1"] = nodeType2
	nodeMap["ChildNode-2"] = nodeType3
	newGraphMetadata.SetNodeTypes(nodeMap)
	nodeById := make(map[int]types.TGNodeType, 0)
	for _, nodeType := range nodeMap {
		nodeById[nodeType.GetEntityTypeId()] = nodeType
	}
	newGraphMetadata.SetNodeTypesById(nodeById)

	edgeType1 := CreateTestEdgeType("Edge-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, model.DefaultEntityType())
	edgeType2 := CreateTestEdgeType("Edge-2", types.DirectionTypeDirected, types.SystemTypeEdge, edgeType1)
	edgeType3 := CreateTestEdgeType("Edge-3", types.DirectionTypeDirected, types.SystemTypeEdge, edgeType1)
	edgeMap := make(map[string]types.TGEdgeType, 3)
	edgeMap["Edge-1"] = edgeType1
	edgeMap["Edge-2"] = edgeType2
	edgeMap["Edge-3"] = edgeType3
	newGraphMetadata.SetEdgeTypes(edgeMap)
	edgeById := make(map[int]types.TGEdgeType, 0)
	for _, edgeType := range edgeMap {
		edgeById[edgeType.GetEntityTypeId()] = edgeType
	}
	newGraphMetadata.SetEdgeTypesById(edgeById)
	return newGraphMetadata
}

func CreateTestAttributeDescriptor(name string, attributeType int) *model.AttributeDescriptor {
	newAttributeDescriptor := model.NewAttributeDescriptorWithType(name, attributeType)
	return newAttributeDescriptor
}

func CreateTestEdgeType(name string, directionType types.TGDirectionType, entityType types.TGSystemType, parent types.TGEntityType) *model.EdgeType {
	newEdgeType := model.DefaultEdgeType()
	newEdgeType.SetName(name)
	newEdgeType.SetSystemType(entityType)
	attributes := make(map[string]*model.AttributeDescriptor, 3)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	attributes["BoolDesc"] = bAttrDesc
	attributes["IntegerDesc"] = iAttrDesc
	attributes["StringDesc"] = sAttrDesc
	newEdgeType.SetAttributeMap(attributes)
	newEdgeType.SetDirectionType(directionType)
	newEdgeType.SetParent(parent)
	return newEdgeType
}

func CreateTestNodeType(name string, entityType types.TGSystemType, parent types.TGEntityType) *model.NodeType {
	newNodeType := model.DefaultNodeType()
	newNodeType.SetName(name)
	newNodeType.SetSystemType(entityType)
	attributes := make(map[string]*model.AttributeDescriptor, 3)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	attributes["BoolDesc"] = bAttrDesc
	attributes["IntegerDesc"] = iAttrDesc
	attributes["StringDesc"] = sAttrDesc
	newNodeType.SetAttributeMap(attributes)
	bPkAttrDesc := CreateTestAttributeDescriptor("BoolPkDesc", types.AttributeTypeBoolean)
	iPkAttrDesc := CreateTestAttributeDescriptor("IntegerPkDesc", types.AttributeTypeInteger)
	sPkAttrDesc := CreateTestAttributeDescriptor("StringPkDesc", types.AttributeTypeString)
	pKeys := []*model.AttributeDescriptor{bPkAttrDesc, iPkAttrDesc, sPkAttrDesc}
	newNodeType.SetPKeyAttributeDescriptors(pKeys)
	newNodeType.SetParent(parent)
	return newNodeType
}

func CreateTestNodeEntity() *model.Node {
	gmd := CreateTestGraphMetadata()
	newNodeEntity := model.NewNode(gmd)
	bAttrDesc := CreateTestAttributeDescriptor("Bool", types.AttributeTypeBoolean)
	boolAttr, _ := model.CreateAttributeWithDesc(newNodeEntity, bAttrDesc, true)
	iAttrDesc := CreateTestAttributeDescriptor("Integer", types.AttributeTypeInteger)
	intAttr, _ := model.CreateAttributeWithDesc(newNodeEntity, iAttrDesc, 33333)
	sAttrDesc := CreateTestAttributeDescriptor("String", types.AttributeTypeString)
	strAttr, _ := model.CreateAttributeWithDesc(newNodeEntity, sAttrDesc, "InsideNodeEntity")
	attrMap := make(map[string]types.TGAttribute, 0)
	attrMap["Bool"] = boolAttr
	attrMap["Integer"] = intAttr
	attrMap["String"] = strAttr
	newNodeEntity.SetAttributes(attrMap)
	return newNodeEntity
}

// TODO: Revisit later - once connection is implemented
// This will test both APIs - (a) AddEntityToResultSet, and (b) Count
func TestResultSetAddEntityToResultSet(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
}

func TestResultSetClose(t *testing.T) {
	testRs := DefaultResultSet()
	t.Logf("Before closing query, this query object is '%+v'", testRs)
	rs := testRs.Close()
	t.Logf("After closing query, this query object is '%+v'", rs)
}

func TestResultSetFirst(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
	rs1 := rs.AddEntityToResultSet(testNode)
	if rs1 != nil {
		t.Logf("This result '%+v' set has %d entries", rs1, rs1.Count())
	}
	rs2 := rs1.AddEntityToResultSet(testNode)
	if rs2 != nil {
		t.Logf("This result '%+v' set has %d entries", rs2, rs2.Count())
	}
	resultEntity := rs2.First()
	t.Logf("The first entry in the result set is '%+v'", resultEntity)
}

func TestResultSetLast(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
	rs1 := rs.AddEntityToResultSet(testNode)
	if rs1 != nil {
		t.Logf("This result '%+v' set has %d entries", rs1, rs1.Count())
	}
	rs2 := rs1.AddEntityToResultSet(testNode)
	if rs2 != nil {
		t.Logf("This result '%+v' set has %d entries", rs2, rs2.Count())
	}
	resultEntity := rs2.Last()
	t.Logf("The last entry in the result set is '%+v'", resultEntity)
}

func TestResultSetGetAt(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
	rs1 := rs.AddEntityToResultSet(testNode)
	if rs1 != nil {
		t.Logf("This result '%+v' set has %d entries", rs1, rs1.Count())
	}
	rs2 := rs1.AddEntityToResultSet(testNode)
	if rs2 != nil {
		t.Logf("This result '%+v' set has %d entries", rs2, rs2.Count())
	}
	resultEntity := rs2.GetAt(2)
	t.Logf("The entry at position # 2 in the result set is '%+v'", resultEntity)
}

func TestResultSetGetPosition(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
	rs1 := rs.AddEntityToResultSet(testNode)
	if rs1 != nil {
		t.Logf("This result '%+v' set has %d entries", rs1, rs1.Count())
	}
	rs2 := rs1.AddEntityToResultSet(testNode)
	if rs2 != nil {
		t.Logf("This result '%+v' set has %d entries", rs2, rs2.Count())
	}
	currentPos := rs2.GetPosition()
	t.Logf("The new entry in the result set will be at position '%+v'", currentPos)
}

// This will test all 3 APIs - (a) Prev, (b) Next, and (c) Skip
func TestResultSetNext(t *testing.T) {
	testRs := DefaultResultSet()
	testNode := CreateTestNodeEntity()
	t.Logf("Before adding entity '%+v' to result set, this result set has %d entries", testNode, testRs.Count())
	rs := testRs.AddEntityToResultSet(testNode)
	if rs != nil {
		t.Logf("This result '%+v' set has %d entries", rs, rs.Count())
	}
	rs1 := rs.AddEntityToResultSet(testNode)
	if rs1 != nil {
		t.Logf("This result '%+v' set has %d entries", rs1, rs1.Count())
	}
	rs2 := rs1.AddEntityToResultSet(testNode)
	if rs2 != nil {
		t.Logf("This result '%+v' set has %d entries", rs2, rs2.Count())
	}
	prevEntry := rs2.Prev()
	t.Logf("The prev entry in the result set was '%+v' and current position is '%+v'", prevEntry, rs2.GetPosition())
	prevEntry2 := rs2.Prev()
	t.Logf("The prev entry in the result set was '%+v' and current position is '%+v'", prevEntry2, rs2.GetPosition())
	rs3 := rs2.Skip(1)
	t.Logf("The current position is '%+v'", rs3.GetPosition())
	nextEntry := rs3.Next()
	t.Logf("The next entry in the result set is '%+v' and current position is '%+v'", nextEntry, rs3.GetPosition())
}
