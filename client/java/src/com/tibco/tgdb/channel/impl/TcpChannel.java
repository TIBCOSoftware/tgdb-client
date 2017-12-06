package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.*;
import com.tibco.tgdb.utils.*;

import java.io.*;
import java.net.Socket;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

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
 * File name :TcpChannel
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TcpChannel.java 1880 2017-11-04 02:29:10Z vchung $
 */
public class TcpChannel extends AbstractChannel {

    final static Lock shutdownLock = new ReentrantLock();

    Socket                  socket         = null;
    DataInputStream         input          = null;
    DataOutputStream        output         = null;
    // set while reconnect attempt is in progress
    volatile boolean        reconnecting   = false;


    /**
     * Construct a Tcp channel on a URL, with a set properties
     *
     * @param linkUrl
     * @param props
     */

    protected TcpChannel(TGChannelUrl linkUrl, TGProperties<String, String> props) {
        super(linkUrl, props);
    }

    /**
     * Create the physical link socket
     *
     * @throws TGException
     */
    void createSocket() throws TGException {
        String host = linkURL.getHost();
        int port = linkURL.getPort();
        String failureMessage = "Failed to connect to the server at " + linkURL.url;

        state = LinkState.NotConnected;

        try {
            synchronized (shutdownLock) {
                socket = null;

                AbstractSocket ts = new DefaultSocketImpl(this.properties);
                socket = ts.createSocket(host, port, TGEnvironment.getInstance().getChannelConnectTimeout());
                socket.setTcpNoDelay(true);
                socket.setSoLinger(false, 0);
                input = new DataInputStream(new BufferedInputStream(socket.getInputStream(), 1024 * 32));
                output = new DataOutputStream(socket.getOutputStream());
                clientId = properties.get(ConfigName.ChannelClientId.getName());
                if (clientId == null) {
                    clientId = properties.get(ConfigName.ChannelClientId.getAlias());
                    if (clientId == null) {
                        clientId = TGEnvironment.getInstance().getChannelClientId();
                    }
                }
                inboxAddr = socket.getInetAddress().toString(); //SS:TODO: Is this correct
            }
        } catch (IOException ioex) {
            if (socket != null) closeSocket();
            throw TGException.buildException(failureMessage, null, ioex);
        }
    }


    /**
     * Close the socket.
     */
    void closeSocket() {
        //Flush the output.
        try {
            sendLock.lock();
            if (output != null) {
                output.flush();
                output.close();
            }
        } catch (Exception e) {
        } //ignore the error
        finally {
            sendLock.unlock();
        }
        Socket s = socket;
        if (s != null) {
            try {s.shutdownInput();}catch(Exception e){}catch(NoSuchMethodError e){}
            try {s.shutdownOutput();}catch(Exception e){}catch(NoSuchMethodError e){}
            try {s.close();}catch(Exception e){}
        }
        socket = null;
        output = null;
        input  = null;
    }

    /**
     * On connect called by the AbstractChannel
     * @throws TGException
     */

    public void onConnect() throws TGException {
        performHandshake(false);
        doAuthenticate();
    }

    public boolean isConnected() {
        return (state == LinkState.Connected);
    }

    public boolean isClosed() {
        return (state == LinkState.Closed && !reconnecting);
    }

    /**
     * Read wire message called from the LinkReader
     * @return
     * @throws IOException
     */
    protected TGMessage readWireMsg() throws TGException, IOException {
        DataInputStream in      = input;
        byte[]          buffer  = TGConstants.EmptyByteArray;

        if (in == null)
            return null;

        disablePing();


        if (state == LinkState.Closed)
            return null;

        try {
            lastActiveTime = System.currentTimeMillis();
            int size = in.readInt();
            if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            	gLogger.log(TGLogger.TGLevel.Debug, "readWireMsg incoming buffer size : %d", size);
            }
            if ((size > 0) && (in.available() > 0)) {
                buffer = new byte[size];
                in.readFully(buffer, 4, size - 4);

            } else {
            	gLogger.log(TGLogger.TGLevel.Debug, "readWireMsg data input stream return with no data");
                throw new EOFException();
            }
            intToBytes(size, buffer, 0);
            if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            	gLogger.log(TGLogger.TGLevel.Debug, "readWireMsg : %s", HexUtils.formatHex(buffer));
            }
            TGMessage msg = TGMessageFactory.getInstance().createMessage(buffer, 0, size);

