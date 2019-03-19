package utils

import "testing"

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
 * File name: TGProtocolVersion_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TODO: Revisit later - for more testing
//func TestGetAsSortedProperties(t *testing.T) {
//	testEnv := NewTGEnvironment()
//	//t.Logf("TestGetChannelClientId has properties as '%+v' - '%+v'", len(testEnv.TGEnv), testEnv.TGEnv)
//	v := testEnv.GetAsSortedProperties()
//	t.Logf("Got '%digits' environment properties", len(v.(*SortedProperties).Properties))
//	for _, v := range v.(*SortedProperties).Properties {
//		t.Logf("Got environment property as '%+v':'%+v'", v.KeyName, v.KeyValue)
//	}
//}

func TestGetChannelClientId(t *testing.T) {
	testEnv := NewTGEnvironment()
	//t.Logf("TestGetChannelClientId has properties as '%+v' - '%+v'", len(testEnv.TGEnv), testEnv.TGEnv)
	v := testEnv.GetEnvironmentProperty("clientId")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelClientId()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestChannelConnectTimeout(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("connectTimeout")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelConnectTimeout()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestChannelDefaultHost(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("defaultHost")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelDefaultHost()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelDefaultPort(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("defaultPort")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelDefaultPort()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelDefaultUser(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("defaultUserID")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelDefaultUser()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelFTHosts(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("ftHosts")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelFTHosts()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelPingInterval(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("pingInterval")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelPingInterval()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelSendSize(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("sendSize")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelSendSize()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelReceiveSize(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("recvSize")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelReceiveSize()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetChannelUser(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("userID")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelUser()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetConnectionPoolDefaultPoolSize(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("defaultPoolSize")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetConnectionPoolDefaultPoolSize()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetEnvironmentProperty(t *testing.T) {
	testEnv := NewTGEnvironment()
	v := testEnv.GetEnvironmentProperty("pingInterval")
	t.Logf("Got Value: %+v", v)
	exp := testEnv.GetChannelPingInterval()
	t.Logf("Expected Value: %+v", exp)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestEnvironmentSetProperty(t *testing.T) {
	testEnv := NewTGEnvironment()
	beforeValue := testEnv.GetEnvironmentProperty("pingInterval")
	testEnv.SetEnvironmentProperty("pingInterval", "60")
	afterValue := testEnv.GetEnvironmentProperty("pingInterval")
	if err := assertEq(beforeValue, afterValue); err != nil {
		t.Logf("Wanted %+v; Got %+v", "60", afterValue)
	}
}
