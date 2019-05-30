/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name : SessionForcefullyTermincated.${EXT}
 * Created on: 3/18/18
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: SessionForcefullyTerminated.java 3141 2019-04-26 00:00:35Z nimish $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

public class SessionForcefullyTerminated extends ExceptionMessage {


    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {

    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        msg = is.readUTF();
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.SessionForcefullyTerminated;
    }

    public String getKillString() { return msg; }


}
