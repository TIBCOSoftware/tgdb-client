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
    this.parseInternal(url);
}

TGChannelURL.prototype.parseInternal = function(url) {
    //Take protocol part out
    var userURL = this.parseProtocol(url);
    //Parse username is present
    var hostPortURL = this.parseUserInfo(userURL);
    //Extract host and port
    this.parseHostPort(hostPortURL);
};

TGChannelURL.prototype.parseProtocol = function(url) {
    var protocolIndex;

    if (url.indexOf(LinkProtocol.TCP, 0) == 0) {
        this._protocol = LinkProtocol.Type.TCP;
        protocolIndex = LinkProtocol.TCP.length;
    }
    else if (url.indexOf(LinkProtocol.TCPS) == 0) {
        this._protocol = LinkProtocol.Type.SSL;
        protocolIndex = LinkProtocol.TCPS.length;
    }
    else if (url.indexOf(LinkProtocol.HTTP) == 0) {
        this._protocol = LinkProtocol.Type.HTTP;
        protocolIndex = LinkProtocol.HTTP.length;
    }
    else if (url.indexOf(LinkProtocol.HTTPS) == 0) {
        this._protocol = LinkProtocol.Type.HTTPS;
        protocolIndex = LinkProtocol.HTTPS.length;
    }
    else if (url.indexOf(LinkProtocol.TEST) == 0) {
        this._protocol = LinkProtocol.Type.TEST;
        protocolIndex = LinkProtocol.TEST.length;
    }
    else {
        throw Error("Invalid protocol specification. for URL. URL is of the form protocol://host:value/{name=value;name=value...}");
    }
    return url.substring(protocolIndex);
};

TGChannelURL.prototype.parseUserInfo = function(userURL) {
    var atPos = userURL.indexOf('@');
    if (atPos == -1) {
        return -1;
    }
    this._username = userURL.substring(0, atPos);
    return userURL.substring(atPos + 1);
};

TGChannelURL.prototype.parseHostPort = function(hostPortURL) {
    var portString;

    if (hostPortURL.length == 0) {
        this._host = 'localhost';
        this._port = 8700;
    } else {
        var offset = 0;
        var ipv6Offset = hostPortURL.indexOf('[');
        if (ipv6Offset != -1) {
            var endIPv6Pos = hostPortURL.indexOf(']');
            if (endIPv6Pos > ipv6Offset + 2) {
                offset = endIPv6Pos + 1;
            }
            else {
                throw new Error('Invalid or missing host name');
            }
        }
        var position = hostPortURL.indexOf(':', offset);
        var lPosition = hostPortURL.indexOf("/");

        if (position < 0) {
            var noPort = true;

            if (hostPortURL.indexOf('.') < 0) {
                this._host = 'localhost';
                portString = hostPortURL;
                noPort = false;
            }
            if (noPort) {
                this._host = hostPortURL;
                portString = '8700';
            }
        }
        else {
            this._host = hostPortURL.substring(0, position);
            portString = hostPortURL.substring(position + 1, (lPosition != -1 ? lPosition : hostPortURL.length));
        }
        if (this._host.length == 0) {
            throw new Error('Invalid or missing host name');
        }
        if (portString.length == 0) {
            throw new Error('Invalid or missing port number');
        }
        this._port = parseInt(portString);
        return lPosition == -1 ? '' : hostPortURL.substring(lPosition + 1);
    }
};

TGChannelURL.prototype.getProtocol = function() {
    return this._protocol;
};

TGChannelURL.prototype.getHost = function() {
    return this._host;
};

TGChannelURL.prototype.getPort = function() {
    return this._port;
};

exports.TGChannelURL = TGChannelURL;