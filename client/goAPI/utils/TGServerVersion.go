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
 * File name: TGServerVersion.go
 * Created on: Feb 20, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**
package main

import (
	"fmt"
)

type TGVersion struct {
	major     byte
	minor     byte
	update    byte
	hotFixNo  byte
	buildNo   uint16
	buildType byte
	edition   byte
}

func main() {
	fmt.Println("Hello, playground")
	obj := TGVersion{2,0,1,0,011,0,1}
	result := int64(obj.major)
	lMinor := int64(obj.minor) << 8
	lUpdate := int64(obj.update) << 16
	lhfNo := int64(obj.hotFixNo) << 24
	lbuildNo := int64(obj.buildNo) << 40
	lbuildType := int64(obj.buildType) << 44
	lEdition := int64(obj.edition) << 48

	result |= lMinor
	//	fmt.Printf("Result: %+v\n", result)
	result |= lUpdate
	//	fmt.Printf("Result1: %+v\n", result)
	result |= lhfNo
	//	fmt.Printf("Result2: %+v\n", result)
	result |= lbuildNo
	//	fmt.Printf("Result3: %+v\n", result)
	result |= lbuildType
	//	fmt.Printf("Result4: %+v\n", result)
	result |= lEdition
	fmt.Printf("Result5: %+v\n", result)

	result = 291370581426178

	maj := byte(result & 0xff)
	fmt.Printf("maj: %+v\n", maj)
	min := byte((result & 0xff00) >> 8)
	fmt.Printf("min: %+v\n", min)
	upd := byte((result & 0xff0000) >> 16)
	fmt.Printf("upd: %+v\n", upd)
	hf := byte((result & 0xff000000) >> 24)
	fmt.Printf("hf: %+v\n", hf)
	bld := uint16((result & 0xff0000000000) >> 40)
	fmt.Printf("bld: %+v\n", bld)
	bType := byte((result & 0x0f00000000000) >> 44)
	fmt.Printf("bType: %+v\n", bType)
	edt := byte((result & 0xf000000000000) >> 48)
	fmt.Printf("edt: %+v\n", edt)
	//unu := byte(int64((result & 0xff00000000000000)) >> 56)
	//fmt.Printf("unu: %+v\n", unu)
}
*/

type TGServerVersion struct {
	lVersion  int64
	major     byte
	minor     byte
	update    byte
	hotFixNo  byte
	buildNo   uint16
	buildType byte
	edition   byte
	unused    byte
}

func DefaultTGServerVersion() *TGServerVersion {
	version := TGServerVersion{
		major:     currentMajor,
		minor:     currentMinor,
		update:    currentUpdate,
		hotFixNo:  currentHotFix,
		buildNo:   currentBuild,
		buildType: buildTypeProduction,
		edition:   editionCommunity,
		unused:    editionCommunity,
	}
	return &version
}

func NewTGServerVersion(ver int64) *TGServerVersion {
	version := TGServerVersion{
		lVersion: ver,
	}
	version.setVersionComponents()
	return &version
}

func (obj *TGServerVersion) GetServerVersion() int64 {
	return obj.lVersion
}

func (obj *TGServerVersion) GetMajor() byte {
	return obj.major
}

func (obj *TGServerVersion) GetMinor() byte {
	return obj.minor
}

func (obj *TGServerVersion) GetUpdate() byte {
	return obj.update
}

func (obj *TGServerVersion) GetHotFixNo() byte {
	return obj.hotFixNo
}

func (obj *TGServerVersion) GetBuildNo() uint16 {
	return obj.buildNo
}

func (obj *TGServerVersion) GetBuildType() byte {
	return obj.buildType
}

func (obj *TGServerVersion) GetEdition() byte {
	return obj.edition
}

func (obj *TGServerVersion) GetUnused() byte {
	return obj.unused
}

func (obj *TGServerVersion) setVersionComponents() {
	obj.major = byte(obj.lVersion & 0xff)
	obj.minor = byte((obj.lVersion & 0xff00) >> 8)
	obj.update = byte((obj.lVersion & 0xff0000) >> 16)
	obj.hotFixNo = byte((obj.lVersion & 0xff000000) >> 24)
	obj.buildNo = uint16((obj.lVersion & 0xff0000000000) >> 40)
	obj.buildType = byte((obj.lVersion & 0x0f00000000000) >> 44)
	obj.edition = byte((obj.lVersion & 0xf000000000000) >> 48)
	//obj.unused = byte((obj.lVersion & 0xff00000000000000) >> 56)
}

func (obj *TGServerVersion) GetVersionString() string {
	strVersion := fmt.Sprintf("ServerVersionInfo [version=%d, major=%d, minor=%d, update=%d, hfNo=%d, buildNo=%d, buildType=%d, edition=%d, unused=%d]",
		obj.lVersion, obj.major, obj.minor, obj.update, obj.hotFixNo, obj.buildNo, obj.buildType, obj.edition, obj.unused)
	return strVersion
}
