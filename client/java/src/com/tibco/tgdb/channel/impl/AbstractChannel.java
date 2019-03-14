package com.tibco.tgdb.channel.impl;

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
 * File name :AbstractChannel
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: AbstractChannel.java 2742 2018-11-15 22:17:42Z nimish $
 */

import com.tibco.tgdb.TGProtocolVersion;
import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGChannelDisconnectedException;
import com.tibco.tgdb.exception.TGConnectionTimeoutException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.DisconnectChannelRequest;
import com.tibco.tgdb.pdu.impl.ExceptionMessage;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.HexUtils;
import com.tibco.tgdb.utils.TGConstants;
import com.tibco.tgdb.utils.TGProperties;

import java.io.IOException;
import java.util.Deque;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.AtomicReference;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import static com.tibco.tgdb.channel.impl.AbstractChannel.ExceptionHandleResult.Disconnected;
import static com.tibco.tgdb.channel.impl.AbstractChannel.ExceptionHandleResult.RethrowException;
import static com.tibco.tgdb.channel.impl.AbstractChannel.ExceptionHandleResult.RetryOperation;

/**
 * Physical link between the Client and the Server. A Connection is a logical abstract which uses the channel.
 */
public abstract class AbstractChannel implements TGChannel {



    enum ExceptionHandleResult {
        RethrowException(TGException.ExceptionType.GeneralException, "TGDB-CHANNEL-FAIL:Failed to reconnect"),
        RetryOperation(TGException.ExceptionType.RetryIOException, "TGDB-CHANNEL-RETRY:Channel Reconnected, Retry Operation"),
        Disconnected(TGException.ExceptionType.DisconnectedException, "TGDB-CHANNEL-FAIL:Failed to reconnect");

        ExceptionHandleResult(TGException.ExceptionType type, String msg) {
            this.type = type;
            this.msg = msg;
        }
        TGException.ExceptionType type;
        String msg;
    }


    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();


    final Lock sendLock = new ReentrantLock();
    final Lock exceptionLock = new ReentrantLock();
    final Condition exceptionCond = exceptionLock.newCondition();

    private int connectionIndex;
    TGChannelUrl        linkURL             = null;
    TGChannelUrl        primaryURL          = null;
    TGProperties<String, String> properties = null;
    AtomicReference<LinkState> state        = new AtomicReference<LinkState>(LinkState.NotConnected);

    AtomicBoolean       bNeedsPing          = new AtomicBoolean(false);
    AtomicInteger       numConnections      = new AtomicInteger(0);  //number of logical connection
    AtomicLong 			requestId = 		new AtomicLong(); //FIXME:  requestId is per channel. How should we handle this?
    long                lastActiveTime      = 0;  //used for server-client HBs.

    ChannelReader       reader              = null;
    ConcurrentMap<Long, TGChannelResponse>  responses    = new ConcurrentHashMap<>();

    String              clientId            = null;
    String              inboxAddr           = null;
    long                sessionId           = -1;
    long                authToken           = -1;



    protected AbstractChannel(TGChannelUrl linkUrl, TGProperties<String, String> props) throws TGException {
        this.linkURL = linkUrl;
        this.primaryURL = linkUrl;
        this.properties = props;
        reader = new ChannelReader(this);
        this.connectionIndex = 0;
    }

    abstract void createSocket() throws TGException;

    abstract void onConnect() throws TGException;

    abstract void closeSocket();


    /**
     * Send Message to the server, compress and or encrypt.
     * Hence it is abstraction, that the Channel knows about it.
     * @param msg       The message that needs to be sent to the server
     * @throws TGException, IOException
     */
    abstract void send(TGMessage msg) throws TGException, IOException;

    /**
     *
     * @return TGMessage
     * @throws TGException
     */
    abstract TGMessage readWireMsg() throws TGException, IOException;


    /**
     * --------------------------------------------------------------
     * return true if connected
     * --------------------------------------------------------------
     */
    public boolean isConnected() {
        return this.getLinkState() == LinkState.Connected;
    }

    /**
     * --------------------------------------------------------------
     * return true if closed
     * --------------------------------------------------------------
     */
    public boolean isClosed() {
        return (
                this.getLinkState() == LinkState.Closed  ||
                this.getLinkState() == LinkState.Closing ||
                this.getLinkState() == LinkState.Terminated
        );
    }

