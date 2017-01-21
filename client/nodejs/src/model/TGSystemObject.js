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

//var TGAttributeDescriptor = require('./TGAttributeDescriptor').TGAttributeDescriptor, 
//    TGProtocolVersion = require('../../TGProtocolVersion').TGProtocolVersion;

function TGSystemObject() {
}

TGSystemObject.TGSystemType = {
	InvalidType         : { value : -1},
	AttributeDescriptor : { value : 0},
	NodeType            : { value : 1},
	EdgeType            : { value : 2},
	Index               : { value : 3},
	Prinicapl           : { value : 4},
	Role                : { value : 5},
	Sequence            : { value : 6},
	MaxSysObjectTypes   : { value : 7},
	fromValue : function (type) {
		switch (type) {
			case this.AttributeDescriptor.value: return this.AttributeDescriptor;
			case this.NodeType.value: return this.NodeType;
			case this.EdgeType.value: return this.EdgeType;
			case this.Index.value: return this.Index;
			case this.Prinicapl.value: return this.Prinicapl;
			case this.Role.value: return this.Role;
			case this.Sequence.value: return this.Sequence;
			case this.MaxSysObjectTypes.value: return this.MaxSysObjectTypes;
		}
		return this.InvalidType;
	}
};

exports.TGSystemObject = TGSystemObject;