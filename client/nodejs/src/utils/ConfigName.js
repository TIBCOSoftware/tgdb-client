/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not
 * use this file except in compliance with the License. A copy of the License is
 * included in the distribution package with this file. You also may obtain a
 * copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

var CONFIG_NAMES = {
	CHANNEL_DEFAULT_HOST : {
		name : 'tgdb.channel.defaultHost',
		alias : null,
		defaultValue : 'localhost'
	},
	CHANNEL_DEFAULT_PORT : {
		name : 'tgdb.channel.defaultPort',
		alias : null,
		defaultValue : '8700'
	},
	CHANNEL_DEFAULT_PROTOCOL : {
		name : 'tgdb.channel.defaultProtocol',
		alias : null,
		defaultValue : 'tcp'
	},
	CHANNEL_SEND_SIZE : {
		name : 'tgdb.channel.sendSize',
		alias : 'sendSize',
		defaultValue : '122'
	},
	CHANNEL_RECV_SIZE : {
		name : 'tgdb.channel.recvSize',
		alias : 'recvSize',
		defaultValue : '128'
	},
	USE_DEDICATED_CHN_PER_CONN : {
		name : 'tgdb.connectionpool.useDedicatedChannelPerConnection',
		alias : 'useDedicatedChannelPerConnection',
		defaultValue : 'false'
	},
	CONN_POOL_DEFAULT_SIZE : {
		name : 'tgdb.connectionpool.defaultPoolSize',
		alias : 'defaultPoolSize',
		defaultValue : '10'
	},
	
    CHANNEL_PING_INTERVAL : {
        name : 'tgdb.channel.pingInterval',
        alias : 'pingInterval',
        defaultValue : '30'
    },
    CHANNEL_CONNECTION_TIMEOUT : {
        name : 'tgdb.channel.connectTimeout',
        alias : 'connectTimeout',
        defaultValue : '30'
    },
	
	CHANNEL_FT_HOSTS : {
		name : 'tgdb.channel.FTHosts',
		alias : 'ftHosts',
		defaultValue : 'null'
	},
	CHANNEL_DEFAULT_USER_ID : {
		name : 'tgdb.channel.defaultUserID',
		alias : 'null',
		defaultValue : null
	},
	CHANNEL_USER_ID : {
		name : 'tgdb.channel.userID',
		alias : 'userID',
		defaultValue : null
	},
	CHANNEL_PASSWD : {
		name : 'tgdb.channel.password',
		alias : 'password',
		defaultValue : null
	},
	CHANNEL_CLIENTID : {
		name : 'tgdb.channel.clientId',
		alias : 'clientId',
		defaultValue : 'tgdb.nodejs-api.client'
	},
	CONNECTION_OPERATION_TIMEOUT : {
		name : 'tgdb.connection.operationTimeout',
		alias : null,
		defaultValue : '10000'
	},
	INVALID_NAME : {
		name : null,
		alias : null,
		defaultValue : null
	},
	fromName : function (name) {
		var config = this.INVALID_NAME;
		for(var key in CONFIG_NAMES) {
			if(CONFIG_NAMES.hasOwnProperty(key)) {
				if(CONFIG_NAMES[key].name === name.toLowerCase()) {
					config = CONFIG_NAMES[key];
					break;
				}
				
				if(!CONFIG_NAMES[key].alias) {
					break;
				}
				
				if(CONFIG_NAMES[key].alias.toLowerCase() === name.toLowerCase()) {
					config =  CONFIG_NAMES[key];
					break;
				}
			}
		}
		return this.INVALID_NAME;
	}
};

exports.CONFIG_NAMES = CONFIG_NAMES;
