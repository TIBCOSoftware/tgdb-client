"""
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
 *  File name :entityimpl.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: entityimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all the Attribute and its related classes
 """

from typing import *

import tgdb.model as tgmodel
import tgdb.pdu as tgpdu
import tgdb.exception as tgexception
from tgdb.impl.atomics import *
import tgdb.impl.attrimpl as tgattrimpl
import tgdb.log as tglog


class AbstractEntity(tgmodel.TGEntity):
    __entityId__: int = -1
    __entityType__: tgmodel.TGEntityType = None
    __isNew__: bool = False
    __isDeleted__: bool = False
    __version__: int = 0
    __graphMetadata__: tgmodel.TGGraphMetadata = None
    __virtualId__: int = -1
    __isInitialized__: bool = True
    __modifiedAttributes__: List[tgmodel.TGAttribute] = None
    __attributes__: Dict[str, tgmodel.TGAttribute] = None

    __gEntitySequencer__ = AtomicReference('q', 0)

    def __init__(self, gmd: tgmodel.TGGraphMetadata = None, entitytype: tgmodel.TGEntityType = None):
        if gmd is None:
            raise tgexception.TGException('Graph metadata is none')
        self.__graphMetadata__: tgmodel.TGGraphMetadata = gmd
        self.__isNew__: bool = True
        self.__version__: int = 0
        self.__isDeleted__: bool = False
        self.__isInitialized__: bool = True
        self.__entityId__: int = -1
        self.__entityType__: tgmodel.TGEntityType = entitytype
        self.__virtualId__: int = AbstractEntity.__gEntitySequencer__.decrement()
        self.__attributes__: Dict[str, tgmodel.TGAttribute] = dict()
        self.__modifiedAttributes__: List[tgattrimpl.AbstractAttribute] = list()

    @property
    def isNew(self) -> bool:
        return self.__isNew__

    @property
    def isDeleted(self) -> bool:
        return self.__isDeleted__

    def markDeleted(self):
        self.__isDeleted__ = True
        self.__isNew__ = False

    def resetModifiedAttributes(self):
        for attr in self.__modifiedAttributes__:
            attr.resetModified()
        self.__modifiedAttributes__.clear()

    @property
    def version(self) -> int:
        return self.__version__


    @property
    def attributes(self) -> List[tgmodel.TGAttribute]:
        attrlist: List[tgmodel.TGAttribute] = list()
        for k, v in self.__attributes__.items():
            if k == '@name':
                continue
            attrlist.append(v)
        return attrlist

    def getAttribute(self, key) -> tgmodel.TGAttribute:
        return self.__attributes__[key]

    def setAttribute(self, key, value):
        if key is None:
            raise tgexception.TGException("Attribute name specified is null")
        if value is None:
            return

        attr: tgmodel.TGAttribute = self.__attributes__[key] if key in self.__attributes__ else None
        if attr is None:
            desc: tgmodel.TGAttributeDescriptor = self.__graphMetadata__.attritubeDescriptors[key]
            if desc is None:
                raise tgexception.TGException("Attribute descriptor for {0} is not defined".format(key))
            attr = tgattrimpl.AbstractAttribute.createAttribute(self, desc, _value=value)

        self.__modifiedAttributes__.append(attr)

        attr._setValue(value)
        self.__attributes__[key] = attr

    def __setitem__(self, key, value):
        self.setAttribute(key, value)

    def __getitem__(self, key):
        return self.getAttribute(key).value

    def __hash__(self):
        return hash(self.primaryKey)

    def __eq__(self, other):
        if isinstance(other, AbstractEntity):
            return (self.entityType is not None and other.entityType is not None) and (self.entityType.id ==\
                           other.entityType.id) and (other.primaryKey == self.primaryKey)
        return False

    @property
    def entityType(self):
        return self.__entityType__

    def _get_entityId(self):
        return self.__entityId__

    def _set_entityId(self, value):
        self.__entityId__ = value
        self.__isNew__ = False

    entityId = property(_get_entityId, _set_entityId)

    @property
    def virtualId(self):
        return self.__virtualId__ if self.__isNew__ else self.__entityId__

    def isInitialized(self, value):
        self.__isInitialized__ = value
        self.__isNew__ = False

    @version.setter
    def version(self, value):
        self.__version__ = value

    @property
    def isNew(self):
        return self.__isNew__

    @property
    def graphMetadata(self) -> tgmodel.TGGraphMetadata:
        return self.__graphMetadata__

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        # Write the header for the entity
        ostream.writeBoolean(self.isNew)
        ostream.writeByte(self.entityKind.value)
        ostream.writeLong(self.virtualId)
        ostream.writeInt(self.version)
        ostream.writeInt(0 if self.entityType is None else self.entityType.id)

        # write the attributes
        ostream.writeInt(len(self.__modifiedAttributes__))
        for v in self.__modifiedAttributes__:
            v.writeExternal(ostream)

    def readExternal(self, istream: tgpdu.TGInputStream):
        _: bool = istream.readBoolean()     # should be False, but apparently not guaranteed.
        self.__isNew__ = False  # Force it to false - TGDB-504

        kind = istream.readByte()
        if self.entityKind.value != kind:
            raise tgexception.TGException('Invalid object for deserialization. Expecting {0}, found {1} with byte {2}.'.
                              format(self.entityKind.name, tgmodel.TGEntityKind.fromKindId(kind).name, kind))
        self.entityId = istream.readLong()
        self.version = istream.readInt()
        typeid = istream.readInt()
        if typeid != 0:
            et: tgmodel.TGEntityType = self.graphMetadata.types[typeid]
            if et is None:
                raise tgexception.TGException('cannot lookup entity desc {0} from graph metadata cache'.format(typeid))
            self.__entityType__ = et
        else:
            raise tgexception.TGException('cannot lookup entity desc {0} from graph metadata cache'.format(typeid))

        count = istream.readInt()
        for i in range(0, count):
            attr: tgmodel.TGAttribute = tgattrimpl.AbstractAttribute.createFromStream(self, istream)
            self.__attributes__[attr.descriptor.name] = attr

        return

    @classmethod
    def readNewEntityFromStream(cls, instream: tgpdu.TGInputStream, gof: tgmodel.TGGraphObjectFactory,
                                ek: tgmodel.TGEntityKind = None) -> tgmodel.TGEntity:
        if ek is None:
            ek = tgmodel.TGEntityKind.fromKindId(instream.readByte())
        ret: AbstractEntity = None
        if ek is tgmodel.TGEntityKind.Node:
            ret = gof.createNode(None)
        elif ek is tgmodel.TGEntityKind.Edge:
            ret = gof.createEdge(None, None, tgmodel.DirectionType.BiDirectional)
        elif ek is tgmodel.TGEntityKind.InvalidKind:
            tglog.gLogger.log(tglog.TGLevel.Warning, "Cannot parse invalid entity!")
        else:
            tglog.gLogger.log(tglog.TGLevel.Warning, "Unsupported entity kind: %s", str(ek))
        if ret is not None:
            pos = instream.position
            ret.readExternal(instream)
            if ret.entityId in instream.referencemap:
                instream.position = pos
                instream.referencemap[ret.entityId].readExternal(instream)
            else:
                instream.referencemap[ret.entityId] = ret
        return ret


