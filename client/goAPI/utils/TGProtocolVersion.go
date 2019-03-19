package utils

import (
	"encoding/binary"
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
 * File name: TGProtocolVersion.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	TgMajorVersion uint8 = 1
	TgMinorVersion uint8 = 0
	TgMagic        int   = 0xdb2d1e4 // TGDecimal: 229822948
)

func GetMagic() int {
	return TgMagic
}

func GetProtocolVersion() uint16 {
	b := []byte{TgMajorVersion, TgMinorVersion}
	return binary.BigEndian.Uint16(b)
}

func IsCompatible(protocolVersion uint16) bool {
	return protocolVersion == GetProtocolVersion()
}
