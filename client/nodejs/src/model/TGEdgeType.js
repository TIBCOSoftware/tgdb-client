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

var TGAbstractEntity = require('./TGAbstractEntity').TGAbstractEntity,
    util             = require('util');


/**
 *
 * @param graphMetadata - #TGGraphMetadata
 * @param name
 * @param directionType - #TGEdge#DirectionType
 * @param parentType - #TGEdgeType
 * @constructor
 */
//Class definition
function TGEdgeType(graphMetadata, name, directionType, parentType) {
    TGAbstractEntity.call(graphMetadata);
    this._graphMetadata  = graphMetadata;
    this._name           = name;
    this._directionType  = directionType;
    this._parentType     = parentType;
}

util.inherits(TGEdgeType, TGAbstractEntity);

TGEdgeType.prototype.getDirectionType = function() {
    return this._directionType;
};


TGEdgeType.prototype.getTypeName = function() {
    return this._name;
};


TGEdgeType.prototype.getDerivedFrom = function() {
    return this._parentType;
};


TGEdgeType.prototype.getEntityKind = function() {
    return TGAbstractEntity.EntityKind.EDGETYPE;
};
