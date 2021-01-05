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
 *  File name :gmdimpl.py
 *  Created on: 5/29/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: gmdimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all the Graph metadata and its related operations
 """
import typing

import tgdb.impl.entityimpl as tgentimpl
import tgdb.model as tgmodel
from typing import *
from abc import *
import tgdb.query as tgquery
import tgdb.impl.gmdimpl as tggmdimpl
import tgdb.pdu as tgpdu
import tgdb.model as tgmodel
import tgdb.impl.entityimpl as tgentimpl
import tgdb.log as tglog
import tgdb.exception as tgexception
import tgdb.impl.queryimpl as tgqueryimpl
import tgdb.impl.attrimpl as tgattrimpl
from tgdb.model import TGEntityType, TGAttributeDescriptor, TGSPParamType, TGSystemType, TGSPParam, TGSPReturnElemType, \
    TGRolePermFlag, TGPermissionType, TGRolePermission, TGPrivilege
from tgdb.pdu import TGInputStream, TGOutputStream


class TypeRegistry(object):

    def __init__(self):
        self.__inttypemap__ = dict()
        self.__strtypemap__ = dict()
        self.__nodetypes__ = dict()
        self.__edgetypes__ = dict()
        self.__attrdescs__ = dict()

    def __getitem__(self, item) -> tgmodel.TGSystemObject:
        '''
        i = self.__inttypemap__
        s = self.__strtypemap__
        n = self.__nodetypes__
        e = self.__edgetypes__
        a = self.__attrdescs__
        '''
        if isinstance(item, str):
            return self.__strtypemap__[item]
        elif isinstance(item, int):
            return self.__inttypemap__[item]
        else:
            raise tgexception.TGException('invalid type: {0} found. Expected string or int '.format(type(item)))

    def addItem(self, v: tgmodel.TGSystemObject):
        self.__inttypemap__[v.id] = v
        self.__strtypemap__[v.name] = v
        if v.systemType == tgmodel.TGSystemType.AttributeDescriptor:
            self.__attrdescs__[v.name] = v
        elif v.systemType == tgmodel.TGSystemType.NodeType:
            self.__nodetypes__[v.name] = v
        elif v.systemType == tgmodel.TGSystemType.EdgeType:
            self.__edgetypes__[v.name] = v
        else:
            return  # do nothing
        return

    def getItem(self, k, default=None):
        try:
            return self[k]
        except KeyError:
            return default


class AbstractEntityType(tgmodel.TGEntityType):

    def __init__(self, registry):
        if registry is None:
            raise tgexception.TGException('registry is none')
        self.__registry__ = registry
        self.__attrdescs__: Dict[str, tgmodel.TGAttributeDescriptor] = {}
        self.__id__: int = 0
        self.__name__: str = None
        self.__pagesize: int = -1

    @property
    def pagesize(self) -> int:
        return self.__pagesize

    @property
    def attributeDescriptors(self) -> List[tgmodel.TGAttributeDescriptor]:
        return list(self.__attrdescs__.values())

    def getAttributeDescriptor(self, name) -> tgmodel.TGAttributeDescriptor:
        return self.__attrdescs__[name]

    @property
    def systemType(self):
        return tgmodel.TGSystemType.InvalidType

    @property
    def name(self):
        return self.__name__

    @property
    def id(self):
        return self.__id__

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        raise tgexception.TGException("writeExternam is not implemented for AbstractEntity ")

    def readExternal(self, istream: tgpdu.TGInputStream):
        value = istream.readByte()
        type: tgmodel.TGSystemType = tgmodel.TGSystemType.fromValue(value)
        if type == tgmodel.TGSystemType.InvalidType:
            tglog.gLogger.log(tglog.TGLevel.Error, "Entity desc input stream has invalid desc value : %d", value)

        self.__id__ = istream.readInt()
        self.__name__ = istream.readUTF()
        self.__pagesize = istream.readInt()    # We are ignoring this

        attrcount = istream.readShort()
        for i in range(0, attrcount):
            name = istream.readUTF()
            self.__attrdescs__[name] = self.__registry__[name]


class NodeTypeImpl(AbstractEntityType, tgmodel.TGNodeType):

    def __init__(self, registry):
        super().__init__(registry)
        self.__pkeylist__: List[tgmodel.TGAttributeDescriptor] = list()
        self.__idxids__: List[int] = list()
        self.__numentries: int = 0

    @property
    def systemType(self):
        return tgmodel.TGSystemType.NodeType

    @property
    def pkeyAttributeDescriptors(self) -> List[tgmodel.TGAttributeDescriptor]:
        return self.__pkeylist__

    def readExternal(self, istream: tgpdu.TGInputStream):
        super().readExternal(istream)

        cnt = istream.readShort()
        for i in range(0, cnt):
            name = istream.readUTF()
            self.__pkeylist__.append(self.__registry__[name])

        cnt = istream.readShort()
        for i in range(0, cnt):
            self.__idxids__.append(istream.readInt())

        self.__numentries = istream.readLong()

    @property
    def numberEntities(self) -> int:
        return self.__numentries


class EdgeTypeImpl(AbstractEntityType, tgmodel.TGEdgeType):

    def __init__(self, registry):
        super().__init__(registry)
        self.__fromtype__: tgmodel.TGNodeType = None
        self.__totype__: tgmodel.TGNodeType = None
        self.__dirtype__: tgmodel.DirectionType = None
        self.__numentries: int = None

    @property
    def systemType(self):
        return tgmodel.TGSystemType.EdgeType

    @property
    def directiontype(self) -> tgmodel.DirectionType:
        return self.__dirtype__

    def fromNodeType(self) -> tgmodel.TGNodeType:
        return self.__fromtype__

    def toNodeType(self) -> tgmodel.TGNodeType:
        return self.__totype__

    def readExternal(self, istream: tgpdu.TGInputStream):
        super().readExternal(istream)
        self.__fromtype__ = self.__registry__.getItem(istream.readInt())
        self.__totype__ = self.__registry__.getItem(istream.readInt())
        dir = istream.readByte()
        if dir == 0:
            self.__dirtype__ = tgmodel.DirectionType.UnDirected
        elif dir == 1:
            self.__dirtype__ = tgmodel.DirectionType.Directed
        else:
            self.__dirtype__ = tgmodel.DirectionType.BiDirectional
        self.__numentries = istream.readLong()      # number of entries

    @property
    def numberEntities(self) -> int:
        return self.__numentries


class IndexImpl(tgmodel.TGIndex):

    def __init__(self, reg: TypeRegistry):
        self.__reg = reg
        self.__unique = False
        self.__attrdescs: typing.List[tgmodel.TGAttributeDescriptor] = []
        self.__types: typing.List[tgmodel.TGEntityType] = []
        self.__blocksize: int = 512
        self.__status = ''
        self.__name = ''
        self.__id = -1
        self.__numents = -1

    @property
    def unique(self) -> bool:
        return self.__unique

    @property
    def attributeDescriptors(self) -> typing.List[TGAttributeDescriptor]:
        return self.__attrdescs

    @property
    def onTypes(self) -> typing.List[TGEntityType]:
        return self.__types

    @property
    def blockSize(self) -> int:
        return self.__blocksize

    @property
    def status(self) -> str:
        return self.__status

    @property
    def systemType(self):
        return tgmodel.TGSystemType.Index

    @property
    def name(self):
        return self.__name

    @property
    def id(self):
        return self.__id

    @property
    def numEntries(self) -> int:
        return self.__numents

    def writeExternal(self, ostream: TGOutputStream):
        raise tgexception.TGProtocolNotSupported("Not allowed to write indices to server!")

    def readExternal(self, istream: TGInputStream):
        sys_type = istream.readByte()
        if sys_type != tgmodel.TGSystemType.Index.value:
            raise tgexception.TGException('Tried to read index when type is not an index!')
        self.__id = istream.readInt()
        self.__name = istream.readUTF()
        self.__unique = istream.readBoolean()
        numAttrs = istream.readInt()
        self.__attrdescs = []
        for i in range(numAttrs):
            attrname = istream.readUTF()
            if attrname not in self.__reg.__attrdescs__:
                raise tgexception.TGException('Unknown attribute descriptor received from server.')
            self.__attrdescs.append(self.__reg.__attrdescs__[attrname])
        numTypes = istream.readInt()
        self.__types = []
        for i in range(numTypes):
            typename = istream.readUTF()
            if typename in self.__reg.__nodetypes__:
                self.__types.append(self.__reg.__nodetypes__[typename])
            elif typename in self.__reg.__edgetypes__:
                self.__types.append(self.__reg.__edgetypes__[typename])
            else:
                raise tgexception.TGException('Unknown type name received from server.')
        self.__blocksize = istream.readInt()
        self.__numents = istream.readLong()

        self.__status = istream.readChars()

"""
class SPReturnTypeImpl(tgmodel.TGSPReturnType):

    def __init__(self, eType: tgmodel.TGSPReturnElemType):
        self.__etype = eType

    @property
    def elemType(self) -> tgmodel.TGSPReturnElemType:
        return self.__etype
