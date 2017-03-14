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
 * SVN Id: $Id: AbstractChannel.java 1345 2017-02-03 03:33:02Z vchung $
 */

import com.tibco.tgdb.TGProtocolVersion;
import com.tibco.tgdb.channel.TGChannel;
import com.tibco.tgdb.channel.TGChannelResponse;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.DisconnectChannelRequest;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGConstants;
import com.tibco.tgdb.utils.TGProperties;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Physical link between the Client and the Server. A Connection is a logical abstract which uses the channel.
 */
public abstract class AbstractChannel implements TGChannel {

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();


    final Lock sendLock = new ReentrantLock();

    LinkUrl             linkURL             = null;
    LinkEventHandler    linkHandler         = null;
    TGProperties<String, String> properties = null;
    LinkState           state               = LinkState.NotConnected;

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


    protected AbstractChannel(TGChannelUrl linkUrl, TGProperties<String, String> props) {
        this.linkURL = (LinkUrl) linkUrl;
        this.properties = props;
        reader = new ChannelReader(this);
    }

    abstract void createSocket() throws TGException;

    abstract void onConnect() throws TGException;

    abstract void closeSocket();


    /**
     * Send Message to the server with the Resend Mode flag set
     *
     * @param data       The message that needs to be sent to the server

     * @throws TGException, IOException
     */
    abstract void send(TGMessage data) throws TGException, IOException;

    /**
     *
     * @return TGMessage
     * @throws TGException
     */
    abstract TGMessage readWireMsg() throws TGException, IOException;

    @Override
    public void setLinkEventHandler(LinkEventHandler handler) {
        this.linkHandler = handler;
        lastActiveTime = System.currentTimeMillis();
    }

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
        return this.getLinkState() == LinkState.Closed;
    }

    protected boolean needsPing() { return bNeedsPing.get(); }

    protected void disablePing() { bNeedsPing.set(false); }

    protected void enablePing() { bNeedsPing.set(true); }

    protected String getClientId() { return properties.getProperty(ConfigName.ChannelClientId);}

    protected String getInboxAddr() { return inboxAddr;}

    protected String getUserName() { return properties.getProperty(ConfigName.ChannelUserID);}

    //Password is encrypted in properties
    protected byte[] getPassword() {

        String s =  properties.getProperty(ConfigName.ChannelPassword);
        return s == null ? TGConstants.EmptyByteArray : s.getBytes();
    }

    public final synchronized void connect() throws TGException {
        if (isConnected()) {
            numConnections.incrementAndGet();
            return;
        }

        if ((isClosed()) || (state == LinkState.NotConnected)) {
            createSocket();
            onConnect();
        }
        else {
            throw new TGException("Connect called on an invalid state :=" + state.name());
        }

        state = LinkState.Connected;
        numConnections.incrementAndGet();
        return;
    }

    public synchronized void start() throws TGException {
        if (state == LinkState.Connected) {
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

                // Send the disconnect request.
                DisconnectChannelRequest request = (DisconnectChannelRequest) TGMessageFactory.getInstance().createMessage(VerbId.DisconnectChannelRequest);
                // sendRequest will not receive a channel response since the channel will be disconnected.
                this.send(request);

                state = LinkState.Closing;
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
        		state = LinkState.Closed;
        	}
            sendLock.unlock();
        }
    }

    //SS:TODO
    //SS:This needs to be revisited
    public final synchronized boolean reconnect() {
        // This is needed here to avoid a FD leak
        closeSocket();

        switch (state) {
            case Closed:
            case Closing:
                return false;

            case FailedOnSend:
            case FailedOnRecv:
                return linkHandler.onReconnect();

        }

        return false;
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

    //---------------------------------------------------------------
    // sendRequestMsg
    //---------------------------------------------------------------

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
        ResendMode resendmode = ResendMode.ReconnectAndRaiseException;

        final long key = channelResponse.getRequestId(); //FIXME: How do we get request id before we send out an request??
        request.setRequestId(key);

        if (resend)
            resendmode = ResendMode.ReconnectAndResend;

        TGException exception ;

        while(true) {
            try {
                exception = null;
                if (state != LinkState.Connected)
                    throw new TGException("Connection is closed");

                responses.put(key,channelResponse);

                //We set it everytime, because, if it reconnected, these values will be changed.
                //request.setAuthToken(this.authToken);
                //request.setSessionId(this.sessionId);
                //request.setClientId(this.clientId);
                //request.getRequestId();

                send(request);
            }
            catch(IOException ioe) {
            	if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            		gLogger.logException("sendRequest failed with IOException", ioe);
            	}
                exception = TGException.buildException("Failed to send due to IO issues", null, ioe);
                responses.remove(key);

                if (isClosed() || isClosing()) throw exception;

                state = LinkState.FailedOnSend;

                if (!reconnect()) {
                    exception =  TGException.buildException("Failed to send & reconnect due to IO issues", null, ioe);
                }

                if (resendmode == ResendMode.ReconnectAndRaiseException) {
                    throw exception;
                }

                // else loop back and resend
            }
            catch(TGException e) {
            	if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            		gLogger.logException("sendRequest failed", e);
            	}
                responses.remove(key);
                throw e;
            }

            if (exception == null) {  //We didnt have any exception
                //Non Blocking Channel - return
                if (!channelResponse.isBlocking())
                    return null;

                try {
                    channelResponse.await(p -> p == TGChannelResponse.Status.Waiting);
                } catch (InterruptedException e) {
                	if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
                		gLogger.logException("Channel response wait interrupted", e);
                	}
                    exception = TGException.buildException("InterruptedException has occurred while waiting for server response", null, e);
                }

                responses.remove(key);

                if (channelResponse.getStatus() == TGChannelResponse.Status.Resend) {
                    if (!resend) {
                        exception = TGException.buildException("Send failed due to fault-tolerant switch", null, null);
                    }
                }

                if (exception != null) {
                    throw exception;
                }

                // in any case we return reply, it'll be null if
                // timeout, connection closed or we were pushed
                return channelResponse.getReply();
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
        	if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) gLogger.log(TGLogger.TGLevel.Debug, "Process msg: %s", msg.toBytes());
        } catch (IOException ioe) {
        	
        }
        channelResponse.setReply(msg);

        return;
    }

    protected int handleException(Exception ex) {
        try {
            sendLock.lock();

            if (isClosed() && ex instanceof IOException) {
                if (reconnect()) return 1;
            }
            gLogger.logException("Aborting channel due to exception", ex);
            state = LinkState.Closed;

            //Notify all blocked Threads of the Status
            for (TGChannelResponse response : responses.values()) {
                response.signal(TGChannelResponse.Status.Disconnected);
            }
            if (linkHandler != null) linkHandler.onException(ex, true);
        }
        finally {
            sendLock.unlock();
        }
        return 0;
    }

    protected boolean isClosing() { return this.getLinkState() == LinkState.Closing; }

    @Override
    public LinkState getLinkState() {
        return state;
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
