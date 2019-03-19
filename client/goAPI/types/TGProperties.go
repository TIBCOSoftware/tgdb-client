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
 * File name: TGProperties.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGProperties interface {
	// AddProperty checks whether a property already exists, else adds a new property in the form of name=value pair
	AddProperty(name, value string)
	// GetProperty gets the property either with value or default value
	GetProperty(cn TGConfigName, value string) string
	// SetProperty sets existing property value in the form of name=value pair
	SetProperty(name, value string)
	// SetUserAndPassword sets urlUser and password
	//SetUserAndPassword(user, pwd string) TGError
	// GetPropertyAsInt gets Property as int value
	GetPropertyAsInt(cn TGConfigName) int
	// GetPropertyAsLong gets Property as long value
	GetPropertyAsLong(cn TGConfigName) int64
	// GetPropertyAsBoolean gets Property as bool value
	GetPropertyAsBoolean(cn TGConfigName) bool
}
