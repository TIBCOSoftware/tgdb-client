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
 *
 *  File name : DumpStacktraceRequest.java
 *  Created on: 03/29/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: DumpStacktraceRequest.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import java.io.IOException;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.pdu.impl.AbstractProtocolMessage;

public class DumpStacktraceRequest extends /*AuthenticatedMessage*/ AbstractProtocolMessage {
	
	@Override
	protected void writePayload(TGOutputStream os) throws TGException, IOException {
		
		/*
		int dataLength = 0;
	    int checksum = 0;
	    
	    int command = 39;
			    
	    //os.writeInt(dataLength);
	    //os.writeInt(checksum);
	    //os.writeInt(command);
	     */
	}
	
	@Override
	protected void readPayload(TGInputStream is) throws TGException, IOException {
	}

	@Override
	public boolean isUpdateable() {
		return false;
	}

	@Override
	public VerbId getVerbId() {
		return VerbId.DumpStacktraceRequest;
	}

}