class NodeImpl(AbstractEntity, tgmodel.TGNode):

    def __init__(self, gmd: tgmodel.TGGraphMetadata, entype):
        super().__init__(gmd, entype)
        self.__edges__: Dict[int, tgmodel.TGEdge] = {}
        return

    def addEdge(self, edge: tgmodel.TGEdge):
        if self not in edge.vertices:
            raise tgexception.TGException("Cannot add edge to a vertex that is not one of the edge's vertices!")
        self.__edges__[edge.virtualId] = edge

    def hasEdge(self, edge: tgmodel.TGEdge) -> bool:
        return edge.virtualId in self.__edges__

    def removeEdge(self, edge: tgmodel.TGEdge):
        if edge.virtualId not in self.__edges__:
            raise tgexception.TGException("Cannot remove edge not already in node!")
        del self.__edges__[edge.virtualId]

    @property
    def primaryKey(self) -> Any:
        ret: Any = self.virtualId
        if self.entityType is not None:                     # Only check if the node has an entity type.
            nodetype: tgmodel.TGNodeType = self.entityType  # if not, then there are no primary keys
            retSet: Set[Any] = set()
            for attrDesc in nodetype.pkeyAttributeDescriptors:
                try:
                    retSet.add(self[attrDesc.name])
                except KeyError:
                    pass
            if len(retSet) > 0:             # only replace if we found some keys or the pkeys where non-empty
                ret = frozenset(retSet)
        return ret

    @property
    def entityKind(self) -> tgmodel.TGEntityKind:
        return tgmodel.TGEntityKind.Node

    """
    * f can be specified as an example 
    f: Callable[[tgmodel.TGEdge], bool] = lambda x: x.directiontype == tgmodel.DirectionType.Directed
    """
    def getEdges(self, f: Callable[[tgmodel.TGEdge], bool] = lambda _: True) -> List[tgmodel.TGEdge]:
        return list(filter(f, self.__edges__.values()))

    @property
    def edges(self) -> List[tgmodel.TGEdge]:
        return self.getEdges()

    @property
    def directedEdges(self):
        f: Callable[[tgmodel.TGEdge], bool] = lambda x: x.directiontype == tgmodel.DirectionType.Directed
        return self.getEdges(f)

    @property
    def undirectedEdges(self):
        f: Callable[[tgmodel.TGEdge], bool] = lambda x: x.directiontype == tgmodel.DirectionType.UnDirected
        return self.getEdges(f)

    @property
    def bidirectedEdges(self):
        f: Callable[[tgmodel.TGEdge], bool] = lambda x: x.directiontype == tgmodel.DirectionType.BiDirectional
        return self.getEdges(f)

    def markDeleted(self):
        super().markDeleted()
        keyset: Set[int] = set(self.__edges__.keys())
        key: int
        for key in keyset:
            self.__edges__[key].markDeleted()

    def __writeIgnoreErrors(self, ostream: tgpdu.TGOutputStream, edge: tgmodel.TGEdge):
        try:
            ostream.writeLong(edge.virtualId)
        except IOError:
            pass

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        pos: int = ostream.position
        ostream.writeInt(0)
        super().writeExternal(ostream)
        ostream.writeInt(sum(1 for edge in self.edges if edge.isNew))
        [self.__writeIgnoreErrors(ostream, edge) for edge in self.edges if edge.isNew]
        curPos: int = ostream.position
        ostream.writeIntAt(pos, curPos - pos)

    def readExternal(self, istream: tgpdu.TGInputStream):
        _: int = istream.readInt()          # Unused buffer length, could be used as a rudimentary checksum maybe?
        super().readExternal(istream)
        edgeCount: int = istream.readInt()
        for i in range(edgeCount):
            edge: EdgeImpl = None
            edgeId: int = istream.readLong()
            refMap: Dict[int, tgmodel.TGEntity] = istream.referencemap
            if refMap is not None and edgeId in refMap:
                edge = refMap[edgeId]
            if edge is None:
                edge = EdgeImpl(self.graphMetadata, None, None)
                edge.entityId = edgeId
                edge.isInitialized(False)
                edge.entityId = edgeId
                if refMap is not None:
                    refMap[edgeId] = edge
            self.__edges__[edge.virtualId] = edge
        self.isInitialized(True)


