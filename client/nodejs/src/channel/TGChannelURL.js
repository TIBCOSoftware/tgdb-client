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

var LinkProtocol = require('./LinkProtocol').LinkProtocol;

function TGChannelURL(url) {
    this._url = url;
    if (url.indexOf(LinkProtocol.TCP.prefix, 0) === 0) {
        this._protocol = LinkProtocol.TCP;
    }
    else if (url.indexOf(LinkProtocol.TCPS.prefix) === 0) {
        this._protocol = LinkProtocol.SSL;
    }
    else if (url.indexOf(LinkProtocol.HTTP.prefix) === 0) {
        this._protocol = LinkProtocol.HTTP;
    }
    else if (url.indexOf(LinkProtocol.HTTPS.prefix) === 0) {
        this._protocol = LinkProtocol.HTTPS;
    }
    else if (url.indexOf(LinkProtocol.TEST.prefix) === 0) {
        this._protocol = LinkProtocol.TEST;
    }
    else {
        throw Error("Invalid protocol specification. for URL. URL is of the form protocol://host:value/{name=value;name=value...}");
    }
}

TGChannelURL.prototype.getUrl = function() {
    return this._url;
};

TGChannelURL.prototype.getProtocol = function() {
    return this._protocol;
};

exports.TGChannelURL = TGChannelURL;