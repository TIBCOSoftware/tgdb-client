package utils

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
 * File name: TGProperties_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func TestAddProperty(t *testing.T) {
	testPropertySet := NewSortedProperties()
	//t.Logf("TestAddProperty has properties as '%+v'", testPropertySet.Properties)
	cn := GetConfigFromKey(ChannelConnectTimeout)
	//t.Logf("TestAddProperty has config as '%+v'", cn)
	t.Logf("TestAddProperty adding property as '%+v':'%+v'", "connectTimeout", 60)
	testPropertySet.AddProperty("connectTimeout", "60")
	//t.Logf("TestAddProperty has properties as '%+v'", testPropertySet.Properties)
	exp := 60
	t.Logf("Expected Value: '%+v'", exp)
	v := testPropertySet.GetProperty(cn, "240")
	t.Logf("Got Value: '%+v'", v)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetProperty(t *testing.T) {
	testPropertySet := NewSortedProperties()
	cn := GetConfigFromKey(ChannelPingInterval)
	testPropertySet.AddProperty("pingInterval", "30")
	exp := 30
	t.Logf("Expected Value: '%+v'", exp)
	v := testPropertySet.GetProperty(cn, "240")
	t.Logf("Got Value: '%+v'", v)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestGetPropertyAsBoolean(t *testing.T) {
	testPropertySet := NewSortedProperties()
	cn := GetConfigFromKey(TlsVerifyDatabaseName)
	testPropertySet.AddProperty("verifyDBName", "true")
	exp := true
	t.Logf("Expected Value: %+v", exp)
	v := testPropertySet.GetPropertyAsBoolean(cn)
	t.Logf("Got Value: %+v", v)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestSetProperty(t *testing.T) {
	testPropertySet := NewSortedProperties()
	cn := GetConfigFromKey(ChannelPingInterval)
	testPropertySet.AddProperty("pingInterval", "30")
	v := testPropertySet.GetProperty(cn, "0")
	//t.Logf("After adding, Got Value: '%+v'", v)
	testPropertySet.SetProperty("pingInterval", "180")
	exp := 180
	t.Logf("Expected Value: '%+v'", exp)
	v = testPropertySet.GetProperty(cn, "0")
	t.Logf("Got Value: '%+v'", v)
	if err := assertEq(exp, v); err != nil {
		t.Logf("Wanted %+v; Got %+v", exp, v)
	}
}

func TestSetUserAndPassword(t *testing.T) {
	testPropertySet := NewSortedProperties()
	//cnu := getConfigName(ChannelUserID)
	//cnp := getConfigName(ChannelPassword)
	user := "ComboUser"
	pwd := "ComboPwd"
	err := SetUserAndPassword(testPropertySet, user, pwd)
	//t.Logf("Got Value: %+v", v)
	if err != nil {
		t.Logf("TGProperties: SetUserAndPassword resulted in '%+v'", err)
	}
}
