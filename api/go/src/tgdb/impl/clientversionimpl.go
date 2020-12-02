/* Copyright (C) 1991-2012 Free Software Foundation, Inc.
   This file is part of the GNU C Library.

   The GNU C Library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   The GNU C Library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with the GNU C Library; if not, see
   <http://www.gnu.org/licenses/>.  */


/* This header is separate from features.h so that the compiler can
   include it implicitly at the start of every compilation.  It must
   not itself include <features.h> or any other header that includes
   <features.h> because the implicit include comes before any feature
   test macros that may be defined in a source file before it first
   explicitly includes a system header.  GCC knows the name of this
   header in order to preinclude it.  */

/* We do support the IEC 559 math functionality, real and complex.  */

/* wchar_t uses ISO/IEC 10646 (2nd ed., published 2011-03-15) /
   Unicode 6.0.  */

/* We do not support C11 <threads.h>.  */

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
 * SVN Id: $Id: clientversionimpl.go 4710 2020-11-13 18:17:30Z relbuild $
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
		major:     3,
		minor:     0,
		update:    0,
		hotFixNo:  0,
		buildNo:   39,
		buildRev:  4709,
		buildType: fromName2BuildType("Production"),
		edition:   fromName2BuildEdition("Enterprise"),
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

