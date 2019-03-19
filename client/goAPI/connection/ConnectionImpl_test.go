package connection

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
 * File name: ConnectionImpl_Test.go
 * Created on: Nov 24, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	url = "tcp://scott@localhost:8222"
	password = "scott"
	prefetchMetaData = false
	testUrl0 = "foo1.bar.com"
	testUrl1 = "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
	testUrl2 = "tcp://scott@foo.bar.com:8700"
	testUrl3 = "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
	testUrl4 = "http://[2001:db8:1f70::999:de8:7648:6e8]:100/{userID=Admin}"
)

/**
	TODO: Please make sure that the server is up and listening on the correct port, before exeucting these tests
*/

//func TestConnectionImpl_GetGraphObjectFactory(t *testing.T) {
//	connFactory := NewTGConnectionFactory()
//	conn, err := connFactory.CreateConnection(url, "", password, nil)
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphObjectFactory - error during CreateConnection")
//		return
//	}
//
//	err = conn.Connect()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphObjectFactory - error during conn.Connect")
//		return
//	}
//
//	gof, err := conn.GetGraphObjectFactory()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphObjectFactory - error during conn.GetGraphObjectFactory")
//	}
//	if gof == nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphObjectFactory - Graph Object Factory is null")
//	}
//	t.Logf("Returning from TestConnectionImpl_GetGraphObjectFactory w/ GOF as '%+v'", gof)
//
//	err = conn.Disconnect()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphObjectFactory - error during conn.Connect")
//		return
//	}
//	t.Log("Returning from TestConnectionImpl_GetGraphObjectFactory - successfully disconnected.")
//}
//
//func TestConnectionImpl_GetGraphMetadata(t *testing.T) {
//	connFactory := NewTGConnectionFactory()
//	conn, err := connFactory.CreateConnection(url, "", password, nil)
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during CreateConnection")
//		return
//	}
//
//	err = conn.Connect()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.Connect")
//		return
//	}
//
//	gof, err := conn.GetGraphObjectFactory()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.GetGraphObjectFactory")
//	}
//	if gof == nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - Graph Object Factory is null")
//	}
//	t.Logf("Returning from TestConnectionImpl_GetGraphMetadata w/ GOF as '%+v'", gof)
//
//	gmd, err := conn.GetGraphMetadata(true)
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.GetGraphMetadata")
//	}
//
//	attrDesc, err := gmd.GetAttributeDescriptor("factor")
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.GetAttributeDescriptor('factor')")
//	}
//	if attrDesc != nil {
//		t.Logf("'factor' has id %d and data desc: %s\n", attrDesc.GetAttributeId(), types.GetAttributeTypeFromId(attrDesc.GetAttrType()).TypeName)
//	} else {
//		t.Log("'factor' is not found from meta data fetch")
//	}
//
//	attrDesc, err = gmd.GetAttributeDescriptor("level")
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.GetAttributeDescriptor('level')")
//	}
//	if attrDesc != nil {
//		t.Logf("'level' has id %d and data desc: %s\n", attrDesc.GetAttributeId(), types.GetAttributeTypeFromId(attrDesc.GetAttrType()).TypeName)
//	} else {
//		t.Log("'level' is not found from meta data fetch")
//	}
//
//	testNodeType, err := gmd.GetNodeType("testnode")
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.GetNodeType('testnode')")
//	}
//	if testNodeType != nil {
//		t.Logf("'testnode' is found with %d attributes\n", len(testNodeType.GetAttributeDescriptors()))
//	} else {
//		t.Log("'testnode' is not found from meta data fetch")
//	}
//
//	err = conn.Disconnect()
//	if err != nil {
//		t.Error("Returning from TestConnectionImpl_GetGraphMetadata - error during conn.Connect")
//	}
//	t.Log("Returning from TestConnectionImpl_GetGraphMetadata - successfully disconnected.")
//}
