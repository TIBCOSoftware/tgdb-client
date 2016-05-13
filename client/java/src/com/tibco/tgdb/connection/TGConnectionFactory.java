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
 * SVN Id: $Id: TGConnectionFactory.java 622 2016-03-19 20:51:12Z ssubrama $
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
     * @param url The url for connection.
     * @param userName The user name for connection
     * @param password The managled or unmanagled password
     * @param env optional environment. This environment will override every other environment values infered, and is specific for this connection only.
     * @return TGConnection - an instance of connection to the server with a dedicated channel
     * @throws com.tibco.tgdb.exception.TGException - If it cannot create a connection to the server successfully
     */
    public abstract TGConnection createConnection(String url, String userName, String password, Map<String, String> env) throws TGException;;

    /**
     * Create a Connection Pool of pool size on the the url using the name and password. Each connection in the pool will default
     * use a shared channel, but this can be overriden by setting the value property tgdb.connectionpool.useDedicatedChannel=true
     * @param url The url for the channel used in the connection pool. @see
     * @param userName  The user name for connection
     * @param password  The password mangled or unmangled
     * @param poolSize the size of the pool
     * @param env optional environment. This environment will override every other environment values infered, and is specific for this pool only
     * @return A Connection Pool
     * @throws com.tibco.tgdb.exception.TGException - If it cannot create a connectionpool to the server successfully
     */
    public abstract TGConnectionPool createConnectionPool(String url, String userName, String password, int poolSize, Map<String, String> env) throws TGException;;


}
