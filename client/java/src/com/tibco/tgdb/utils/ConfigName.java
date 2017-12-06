package com.tibco.tgdb.utils;

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

 * File name :ConfigName
 * Created on: 12/31/14
 * Created by: suresh

 * SVN Id: $Id: ConfigName.java 1678 2017-09-19 22:49:32Z ssubrama $
 */

/**
 * Enumeration class that provides a configuration name, alias name and a defaultValue;
 */
public enum ConfigName {

    ChannelDefaultHost
            ("tgdb.channel.defaultHost", null, "localhost"),

    ChannelDefaultPort
            ("tgdb.channel.defaultPort", null, "8700"),

    ChannelDefaultProtocol
            ("tgdb.channel.defaultProtocol", null, "tcp"),

    ChannelSendSize
            ("tgdb.channel.sendSize", "sendSize", "122"),

    ChannelRecvSize
            ("tgdb.channel.recvSize", "recvSize", "128"),

    ChannelPingInterval
            ("tgdb.channel.pingInterval", "pingInterval", "30"),

    ChannelConnectTimeout
            ("tgdb.channel.connectTimeout", "connectTimeout", "1000"), //1 sec timeout.

    ChannelFTHosts
            ("tgdb.channel.ftHosts", "ftHosts", null),

    ChannelDefaultUserID
            ("tgdb.channel.defaultUserID", null, null),

    ChannelUserID
            ("tgdb.channel.userID", "userID", null),

    ChannelPassword
            ("tgdb.channel.password", "password", null),

    ChannelClientId("tgdb.channel.clientId", "clientId",null),

    ConnectionPoolUseDedicatedChannelPerConnection
            ("tgdb.connectionpool.useDedicatedChannelPerConnection", "useDedicatedChannelPerConnection", "false"),

    ConnectionPoolDefaultPoolSize
            ("tgdb.connectionpool.defaultPoolSize", "defaultPoolSize", "10"),

    ConnectionOperationTimeout
            ("tgdb.connection.operationTimeout", "connectionOperationTimeout", "10000"),  //Represented in ms. Default Value is 10sec


    
    InvalidName(null, null, null);

    String propName;
    String defaultValue;
    String aliasName;

    ConfigName(String propName, String aliasName, String defaultValue)
    {
        this.propName       = propName;
        this.defaultValue   = defaultValue;
        this.aliasName      = aliasName;
    }

    /**
     * Return the ConfigName given its full qualified string form or its alias name.
     * @param name property config name
     * @return ConfigName associated to the name
     */
    public static ConfigName fromName(String name) {
        for (ConfigName cn : ConfigName.values()) {

            if (name.equalsIgnoreCase(cn.propName)) return cn;

            if (cn.aliasName != null) {
                if (name.equalsIgnoreCase(cn.aliasName)) return cn;
            }
        }

        return ConfigName.InvalidName;
    }

    /**
     * @return the Alias for this ConfigName
     */
    public String getAlias() {
        return aliasName;
    }


    /**
     * @return Property Name for this Config
     */
    public String getName() {
        return propName;
    }
}
