package connection

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/channel"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"math"
	"strconv"
	"sync"
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
 * File name: TGConnectionPool.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Various Connection Pool States =======
const (
	ConnectionPoolInitialized = iota
	ConnectionPoolConnecting
	ConnectionPoolConnected
	ConnectionPoolInUse
	ConnectionPoolDisconnecting
	ConnectionPoolDisconnected
	ConnectionPoolStopped
)

// 0 :      Indefinite
// -1 :     Immediate
// &gt; :   That many seconds
const (
	Immediate  = "-1"
	Indefinite = "0"
	IMMEDIATE  = 0
	INFINITE   = math.MaxInt32 // math.MaxInt64
)

type ConnectionPoolImpl struct {
	adminLock             sync.RWMutex // rw-lock for synchronizing read-n-update of connection pool properties
	connectReserveTimeOut time.Duration
	connList              []*TGDBConnection // Total Available Connections (Active + Dead/ToBeReused)
	connType              TypeConnection
	chanPool              chan *TGDBConnection
	poolProperties        *utils.SortedProperties
	consumers             map[int64]*TGDBConnection           // Active/In-Use Connections
	exceptionListener     types.TGConnectionExceptionListener // Function Pointer
	poolSize              int
	poolState             int
	useDedicateChannel    bool
}

var gInstance *ConnectionPoolImpl
var once sync.Once

func defaultTGConnectionPool() *ConnectionPoolImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ConnectionPoolImpl{})

	//once.Do(func() {
		gInstance := &ConnectionPoolImpl{
			connList:  make([]*TGDBConnection, 0),
			connType:  TypeConventional,
			consumers: make(map[int64]*TGDBConnection, 0),
		}
		gInstance.poolSize, _ = strconv.Atoi(utils.GetConfigFromKey(utils.ConnectionPoolDefaultPoolSize).GetDefaultValue())
		gInstance.useDedicateChannel, _ = strconv.ParseBool(utils.GetConfigFromKey(utils.ConnectionPoolUseDedicatedChannelPerConnection).GetDefaultValue())
	//})
	return gInstance
}

