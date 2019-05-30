/**
 * Copyright (c) 2019 TIBCO Software Inc.
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
 * File name: ChannelTracerImpl.java
 * Created on: 2019-02-07
 * Created by: nimish
 * <p/>
 * SVN Id: $Id: ChannelTracerImpl.java 3158 2019-04-26 20:49:24Z kattaylo $
 */

package com.tibco.tgdb.channel.impl;

import java.io.IOException;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;

import com.tibco.tgdb.channel.TGTracer;
import com.tibco.tgdb.pdu.TGMessage;

public class ChannelTracerImpl implements TGTracer {

	protected Queue<TGMessage> q = new ConcurrentLinkedQueue<TGMessage>();
	protected ChannelMessageTracer qs = null;
	protected String clientId;
	
	public ChannelTracerImpl (String _clientId, String commitTraceDir) throws IOException {
		this.clientId = _clientId;
		qs = new ChannelMessageTracer(q, clientId, commitTraceDir);
		qs.start();
	}

	
	@Override
	public void trace(TGMessage message) {
		q.add(message);
	}
	
	
	
}


