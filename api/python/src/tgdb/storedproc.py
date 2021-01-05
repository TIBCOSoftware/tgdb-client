"""
.. Necessary to reject for documentation building
 * Copyright (c) 2020 TIBCO Software Inc. All rights reserved.
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
 *   File name   :storedproc.py
 *   Created on  :01/21/20
 *   Created by  :katie
 *   SVN Id      :$Id:$
"""


from functools import wraps
from abc import ABC, abstractmethod
from tgdb.model import TGNode, TGKey, TGEdge, TGEntity, DirectionType
from tgdb.query import TGQueryOption, DefaultQueryOption, TGResultSet
import typing


def tgstoredproc(fn):
    """
    This decorates a stored procedure and marks it so that the TIBCO Graph Database can execute it when called. The \
    name of the function is what you should pass to the ``execsp`` Gremlin step.

    :param fn: A function that should be added to the stored procedures that a client may execute. Must have the first\
    argument be annotated as a TGGraphContext. It must also have properly annotated return types that match up with the\
    TGDB return language. The return type must match the annotated return type.
    :type fn: typing.Callable[[TGGraphContext, .... ], typing.Any]
    """

    @wraps(fn)
    def __sp_wrap__1z(*args, **kwargs):
        if not isinstance(args[0], TGGraphContext):
            raise SystemError("Stored procedure definition must include TGGraphContext variable as first argument")
        return fn(*args, **kwargs)
    return __sp_wrap__1z


def tgstoredprocfull(fn):
    """
    This function decorates a stored procedure and marks it so that the TIBCO Graph Database can execute it when\
    called. The name of the function is what you should pass to the ``execsp`` Gremlin step.\
    This function decorator is like tgstoredproc, but it also indicates that the results from the Gremlin steps\
    preceding the execsp() step should be aggregated before calling the stored procedure.

    :param fn: A function that should be added to the stored procedures that a client may execute. Must have the first\
    argument be annotated as a TGGraphContext. It must also have properly annotated return types that match up with the\
    TGDB return language. The return type must match the annotated return type.
    :type fn: typing.Callable[[TGGraphContext, .... ], typing.Any]
    """

    @wraps(fn)
    def __sp_wrap__1zz(*args, **kwargs):
        if not isinstance(args[0], TGGraphContext):
            raise SystemError("Stored procedure definition must include TGGraphContext variable as first argument")
        return fn(*args, **kwargs)
    return __sp_wrap__1zz


def tgtrigger(fn):
    """
    This function decorates a trigger and marks it so that the TIBCO Graph Database can execute it when a transaction\
    is committed.

    :param fn: A function that should be added to the triggers that execute on a transaction. Must have the first\
    argument be annotated as a TGGraphContext. The second argument must be annotated as a TGTransactionChangeList.
    :type fn: typing.Callable[[TGGraphContext, TGTransactionChangelist], None]
    """

    @wraps(fn)
    def __sp_wrap__1y(*args, **kwargs):
        if not isinstance(args[0], TGGraphContext):
            raise SystemError("Trigger function must pass TGGraphContext variable")
        if not isinstance(args[1], TGTransactionChangelist):
            raise SystemError("Trigger function must pass TGTransactionChangelist variable")
        return fn(*args, **kwargs)
    return __sp_wrap__1y


class TGTransactionChangelist(ABC):
    """Represents all of the entities that are about to be inserted, deleted, or changed on a transaction."""

    @abstractmethod
    def addEntity(self, ent: TGEntity, opcode: str):
        """
        Adds the entity to this transaction.

        :param ent: The entity to add to a change list in the transaction.
        :type ent: tgdb.model.TGEntity
        :param opcode: The operation code that determines whether to insert, change, or delete the entity. Must be one\
        of "insert", "update", or "delete"
        :type opcode: str
        """

    @abstractmethod
    def removeEntity(self, ent: TGEntity, opcode):
        """
        Remove the entity from this transaction.

        :param ent: The entity to remove from a change list in the transaction.
        :type ent: tgdb.model.TGEntity
        :param opcode: The operation code that determines whether to remove the entity from the insert, change, or\
        delete changelist. Must be one of "insert", "update", or "delete"
        :type opcode: str
        """

    @abstractmethod
    def getInsertedEntities(self) -> typing.Iterable[TGEntity]:
        """
        Gets all newly inserted entities for this transaction. See :ref:`stored procedure entity<sp-tgent>` for \
        potential restrictions.

        :returns: All of the newly inserted entities for this transaction.
        :rtype: typing.Iterable[tgdb.model.TGEntity]
        """

    @abstractmethod
    def getModifiedEntities(self) -> typing.Iterable[TGEntity]:
        """
        Get all modified entities for this transaction. See :ref:`stored procedure entity<sp-tgent>` for potential \
        restrictions.

        :returns: All of the updated entities for this transaction.
        :rtype: typing.Iterable[tgdb.model.TGEntity]
        """

    @abstractmethod
    def getDeletedEntities(self) -> typing.Iterable[TGEntity]:
        """
        Get all removed entities for this transaction. See :ref:`stored procedure entity<sp-tgent>` for potential \
        restrictions.

        :returns: All of the deleted entities for this transaction.
        :rtype: typing.Iterable[tgdb.model.TGEntity]
        """


