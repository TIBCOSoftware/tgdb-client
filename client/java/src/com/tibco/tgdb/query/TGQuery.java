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

 * File name : TGQuery.java
 * Created on: 1/22/15
 * Created by: suresh 

 * SVN Id: $Id: TGQuery.java 622 2016-03-19 20:51:12Z ssubrama $
 */


package com.tibco.tgdb.query;

import com.tibco.tgdb.exception.TGException;

public interface TGQuery {

    /**
     * Set Boolean parameter
     * @param name
     * @param value
     */
    void setBoolean(String name, boolean value);
    /**
     * Set an character Parameter
     * @param name
     * @param value
     */
    void setChar(String name, char value);

    /**
     * 
     * @param name
     * @param value
     */
    void setShort(String name, short value);

    /**
     * Set Integer parameter
     * @param name
     * @param value
     */
    void setInt(String name, int value);

    /**
     * Set Long parameter
     * @param name
     * @param value
     */
    void setLong(String name, long value);

    /**
     * Set Float Parameter
     * @param name
     * @param value
     */
    void setFloat(String name, float value);

    /**
     * Set Double Parameter
     * @param name
     * @param value
     */
    void setDouble(String name, double value);

    /**
     * Set String parameter
     * @param name
     * @param value
     */
    void setString(String name, String value);

    /**
     * Set Date parameter
     * @param name
     * @param value
     */
    void setDate(String name, java.util.Date value);

    /**
     * Set Bytes
     * @param name
     * @param bos
     */
    void setBytes(String name, byte[] bos);

    /**
     * Set the parameter to null
     * @param name
     */
    void setNull(String name);

    /**
     * Set the Query Option
     * @param option query option to be used for execution
     */
    public void setOption(TGQueryOption option);

    /**
     * Execute the Query
     * @return
     * @throws TGException
     */
    TGResultSet execute() throws TGException;

    /**
     * Close the Query
     * @throws TGException 
     */
    void close() throws TGException;
}
