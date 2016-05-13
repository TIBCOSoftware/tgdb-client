package com.tibco.tgdb.query.impl;

import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;

import java.util.Date;
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
 *
 *		SVN Id: $Id:$
 *
 */

public class QueryImpl implements TGQuery {

    private ConnectionImpl connection;
    private long queryHashId;
    private TGQueryOption option;
    private Map<String, Object> parameters;

    private static NullParameter NULLPARAMETER = new NullParameter();

    static private class NullParameter {

    }
   
    // TGConneciton, String query, byte[] buff.
    // CHange TGConnection to ConnectionImpl, add long hashId,
    // Within buff, queryId, parameter names.
    // ConnectionImpl should be able to register query, for query close when connection closes.
    public QueryImpl(ConnectionImpl connection, long queryHashId) {
        this.connection = connection;
        this.queryHashId = queryHashId;
        option = TGQueryOption.DEFAULT_QUERY_OPTION;
    }
    
    @Override
    public void setBoolean(String name, boolean value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setChar(String name, char value) {
        parameters.put(name, value);
    }

    @Override
    public void setShort(String name, short value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setInt(String name, int value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setLong(String name, long value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setFloat(String name, float value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setDouble(String name, double value) {
        parameters.put(name, value);
        
    }

    @Override
    public void setString(String name, String value) {
        parameters.put(name, value);
        
    }
    
    @Override
    public void setDate(String name, Date value) {
        parameters.put(name, value);
        
    }
    
    @Override
    public void setBytes(String name, byte[] bos) {
        parameters.put(name, bos);
        
    }

    @Override
    public void setNull(String name) {
        parameters.put(name, NULLPARAMETER);
        
    }

    public void setOption(TGQueryOption option) {
        this.option = option;
    }

    @Override
    public TGResultSet execute() throws TGException {
        this.connection.executeQueryWithId(this.queryHashId, option);
        return null;
    }

    @Override
    public void close() throws TGException {
    	this.connection.closeQuery(this.queryHashId);
    }
    
}
