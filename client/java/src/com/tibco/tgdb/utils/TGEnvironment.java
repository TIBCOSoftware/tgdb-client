package com.tibco.tgdb.utils;

import java.util.EnumMap;
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

 * File name :TGEnvironment
 * Created on: 12/16/14
 * Created by: suresh

 * SVN Id: $Id: TGEnvironment.java 2344 2018-06-11 23:21:45Z ssubrama $
 */
public class TGEnvironment {



    static TGEnvironment gInstance = new TGEnvironment();


    public static TGEnvironment getInstance() { return gInstance;}

    private EnumMap<ConfigName, String> env;

    private TGEnvironment() {

        env = new EnumMap<ConfigName, String>(ConfigName.class);
        for(ConfigName cn : ConfigName.values()) {
            env.put(cn, cn.defaultValue);
        }
        for (Map.Entry entry : System.getProperties().entrySet())
        {
            this.setProperty(entry.getKey().toString(), entry.getValue().toString());
        }

    }
    public int  getChannelSendSize() {

        String value = env.get(ConfigName.ChannelSendSize);
        if (value != null) {
            return Integer.parseInt(value);
        }
        return 122;
    }

    public int getChannelReceiveSize() {
        String value = env.get(ConfigName.ChannelRecvSize);
        if (value != null) {
            return Integer.parseInt(value);
        }
        return 128;

    }



    public int getChannelPingInterval() {

        String value = env.get(ConfigName.ChannelPingInterval);
        if (value != null) {
            return Integer.parseInt(value);
        }
        return 30;

    }

    public int getChannelConnectTimeout() {
        String value = env.get(ConfigName.ChannelConnectTimeout);
        if (value != null) {
            return Integer.parseInt(value);
        }
        return 1000;

    }

    public String getChannelFTHosts() {

        return env.get(ConfigName.ChannelFTHosts);

    }

    public String getChannelDefaultUser() {
        return env.get(ConfigName.ChannelDefaultUserID);
    }

    /**
     * Get the Connection User as specified in the Environment
     * @return the Channel User
     */
    public String getChannelUser() {
        return env.get(ConfigName.ChannelUserID);
    }


    public void setProperty(String name, String value) {
        ConfigName cn = ConfigName.fromName(name);
        if (cn == ConfigName.InvalidName) return;
        env.put(cn, value);
    }

    public String getProperty(String name) {

        ConfigName cn = ConfigName.fromName(name);
        if (cn == ConfigName.InvalidName) return null;

        return env.get(cn);

    }

    public String getProperty(String name, String defaultValue) {

        String value = getProperty(name);

        return value == null ? defaultValue : value;

    }

    public int getChannelDefaultPort() {

        String value = env.get(ConfigName.ChannelDefaultPort);

        return Integer.parseInt(value);
    }

    public String getChannelDefaultHost() {
        return env.get(ConfigName.ChannelDefaultHost);
    }

    public int getConnectionPoolDefaultPoolSize() {
        String value = env.get(ConfigName.ConnectionPoolDefaultPoolSize);
        return Integer.parseInt(value);
    }

    public String getChannelClientId() {
        String value = env.get(ConfigName.ChannelClientId);
        if (value == null) {
            return "tgdb-client";
        }
        return value;
    }

    public SortedProperties<String, String> getAsSortedProperties()
    {
        SortedProperties<String, String> sp = new SortedProperties<String, String>();
        for (Map.Entry<ConfigName, String> e : env.entrySet()) {
            ConfigName name = e.getKey();
            String value = e.getValue();
            if (name.equals(ConfigName.InvalidName)) continue;
            sp.put(name.getName(), value);
        }
        return sp;
    }

    public String getDefaultDateTimeFormat()
    {
        //SS:TODO
        return "mm-dd-yyyy hh:mm:ss";
    }
}
