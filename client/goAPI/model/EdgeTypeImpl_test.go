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
 * File name: TGEdgeType_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestEdgeType(name string, directionType types.TGDirectionType, entityType types.TGSystemType, parent types.TGEntityType) *EdgeType {
	newEdgeType := DefaultEdgeType()
	newEdgeType.SetName(name)
	newEdgeType.sysType = entityType
	attributes := make(map[string]*AttributeDescriptor, 3)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	attributes["BoolDesc"] = bAttrDesc
	attributes["IntegerDesc"] = iAttrDesc
	attributes["StringDesc"] = sAttrDesc
	newEdgeType.attributes = attributes
	newEdgeType.directionType = directionType
	newEdgeType.parent = parent
	return newEdgeType
}

func TestEdgeTypeDerivedFrom(t *testing.T) {
	parentEntityType := CreateTestEdgeType("Node-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, DefaultEntityType())

	newEdgeType1 := CreateTestEdgeType("Edge-1", types.DirectionTypeDirected, types.SystemTypeEdge, parentEntityType)
	newEdgeType2 := CreateTestEdgeType("Edge-2", types.DirectionTypeDirected, types.SystemTypeEdge, parentEntityType)

	parent := newEdgeType1.DerivedFrom()
	t.Logf("Entitytype '%+v' is derived from '%+v' and has parent as '%+v'", newEdgeType1, parentEntityType, parent)

	parent = newEdgeType2.DerivedFrom()
	t.Logf("Entitytype '%+v' is derived from '%+v' and has parent as '%+v'", newEdgeType2, parentEntityType, parent)
}

func TestGetFromNodeType(t *testing.T) {
	parentEntityType := CreateTestEdgeType("Node-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, DefaultEntityType())

	newEdgeType1 := CreateTestEdgeType("Edge-1", types.DirectionTypeDirected, types.SystemTypeEdge, parentEntityType)
	//newEdgeType2 := CreateTestEdgeType("Edge-2", types.DirectionOutbound, types.SystemTypeEdge, parentEntityType)

	fromNodeType := newEdgeType1.GetFromNodeType()
	t.Logf("TestGetFromNodeType returned fromNodeType as '%+v'", fromNodeType)
}

func TestGetToNodeType(t *testing.T) {
	parentEntityType := CreateTestEdgeType("Node-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, DefaultEntityType())

	newEdgeType1 := CreateTestEdgeType("Edge-1", types.DirectionTypeDirected, types.SystemTypeEdge, parentEntityType)
	//newEdgeType2 := CreateTestEdgeType("Edge-2", types.DirectionOutbound, types.SystemTypeEdge, parentEntityType)

	toNodeType := newEdgeType1.GetToNodeType()
	t.Logf("TestGetFromNodeType returned toNodeType as '%+v'", toNodeType)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestEdgeTypeWriteExternal(t *testing.T) {
	parentEntityType := CreateTestEdgeType("Node-1", types.DirectionTypeBiDirectional, types.SystemTypeNode, DefaultEntityType())

	ToBeExportedEdgeType := CreateTestEdgeType("Edge-1", types.DirectionTypeDirected, types.SystemTypeEdge, parentEntityType)
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = ToBeExportedEdgeType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedEdgeType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	_ = ToBeExportedEdgeType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedEdgeType)
}
