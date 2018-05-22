package com.tibco.tgdb.channel;

import java.util.List;
import java.util.Map;
import java.util.Queue;

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
 * File name :TGChannelUrl
 * Created on: 12/25/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGChannelUrl.java 2203 2018-04-04 01:58:36Z ssubrama $
 */
public interface TGChannelUrl {

    public enum Protocol {
        TCP,
        SSL,
        HTTP,
        HTTPS
    }

    /**
     * Get the Protocol
     * @return
     */
    Protocol getProtocol();

    /**
     * Get the host part of the URL
     * @return
     */
    String getHost();

    /**
     * Get the User associated with the URL
     * @return
     */
    String getUser();

    /**
     * Get the port to which it is connected
     * @return
     */
    int getPort();


    /**
     * Get the URL Properties
     * @return
     */
    Map<String,String> getProperties();

    /**
     * Get the Fault Tolerant URLs
     * @return
     */
    List<TGChannelUrl> getFTUrls();

    /**
     * Get the String form of the URL.
     * @return
     */
    String getUrlAsString();





}
