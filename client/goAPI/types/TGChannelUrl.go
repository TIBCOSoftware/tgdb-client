package types

import "bytes"

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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGChannelUrl.go
 * Created on: Oct 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Protocol Types Supported =======
type TGProtocol int

const (
	ProtocolTCP TGProtocol = 1 << iota
	ProtocolSSL
	ProtocolHTTP
	ProtocolHTTPS
)

func (proType TGProtocol) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if proType&ProtocolTCP == ProtocolTCP {
		buffer.WriteString("tcp")
	} else if proType&ProtocolSSL == ProtocolSSL {
		buffer.WriteString("ssl")
	} else if proType&ProtocolHTTP == ProtocolHTTP {
		buffer.WriteString("http")
	} else if proType&ProtocolHTTPS == ProtocolHTTPS {
		buffer.WriteString("https")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

// Channel URL is an encapsulation of all the attributes necessary to construct a valid and meaningful URL to connect to TGDB server
type TGChannelUrl interface {
	// GetFTUrls gets the Fault Tolerant URLs
	GetFTUrls() []TGChannelUrl
	// GetHost gets the host part of the URL
	GetHost() string
	// GetPort gets the port on which it is connected
	GetPort() int
	// GetProperties gets the URL Properties
	GetProperties() TGProperties
	// GetProtocol gets the protocol used as part of the URL
	GetProtocol() TGProtocol
	// GetUrlAsString gets the string form of the URL
	GetUrlAsString() string
	// GetUser gets the user associated with the URL
	GetUser() string
}
