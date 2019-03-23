package model

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
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
 * File name: TGNodeType_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestNodeType(name string, entityType types.TGSystemType, parent types.TGEntityType) *NodeType {
	newNodeType := DefaultNodeType()
	newNodeType.name = name
	newNodeType.sysType = entityType
	attributes := make(map[string]*AttributeDescriptor, 3)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	attributes["BoolDesc"] = bAttrDesc
	attributes["IntegerDesc"] = iAttrDesc
	attributes["StringDesc"] = sAttrDesc
	newNodeType.attributes = attributes
	bPkAttrDesc := CreateTestAttributeDescriptor("BoolPkDesc", types.AttributeTypeBoolean)
	iPkAttrDesc := CreateTestAttributeDescriptor("IntegerPkDesc", types.AttributeTypeInteger)
	sPkAttrDesc := CreateTestAttributeDescriptor("StringPkDesc", types.AttributeTypeString)
	pKeys := []*AttributeDescriptor{bPkAttrDesc, iPkAttrDesc, sPkAttrDesc}
	newNodeType.SetPKeyAttributeDescriptors(pKeys)
	newNodeType.parent = parent
	return newNodeType
}

func TestGetPKeyAttributeDescriptors(t *testing.T) {
	newNodeType := CreateTestNodeType("GrandChild-1", types.SystemTypeNode, DefaultNodeType())

	pkDesc := newNodeType.GetPKeyAttributeDescriptors()
	for _, pk := range pkDesc {
		t.Logf("TestGetPKeyAttributeDescriptors returned attribute descriptor as '%+v'", pk)
	}
}

func TestNodeTypeDerivedFrom(t *testing.T) {
	parentEntityType := CreateTestNodeType("Node-1", types.SystemTypeNode, DefaultNodeType())

	newNodeType1 := CreateTestNodeType("ChildNode-1", types.SystemTypeNode, parentEntityType)
	newNodeType2 := CreateTestNodeType("GrandChild-1", types.SystemTypeNode, newNodeType1)

	parent := newNodeType1.DerivedFrom()
	t.Logf("Entitytype '%+v' is derived from '%+v' and has parent as '%+v'", newNodeType1, parentEntityType, parent.(*NodeType))

	parent = newNodeType2.DerivedFrom()
	t.Logf("Entitytype '%+v' is derived from '%+v' and has parent as '%+v'", newNodeType2, parentEntityType, parent.(*NodeType))
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestNodeTypeWriteExternal(t *testing.T) {
	parentNodeType := CreateTestNodeType("Node-1", types.SystemTypeNode, DefaultNodeType())
	ToBeExportedNodeType := CreateTestNodeType("Edge-1", types.SystemTypeNode, parentNodeType)
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = ToBeExportedNodeType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedNodeType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	_ = ToBeExportedNodeType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedNodeType)
}
