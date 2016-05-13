package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.utils.TGProperties;

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
 * File name :SslChannel
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: SslChannel.java 583 2016-03-15 02:02:39Z vchung $
 */
public class SslChannel extends TcpChannel {

    protected SslChannel(TGChannelUrl linkUrl, TGProperties props)
    {
        super(linkUrl, props);

    }







}
