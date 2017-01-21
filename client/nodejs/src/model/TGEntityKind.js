/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 */

exports.TGEntityKind = {
	INVALIDKIND : {value : 0, name : 'Invalid'},
	ENTITY      : {value : 1, name : 'Entity'},
	NODE        : {value : 2, name : 'Node'},
	EDGE        : {value : 3, name : 'Edge'},
	GRAPH       : {value : 4, name : 'Graph'},
	HYPEREDGE   : {value : 5, name : 'HyperEdge'},
	isNode      : function (kind) {
		return kind.value === 2;
	}, 
	isEdge      : function (kind) {
		return kind.value === 3;
	},
	fromValue   : function (kindValue) {
		switch(kindValue) {
			case 0 : return this.INVALIDKIND;
			case 1 : return this.ENTITY;
			case 2 : return this.NODE;
			case 3 : return this.EDGE;
			case 4 : return this.GRAPH;
			case 5 : return this.HYPEREDGE;
		}
	}
};