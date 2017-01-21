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

var inherits      = require('util').inherits,
    LinkProtocol  = require('../LinkProtocol').LinkProtocol,
    TGChannelURL  = require('../TGChannelURL').TGChannelURL,
    TGException   = require('../../exception/TGException').TGException,
    TGEnvironment = require('../../utils/TGEnvironment').TGEnvironment,
    CONFIG_NAMES  = require('../../utils/ConfigName').CONFIG_NAMES,
    TGLogManager  = require('../../log/TGLogManager'),
    TGLogLevel    = require('../../log/TGLogger').TGLogLevel;

var logger = TGLogManager.getLogger();
var defaultHost = "localhost";
var defaultPort = 8222;
var defaultProtocol = LinkProtocol.TCP;

function TCPChannelURL (tgChannelURL) {
	TCPChannelURL.super_.call(this, tgChannelURL.getUrl());

    this._user     = null;
    this._host     = null;
    this._portStr  = null;
    this._port     = 0;
    this._isIPv6   = false;
    this._props    = {};
    this._ftUrls   = null;
    
    this._protocol = tgChannelURL.getProtocol();
    this._url      = tgChannelURL.getUrl();
    
    this.parseInternal(this._url.toLowerCase().trim());
}

inherits(TCPChannelURL, TGChannelURL);

TCPChannelURL.prototype.parseInternal = function (url) {
	
    if (!url || url.length === 0) {
    	url = defaultProtocol.prefix.concat(defaultHost).concat(defaultPort);
    }
	
	logger.logDebugWire('[TCPChannelURL.prototype.parseInternal] url : %s', url);
	
    var hostAndPort = this.parseUser(this.stripProtocol(url));
    logger.logDebugWire('[TCPChannelURL.prototype.parseInternal] hostAndPort : %s', hostAndPort);
        
    var properties = this.parseHostAndPort(hostAndPort);
    logger.logDebugWire('[TCPChannelURL.prototype.parseInternal] properties : %s', properties);
};

TCPChannelURL.prototype.parseUser = function (url) {
    logger.logDebugWire('[TCPChannelURL.prototype.parseUser] url : ' + url);
    var posAt = url.indexOf('@');
    if (posAt === -1) {
        return url;
    }
        
    this._user = url.substr(0, posAt);
    
    logger.logDebugWire('[TCPChannelURL.prototype.parseUser] User : ' + this._user);
    
    this._props[CONFIG_NAMES.CHANNEL_USER_ID.name] = this._user;
    return url.substr(posAt+1);
};

TCPChannelURL.prototype.parseHostAndPort = function (hostAndPort) {
    if (hostAndPort.length === 0) {
        this._host = "localhost";
        this._port = 8700;  //default values
        return null;
    }

    var offset = 0;
    var posIPv6 = hostAndPort.indexOf('[');
    if (posIPv6 !== -1) {
        var endIPv6 = hostAndPort.indexOf(']');
        if (endIPv6 > posIPv6 + 2) {
            offset = endIPv6 + 1;
            this.isIPv6 = true;
        } else {
            throw new TGException("Invalid or missing host name");
        }
    }

    var pos = 0;
    if (this.isIPv6) {
        pos = hostAndPort.lastIndexOf(':');
    } else {
        pos = hostAndPort.indexOf(':', offset);
    }
    
    var lpos = hostAndPort.indexOf("/");
    if (pos < 0) {
        var noPort = true;

        if (hostAndPort.indexOf('.') < 0) {
            this._host = this._gDefaultHost;
            this._portStr = hostAndPort;
            noPort = false;
        }

        if (noPort) {
            this._host = hostAndPort;
            this._portStr = this._gDefaultPort.toString();
        }
    } else {
            var startpos = 0;
            var endpos = lpos !== -1 ? lpos : hostAndPort.length;
            if (this._isIPv6) {
                startpos = 1;
                endpos = offset - 1;
            }

            this._host = hostAndPort.substr(startpos, pos);
            this._portStr = hostAndPort.substr(pos+1, endpos);
        }

        if (this._host.length === 0) {
            throw new TGException("Invalid or missing host name");
        }
        if (this._portStr.length === 0) {
            throw new TGException("Invalid or missing port number");
        }


        this._port = this._portStr.valueOf();

        return lpos === -1 ? "" : hostAndPort.substr(lpos+1);
};

TCPChannelURL.prototype.stripProtocol = function () {
    return this._url.substr(this._protocol.prefix.length);
};

