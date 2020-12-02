/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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

 * File name : ConfigName.java
 * Created on: 12/31/2014
 * Created by: suresh
 * SVN Id: $Id: ConfigName.java 4175 2020-07-17 21:17:09Z ssubrama $
 */

package com.tibco.tgdb.utils;

/**
 * Enumeration class that provides a configuration name, alias name and a defaultValue;
 */
public enum ConfigName {
    ChannelDefaultHost(
            "tgdb.channel.defaultHost",
            "defaultHost",
            "localhost",
            "The default host specifier"
    ),

    ChannelDefaultPort(
            "tgdb.channel.defaultPort",
            "defaultPort",
            "8222",
            "The default port specifier"
    ),

    ChannelDefaultProtocol(
            "tgdb.channel.defaultProtocol",
            "defaultProtocol",
            "tcp",
            "The default protocol"
    ),

    ChannelSendSize(
            "tgdb.channel.sendSize",
            "sendSize",
            "122",
            "TCP send packet size in KBs"
    ),

    ChannelRecvSize(
            "tgdb.channel.recvSize",
            "recvSize",
            "128",
            "TCP recv packet size in KB"
    ),

    ChannelPingInterval(
            "tgdb.channel.pingInterval",
            "pingInterval",
            "30",
            "Keep alive ping intervals"
    ),

    ChannelConnectTimeout( //1 sec timeout.
            "tgdb.channel.connectTimeout",
            "connectTimeout",
            "1000",
            "Timeout for connection to establish, before it gives up and tries the ftUrls if specified"
    ),

    ChannelFTHosts(
            "tgdb.channel.ftHosts",
            "ftHosts",
            null,
            "Alternate fault tolerant list of &lt;host:port&gt; pair seperated by comma. "
    ),

    ChannelFTRetryIntervalSeconds(
            "tgdb.channel.ftRetryIntervalSeconds",
            "ftRetryIntervalSeconds",
            "10",
            "The connect retry interval to ftHosts"
    ),

    ChannelFTRetryCount(
            "tgdb.channel.ftRetryCount",
            "ftRetryCount",
            "3",
            "The number of times ro retry "
    ),

    ChannelDefaultUserID(
            "tgdb.channel.defaultUserID",
            "defaultUserID",
            null,
            "The default user Id for the connection"
    ),

    ChannelUserID(
            "tgdb.channel.userID",
            "userID",
            null,
            "The user id for the connection if it is not specified in the API. See the rules for picking the user name"
    ),

    ChannelPassword(
            "tgdb.channel.password",
            "password",
            null,
            "The password for the username"
    ),

    ChannelClientId(
            "tgdb.channel.clientId",
            "clientId",
            "tgdb.java-api.client",
            "The client id to be used for the connection"
    ),


    ConnectionPoolUseDedicatedChannelPerConnection(
            "tgdb.connectionpool.useDedicatedChannelPerConnection",
            "useDedicatedChannelPerConnection",
            "false",
            "A boolean value indicating either to multiplex mulitple connections on a single tcp socket or use dedicate " +
                    "socket per connection. A true value consumes resource but provides good performance. Also check " +
                    "the max number of connections"
    ),

    ConnectionPoolDefaultPoolSize(
            "tgdb.connectionpool.defaultPoolSize",
            "defaultPoolSize",
            "10",
            "The default connection pool size to use when creating a ConnectionPool"
    ),

    ConnectionPoolReserveTimeoutSeconds(
            "tgdb.connectionpool.connectionReserveTimeoutSeconds",
            "connectionpoolReserveTimeoutSeconds",
            "10",
            "A timeout parameter indicating how long to wait before getting a connection from the pool"
    ),
    //0 = mean immediate, Integer Max for indefinite.
    ConnectionDatabaseName(
            "tgdb.connection.dbName",
            "dbName",
            null,
            "The database name the client is connecting to. It is used as part of verification for ssl channels"
    ),

    ConnectionSpecifiedRoles(
            "tgdb.connection.specifiedRoles",
            "roles",
            null,
            "The role name(s) that the user wants to log in as."
    ),

