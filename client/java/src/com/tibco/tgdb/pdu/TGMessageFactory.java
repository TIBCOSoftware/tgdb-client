package com.tibco.tgdb.pdu;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.impl.ProtocolMessageFactory;

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
 * File name :TGMessageFactory
 * Created on: 12/22/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGMessageFactory.java 583 2016-03-15 02:02:39Z vchung $
 */
public abstract class TGMessageFactory {

    private static TGMessageFactory gFactory = new ProtocolMessageFactory();

    public static TGMessageFactory getInstance() {
        return gFactory;
    }

    public abstract TGMessage createMessage(VerbId verbId) throws TGException;

    public abstract TGMessage createMessage(VerbId verbId, long authToken, long sessionId) throws TGException;

    public abstract TGMessage createMessage(byte[] buffer, int offset, int length) throws TGException, IOException;

}
