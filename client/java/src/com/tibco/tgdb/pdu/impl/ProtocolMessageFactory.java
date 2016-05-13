package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGBadVerb;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.util.Arrays;
import java.util.HashMap;

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
 * File name :ProtocolMessageFactory
 * Created on: 1/31/15
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: ProtocolMessageFactory.java 583 2016-03-15 02:02:39Z vchung $
 */
public class ProtocolMessageFactory extends TGMessageFactory {

    HashMap<VerbId, Class<? extends TGMessage>> msgtypeMap = new HashMap<VerbId, Class<? extends TGMessage>>();

    public ProtocolMessageFactory()  {
        for (VerbId vid : VerbId.values()) {
            msgtypeMap.put(vid, vid.getMessageClass());
        }
    }

    public TGMessage createMessage(VerbId verbId) throws TGException {
        Class<? extends TGMessage> klazz = null;
        try {
            klazz = msgtypeMap.get(verbId);

            if (klazz == null)
                throw new TGBadVerb(String.format("Invalid verbid:%s specified for Message construction", verbId), null);

            TGMessage msg = klazz.newInstance();
            return msg;
        } catch (InstantiationException | IllegalAccessException e) {
            throw TGException.buildException(String.format("Could not create message object from class:%s", klazz.getName()), null, e );
        }
    }

    public TGMessage createMessage(VerbId verbId, long authToken, long sessionId) throws TGException {
        Class<? extends TGMessage> klazz = null;
        try {
            klazz = msgtypeMap.get(verbId);

            if (klazz == null)
                throw new TGBadVerb(String.format("Invalid verbid:%s specified for Message construction", verbId), null);

            TGMessage msg = klazz.getConstructor(long.class, long.class).newInstance(authToken, sessionId);
            return msg;
        } catch (InstantiationException | IllegalAccessException | NoSuchMethodException | InvocationTargetException e) {
            throw TGException.buildException(String.format("Could not create message object from class:%s", klazz.getName()), null, e );
        }
    }

    public TGMessage createMessage(byte[] buffer, int offset, int length) throws TGException, IOException
    {
        byte[] buf;

        if (buffer.length == length) {
            buf = buffer;
        } else {
            buf = Arrays.copyOfRange(buffer, offset, offset + length);
        }
        TGMessage msg = createMessage(AbstractProtocolMessage.verbIdFromBytes(buf));
        msg.fromBytes(buffer);

        return msg;
    }
}
