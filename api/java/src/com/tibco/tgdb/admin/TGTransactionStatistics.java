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
 *  File name :TGTransactionStatistics.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  <p>This interface allows users to retrieve the transaction statistics from server
 *  
 *  SVN Id: $Id: TGTransactionStatistics.java 3120 2019-04-25 21:21:48Z nimish $  
 */

package com.tibco.tgdb.admin;

public interface TGTransactionStatistics {
	
	
	/**
	 * Get the transaction processors count
	 * @return transaction processors count
	 */
	long getTransactionProcessorsCount ();
	
	
	/**
	 * Get the processed transaction count 
	 * @return processed tnansaction count
	 */
	long getTransactionProcessedCount ();
	
	
	/**
	 * Get the successful transactions count
	 * @return successful transactions count
	 */
	long getTransactionSuccessfulCount ();
	
	
	/**
	 * Get the average processing time
	 * @return average processing time for the transactions
	 */
	double getAverageProcessingTime ();
	
	
	/**
	 * Get the pending transactions count
	 * @return the pending transactions count
	 */
	long getPendingTransactionsCount ();
	
	
	/**
	 * Get the queue depth of transactionLogger
	 * @return queue depth of transactionLogger
	 */
	int getTransactionLoggerQueueDepth ();

}
