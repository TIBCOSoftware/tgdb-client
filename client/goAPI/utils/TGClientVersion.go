package utils

import "fmt"

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
 * File name: TGClientVersion.go
 * Created on: Feb 20, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	buildTypeProduction byte = iota
	buildTypeEngineering
	buildTypeBeta
)

const (
	editionEvaluation byte = iota
	editionCommunity
	editionEnterprise
	editionDeveloper
)

const (
	currentMajor  = 2
	currentMinor  = 0
	currentUpdate = 1
	currentHotFix = 0
	currentBuild  = 011
)

type TGClientVersion struct {
	major     byte
	minor     byte
	update    byte
	hotFixNo  byte
	buildNo   uint16
	buildType byte
	edition   byte
}

func DefaultTGClientVersion() *TGClientVersion {
	version := TGClientVersion{
		major:     currentMajor,
		minor:     currentMinor,
		update:    currentUpdate,
		hotFixNo:  currentHotFix,
		buildNo:   currentBuild,
		buildType: buildTypeProduction,
		edition:   editionCommunity,
	}
	return &version
}

func NewTGClientVersion(maj, min, upd, hf byte, bld uint16, bType, edt byte) *TGClientVersion {
	version := TGClientVersion{
		major:     maj,
		minor:     min,
		update:    upd,
		hotFixNo:  hf,
		buildNo:   bld,
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
	strVersion := fmt.Sprintf("ClientVersionInfo [Major=%d, minor=%d, update=%d, hotFix=%d, buildNo=%d, buildType=%d, edition=%d]",
		obj.major, obj.minor, obj.update, obj.hotFixNo, obj.buildNo, obj.buildType, obj.edition)
	return strVersion
}
