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

var inherits     = require('util').inherits,
    TGChannel    = require('../TGChannel').TGChannel,
    TGProperties = require('../../utils/TGProperties').TGProperties;

//Class definition
/**
 *
 * @param serverURL - #TGChannelURL
 * @param properties - Dictionary for configuration
 * @constructor
 */
function TESTChannel(serverURL, properties) {
    //AbstractChannel.call(serverURL, properties);
    this._serverURL  = serverURL;
    this._properties = new TGProperties(properties);
}

inherits(TESTChannel, TGChannel);

/**
 * Connect using TEST.
 * This operation is asynchronous, upon completion
 * will invoke the callback specified in the function.
 */
TESTChannel.prototype.makeConnection = function(callback) {
    var options = {
        host: this._serverURL.getHost(),
        port: this._serverURL.getPort()
    };
    //var eventEmitter = new TCPConnectionEventEmitter(this);
    var channel = this;
    ////console.log("[TESTChannel::makeConnection] here I am ........ ");
};

exports.TESTChannel = TESTChannel;