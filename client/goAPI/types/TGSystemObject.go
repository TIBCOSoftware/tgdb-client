package types

import "bytes"

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
 * File name: TGSystemObject.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= System Types =======
type TGSystemType int

const (
	SystemTypeInvalid TGSystemType = -1
	SystemTypeEntity  TGSystemType = -2
)
const (
	SystemTypeAttributeDescriptor TGSystemType = iota
	SystemTypeNode
	SystemTypeEdge
	SystemTypeIndex
	SystemTypePrincipal
	SystemTypeRole
	SystemTypeSequence
	SystemTypeMaxSysObject
)

func (systemType TGSystemType) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if systemType&SystemTypeInvalid == SystemTypeInvalid {
		buffer.WriteString("SystemTypeInvalid")
	} else if systemType&SystemTypeAttributeDescriptor == SystemTypeAttributeDescriptor {
		buffer.WriteString("SystemTypeAttributeDescriptor")
	} else if systemType&SystemTypeEntity == SystemTypeEntity {
		buffer.WriteString("SystemTypeEntity")
	} else if systemType&SystemTypeNode == SystemTypeNode {
		buffer.WriteString("SystemTypeNode")
	} else if systemType&SystemTypeEdge == SystemTypeEdge {
		buffer.WriteString("SystemTypeEdge")
	} else if systemType&SystemTypeIndex == SystemTypeIndex {
		buffer.WriteString("SystemTypeIndex")
	} else if systemType&SystemTypePrincipal == SystemTypePrincipal {
		buffer.WriteString("SystemTypePrincipal")
	} else if systemType&SystemTypeRole == SystemTypeRole {
		buffer.WriteString("SystemTypeRole")
	} else if systemType&SystemTypeSequence == SystemTypeSequence {
		buffer.WriteString("SystemTypeSequence")
	} else if systemType&SystemTypeMaxSysObject == SystemTypeMaxSysObject {
		buffer.WriteString("SystemTypeMaxSysObject")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

type TGSystemObject interface {
	TGSerializable
	// GetName gets the system object's name
	GetName() string
	// GetSystemType gets the system object's type
	GetSystemType() TGSystemType
}
