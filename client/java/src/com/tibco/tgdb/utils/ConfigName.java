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

 * SVN Id: $Id: ConfigName.java 2686 2018-11-09 20:55:57Z ssubrama $
 */



/**
 * Enumeration class that provides a configuration name, alias name and a defaultValue;
 */
public enum ConfigName {


    ChannelDefaultHost
            ("tgdb.channel.defaultHost",
                    null,
                    "localhost",
                    "The default host specifier"),

    ChannelDefaultPort
            ("tgdb.channel.defaultPort",
                    null,
                    "8700",
                    "The default port specifier"),

    ChannelDefaultProtocol
            ("tgdb.channel.defaultProtocol",
                    null,
                    "tcp",
                    "The default protocol"),

    ChannelSendSize
            ("tgdb.channel.sendSize",
                    "sendSize",
                    "122",
                    "TCP send packet size in KBs"),

    ChannelRecvSize
            ("tgdb.channel.recvSize",
                    "recvSize",
                    "128",
                    "TCP recv packet size in KB"),

    ChannelPingInterval
            ("tgdb.channel.pingInterval", "pingInterval", "30", "Keep alive ping intervals"),

    ChannelConnectTimeout
            ("tgdb.channel.connectTimeout", "connectTimeout", "1000", "Timeout for connection to establish, before it gives up and tries the ftUrls if specified"), //1 sec timeout.

    ChannelFTHosts
            ("tgdb.channel.ftHosts", "ftHosts", null, "Alternate fault tolerant list of &lt;host:port&gt; pair seperated by comma. "),

    ChannelFTRetryIntervalSeconds
            ("tgdb.channel.ftRetryIntervalSeconds", "ftRetryIntervalSeconds", "10", "The connect retry interval to ftHosts"),

    ChannelFTRetryCount
            ("tgdb.channel.ftRetryCount", "ftRetryCount", "3", "The number of times ro retry "),

    ChannelDefaultUserID
            ("tgdb.channel.defaultUserID", null, null, "The default user Id for the connection"),

    ChannelUserID
            ("tgdb.channel.userID", "userID", null, "The user id for the connection if it is not specified in the API. See the rules for picking the user name"),

    ChannelPassword
            ("tgdb.channel.password", "password", null, "The password for the username"),

    ChannelClientId("tgdb.channel.clientId", "clientId","tgdb.java-api.client","The client id to be used for the connection"),

    ConnectionDatabaseName("tgdb.connection.dbName", "dbName", null,"The database name the client is connecting to. It is used as part of verification for ssl channels"),

    ConnectionPoolUseDedicatedChannelPerConnection
            ("tgdb.connectionpool.useDedicatedChannelPerConnection", "useDedicatedChannelPerConnection", "false", "A boolean value indicating either to multiplex mulitple connections on a single tcp socket or use dedicate socket per connection. A true value consumes resource but provides good performance. Also check the max number of connections"),

    ConnectionPoolDefaultPoolSize
            ("tgdb.connectionpool.defaultPoolSize", "defaultPoolSize", "10","The default connection pool size to use when creating a ConnectionPool"),

    ConnectionReserveTimeoutSeconds
            ("tgdb.connectionpool.connectionReserveTimeoutSeconds", "connectionReserveTimeoutSeconds", "10","A timeout parameter indicating how long to wait before getting a connection from the pool"),
            //0 = mean immediate, Integer Max for indefinite.

    ConnectionOperationTimeoutSeconds
            ("tgdb.connection.operationTimeoutSeconds", "connectionOperationTimeoutSeconds", "10","A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior."),  //Represented in ms. Default Value is 10sec

    //TSL Parameters
    TlsProviderName //The default is the Sun JSSE.
            ("tgdb.tls.provider.name", "tlsProviderName", "SunJSSE","Transport level Security provider. Work with your InfoSec team to change this value"),

    TlsProviderClassName //The default is the Sun JSSE.
            ("tgdb.tls.provider.className", "tlsProviderClassName", "com.sun.net.ssl.internal.ssl.Provider", "The underlying Provider implementation. Work with your InfoSec team to change this value."),

    TlsProviderConfigFile
            ("tgdb.tls.provider.configFile", "tlsProviderConfigFile", null,"Some providers require extra configuration paramters, and it can be passed as a file"),

    TlsProtocol
            ("tgdb.tls.protocol", "tlsProtocol", "TLSv1.2","tlsProtocol version. The system only supports 1.2+"),

    TlsCipherSuites
            ("tgdb.tls.cipherSuites", "cipherSuites", null,"A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher list and Openssl list that supports 1.2 protocol "), //Use the Default Cipher Suites

    TlsVerifyDatabaseName
            ("tgdb.tls.verifyDBName", "verifyDBName", "false","Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL."),

    TlsExpectedHostName
            ("tgdb.tls.expectedHostName", "expectedHostName", null,"The expected hostName for the certificate. This is for future use"),

    TlsTrustedCertificates
            ("tgdb.tls.trustedCertificates", "trustedCertificates", null,"The list of trusted Certificates"),

    KeyStorePassword
            ("tgdb.security.keyStorePassword", "keyStorePassword", null, "The Keystore for the password"),



    InvalidName(null, null, null, null);

    String propName;
    String defaultValue;
    String aliasName;
    String desc;

    //public static String TG_SUN_FIPSWRAPPER_CLASSNM = "com.tibco.tgdb.crypto.TGSunFIPsProviderFacade";
    //public static String TG_IBM_FIPSWRAPPER_CLASSNM = "com.tibco.tgdb.crypto.TGIBMFIPsProviderFacade";


    ConfigName(String propName, String aliasName, String defaultValue, String desc)
    {
        this.propName       = propName;
        this.defaultValue   = defaultValue;
        this.aliasName      = aliasName;
        this.desc           = desc;
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

    public static void main(String[] args) {
        StringBuilder builder = new StringBuilder();
        //<tr><td>User</td><td>Typical end-user</td></tr>
        builder.append("* <table>\n");
        builder.append("* \t<thead>\n");
        builder.append("* \t\t<tr>\n");
        builder.append("* \t\t\t<th style=\"width:auto;text-align:left\">").append("Full Name").append("</th>\n");
        builder.append("* \t\t\t<th style=\"width:auto;text-align:left\">").append("Alias").append("</th>\n");
        builder.append("* \t\t\t<th style=\"width:auto;text-align:left\">").append("Default Value").append("</th>\n");
        builder.append("* \t\t\t<th style=\"width:auto;text-align:left\">").append("Description").append("</th>\n");
        builder.append("* \t\t</tr>\n");
        builder.append("*\t</thead>\n");
        builder.append("* \t<tbody>\n");
        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) break;
            builder.append("* \t\t<tr>\n");
            builder.append("* \t\t\t<td>").append(cn.propName).append("</td>\n");
            builder.append("* \t\t\t<td>").append(cn.aliasName == null ? "-" : cn.aliasName).append("</td>\n");
            builder.append("* \t\t\t<td>").append(cn.defaultValue == null ? "-" : cn.defaultValue).append("</td>\n");
            builder.append("* \t\t\t<td>").append(cn.desc == null ? '-' : cn.desc).append("</td>\n");
            builder.append("* \t\t</tr>\n");
        }
        builder.append("* \t</tbody>\n");
        builder.append("* </table>\n");
        System.out.println(builder.toString());
    }


}