"""


class SPParamImpl(tgmodel.TGSPParam):

    def __init__(self, name: str, pType: tgmodel.TGSPParamType):
        self.__name = name
        self.__ptype = pType

    @property
    def name(self) -> str:
        return self.__name

    @property
    def paramType(self) -> TGSPParamType:
        return self.__ptype


class StoredProcedureImpl(tgmodel.TGStoredProcedure):

    def __init__(self):
        self.__module: str = ''
        self.__returnType = TGSPReturnElemType.SPRInvalid  # todo set to SPReturnType object
        self.__params: typing.List[SPParamImpl] = []
        self.__name: str = ''
        self.__id: int = -1

    @property
    def module(self) -> str:
        return self.__module

    @property
    def returnType(self) -> TGSPReturnElemType:
        return self.__returnType

    @property
    def params(self) -> typing.List[TGSPParam]:
        return self.__params

    @property
    def systemType(self) -> TGSystemType:
        return TGSystemType.StoredProcedure

    @property
    def name(self) -> str:
        return self.__name

    @property
    def id(self):
        return self.__id

    def writeExternal(self, ostream: TGOutputStream):
        raise tgexception.TGException("Cannot write stored procedure as output (at least not yet...).")

    def readExternal(self, istream: TGInputStream):
        typ = TGSystemType.fromValue(istream.readByte())

        if typ != self.systemType:
            raise tgexception.TGException("Illegal type passed into stored procedure read!")

        self.__id = istream.readInt()
        self.__name = istream.readUTF()

        self.__module = istream.readChars()

        self.__returnType = istream.readChars()

        numParams = istream.readLong()

        self.__params = []

        for i in range(numParams):
            name = istream.readChars()
            pType = tgmodel.TGSPParamType.fromId(istream.readInt())
            self.__params.append(SPParamImpl(name, pType))


class RolePermissionImpl(tgmodel.TGRolePermission, tgpdu.TGSerializable):

    def __init__(self):

        self.__granted = set()
        self.__revoked = set()
        self.__sysID = -1
        self.__flag = TGRolePermFlag.Invalid

    def __hash__(self):
        return self.__sysID

    def __eq__(self, other):
        if isinstance(other, RolePermissionImpl):
            return other.__sysID == self.__sysID
        return TypeError

    @property
    def grantedPerms(self) -> typing.FrozenSet[TGPermissionType]:
        return frozenset(self.__granted)

    @property
    def revokedPerms(self) -> typing.FrozenSet[TGPermissionType]:
        return frozenset(self.__revoked)

    @property
    def sysID(self):
        return self.__sysID

    @property
    def flag(self) -> TGRolePermFlag:
        return self.__flag

    def writeExternal(self, ostream: TGOutputStream):
        raise tgexception.TGException("Cannot write role permission as output (at least not yet...).")

    def readExternal(self, istream: TGInputStream):
        self.__sysID = istream.readInt()
        granted = istream.readUnsignedLong()
        revoked = istream.readUnsignedLong()
        flag = istream.readUnsignedByte()
        self.__granted = TGPermissionType.permissionsFromByte(granted)
        self.__revoked = TGPermissionType.permissionsFromByte(revoked)
        self.__flag = TGRolePermFlag.flagFromInt(flag)


class RoleImpl(tgmodel.TGRole, tgpdu.TGSerializable):

    def __init__(self):
        self.__id: int = -1
        self.__name: str = ''
        self.__privileges: typing.Set[tgmodel.TGPrivilege] = set()
        self.__permissions: typing.List[RolePermissionImpl] = []
        self.__default: RolePermissionImpl = None
        self.__syscatalog: RolePermissionImpl = None

    @property
    def privileges(self) -> typing.FrozenSet[TGPrivilege]:
        return frozenset(self.__privileges)

    @property
    def permissions(self) -> typing.List[TGRolePermission]:
        return self.__permissions

    @property
    def default(self) -> TGRolePermission:
        return self.__default

    @property
    def syscatalog(self) -> TGRolePermission:
        return self.__syscatalog

    @property
    def systemType(self) -> TGSystemType:
        return TGSystemType.Role

    @property
    def name(self) -> str:
        return self.__name

    @property
    def id(self):
        return self.__id

    def writeExternal(self, ostream: TGOutputStream):
        raise tgexception.TGProtocolNotSupported("Not allowed to write roles to server!")

    def readExternal(self, istream: TGInputStream):
        sys_type = istream.readByte()
        if sys_type != tgmodel.TGSystemType.Role.value:
            raise tgexception.TGException('Tried to read role when type is not an role!')
        self.__id = istream.readInt()
        self.__name = istream.readUTF()

        privileges = istream.readUnsignedLong()

        self.__privileges = TGPrivilege.privilegesFromInt(privileges)

        self.__syscatalog = RolePermissionImpl()
        self.__default = RolePermissionImpl()

        self.__syscatalog.readExternal(istream)
        self.__default.readExternal(istream)

        number_role_perms = istream.readUnsignedInt()

        self.__permissions = []

        for i in range(number_role_perms):
            perm = RolePermissionImpl()

            perm.readExternal(istream)

            self.__permissions.append(perm)


class GraphMetadataImpl(tgmodel.TGGraphMetadata):

    def __init__(self, gof: tgmodel.TGGraphObjectFactory):
        self.__gof__ = gof
        self.__registry__ = TypeRegistry()

    @property
    def registry(self) -> TypeRegistry:
        return self.__registry__

    @registry.setter
    def registry(self, value):
        if not isinstance(value, TypeRegistry):
            raise tgexception.TGException('Expecting a typeregistry instance, found:{0}'.format(type(value)))
        self.__registry__ = value

    def defaultEdgetype(self, edgeDir: tgmodel.DirectionType):
        edgetype: tgmodel.TGEdgeType
        if edgeDir == tgmodel.DirectionType.Directed:
            edgetype = self.edgetypes['$def-directed-edge$']
        elif edgeDir == tgmodel.DirectionType.BiDirectional:
            edgetype = self.edgetypes['$def-bidirected-edge$']
        elif edgeDir == tgmodel.DirectionType.UnDirected:
            edgetype = self.edgetypes['$def-undirected-edge$']
        else:
            raise tgexception.TGException("Unknown direction type: {0}".format(str(self)))
        return edgetype

    @property
    def nodetypes(self) -> Dict[str, tgmodel.TGNodeType]:
        return self.__registry__.__nodetypes__

    @property
    def edgetypes(self) -> Dict[str, tgmodel.TGEdgeType]:
        return self.__registry__.__edgetypes__

    @property
    def attritubeDescriptors(self) -> Dict[str, tgmodel.TGAttributeDescriptor]:
        return self.__registry__.__attrdescs__

    @property
    def types(self) -> Dict[int, tgmodel.TGSystemObject]:
        return self.registry.__inttypemap__

    @property
    def graphObjectFactory(self) -> tgmodel.TGGraphObjectFactory:
        return self.__gof__

    def writeExternal(self, ostream: tgpdu.TGOutputStream):
        return

    def readExternal(self, istream: tgpdu.TGInputStream):
        return


class GraphObjectFactoryImpl(tgmodel.TGGraphObjectFactory):

    def __init__(self, connection):
        self.__connection__ = connection
        self.__gmd__ = GraphMetadataImpl(self)
        self.__managedEntities__: Dict[int, List[tgmodel.TGEntity]] = {}

    def createNode(self, nodetype: Union[str, tgmodel.TGNodeType]) -> tgmodel.TGNode:
        if nodetype is not None and isinstance(nodetype, str):
            nodetype = self.__gmd__.nodetypes[nodetype]
        return tgentimpl.NodeImpl(self.__gmd__, nodetype)

    def createEdge(self, fromnode: tgmodel.TGNode, tonode: tgmodel.TGNode, dirtype: tgmodel.DirectionType =
                    tgmodel.DirectionType.Directed, edgetype: Union[tgmodel.TGEdgeType, str] = None) -> tgmodel.TGEdge:
        if edgetype is not None and isinstance(edgetype, str):
            edgetype = self.__gmd__.edgetypes[edgetype]
        elif edgetype is None:
            edgetype = self.__gmd__.defaultEdgetype(dirtype)
        return tgentimpl.EdgeImpl(self.__gmd__, fromnode, tonode, dirtype, edgetype)

    def createCompositeKey(self, nodetype: Union[str, tgmodel.TGEntityType]) -> tgmodel.TGKey:
        if nodetype is None:
            raise tgexception.TGException("Must have keys operating on a particular entitytype.")
        if isinstance(nodetype, str):
            nodetype = self.__gmd__.registry[nodetype]
        return tgentimpl.CompositeKeyImpl(self.__gmd__, nodetype)

    @property
    def connection(self):
        return self.__connection__

    @property
    def graphmetadata(self):
        return self.__gmd__
