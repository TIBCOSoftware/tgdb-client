package admin

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGTransactionStatistics.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// TGTransactionStatistics allows users to retrieve the transaction statistics from server
type TGTransactionStatistics interface {
	// GetAverageProcessingTime returns the average processing time for the transactions
	GetAverageProcessingTime() float64
	// GetPendingTransactionsCount returns the pending transactions count
	GetPendingTransactionsCount() int64
	// GetTransactionLoggerQueueDepth returns the queue depth of transactionLogger
	GetTransactionLoggerQueueDepth() int
	// GetTransactionProcessorsCount returns the transaction processors count
	GetTransactionProcessorsCount() int64
	// GetTransactionProcessedCount returns the processed transaction count
	GetTransactionProcessedCount() int64
	// GetTransactionSuccessfulCount returns the successful transactions count
	GetTransactionSuccessfulCount() int64
}
