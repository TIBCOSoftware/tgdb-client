/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not
 * use this file except in compliance with the License. A copy of the License is
 * included in the distribution package with this file. You also may obtain a
 * copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

function TGProtocolVersion() {
}

TGProtocolVersion.MAJOR_VERSION = 1;
TGProtocolVersion.MINOR_VERSION = 0;
// Magic number static
TGProtocolVersion.TG_MAGIC = 0xdb2d1e4;

TGProtocolVersion.getProtocolVersion = function() {
	var version = TGProtocolVersion.MAJOR_VERSION;
	version = version << 8;
	version = version + TGProtocolVersion.MINOR_VERSION;
	return version;
};

TGProtocolVersion.getMagic = function() {
	return TGProtocolVersion.TG_MAGIC;
};

TGProtocolVersion.isCompatible = function(protocolVersion) {
	return protocolVersion == TGProtocolVersion.getProtocolVersion();
};

exports.TGProtocolVersion = TGProtocolVersion;
