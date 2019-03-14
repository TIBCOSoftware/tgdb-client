package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelFactory;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.SortedProperties;
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
 * SVN Id: $Id: ChannelFactoryImpl.java 2214 2018-04-05 18:21:28Z ssubrama $
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

        TGProperties<String,String> properties = new SortedProperties<String,String>(String.CASE_INSENSITIVE_ORDER);  //SS:All property names are case insensitive.
        if (props != null) {
            properties.putAll(props); //Override the value from userdefined props.
        }

        //System Properties are overridden by the URL properties.
        TGChannelUrl url = LinkUrl.parse(urlPath);
        properties.putAll(url.getProperties());

        //The API parameters override the URL Properties
        TGProperties.setUserAndPassword(properties, userName, password);

        return createChannel(url, properties);

    }

    public TGChannel createChannel(TGChannelUrl url, TGProperties<String, String> properties) throws TGException
    {
        TGChannelUrl.Protocol  protocol = url.getProtocol();

        switch(protocol) {
            case SSL:
                return new SSLChannel(url, properties);

            case TCP:
                return new TcpChannel(url, properties);

            default:
                throw new TGException("Protocol not supported");
        }
    }


}
