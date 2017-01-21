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

var LinkProtocol = require('./LinkProtocol').LinkProtocol,
    TCPChannel   = require('./tcp/TCPChannel').TCPChannel,
    TESTChannel   = require('./test/TESTChannel').TESTChannel,
    CONFIG_NAMES = require('../utils/ConfigName').CONFIG_NAMES;

function setUserAndPassword(properties, username, password) {
    properties[CONFIG_NAMES.CHANNEL_USER_ID.name] = username;

    if (password !== null && password.length !== 0) {
        properties[CONFIG_NAMES.CHANNEL_PASSWD.name] = password;
    }
}

function TGChannelFactory() {

}

/**
 *
 * @param channelURL - #TGChannelURL
 * @param username
 * @param password
 * @param properties
 * @returns {*|TCPChannel}
 */
TGChannelFactory.createChannel = function(channelURL, username, password, properties) {
    if (properties === null) {
        properties = {};
    }
    if (channelURL === undefined) {
        throw new Error('Invalid url');
    }
    setUserAndPassword(properties, username, password);
    //Get protocol
    var protocol = channelURL.getProtocol();

    switch (protocol) {
        case LinkProtocol.TCP :
            return new TCPChannel(channelURL, properties);
        case LinkProtocol.TEST :
            return new TESTChannel(channelURL, properties);
    }
};

exports.TGChannelFactory = TGChannelFactory;
