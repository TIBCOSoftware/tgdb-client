package com.tibco.tgdb.channel.impl;

import java.io.IOException;
import java.net.*;
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
 * File name :DefaultSocketImpl
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: DefaultSocketImpl.java 583 2016-03-15 02:02:39Z vchung $
 */
public class DefaultSocketImpl extends AbstractSocket {

    DefaultSocketImpl(Map<String, String> props) {
        super(props);
    }

    Socket createSocket(String host, int port, int timeoutMillis)
            throws IOException
    {

        /*
         * Loop through IPV6 Addresses and IPV4 Address, and connect them.
         */
        InetAddress addresses[];
        InetAddress         addr;
        InetSocketAddress sockAddr;
        Socket              s = null;
        addresses = InetAddress.getAllByName(host);
        InetAddress[]        picks = new InetAddress[2];
        int numPicks = 0;

        for (int i=0; i<addresses.length; i++)
        {
            addr = addresses[i];
            if ((addr instanceof Inet6Address))
            {
                picks[numPicks++] = addr;
                break;
            }
        }
        for (int i=0; i<addresses.length; i++)
        {
            addr = addresses[i];
            if ((addr instanceof Inet4Address))
            {
                picks[numPicks++] = addr;
                break;
            }
        }

        for (int i=0; i<numPicks; i++) {
            addr = picks[i];

            s = new Socket();

            // before connect() for TCP Window Scaling to be turned on
            setBuffers(s);


            sockAddr = new InetSocketAddress(addr, port);
            try {
                s.connect(sockAddr,timeoutMillis);
                break;
            } catch (IOException e) {
                s.close();
                s = null;
                if (i == numPicks - 1)
                    throw e;
            }
        }

        if (s == null)
            throw new IOException("No supported IP address found");

        return s;
    }
}
