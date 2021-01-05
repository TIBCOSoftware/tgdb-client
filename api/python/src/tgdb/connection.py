"""
.. Necessary to reject for documentation building
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
 *  File name :connection.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: connection.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates channel interfaces
 """

import abc
import typing

from tgdb.utils import TGProperties
import tgdb.bulkio as tgbulk
import tgdb.query as tgquery
import tgdb.model as tgmodel
import tgdb.channel as tgchan


class TGConnection(abc.ABC):
    """ This connection object represents a single session with the database.

    First, to acquire a TGConnection object, call the TGConnectionFactory's class method createConnection with the
    appropriate parameters. Do not directly instantiate a TGConnection object.

    Next, to connect to the database, simply call the TGConnection object's connect method. To insert an entity into the
    database on the next commit, call insertEntity. To update an entity already in the database, call updateEntity. To
    delete an entity, call deleteEntity. When you are ready to commit all of the database changes, call commit. To
    query, call executeQuery.

    Properties
    -------------------------------
    graphMetadata: tgdb.model.TGGraphMetadata
        Gets this session's graph metadata instance. Not settable or deletable. Useful for inquries about what types
        (nodetypes, edgetypes, and attrdescs) that are already in the database (and therefore what you can use).
    graphObjectFactory: tgdb.model.TGGraphObjectFactory
        Gets this session's graph object factory instance. Not settable or deletable. Useful for creating other
        entities.
    """

    @abc.abstractmethod
    def connect(self):
        """Connects to the database.

        Does not support reconnect, will need to create a new connection object for that.

        Tip: Use this in a try-finally block with disconnect to always cleanly sever the connection with the database
        regardless of any exceptions.
        """
        pass

    @abc.abstractmethod
    def disconnect(self):
        """Disconnects from the database.

        Once called, another connection object needs to be created to connect to the database again.

        Tip: Use this in a try-finally block with connect to always cleanly sever the connection with the database
        regardless of any exceptions.
        """
        pass

    @abc.abstractmethod
    def commit(self):
        """Commit any changes to the database."""
        pass

    @abc.abstractmethod
    def refreshMetadata(self):
        """Refreshes this connection's metadata."""

    @abc.abstractmethod
    def getEntity(self, key: tgmodel.TGKey, option: tgquery.TGQueryOption = tgquery.DefaultQueryOption) ->\
            tgmodel.TGEntity:
        """Get an entity from the database.

        :param key: Key representing the primary key and/or the attributes set on a unique index.
        :param option: A list of query options to determine how to acquire the entity (default:
            tgdb.query.DefaultQueryOption)

        :returns: the entity if one was found.
        """
        pass

    @abc.abstractmethod
    def insertEntity(self, entity: tgmodel.TGEntity):
        """Inserts an entity into the database.

        Newly created entities will not automatically be inserted into the database.

        Need to call commit before any changes are finalized.

        :param entity: The entity to insert into the database.
        """

    @abc.abstractmethod
    def updateEntity(self, entity: tgmodel.TGEntity):
        """Updates an entity already in the database

        Need to call commit before any changes are finalized.

        :param entity: The entity to update in the database.
        """

    @abc.abstractmethod
    def deleteEntity(self, entity: tgmodel.TGEntity):
        """Updates an entity already in the database

        Need to call commit before any changes are finalized.

        :param entity: The entity to delete from the database.
        """
    
    @abc.abstractmethod
    def createQuery(self, query: str) -> tgquery.TGQuery:
        """Creates a query for the database.

        Will be used in the future as a way to write compiled/parameterized queries.
        """

    @abc.abstractmethod
    def executeQuery(self, query: str, option: tgquery.TGQueryOption = tgquery.DefaultQueryOption) ->\
            tgquery.TGResultSet:
        """Executes a query that is not parameterized.

        :param query: The query to execute. Should be of a supported format like TGQL or Gremlin. If its of a TGQL
            format, please prepend with tgql:// and if it is of a Gremlin, then prepend with a gremlin://.
        :param option: The options for this query.
        :returns: The result of the query. Will handle any cursoring if necessary
        """

    @abc.abstractmethod
    def startImport(self, loadopt: typing.Union[str, tgbulk.TGLoadOptions] = tgbulk.TGLoadOptions.Insert,
                    erroropt: typing.Union[str, tgbulk.TGErrorOptions] = tgbulk.TGErrorOptions.Stop,
                    dateformat: typing.Union[str, tgbulk.TGDateFormat] = tgbulk.TGDateFormat.YMD,
                    props: typing.Optional[TGProperties] = None) -> tgbulk.TGBulkImport:
        """Starts the bulk import process with the server.

        :param loadopt: Whether to only insert entities or update entities if that primary key already is in the
            database.
        :param erroropt: Whether to fail or ignore errors.
        :param dateformat: How to interpret dates.
        :param props: TGProperties for starting the import process.

        :returns: Returns a bulk import instance that handles the remainder of the bulk import.
        """

    @abc.abstractmethod
    def startExport(self, props: typing.Optional[TGProperties] = None) -> "tgbulk.TGBulkExport":
        """Starts the bulk export process with the server.

        :param props: TGProperties for starting the export process.
        :returns: Returns a bulk export instance that handles the remainder of the bulk export.
        """

    @abc.abstractmethod
    def rollback(self):
        """Rolls back the database."""

    @property
    @abc.abstractmethod
    def graphMetadata(self) -> tgmodel.TGGraphMetadata:
        """Gets the graph's metadata."""

    @property
    @abc.abstractmethod
    def graphObjectFactory(self) -> tgmodel.TGGraphObjectFactory:
        """Gets a factory object for creating new entities."""

    @property
    @abc.abstractmethod
    def linkState(self) -> tgchan.LinkState:
        """Gets the link state."""

    @property
    @abc.abstractmethod
    def outboxaddr(self) -> str:
        """Gets outbox (server) address."""

    @property
    @abc.abstractmethod
    def connectedUsername(self) -> str:
        """Gets the connected user's name."""

    def __enter__(self):
        """Required for using the 'with' paradigm."""
        self.connect()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """Required for using the 'with' paradigm."""
        self.disconnect()


class TGConnectionFactory:
    """Factory for creating the TGConnection instance."""

    @classmethod
    def createConnection(cls, url: str, username: str, password: str, dbName: str = None, env: typing.Dict[str, str] = None) \
            -> TGConnection:
        """Creates a single connection object.

        :param dbName:
        :param url: The URL for the database to connect to. When running on localhost, it should look like
            tcp://localhost:8222 for an insecure connection.
        :param username: The user's name of the database to use to connect to the database.
        :param password: The user's password for the database server.
        :param env: The environment to set up the connection with, contains any additional properties.

        :returns: The connection object requested.
        """

        import tgdb.impl.connectionimpl as tgconn

        return tgconn.ConnectionImpl(url, username, password, dbName, env)

    @classmethod
    def createAdminConnection(cls, url: str, username: str, password: str, dbName: str = None, env: typing.Dict[str, str] = None):
        """Creates a single administrative connection object.

        :param url: The URL for the database to connect to. When running on localhost, it should look like
            tcp://localhost:8222 for an insecure connection.
        :param username: The user's name of the database to use to connect to the database.
        :param password: The user's password for the database server.
        :param dbName: The database to connect to
        :param env: The environment to set up the connection with, contains any additional properties.

        :returns: The administrative connection object requested.
        :rtype: tgdb.admin.TGAdminConnection
        """

        import tgdb.impl.adminimpl as tgadm

        return tgadm.AdminConnectionImpl(url, username, password, dbName, env)
