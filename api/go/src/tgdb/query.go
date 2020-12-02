/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name: query.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: query.go 4144 2020-07-09 18:17:49Z nimish $
 */

package tgdb

import "time"

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
	First() interface{}
	// Last returns the last Entity in the result set
	Last() interface{}
	// GetAt gets the entity at the position.
	GetAt(position int) interface{}
	// GetExceptions gets the Exceptions in the result set
	GetExceptions() []TGError
	// GetPosition gets the Current cursor position. A result set upon creation is set to the position 0.
	GetPosition() int
	// HasExceptions checks whether the result set has any exceptions
	HasExceptions() bool
	// HasNext Check whether there is next entry in result set
	HasNext() bool
	// Next returns the next entity w.r.t to the current cursor position in the result set
	Next() interface{}
	// Prev returns the prev entity w.r.t to the current cursor position in the result set
	Prev() interface{}
	// Skip skips a number of position
	Skip(position int) TGResultSet
	// Additional Method to help debugging
	String() string
	// ToCollection converts the result set into a collection
	ToCollection() []interface{}

	GetMetadata() TGResultSetMetaData
}

// A Set of QueryOption that allows the user manipulate the results of the query
type TGQueryOption interface {
	TGProperties
	// GetBatchSize gets the current value of the batch size
	GetBatchSize() int
	// SetBatchSize sets a limit on the batch. Default is 50
	SetBatchSize(size int) TGError
	// GetPreFetchSize gets the current value of the pre-fetch size
	GetPreFetchSize() int
	// SetPreFetchSize sets a limit on the number of entities(nodes and edges) return in a query. Default is 1000
	SetPreFetchSize(size int) TGError
	// GetTraversalDepth gets the current value of traversal depth
	GetTraversalDepth() int
	// SetTraversalDepth sets the additional level of traversal from the query result set. Default is 3.
	SetTraversalDepth(depth int) TGError
	// GetEdgeLimit gets the current value of edge limit
	GetEdgeLimit() int
	// SetEdgeLimit sets the number of edges per node to be returned in a query.  Default is 0 which means unlimited.
	SetEdgeLimit(limit int) TGError
	// GetSortAttrName gets sort attribute name
	GetSortAttrName() string
	// SetSortAttrName sets sort attribute name
	SetSortAttrName(name string) TGError
	// IsSortOrderDsc gets sort order desc
	IsSortOrderDsc() bool
	// SetSortOrderDsc sets sort order desc
	SetSortOrderDsc(isDsc bool) TGError
	// GetSortResultLimit gets sort result limit
	GetSortResultLimit() int
	// SetSortResultLimit sets sort result limit
	SetSortResultLimit(limit int) TGError
	// Additional Method to help debugging
	String() string
}



type TGQuery interface {
	// Close closes the Query
	Close()
	// Execute executes the Query
	Execute() TGResultSet
	// SetBoolean sets Boolean parameter
	SetBoolean(name string, value bool)
	// SetBytes sets Byte Parameter
	SetBytes(name string, bos []byte)
	// SetChar sets Character Parameter
	SetChar(name string, value string)
	// SetDate sets Date Parameter
	SetDate(name string, value time.Time)
	// SetDouble sets Double Parameter
	SetDouble(name string, value float64)
	// SetFloat sets Float Parameter
	SetFloat(name string, value float32)
	// SetInt sets Integer Parameter
	SetInt(name string, value int)
	// SetLong sets Long Parameter
	SetLong(name string, value int64)
	// SetNull sets the parameter to null
	SetNull(name string)
	// SetOption sets the Query Option
	SetOption(options TGQueryOption)
	// SetShort sets Short Parameter
	SetShort(name string, value int16)
	// SetString sets String Parameter
	SetString(name string, value string)
	// Additional Method to help debugging
	String() string
}

type TGFilter interface {
}

type TGTraversalDescriptor interface {
	// Traverse the graph using starting points provided
	Traverse(startingPoints []TGNode) TGResultSet
}








const (
	TYPE_UNKNOWN = iota
	TYPE_OBJECT
	TYPE_ENTITY
	TYPE_ATTR
	TYPE_NODE
	TYPE_EDGE
	TYPE_LIST
	TYPE_MAP
	TYPE_TUPLE
	TYPE_SCALAR
	TYPE_PATH
)


type TGResultDataDescriptor interface {
	GetDataType() int
	GetContainedDataSize() int
	IsMap() bool
	IsArray() bool
	HasConcreteType() bool
	GetScalarType() int//TGAttributeType
	GetSystemObject() TGSystemObject
	GetKeyDescriptor() TGResultDataDescriptor
	GetValueDescriptor() TGResultDataDescriptor
	GetContainedDescriptors() []TGResultDataDescriptor
	GetContainedDescriptor(position int) TGResultDataDescriptor
}

type TGResultSetMetaData interface {
	GetResultDataDescriptor () *TGResultDataDescriptor
	GetResultType() int
	GetAnnot () string
}
