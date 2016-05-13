package com.tibco.tgdb.pdu;

import com.tibco.tgdb.exception.TGException;

import java.io.IOException;

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
 * File name :TGMessage
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGMessage.java 583 2016-03-15 02:02:39Z vchung $
 */
public interface  TGMessage {

    /**
     * Get VerbId of the Message
     * @return
     */
    public VerbId getVerbId();

    /**
     * Get the SequenceNo of the Message
     * @return
     */
    public long getSequenceNo();

    /**
     * Get the RequestId for the Message. This will be used as the CoorelationId
     * @return
     */
    public long getRequestId();

    /** FIXME:  Is it ok to expose this?
     * Set the request Id
     * @param requestId
     */
    public void setRequestId(long requestId);

    /**
     * Get the TimeStamp of the Message
     * @return
     */
    public long getTimestamp();

    /**
     * Set the Timestamp
     * @param timestamp
     */
    public void setTimestamp(long timestamp) throws TGException;

    /**
     * Get the AuthToken
     * @return
     */
    public long getAuthToken();

    /**
     * Set the AuthToken
     * @param authToken
     */
    public void setAuthToken(long authToken);

    /**
     * Get the Session Id
     * @return
     */
    public long getSessionId();

    /**
     * Set the session Id
     * @param sessionId
     */
    public void setSessionId(long sessionId);

    /**
     * Is this message updateable
     * @return
     */
    public boolean isUpdateable();

    /**
     * If Message is Mutable, then update the SequenceAndTimeStamp
     * @param timestamp
     * @throws TGException
     */
    public void updateSequenceAndTimeStamp(long timestamp) throws TGException;


    /**
     * Get the bytes from the message
     * @return
     * @throws Exception
     */
    public byte[] toBytes() throws TGException, IOException;

    /**
     * Get the MessageByteBufLength, Call this method after the toBytes is called.
     * @return
     */
    public int getMessageByteBufLength();

    /**
     * Reconstruct the message from the buffer.
     * @param buffer
     * @throws TGException
     * @throws IOException
     */
    public void fromBytes(byte[] buffer) throws TGException, IOException;




}
