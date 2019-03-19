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
 * File name: TGEdge.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Various Direction Types for EDGE (Entity) type =======
type TGDirectionType int

const (
	DirectionTypeUnDirected TGDirectionType = iota
	DirectionTypeDirected
	DirectionTypeBiDirectional
)

func (directionType TGDirectionType) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if directionType&DirectionTypeUnDirected == DirectionTypeUnDirected {
		buffer.WriteString("UnDirected")
	} else if directionType&DirectionTypeDirected == DirectionTypeDirected {
		buffer.WriteString("Directed")
	} else if directionType&DirectionTypeBiDirectional == DirectionTypeBiDirectional {
		buffer.WriteString("BiDirectional")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

// ======= Various Directions associated with EDGE (Entity) type =======
type TGDirection int

const (
	DirectionInbound TGDirection = iota
	DirectionOutbound
	DirectionAny
)

func (direction TGDirection) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if direction&DirectionInbound == DirectionInbound {
		buffer.WriteString("Inbound")
	}
	if direction&DirectionOutbound == DirectionOutbound {
		buffer.WriteString("Outbound")
	}
	if direction&DirectionAny == DirectionAny {
		buffer.WriteString("Any")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

type TGEdge interface {
	TGEntity
	// GetDirectionType gets direction type as one of the constants
	GetDirectionType() TGDirectionType
	// GetVertices gets array of NODE (Entity) types for this EDGE (Entity) type
	GetVertices() []TGNode
}
