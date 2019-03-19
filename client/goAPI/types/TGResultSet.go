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
 * File name: TGResultSet.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGResultSet interface {
	// AddEntityToResultSet adds another entity to the result set
	AddEntityToResultSet(entity TGEntity) TGResultSet
	// Close closes the result set
	Close() TGResultSet
	// Count returns nos of entities returned by the query. The result set has a cursor which prefetches
	// "n" rows as per the query constraint. If the nos of entities returned by the query is less
	// than prefetch count, then all are returned.
	Count() int
	// First returns the first entity in the result set
	First() TGEntity
	// Last returns the last Entity in the result set
	Last() TGEntity
	// GetAt gets the entity at the position.
	GetAt(position int) TGEntity
	// GetExceptions gets the Exceptions in the result set
	GetExceptions() []TGError
	// GetPosition gets the Current cursor position. A result set upon creation is set to the position 0.
	GetPosition() int
	// HasExceptions checks whether the result set has any exceptions
	HasExceptions() bool
	// HasNext Check whether there is next entry in result set
	HasNext() bool
	// Next returns the next entity w.r.t to the current cursor position in the result set
	Next() TGEntity
	// Prev returns the prev entity w.r.t to the current cursor position in the result set
	Prev() TGEntity
	// Skip skips a number of position
	Skip(position int) TGResultSet
	// Additional Method to help debugging
	String() string
}
