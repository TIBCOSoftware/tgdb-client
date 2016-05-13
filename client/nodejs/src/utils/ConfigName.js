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
		name : 'sgdb.channel.defaultHost',
		alias : null,
		defaultValue : 'localhost'
	},
	CHANNEL_DEFAULT_PORT : {
		name : 'sgdb.channel.defaultPort',
		alias : null,
		defaultValue : '8700'
	},
	CHANNEL_DEFAULT_PROTOCOL : {
		name : 'sgdb.channel.defaultProtocol',
		alias : null,
		defaultValue : 'tcp'
	},
	CHANNEL_SEND_SIZE : {
		name : 'sgdb.channel.sendSize',
		alias : 'sendSize',
		defaultValue : '122'
	},
	CHANNEL_RECV_SIZE : {
		name : 'sgdb.channel.recvSize',
		alias : 'recvSize',
		defaultValue : '128'
	},
	USE_DEDICATED_CHN_PER_CONN : {
		name : 'sgdb.connectionpool.useDedicatedChannelPerConnection',
		alias : 'useDedicatedChannelPerConnection',
		defaultValue : 'false'
	},
	CONN_POOL_DEFAULT_SIZE : {
		name : 'sgdb.connectionpool.defaultPoolSize',
		alias : 'defaultPoolSize',
		defaultValue : '10'
	},
	CHANNEL_USER_ID : {
		name : 'sgdb.channel.userID',
		alias : 'userID',
		defaultValue : null
	},
	CHANNEL_PASSWD : {
		name : 'sgdb.channel.password',
		alias : 'password',
		defaultValue : null
	},
	CHANNEL_CLIENTID : {
		name : 'sgdb.channel.clientId',
		alias : 'clientId',
		defaultValue : 'sgdb.nodejs-api.client'
	},

	CONNECTION_OPERATION_TIMEOUT : {
		name : 'sgdb.connection.operationTimeout',
		alias : null,
		defaultValue : '10000'
	}
};

exports.CONFIG_NAMES = CONFIG_NAMES;