    protected boolean needsPing() { return bNeedsPing.get(); }

    protected void disablePing() { bNeedsPing.set(false); }

    protected void enablePing() { bNeedsPing.set(true); }

    protected String getClientId() {
     /*SSB: Fix for JIRA 356*/
     return this.clientId;
    }

    protected String getInboxAddr() { return inboxAddr;}

    protected String getUserName() { return properties.getProperty(ConfigName.ChannelUserID);}

    //Password is encrypted in properties
    protected byte[] getPassword() {

        String s =  properties.getProperty(ConfigName.ChannelPassword);
        return s == null ? TGConstants.EmptyByteArray : s.getBytes();
    }

    protected String getHost() {
        return linkURL.getHost();
    }

    protected int getPort() {
        return linkURL.getPort();
    }

    private void _connect(boolean sleepOnFirstInvocation) throws TGException
    {
        int connectInterval =  Integer.parseInt(properties.getProperty(ConfigName.ChannelFTRetryIntervalSeconds));
        int retryCount      =  Integer.parseInt(properties.getProperty(ConfigName.ChannelFTRetryCount));
        List<TGChannelUrl> ftUrls =  this.primaryURL.getFTUrls();
        int idx = this.connectionIndex;
        int urlCount = ftUrls.size();

        do {
            this.linkURL = ftUrls.get(idx);
            for (int i=0; i<retryCount; i++) {
                try {
                    gLogger.log(TGLogger.TGLevel.Info, "Attempt:%d to connect to url:%s", i, this.linkURL.getUrlAsString());
                    if (sleepOnFirstInvocation) {
                        Thread.sleep(connectInterval * 1000);
                        sleepOnFirstInvocation = false;
                    }
                    createSocket();
                    onConnect();
                    this.connectionIndex = idx;
                    return;
                } catch (TGAuthenticationException | TGChannelDisconnectedException te) {
                    throw te;
                } catch (TGException tge) {
                    if (null == properties.getProperty(ConfigName.ChannelFTHosts))
                    {
                    	gLogger.logException(String.format("Failed connecting to urlstr:%s", this.linkURL), tge);
                    	closeSocket();
                    	throw tge;
                    }
                    else {
                        gLogger.logException(String.format("Failed connecting to urlstr:%s, reattempting", this.linkURL), tge);
                        closeSocket();
                    }
                } catch (Exception e) {
                    gLogger.logException(String.format("Failed connecting to urlstr:%s, reattempting", this.linkURL), e);
                    closeSocket();
                }
            }
            idx = (idx + 1) % urlCount;

        } while (idx != this.connectionIndex);

        throw new TGConnectionTimeoutException(String.format("%s:Failed %d attempts to connect to TGDB Server.", "TGDB-CONNECT-ERR", retryCount));

    }
    public final synchronized void connect() throws TGException {
        if (isConnected()) {
            numConnections.incrementAndGet();
            return;
        }

        if ((isClosed()) || (state.get() == LinkState.NotConnected)) {
//            createSocket();
//            onConnect();
            _connect(false);
        }
        else {
            throw new TGException("Connect called on an invalid state :=" + state.get().name());
        }
        state.set(LinkState.Connected);
        numConnections.incrementAndGet();
        return;
    }

    public synchronized void start() throws TGException {
        if (isConnected()) {
            enablePing();
            reader.start();
            MultiChannelPinger.getInstance().addChannel(this);
            MultiChannelPinger.getInstance().start();
        }
        else {
            throw new TGException("Channel not connected");
        }
    }

    public void stop() {
        stop(false);
    }

    public void stop(boolean bForcefully) {

        sendLock.lock(); //Ensure nobody can send.

        try {
            if (!isConnected()) {
                return;
            };

            if ((bForcefully) || (numConnections.get() == 0))
            {
                gLogger.log(TGLogger.TGLevel.Debug, "Stopping channel");
                
                disablePing();
                reader.stop();

                // Send the disconnect request. sendRequest will not receive a channel response since the channel will be disconnected.
                DisconnectChannelRequest request = (DisconnectChannelRequest) TGMessageFactory.getInstance().createMessage(VerbId.DisconnectChannelRequest);
                this.send(request);
                state.set(LinkState.Closing);
                closeSocket();
                MultiChannelPinger.getInstance().removeChannel(this);
                MultiChannelPinger.getInstance().stop();
            }
        } catch (Exception ioe) {
            gLogger.logException("DisconnectChannelRequest send failed", ioe);
            closeSocket();
        }

        finally {
        	if (isClosing()) {
        		state.set(LinkState.Closed);
        	}
            sendLock.unlock();
        }
    }

