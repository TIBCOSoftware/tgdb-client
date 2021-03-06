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
 *  File name : TransactionStatisticsImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: TransactionStatisticsImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGTransactionStatistics;

public class TransactionStatisticsImpl implements TGTransactionStatistics {
	
	

	@Override
	public String toString() {
		return "TransactionStatisticsImpl [transactionProcessorsCount=" + transactionProcessorsCount
				+ ", transactionProcessedCount=" + transactionProcessedCount + ", transactionSuccessfulCount="
				+ transactionSuccessfulCount + ", averageProcessingTime=" + averageProcessingTime
				+ ", pendingTransactionsCount=" + pendingTransactionsCount + ", transactionLoggerQueueDepth="
				+ transactionLoggerQueueDepth + "]";
	}

	protected long transactionProcessorsCount;
	protected long transactionProcessedCount;
	protected long transactionSuccessfulCount;
	protected double averageProcessingTime;
	protected long pendingTransactionsCount;
	protected int transactionLoggerQueueDepth;

	public TransactionStatisticsImpl (
		long _transactionProcessorsCount,
		long _transactionProcessedCount,
		long _transactionSuccessfulCount,
		double _averageProcessingTime,
		long _pendingTransactionsCount,
		int _transactionLoggerQueueDepth) 
	{
		transactionProcessorsCount = _transactionProcessorsCount;
		transactionProcessedCount = _transactionProcessedCount;
		transactionSuccessfulCount = _transactionSuccessfulCount;
		averageProcessingTime = _averageProcessingTime;
		pendingTransactionsCount = _pendingTransactionsCount;
		transactionLoggerQueueDepth = _transactionLoggerQueueDepth;
	} 

	@Override
	public long getTransactionProcessorsCount() {
		return this.transactionProcessorsCount;
	}

	@Override
	public long getTransactionProcessedCount() {
		return this.transactionProcessedCount;
	}

	@Override
	public long getTransactionSuccessfulCount() {
		return this.transactionSuccessfulCount;
	}

	@Override
	public double getAverageProcessingTime() {
		return this.averageProcessingTime;
	}

	@Override
	public long getPendingTransactionsCount() {
		return this.pendingTransactionsCount;
	}

	@Override
	public int getTransactionLoggerQueueDepth() {
		return this.transactionLoggerQueueDepth;
	}

}
