package com.tibco.tgdb.connection;

import com.tibco.tgdb.exception.TGException;

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
 *  File name :TGConnectionPool.java
 *  Created on: 3/18/14
 *  Created by: suresh
 *
 *		SVN Id: $Id: TGConnectionPool.java 2180 2018-03-29 23:19:38Z ssubrama $
 *
 */

public interface TGConnectionPool {


    /**
     * Connect to the Server
     * @throws Exception - could be BadAuthentication or BadUrl
     */
    void connect() throws Exception;

    /**
     * Set Exception Listener
     * @param lsnr - an instance of exception listener
     */
    void setExceptionListener(TGConnectionExceptionListener lsnr);

    /**
     * Disconnect from the Server
     * @throws Exception if it cannot terminate and flush the sockets. One can safely ignore the exception
     */
    void disconnect() throws Exception;

    /**
     * Get a free Connection
     * @return a free Connection object
     * @throws com.tibco.tgdb.exception.TGConnectionTimeoutException is thrown when there is no connection in the pool
     * and the time specified has elapsed.
     * The property ConnectionReserveTimeoutSeconds or tgdb.connectionpool.ConnectionReserveTimeoutSeconds specifies the time
     * to wait in seconds. It has the following meaning
     * 0 :      Indefinite
     * -1 :     Immediate
     * &gt; :   That many seconds
     */
    TGConnection get() throws TGException;

    /**
     * Release the connection back to the pool
     * @param conn a TGConnection
     */
    void release(TGConnection conn);

    /**
     * Get the pool size
     * @return the pool size
     */
    int getPoolSize();




}
