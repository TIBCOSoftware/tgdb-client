package com.tibco.tgdb.channel;

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
 *
 * File name : TGChannel
 * Created by: Suresh
 * SVN Id : $Id: TGChannel.java 2155 2018-03-19 04:31:41Z ssubrama $
 */

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGMessage;

import java.util.Map;

/**
 *
 */
public interface TGChannel {

    public enum LinkState {
        NotConnected,
        Connected,
        Closing,
        Closed,
        FailedOnSend,
        FailedOnRecv,
        FailedOnProcessing,
        Reconnecting,
        Terminated,
    }

    public enum ResendMode {
        DontReconnectAndIgnore,
        ReconnectAndResend,
        ReconnectAndRaiseException,
        ReconnectAndIgnore,
    }

    public enum ReconnectState {
        ReconnectChannelClosed,
        ReconnectSuccess,
        ReconnectFailed,
        ReconnectFailedAllAttempts,

    }



    public interface LinkEventHandler {
        public void     onException(Exception ex, boolean duringClose);
        public boolean  onReconnect();
        public String   getTerminatedText();
    }

    /**
     * Get the Link/Channel State
     * @return
     */
    LinkState getLinkState();

    /**
     * Get the Channel Properties
     * @return
     */
    Map getProperties();

    /**
     * Get the server protocol version
     * @return
     */
    int getServerProtocolVersion();

    /**
     * Get the client protocol version
     * @return
     */
    int getClientProtocolVersion();

    /**
     * Get Authorization Token
     * @return
     */
    long getAuthToken();
    
    /**
     * Get Session Id
     * @return
     */
    long getSessionId();

    /**
     * Connect the underlying Channel using the URL
     * @throws TGException
     */

    void connect() throws TGException;

    /**
     * Start the channel for send and recving messages. Set the LinkEventHandler before starting the Channel.
     *
     * @throws TGException
     */
    void start()   throws TGException;


    /**
     * Disconnect the Channel from its endpoint
     */
    void disconnect();

    /**
     * Stop the channel forcefully or gracefully.
     * @param bForcefully
     */
    void stop(boolean bForcefully);


    /**
     * Send a Message on this channel, and returns immediately.
     * @param msg
     * @throws TGException
     */
    public void sendMessage(TGMessage msg) throws TGException;

    /**
     * Send a Message waiting for a response. The thread blocks till it gets the response.
     * @param msg
     * @param response
     * @return
     * @throws TGException
     */
    public TGMessage sendRequest(TGMessage msg, TGChannelResponse response) throws TGException;

}
