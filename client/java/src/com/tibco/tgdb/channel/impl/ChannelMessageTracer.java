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
 * File name: ChannelMessageTracer.java
 * Created on: 2019-02-07
 * Created by: nimish
 * <p/>
 * SVN Id: $Id: ChannelMessageTracer.java 3158 2019-04-26 20:49:24Z kattaylo $
 */

package com.tibco.tgdb.channel.impl;

import java.io.DataOutputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.Queue;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGMessage;

//
// Possibly use the executor
//
public class ChannelMessageTracer extends Thread {

	

	public static final int MB = 1 << 20;
	
	protected Queue<TGMessage> queue;
	protected String outputFileName;
	int currentSuffix = 0;
	DataOutputStream dos;
	FileOutputStream fos;
	int foslen = 0;

	public ChannelMessageTracer (Queue<TGMessage> _queue, String _clientId, String commitTraceDir) throws IOException
	{
		queue = _queue;
		String filepathSeperator = System.getProperty("file.separator");
		
		File dir = new File(commitTraceDir);
		if (!dir.exists())
		{
			boolean mkdir = dir.mkdir();
			if (mkdir)
			{
				//
				// Replace the System.out with corresponding log message
				//
				System.out.println("Directory: " + commitTraceDir + " created successfully.");
			}
		}
		
		if (commitTraceDir.endsWith(filepathSeperator))
		{
			outputFileName = commitTraceDir + _clientId + ".trace";
		}
		else 
		{
			outputFileName = commitTraceDir + filepathSeperator + _clientId + ".trace";
		}
		String fileName = String.format("%s.%d", outputFileName, currentSuffix);
		fos = new FileOutputStream(fileName);
		dos = new DataOutputStream(fos);

	}
	
	
	@Override
	public void run() {
		

		for (;;)
		{
			TGMessage message = queue.poll();
			if (message != null)
			{
				try {
					byte[] bytes = message.toBytes();
					int buflen = message.getMessageByteBufLength();
					
					tryRollover();
//					String countAdjustedFileName = outputFileName + "." + currentSuffix;
//					DataOutputStream dos = new DataOutputStream (new FileOutputStream(countAdjustedFileName, true));
					dos.writeInt(buflen);
					dos.write(bytes,0, buflen);
					foslen += (buflen + 4);
//					dos.close();
				} catch (TGException | IOException e) {
					e.printStackTrace();
				}
			}
			else {
				try {
					sleep(1000);
				} catch (InterruptedException e) {
					e.printStackTrace();
				}
			}
		}
	}

	private void tryRollover() throws IOException
	{
		if (foslen > MB) {
			dos.flush();
			fos.close();
			String fileName = String.format("%s.%d", outputFileName, ++currentSuffix);
			fos = new FileOutputStream(fileName);
			dos = new DataOutputStream(fos);
			foslen = 0;
		}
		return;
	}
	private int getNewSuffix(int currentSuffix) {
		for (;;)
		{
			String possibleFileName = outputFileName + "." + currentSuffix;
			File file = new File(possibleFileName);
			if (!file.exists())
			{
				break;
			}
			
			if (file.length() < MB)
			{
				break;
			}
			currentSuffix++;
		}
		return currentSuffix;
	}


	public static void main(String[] args) throws Exception {
		/*
		ChannelTracerImpl tracer = new ChannelTracerImpl("clientidtest", "g:\\temp\\abcdef");
		
		for (int i = 0; i < 1000*1000; ++i)
		{
			TempTGMessage message = new TempTGMessage(i);
			
			tracer.trace(message);
			sleep(100);
		}
		*/
	}
	
}