










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
* File name : version.tag
* Created on: 11/01/2019
* Created by: Derek
* SVN Id: $Id$
*
* Determines version information.
*
* This file is processed by GCC-Processor using -E and options.
* Use C/C++ Style comments, both multiline and single line
* This file is passed through C-Preprocesor, and Python comments are interpreted as C-Processor commands.
"""


import enum
import typing
import struct

@enum.unique
class TGBuildEdition(enum.Enum):
    """
    Represents the type of edition which determines what the usage for this API and server are.
    """
    Evaluation = 0
    Community = 1
    Enterprise = 2
    Developer = 3

    @classmethod
    def fromId(cls, id: int):
        for v in TGBuildEdition:
            if v.value == id:
                return v
        return TGBuildEdition.Developer

    @classmethod
    def fromName(cls, id: str):
        for name, member in TGBuildEdition.__members__.items():
            if name == id:
                return member
        return TGBuildEdition.Developer

@enum.unique
class TGBuildType(enum.Enum):
    """
    The type of build. This determines how feature complete and thoroughly tested this particular release.
    """
    Production = 0
    Engineering = 1
    Beta = 2
    Alpha = 3

    @classmethod
    def fromId(cls, id: int):
        for v in TGBuildType:
            if v.value == id:
                return v
        return TGBuildType.Engineering

    @classmethod
    def fromName(cls, nm: str):
        for name, member in TGBuildType.__members__.items():
            if name == nm:
                return member
        return TGBuildType.Engineering


class TGVersion:
    """
    Represents the over-all version of this API and server. This determines whether the server and client can    communicate and what features are supported and how thoroughly tested this release is.
    """
    _version = None
    _MAJOR = 3
    _MINOR = 0
    _UPDATE = 0
    _HFNO = 0
    _BUILD_NO = 39
    _BUILD_REV = 4709
    _BUILD_TYPE = TGBuildType.fromName("Production")
    _EDITION = TGBuildEdition.fromName("Enterprise")

    def __init__(self, edition: TGBuildEdition = None, typ: TGBuildType = None):
        self.__major = TGVersion._MAJOR
        self.__minor = TGVersion._MINOR
        self.__update = TGVersion._UPDATE
        self.__hfNo = TGVersion._HFNO
        self.__buildNo = TGVersion._BUILD_NO
        self.__buildRev = TGVersion._BUILD_REV

        self.__edition = edition if edition is not None else TGVersion._EDITION
        self.__type = typ if typ is not None else TGVersion._BUILD_TYPE

    @property
    def major(self) -> int:
        """
        This is the major version number. This indicates what major features are introduced and that compatability is        not necessarily guaranteed.

        :return: The major version number.
        :rtype: int
        """
        return self.__major

    @property
    def minor(self) -> int:
        """
        This is the minor version number. This indicates what minor features are introduced.

        :return: The minor version number.
        :rtype: int
        """
        return self.__minor

    @property
    def update(self) -> int:
        """
        This is the update version number. This indicates what incremental changes are introduced.

        :return: The update version number.
        :rtype: int
        """
        return self.__update

    @property
    def hfNo(self) -> int:
        return self.__hfNo

    @property
    def buildNo(self) -> int:
        """
        The build number for this release.

        :return: The build number.
        :rtype: int
        """
        return self.__buildNo

    @property
    def buildRev(self) -> int:
        """
        The content repository revision for this build.

        :return: The repository revision.
        :rtype: int
        """
        return self.__buildRev

    @property
    def edition(self) -> TGBuildEdition:
        """
        The edition for this build.

        :return: The build edition.
        :rtype: TGBuildEdition
        """
        return self.__edition

    @property
    def buildType(self) -> TGBuildType:
        """
        The type of this build.

        :return: The build edition.
        :rtype: TGBuildType
        """
        return self.__type

    @classmethod
    def getVersion(cls, edition: typing.Optional[TGBuildEdition] = None, typ: typing.Optional[TGBuildType] = None):
        """
        Gets the current version.

        :param edition: The edition to set for the returned version.
        :type edition: TGBuildEdition
        :param typ: The type to set for the returned version.
        :type typ: TGBuildType
        :return: The version.
        :rtype: TGVersion
        """
        if edition is not None and typ is not None:
            cls._version = TGVersion(edition=edition, typ=typ)
        elif edition is not None:
            cls._version = TGVersion(edition=edition)
        elif typ is not None:
            cls._version = TGVersion(typ=typ)
        elif cls._version is None:
            cls._version = TGVersion()

        return cls._version

    def __str__(self):
        return "%d.%d.%d Build(%d) Revision(%d) %s Edition." % (self.__major, self.__minor, self.__update,
                            self.__buildNo, self.__buildRev, self.__edition.name)

    def getVersionStr(self):
        """
        Gets an acceptable build version that includes the major, minor, update, and build number.

        :returns: A string with the build version.
        :rtype: str
        """
        return "%d.%d.%d Build(%d)" % (self.__major, self.__minor, self.__update, self.__buildNo)

    def __eq__(self, other) -> bool:
        ret = False
        if isinstance(other, TGVersion):
            other: TGVersion = other
            ret = True
            ret = ret and (self.__major == other.__major)
            ret = ret and (self.__minor == other.__minor)
            ret = ret and (self.__update == other.__update)
            ret = ret and (self.__hfNo == other.__hfNo)
            ret = ret and (self.__buildNo == other.__buildNo)
            ret = ret and (self.__type == other.__type)
            ret = ret and (self.__edition == other.__edition)
            ret = ret and (self.__buildRev == other.__buildRev)

        return ret

    def writeExternal(self, ostream):
        """
        Only intended for use within the Python API.

        :param ostream: The out stream.
        :return: Nothing
        """
        import tgdb.pdu as tgpdu

        ostream: tgpdu.TGOutputStream = ostream

        ostream.writeByte(self.__major)
        ostream.writeByte(self.__minor)
        ostream.writeByte(self.__update)
        ostream.writeByte(self.__hfNo)
        ostream.writeShort(self.__buildNo)
        ostream.writeByte(((0xF & self.__type.value) << 4) | (0xF & self.__edition.value))
        ostream.writeByte(0)

    def toUInt64(self) -> int:
        """
        Converts this version to a bit string.

        :return: The bit string representation of this version.
        :rtype: int
        """
        ret: int = 0
        ret |= (0xFF & self.__major)
        ret = ret << 8
        ret |= (0xFF & self.__minor)
        ret = ret << 8
        ret |= (0xFF & self.__update)
        ret = ret << 8
        ret |= (0xFF & self.__hfNo)
        ret = ret << 8
        ret |= (0xFFFF & self.__buildNo)
        ret = ret << 16
        ret |= ((0xF & self.__type.value) << 4) | (0xF & self.__edition.value)
        ret = ret << 8
        buf = struct.pack(">Q", ret)

        return struct.unpack_from("<Q", buf)[0]



    @classmethod
    def readExternal(cls, instream, reverse: bool = False):
        """
        Only intended for use within the Python API.

        :param instream: The in stream.
        :param reverse: Whether to reverse the version bit-string.
        :return: Nothing
        """
        import tgdb.pdu as tgpdu

        instream: tgpdu.TGInputStream = instream

        version = TGVersion()
        if reverse:
            _ = instream.readByte()

            typeAndEdition = instream.readByte()

            typeNum = typeAndEdition & 0xF

            editionNum = (typeAndEdition >> 4) & 0xF

            version.__type = TGBuildType.fromId(typeNum)
            version.__edition = TGBuildEdition.fromId(editionNum)

            version.__buildNo = instream.readShort()
            version.__hfNo = instream.readByte()
            version.__update = instream.readByte()
            version.__minor = instream.readByte()
            version.__major = instream.readByte()

        else:
            version.__major = instream.readByte()
            version.__minor = instream.readByte()
            version.__update = instream.readByte()
            version.__hfNo = instream.readByte()
            version.__buildNo = instream.readShort()

            typeAndEdition = instream.readByte()

            typeNum = (typeAndEdition >> 4) & 0xF

            editionNum = typeAndEdition & 0xF

            version.__type = TGBuildType.fromId(typeNum)
            version.__edition = TGBuildEdition.fromId(editionNum)

            _ = instream.readByte()

        return version

