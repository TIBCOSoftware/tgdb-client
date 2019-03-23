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
 * File name: TGEntityType_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func CreateTestEntityType(name string, entityType types.TGSystemType) *EntityType {
	newEntityType := DefaultEntityType()
	newEntityType.name = name
	newEntityType.sysType = entityType
	attributes := make(map[string]*AttributeDescriptor, 3)
	bAttrDesc := CreateTestAttributeDescriptor("BoolDesc", types.AttributeTypeBoolean)
	iAttrDesc := CreateTestAttributeDescriptor("IntegerDesc", types.AttributeTypeInteger)
	sAttrDesc := CreateTestAttributeDescriptor("StringDesc", types.AttributeTypeString)
	attributes["BoolDesc"] = bAttrDesc
	attributes["IntegerDesc"] = iAttrDesc
	attributes["StringDesc"] = sAttrDesc
	newEntityType.attributes = attributes
	return newEntityType
}

func TestEntityTypeDerivedFrom(t *testing.T) {
	parentEntityType := CreateTestEntityType("Node-1", types.SystemTypeNode)

	newEntityType := CreateTestEntityType("Edge-1", types.SystemTypeEntity)
	newEntityType.parent = parentEntityType

	parent := newEntityType.DerivedFrom()
	t.Logf("Entitytype '%+v' is derived from '%+v' and has parent as '%+v'", newEntityType, parentEntityType, parent)
}

func TestGetAttributeDescriptor(t *testing.T) {
	parentEntityType := CreateTestEntityType("Node-1", types.SystemTypeNode)

	newEntityType := CreateTestEntityType("Edge-1", types.SystemTypeEntity)
	newEntityType.parent = parentEntityType

	desc := newEntityType.GetAttributeDescriptor("StringDesc")
	t.Logf("Entitytype '%+v' returned descriptor as '%+v'", newEntityType, desc.(*AttributeDescriptor))
}

func TestGetAttributeDescriptors(t *testing.T) {
	parentEntityType := CreateTestEntityType("Node-1", types.SystemTypeNode)

	newEntityType := CreateTestEntityType("Edge-1", types.SystemTypeEntity)
	newEntityType.parent = parentEntityType

	descList := newEntityType.GetAttributeDescriptors()
	t.Logf("Entitytype '%+v' returned descriptors as '%+v'", newEntityType, descList)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestEntityTypeWriteExternal(t *testing.T) {
	ToBeExportedEntityType := CreateTestEntityType("Node-1", types.SystemTypeNode)
	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = ToBeExportedEntityType.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", ToBeExportedEntityType, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	//TobeImportedEntityType := DefaultEntityType()
	_ = ToBeExportedEntityType.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", ToBeExportedEntityType)
}
