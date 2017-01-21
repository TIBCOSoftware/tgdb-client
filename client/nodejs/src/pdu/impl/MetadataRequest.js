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
var util            = require('util'),
    VerbId          = require('./VerbId').VerbId,
    CallableRequest = require('./CallableRequest').CallableRequest;

function MetadataRequest (callback) {
	MetadataRequest.super_.call(this, callback);
}

util.inherits(MetadataRequest, CallableRequest);

MetadataRequest.prototype.getVerbId = function () {
    return VerbId.METADATA_REQUEST;
};

MetadataRequest.prototype.isUpdateable = function () {
    return true;
};

MetadataRequest.prototype.writePayload = function (os) {
};

MetadataRequest.prototype.readPayload = function (is) {	
};

exports.MetadataRequest = MetadataRequest;