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
 * File name: TGProtocolVersion_Test.go
 * Created on: Nov 10, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func TestGetConfigFromKey(t *testing.T) {
	cn := GetConfigFromKey(ChannelPingInterval)
	if err := assertEq(cn.GetName(), "ChannelPingInterval"); err != nil {
		t.Logf("Wanted %+v; Got %+v", "ChannelPingInterval", cn.GetName())
	}
}

func TestGetConfigFromName(t *testing.T) {
	name := "defaultUserID"
	cn := GetConfigFromName(name)
	if err := assertEq(cn.GetName(), name); err != nil {
		t.Logf("Wanted %+v; Got %+v", name, cn.GetName())
	}
}
