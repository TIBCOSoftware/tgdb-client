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
 *  File name : TGServerLogDetails.java
 *  Created on: 04/22/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: TGServerLogDetails.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin;


public class TGServerLogDetails {

	public enum LogComponent {
		
		TGLC_COMMON_COREMEMORY (1L << 0),
		TGLC_COMMON_CORECOLLECTIONS (1L << 1),
		TGLC_COMMON_COREPLATFORM (1L << 2),
		TGLC_COMMON_CORESTRING (1L << 3),
		TGLC_COMMON_UTILS (1L << 4),
		TGLC_COMMON_GRAPH (1L << 5),
		TGLC_COMMON_MODEL (1L << 6),
		TGLC_COMMON_NET (1L << 7),
		TGLC_COMMON_PDU (1L << 8),
		TGLC_COMMON_SEC (1L << 9),
		TGLC_COMMON_FILES (1L << 10),
		TGLC_COMMON_RESV2 (1L << 11),
		
		//Server Components
		TGLC_SERVER_CDMP (1L << 12),
		TGLC_SERVER_DB (1L << 13),
		TGLC_SERVER_EXPIMP (1L << 14),
		TGLC_SERVER_INDEX (1L << 15),
		TGLC_SERVER_INDEXBTREE (1L << 16),
		TGLC_SERVER_INDEXISAM (1L << 17),
		TGLC_SERVER_QUERY (1L << 18),
		TGLC_SERVER_QUERY_RESV1 (1L << 19),
		TGLC_SERVER_QUERY_RESV2 (1L << 20),
		TGLC_SERVER_TXN (1L << 21),
		TGLC_SERVER_TXNLOG (1L << 22),
		TGLC_SERVER_TXNWRITER (1L << 23),
		TGLC_SERVER_STORAGE (1L << 24),
		TGLC_SERVER_STORAGEPAGEMANAGER (1L << 25),
		TGLC_SERVER_GRAPH (1L << 26),
		TGLC_SERVER_MAIN (1L << 27),
		TGLC_SERVER_RESV2 (1L << 28),
		TGLC_SERVER_RESV3 (1L << 29),
		TGLC_SERVER_RESV4 (1L << 30),
		
		//Security Components
		TGLC_SECURITY_DATA (1L << 31),
		TGLC_SECURITY_NET (1L << 32),
		TGLC_SECURITY_RESV1 (1L << 33),
		TGLC_SECURITY_RESV2 (1L << 34),
		
		TGLC_ADMIN_LANG (1L << 35),
		TGLC_ADMIN_CMD (1L << 36),
		TGLC_ADMIN_MAIN (1L << 37),
		TGLC_ADMIN_AST (1L << 38),
		TGLC_ADMIN_GREMLIN (1L << 39),
		
		TGLC_CUDA_GRAPHMGR (1L << 40),
		TGLC_CUDA_KERNELEXECUTIVE (1L << 41),
		TGLC_CUDA_RESV1 (1L << 42),
		
		TGLC_LOG_GLOBAL (0xFFFFFFFFFFFFFFFFL),

		// User Defined Components
		TGLC_LOG_COREALL (TGLC_COMMON_COREMEMORY.getLogComponent() | TGLC_COMMON_CORECOLLECTIONS.getLogComponent() | TGLC_COMMON_COREPLATFORM.getLogComponent() | TGLC_COMMON_CORESTRING.getLogComponent()),
		
		TGLC_LOG_GRAPHALL (TGLC_COMMON_GRAPH.getLogComponent() | TGLC_SERVER_GRAPH.getLogComponent()),
		TGLC_LOG_MODEL (TGLC_COMMON_MODEL.getLogComponent()),
		TGLC_LOG_NET (TGLC_COMMON_NET.getLogComponent()),
		TGLC_LOG_PDUALL (TGLC_COMMON_PDU.getLogComponent() | TGLC_SERVER_CDMP.getLogComponent()),
		TGLC_LOG_SECALL (TGLC_COMMON_SEC.getLogComponent() | TGLC_SECURITY_DATA.getLogComponent() | TGLC_SECURITY_NET.getLogComponent()),
		TGLC_LOG_CUDAALL (TGLC_LOG_GRAPHALL.getLogComponent() | TGLC_CUDA_GRAPHMGR.getLogComponent() | TGLC_CUDA_KERNELEXECUTIVE.getLogComponent()),
		TGLC_LOG_TXNALL (TGLC_SERVER_TXN.getLogComponent() | TGLC_SERVER_TXNLOG.getLogComponent() | TGLC_SERVER_TXNWRITER.getLogComponent()),
		TGLC_LOG_STORAGEALL (TGLC_SERVER_STORAGE.getLogComponent() | TGLC_SERVER_STORAGEPAGEMANAGER.getLogComponent()),
		TGLC_LOG_PAGEMANAGER (TGLC_SERVER_STORAGEPAGEMANAGER.getLogComponent()),
		TGLC_LOG_ADMINALL (TGLC_ADMIN_LANG.getLogComponent() | TGLC_ADMIN_CMD.getLogComponent() | TGLC_ADMIN_MAIN.getLogComponent() | TGLC_ADMIN_AST.getLogComponent() | TGLC_ADMIN_GREMLIN.getLogComponent()),
		TGLC_LOG_MAIN (TGLC_SERVER_MAIN.getLogComponent() | TGLC_ADMIN_MAIN.getLogComponent());
		
		
		protected long lc;
		LogComponent(long _lc)
		{
			lc = _lc;
		}
		
		public long getLogComponent ()
		{
			return lc;
		}

	}
	
	public enum LogLevel {

		TGLL_Console(-2),
		TGLL_Invalid(-1),
		TGLL_Fatal(0),
		TGLL_Error(1),
		TGLL_Warn(2),
		TGLL_Info(3),
		TGLL_User(4),
		TGLL_Debug(5),
		TGLL_DebugFine(6),
		TGLL_DebugFiner(7),
		TGLL_MaxLogLevel(8);
		
		protected int ll;

		LogLevel(int _ll)
		{
			ll = _ll;
		}
		
		public int getLogLevel () 
		{
			return ll;
		}
	}

	
	
	protected LogLevel logLevel;
	protected LogComponent logComponent;
	
	public TGServerLogDetails(LogComponent _logComponent, LogLevel _logLevel)
	{
		logComponent = _logComponent;
		logLevel = _logLevel;
	}
	
	public LogLevel getLogLevel() {
		return logLevel;
	}

	public LogComponent getLogComponent() {
		return logComponent;
	}	
	
}