class TGGraphContext(ABC):
    """
    Acts as a combination of TGGraphObjectFactory and TGConnection. This is the class that is the primary interface\
    between the stored procedure/trigger and the server. Since this code is executed on the server, the interface does\
    not include any of the connection details required by the server. It also has access to all of the nodes and edges\
    in this step in the query or across the whole database when a trigger??.
    """
    
    @abstractmethod
    def getNodes(self, nodetype: str = None) -> typing.Iterable[TGNode]:
        """
        Gets all nodes a part of this Graph Context. See :ref:`stored procedure node<sp-tgnode>` for potential \
        restrictions.
        
        :param nodetype: The node type of the nodes to get. If set to None, then gets all of the nodes available. If it\
        is a string, then it must be a valid nodetype name.
        :type nodetype: str
        :return: An iterator of all of the nodes that are of that nodetype or all that are available.
        :rtype: typing.Iterable[tgdb.model.TGNode]
        """

    @abstractmethod
    def getEdges(self, edgetype: str = None) -> typing.Iterable[TGEdge]:
        """
        Gets all edges a part of this Graph Context. See :ref:`stored procedure edge<sp-tgedge>` for potential \
        restrictions.

        :param edgetype: The edgetype of the edges to get. If set to None, then gets all of the edges available. If it\
        is a string, then it must be a valid edgetype name.
        :type edgetype: str
        :return: An iterator of all of the edges that are of that edgetype or all that are available.
        :rtype: typing.Iterable[tgdb.model.TGEdge]
        """

    @abstractmethod
    def createNode(self, nodetype: str) -> TGNode:
        """
        Creates a new node with the given nodetype. Does not insert the node into the database. See \
        :ref:`stored procedure node<sp-tgnode>` for potential restrictions.

        :param nodetype: The nodetype for the new node. Must be a valid nodetype.
        :type nodetype: str
        :returns: The new node.
        :rtype: tgdb.model.TGNode
        """

    @abstractmethod
    def createEdge(self, fromnode: TGNode, tonode: TGNode,
                   dirtype=DirectionType.Directed, edgetype: str = None) -> TGEdge:
        """
        Creates an edge from the given fromnode to the given tonode with the direction type and edgetype if provided.\
        Does not insert the edge into the database. See :ref:`stored procedure edge<sp-tgedge>` for potential \
        restrictions.

        :param fromnode: The from node for this new edge.
        :type fromnode: tgdb.model.TGNode
        :param tonode: The to node for this new edge.
        :type tonode: tgdb.model.TGNode
        :param dirtype: The direction for the new edge. Should only set when edgetype is set to None.
        :type dirtype: tgdb.model.DirectionType
        :param edgetype: The edgetype for the new node. If set to None, then uses the default edgetype for the\
        direction type set. If set to a string, it must be a valid edgetype.
        :type edgetype: str
        :returns: The new edge.
        :rtype: tgdb.model.TGEdge
        """

    @abstractmethod
    def createCompositeKey(self, name: str) -> TGKey:
        """
        Creates a composite key on the nodetype given by name. See :ref:`stored procedure key<sp-tgkey>`.

        :param name: The nodetype name. Must be a valid nodetype.
        :type name: str
        :returns: A new key to use for getting a particular node of a particular nodetype.
        :rtype: tgdb.model.TGKey
        """

    @abstractmethod
    def getEntity(self, key: TGKey, option: TGQueryOption = DefaultQueryOption) -> TGEntity:
        """
        Gets the entity with the given key and query options specified.

        :param key: The key to use when looking up a particular entity in the database. Currently only supports looking\
        up nodes.
        :type key: tgdb.model.TGKey
        :param option: The options for acquiring the entity requested.
        :type option: tgdb.query.TGQueryOption
        :returns: The entity requested.
        :rtype: tgdb.model.TGEntity
        """

    @abstractmethod
    def executeQuery(self, query: str, option: TGQueryOption = DefaultQueryOption) -> TGResultSet:
        """
        Executes the query string with the query options specified.

        :param query: The query to run on the server.
        :type query: str
        :param option: The query options to use.
        :type option: tgdb.query.TGQueryOption
        :returns: The result of the query.
        :rtype: tgdb.query.TGResultSet
        """

    @abstractmethod
    def insertEntity(self, entity: TGEntity):
        """
        Inserts the entity. Will only be committed after the `commit` method is called.

        :param entity: The entity to insert into the server on the next commit call.
        :type entity: tgdb.model.TGEntity
        :returns: Nothing
        :rtype: None
        """

    @abstractmethod
    def updateEntity(self, entity: TGEntity):
        """
        Updates the entity. Will only be committed after the `commit` method is called.

        :param entity: The entity to update on the server on the next commit call.
        :type entity: tgdb.model.TGEntity
        :returns: Nothing
        :rtype: None
        """

    @abstractmethod
    def deleteEntity(self, entity: TGEntity):
        """
        Deletes the entity. Will only be committed after the `commit` method is called.

        :param entity: The entity to delete from the server on the next commit call.
        :type entity: tgdb.model.TGEntity
        :returns: Nothing
        :rtype: None
        """

    @abstractmethod
    def commit(self):
        """
        Commits all changes scheduled in the changelist.

        :returns: Nothing
        :rtype: None
        """
