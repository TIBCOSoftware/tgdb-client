package types

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
 * File name: TGEdgeType.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGEdgeType interface {
	TGEntityType
	// GetDirectionType gets direction type as one of the constants
	GetDirectionType() TGDirectionType
	// GetFromNodeType gets From-Node Type
	GetFromNodeType() TGNodeType
	// GetFromTypeId gets From-Node ID
	GetFromTypeId() int
	// GetToNodeType gets To-Node Type
	GetToNodeType() TGNodeType
	// GetToTypeId gets To-Node ID
	GetToTypeId() int
	// SetFromNodeType sets From-Node Type
	SetFromNodeType(fromNode TGNodeType)
	// SetFromTypeId sets From-Node ID
	SetFromTypeId(fromTypeId int)
	// SetToNodeType sets From-Node Type
	SetToNodeType(toNode TGNodeType)
	// SetToTypeId sets To-Node ID
	SetToTypeId(toTypeId int)
}
