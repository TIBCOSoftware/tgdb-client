package model

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"testing"
	"time"
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
 * File name: TGKey_Test.go
 * Created on: Nov 17, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

func TestSetOrCreateAttribute(t *testing.T) {
	gmd := CreateTestGraphMetadata()
	newKey := NewCompositeKey(gmd, "key-1")
	t.Logf("New Composite key '%+v' is created from '%+v'", newKey, gmd)

	_ = newKey.SetOrCreateAttribute("ShortDesc", 1)
	_ = newKey.SetOrCreateAttribute("LongDesc", 222)
	_ = newKey.SetOrCreateAttribute("BoolDesc", false)
	_ = newKey.SetOrCreateAttribute("NumberDesc", 5.0)
	_ = newKey.SetOrCreateAttribute("TimeDesc", time.Now().Unix())

	t.Logf("Composite key has been modified to add '%d' attributes as '%+v'", len(newKey.attributes), newKey)
}

// This automatically will test both APIs - (a) ReadExternal and (b) WriteExternal
func TestKeyWriteExternal(t *testing.T) {
	gmd := CreateTestGraphMetadata()
	newKeyToBeExported := NewCompositeKey(gmd, "key-1")
	t.Logf("New Composite key '%+v' is created from '%+v'", newKeyToBeExported, gmd)

	_ = newKeyToBeExported.SetOrCreateAttribute("LongDesc", 1)
	_ = newKeyToBeExported.SetOrCreateAttribute("StringDesc", "2")
	_ = newKeyToBeExported.SetOrCreateAttribute("BoolDesc", true)

	//var network bytes.Buffer
	oNetwork := iostream.DefaultProtocolDataOutputStream()

	_ = newKeyToBeExported.WriteExternal(oNetwork)
	t.Logf("EntityType WriteExternal exported entity type value '%+v' as '%+v'", newKeyToBeExported, string(oNetwork.Buf))

	iNetwork := iostream.DefaultProtocolDataInputStream()
	_ = newKeyToBeExported.ReadExternal(iNetwork)
	t.Logf("EntityType ReadExternal imported entity type as '%+v'", newKeyToBeExported)
}
