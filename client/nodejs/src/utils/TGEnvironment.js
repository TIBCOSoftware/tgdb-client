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

var ConfigName = require('./ConfigName').CONFIG_NAMES;

var env = {};

 function TGEnvironment () {
	 
	for(var key in ConfigName) {
		if(ConfigName.hasOwnProperty(key)) {
			var config = ConfigName[key];
		    env[config.name] = config.defaultValue;
		}
	}
    
    this.getChannelSendSize = function () {
        var value = env[ConfigName.CHANNEL_SEND_SIZE.name];
        
        if (value !== null) {
            return value.valueOf();
        }
        return 122;
    };

    this.getChannelReceiveSize = function () {
        var value = env[ConfigName.CHANNEL_RECV_SIZE.name];
        if (value !== null) {
            return value.valueOf();
        }
        return 128;

    };
    
    this.getChannelPingInterval = function () {
        var value = env[ConfigName.CHANNEL_PING_INTERVAL.name];
        if (value !== null) {
            return value.valueOf();
        }
        return 30;

    };
    
    this.getChannelConnectTimeout = function () {
        var value = env[ConfigName.CHANNEL_CONNECTION_TIMEOUT.name];
        if (value !== null) {
            return value.valueOf();
        }
        return 30;

    };

    this.getChannelFTHosts = function () {
        return env[ConfigName.CHANNEL_FT_HOSTS.name];
    };

    this.getChannelDefaultUser = function () {
        return env[ConfigName.CHANNEL_DEFAULT_USER_ID.name];
    };

    /**
     * Get the Connection User as specified in the Environment
     * @return the Channel User
     */
    this.getChannelUser = function () {
        return env[ConfigName.CHANNEL_USER_ID.name];
    };

    this.setProperty = function (name, value) {
        var cn = ConfigName.fromName(name);
        if (cn === ConfigName.INVALID_NAME) {
        	return;
        }
        env.put(cn.name, value);
    };

    this.getProperty = function (name, defaultValue) {
    	defaultValue = (!defaultValue) ? null : defaultValue;
        var cn = ConfigName.fromName(name);
        if (cn === ConfigName.INVALID_NAME) {
        	return defaultValue;
        }

        var value = env[cn.name];
        return value === null ? defaultValue : value;
    };

    this.getChannelDefaultPort = function () {
        var value = env[ConfigName.CHANNEL_DEFAULT_PORT.name];

        return value.valueOf();
    };

    this.getChannelDefaultHost = function () {
        return env[ConfigName.CHANNEL_DEFAULT_HOST.name];
    };

    this.getConnectionPoolDefaultPoolSize = function () {
        var value = env[ConfigName.CONN_POOL_DEFAULT_SIZE.name];
        return value.valueOf();
    };

    this.getChannelClientId = function () {
        var value = env[ConfigName.CHANNEL_CLIENTID.name];
        if (value === null) {
            return "tgdb.java-api.client";
        }
        return value;
    };
}

exports.TGEnvironment = new TGEnvironment();

function testStub () {
	var tge = new TGEnvironment();
	
	console.log('getChannelSendSize:' + tge.getChannelSendSize());

	console.log('getChannelReceiveSize' + tge.getChannelReceiveSize());
    
	console.log('getChannelPingInterval' + tge.getChannelPingInterval());
    
	console.log('getChannelConnectTimeout' + tge.getChannelConnectTimeout());

	console.log('getChannelFTHosts' + tge.getChannelFTHosts());

	console.log('getChannelDefaultUser' + tge.getChannelDefaultUser());

	console.log('getChannelUser' + tge.getChannelUser());

	tge.setProperty('tgdb.a.b', 'value');

	console.log('getProperty(name)' + tge.getProperty('tgdb.a.b'));

	console.log('getProperty(name, defaultValue)' + tge.getProperty('tgdb.a.c', 'defaultValue'));

	console.log('getChannelDefaultPort' + tge.getChannelDefaultPort());

	console.log('getChannelDefaultHost' + tge.getChannelDefaultHost());

	console.log('getConnectionPoolDefaultPoolSize' + tge.getConnectionPoolDefaultPoolSize());

	console.log('getChannelClientId' + tge.getChannelClientId());

}

//testStub();