            //FIXME: channel level exception??  otherwise, command error returns as command response?
            if (msg.getVerbId() == VerbId.ExceptionMessage) {
                throw TGException.buildException((ExceptionMessage) msg);
            }
            return msg;
        } catch (TGException  be) {
            //gLogger.logException(String.format("readWireMsg TGException : %s(url=%s)", be.toString(), linkURL.toString()), be);
            logMessage(buffer);
            throw be;
        } catch (IOException ie) {
            //gLogger.logException(String.format("readWireMsg IOException : %s(url=%s)", ie.toString(), linkURL.toString()), ie);
            logMessage(buffer);
            throw ie;
        }


    }

    void intToBytes(int value, byte[] bytes, int offset) {
        for (int i=0; i<4; i++)
            bytes[offset+i] = (byte)((value>>>(8*(3-i)))&0xff);
    }


    /**
     * performHandshake
     *
     * @param sslMode
     * @throws TGException
     */
    protected void performHandshake(boolean sslMode) throws TGException {

        try {

            HandshakeRequest request = (HandshakeRequest) TGMessageFactory.getInstance().createMessage(VerbId.HandShakeRequest);
            request.setRequestType(HandshakeRequest.RequestType.Initiate);
            this.send(request);

            TGMessage msg = readWireMsg();
            if (!(msg instanceof HandshakeResponse)) {
                throw new TGException("Expecting a HandshakeResponse message, and received :" + msg.getClass());
            }

            HandshakeResponse response = (HandshakeResponse) msg;
            if (response.getResponseStatus() != HandshakeResponse.ResponseStatus.AcceptChallenge) {
                throw new TGException("Handshake failed with DB Server.");
            }
            int challenge   = response.getChallenge();

            challenge = challenge*2/3;

            request.updateSequenceAndTimeStamp(-1);
            request.setRequestType(HandshakeRequest.RequestType.ChallengeAccepted);
            request.setSslMode(sslMode);
            request.setChallenge(challenge);

            this.send(request);

            response = (HandshakeResponse) readWireMsg();

            if (response.getResponseStatus() != HandshakeResponse.ResponseStatus.ProceedWithAuthentication) {
                throw new TGException("Handshake failed with DB Server.");
            }

            gLogger.log(TGLogger.TGLevel.Info, "Hand shake successfull.");


        }
        catch(IOException e) {
            gLogger.logException("performHandshake failed", e);
            closeSocket();
            throw TGException.buildException("Failed to connect to the server at "+  linkURL.url, null, e);
        }
        catch (Exception e) {
            gLogger.logException("performHandshake failed", e);
            closeSocket();
            if (e instanceof TGException) throw e;
            throw TGException.buildException("Failed to create Handshake message", null, e);
        }
    }

    protected void doAuthenticate() throws TGException {
        try {
            AuthenticateRequest request = (AuthenticateRequest) TGMessageFactory.getInstance().createMessage(VerbId.AuthenticateRequest);
            request.setClientId(this.getClientId());
            request.setInboxAddr(this.getInboxAddr());
            request.setUserName(this.getUserName());
            request.setPassword(this.getPassword());

            this.send(request);

            AuthenticateResponse response = (AuthenticateResponse) readWireMsg();

            if (response.isSuccess()) {
                this.authToken = response.getAuthToken();
                this.sessionId = response.getSessionId();
                gLogger.log(TGLogger.TGLevel.Info, "Connected successfully using user:" + this.getUserName());
                return;
            }
            throw new TGAuthenticationException("Bad username/password combination", -1, "Bad username/password combination", "tgdb");
        }
        catch (IOException ioe) {
            gLogger.logException("doAuthenticate failed", ioe);
            closeSocket();
            throw TGException.buildException("Failed to connect to the server. Bad username and/or password combination "+  linkURL.url, null, ioe);
        }
    }

    /**
     *
     * @param msg
     * @throws TGException
     */
    public void send(TGMessage msg) throws TGException, IOException {
        if (state == LinkState.Closed)
            throw new TGException(linkHandler.getTerminatedText());

        if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
            logMessage(msg);
        }

        // loop while we can successfully reconnect.
        // if can not reconnect then we throw exception
        sendLock.lock();
        try {
        	if (output == null || state == LinkState.Closed)
        		throw new TGException(linkHandler.getTerminatedText());

            disablePing();
            byte[] buf = msg.toBytes();
            int bufLen = msg.getMessageByteBufLength();
            if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            	gLogger.log(TGLogger.TGLevel.Debug, "Send buf : %s", HexUtils.formatHex(buf, bufLen));
            }
            output.write(buf, 0, bufLen);
            output.flush();
            return;
        }
        finally {
            sendLock.unlock();
        }
    }

    private void logMessage(byte[] buf) {
        if (buf == null) {
            gLogger.log(TGLogger.TGLevel.Warning, "Unrecognized object, can not print it...");
        }

        if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            gLogger.log(TGLogger.TGLevel.Debug, "----------------- byte array --------------------");
            gLogger.log(TGLogger.TGLevel.Debug, HexUtils.formatHex(buf));
        }
    }

    private void logMessage(TGMessage m) {
        if (m == null) {
            if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
                gLogger.log(TGLogger.TGLevel.Debug, "Unrecognized object, can not print it...");
            }
        }
        if (gLogger.isEnabled(TGLogger.TGLevel.Debug)) {
            gLogger.log(TGLogger.TGLevel.Debug, "----------------- outgoing message --------------------");
            gLogger.log(TGLogger.TGLevel.Debug, m.toString());
        }
    }
}

