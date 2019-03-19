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
 * File name: TGQueryOption.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
