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
 * <p/>
 * File name : DefaultConnectionFactory.${EXT}
 * Created on: 1/10/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: DefaultConnectionFactory.java 2179 2018-03-29 21:49:54Z ssubrama $
 */


package com.tibco.tgdb.connection.impl;

import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelFactory;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.channel.impl.LinkUrl;
import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.connection.TGConnectionPool;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGEnvironment;
import com.tibco.tgdb.utils.TGProperties;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 *
 */
public class DefaultConnectionFactory extends TGConnectionFactory {

    @Override
    public TGConnection createConnection(String url, String userName, String password, Map<String, String> env) throws TGException
    {
        ConnectionPoolImpl connPool = (ConnectionPoolImpl) createConnectionPool(url, userName, password, 1, env);
        return connPool.getConnection();
    }

    @Override
    public TGConnectionPool createConnectionPool(String url, String userName, String password, int poolSize, Map<String, String> env) throws TGException
    {
        if (poolSize <= 0 ) {
            poolSize = TGEnvironment.getInstance().getConnectionPoolDefaultPoolSize();
        }

        TGProperties<String, String> properties = new SortedProperties<>(String.CASE_INSENSITIVE_ORDER);
        TGProperties<String, String> defProps = TGEnvironment.getInstance().getAsSortedProperties();
        properties.putAll(defProps);

        if (env != null) {
           properties.putAll(env);
        }
        TGChannelUrl channelUrl = LinkUrl.parse(url);
        properties.putAll(channelUrl.getProperties());
        TGProperties.setUserAndPassword(properties, userName, password);

        return new ConnectionPoolImpl(channelUrl, poolSize, properties);
    }

}
