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
var QUERY_OPTION_FETCHSIZE = "fetchsize";
var QUERY_OPTION_TRAVERSALDEPTH = "traversaldepth";
var QUERY_OPTION_EDGELIMIT = "edgelimit";

function TGQueryOptionImpl (mutable) {
    this._mutable = mutable;
    this._properties = {};
    this._properties[QUERY_OPTION_FETCHSIZE] = -1;
    this._properties[QUERY_OPTION_TRAVERSALDEPTH] = -1;
    this._properties[QUERY_OPTION_EDGELIMIT] = -1;
}

TGQueryOptionImpl.prototype.setPrefetchSize = function(size) {
    if (!this._mutable) {
    	throw new Error("Can't modify a immutable Option");
    }
    if (size === 0) {
    	size = -1;
    }
    this._properties[QUERY_OPTION_FETCHSIZE] = size;
};

TGQueryOptionImpl.prototype.getPrefetchSize = function() {
    return this._properties[QUERY_OPTION_FETCHSIZE];
};

TGQueryOptionImpl.prototype.setTraversalDepth = function(depth) {
    if (!this._mutable) {
        throw new Error("Can't modify a immutable Option");
    }
    if (depth === 0) {
        depth = -1;
    }
    this._properties[QUERY_OPTION_TRAVERSALDEPTH] = depth;
};

TGQueryOptionImpl.prototype.getTraversalDepth = function() {
    return this._properties[QUERY_OPTION_TRAVERSALDEPTH];
};

TGQueryOptionImpl.prototype.setEdgeLimit = function(limit) {
    if (!this._mutable) {
    	throw new Error("Can't modify a immutable Option");
    }
    if (limit === 0) {
    	limit = -1;
    }
    this._properties[QUERY_OPTION_EDGELIMIT] = limit;
};

TGQueryOptionImpl.prototype.getEdgeLimit = function() {
    return this._properties[QUERY_OPTION_EDGELIMIT];
};

exports.TGQueryOptionImpl = TGQueryOptionImpl;

function test() {
	var option = new TGQueryOptionImpl(false);
}

//test();