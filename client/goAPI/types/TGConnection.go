package types

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
 * File name: TGConnection.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
	DecryptBuffer(encryptedBuf []byte) ([]byte, TGError)
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
	// ExecuteQuery executes an immediate query with associated query options
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
	// InsertEntity marks an ENTITY for insert operation. Upon commit, the entity will be inserted in the database
	InsertEntity(entity TGEntity) TGError
	// Rollback rolls back the current transaction on this connection
	Rollback() TGError
	// SetExceptionListener sets exception listener
	SetExceptionListener(listener TGConnectionExceptionListener)
	// UpdateEntity marks an ENTITY for update operation. Upon commit, the entity will be updated in the database
	// When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
	// dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
	// Calling multiple times, does not change the behavior.
	// The same entity cannot be updated on multiple connections. It will result an TGException of already associated to a connection.
	UpdateEntity(entity TGEntity) TGError
}
