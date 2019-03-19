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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGConstants.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	//public static final List<?> EmptyList = new ArrayList<>();
	//EmptyByteArray []byte = []byte{}
	EmptyString string = ""

	U64_NULL       int64 = 0xfffffffffffffff
	U64PACKED_NULL byte  = 0xf0

	INTERNAL_SERVER_ERROR    string = "TGDB-00001"
	TGDB_HNDSHKRESP_ERROR    string = "TGDB-HNDSHKRESP-ERR"
	TGDB_CHANNEL_ERROR       string = "TGDB-CHANNEL-ERR"
	TGDB_SEND_ERROR          string = "TGDB-SENDL-ERR"
	TGDB_CLIENT_READEXTERNAL string = "TGDB-CLIENT-READEXTERNAL"

	DebugEnabled bool = false
)
