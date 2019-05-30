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
 * <p/>
 * File name: TGConnection.java
 * Created on: 3/18/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGConnection.java 3158 2019-04-26 20:49:24Z kattaylo $
 */

package com.tibco.tgdb.connection;

import java.util.Collection;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGQueryOption;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.utils.TGProperties;

public interface TGConnection {


    /**
     * Connect to the Server
     * @throws com.tibco.tgdb.exception.TGException Any connection related exception.
     */
    void connect() throws TGException;

    /**
     * Disconnect from the Server.
     */
    void disconnect();

    /**
     * Set Exception Listener
     * @param lsnr - An Exception Listener object which receive connection related failures.
     */
    void setExceptionListener(TGConnectionExceptionListener lsnr);

    /**
     * Commit the Transaction on this Connection
     * @return TGResultSet indicating how many nodes/edges were inserted/updated/deleted.
     * @throws TGException - If the commit was successfull, then the ResultSet provides detailed error message
     */
    TGResultSet commit() throws TGException;

    /**
     * Rollback the Transaction on this Connection
     */
    void rollback();

    /**
     * Get an Entity given an UniqueKey for the Object.
     * If there are more than one entries, then the first one in the list is returned.
     * The key needs to have all the attributes set for which the index is defined.
     * see defining index for a Node, Edge or a Graph
     * This is synchronous non transactional operation. It does not hold any locks.
     * @param tgkey - An instance of a key using @see com.tibco.tgdb.TGGraphObjectFactory.createCompositeKey
     * @param option - properties affect the request such as batchsize, fetchsize and traversaldepth
     * @return TGEntity for the key specified
     * @throws com.tibco.tgdb.exception.TGException Throws an exception if there was any error while retrieving the object
     */
    TGEntity getEntity(TGKey tgkey, TGQueryOption option) throws TGException;

    /**
     * Get a result set of entities given an non-uniqueKey.
     * The key needs to have all the attributes set for which the index is defined.
     * see defining index for a Node, Edge or a Graph
     * This is synchronous non transactional operation. It does not hold any locks.
     * @param tgkey - An instance of a key using @see com.tibco.tgdb.TGGraphObjectFactory.createCompositeKey
     * @param properties - properties affect the request such as batchsize, fetchsize and traversaldepth
     * @return result set of entities
     * @throws com.tibco.tgdb.exception.TGException Throws an exception if there was any error while retrieving the object
     */
    TGResultSet getEntities(TGKey tgkey, TGProperties<String, String> properties) throws TGException;

    /**
     * Insert a new entity constructed using createNode/Edge 
     * @param entity
     * @throws TGException
     */
    void insertEntity(TGEntity entity) throws TGException;
    
    /**
     * Mark an Entity for Update operation. Call this method to associate the entity with a Connection
     * When commit is called, the object is resolved to check if it is dirty. Entity.setAttribute calls make the entity
     * dirty. If it is dirty, then the object is send to the server for update, otherwise it is ignored.
     *
     * Calling multiple times, does not change the behavior.
     * The same entity cannot be updated on multiple connections. It will throw a TGException of already associated to a connection.
     *
     * @param entity - The entity that was updated
     * @throws com.tibco.tgdb.exception.TGException Throws an exception if there was any error while updating the object
     */
    void updateEntity(TGEntity entity) throws TGException;

    /**
     * Mark an Entity for Delete for delete operation. Upon commit, the entity will be deleted.
     * @param entity the entity for delete
     * @throws TGException if could not be marked
     */
    void deleteEntity(TGEntity entity) throws TGException;

    /**
     * Execute a query in either tqql or gremlin format.  
     * Format is determined by either 'tgql://' or 'gremlin://' prefixes in the 'expr' argument. 
     * The format can also be specified by using 'tgdb.connection.defaultQueryLanguage' 
     * connection property with value 'tgql' or 'gremlin'. Prefix in the query expression is no needed
     * if connection property is used.
     *
     * @param expr query in tgql format or gremlin format
     * @param option Query options for executing. Can be null, then it will use the default option
     *
     * @return A navigatable ResultSet if successful in executing the query. The result set will indicate errors if
     * the query had any exceptions
     * @throws com.tibco.tgdb.exception.TGException Throws an exception if there was any error while querying the object
     */
	TGResultSet executeQuery(String expr, TGQueryOption option) throws TGException;
	
    /**
     * Execute an immediate Query with query options
     * The query option is place holder at this time
     *
     * @param expr A subset of SQL-92 where clause. @see com.tibco.tgdb.query.TGQuery
     * @param edgeFilter filter used for selecting edges to be returned
     * @param traversalCondition condition used for selecting edges to be traversed and returned
     * @param endCondition condition used to stop the traversal
     * @param option Query options for executing. Can be null, then it will use the default option
     *
     * @return A navigatable ResultSet if successful in executing the query. The result set will indicate errors if
     * the query had any exceptions
     * @throws com.tibco.tgdb.exception.TGException Throws an exception if there was any error while querying the object
     */
	TGResultSet executeQuery(String expr, String edgeFilter, String traversalCondition, String endCondition, TGQueryOption option) throws TGException;

    /**
     * Create a Resuable Query
     * @param expr A subset of SQL-92 where clause. @see com.tibco.tgdb.query.TGQuery
     * @return A resuable query object
     * @throws TGException if it can't create the Query
     */
    TGQuery createQuery(String expr) throws TGException;

    /**
     * Get the Graph Metadata
     * @param refresh meta data from the server if set to true
     * @return the Graph Metadata associated with the connection
     */
    TGGraphMetadata getGraphMetadata(boolean refresh) throws TGException;

    /**
     * Get the Graph Object Factory for Object creation.
     * @return TGGraphObjectFactory for creating objects
     * @throws TGException 
     */
    TGGraphObjectFactory getGraphObjectFactory() throws TGException;

}
