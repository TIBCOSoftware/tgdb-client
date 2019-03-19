package query

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"time"
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
 * File name: TGQuery.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

var logger = logging.DefaultTGLogManager().GetLogger()

type TGQueryImpl struct {
	qryConnection types.TGConnection
	qryHashId     int64
	qryOption     *TGQueryOptionImpl
	// These parameters are for tokenized query parameters specified as '?, ?, ...'
	qryParameters map[string]interface{}
}

func DefaultQuery() *TGQueryImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGQueryImpl{})

	newQuery := TGQueryImpl{
		qryHashId:     -1,
		qryParameters: make(map[string]interface{}, 0),
	}
	newQuery.qryOption = DefaultQueryOption()
	return &newQuery
}

// Make sure that the TGQueryImpl implements the TGQuery interface
var _ types.TGQuery = (*TGQueryImpl)(nil)

func NewQuery(conn types.TGConnection, queryHashId int64) *TGQueryImpl {
	newQuery := DefaultQuery()
	newQuery.qryConnection = conn
	newQuery.qryHashId = queryHashId
	return newQuery
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) GetQueryConnection() types.TGConnection {
	return obj.qryConnection
}

func (obj *TGQueryImpl) GetQueryId() int64 {
	return obj.qryHashId
}

func (obj *TGQueryImpl) GetQueryOption() *TGQueryOptionImpl {
	return obj.qryOption
}

func (obj *TGQueryImpl) GetQueryParameters() map[string]interface{} {
	return obj.qryParameters
}

func (obj *TGQueryImpl) SetQueryId(qId int64) {
	obj.qryHashId = qId
}

func (obj *TGQueryImpl) SetQueryOption(queryOptions *TGQueryOptionImpl) {
	obj.qryOption = queryOptions
}

func (obj *TGQueryImpl) SetQueryParameters(params map[string]interface{}) {
	obj.qryParameters = params
}

/////////////////////////////////////////////////////////////////
// Private functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) setQueryParameter(name string, value interface{}) {
	obj.qryParameters[name] = value
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

// Close closes the Query
func (obj *TGQueryImpl) Close() {
	_, _ = obj.qryConnection.CloseQuery(obj.qryHashId)
}

// Execute executes the Query
func (obj *TGQueryImpl) Execute() types.TGResultSet {
	rSet, _ := obj.qryConnection.ExecuteQueryWithId(obj.qryHashId, obj.qryOption)
	return rSet
}

// SetBoolean sets Boolean parameter
func (obj *TGQueryImpl) SetBoolean(name string, value bool) {
	obj.setQueryParameter(name, value)
}

// SetBytes sets Byte Parameter
func (obj *TGQueryImpl) SetBytes(name string, value []byte) {
	obj.setQueryParameter(name, value)
}

// SetChar sets Character Parameter
func (obj *TGQueryImpl) SetChar(name string, value string) {
	obj.setQueryParameter(name, value)
}

// SetDate sets Date Parameter
func (obj *TGQueryImpl) SetDate(name string, value time.Time) {
	obj.setQueryParameter(name, value)
}

// SetDouble sets Double Parameter
func (obj *TGQueryImpl) SetDouble(name string, value float64) {
	obj.setQueryParameter(name, value)
}

// SetFloat sets Float Parameter
func (obj *TGQueryImpl) SetFloat(name string, value float32) {
	obj.setQueryParameter(name, value)
}

// SetInt sets Integer Parameter
func (obj *TGQueryImpl) SetInt(name string, value int) {
	obj.setQueryParameter(name, value)
}

// SetLong sets Long Parameter
func (obj *TGQueryImpl) SetLong(name string, value int64) {
	obj.setQueryParameter(name, value)
}

// SetNull sets the parameter to null
func (obj *TGQueryImpl) SetNull(name string) {
	obj.setQueryParameter(name, nil)
}

// SetOption sets the Query Option
func (obj *TGQueryImpl) SetOption(options types.TGQueryOption) {
	obj.qryOption = options.(*TGQueryOptionImpl)
}

// SetShort sets Short Parameter
func (obj *TGQueryImpl) SetShort(name string, value int16) {
	obj.setQueryParameter(name, value)
}

// SetString sets String Parameter
func (obj *TGQueryImpl) SetString(name string, value string) {
	obj.setQueryParameter(name, value)
}

func (obj *TGQueryImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TGQueryImpl:{")
	buffer.WriteString(fmt.Sprintf("QryConnection: %+v", obj.qryConnection))
	buffer.WriteString(fmt.Sprintf(", QryHashId: %+v", obj.qryHashId))
	buffer.WriteString(fmt.Sprintf(", QryOption: %+v", obj.qryOption.String()))
	buffer.WriteString(fmt.Sprintf(", QryParameters: %+v", obj.qryParameters))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.qryConnection, obj.qryHashId, obj.qryOption, obj.qryParameters)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.qryConnection, &obj.qryHashId, &obj.qryOption, &obj.qryParameters)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
