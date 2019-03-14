package com.tibco.tgdb.connection;

import com.tibco.tgdb.exception.TGException;

import java.util.Map;

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
 *
 * File name :TGConnectionFactory
 * Created by: suresh
 *
 * SVN Id: $Id: TGConnectionFactory.java 2692 2018-11-11 18:32:34Z ssubrama $
 */

public abstract class TGConnectionFactory {
    private static final TGConnectionFactory gInstance = createConnectionFactory();
    private static final String TG_CONNECTIONFACTORY_PROVIDER = "com.tibco.tgdb.connection.TGConnectionFactory.Provider";


    /**
     * @return A global instance of Connection Factory
     */
    public static TGConnectionFactory getInstance() {
        return gInstance;
    }



    private static TGConnectionFactory createConnectionFactory()  {

        try {
            String clsName = System.getProperty(TG_CONNECTIONFACTORY_PROVIDER, "com.tibco.tgdb.connection.impl.DefaultConnectionFactory");
            Class<TGConnectionFactory> klazz = (Class<TGConnectionFactory>) Class.forName(clsName);
            return klazz.newInstance();
        }
        catch (Exception e) {
            e.printStackTrace();
            throw new RuntimeException(e);
        }

    }

    /**
     * Create a connection on the url using the name and password
     * Each connection will create a dedicated Channel for connection.
     * @param url The url for connection.  A URL is represented as a string of the form <BR>
     *            &lt;protocol&gt;://[user@]['['ipv6']'] | ipv4 [:][port][/]'{' name:value;... '}' <BR>
     *            protocol can be tcp or ssl.<BR>
     *            For ssl protocol, the following properties must also be present <BR>
     *            dbName:the database name <BR>
     *            verifyDBName:true<BR>
     *            <BR>
     *
     * @param userName The user name for connection. The userId provided overrides all other userIds that can be infered.
     *                 The rules for overriding are in this order<BR>
     *                 a. The argument 'userId' is the highest priority. If Null then <BR>
     *                 b. The user@url is considered. If that is Null <BR>
     *                 c. the "userID=value" from the URL string is considered.<BR>
     *                 d. If all of them is Null, then the default User associated to the installation will be taken.<BR>
     *                 <BR>
     *
     * @param password The managled or unmanagled password
     *                 <BR>
     * @param env optional environment. This environment will override every other environment values infered, and is
     *            specific for this connection only. The following properties are supported
     * <table>
     *  <caption>Connection Properties</caption>
     * 	<thead>
     * 		<tr>
     * 			<th style="width:auto;text-align:left">Full Name</th>
     * 			<th style="width:auto;text-align:left">Alias</th>
     * 			<th style="width:auto;text-align:left">Default Value</th>
     * 			<th style="width:auto;text-align:left">Description</th>
     * 		</tr>
     *	</thead>
     * 	<tbody>
     * 		<tr>
     * 			<td>tgdb.channel.defaultHost</td>
     * 			<td>-</td>
     * 			<td>localhost</td>
     * 			<td>The default host specifier</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.defaultPort</td>
     * 			<td>-</td>
     * 			<td>8700</td>
     * 			<td>The default port specifier</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.defaultProtocol</td>
     * 			<td>-</td>
     * 			<td>tcp</td>
     * 			<td>The default protocol</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.sendSize</td>
     * 			<td>sendSize</td>
     * 			<td>122</td>
     * 			<td>TCP send packet size in KBs</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.recvSize</td>
     * 			<td>recvSize</td>
     * 			<td>128</td>
     * 			<td>TCP recv packet size in KB</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.pingInterval</td>
     * 			<td>pingInterval</td>
     * 			<td>30</td>
     * 			<td>Keep alive ping intervals</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.connectTimeout</td>
     * 			<td>connectTimeout</td>
     * 			<td>1000</td>
     * 			<td>Timeout for connection to establish, before it gives up and tries the ftUrls if specified</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.ftHosts</td>
     * 			<td>ftHosts</td>
     * 			<td>-</td>
     * 			<td>Alternate fault tolerant list of &lt;host:port&gt; pair seperated by comma. </td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.ftRetryIntervalSeconds</td>
     * 			<td>ftRetryIntervalSeconds</td>
     * 			<td>10</td>
     * 			<td>The connect retry interval to ftHosts</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.ftRetryCount</td>
     * 			<td>ftRetryCount</td>
     * 			<td>3</td>
     * 			<td>The number of times ro retry </td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.defaultUserID</td>
     * 			<td>-</td>
     * 			<td>-</td>
     * 			<td>The default user Id for the connection</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.userID</td>
     * 			<td>userID</td>
     * 			<td>-</td>
     * 			<td>The user id for the connection if it is not specified in the API. See the rules for picking the user name</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.password</td>
     * 			<td>password</td>
     * 			<td>-</td>
     * 			<td>The password for the username</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.channel.clientId</td>
     * 			<td>clientId</td>
     * 			<td>tgdb.java-api.client</td>
     * 			<td>The client id to be used for the connection</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.connection.dbName</td>
     * 			<td>dbName</td>
     * 			<td>-</td>
     * 			<td>The database name the client is connecting to. It is used as part of verification for ssl channels</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.connectionpool.useDedicatedChannelPerConnection</td>
     * 			<td>useDedicatedChannelPerConnection</td>
     * 			<td>false</td>
     * 			<td>A boolean value indicating either to multiplex mulitple connections on a single tcp socket or use dedicate socket per connection. A true value consumes resource but provides good performance. Also check the max number of connections</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.connectionpool.defaultPoolSize</td>
     * 			<td>defaultPoolSize</td>
     * 			<td>10</td>
     * 			<td>The default connection pool size to use when creating a ConnectionPool</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.connectionpool.connectionReserveTimeoutSeconds</td>
     * 			<td>connectionReserveTimeoutSeconds</td>
     * 			<td>10</td>
     * 			<td>A timeout parameter indicating how long to wait before getting a connection from the pool</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.connection.operationTimeoutSeconds</td>
     * 			<td>connectionOperationTimeoutSeconds</td>
     * 			<td>10</td>
     * 			<td>A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior.</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.provider.name</td>
     * 			<td>tlsProviderName</td>
     * 			<td>SunJSSE</td>
     * 			<td>Transport level Security provider. Work with your InfoSec team to change this value</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.provider.className</td>
     * 			<td>tlsProviderClassName</td>
     * 			<td>com.sun.net.ssl.internal.ssl.Provider</td>
     * 			<td>The underlying Provider implementation. Work with your InfoSec team to change this value.</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.provider.configFile</td>
     * 			<td>tlsProviderConfigFile</td>
     * 			<td>-</td>
     * 			<td>Some providers require extra configuration paramters, and it can be passed as a file</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.protocol</td>
     * 			<td>tlsProtocol</td>
     * 			<td>TLSv1.2</td>
     * 			<td>tlsProtocol version. The system only supports 1.2+</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.cipherSuites</td>
     * 			<td>cipherSuites</td>
     * 			<td>-</td>
     * 			<td>A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher list and Openssl list that supports 1.2 protocol </td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.verifyDBName</td>
     * 			<td>verifyDBName</td>
     * 			<td>false</td>
     * 			<td>Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL.</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.expectedHostName</td>
     * 			<td>expectedHostName</td>
     * 			<td>-</td>
     * 			<td>The expected hostName for the certificate. This is for future use</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.tls.trustedCertificates</td>
     * 			<td>trustedCertificates</td>
     * 			<td>-</td>
     * 			<td>The list of trusted Certificates</td>
     * 		</tr>
     * 		<tr>
     * 			<td>tgdb.security.keyStorePassword</td>
     * 			<td>keyStorePassword</td>
     * 			<td>-</td>
     * 			<td>The Keystore for the password</td>
     * 		</tr>
     * 	</tbody>
     * </table>
     *
     * @return TGConnection - an instance of connection to the server with a dedicated channel
     * @throws com.tibco.tgdb.exception.TGException - If it cannot create a connection to the server successfully
     */
    public abstract TGConnection createConnection(String url, String userName, String password, Map<String, String> env) throws TGException;

    /**
     * Create a Connection Pool of pool size on the the url using the name and password. Each connection in the pool will default
     * use a shared channel, but this can be overriden by setting the value property tgdb.connectionpool.useDedicatedChannel=true
     * @param url The url for the channel used in the connection pool.
     * @param userName  The user name for connection.
     * @param password  The password mangled or unmangled
     * @param poolSize the size of the pool
     * @param env optional environment. This environment will override every other environment values infered, and is specific for this pool only
     * @return A Connection Pool
     * @throws com.tibco.tgdb.exception.TGException - If it cannot create a connectionpool to the server successfully
     * @see TGConnectionFactory#createConnection(java.lang.String, java.lang.String, java.lang.String, java.util.Map)
     */
    public abstract TGConnectionPool createConnectionPool(String url, String userName, String password, int poolSize, Map<String, String> env) throws TGException;;


}
