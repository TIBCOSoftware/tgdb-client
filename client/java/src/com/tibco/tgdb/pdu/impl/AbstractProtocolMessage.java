package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.TGProtocolVersion;
import com.tibco.tgdb.exception.TGBadMagic;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGInvalidMessageLength;
import com.tibco.tgdb.exception.TGProtocolNotSupported;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.util.concurrent.atomic.AtomicLong;

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
 * File name :AbstractProtocolMessage
 * Created on: 1/31/15
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: AbstractProtocolMessage.java 583 2016-03-15 02:02:39Z vchung $
 */
public abstract class AbstractProtocolMessage implements TGMessage {

    static AtomicLong gAtomicSequenceNumber = new AtomicLong();


    long sequenceNo;
    long timestamp = -1;
    long requestId = -1;
    long authToken = 0;
    long sessionId = 0;

    int bufLength = -1;
    short dataOffset = -1;


    protected AbstractProtocolMessage() {
        sequenceNo = gAtomicSequenceNumber.getAndIncrement();
    }

    protected AbstractProtocolMessage(long authToken, long sessionId) {
    	this.authToken = authToken;
    	this.sessionId = sessionId;
    }

    public static VerbId verbIdFromBytes(byte[] buffer) throws TGException, IOException {


        TGInputStream is = new ProtocolDataInputStream(buffer);
        int len = is.readInt();
        if (len != buffer.length) {
            throw new TGInvalidMessageLength("buffer length mismatch", null);
        }

        int magic = is.readInt();

        if (magic != TGProtocolVersion.getMagic()) {
            throw new TGBadMagic("Bad Message Magic. ", null);
        }

        int protocolVersion = is.readShort();

        if (!TGProtocolVersion.isCompatible(protocolVersion)) {
            throw new TGProtocolNotSupported("Unsupported Protocol version", null);
        }

        return VerbId.fromId(is.readShort());
    }



    public long getSequenceNo() {
        return sequenceNo;
    }

    public long getTimestamp() {
        if (timestamp == -1) timestamp = System.currentTimeMillis();
        return timestamp;
    }

    public void setTimestamp(long timestamp) throws TGException {
        if ((timestamp == -1) || isUpdateable())
            this.timestamp = timestamp;
        else {
            throw new TGException("Mutating a readonly message");
        }
    }

    public long getRequestId() {
        return requestId;
    }

    public void setRequestId(long requestId) {
        this.requestId = requestId;
    }

    public long getAuthToken() {
        return authToken;
    }

    public void setAuthToken(long authToken) {
        this.authToken = authToken;
    }

    public long getSessionId() {
        return sessionId;
    }

    public void setSessionId(long sessionId) {
        this.sessionId = sessionId;
    }

    public void updateSequenceAndTimeStamp(long timestamp) throws TGException
    {
        if (isUpdateable())  {
            this.sequenceNo = gAtomicSequenceNumber.getAndIncrement();
            this.timestamp = timestamp;
            bufLength = -1;
            return;
        }
        throw new TGException("Mutating a readonly message");

    }

    public final byte[] toBytes() throws TGException, IOException
    {
        TGOutputStream os = new ProtocolDataOutputStream();

        writeHeader(os);
        writePayload(os);

        bufLength = os.getLength();

        os.writeIntAt(0, bufLength);



        return os.getBuffer();
    }

    public final int getMessageByteBufLength() {
        return bufLength;
    }

    public final void fromBytes(byte[] buffer) throws TGException, IOException {
        TGInputStream is = new ProtocolDataInputStream(buffer);
        int len = is.readInt();
        if (len != buffer.length) {
            throw new TGInvalidMessageLength("buffer length mismatch", null);
        }

        readHeader(is);

        readPayload(is);
    }

    /**
     * The Server describes the pdu header as below.
     * struct _tg_pduheader_t_ {
         tg_int32    length;         //length of the message including the header
         tg_int32    magic;          //Magic to recognize this is our message
         tg_int16    protVersion;    //protocol version
         tg_pduverb  verbId;         //we write the verb as a short value
         tg_uint64   sequenceNo;     //message Sequence No from the client
         tg_uint64   timestamp;      //Timestamp of the message sent.
         tg_uint64   requestId;      //Unique _request Identifier from the client, which is returned
         tg_int32    dataOffset;     //Offset from where the payload begins
     }
     * @param os
     * @throws TGException
     * @throws IOException
     */
    protected void writeHeader(TGOutputStream os) throws TGException, IOException
    {
        os.writeInt(0); //The length is written later.
        os.writeInt(TGProtocolVersion.getMagic());
        os.writeShort(TGProtocolVersion.getProtocolVersion());
        os.writeShort(getVerbId().getId());

        os.writeLong(getSequenceNo());
        os.writeLong(getTimestamp());
        os.writeLong(getRequestId());
        
        os.writeLong(getAuthToken());
        os.writeLong(getSessionId());
        os.writeShort(os.getPosition() + 2); //DataOffset.
    }

    protected void readHeader(TGInputStream is) throws TGException, IOException
    {
        int magic = is.readInt();

        if (magic != TGProtocolVersion.getMagic()) {
            throw new TGBadMagic("Bad Message Magic. ", null);
        }

        short protocolVersion = is.readShort();

        if (!TGProtocolVersion.isCompatible(protocolVersion)) {
            throw new TGProtocolNotSupported("Unsupported Protocol version", null);
        }


        VerbId vid = VerbId.fromId(is.readShort());
        if (vid != this.getVerbId()) {
            throw new TGException("Incorrect Message Type");
        }


        sequenceNo   = is.readLong(); //we should think hard
        timestamp    = is.readLong();
        requestId    = is.readLong();
        authToken 	 = is.readLong();
        sessionId    = is.readLong();
        dataOffset   = is.readShort();
    }

    protected abstract void writePayload(TGOutputStream os) throws TGException, IOException;

    protected abstract void readPayload(TGInputStream is) throws TGException, IOException;

    public abstract boolean isUpdateable();

    public abstract VerbId getVerbId();



}
