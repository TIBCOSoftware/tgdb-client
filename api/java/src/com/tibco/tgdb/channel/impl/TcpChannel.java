/**
 * Copyright 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: DefaultSocketImpl.java
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TcpChannel.java 4626 2020-11-02 01:10:07Z ssubrama $
 */

package com.tibco.tgdb.channel.impl;

import java.io.BufferedInputStream;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.EOFException;
import java.io.IOException;
import java.net.Socket;
import java.util.Arrays;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import com.tibco.tgdb.TGVersion;
import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGChannelDisconnectedException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGVersionMismatchException;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.AuthenticateRequest;
import com.tibco.tgdb.pdu.impl.AuthenticateResponse;
import com.tibco.tgdb.pdu.impl.ExceptionMessage;
import com.tibco.tgdb.pdu.impl.HandshakeRequest;
import com.tibco.tgdb.pdu.impl.HandshakeResponse;
import com.tibco.tgdb.pdu.impl.HandshakeResponse.ResponseStatus;
import com.tibco.tgdb.pdu.impl.SessionForcefullyTerminated;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.HexUtils;
import com.tibco.tgdb.utils.TGConstants;
import com.tibco.tgdb.utils.TGEnvironment;
import com.tibco.tgdb.utils.TGProperties;

import static com.tibco.tgdb.exception.TGException.TGExceptionType.HandshakeResponseError;


public class TcpChannel extends AbstractChannel {

    final  Lock shutdownLock = new ReentrantLock();

    Socket                  socket         = null;
    DataInputStream         input          = null;
    DataOutputStream        output         = null;
    // set while reconnect attempt is in progress
    volatile boolean        reconnecting   = false;
    private Socket originalSocket;


    /**
     * Construct a Tcp channel on a URL, with a set properties
     *
     * @param linkUrl
     * @param props
     */

    protected TcpChannel(TGChannelUrl linkUrl, TGProperties<String, String> props) throws TGException
    {
        super(linkUrl, props);
    }

    /**
     * Create the physical link socket
     *
     * @throws TGException
     */
    protected void createSocket() throws TGException {
        String host = linkURL.getHost();
        int port = linkURL.getPort();
        String failureMessage = "Failed to connect to the server at " + linkURL.getUrlAsString();

        state.set(LinkState.NotConnected);

        try {
            shutdownLock.lock();
            socket = null;

            AbstractSocket ts = new DefaultSocketImpl(this.properties);
            socket = ts.createSocket(host, port, TGEnvironment.getInstance().getChannelConnectTimeout());
            this.setSocket(socket);


        } catch (IOException ioex) {
            if (socket != null) closeSocket();
            throw TGException.buildException(failureMessage, null, ioex);
        }
        finally {
            shutdownLock.unlock();
        }
    }


    /**
     * Close the socket.
     */
    void closeSocket() {
        //Flush the output.
        try {
            shutdownLock.lock();
            if (output != null) {
                output.flush();
                output.close();
            }
            Socket s = socket;
            if (s != null) {
                s.shutdownInput();
                s.shutdownOutput();
                s.close();
            }
        } catch (Exception e) { } //ignore the error
        finally {
            shutdownLock.unlock();
            socket = null;
            output = null;
            input  = null;
        }

    }

    /**
     * On connect called by the AbstractChannel
     * @throws TGException
     */

    public void onConnect() throws TGException {
        TGMessage msg = tryRead(); //For ChannelDisconnected message
        if (msg instanceof SessionForcefullyTerminated) {
            throw new TGChannelDisconnectedException(((SessionForcefullyTerminated)msg).getKillString());
        }
        performHandshake(false);
        doAuthenticate();
    }

    /**
     * Read wire message called from the ChannelReader
     * @return
     * @throws IOException
     */
    protected TGMessage readWireMsg() throws TGException, IOException {
        DataInputStream in      = input;
        byte[]          buffer  = TGConstants.EmptyByteArray;

        if (in == null)
            return null;

        disablePing();

        if (isClosed()) return null;
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
            //throw TGException.buildException((ExceptionMessage) msg);
        }
        
