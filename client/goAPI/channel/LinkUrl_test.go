package channel

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
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
 * File name: LinkUrl_Test.go
 * Created on: Nov 24, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const testUrl0 = "foo1.bar.com"
const testUrl1 = "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
const testUrl2 = "tcp://scott@foo.bar.com:8700"
const testUrl3 = "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
const testUrl4 = "http://[2001:db8:1f70::999:de8:7648:6e8]:100/{userID=Admin}"

func TestGetFTUrls(t *testing.T) {
	linkUrl := NewLinkUrl(testUrl1)

	testFtUrls := linkUrl.GetFTUrls()
	if len(testFtUrls) > 0 {
		for _, testFtUrl := range testFtUrls {
			t.Logf("Test FT Urls retrieved from linkUrl '%s' are: '%+v'", testUrl1, testFtUrl)
		}
	}
}

func TestGetProperties(t *testing.T) {
	linkUrl := NewLinkUrl(testUrl3)

	urlProperties := linkUrl.GetProperties()
	nvPairs := urlProperties.(*utils.SortedProperties).GetAllProperties()
	for _, nvPair := range nvPairs {
		t.Logf("Test FT Urls retrieved from linkUrl '%s' are: '%+v'", testUrl1, nvPair)
	}
}

func TestGetProtocol(t *testing.T) {
	linkUrl := NewLinkUrl(testUrl3)

	urlProtocol := linkUrl.GetProtocol()
	t.Logf("Test URL protocol retrieved from linkUrl '%s' are: '%+v'", testUrl3, urlProtocol)
}

func TestGetUrlAsString(t *testing.T) {
	linkUrl := NewLinkUrl(testUrl2)

	urlString := linkUrl.GetUrlAsString()
	t.Logf("Test URL string retrieved from linkUrl '%s' are: '%+v'", testUrl2, urlString)
}
