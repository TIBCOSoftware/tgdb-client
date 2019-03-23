package model

import (
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
 * File name: TGGraphMetadata_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestGraphMetadata() *GraphMetadata {
	newGraphMetadata := NewGraphMetadata(nil)
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
	newGraphMetadata.descriptors = descMap
	for _, attrDesc := range descMap {
		newGraphMetadata.descriptorsById[attrDesc.GetAttributeId()] = attrDesc
	}

	nodeType1 := CreateTestNodeType("Node-1", types.SystemTypeNode, DefaultNodeType())
	nodeType2 := CreateTestNodeType("ChildNode-1", types.SystemTypeNode, nodeType1)
	nodeType3 := CreateTestNodeType("ChildNode-2", types.SystemTypeNode, nodeType1)
	nodeMap := make(map[string]types.TGNodeType, 3)
	nodeMap["Node-1"] = nodeType1
	nodeMap["ChildNode-1"] = nodeType2
	nodeMap["ChildNode-2"] = nodeType3
	newGraphMetadata.nodeTypes = nodeMap
	for _, nodeType := range nodeMap {
		newGraphMetadata.nodeTypesById[nodeType.GetEntityTypeId()] = nodeType
	}

	edgeType1 := CreateTestEdgeType("Edge-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, DefaultEntityType())
	edgeType2 := CreateTestEdgeType("Edge-2", types.DirectionTypeDirected, types.SystemTypeEdge, edgeType1)
	edgeType3 := CreateTestEdgeType("Edge-3", types.DirectionTypeDirected, types.SystemTypeEdge, edgeType1)
	edgeMap := make(map[string]types.TGEdgeType, 3)
	edgeMap["Edge-1"] = edgeType1
	edgeMap["Edge-2"] = edgeType2
	edgeMap["Edge-3"] = edgeType3
	newGraphMetadata.edgeTypes = edgeMap
	for _, edgeType := range edgeMap {
		newGraphMetadata.edgeTypesById[edgeType.GetEntityTypeId()] = edgeType
	}
	return newGraphMetadata
}

func TestMetadataCreateAttributeDescriptor(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	t.Logf("TestGraphMetadata has following %d descriptors", len(testGmd.descriptors))
	dAttrDesc := testGmd.CreateAttributeDescriptor("First", types.AttributeTypeString, false)
	t.Logf("TestGraphMetadata has added descriptor: '%+v'", dAttrDesc)
	descList, err := testGmd.GetAttributeDescriptors()
	if err != nil {
		errMsg := "Error retrieving descriptors of TestAbstractEntity"
		t.Errorf("TestMetadataCreateAttributeDescriptor returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d descriptors", len(descList))
	t.Logf("TestGraphMetadata has following descriptors: '%+v'", descList[len(descList)-1])
}

func TestMetadataGetAttributeDescriptor(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	desc, err := testGmd.GetAttributeDescriptor("TimeDesc")
	if err != nil {
		errMsg := "Error retrieving descriptor of TestMetadata"
		t.Errorf("TestMetadataGetAttributeDescriptor returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following descriptor: '%+v'", desc)
}

func TestMetadataGetAttributeDescriptorById(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	dAttrDesc := testGmd.CreateAttributeDescriptor("First", types.AttributeTypeString, false)
	t.Logf("TestGraphMetadata has added descriptor: '%+v'", dAttrDesc)
	desc, err := testGmd.GetAttributeDescriptorById(dAttrDesc.GetAttributeId())
	if err != nil {
		errMsg := "Error retrieving descriptor of TestMetadata"
		t.Errorf("TestMetadataGetAttributeDescriptorById returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following descriptor: '%+v'", desc)
}

func TestMetadataGetAttributeDescriptors(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	descList, err := testGmd.GetAttributeDescriptors()
	if err != nil {
		errMsg := "Error retrieving descriptors of TestMetadata"
		t.Errorf("TestMetadataGetAttributeDescriptors returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d descriptors", len(descList))
	for _, desc := range descList {
		t.Logf("TestGraphMetadata has following descriptor: '%+v'", desc)
	}
}

func TestMetadataCreateCompositeKey(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	cKey := testGmd.CreateCompositeKey("Node-1")
	t.Logf("TestGraphMetadata has retrieved the following composite key: '%+v'", cKey)
}

func TestMetadataCreateEdgeType(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	t.Logf("TestGraphMetadata has following %d edges", len(testGmd.edgeTypes))
	pEdge, err := testGmd.GetEdgeType("Edge-3")
	if err != nil {
		errMsg := "Error retrieving edge of TestMetadata"
		t.Errorf("TestMetadataCreateEdgeType returned error message %s", errMsg)
	}
	newEdgeType := testGmd.CreateEdgeType("First", pEdge)
	t.Logf("TestGraphMetadata has added edge type: '%+v'", newEdgeType)
	edgeList, err := testGmd.GetEdgeTypes()
	if err != nil {
		errMsg := "Error retrieving edges of TestMetadata"
		t.Errorf("TestMetadataCreateEdgeType returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d edges", len(edgeList))
	t.Logf("TestGraphMetadata has following edges: '%+v'", edgeList[len(edgeList)-1])
}

func TestMetadataGetEdgeType(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	edge, err := testGmd.GetEdgeType("Edge-1")
	if err != nil {
		errMsg := "Error retrieving edge of TestMetadata"
		t.Errorf("TestMetadataGetEdgeType returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following edge: '%+v'", edge)
}

func TestMetadataGetEdgeTypeById(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	pEdge, err := testGmd.GetEdgeType("Edge-3")
	if err != nil {
		errMsg := "Error retrieving edge of TestMetadata"
		t.Errorf("TestMetadataGetEdgeTypeById returned error message %s", errMsg)
	}
	newEdgeType := testGmd.CreateEdgeType("First", pEdge)
	t.Logf("TestGraphMetadata has added edge type: '%+v'", newEdgeType)
	edge, err := testGmd.GetEdgeTypeById(newEdgeType.GetEntityTypeId())
	if err != nil {
		errMsg := "Error retrieving edge of TestMetadata"
		t.Errorf("TestMetadataGetEdgeTypeById returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following edge: '%+v'", edge)
}

func TestMetadataGetEdgeTypes(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	edgeList, err := testGmd.GetEdgeTypes()
	if err != nil {
		errMsg := "Error retrieving edges of TestMetadata"
		t.Errorf("TestMetadataGetEdgeTypes returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d edges", len(edgeList))
	for _, edge := range edgeList {
		t.Logf("TestGraphMetadata has following edge: '%+v'", edge)
	}
}

func TestMetadataCreateNodeType(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	t.Logf("TestGraphMetadata has following %d nodes", len(testGmd.nodeTypes))
	pNode, err := testGmd.GetNodeType("Node-1")
	if err != nil {
		errMsg := "Error retrieving node of TestMetadata"
		t.Errorf("TestMetadataCreateNodeType returned error message %s", errMsg)
	}
	newNodeType := testGmd.CreateNodeType("ChildNode-3", pNode)
	t.Logf("TestGraphMetadata has added node type: '%+v'", newNodeType)
	nodeList, err := testGmd.GetNodeTypes()
	if err != nil {
		errMsg := "Error retrieving nodes of TestMetadata"
		t.Errorf("TestMetadataCreateNodeType returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d nodes", len(nodeList))
	t.Logf("TestGraphMetadata has following nodes: '%+v'", nodeList[len(nodeList)-1])
}

func TestMetadataGetNodeType(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	node, err := testGmd.GetNodeType("ChildNode-1")
	if err != nil {
		errMsg := "Error retrieving node of TestMetadata"
		t.Errorf("TestMetadataGetNodeType returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following node: '%+v'", node)
}

func TestMetadataGetNodeTypeById(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	pNode, err := testGmd.GetNodeType("Node-1")
	if err != nil {
		errMsg := "Error retrieving node of TestMetadata"
		t.Errorf("TestMetadataCreateNodeType returned error message %s", errMsg)
	}
	newNodeType := testGmd.CreateNodeType("ChildNode-3", pNode)
	t.Logf("TestGraphMetadata has added node type: '%+v'", newNodeType)
	node, err := testGmd.GetNodeTypeById(newNodeType.GetEntityTypeId())
	if err != nil {
		errMsg := "Error retrieving node of TestMetadata"
		t.Errorf("TestMetadataGetNodeTypeById returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has retrieved the following node: '%+v'", node)
}

func TestMetadataGetNodeTypes(t *testing.T) {
	testGmd := CreateTestGraphMetadata()
	nodeList, err := testGmd.GetNodeTypes()
	if err != nil {
		errMsg := "Error retrieving nodes of TestMetadata"
		t.Errorf("TestMetadataGetNodeTypes returned error message %s", errMsg)
	}
	t.Logf("TestGraphMetadata has following %d nodes", len(nodeList))
	for _, node := range nodeList {
		t.Logf("TestGraphMetadata has following node: '%+v'", node)
	}
}
