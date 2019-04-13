package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

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
 * File name: TransactionStatisticsImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TransactionStatisticsImpl struct {
	averageProcessingTime       float64
	pendingTransactionsCount    int64
	transactionLoggerQueueDepth int
	transactionProcessorCount   int64
	transactionProcessedCount   int64
	transactionSuccessfulCount  int64
}

// Make sure that the TransactionStatisticsImpl implements the TGTransactionStatistics interface
var _ TGTransactionStatistics = (*TransactionStatisticsImpl)(nil)

func DefaultTransactionStatisticsImpl() *TransactionStatisticsImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TransactionStatisticsImpl{})

	return &TransactionStatisticsImpl{}
}

func NewTransactionStatisticsImpl(_averageProcessingTime float64, _pendingTransactionsCount int64, _transactionLoggerQueueDepth int,
	_transactionProcessorCount, _transactionProcessedCount, _transactionSuccessfulCount int64) *TransactionStatisticsImpl {
	newConnectionInfo := DefaultTransactionStatisticsImpl()
	newConnectionInfo.averageProcessingTime = _averageProcessingTime
	newConnectionInfo.pendingTransactionsCount = _pendingTransactionsCount
	newConnectionInfo.transactionLoggerQueueDepth = _transactionLoggerQueueDepth
	newConnectionInfo.transactionProcessorCount = _transactionProcessorCount
	newConnectionInfo.transactionProcessedCount = _transactionProcessedCount
	newConnectionInfo.transactionSuccessfulCount = _transactionSuccessfulCount
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGTransactionStatisticsImpl
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TransactionStatisticsImpl:{")
	buffer.WriteString(fmt.Sprintf("AverageProcessingTime: '%+v'", obj.averageProcessingTime))
	buffer.WriteString(fmt.Sprintf(", PendingTransactionsCount: '%d'", obj.pendingTransactionsCount))
	buffer.WriteString(fmt.Sprintf(", TransactionLoggerQueueDepth: '%d'", obj.transactionLoggerQueueDepth))
	buffer.WriteString(fmt.Sprintf(", TransactionProcessorCount: '%d'", obj.transactionProcessorCount))
	buffer.WriteString(fmt.Sprintf(", TransactionProcessedCount: '%d'", obj.transactionProcessedCount))
	buffer.WriteString(fmt.Sprintf(", TransactionSuccessfulCount: '%d'", obj.transactionSuccessfulCount))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGTransactionStatistics
/////////////////////////////////////////////////////////////////

// GetAverageProcessingTime returns the average processing time for the transactions
func (obj *TransactionStatisticsImpl) GetAverageProcessingTime() float64 {
	return obj.averageProcessingTime
}

// GetPendingTransactionsCount returns the pending transactions count
func (obj *TransactionStatisticsImpl) GetPendingTransactionsCount() int64 {
	return obj.pendingTransactionsCount
}

// GetTransactionLoggerQueueDepth returns the queue depth of transactionLogger
func (obj *TransactionStatisticsImpl) GetTransactionLoggerQueueDepth() int {
	return obj.transactionLoggerQueueDepth
}

// GetTransactionProcessorsCount returns the transaction processors count
func (obj *TransactionStatisticsImpl) GetTransactionProcessorsCount() int64 {
	return obj.transactionProcessorCount
}

// GetTransactionProcessedCount returns the processed transaction count
func (obj *TransactionStatisticsImpl) GetTransactionProcessedCount() int64 {
	return obj.transactionProcessedCount
}

// GetTransactionSuccessfulCount returns the successful transactions count
func (obj *TransactionStatisticsImpl) GetTransactionSuccessfulCount() int64 {
	return obj.transactionSuccessfulCount
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.averageProcessingTime, obj.pendingTransactionsCount, obj.transactionLoggerQueueDepth,
		obj.transactionProcessorCount, obj.transactionProcessedCount, obj.transactionSuccessfulCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionStatisticsImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.averageProcessingTime, &obj.pendingTransactionsCount, &obj.transactionLoggerQueueDepth,
		&obj.transactionProcessorCount, &obj.transactionProcessedCount, &obj.transactionSuccessfulCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionStatisticsImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
