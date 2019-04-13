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
 * File name: ServerInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ServerInfoImpl struct {
	cacheInfo        *CacheStatisticsImpl
	databaseInfo     *DatabaseStatisticsImpl
	memoryInfo       *ServerMemoryInfoImpl
	netListenersInfo []TGNetListenerInfo
	serverStatusInfo *ServerStatusImpl
	transactionsInfo *TransactionStatisticsImpl
}

// Make sure that the ServerInfoImpl implements the TGServerInfo interface
var _ TGServerInfo = (*ServerInfoImpl)(nil)

func DefaultServerInfoImpl() *ServerInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerInfoImpl{})

	return &ServerInfoImpl{}
}

func NewServerInfoImpl(_cacheInfo *CacheStatisticsImpl, _databaseInfo *DatabaseStatisticsImpl,
	_memoryInfo *ServerMemoryInfoImpl, _netListenersInfo []TGNetListenerInfo,
	_serverStatusInfo *ServerStatusImpl, _transactionsInfo *TransactionStatisticsImpl) *ServerInfoImpl {
	newServerInfo := DefaultServerInfoImpl()
	newServerInfo.cacheInfo = _cacheInfo
	newServerInfo.databaseInfo = _databaseInfo
	newServerInfo.memoryInfo = _memoryInfo
	newServerInfo.netListenersInfo = _netListenersInfo
	newServerInfo.serverStatusInfo = _serverStatusInfo
	newServerInfo.transactionsInfo = _transactionsInfo
	return newServerInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("CacheInfo: '%+v'", obj.cacheInfo))
	buffer.WriteString(fmt.Sprintf(", DatabaseInfo: '%+v'", obj.databaseInfo))
	buffer.WriteString(fmt.Sprintf(", MemoryInfo: '%+v'", obj.memoryInfo))
	buffer.WriteString(fmt.Sprintf(", NetListenersInfo: '%+v'", obj.netListenersInfo))
	buffer.WriteString(fmt.Sprintf(", ServerStatusInfo: '%+v'", obj.serverStatusInfo))
	buffer.WriteString(fmt.Sprintf(", TransactionsInfo: '%+v'", obj.transactionsInfo))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerInfo
/////////////////////////////////////////////////////////////////

// GetCacheInfo returns cache statistics information from server
func (obj *ServerInfoImpl) GetCacheInfo() TGCacheStatistics {
	return obj.cacheInfo
}

// GetDatabaseInfo returns database statistics information from server
func (obj *ServerInfoImpl) GetDatabaseInfo() TGDatabaseStatistics {
	return obj.databaseInfo
}

// GetMemoryInfo returns object corresponding to specific memory type
func (obj *ServerInfoImpl) GetMemoryInfo(memType MemType) TGMemoryInfo {
	return obj.memoryInfo.GetServerMemoryInfo(memType)
}

// GetNetListenersInfo returns a collection of information on NetListeners
func (obj *ServerInfoImpl) GetNetListenersInfo() []TGNetListenerInfo {
	return obj.netListenersInfo
}

// GetServerStatus returns the information on Server Status including name, version etc.
func (obj *ServerInfoImpl) GetServerStatus() TGServerStatus {
	return obj.serverStatusInfo
}

// GetTransactionsInfo returns transaction statistics from server including processed transaction count, successful transaction count, average processing time etc.
func (obj *ServerInfoImpl) GetTransactionsInfo() TGTransactionStatistics {
	return obj.transactionsInfo
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.cacheInfo, obj.databaseInfo, obj.memoryInfo,
		obj.netListenersInfo, obj.serverStatusInfo, obj.transactionsInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.cacheInfo, &obj.databaseInfo, &obj.memoryInfo,
		&obj.netListenersInfo, &obj.serverStatusInfo, &obj.transactionsInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
