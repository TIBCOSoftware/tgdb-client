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
 * File name: ChannelMessageTraceReader.java
 * Created on: 2019-02-07
 * Created by: nimish
 * <p/>
 * SVN Id: $Id: ChannelMessageTraceReader.java 3158 2019-04-26 20:49:24Z kattaylo $
 */

package com.tibco.tgdb.channel.impl;

import java.io.ByteArrayInputStream;
import java.io.DataInputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.TGMessageFactory;

public class ChannelMessageTraceReader extends Thread {

	protected String inputFileName;
	protected String commitTraceDir;
	protected ConnectionImpl conn;
	
	static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
	
	public ChannelMessageTraceReader (String _clientId, String _commitTraceDir, ConnectionImpl _conn)
	{
		inputFileName = _clientId + ".trace";
		commitTraceDir = _commitTraceDir;
		conn = _conn;
		
		String filepathSeperator = System.getProperty("file.separator");
		
		if (commitTraceDir.endsWith(filepathSeperator))
		{
			inputFileName = commitTraceDir + _clientId + ".trace";
		}
		else 
		{
			inputFileName = commitTraceDir + filepathSeperator + _clientId + ".trace";
		}

	}
	
	
	@Override
	public void run() {
		
		File dir = new File(commitTraceDir);
		int nCountOfFiles = dir.listFiles().length;
		int mCount = 0;
		
		for (int i = 0; i < nCountOfFiles; ++i)
		{
			
			DataInputStream dis = null;
			String currentFileName = inputFileName + "." + i;
			try {
				dis = new DataInputStream (new FileInputStream(currentFileName));
			} catch (FileNotFoundException e2) {
				continue;
			}
				
			for (;;)
			{
				try {
					int nBytesToRead = dis.readInt();
					byte[] bytesRead = new byte[nBytesToRead];
					int nBytesRead = dis.read(bytesRead, 0, nBytesToRead);
					
					TGMessage requestMessage = TGMessageFactory.getInstance().createMessage(bytesRead, 0, nBytesRead);

					
					//System.out.println("Request Message: Verb ID = " + requestMessage.getVerbId());
	                gLogger.log(TGLogger.TGLevel.Info, "Request Message: Verb ID = " + requestMessage.getVerbId());

	                
					if (conn != null)
					{
						TGMessage responseMessage = conn.sendTraceMessage(requestMessage);
						//System.out.println("Response Message: Verb ID = " + responseMessage.getVerbId());
						gLogger.log(TGLogger.TGLevel.Info, "Response Message: Verb ID = " + responseMessage.getVerbId());
					}
					
					
					
					
					//System.out.println("File: " + i + " Message Value = " + tgMessage.toString());
					
					++mCount;
				}
				catch (Exception e1)
				{
					break;
				}
			}
			
			if (dis != null)
			{
				try {
					dis.close();
				} catch (IOException e) {
					e.printStackTrace();
				}
			}
		}
		
		if (conn != null) conn.disconnect();
	}
	
	private int readInternalMessageLength(byte[] bytesRead) throws IOException {
		
		DataInputStream dis = new DataInputStream (new ByteArrayInputStream(bytesRead));
		int internalMessageLength = dis.readInt();
		dis.close();
		return internalMessageLength;
	
	}


	public static void main(String[] args) throws Exception {
		//ChannelMessageTraceReader qD = new ChannelMessageTraceReader("clientidtest", "g:\\temp\\abcd");
		
		ConnectionImpl conn = getConnectionHandle();
		//ChannelMessageTraceReader qD = new ChannelMessageTraceReader("tgdb.java-api.client", "G:\\temp\\tgdb-565\\tgdb-test-20190206\\tracefile", conn);
		
		ChannelMessageTraceReader qD = new ChannelMessageTraceReader("tgdb.java-api.client", "G:\\temp\\traceMessages0", conn);
		
		
		qD.start();
	}
	
	private static ConnectionImpl getConnectionHandle () {
		
		String url = "tcp://scott@localhost:8222";
		String user = "scott";
	    String passwd = "scott";

		TGConnection conn = null;
		try {
			conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
			conn.connect();
			
			return (ConnectionImpl)conn;
		}
		catch (TGException e) {
			e.printStackTrace();
		} finally {
			if (conn != null)
			{
				//conn.disconnect();
			}
		}
		return null;
	}
	
}