TCPChannelURL.prototype.parseProperties = function (properties) {
    if (properties.length === 0) {
    	return properties;
    }
        
    if (!((properties.startsWith("{")) && (properties.endsWith("}")))) {
        throw new TGException("Malformed URL property specification - Must begin with { and end with }. All key=value must be seperated with ;");
    }

    var kvStr = properties.substr(1, properties.length-1);

    var convtBuf = [];
    var limit;
    var keyLen;
    var valueStart;
    var hasSep;
    var precedingBackslash;
    var scr = new SemiColonReader(kvStr.split(''));
    var c;
        
    while ((limit = scr.readSemiColon()) > 0) {
        keyLen = 0;
        valueStart = limit;
        hasSep = false;

        precedingBackslash = false;
        while (keyLen < limit) {
            c = scr.kvBuf[keyLen];
            //need check if escaped.
            if ((c === '=' || c === ':') && !precedingBackslash) {
                valueStart = keyLen + 1;
                hasSep = true;
                break;
            } else if ((c === ' ' || c === '\t' || c === '\f') && !precedingBackslash) {
                valueStart = keyLen + 1;
                break;
            }
            if (c === '\\') {
                precedingBackslash = !precedingBackslash;
            } else {
                precedingBackslash = false;
            }
            keyLen++;
        }
        while (valueStart < limit) {
            c = scr.kvBuf[valueStart];
            if (c !== ' ' && c !== '\t' && c !== '\f') {
                if (!hasSep && (c === '=' || c === ':')) {
                    hasSep = true;
                } else {
                    break;
                }
            }
            valueStart++;
        }
        var key = scr.kvBuf.slice(0, keyLen).join('');
        var value = scr.kvBuf(valueStart, limit - valueStart).join('');
        this._props[key] = value;
    }

    return properties;
};

TCPChannelURL.prototype.getProtocol = function () {
    return this._protocol;
};

TCPChannelURL.prototype.getHost = function () {
        if(this._host !== null){
        	return this._host;
        }
        
        this._host = TGEnvironment.getChannelDefaultHost();

        return this._host;
};

TCPChannelURL.prototype.getPort = function () {
    if (this._port !== -1) {
    	return this._port;
    }

    this._port = TGEnvironment.getChannelDefaultPort();
    return this._port;
};

TCPChannelURL.prototype.getProperties = function () {
    return this._props;
};

TCPChannelURL.prototype.getUser = function () {
    if (this._user !== null) {
    	return this._user;
    }

    if ((this._user === null) || (this._user.length === 0)) {
        this._user = this._props[CONFIG_NAMES.CHANNEL_USER_ID.alias];

        if ((this._user === null) || (this._user.length === 0)) {
            this._user = this._props[CONFIG_NAMES.CHANNEL_USER_ID.name];

            if ((this._user === null) || (this._user.length === 0)) {
                this._user = TGEnvironment.getChannelUser();

                if ((this._user === null) || (this._user.length === 0)) {
                    	this._user = TGEnvironment.getChannelDefaultUser();
                }
            }
        }
    }
    return this._user;
};
/*
TCPChannelURL.prototype.getFTUrls = function () {
    if (this._ftUrls !== null) {
    	return this._ftUrls;
    }
    
    try {
        var ftHosts = this._props.get(CONFIG_NAMES.CHANNEL_FT_HOSTS.alias);
        if ((ftHosts === null) || (ftHosts.length === 0)) {
            ftHosts = this._props.get(CONFIG_NAMES.CHANNEL_FT_HOSTS.name);
        }

        if ((ftHosts === null) || (ftHosts.length === 0)) {
        	var emptyList = [];
        	return emptyList;
        }
        
        var ftHostStrs = ftHosts.trim().split(",");
        var ftUrls = [];
        for(var index in ftHostStrs) {
        	LinkURL url = new LinkURL(this.protocol, this.host, this.port);
            url.parseHostAndPort(ftHost);
            url.user = this.user;
            url.props = this.props;
            ftUrls.add(url);
        }
        this.ftUrls = ftUrls;
    }
    catch (error) {
        return Collections.emptyList();
    }

    return ftUrls;
};
*/

exports.TCPChannelURL = TCPChannelURL;

TCPChannelURL.prototype.toString = function () {
	return "user:".concat(this.getUser()).concat(", ").
	       concat("protocol:").concat(this.getProtocol()).concat(", ").
	       concat("host:").concat(this.getHost()).concat(", ").
	       concat("port:").concat(this.getPort()).concat(", ").
	       concat("props:").concat(this.getProperties());
};

function SemiColonReader (inBuf) {

    this._curPos = 0;
    this._kvBuf = null;

    this.getKvBuf = function() {
    	return this._kvBuf;
    };
    
    this.readSemiColon = function () {
        var kvLen = 0;
        this._kvBuf = [];
        var precedingBackslash = false;
        while (true) {
            if (this._curPos >= inBuf.length){
                return kvLen;
            }

            var c = inBuf[this._curPos++];

            switch (c) {
                case ' ':
                case '\r':
                case '\n':
                case '\t':
                case '\f':
                    break;
                case '\\':
                    precedingBackslash = true;
                    break;
                case ';':
                    if (precedingBackslash) {
                    	this._kvBuf.push(c);
                        ++kvLen;
                        precedingBackslash = false;
                        break;
                    }
                    else {
                        return kvLen;  //come out of the loop
                    }

                default:
                	this._kvBuf.push(c);
                    precedingBackslash = false;
                    ++kvLen;
                    break;
             }

        }

    };
}