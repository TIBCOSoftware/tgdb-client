/*
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
 * File name: connection.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: connection.go 3542 2019-11-17 17:18:04Z nimish $
 */

package tgdb

type TGConnection interface {
	// Commit commits the current transaction on this connection
	Commit() (TGResultSet, TGError)
	// Connect establishes a network connection to the TGDB server
	Connect() TGError
	// CloseQuery closes a specific query and associated objects
	CloseQuery(queryHashId int64) (TGQuery, TGError)
	// CreateQuery creates a reusable query object that can be used to execute one or more statement
	CreateQuery(expr string) (TGQuery, TGError)
	// DecryptBuffer decrypts the encrypted buffer by sending a DecryptBufferRequest to the server
	DecryptBuffer(is TGInputStream) ([]byte, TGError)
	// DecryptEntity decrypts the encrypted entity using channel's data cryptographer
	DecryptEntity(entityId int64) ([]byte, TGError)
	// DeleteEntity marks an ENTITY for delete operation. Upon commit, the entity will be deleted from the database
	DeleteEntity(entity TGEntity) TGError
	// Disconnect breaks the connection from the TGDB server
	Disconnect() TGError
	// EncryptEntity encrypts the encrypted entity using channel's data cryptographer
	EncryptEntity(rawBuffer []byte) ([]byte, TGError)
	// ExecuteGremlinQuery executes a Gremlin Grammer-Based query with  query options
	ExecuteGremlinQuery(expr string, collection []interface{}, options TGQueryOption) ([]interface{}, TGError)
	// ExecuteQuery executes a query in either tqql or gremlin format.
	// Format is determined by either 'tgql://' or 'gremlin://' prefixes in the 'expr' argument.
	// The format can also be specified by using 'tgdb.connection.defaultQueryLanguage'
	// connection property with value 'tgql' or 'gremlin'. Prefix in the query expression is no needed
	// if connection property is used.
	ExecuteQuery(expr string, options TGQueryOption) (TGResultSet, TGError)
	// ExecuteQueryWithFilter executes an immediate query with specified filter & query options
	// The query option is place holder at this time
	// @param expr A subset of SQL-92 where clause
	// @param edgeFilter filter used for selecting edges to be returned
	// @param traversalCondition condition used for selecting edges to be traversed and returned
	// @param endCondition condition used to stop the traversal
	// @param option Query options for executing. Can be null, then it will use the default option
	ExecuteQueryWithFilter(expr string, edgeFilter string, traversalCondition string, endCondition string, options TGQueryOption) (TGResultSet, TGError)
	// ExecuteQueryWithId executes an immediate query for specified id & query options
	ExecuteQueryWithId(queryHashId int64, option TGQueryOption) (TGResultSet, TGError)
	// GetAddedList gets a list of added entities
	GetAddedList() map[int64]TGEntity
	// GetChangedList gets a list of changed entities
	GetChangedList() map[int64]TGEntity
	// GetChangedList gets the communication channel associated with this connection
	GetChannel() TGChannel
	// GetConnectionId gets connection identifier
	GetConnectionId() int64
	// GetConnectionProperties gets a list of connection properties
	GetConnectionProperties() TGProperties
	// GetEntities gets a result set of entities given an non-uniqueKey
	GetEntities(key TGKey, properties TGProperties) (TGResultSet, TGError)
	// GetEntity gets an Entity given an UniqueKey for the Object
	GetEntity(key TGKey, options TGQueryOption) (TGEntity, TGError)
	// GetGraphMetadata gets the Graph Metadata
	GetGraphMetadata(refresh bool) (TGGraphMetadata, TGError)
	// GetGraphObjectFactory gets the Graph Object Factory for Object creation
	GetGraphObjectFactory() (TGGraphObjectFactory, TGError)
	// GetLargeObjectAsBytes gets an Binary Large Object Entity given an UniqueKey for the Object
	GetLargeObjectAsBytes(entityId int64, decryptFlag bool) ([]byte, TGError)
	// GetRemovedList gets a list of removed entities
	GetRemovedList() map[int64]TGEntity
	// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
	InsertEntity(entity TGEntity) TGError
	// Rollback rolls back the current transaction on this connection
	Rollback() TGError
	// SetConnectionPool sets connection pool
	SetConnectionPool(connPool TGConnectionPool)
	// SetConnectionProperties sets connection properties
	SetConnectionProperties(connProps TGProperties)
	// SetExceptionListener sets exception listener
	SetExceptionListener(listener TGConnectionExceptionListener)
	// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
	// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
	// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
	// Calling multiple times, does not change the behavior.
	// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
	UpdateEntity(entity TGEntity) TGError
}

type TGConnectionPool interface {
	// AdminLock locks the connection pool so that the list of connections can be updated
	AdminLock()
	// AdminUnlock unlocks the connection pool so that the list of connections can be updated
	AdminUnlock()
	// Connect establishes connection from this pool of available/configured connections to the TGDB server
	// Exception could be BadAuthentication or BadUrl
	Connect() TGError
	// Disconnect breaks the connection from the TGDB server and returns the connection back to this connection pool for reuse
	Disconnect() TGError
	// Get returns a free connection from the connection pool
	// The property ConnectionReserveTimeoutSeconds or tgdb.connectionpool.ConnectionReserveTimeoutSeconds specifies the time
	// to wait in seconds. It has the following meaning
	// 0 :      Indefinite
	// -1 :     Immediate
	// &gt; :   That many seconds
	Get() (TGConnection, TGError)
	// GetPoolSize gets pool size
	GetPoolSize() int
	// ReleaseConnection frees the connection and sends back to the pool
	ReleaseConnection(conn TGConnection) (TGConnectionPool, TGError)
	// SetExceptionListener sets exception listener
	SetExceptionListener(lsnr TGConnectionExceptionListener)
}

type TGConnectionExceptionListener interface {
	// OnException registers a callback method with the exception
	OnException(ex TGError)
}

type TGTransaction interface {
	TGSerializable
	String() string
}

// ======= Various Connection Types =======
type TypeConnection int
const (
	TypeConventional TypeConnection = iota
	TypeAdmin
)

type TGConnectionFactory interface {
	CreateConnection(url, user, pwd string, env map[string]string) (TGConnection, TGError)
	CreateAdminConnection(url, user, pwd string, env map[string]string) (TGConnection, TGError)
	CreateConnectionPool(url, user, pwd string, poolSize int, env map[string]string) (TGConnectionPool, TGError)
	CreateConnectionPoolWithType(url, user, pwd string, poolSize int, env map[string]string, connType TypeConnection) (TGConnectionPool, TGError)
}




//var logger = logging.DefaultTGLogManager().GetLogger()

//

