package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelFactory;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGEnvironment;
import com.tibco.tgdb.utils.TGProperties;

import java.util.Collections;
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
 * <p/>
 * File name :ChannelFactoryImpl
 * Created on: 12/26/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: ChannelFactoryImpl.java 583 2016-03-15 02:02:39Z vchung $
 */
public class ChannelFactoryImpl extends TGChannelFactory {

    @Override
    public  TGChannel createChannel(String urlPath, String userName, String password) throws  TGException
    {
        Map<String,String> emptyMap = Collections.emptyMap();
        return createChannel(urlPath, userName, password, emptyMap);
    }



    @Override
    public  TGChannel createChannel(String urlPath, String userName, String password, Map<String,String> props ) throws  TGException
    {

        TGProperties<String,String> properties;

        TGChannelUrl url = LinkUrl.parse(urlPath);


        properties = new SortedProperties(String.CASE_INSENSITIVE_ORDER);  //SS:All property names are case insensitive.
        properties.putAll(url.getProperties());

        if (props != null) {
            properties.putAll(props); //Override the value from userdefined props.
        }


        setUserAndPassword(properties, userName, password);


        TGChannelUrl.Protocol  protocol = url.getProtocol();

        switch(protocol) {
            case SSL:
                return new SslChannel(url, properties);

            case TCP:
                return new TcpChannel(url, properties);


            default:
                throw new TGException("Protocol not supported");
        }

    }

    private void setUserAndPassword(TGProperties<String,String> properties, String userName, String password) throws TGException
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
}