    //SS:TODO
    //SS:This needs to be revisited
     boolean reconnect() {
        // This is needed here to avoid a FD leak
        closeSocket();
        
        if (null == properties.getProperty(ConfigName.ChannelFTHosts))
        {
        	return false;
        }

        TGChannelUrl oldurl    = this.linkURL;
        int connectInterval =  Integer.parseInt(properties.getProperty(ConfigName.ChannelFTRetryIntervalSeconds));
        int retryCount      =  Integer.parseInt(properties.getProperty(ConfigName.ChannelFTRetryCount));
        gLogger.log(TGLogger.TGLevel.Info, "Retrying to reconnnect %d times at interval of %d seconds to FTUrls.", retryCount, connectInterval);

        this.state.set(LinkState.Reconnecting);
        MultiChannelPinger.getInstance().removeChannel(this);
        try {
          _connect(true);
            this.state.set(LinkState.Connected);
            MultiChannelPinger.getInstance().addChannel(this);
            return true;
        }
        catch (TGException e) {
            this.linkURL = (LinkUrl) oldurl;
            this.state.set(LinkState.Closed);
            return false;
        }
    }

    public final synchronized void disconnect() {

        sendLock.lock();  //Ensure nobody is sending when we are disconnecting.
        if (!isConnected()) {
        	return;
        }
        try {
            if (numConnections.get() == 0) {
                gLogger.log(TGLogger.TGLevel.Error, "Calling disconnect more than the number of connects.");
                return;
            }
            numConnections.decrementAndGet();
        }
        catch (Exception e) {
            gLogger.logException("Channel disconnect failed", e);
        }
        finally {
            sendLock.unlock();
        }
    }


    public void sendMessage(TGMessage msg, boolean resend) throws TGException
    {
        ResendMode resendmode = resend ? ResendMode.ReconnectAndResend : ResendMode.ReconnectAndRaiseException;

        while (true) {
            try {
                if (!isConnected()) throw new TGException("Channel Closed", "TGDB-CHANNEL-ERR" );
                sendLock.lock();
                send(msg);
                return;
            }

            catch (Exception e) {
                ExceptionHandleResult ehr = handleException(e, false);
                if (ehr == RethrowException) {
                    if (e instanceof TGException) throw (TGException)e;
                    throw TGException.buildException("Failed to send message", "TGDB-SND-ERR", e);
                }
                else if (ehr == Disconnected) {
                    throw new TGChannelDisconnectedException(e);
                }
                else {
                    gLogger.log(TGLogger.TGLevel.Info, "Retrying to send message on urlstr:%s", this.linkURL);
                }
            }
            finally {
                sendLock.unlock();
            }

        }

    }

    public void sendMessage(TGMessage msg) throws TGException
    {
        sendMessage(msg, true);
    }

    protected TGMessage requestReply(TGMessage request) throws TGException
    {
        while (true) {
            try {
                sendLock.lock();
                this.send(request);
                TGMessage msg = readWireMsg();
                return msg;
            } catch (Exception e) {
                ExceptionHandleResult ehr = handleException(e, true);
                if (ehr == RethrowException) {
                    if (e instanceof TGException) throw (TGException) e;
                    throw TGException.buildException("Failed to send message", "TGDB-SND-ERR", e);
                } else if (ehr == Disconnected) {
                    throw new TGChannelDisconnectedException(e);
                } else {
                    gLogger.log(TGLogger.TGLevel.Info, "Retrying to send message on urlstr:%s", this.linkURL);
                }
            } finally {
                sendLock.unlock();
            }
        }
    }
    public TGMessage sendRequest(TGMessage msg, TGChannelResponse response) throws TGException {
        return sendRequest(msg, response, true);
    }

