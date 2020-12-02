/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name: clientversionimpl.go
 * Created on: 11/21/2019
 * Created by: nimish
 *
 * SVN Id: $Id: clientversionimpl.tag 4146 2020-07-10 01:12:51Z nimish $
 */

package impl

import (
	"fmt"
	"strings"
)

const (
	BuildTypeProduction byte = iota
	BuildTypeEngineering
	BuildTypeBeta
)

const (
	EditionEvaluation byte = iota
	EditionCommunity
	EditionEnterprise
	EditionDeveloper
)

const (
	currentMajor  = 3
	currentMinor  = 0
	currentUpdate = 0
	currentHotFix = 0
	currentBuild  = 01
)

type TGClientVersion struct {
	major     byte
	minor     byte
	update    byte
	hotFixNo  byte
	buildNo   uint16
	buildRev  uint16
	buildType byte
	edition   byte
}

func DefaultTGClientVersion() *TGClientVersion {
	version := TGClientVersion{
		major:     VERS_MAJOR,
		minor:     VERS_MINOR,
		update:    VERS_UPDATE,
		hotFixNo:  VERS_HFNO,
		buildNo:   VERS_BUILDNO,
		buildRev:  VERS_REV,
		buildType: fromName2BuildType(VERS_BUILDTYPE_STR),
		edition:   fromName2BuildEdition(VERS_EDITION_STR),
	}
	return &version
}

func fromName2BuildType (nameIn string) byte {
	name := strings.ToLower(nameIn);
	if strings.Compare("production", name) == 0 {
		return BuildTypeProduction
	}
	if strings.Compare("engineering", name) == 0 {
		return BuildTypeEngineering
	}
	if strings.Compare("beta", name) == 0 {
		return BuildTypeBeta
	}
	return BuildTypeEngineering
}

func fromName2BuildEdition (nameIn string) byte {
	name := strings.ToLower(nameIn);

	if strings.Compare("evaluation", name) == 0 {
		return EditionEvaluation
	}
	if strings.Compare("community", name) == 0 {
		return EditionCommunity
	}
	if strings.Compare("enterprise", name) == 0 {
		return EditionEnterprise
	}
	if strings.Compare("developer", name) == 0 {
		return EditionDeveloper
	}
	return EditionCommunity
}

func NewTGClientVersion(maj, min, upd, hf byte, bld uint16, bldRev uint16, bType, edt byte) *TGClientVersion {
	version := TGClientVersion{
		major:     maj,
		minor:     min,
		update:    upd,
		hotFixNo:  hf,
		buildNo:   bld,
		buildRev:  bldRev,
		buildType: bType,
		edition:   edt,
	}
	return &version
}

func (obj *TGClientVersion) GetMajor() byte {
	return obj.major
}

func (obj *TGClientVersion) GetMinor() byte {
	return obj.minor
}

func (obj *TGClientVersion) GetUpdate() byte {
	return obj.update
}

func (obj *TGClientVersion) GetBuildNo() uint16 {
	return obj.buildNo
}

func (obj *TGClientVersion) GetBuildRevision() uint16 {
	return obj.buildRev
}

func (obj *TGClientVersion) GetEdition() byte {
	return obj.edition
}

func (obj *TGClientVersion) GetBuildType() byte {
	return obj.buildType
}

func (obj *TGClientVersion) GetVersionAsLong() int64 {
	result := int64(obj.major)
	lMinor := int64(obj.minor) << 8
	lUpdate := int64(obj.update) << 16
	lhfNo := int64(obj.hotFixNo) << 24
	lbuildNo := int64(obj.buildNo) << 40
	lbuildType := int64(obj.buildType) << 44
	lEdition := int64(obj.edition) << 48

	result |= lMinor
	result |= lUpdate
	result |= lhfNo
	result |= lbuildNo
	result |= lbuildType
	result |= lEdition

	return int64(result)
}

func GetClientVersion() *TGClientVersion {
	return DefaultTGClientVersion()
}

func (obj *TGClientVersion) GetVersionString() string {
	strVersion := fmt.Sprintf("ClientVersionInfo [Major=%d, minor=%d, update=%d, hotFix=%d, buildNo=%d, buildRevision=%d, buildType=%d, edition=%d]",
		obj.major, obj.minor, obj.update, obj.hotFixNo, obj.buildNo, obj.buildRev, obj.buildType, obj.edition)
	return strVersion
}

