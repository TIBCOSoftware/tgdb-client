package utils

import (
	"encoding/binary"
	"fmt"
	"reflect"
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

func assertEq(exp, got interface{}) error {
	if !reflect.DeepEqual(exp, got) {
		return fmt.Errorf("wanted %v; got %v", exp, got)
	}
	return nil
}

func assertNotEq(exp, got interface{}) error {
	if reflect.DeepEqual(exp, got) {
		return fmt.Errorf("wanted %v; got %v", exp, got)
	}
	return nil
}

func TestIsCompatibleSuccess(t *testing.T) {
	b := []byte{1, 1}
	version := binary.BigEndian.Uint16(b)
	if IsCompatible(version) {
		t.Logf("TGProtocolVersion: %+v is compatible w/ TGDB GO Client API version.", b)
	}
}

func TestIsCompatibleError(t *testing.T) {
	b := []byte{2, 0}
	version := binary.BigEndian.Uint16(b)
	if err := assertEq(GetProtocolVersion(), version); err != nil {
		t.Logf("Wanted %+v; Got %+v", GetProtocolVersion(), version)
	}
}
