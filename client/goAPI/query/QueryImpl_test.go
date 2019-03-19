package query

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
 * File name: TGQuery_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TODO: Revisit later - once connection is implemented - to create proper test queries against default meta data
//func TestQueryClose(t *testing.T) {
//	testQry := DefaultQuery()
//	testQry.queryHashId = 1234567890
//	t.Logf("Before closing connection, this query object has %d connections", testQry.connection)
//	testQry.Close()
//}
//
//func TestQueryExecute(t *testing.T) {
//	testQry := DefaultQuery()
//	testQry.queryHashId = 1234567890
//	//	testQry.qryOption.AddProperty()
//	t.Logf("Before executing query, this query object has %d connections", testQry.connection)
//	rs := testQry.Execute()
//	if rs != nil {
//		t.Logf("Before executing query, the result set has %d records", rs.Count())
//	}
//}
//func TestQuerySetOption(t *testing.T) {
//	testQry := DefaultQuery()
//	testQry.queryHashId = 1234567890
//	t.Logf("Before closing connection, this query object has %d connections", testQry.connection)
//	testQry.Close()
//}
//