    ConnectionOperationTimeoutSeconds(//Represented in ms. Default Value is 10sec
            "tgdb.connection.operationTimeoutSeconds",
            "connectionOperationTimeoutSeconds",
            "10",
            "A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior."
    ),

    ConnectionIdleTimeoutSeconds(//Represented in secs. Default Value is 3600 : 1hr
            "tgdb.connection.idleTimeoutSeconds",
            "connectionIdleTimeoutSeconds",
            "3600",
            "An idle timeout parameter requested to server, before the server disconnects. This may/may not be honored by the server"
    ),

    ConnectionDateFormat(
            "tgdb.connection.dateFormat",
            "dateFormat",
            "YYYY-MM-DD",
            "Date format for this connection"
    ),

    ConnectionTimeFormat(
            "tgdb.connection.timeFormat",
            "timeFormat",
            "HH:mm:ss",
            "Date format for this connection"
    ),
    ConnectionTimeStampFormat(
            "tgdb.connection.timeStampFormat",
            "timeStampFormat",
            "YYYY-MM-DD HH:mm:ss.zzz",
            "Timestamp format for this connection"
    ),

    ConnectionLocale(
            "tgdb.connection.locale",
            "locale",
            "en_US",
            "Locale for this connection"
    ),
    ConnectionTimezone(
            "tgdb.connection.timezone",
            "timezone",
            "Americas/Los_Angeles",
            "Timezone to use for this connection"
    ),

    ConnectionDefaultQueryLanguage(
            "tgdb.connection.defaultQueryLanguage",
            "queryLanguage",
            "tgql",
            "Default query lanaguge format for this connection"
    ),

    //TSL Parameters
    TlsProviderName(
            "tgdb.tls.provider.name",
            "tlsProviderName",
            "SunJSSE", //The default is the Sun JSSE.
            "Transport level Security provider. Work with your InfoSec team to change this value"
    ),

    TlsProviderClassName(
            "tgdb.tls.provider.className",
            "tlsProviderClassName",
            "com.sun.net.ssl.internal.ssl.Provider", //The default is the Sun JSSE.
            "The underlying Provider implementation. Work with your InfoSec team to change this value."
    ),

    TlsProviderConfigFile(
            "tgdb.tls.provider.configFile",
            "tlsProviderConfigFile",
            null,
            "Some providers require extra configuration paramters, and it can be passed as a file"
    ),

    TlsProtocol(
            "tgdb.tls.protocol",
            "tlsProtocol",
            "TLSv1.2",
            "tlsProtocol version. The system only supports 1.2+"
    ),

    TlsCipherSuites( //Use the Default Cipher Suites
            "tgdb.tls.cipherSuites",
            "cipherSuites",
            null,
            "A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher " +
                    "list and Openssl list that supports 1.2 protocol "
    ),

    TlsVerifyDatabaseName(
            "tgdb.tls.verifyDBName",
            "verifyDBName",
            "false",
            "Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL."
    ),

    TlsExpectedHostName(
            "tgdb.tls.expectedHostName",
            "expectedHostName",
            null,
            "The expected hostName for the certificate. This is for future use"
    ),

    TlsTrustedCertificates(
            "tgdb.tls.trustedCertificates",
            "trustedCertificates",
            null,
            "The list of trusted Certificates"
    ),

    KeyStorePassword(
            "tgdb.security.keyStorePassword",
            "keyStorePassword",
            null,
            "The Keystore for the password"
    ),
    EnableConnectionTrace(
            "tgdb.connection.enableTrace",
            "enableTrace",
            "false",
            "The flag for debugging purpose, to enable the commit trace"
    ),
    ConnectionTraceDir(
            "tgdb.connection.enableTraceDir",
            "enableTraceDir",
            ".",
            "The base directory to hold commit trace log"
    ),

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

    public static String toJavaDoc() {
        StringBuilder builder = new StringBuilder();
        //<tr><td>User</td><td>Typical end-user</td></tr>
        builder.append("* <table>\n");
        builder.append("* <caption>Connection Properties</caption>\n");
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
        return builder.toString();
    }

