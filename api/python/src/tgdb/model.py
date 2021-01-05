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
 *  File name :model.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: model.py 3245 2019-05-31 16:13:47Z ssubrama $
 *
 *  This file encapsulates all the interfaces defined the model java package.
 """

import abc
import enum
import typing
import tgdb.pdu as tgpdu
import tgdb.exception as tgexception


class TGAttributeType(enum.Enum):
    """Represents the different type of attributes."""
    Invalid = (0, None, None)
    Boolean = (1, '?', None)         # A single bit representing the truth value
    Byte    = (2, 'b', None)         # 8bit octet
    Char    = (3, 'c', None)         # Fixed 8-Bit octet of N length
    Short   = (4, 'h', None)         # 16bit
    Integer = (5, 'i', 'int')        # 32bit signed integer
    Long    = (6, 'l', None)         # 64bit
    Float   = (7, 'f', None)         # 32bit float
    Double  = (8, 'd', None)         # 64bit float
    Number  = (9, 'n', None)         # Number with precision
    String  = (10, 's', None)        # Varying length String < 64K
    Date    = (11, 'D', None)        # Only the Date part from the Calendar class is considered
    Time    = (12, 'e', None)        # Time hh:mm: ss.nnn
    TimeStamp = (13, 'a', None)      # Compatible with the SQL Timestamp field.
    Clob    = (14, None, None)       # Character -UTF-8 encoded string or large length > 64K
    Blob    = (15, None, None)       # Binary object - a stream of octets, unsigned 8bit char with length.

    def __init__(self, val: int, char: str, shortname: str):
        self.__val = val
        self.__char = char
        self.__shortname = shortname

    @property
    def identifer(self) -> int:
        """Gets the identifier for this attribute."""
        return self.__val

    @property
    def char(self) -> str:
        """
        Gets a single-character representation of this attribute. Should correspond to something that a query returns.
        """
        return self.__char

    @property
    def shortname(self) -> str:
        """Gets the short-hand name of this attribute type."""
        return self.__shortname

    @classmethod
    def fromTypeId(cls, typeId: int):
        """Gets the corresponding TGAttributeType based on the typeid"""
        for attrtype in TGAttributeType:
            if attrtype.identifer == typeId:
                return attrtype
        return TGAttributeType.Invalid

    @classmethod
    def fromTypeChar(cls, charId: str):
        """Gets the corresponding TGAttributeType based on the character identifier"""
        for attrtype in TGAttributeType:
            if attrtype.char == charId:
                return attrtype
        return TGAttributeType.Invalid

    @classmethod
    def fromTypeName(cls, name):
        """Gets the corresponding TGAttributeType based on the name of the attribute."""
        for attrname, value in TGAttributeType.__members__.items():
            if value.shortname is None and attrname.lower() == name.lower():
                return value
            elif value.shortname is not None and value.shortname.lower() == name.lower():
                return value
        return TGAttributeType.Invalid


class TGSystemType(enum.Enum):
    """Describes the various kinds of types that the server recognizes."""
    InvalidType = -1
    AttributeDescriptor = 0
    NodeType = 1
    EdgeType = 2
    Index = 3
    Principal = 4
    Role = 5
    Sequence = 6
    StoredProcedure = 7

    MaxSysObjectTypes = 64

    @classmethod
    def fromValue(cls, value):
        """Gets the kind of type from the integer identifier."""
        for systype in TGSystemType:
            if systype.value == value:
                return systype
        return TGSystemType.InvalidType


class TGEntityKind(enum.Enum):
    """Describes the different kinds of the entities in the database."""
    InvalidKind = 0
    Entity      = 1
    Node        = 2
    Edge        = 3
    Graph       = 4
    HyperEdge   = 5
    EdgeRef     = 6

    @classmethod
    def fromKindId(cls, kindid):
        """Gets the entity kind from an identifier."""
        for kind in TGEntityKind:
            if kind.value == kindid:
                return kind
        return TGEntityKind.InvalidKind


class TGSystemObject(tgpdu.TGSerializable):
    """Represents a single object type in the graph database."""

    @property
    @abc.abstractmethod
    def systemType(self) -> TGSystemType:
        """Returns the kind of type of this object"""

    @property
    @abc.abstractmethod
    def name(self) -> str:
        """Returns the name of this object type"""

    @property
    @abc.abstractmethod
    def id(self):
        """Returns the identifier of this object type"""


class TGAttributeDescriptor(TGSystemObject):
    """Each object of this class represents an attribute descriptor in the database.

    An attribute descriptor in TIBCO Graph Database acts like a column - a combination of a name and a type. Unlike more
    traditional relational databases, this attribute descriptor can fit on different types and kinds of objects. For
    example one attribute descriptor can be used on both edges and nodes, and on different node types and edge types.
    Also, you can set an attribute on node or edge of a type, even when that type does not include the attribute
    descriptor.
    """

    @property
    @abc.abstractmethod
    def id(self):
        """Returns the identifier of this attribute descriptor"""

    @property
    @abc.abstractmethod
    def type(self) -> TGAttributeType:
        """Returns the attribute type of this object"""

    @property
    @abc.abstractmethod
    def isArray(self) -> bool:
        """Returns whether this attribute descriptor is an array"""

    @property
    @abc.abstractmethod
    def isEncrypted(self) -> bool:
        """Returns whether this attribute descriptor has encrypted values on the database side.

        NOTE: When receiving encrypted attributes from the server, they are only obfuscated, NOT ENCRYPTED. They were
        only obfuscated to reduce overhead on client's side. To ensure proper data protection from the server, be sure
        to use SSL channels instead of plain-text tcp channels. Sending encrypted data to the server is encrypted, even
        when going over an insecure channel.
        """

    @property
    @abc.abstractmethod
    def precision(self) -> int:
        """Returns the precision of this attribute descriptor (only if it is of Number type)"""

    @property
    @abc.abstractmethod
    def scale(self) -> int:
        """Returns the scale of this attribute descriptor."""


class TGAttribute(tgpdu.TGSerializable):
    """This represents an attribute on the database server.

    This is analogous to a cell in a relational model.
    """

    @property
    @abc.abstractmethod
    def descriptor(self) -> TGAttributeDescriptor:
        """Returns the attribute descriptor of this attribute (analogous to the attribute's type)."""

    @property
    @abc.abstractmethod
    def isNull(self) -> bool:
        """Returns whether this attribute has a null value (True if it's value is null)."""

    @property
    @abc.abstractmethod
    def isModified(self) -> bool:
        """Returns whether this attribute is modifiable."""

    @property
    @abc.abstractmethod
    def value(self) -> object:
        """Returns the value of this attribute"""

    @value.setter
    def value(self, value):
        """sets the value of this object"""
        self._setValue(value)

    '''
    Need to do this redirection because python getters/setters + the after thought of abstract methods = weird conflicts
    ala: https://stackoverflow.com/questions/35344209/python-abstract-property-setter-with-concrete-getter
    '''
    @abc.abstractmethod
    def _setValue(self, value):
        """Sets the value. Should NOT be called by client code."""

    @property
    def booleanValue(self) -> bool:
        """"Gets the value of this object as a boolean.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Boolean)

    @property
    def byteValue(self) -> bool:
        """"Gets the value of this object as a byte.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Byte)

    @property
    def charValue(self):
        """"Gets the value of this object as a character.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Char)


    @property
    def shortValue(self):
        """"Gets the value of this object as a 16-bit integer.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Short)

    @property
    def intValue(self):
        """"Gets the value of this object as a 32-bit integer.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Integer)

    @property
    def longValue(self):
        """"Gets the value of this object as a 64-bit integer.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Long)

    @property
    def floatValue(self):
        """"Gets the value of this object as a 32-bit float.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Float)

    @property
    def doubleValue(self):
        """"Gets the value of this object as a 64-bit float.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Double)

    @property
    def stringValue(self):
        """"Gets the value of this object as a string.

        Will throw an exception if that conversion does not make sense.
        """
        value: object = self.value
        return "None" if value is None else value

    @property
    def dateValue(self):
        """"Gets the value of this object as a date.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Date)

    @property
    def timeValue(self):
        """"Gets the value of this object as a time.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Time)

    @property
    def timestampValue(self):
        """"Gets the value of this object as a timestamp.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.TimeStamp)

    @property
    def numberValue(self):
        """"Gets the value of this object as a number.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Number)

    def asBytes(self):
        """"Gets the value of this object as a byte array.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Byte, True)

    def asChars(self):
        """"Gets the value of this object as a character array.

        Will throw an exception if that conversion does not make sense.
        """
        raise tgexception.TGTypeCoercionNotSupported(self.descriptor.type, TGAttributeType.Char, True)


class DirectionType(enum.Enum):
    """The direction for an edge."""

    UnDirected = 0
    Directed = 1
    BiDirectional = 2
    All = 3

    @classmethod
    def fromId(cls, id: int):
        """Gets the direction from the integer identifier"""
        for direction in DirectionType:
            if direction.value == id:
                return direction
        return DirectionType.UnDirected


class Direction(enum.Enum):
    """Direction enum relative to the nodes."""
    Inbound = 0
    Outbound = 1
    Any = 2


class TGEntityType(TGSystemObject):
    """Represents a type for an entity on the server."""
    @property
    @abc.abstractmethod
    def attributeDescriptors(self) -> typing.List[TGAttributeDescriptor]:
        """Gets the attribute descriptors for this type.

        These are only defined because the all of the entities of this type should conceptually share many of these
        attributes.
        """

    @abc.abstractmethod
    def getAttributeDescriptor(self, name) -> TGAttributeDescriptor:
        """Gets a particular attribute descriptor (indexed by that attribute descriptor's name)."""

    @property
    @abc.abstractmethod
    def numberEntities(self) -> int:
        """Gets the number of entities """


class TGNodeType(TGEntityType):
    """This represents a node type.

    Specifically for nodes and should not be used with edges or any other kind of entity.
    """

    @property
    @abc.abstractmethod
    def pkeyAttributeDescriptors(self) -> typing.List[TGAttributeDescriptor]:
        """Gets the primary key attribute descriptors."""


class TGEdgeType(TGEntityType):
    """This represents an edge type.

    Do not use with anything other than edge entities.
    """

    @property
    @abc.abstractmethod
    def directiontype(self) -> DirectionType:
        """Gets the direction of this edge type."""

    @property
    @abc.abstractmethod
    def fromNodeType(self) -> TGNodeType:
        """Gets the from node type (or none if one does not exist)."""

    @property
    @abc.abstractmethod
    def toNodeType(self) -> TGNodeType:
        """Gets the to node type (or none if one does not exist)."""


class TGIndex(TGSystemObject):
    """This represents an index definition and all of the configuration aspects for this index."""

    @property
    @abc.abstractmethod
    def unique(self) -> bool:
        """Gets whether this index is unique"""

    @property
    @abc.abstractmethod
    def attributeDescriptors(self) -> typing.List[TGAttributeDescriptor]:
        """Gets the attribute descriptors that form this index (in order)."""

    @property
    @abc.abstractmethod
    def onTypes(self) -> typing.List[TGEntityType]:
        """Gets the types that are currently in this index."""
    
    @property
    @abc.abstractmethod
    def blockSize(self) -> int:
        """Gets the size of the data block when the index is stored on disk."""

    @property
    @abc.abstractmethod
    def status(self) -> str:
        """Gets the status of the database"""

    @property
    @abc.abstractmethod
    def numEntries(self) -> int:
        """Gets the number of entities in the index."""


class TGSPReturnType(enum.Enum):
    Invalid = 0
    Int = 1
    String = 2
    Node = 4
    Edge = 5

    @classmethod
    def fromId(cls, id_: int):
        for sp in TGSPReturnType:
            if sp.value == id_:
                return sp
        return TGSPReturnType.Invalid


class TGSPReturnElemType(enum.Enum):
    SPRInvalid = 0
    SPRNode = 1
    SPREdge = 2
    SPRString = 3
    SPRChar = 4
    SPRByte = 5
    SPRuByte = 6
    SPRBool = 7
    SPRShort = 8
    SPRuShort = 9
    SPRInt = 10
    SPRuInt = 11
    SPRLong = 12
    SPRuLong = 13
    SPRFloat = 14
    SPRDouble = 15
    SPRNumber = 16
    SPRLLong = 17
    SPRuLLong = 18

    @classmethod
    def fromId(cls, id_: int):
        for sp in TGSPReturnElemType:
            if sp.value == id_:
                return sp
        return TGSPReturnElemType.SPRInvalid


"""
class TGSPReturnType(abc.ABC):
    isList: bool = False
    typeName: str = None

    @property
    @abc.abstractmethod
    def elemType(self) -> TGSPReturnElemType:
        #""TODO""
"""

class TGSPParamType(enum.Enum):
    Invalid = 0
    GraphContext = 1
    TransactionChangelist = 2
    Boolean = 3
    Long = 4
    String = 5
    Float = 6

    @classmethod
    def fromId(cls, id_: int):
        for sp in TGSPParamType:
            if sp.value == id_:
                return sp
        return TGSPParamType.Invalid


class TGSPParam(abc.ABC):
    
    @property
    @abc.abstractmethod
    def name(self) -> str:
        """TODO"""

    @property
    @abc.abstractmethod
    def paramType(self) -> TGSPParamType:
        """TODO"""


class TGStoredProcedure(TGSystemObject, abc.ABC):

    @property
    @abc.abstractmethod
    def module(self) -> str:
        """TODO"""

    @property
    @abc.abstractmethod
    def returnType(self) -> str:
    #def returnType(self) -> TGSPReturnElemType:
        """TODO"""

    @property
    @abc.abstractmethod
    def params(self) -> typing.List[TGSPParam]:
        """TODO"""


class TGPrivilege(enum.Enum):
    """
    This class represents the privileges that a role could possibly have access to. Privileges represent a database-wide
    entitlement granted to a role that allows database diagnostics, granting/revoking of privileges or permissions, and
    other entitlements that are not specific to any system type.
    """
    Grant = (0, 'g')
    Revoke = (1, 'r')
    Operate = (2, 'o')
    Diagnostic = (3, 'd')
    Import = (4, 'i')
    Export = (5, 'e')
    APIProxy = (6, 'p')

    def __init__(self, flag: int, shorthand: str):
        self.__flag = flag
        self.__sh = shorthand

    @property
    def flag(self) -> int:
        """Gets the bitfield index for this privilege."""
        return self.__flag

    @property
    def shorthand(self) -> str:
        """Gets the short-hand string name for this privilege."""
        return self.__sh

    @classmethod
    def privilegesToInt(cls, data: typing.Union[typing.Iterable, typing.Any]):
        """Convets one or many privileges into a set bit string."""
        ret = 0

        try:
            for priv in data:
                ret |= 1 << priv.flag
        except TypeError:
            ret = 1 << data.flag

        return ret

    @classmethod
    def flagToStr(cls, values: typing.FrozenSet) -> str:
        """Converts a set of privileges into a more-readable string"""
        ret: str = ''
        for val in cls:
            if val.flag >= 0:
                if val in values:
                    ret += val.shorthand
                else:
                    ret += ' '
        return ret

    @classmethod
    def flagFromStr(cls, val: str) -> int:
        """Converts a collection of bits in a string into a bit string representation of those privileges."""
        import re
        import tgdb.exception as tgexcept
        ret = 0
        val = val.lower()
        if re.match("[grodie]+", val) and (len(set(val)) == len(val)):
            for perm in TGPrivilege:
                if perm.shorthand in val:
                    ret = ret | (1 << perm.flag)
        else:
            raise tgexcept.TGException('Illegal permission string: %s' % val)
        return ret

    @classmethod
    def privilegesFromInt(cls, data: int):
        """Converts a bit string into a set of privileges."""
        ret = set()

        for priv in cls:
            if (data & (1 << priv.flag)) > 0:
                ret.add(priv)

        return ret


class TGPermissionType(enum.Enum):
    """This class represents a single role permission, or all of them."""
    Create = (1 << 0, 'c')
    Read = (1 << 1, 'r')
    Update = (1 << 2, 'u')
    Delete = (1 << 3, 'd')
    Encrypted = (1 << 4, 'e')
    Execute = (1 << 5, 'x')
    All = (0xFFFFFFFF, 'all')

    def __init__(self, flag: int, shorthand: str):
        self.__flag = flag
        self.__sh = shorthand

    def __hash__(self):
        """Intended only for using in a set."""
        return self.__flag

    def __eq__(self, other):
        """Intended only for using in a set."""
        if isinstance(other, TGPermissionType):
            return other.__flag == self.__flag
        return TypeError

    @property
    def flag(self) -> int:
        """Gets the bit field that this permission."""
        return self.__flag

    @property
    def shorthand(self) -> str:
        """Gets the shorthand name of this permission"""
        return self.__sh

    @classmethod
    def flagFromStr(cls, val: str) -> int:
        """Converts a flag from a string of shorthand characters into a bit-string of valid permission flags."""
        import re
        import tgdb.exception as tgexcept
        ret = 0
        val = val.lower()
        if val == TGPermissionType.All.shorthand:
            ret = TGPermissionType.All.flag
        elif re.match("[crduxe]+", val) and (len(set(val)) == len(val)):
            for perm in TGPermissionType:
                if perm.shorthand in val:
                    ret = ret | perm.flag
        else:
            raise tgexcept.TGException('Illegal permission string: %s' % val)
        return ret

    @classmethod
    def flagToStr(cls, values: typing.FrozenSet) -> str:
        """Converts a set of permission"""
        ret: str = ''
        for val in cls:
            if val.flag > 0 and val.flag & 3 < 3:
                if val in values:
                    ret += val.shorthand
                else:
                    ret += ' '
        return ret

    @classmethod
    def permissionsFromByte(cls, data: int):
        """Gets the permission type from the byte"""
        ret = set()

        for perm in cls:
            if (perm.__flag > 0) and \
                    ((data & perm.__flag) > 0):
                ret.add(perm)

        return ret


class TGRolePermFlag(enum.Enum):
    """This class determines if a role permission designates any role permission."""
    Invalid = -1
    NoFlags = 0
    SYSCATALOG = 1
    Default = 2

    def __hash__(self):
        """Intended only for using in a set."""
        return self.value

    def __eq__(self, other):
        """Intended only for using in a set."""
        if isinstance(other, TGRolePermFlag):
            return other.value == self.value
        return TypeError

    @classmethod
    def flagFromInt(cls, flag: int):
        """Intended for internal use only."""
        for rpf in cls:
            if rpf.value == flag:
                return rpf
        return cls.Invalid


class TGRolePermission(abc.ABC):
    """Acts as a role's permissions on a specific attribute descriptor, nodetype, or edgetype."""

    @property
    @abc.abstractmethod
    def grantedPerms(self) -> typing.FrozenSet[TGPermissionType]:
        """Gets the granted permission values for this role-permission."""

    @property
    @abc.abstractmethod
    def revokedPerms(self) -> typing.FrozenSet[TGPermissionType]:
        """Gets the revoked permission values for this role-permission."""

    @property
    @abc.abstractmethod
    def sysID(self):
        """Gets the system identifier that corresponds to an attribute descriptor, a nodetype, or an edgetype."""

    @property
    @abc.abstractmethod
    def flag(self) -> TGRolePermFlag:
        """Gets the flag for this role permission."""


class TGRole(TGSystemObject):
    """A role is a collection of permissions on specific system types."""
    
    @property
    @abc.abstractmethod
    def privileges(self) -> typing.FrozenSet[TGPrivilege]:
        """
        Gets the privileges. Privileges differ from permissions in that privileges determine database wide abilities
        like importing/exporting, granting/revoking permissions or privileges, etc.
        """

    @property
    @abc.abstractmethod
    def permissions(self) -> typing.List[TGRolePermission]:
        """Gets a list of all of the general (non-special) permissions."""

    @property
    @abc.abstractmethod
    def default(self) -> TGRolePermission:
        """Gets the default permissions for this role when no type-specific permission is set."""

    @property
    @abc.abstractmethod
    def syscatalog(self) -> TGRolePermission:
        """
        Gets the SYSCATALOG permissions for this role. These specify the permissions on creation/reading/updating/
        deleting/etc on system types.
        """


class TGGraphMetadata(tgpdu.TGSerializable):
    """An object representing all of the metadata for this server instance.

    This includes all of the node types, edge types, and attribute descriptors.
    """

    @property
    @abc.abstractmethod
    def nodetypes(self) -> typing.Dict[str, TGNodeType]:
        """Gets the node types as dictionary indexed by their name."""

    @property
    @abc.abstractmethod
    def edgetypes(self) -> typing.Dict[str, TGEdgeType]:
        """Gets the edge types as dictionary indexed by their name."""

    @property
    @abc.abstractmethod
    def attritubeDescriptors(self) -> typing.Dict[str, TGAttributeDescriptor]:
        """Gets the attribute descriptors as dictionary indexed by their name."""

    @property
    @abc.abstractmethod
    def types(self) -> typing.Dict[int, TGSystemObject]:
        """Gets the types as dictionary indexed by their identifier."""


class TGEntity(tgpdu.TGSerializable):
    """Represents a single entity on the server."""

    @property
    @abc.abstractmethod
    def entityKind(self) -> TGEntityKind:
        """Gets the kind of entity."""

    @property
    @abc.abstractmethod
    def entityType(self) -> TGEntityType:
        """Gets the type of this entity."""

    @property
    @abc.abstractmethod
    def isNew(self) -> bool:
        """
        Determines whether this entity is a newly created one on this client's session and has yet to be sent to the
        server.
        """

    @property
    @abc.abstractmethod
    def isDeleted(self) -> bool:
        """Is true if this entity was deleted from the server."""

    @property
    @abc.abstractmethod
    def version(self) -> int:
        """Gets the version of this entity."""

    @property
    @abc.abstractmethod
    def graphMetadata(self) -> TGGraphMetadata:
        """Gets the graphmetadata."""

    @property
    @abc.abstractmethod
    def attributes(self) -> typing.List[TGAttribute]:
        """Gets the attributes of entity."""

    @abc.abstractmethod
    def getAttribute(self, key: str) -> TGAttribute:
        """Gets a particular attribute of this entity."""

    @abc.abstractmethod
    def setAttribute(self, key: str, value: TGAttribute):
        """Sets a particular attribute of this entity."""

    @property
    @abc.abstractmethod
    def primaryKey(self) -> typing.Any:
        """
        Gets the primary key for this entity.

        Only guaranteed to compare equal to another entity when their primary keys match. Other than that, nothing else
        about the type or what it compares to is guaranteed.
        """

    @abc.abstractmethod
    def __setitem__(self, key: str, value: typing.Any):
        """Sets the value of the attribute.

        The key is the name of the attribute descriptor and the value is what to set the attribute for this entity.
        """

    @abc.abstractmethod
    def __getitem__(self, key: str) -> typing.Any:
        """Gets the value of the attribute indexed by the name of the attribute descriptor."""

    @abc.abstractmethod
    def __hash__(self) -> int:
        """Gets the hash value for this entity."""

    @abc.abstractmethod
    def __eq__(self, other):
        """Gets whether this entity equals the other entity."""


class TGNode(TGEntity):
    """Represents a single node in the server (or on the client)."""

    # get Edges and pass a filter function
    @abc.abstractmethod
    def getEdges(self, filter=None) -> typing.List:
        """Gets the edges associated with this node.

        :param filter: is a filter to base which edges are included.
        """

    @property
    @abc.abstractmethod
    def primaryKey(self) -> typing.Any:
        """
        Gets the primary key of this node. Only guarantee is that it will always compare true to a node with the
        same key, and false to a node with a different key.
        """

    # this property is same as the getEdges()
    @property
    @abc.abstractmethod
    def edges(self) -> typing.List:
        """Gets all of the edges with an endpoint on this node."""

    @property
    @abc.abstractmethod
    def directedEdges(self):
        """Gets all of the directed edges with an endpoint on this node."""

    @property
    @abc.abstractmethod
    def undirectedEdges(self):
        """Gets all of the undirected edges with an endpoint on this node."""

    @property
    @abc.abstractmethod
    def bidirectedEdges(self):
        """Gets all of the bidirected edges with an endpoint on this node."""


class TGEdge(TGEntity):
    """Represents a single edge."""

    @property
    @abc.abstractmethod
    def vertices(self) -> typing.Tuple[TGNode, TGNode]:
        """Gets the vertex endpoints of this edge.

        The first node is the from vertex, the second is the to vertex.
        """

    @property
    @abc.abstractmethod
    def fromVertex(self) -> TGNode:
        """
        Gets the from vertex of this edge, or the vertex that was listed first when adding to the database if this is an
        undirected or bidirected edge.
        """

    @property
    @abc.abstractmethod
    def toVertex(self) -> TGNode:
        """
        Gets the to vertex of this edge, or the vertex that was listed second when adding to the database if this is an
        undirected or bidirected edge.
        """

    @property
    @abc.abstractmethod
    def directiontype(self) -> DirectionType:
        """Gets the direction of this edge."""


class TGKey(tgpdu.TGSerializable):
    """Represents a key into the database.

    Keys must be only on uniquely identifiable attributes with an explicit index or one implicitly created for
    attributes that are primary keys.
    """

    @abc.abstractmethod
    def __setitem__(self, key, value):
        """Sets the attribute values for this key.

        :param key: The attribute descriptor to set.
        :param value: The value to set the attribute's value to.
        """

    @abc.abstractmethod
    def matches(self, entity: TGEntity) -> bool:
        """Determine whether this key matches with the passed entity."""


class TGGraphObjectFactory(abc.ABC):
    """Is a factory for nodes, edges, and keys for the client.

    Important to use this instead of directly instantiating the objects on your own.
    """

    @property
    @abc.abstractmethod
    def connection(self):
        """Gets the connection GOF."""

    @property
    @abc.abstractmethod
    def graphmetadata(self) -> TGGraphMetadata:
        """Gets graph metadata of this GOF."""

    @abc.abstractmethod
    def createNode(self, nodetype: typing.Union[str, TGNodeType]) -> TGNode:
        """Create a node.

        :param nodetype: The nodetype of the node to create. Either the name of the nodetype (as a string) or the
        nodetype itself.
        """

    @abc.abstractmethod
    def createEdge(self, fromnode: TGNode, tonode: TGNode, dirtype: DirectionType = DirectionType.Directed,
                   edgetype: typing.Union[str, TGEdgeType] = None) -> TGEdge:
        """Create an edge.

        :param fromnode: The from node for this edge.
        :param tonode: The to node for this edge.
        :param dirtype: The direction of this node (default: Directed).
        :param edgetype: The edgetype of the node to create. Either the name of the edgetype (as a string) or the
        edgetype itself. Will overright the dirtype parameter if set to a non-null value (default: None).
        """

    @abc.abstractmethod
    def createCompositeKey(self, nodetype: typing.Union[str, TGNodeType]) -> TGKey:
        """Create a key.

        :param nodetype: The nodetype of the node to search for. Either the name of the nodetype (as a string) or the
        nodetype itself.
        """