class EdgeImpl(AbstractEntity, tgmodel.TGEdge):
    __fromnode__: tgmodel.TGNode = None
    __tonode__: tgmodel.TGNode = None
    __dirtype__: tgmodel.DirectionType = tgmodel.DirectionType.Directed

    def __init__(self, gmd: tgmodel.TGGraphMetadata, fromnode: tgmodel.TGNode, tonode: tgmodel.TGNode, dirtype=tgmodel.DirectionType.Directed, entype: tgmodel.TGEdgeType=None):
        super().__init__(gmd, entype)
        self.__fromnode__: NodeImpl = fromnode
        self.__tonode__: NodeImpl = tonode
        self.__dirtype__: tgmodel.DirectionType = dirtype if entype is None else entype.directiontype
        self.__entityId__: int = -1
        self.__addSelfToNodes()

    @property
    def primaryKey(self) -> Any:
        ret: Any = self.virtualId
        return ret

    def _get_entityId(self):
        return super()._get_entityId

    def _set_entityId(self, id: int):
        # Need to do this so that the node.__edge__ dictionary gets updated, but also need to check if self is in
        # node.edges because
        self.__deleteSelfFromNodes()
        super()._set_entityId(id)
        self.__addSelfToNodes()

    entityId = property(_get_entityId, _set_entityId)

    def __deleteSelfFromNodes(self):
        if self.__fromnode__ is not None and self.__fromnode__.hasEdge(self):
            self.__fromnode__.removeEdge(self)
        # Need to check if the nodes are the same, because they might be. For example: self-loops.
        if self.__tonode__ is not None and self.__tonode__.hasEdge(self):
            self.__tonode__.removeEdge(self)

    def __addSelfToNodes(self):
        if self.__fromnode__ is not None:
            self.__fromnode__.addEdge(self)
        if self.__tonode__ is not None:
            self.__tonode__.addEdge(self)

    @property
    def vertices(self) -> tuple:
        return self.__fromnode__, self.__tonode__

    @property
    def fromVertex(self) -> tgmodel.TGNode:
        return self.__fromnode__

    @property
    def toVertex(self) -> tgmodel.TGNode:
        return self.__tonode__

    @property
    def directiontype(self) -> tgmodel.DirectionType:
        return self.__dirtype__

    @property
    def entityKind(self) -> tgmodel.TGEntityKind:
        return tgmodel.TGEntityKind.Edge

    def markDeleted(self):
        self.__deleteSelfFromNodes()
        super().markDeleted()
        self.__fromnode__ = None
        self.__tonode__ = None

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        pos = ostream.position
        ostream.writeInt(0)
        super().writeExternal(ostream)
        ostream.writeByte(self.__dirtype__.value)
        ostream.writeLong(self.__fromnode__.virtualId)
        ostream.writeLong(self.__tonode__.virtualId)
        curpos = ostream.position
        ostream.writeIntAt(pos, curpos - pos)

    def __upsertEdgeToNode(self, id: int, refMap: Dict[int, AbstractEntity]) -> NodeImpl:
        node: NodeImpl = None
        if refMap is not None and id in refMap:
            node = refMap[id]
        if node is None:
            node = NodeImpl(self.__graphMetadata__, None)
            node.entityId = id
            node.isInitialized(False)
            if refMap is not None:
                refMap[id] = node
        return node

    def readExternal(self, istream: tgpdu.TGInputStream):
        _: int = istream.readInt()      # Unused buffer length, at least for now.
        super().readExternal(istream)
        dir: int = istream.readByte()

        if dir is tgmodel.DirectionType.UnDirected.value:
            self.__dirtype__ = tgmodel.DirectionType.UnDirected
        elif dir is tgmodel.DirectionType.Directed.value:
            self.__dirtype__ = tgmodel.DirectionType.Directed
        elif dir is tgmodel.DirectionType.BiDirectional.value:
            self.__dirtype__ = tgmodel.DirectionType.BiDirectional

        refMap: Dict[int, AbstractEntity] = istream.referencemap
        self.__fromnode__ = self.__upsertEdgeToNode(istream.readLong(), refMap)
        self.__tonode__ = self.__upsertEdgeToNode(istream.readLong(), refMap)
        self.isInitialized(True)


