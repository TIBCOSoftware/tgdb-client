package connection

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: TGTransaction.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TransactionImpl struct {
	transactionId int64
}

// Make sure that the transactionImpl implements the TGTransaction interface
var _ types.TGTransaction = (*TransactionImpl)(nil)

func DefaultTransaction() *TransactionImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TransactionImpl{})

	return &TransactionImpl{transactionId: -1}
}

func NewTransaction(txnId int64) *TransactionImpl {
	newTxn := DefaultTransaction()
	newTxn.transactionId = txnId
	return newTxn
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGTransaction
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) GetTransactionId() int64 {
	return obj.transactionId
}

func (obj *TransactionImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TransactionImpl:{")
	buffer.WriteString(fmt.Sprintf("TransactionId: '%d'", obj.transactionId))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGSerializable
/////////////////////////////////////////////////////////////////

// ReadExternal reads the byte format from an external input stream and constructs a system object
func (obj *TransactionImpl) ReadExternal(iStream types.TGInputStream) types.TGError {
	// No-op for Now
	return nil
}

// WriteExternal writes a system object into an appropriate byte format onto an external output stream
func (obj *TransactionImpl) WriteExternal(oStream types.TGOutputStream) types.TGError {
	// No-op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.transactionId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