        if (msg instanceof HandshakeResponse)
        {
        	if (((HandshakeResponse)(msg)).getResponseStatus() == ResponseStatus.ChallengeFailed)
        	{
        		String errorMessage = ((HandshakeResponse)(msg)).getErrorMessage();
        		throw new TGVersionMismatchException(errorMessage);
        	}
        }
        return msg;
    }

    protected TGMessage tryRead()
    {
        try {
            if (input.available() > 0) {
                return readWireMsg();
            }
        }
        catch (Exception e) { }
        return null;
    }

    void intToBytes(int value, byte[] bytes, int offset) {
        for (int i=0; i<4; i++)
            bytes[offset+i] = (byte)((value>>>(8*(3-i)))&0xff);
    }

    protected void setSocket(Socket newsocket) throws IOException
    {
        this.originalSocket = this.socket;
        this.socket = newsocket;
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


    /**
     * performHandshake
     *
     * @param sslMode
     * @throws TGException
     */
    protected void performHandshake(boolean sslMode) throws TGException {

        HandshakeRequest request = (HandshakeRequest) TGMessageFactory.getInstance().createMessage(VerbId.HandShakeRequest);
        request.setRequestType(HandshakeRequest.RequestType.Initiate);

        TGMessage msg = this.requestReply(request);
        if (!(msg instanceof HandshakeResponse)) {
            if (msg instanceof SessionForcefullyTerminated) {
                throw new TGChannelDisconnectedException(((SessionForcefullyTerminated)msg).getKillString());
            }
            if (msg instanceof ExceptionMessage) {
                ExceptionMessage exmsg = (ExceptionMessage)msg;
                throw new TGException(String.format("Handshake Failed:%s. Cannot connect to the server at:%s",exmsg.getMessage(), linkURL.getUrlAsString()));
            }
            throw TGException.buildException(
                    String.format("Expecting a HandshakeResponse message, and received:%s. Cannot connect to the server at:%s",
                            msg.getClass(), linkURL.getUrlAsString()), HandshakeResponseError, null);
        }

        HandshakeResponse response = (HandshakeResponse) msg;
        if (response.getResponseStatus() != HandshakeResponse.ResponseStatus.AcceptChallenge) {
            throw new TGException(String.format("%s:Handshake Failed. Cannot connect to the server at:%s","TGDB-HANDSHAKE-ERR", linkURL.getUrlAsString()));
        }
        
        //
        // Validate the version specific information on the response object
        //
        validateHandshakeResponseVersion (response);
        
        
        long challenge = response.getChallenge();

        long jcVersion = TGVersion.getInstance().getAsLong();
        challenge = jcVersion;
        

        request.updateSequenceAndTimeStamp(-1);
        request.setRequestType(HandshakeRequest.RequestType.ChallengeAccepted);
        request.setSslMode(sslMode);
        request.setChallenge(challenge);
        
        msg = this.requestReply(request);
        if (!(msg instanceof HandshakeResponse)) {
            if (msg instanceof SessionForcefullyTerminated) {
                throw new TGChannelDisconnectedException(((SessionForcefullyTerminated)msg).getKillString());
            }
            if (msg instanceof ExceptionMessage) {
                ExceptionMessage exmsg = (ExceptionMessage)msg;
                throw new TGException(String.format("Handshake Failed:%s. Cannot connect to the server at:%s",exmsg.getMessage(), linkURL.getUrlAsString()));
            }
            throw TGException.buildException(
                    String.format("Expecting a HandshakeResponse message, and received:%s. Cannot connect to the server at:%s",
                            msg.getClass(), linkURL.getUrlAsString()), HandshakeResponseError, null);
        }

        response = (HandshakeResponse) msg;

        if (response.getResponseStatus() != HandshakeResponse.ResponseStatus.ProceedWithAuthentication) {
            throw new TGException(String.format("%s:Handshake Failed. Cannot connect to the server at:%s","TGDB-HANDSHAKE-ERR", linkURL.getUrlAsString()));
        }
        gLogger.log(TGLogger.TGLevel.Info, "Handshake successfull.");
        return;

    }

    private void validateHandshakeResponseVersion(HandshakeResponse response) throws TGVersionMismatchException {
    	long serverVersion = response.getChallenge();
		//ServerVersionInfo serverVersionInfo = new ServerVersionInfo(serverVersion);
    	TGVersion serverVersionInfo = TGVersion.getInstanceFromLong(serverVersion);
		
		//System.out.println ("Server Version Info: " + serverVersionInfo.toString());
		
		TGVersion javaClientVersion = TGVersion.getInstance();
		
		byte jcMajor = javaClientVersion.getMajor();
		byte jcMinor = javaClientVersion.getMinor();
		byte jcUpdate = javaClientVersion.getUpdate();
		
		byte sMajor = serverVersionInfo.getMajor();
		byte sMinor = serverVersionInfo.getMinor();
		byte sUpdate = serverVersionInfo.getUpdate();
		
		
		//
		// Currently, the check is for the exact version match between the server, and Java-Client;
		// in future, the validation will be enhanced to support the range of versions.
		//
		//if ((jcMajor == sMajor) && (jcMinor == sMinor) && (jcUpdate == sUpdate))
		
		if (javaClientVersion.equals(serverVersionInfo))
		{}
		else {
			String jcVersionString = "Major: " + jcMajor + " Minor: " + jcMinor + " Update: " + jcUpdate;
			String serverVersionString = "Major: " + sMajor + " Minor: " + sMinor + " Update: " + sUpdate;
			String strExceptionMessage = "Java-Client-Version and Server-Version are not the exact match. Java-Client-Version-Detail: " + jcVersionString + " and Server-Version-Detail: " + serverVersionString; 
			throw new TGVersionMismatchException(strExceptionMessage);
		}
		
		
	}

	protected void doAuthenticate() throws TGException
    {
        AuthenticateRequest request = (AuthenticateRequest) TGMessageFactory.getInstance().createMessage(VerbId.AuthenticateRequest);
        AuthenticateResponse response = null;
        request.setClientId(this.getClientId());
        request.setInboxAddr(this.getInboxAddr());
        request.setUserName(this.getUserName());
        request.setPassword(this.getPassword());
        request.setDatabaseName(this.getDatabaseName());

         //SSB: ACL related changes to add roles
        String specifiedRoles = properties.get(ConfigName.ConnectionSpecifiedRoles.getName());
        List<String> roleList = null;
        //if the property value is not null
        if(specifiedRoles!=null)
        {
            if(specifiedRoles.isEmpty())
               roleList = new ArrayList<String>();
            else
               roleList  = Arrays.asList(specifiedRoles.split(","));//split on the basis of comma to get a list of role names
        }
        for(int i=0; roleList!=null && i<roleList.size(); i++)
                roleList.set(i,roleList.get(i).trim());

        if(roleList!=null) request.setRoles(roleList);
         //SSB: ACL related changes to add roles

        response = (AuthenticateResponse) this.requestReply(request);
        if (response.isSuccess()) {
            this.authToken = response.getAuthToken();
            this.sessionId = response.getSessionId();
            this.dataCryptographer = new DataCryptoGrapher(sessionId, response.getServerCertificateBuffer());

            gLogger.log(TGLogger.TGLevel.Info, "Connected successfully using user:" + this.getUserName());
            return;
        }
        throw new TGAuthenticationException("Bad username/password combination", -1, "tgdb");
    }

    /**
     * Basic send implementation. No locking and no exception management. It is taken care at the Abstract Channel level
     * which calls this method.
     *
     * @param msg
     * @throws TGException
     */
    protected void send(TGMessage msg) throws TGException, IOException {
        if ((output == null) || (isClosed()))  throw new TGException("Channel is Closed");

        if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
            logMessage(msg);
        }
        disablePing();
        byte[] buf = msg.toBytes();
        int bufLen = msg.getMessageByteBufLength();
        if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
            gLogger.log(TGLogger.TGLevel.DebugWire, "Send buf : %s", HexUtils.formatHex(buf, bufLen));
        }
        output.write(buf, 0, bufLen);
        output.flush();
        return;

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
            if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
                gLogger.log(TGLogger.TGLevel.DebugWire, "Unrecognized object, can not print it...");
            }
            return;
        }
        if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
            gLogger.log(TGLogger.TGLevel.DebugWire, "----------------- outgoing message --------------------");
            gLogger.log(TGLogger.TGLevel.DebugWire, m.toString());
        }
    }
}

