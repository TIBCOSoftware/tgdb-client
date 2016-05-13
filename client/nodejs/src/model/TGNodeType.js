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

/**
 *
 * @param graphMetadata - #TGGraphMetadata
 * @param name
 * @param parentNodeType - #TGNodeType
 * @constructor
 */
//Class definition
function TGNodeType(graphMetadata, name, parentNodeType) {
    this._graphMetadata  = graphMetadata;
    this._name           = name;
    this._parentNodeType = parentNodeType;
}

TGNodeType.prototype.getTypeName = function() {
    return this._name;
};


TGNodeType.prototype.getDerivedFrom = function() {
    return this._parentNodeType;
};


TGNodeType.prototype.getEntityKind = function() {
    return TGAbstractEntity.EntityKind.NODETYPE;
};