package query

import (
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
 * File name: TGResultSet_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestSetPreFetchSize(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.GetPreFetchSize()
	t.Logf("Before setting new value of prefetch size, the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetPreFetchSize(10)
	if err != nil {
		t.Errorf("SetPreFetchSize returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.GetPreFetchSize()
	t.Logf("After setting prefetch size to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}

// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestSetTraversalDepth(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.GetTraversalDepth()
	t.Logf("Before setting new value of traversal depth, the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetTraversalDepth(20)
	if err != nil {
		t.Errorf("SetTraversalDepth returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.GetTraversalDepth()
	t.Logf("After setting traversal depth to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}

//// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestSetEdgeLimit(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.GetEdgeLimit()
	t.Logf("Before setting new value of edge limit, the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetEdgeLimit(25)
	if err != nil {
		t.Errorf("SetEdgeLimit returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.GetEdgeLimit()
	t.Logf("After setting edge limit to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}

//// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestSetSortAttrName(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.GetSortAttrName()
	t.Logf("Before setting new value of sort attr name, the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetSortAttrName("EmployeeId")
	if err != nil {
		t.Errorf("SetPreFetchSize returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.GetSortAttrName()
	t.Logf("After setting sort attr name to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}

//// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestIsSortOrderDsc(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.IsSortOrderDsc()
	t.Logf("Before setting new value of sort order (desc), the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetSortOrderDsc(true)
	if err != nil {
		t.Errorf("SetSortOrderDsc returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.IsSortOrderDsc()
	t.Logf("After setting sort order (desc) to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}

//// This will test both APIs in the following order - (a) Get, (b) Set, and (c) Get
func TestSetSortResultLimit(t *testing.T) {
	testQryOption := NewQueryOption()
	beforeVal := testQryOption.GetSortResultLimit()
	t.Logf("Before setting new value of sort result limit, the old value is '%+v'\n", beforeVal)
	err := testQryOption.SetSortResultLimit(100)
	if err != nil {
		t.Errorf("SetSortResultLimit returned error as '%+v'\n", err)
	}
	afterVal := testQryOption.GetSortResultLimit()
	t.Logf("After setting sort result limit to '%+v', the old value is '%+v'\n", beforeVal, afterVal)
}
