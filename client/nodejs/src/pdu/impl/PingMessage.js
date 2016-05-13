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

var util                    = require('util'),
    AbstractProtocolMessage = require('../AbstractProtocolMessage').AbstractProtocolMessage,
    VerbId                  = require('./VerbId').VerbId;

//Class Definition
function PingMessage() {
    AbstractProtocolMessage.call();
}

util.inherits(PingMessage, AbstractProtocolMessage);

PingMessage.prototype.getVerbId = function() {
    return VerbId.PING_MESSAGE;
};

PingMessage.prototype.writePayload = function(dynamicBuffer) {

};

PingMessage.prototype.readPayload = function() {

};

PingMessage.prototype.isUpdateable = function() {
    return false;
};

exports.PingMessage = PingMessage;