# TODO delete the following class later, but for now, it is only in because we need it for bulk import/export
class SparseEdgeImpl(AbstractEntity, tgmodel.TGEdge):

    def __init__(self, gmd: tgmodel.TGGraphMetadata, fromnode: int, tonode: int, dirtype: int,
                 entype: tgmodel.TGEdgeType):
        super().__init__(gmd, entype)
        self.__fromId__: int = fromnode
        self.__toId__: int = tonode
        self.__dirType__: int = dirtype

    @property
    def vertices(self) -> tuple:
        return self.__fromId__, self.__toId__

    @property
    def directiontype(self) -> tgmodel.DirectionType:
        return tgmodel.DirectionType[self.__dirType__]

    @property
    def entityKind(self) -> tgmodel.TGEntityKind:
        return tgmodel.TGEntityKind.Edge

    @property
    def primaryKey(self) -> Any:
        ret: Any = self.virtualId
        return ret

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        pos = ostream.position
        ostream.writeInt(0)
        super().writeExternal(ostream)
        ostream.writeByte(self.__dirType__)
        ostream.writeLong(self.__fromId__)
        ostream.writeLong(self.__toId__)
        curpos = ostream.position
        ostream.writeIntAt(pos, curpos - pos)


class CompositeKeyImpl(tgmodel.TGKey):

    def __init__(self, gmd: tgmodel.TGGraphMetadata, nodetype: tgmodel.TGEntityType):
        self.__gmd__ = gmd
        self.__typename__: str = None
        self.__attrmap__: Dict[str, tgmodel.TGAttribute] = {}
        entype = nodetype
        self.__typename__ = entype.name

    def __setitem__(self, key, value):
        if key is None or value is None:
            raise tgexception.TGException('Invalid args specified. key={0}, value={1}'.format(key, value))

        desc = self.__gmd__.registry[key]
        if desc is None or not isinstance(desc, tgmodel.TGAttributeDescriptor):
            raise tgexception.TGException('desc {0} is null or not a valid attribute descriptor'.format(str(desc)))

        attr = tgattrimpl.AbstractAttribute.createAttribute(None, desc)
        attr._setValue(value)
        self.__attrmap__[key] = attr

        return

    def matches(self, entity: tgmodel.TGEntity) -> bool:
        retV: bool = True
        try:
            for key in self.__attrmap__:
                selfVal = self.__attrmap__[key].value
                otherVal = entity.getAttribute(key).value
                if not (selfVal == otherVal):
                    retV = False
        except KeyError:
            retV = False
        return retV


    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        if self.__typename__ is None:
            ostream.writeBoolean(False)
        else:
            ostream.writeBoolean(True)
            ostream.writeUTF(self.__typename__)
        ostream.writeShort(len(self.__attrmap__))
        for k in self.__attrmap__:
            self.__attrmap__[k].writeExternal(ostream)

    def readExternal(self, istream: tgpdu.TGInputStream):
        raise tgexception.TGException('Not supported operation')