    public static String toCCode() {
        StringBuilder builder = new StringBuilder();
        System.out.println("Writing Declaration....");
        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) break;
            if (cn.aliasName == null) continue;
            builder.append("/**").append("\n");
            builder.append("* Get ").append(toCamelCase(cn.aliasName)).append("\n");
            builder.append("* @param confignames Retrieve the value for the property ").append(cn.aliasName).append(" from confignames").append("\n");
            builder.append("* @return the value").append("\n");
            builder.append("*/").append("\n");
            builder.append("tg_char*\t").append("tg_configname_get").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames);").append("\n");
            builder.append("\n");

            builder.append("/**").append("\n");
            builder.append("* Set ").append(toCamelCase(cn.aliasName)).append("\n");
            builder.append("* @param confignames Set the value for the property ").append(cn.aliasName).append("\n");
            builder.append("*/").append("\n");
            builder.append("void   \t").append("tg_configname_set").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames, char* value);").append("\n");
            builder.append("\n");

        }

        System.out.println("Writing Definition....");


        builder.append("#include <net/tgconfigname.h>").append("\n");
        builder.append("\n");
        builder.append("/////////////////////////////////////////////////////////////////////////////////////////////////\n" +
                "/////// The below code is generated from ConfigName.java. Use that to keep it in sync. /////////\n" +
                "///////////////////////////////////////////////////////////////////////////////////////////////\n");
        builder.append("static tg_configname systemDefaults[] =").append("\n");
        builder.append("{").append("\n");
        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) {
                builder.append("\t{ NULL, NULL, NULL, NULL }").append(";");
                break;
            }
            String defaultValue = cn.defaultValue == null ? "NULL" : "\"" + cn.defaultValue + "\"";
            builder.append("\t").append("{").append("\n");
            builder.append("\t").append("\t").append('"').append(cn.propName).append('"').append(',').append("\n");
            builder.append("\t").append("\t").append('"').append(cn.aliasName).append('"').append(',').append("\n");
            builder.append("\t").append("\t").append(defaultValue).append(',').append("\n");
            builder.append("\t").append("\t").append('"').append(cn.desc).append('"').append(',').append("\n");
            builder.append("\t").append("},").append("\n");
        }
        builder.append("};").append("\n");
        builder.append("\n");


        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) break;
            if (cn.aliasName == null) continue;
            String defaultValue = cn.defaultValue == null ? "NULL" : "\"" + cn.defaultValue + "\"";
            builder.append("tg_char* ").append("tg_configname_get").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames)\n");
            builder.append("{\n");
            builder.append("\t").append("if (!confignames) return NULL;  \n");
            builder.append("\t").append("char* value = NULL;").append("\n");
            builder.append("\t").append("value = tg_properties_getProperty(confignames, ").
                                 append('"').append(cn.propName).append('"').append(");").append("\n");
            builder.append("\t").append("if (value) return value;").append("\n");
            builder.append("\t").append("return tg_properties_getPropertyEx(confignames, ").
                                 append('"').append(cn.aliasName).append('"').
                                 append(", ").
                                 append(defaultValue).append(");").append("\n");
            builder.append("}\n");
            builder.append("\n");

            builder.append("void ").append("tg_configname_set").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames, char* value)\n");
            builder.append("{\n");
            builder.append("\t").append("if (!confignames) return ;  \n");
            builder.append("\t").append("tg_properties_setProperty(confignames, ").
                                 append('"').append(cn.propName).append('"').append(", ").
                                 append("value);").append("\n");
            builder.append("\t").append("tg_properties_setProperty(confignames, ").
                                 append('"').append(cn.aliasName).append('"').append(", ").
                                 append("value);").append("\n");
            builder.append("}\n");
            builder.append("\n");
        }


        return builder.toString();
    }

    public static String toCCodeV2() {
        StringBuilder builder = new StringBuilder();
        System.out.println("Writing Header file");
        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) break;
            if (cn.aliasName == null) continue;
            builder.append("/**").append("\n");
            builder.append("* Get ").append(toCamelCase(cn.aliasName)).append("\n");
            builder.append("* @param confignames Retrieve the value for the property ").append(cn.aliasName).append(" from confignames").append("\n");
            builder.append("* @return the value").append("\n");
            builder.append("*/").append("\n");
            builder.append("tg_char*\t").append("tg_configname_get").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames);").append("\n");
            builder.append("\n");

            builder.append("/**").append("\n");
            builder.append("* Set ").append(toCamelCase(cn.aliasName)).append("\n");
            builder.append("* @param confignames Set the value for the property ").append(cn.aliasName).append("\n");
            builder.append("*/").append("\n");
            builder.append("void   \t").append("tg_configname_set").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames, char* value);").append("\n");
            builder.append("\n");

        }

        System.out.println("Writing Definition....");


        builder.append("#include <net/tgconfigname.h>").append("\n");
        builder.append("\n");
        builder.append("/////////////////////////////////////////////////////////////////////////////////////////////////\n" +
                "/////// The below code is generated from ConfigName.java. Use that to keep it in sync. /////////\n" +
                "///////////////////////////////////////////////////////////////////////////////////////////////\n");
        builder.append("static tg_configname systemDefaults[] =").append("\n");
        builder.append("{").append("\n");
        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) {
                builder.append("\t{ NULL, NULL, NULL, NULL }").append(";");
                break;
            }
            String defaultValue = cn.defaultValue == null ? "NULL" : "\"" + cn.defaultValue + "\"";
            builder.append("\t").append("{").append("\n");
            builder.append("\t").append("\t").append('"').append(cn.propName).append('"').append(',').append("\n");
            builder.append("\t").append("\t").append('"').append(cn.aliasName).append('"').append(',').append("\n");
            builder.append("\t").append("\t").append(defaultValue).append(',').append("\n");
            builder.append("\t").append("\t").append('"').append(cn.desc).append('"').append(',').append("\n");
            builder.append("\t").append("},").append("\n");
        }
        builder.append("};").append("\n");
        builder.append("\n");


        for (ConfigName cn : ConfigName.values()) {
            if (cn.propName == null) break;
            if (cn.aliasName == null) continue;
            String defaultValue = cn.defaultValue == null ? "NULL" : "\"" + cn.defaultValue + "\"";
            builder.append("tg_char* ").append("tg_configname_get").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames)\n");
            builder.append("{\n");
            builder.append("\t").append("if (!confignames) return NULL;  \n");
            builder.append("\t").append("char* value = NULL;").append("\n");
            builder.append("\t").append("value = tg_properties_getProperty(confignames, ").
                    append('"').append(cn.propName).append('"').append(");").append("\n");
            builder.append("\t").append("if (value) return value;").append("\n");
            builder.append("\t").append("return tg_properties_getPropertyEx(confignames, ").
                    append('"').append(cn.aliasName).append('"').
                    append(", ").
                    append(defaultValue).append(");").append("\n");
            builder.append("}\n");
            builder.append("\n");

            builder.append("void ").append("tg_configname_set").append(toCamelCase(cn.aliasName)).append("(tg_properties_ptr confignames, char* value)\n");
            builder.append("{\n");
            builder.append("\t").append("if (!confignames) return ;  \n");
            builder.append("\t").append("tg_properties_setProperty(confignames, ").
                    append('"').append(cn.propName).append('"').append(", ").
                    append("value);").append("\n");
            builder.append("\t").append("tg_properties_setProperty(confignames, ").
                    append('"').append(cn.aliasName).append('"').append(", ").
                    append("value);").append("\n");
            builder.append("}\n");
            builder.append("\n");
        }


        return builder.toString();
    }

    static String toCamelCase(String string) {
        StringBuilder sb = new StringBuilder(string);
        sb.replace(0, 1, string.substring(0, 1).toUpperCase());
        return sb.toString();

    }
    public static void main(String[] args) {
        if (args[0].equalsIgnoreCase("javadoc")) {
            System.out.println(toJavaDoc());
            return;
        }

        if (args[0].equalsIgnoreCase("C")) {
            System.out.println(toCCode());
            return;
        }

    }


}
