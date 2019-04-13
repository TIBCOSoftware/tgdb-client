package admin

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
 * File name: TGIndexInfo.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TGIndexInfo users to retrieve the index information from server
type TGIndexInfo interface {
	// GetAttributes returns a collection of attribute names
	GetAttributeNames() []string
	// GetName returns the index name
	GetName() string
	// GetType returns the index type
	GetType() byte
	// GetSystemId returns the system ID
	GetSystemId() int
	// GetNodeTypes returns a collection of node types
	GetNodeTypes() []string
	// IsUnique returns the information whether the index is unique
	IsUnique() bool
}
