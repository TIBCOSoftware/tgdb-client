package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGEnvironment;

import java.io.IOException;
import java.net.Socket;
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
 * File name :AbstractSocket
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: AbstractSocket.java 583 2016-03-15 02:02:39Z vchung $
 */
public abstract class AbstractSocket {

    Map<String, String> props;

    AbstractSocket(Map<String, String> props) {
        this.props = props;
    }

    abstract Socket createSocket(String url, int port, int timeoutMillis) throws IOException;


    protected void setBuffers(Socket s)
    {

        int sendSize = Integer.parseInt(props.getOrDefault(ConfigName.ChannelSendSize.getName(), "-1"));

        if (sendSize == -1) {
            sendSize = Integer.parseInt(props.getOrDefault(ConfigName.ChannelSendSize.getAlias(),
                    String.valueOf(TGEnvironment.getInstance().getChannelSendSize())));
        }

        int receiveSize = Integer.parseInt(props.getOrDefault(ConfigName.ChannelRecvSize.getName(), "-1"));

        if (receiveSize == -1) {
            receiveSize = Integer.parseInt(props.getOrDefault(ConfigName.ChannelRecvSize.getAlias(),
                    String.valueOf(TGEnvironment.getInstance().getChannelReceiveSize())));
        }
        // Hmm...Oracle JVM does not implement these methods and throws an exception.
        // NB: size == -1 means don't call these methods
        if (sendSize > 0)
        {
            try { s.setSendBufferSize(sendSize*1024);  }
            catch(Throwable ignore) {}
        }

        if (receiveSize > 0)
        {
            try { s.setReceiveBufferSize(receiveSize*1024);  }
            catch(Throwable ignore) {}
        }
    }
}
