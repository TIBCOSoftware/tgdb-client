package channel

import (
	"crypto/tls"
	"strings"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGCipherSuite.go
 * Created on: Feb 10, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGCipherSuite struct {
	suiteId     uint16
	opensslName string
	keyExch     string
	encryption  string
	bits        string
}

var PreDefinedCipherSuites = map[string]TGCipherSuite{
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         {tls.TLS_RSA_WITH_AES_128_CBC_SHA256, "AES128-SHA256", "RSA", "AES", "128"},
	"TLS_RSA_WITH_AES_256_CBC_SHA256":         {0x3d, "AES256-SHA256", "RSA", "AES", "256"},
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256":     {0x67, "DHE-RSA-AES128-SHA256", "DH", "AES", "128"},
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256":     {0x6b, "DHE-RSA-AES256-SHA256", "DH", "AES", "256"},
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         {tls.TLS_RSA_WITH_AES_128_GCM_SHA256, "AES128-GCM-SHA256", "RSA", "AESGCM", "128"},
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         {tls.TLS_RSA_WITH_AES_256_GCM_SHA384, "AES256-GCM-SHA384", "RSA", "AESGCM", "256"},
	"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256":     {0x9e, "DHE-RSA-AES128-GCM-SHA256", "DH", "AESGCM", "128"},
	"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384":     {0x9f, "DHE-RSA-AES256-GCM-SHA384", "DH", "AESGCM", "256"},
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": {tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, "ECDHE-ECDSA-AES128-SHA256", "ECDH", "AES", "128"},
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384": {0xc024, "ECDHE-ECDSA-AES256-SHA384", "ECDH", "AES", "256"},
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   {tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, "ECDHE-RSA-AES128-SHA256", "ECDH", "AES", "128"},
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384":   {0xc028, "ECDHE-RSA-AES256-SHA384", "ECDH", "AES", "256"},
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": {tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, "ECDHE-ECDSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128"},
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": {tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, "ECDHE-ECDSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256"},
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   {tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, "ECDHE-RSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128"},
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   {tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, "ECDHE-RSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256"},
	"TLS_INVALID_CIPHER":                      {0, "", "", "", ""},
}

func NewCipherSuite(id uint16, name, key, encr, bitSize string) *TGCipherSuite {
	return &TGCipherSuite{suiteId: id, opensslName: name, keyExch: key, encryption: encr, bits: bitSize}
}

/////////////////////////////////////////////////////////////////
// Helper Public functions for TGCipherSuite
/////////////////////////////////////////////////////////////////

// GetCipherSuite returns the TGCipherSuite given its full qualified string form or its alias name.
func GetCipherSuite(nameOrAlias string) *TGCipherSuite {
	for name, suite := range PreDefinedCipherSuites {
		if 	strings.ToLower(name) == strings.ToLower(nameOrAlias) ||
			strings.ToLower(suite.opensslName) == strings.ToLower(nameOrAlias) {
			return &suite
		}
	}
	invalid := PreDefinedCipherSuites["TLS_INVALID_CIPHER"]
	return &invalid
}

// GetCipherSuiteById returns the TGCipherSuite given its ID.
func GetCipherSuiteById(id uint16) *TGCipherSuite {
	for _, suite := range PreDefinedCipherSuites {
		if 	suite.suiteId == id {
			return &suite
		}
	}
	invalid := PreDefinedCipherSuites["TLS_INVALID_CIPHER"]
	return &invalid
}

// FilterSuites returns CipherSuites that are supported by TGDB client
func FilterSuites(suites []string) []string {
	supportedSuites := make([]string, 0)
	for _, inputSuite := range suites {
		cs := GetCipherSuite(inputSuite)
		// Ignore "TLS_INVALID_CIPHER"
		if cs.suiteId != 0 {
			supportedSuites = append(supportedSuites, cs.opensslName)
		}
	}
	return supportedSuites
}

// FilterSuitesById returns CipherSuites that are supported by TGDB client
func FilterSuitesById(suites []uint16) []uint16 {
	supportedSuites := make([]uint16, 0)
	for _, inputSuite := range suites {
		cs := GetCipherSuiteById(inputSuite)
		// Ignore "TLS_INVALID_CIPHER"
		if cs.suiteId != 0 {
			supportedSuites = append(supportedSuites, cs.suiteId)
		}
	}
	return supportedSuites
}