func NewTGConnectionPool(url types.TGChannelUrl, poolSize int, props *utils.SortedProperties, connType TypeConnection) *ConnectionPoolImpl {
	logger.Log(fmt.Sprintf("Entering ConnectionPoolImpl:NewTGConnectionPool w/ ChannelURL: '%+v', Poolsize: '%d'", url.GetUrlAsString(), poolSize))
	cp := defaultTGConnectionPool()
	logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool w/ Default Connection Pool: '%s'", cp.String()))
	cp.connType = connType
	cp.poolProperties = props
	cp.chanPool = make(chan *TGDBConnection, poolSize+2)
	cp.poolSize = poolSize
	timeoutStr := utils.GetConfigFromKey(utils.ConnectionReserveTimeoutSeconds).GetDefaultValue()
	if timeoutStr == Immediate {
		cp.connectReserveTimeOut = time.Second * IMMEDIATE
	} else if timeoutStr == Indefinite {
		cp.connectReserveTimeOut = time.Second * INFINITE
	} else {
		timeout, _ := strconv.Atoi(timeoutStr)
		cp.connectReserveTimeOut = time.Second * time.Duration(timeout)
	}
	var ch types.TGChannel
	channelFactory := channel.GetChannelFactoryInstance()
	for i := 0; i < cp.poolSize; i++ {
		if ch == nil || cp.useDedicateChannel {
			logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to channelFactory.CreateChannelWithUrlProperties() for URL: '%s'", url.GetUrlAsString()))
			// Create a channel from channel factory
			channel1, err := channelFactory.CreateChannelWithUrlProperties(url, props)
			if err != nil {
				errMsg := fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:NewTGConnectionPool Unable to create a channel for URL: '%s' via channel factory - '%+v'", url, err.Error())
				logger.Error(errMsg)
				continue
			}
			ch = channel1
		}
		//logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to NewTGDBConnection() for Channel: '%+v' and Properties: '%+v'", ch, props))
		// Create a connection
		var conn *TGDBConnection
		switch connType {
		case TypeConventional:
			conn = NewTGDBConnection(cp, ch, props)
		//case TypeAdmin:
			//conn = NewTGDBAdminConnection(cp, ch, props)
		default:
			conn = NewTGDBConnection(cp, ch, props)
		}
		conn.SetConnectionPool(cp)
		conn.SetConnectionProperties(props)
		// Add it in the pool to initialize the pool with a set number of initialized connections
		cp.connList = append(cp.connList, conn)
		//logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl:NewTGConnectionPool - about to add conn: '%+v' to the pool", conn.String()))
		cp.chanPool <- conn
	} // End of For loop for Pool Size
	cp.poolState = ConnectionPoolInitialized
	logger.Log(fmt.Sprintf("Returning ConnectionPoolImpl:NewTGConnectionPool w/ Connection Pool: '%s'", cp.String()))
	return cp
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGConnectionPool
/////////////////////////////////////////////////////////////////

func (obj *ConnectionPoolImpl) AdminLock() {
	obj.adminLock.RLock()
}

func (obj *ConnectionPoolImpl) AdminUnlock() {
	obj.adminLock.RUnlock()
}

// GetConnectionList returns all the connections = Active/In-Use + Un-used/Initialized
func (obj *ConnectionPoolImpl) GetConnectionList() []*TGDBConnection {
	return obj.connList
}

// GetConnectionProperties returns all the connection properties
func (obj *ConnectionPoolImpl) GetConnectionProperties() types.TGProperties {
	return obj.poolProperties
}

// GetActiveConnections returns all the Active/In-Use connections
func (obj *ConnectionPoolImpl) GetActiveConnections() map[int64]*TGDBConnection {
	return obj.consumers
}

// GetNoOfActiveConnections returns the count of Active/In-Use connections
func (obj *ConnectionPoolImpl) GetNoOfActiveConnections() int {
	return len(obj.consumers)
}

// GetPoolState returns current state of the connection pool
func (obj *ConnectionPoolImpl) GetPoolState() int {
	return obj.poolState
}

// GetConnection returns an available connection from the pool that is NOT being used or nil if timeout elapses
func (obj *ConnectionPoolImpl) GetConnection() (types.TGConnection, types.TGError) {
	logger.Log(fmt.Sprintf("Entering ConnectionPoolImpl:GetConnection for Pool Type: '%+v'", obj.connType))
	obj.adminLock.Lock()
	defer obj.adminLock.Unlock()

	if obj.GetNoOfActiveConnections() == obj.GetPoolSize() {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:GetConnection - as all the connections in the pool are in use."))
		errMsg := "ConnectionPoolImpl has already exhausted its limit. All the connections in the pool are in use. Please wait and retry."
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl:GetConnection - about to pull connection from obj.ConnPool"))
	var conn *TGDBConnection
	// Search through all available connections within the pooled set to find out which one is free to service next
	select {
	case conn = <-obj.chanPool:
		//logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl::GetConnection about to verify the state of the connection: '%+v' pulled from obj.ConnPool:", conn))
		// Proceed with the connection that is NOT being used already a.k.a. part of consumer connection map
		if _, ok := obj.consumers[conn.GetConnectionId()]; !ok {
			// Add it in the pool to initialize the pool with a set number of initialized connections
			obj.consumers[conn.GetConnectionId()] = conn
			break
		}
	case <-time.After(obj.connectReserveTimeOut):
		logger.Warning(fmt.Sprintf("WARNING: Returning ConnectionPoolImpl:GetConnection - trying to get a connection after wating for '%+v'", obj.connectReserveTimeOut))
		//log.Warning("Timed out trying to get connection from the pool")
		errMsg := fmt.Sprintf("Timed out trying to get a connection after wating for %v", obj.connectReserveTimeOut)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}

	logger.Log(fmt.Sprintf("Returning ConnectionPoolImpl:GetConnection w/ Connection: '%+v'", conn.String()))
	return conn, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnectionPool
/////////////////////////////////////////////////////////////////

// Connect establishes connection from this pool of available/configured connections to the TGDB server
// Exception could be BadAuthentication or BadUrl
func (obj *ConnectionPoolImpl) Connect() types.TGError {
	logger.Log(fmt.Sprint("Entering ConnectionPoolImpl:Connect"))
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	if obj.poolState == ConnectionPoolConnected {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:Connect - ConnectionPoolImpl is already connected. Disconnect and then reconnect."))
		errMsg := "ConnectionPoolImpl is already connected. Disconnect and then reconnect"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}
	// Set the state to connecting
	obj.poolState = ConnectionPoolConnecting

	// Attempt to connect using each of the available connections in the pool
	for connState, conn := range obj.consumers {
		if connState == ConnectionPoolConnecting || connState == ConnectionPoolConnected ||
			connState == ConnectionPoolInUse || connState == ConnectionPoolDisconnecting {
			continue // Skip this connection and go to next one in the pool
		}
		// Proceed when the connection is either in Initialized OR Disconnected (OR Stopped???) state
		logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl::Connect Consumer Loop about to conn.Connect() using connection: '%+v'", conn))
		err := conn.Connect()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:Connect - unable to conn.Connect() w/ '%s'", err.Error()))
			// TODO: Revisit later - Decide what error to throw and whether to continue / break
			return err
		}
		break
	}

	// Set the state to connecting
	obj.poolState = ConnectionPoolConnected
	logger.Log(fmt.Sprint("Returning ConnectionPoolImpl:Connect"))
	return nil
}

// Disconnect breaks the connection from the TGDB server and returns the connection back to this connection pool for reuse
func (obj *ConnectionPoolImpl) Disconnect() types.TGError {
	logger.Log(fmt.Sprint("Entering ConnectionPoolImpl:Disconnect"))
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	if obj.poolState != ConnectionPoolConnected {
		logger.Error(fmt.Sprint("ERROR: Returning ConnectionPoolImpl:Disconnect - ConnectionPoolImpl is NOT connected."))
		errMsg := fmt.Sprintf("ConnectionPoolImpl is not connected. State is: %d", obj.poolState)
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}
	// Set the state to connecting
	obj.poolState = ConnectionPoolDisconnecting

	// Attempt to connect using each of the active connections in the pool
	for _, conn := range obj.connList {
		logger.Log(fmt.Sprintf("Inside ConnectionPoolImpl::Disconnect active connection Loop about to conn.Disconnect() using connection: '%+v'", conn))
		err := conn.Disconnect()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning ConnectionPoolImpl:Disconnect - unable to conn.Disconnect() w/ '%s'", err.Error()))
			// TODO: Revisit later - Decide what error to throw and whether to continue / break
			return err
		}
	}

	// Set the state to connecting
	obj.poolState = ConnectionPoolDisconnected
	logger.Log(fmt.Sprint("Returning ConnectionPoolImpl:Disconnect"))
	return nil
}

// Get a free connection
// The property ConnectionReserveTimeoutSeconds or tgdb.connectionpool.ConnectionReserveTimeoutSeconds specifies the time
// to wait in seconds. It has the following meaning
// 0 :      Indefinite
// -1 :     Immediate
// &gt; :   That many seconds
func (obj *ConnectionPoolImpl) Get() (types.TGConnection, types.TGError) {
	return obj.GetConnection()
}

// GetPoolSize gets pool size
func (obj *ConnectionPoolImpl) GetPoolSize() int {
	return obj.poolSize
}

// ReleaseConnection frees the connection and sends back to the pool
func (obj *ConnectionPoolImpl) ReleaseConnection(conn types.TGConnection) (types.TGConnectionPool, types.TGError) {
	logger.Log(fmt.Sprint("Entering ConnectionPoolImpl:ReleaseConnection"))
	obj.adminLock.RLock()
	defer obj.adminLock.RUnlock()

	logger.Log(fmt.Sprint("Inside ConnectionPoolImpl::ReleaseConnection Consumer Loop about to remove connection from consumer list"))
	delete(obj.consumers, conn.(*TGDBConnection).GetConnectionId())
	logger.Log(fmt.Sprint("Inside ConnectionPoolImpl::ReleaseConnection Consumer Loop about to return connection to obj.ConnPool"))
	obj.chanPool <- conn.(*TGDBConnection)

	logger.Log(fmt.Sprint("Returning ConnectionPoolImpl:ReleaseConnection"))
	return obj, nil
}

// SetExceptionListener sets exception listener
func (obj *ConnectionPoolImpl) SetExceptionListener(listener types.TGConnectionExceptionListener) {
	obj.exceptionListener = listener
}

func (obj *ConnectionPoolImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ConnectionPoolImpl:{")
	buffer.WriteString(fmt.Sprintf("ConnectReserveTimeOut: %+v", obj.connectReserveTimeOut))
	buffer.WriteString(fmt.Sprintf(", ConnType: %+v", obj.connType))
	buffer.WriteString(fmt.Sprintf(", ConnList: %+v", obj.connList))
	buffer.WriteString(fmt.Sprintf(", ConnPool: %+v", obj.chanPool))
	//buffer.WriteString(fmt.Sprintf(", PoolProperties: %+v", obj.poolProperties))
	buffer.WriteString(fmt.Sprintf(", Consumers: %+v", obj.consumers))
	buffer.WriteString(fmt.Sprintf(", PoolSize: %d", obj.poolSize))
	buffer.WriteString(fmt.Sprintf(", PoolState: %d", obj.poolState))
	buffer.WriteString(fmt.Sprintf(", UseDedicateChannel: %+v", obj.useDedicateChannel))
	buffer.WriteString("}")
	return buffer.String()
}