    /**
     *
     * @param request
     * @param channelResponse
     * @param resend - A boolean flag that emulates resending properties if the send fails. If true, will try to
     *               reconnect and attempt to resend
     *               If false, will try to reconnect and throw an exception.
     * @return
     * @throws TGException
     *
     */
    public TGMessage sendRequest(TGMessage request, TGChannelResponse channelResponse, boolean resend) throws TGException {

        final long key = channelResponse.getRequestId(); //FIXME: How do we get request id before we send out an request??
        request.setRequestId(key);

        while(true) {
            try {
                if (!isConnected()) throw new TGException("Connection is closed", "TGDB-CHANNEL-ERR");
                sendLock.lock();
                responses.put(key, channelResponse);
                send(request);
                if (!channelResponse.isBlocking()) return null;
                channelResponse.await(p -> p == TGChannelResponse.Status.Waiting);
                responses.remove(key);
                TGMessage msg =  channelResponse.getReply();
                if (msg instanceof ExceptionMessage) {
                    ExceptionMessage exMsg = (ExceptionMessage) msg;
                    if (exMsg.getExceptionType() == TGException.ExceptionType.RetryIOException) continue;
                    throw TGException.buildException(exMsg);
                }
                return msg;
            } catch (Exception e) {
                ExceptionHandleResult ehr = handleException(e, false);
                if (ehr == RethrowException) {
                    if (e instanceof TGException) throw (TGException) e;
                    throw TGException.buildException("Failed to send message", "TGDB-SND-ERR", e);
                } else if (ehr == Disconnected) {
                    throw new TGChannelDisconnectedException(e);
                } else {
                    gLogger.log(TGLogger.TGLevel.Info, "Retrying to send message on urlstr:%s", this.linkURL);
                }
            } finally {
                sendLock.unlock();
            }
        }


    }

    /**
     * process a message received on the channel. This is called from the ChannelReader.
     * @param msg
     * @throws TGException
     */
    protected void processMessage(TGMessage msg) throws TGException {
        long requestId = msg.getRequestId();
        TGChannelResponse channelResponse = responses.get(requestId);

        if (channelResponse == null) {
            gLogger.log(TGLogger.TGLevel.Error, "Received no response message for corresponding request :%d", requestId);
            return;
        }
        try {
        	if (gLogger.isEnabled(TGLogger.TGLevel.Debug))
        	    gLogger.log(TGLogger.TGLevel.Debug, "Process msg: %s", HexUtils.formatHex(msg.toBytes()));
        } catch (IOException ioe) {}
        channelResponse.setReply(msg);

        return;
    }

    protected ExceptionHandleResult handleException(Exception ex, boolean bReconnect)
    {
        int connectionOpTimeout =  Integer.parseInt(properties.getProperty(ConfigName.ConnectionOperationTimeoutSeconds));
        if (!(ex instanceof IOException)) return RethrowException;
        try {

            final ReentrantLock lock = (ReentrantLock) this.exceptionLock;
            lock.lock();
            while (!bReconnect) {
                try {
                    if (exceptionCond.await(connectionOpTimeout, TimeUnit.SECONDS)) {
                        if (isConnected()) return RetryOperation;
                    }
                    if (isClosed()) return Disconnected;
                }
                catch (InterruptedException ie) {}
            }
            if(reconnect()) return RetryOperation;
            return Disconnected;
        }

        finally {
            if (bReconnect) this.exceptionCond.signalAll();
            this.exceptionLock.unlock();
        }
    }

    void terminated(String killMsg)
    {
        exceptionLock.lock();
        try {
            this.state.set(LinkState.Terminated);
            closeSocket();
            gLogger.log(TGLogger.TGLevel.Error, String.format("Session killed by :%s", killMsg));
        }
        finally {
            exceptionLock.unlock();
        }
    }

    protected boolean isClosing() { return this.getLinkState() == LinkState.Closing; }

    @Override
    public LinkState getLinkState() {
        return state.get();
    }

    @Override
    public Map getProperties() {
        return properties;
    }

    @Override
    public int getServerProtocolVersion() {
        return 0;
    }

    @Override
    public int getClientProtocolVersion() {
        return TGProtocolVersion.getProtocolVersion();
    }

    @Override
    public long getAuthToken() {
    	return authToken;
    }
    
    @Override
    public long getSessionId() {
    	return sessionId;
    }
}
