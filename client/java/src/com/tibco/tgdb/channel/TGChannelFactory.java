package com.tibco.tgdb.channel;

import com.tibco.tgdb.channel.impl.ChannelFactoryImpl;
import com.tibco.tgdb.exception.TGException;

import java.io.IOException;
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
 * File name :TGChannelFactory
 * Created on: 12/25/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGChannelFactory.java 2135 2018-03-07 23:42:34Z ssubrama $
 */
public abstract class TGChannelFactory {

    private static TGChannelFactory gFactory = new ChannelFactoryImpl();

    /**
     * Get an instance of the Channel Factory
     * @return
     */
    public static TGChannelFactory getInstance() {
        return gFactory;
    }


    /**
     * Create a channel on the URL specified using the userName and password.
     * A URL is represented as a string of the form
     * <protocol>://[user@]['['ipv6']'] | ipv4 [:][port][/]'{' name:value;... '}'
     * @param urlPath A url string.
     * @param userName The userName for the channel. The userId provided overrides all other userIds that can be infered.
     *               The rules for overriding are in this order
     *               a. The argument 'userId' is the highest priority. If Null then
     *               b. The user@url is considered. If that is Null
     *               c. the "userID=value" from the URL string is considered.
     *               d. If all of them is Null, then the default User associated to the installation will be taken.
     *
     * @param password An encrypted password associated with the userName
     * @return a Channel
     * @throws IOException
     * @throws
     * @see com.tibco.tgdb.channel.TGChannelUrl for detailed URL specification
     */
    public abstract TGChannel createChannel(String urlPath, String userName, String password) throws  TGException;


    /**
     * Create a channel on the URL specified using the user Name and password
     * @param urlPath A url as a string form
     * @param userName The userName for the channel. The userId provided overrides all other userIds that can be infered.
     *               The rules for overriding are in this order
     *               a. The argument 'userId' is the highest priority. If Null then
     *               b. The user@url is considered. If that is Null
     *               c. the "userID=value" from the URL string is considered.
     *               d. The user retrieved from the Properties is considered
     *               e. If all of them is Null, then the default User associated to the installation will be taken.
     * @param password Encrypted password
     * @param props A properties bag with Connection Properties. The URL infered properties override this property bag.
     * @return a connected channel
     * @throws IOException
     * @throws TGException
     */
    public abstract  TGChannel createChannel(String urlPath, String userName, String password, Map<String,String> props ) throws  TGException;

}
