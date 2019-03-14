package com.tibco.tgdb.utils;

import com.tibco.tgdb.exception.TGException;

import java.util.Properties;
import java.util.SortedMap;

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

 * File name : TGProperties.java
 * Created on: 2/5/15
 * Created by: suresh
  * SVN Id: $Id: TGProperties.java 2316 2018-04-26 23:49:37Z ssubrama $
 */

public interface TGProperties<K,V> extends SortedMap<K,V> {

    V getProperty(ConfigName cn);

    V getProperty(ConfigName cn, V defaultValue);

    public static void setUserAndPassword(TGProperties<String,String> properties, String userName, String password) throws TGException
    {
        if ((userName == null) || (userName.length() == 0)) {
            userName = properties.getProperty(ConfigName.ChannelUserID, TGEnvironment.getInstance().getChannelDefaultUser());
            if ((userName == null) || (userName.length() == 0)) {
                throw new TGException("user name not specified.");
            }
        }
        if ((password == null) || (password.length() == 0)) {
            password = properties.getProperty(ConfigName.ChannelPassword);
        }
        properties.put(ConfigName.ChannelUserID.getName(), userName);
        if ((password != null) || (password.length() != 0)) {
            properties.put(ConfigName.ChannelPassword.getName(), password);
        }

    }

    /**
     * Get the Property As a int value
     * @param cn
     * @return
     */
    int getPropertyAsInt(ConfigName cn);

    int getPropertyAsInt(ConfigName cn, int defaultValue);

    /**
     * Get the Property as a Long value
     * @param cn
     * @return
     */
    long getPropertyAsLong(ConfigName cn);

    long getPropertyAsLong(ConfigName cn, long defaultValue);

    /**
     * Get the Property as a Boolean
     * @param cn
     * @return
     */
    boolean getPropertyAsBoolean(ConfigName cn);

    boolean getPropertyAsBoolean(ConfigName cn, boolean defaultValue);
}
