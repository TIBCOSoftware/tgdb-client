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
 *  File name :channel.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: utils.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates the utils package of java
 """

import enum
import typing
import struct


def tg_checkSum64(data: typing.Union[bytes, bytearray], seed: int = 0):
    """This function creates a TGDB standard checksum.

    This checksum is based on the MurMur64 hashing algorithm.

    :param data: The data bytes to find the hash of.
    :type data: bytes-like
    :param seed: a 64-bit integer that represents a fixed seed to base the checksum off of (default is 0).
    :type seed: int

    :returns: a 64-bit integer representing the MurMur64 hash of the input data and the seed.
    :rtype: int
    """

    def __limit(input: int):
        return input & 0xFFFFFFFFFFFFFFFF

    m = 0xc6a4a7935bd1e995
    r = 47

    length = len(data)

    hash = __limit(seed) ^ __limit(length*m)

    start = 0
    end = length - (length % 8)

    while start != end:
        cur_data: int = struct.unpack_from("<Q", data, start)[0]

        cur_data = __limit(cur_data * m)
        cur_data = __limit(cur_data ^ (cur_data >> r))
        cur_data = __limit(cur_data * m)

        hash = __limit(cur_data ^ hash)
        hash = __limit(hash * m)

        start += 8

    rem = length & 0b111

    # Quite a bit less tidy than the switch statement in the C version, but Python is particular about some things...
    # Maybe move to same line to same line to reflect C version?
    if rem > 6:
        hash ^= (data[end + 6] << 48)
    if rem > 5:
        hash ^= (data[end + 5] << 40)
    if rem > 4:
        hash ^= (data[end + 4] << 32)
    if rem > 3:
        hash ^= (data[end + 3] << 24)
    if rem > 2:
        hash ^= (data[end + 2] << 16)
    if rem > 1:
        hash ^= (data[end + 1] << 8)
    if rem > 0:
        hash ^= (data[end + 0] << 0)
        hash = __limit(hash * m)

    hash = __limit(hash)

    hash ^= hash >> r
    hash = __limit(hash * m)
    hash ^= hash >> r

    return __limit(hash)


class TGProtocolVersion:
    """ TGProtocolVersion keeps track of protocol constants."""

    __majorver = 3
    __minorver = 0
    __magic = 0xdb2d1e4

    @classmethod
    def getProtocolVersion(cls) -> int:
        """
        Gets the protocol version as a single integer

        :returns: The protocol major version and minor version.
        :rtype: int
        """
        ver = cls.__majorver << 8 + cls.__minorver
        return ver

    @classmethod
    def magic(cls):
        """
        Gets the magic number

        :returns: The magic number.
        :rtype: int
        """
        return cls.__magic

    @classmethod
    def isCompatible(cls, protver: int) -> bool:
        """
        Checks if this version is compatible with the server's version.

        :returns: Whether this API's protocol version is compatible with the server's protocol version.
        :rtype: bool
        """
        return cls.getProtocolVersion() == protver


class ConfigName(enum.Enum):
    """ Stores Configuration Related Keys, Default Values, and Descriptions"""

    BulkIONodeBatchSize = (
        "tgdb.bulkIO.nodeBatchSize",
        "bulkIONodeBatchSize",
        "1000",
        "The maximum number of nodes to send at once in a bulk import"
    )

    BulkIOEntityBatchSize = (
        "tgdb.bulkIO.entityBatchSize",
        "bulkIOEntityBatchSize",
        "1000",
        "The maximum number of nodes to send at once in a bulk import"
    )

    BulkIOEdgeBatchSize = (
        "tgdb.bulkIO.edgeBatchSize",
        "bulkIOEdgeBatchSize",
        "2500",
        "The maximum number of edges to send at once in a bulk import"
    )
    
    BulkIOIdColName = (
        "tgdb.bulkIO.idColName",
        "bulkIOIdColName",
        "id",
        "The column name to get the identifiers for the nodes imported"
    )

    BulkIOFromColName = (
        "tgdb.bulkIO.fromColName",
        "bulkIOFromColName",
        "from",
        "The column name to get the 'from' node identifiers for the edges imported"
    )

    BulkIOToColName = (
        "tgdb.bulkIO.toColName",
        "bulkIOToColName",
        "to",
        "The column name to get the 'to' node identifiers for the nodes imported"
    )

    ChannelDefaultHost = (
        "tgdb.channel.defaultHost",
        "defaultHost",
        "localhost",
        "The default host specifier"
    )

    ChannelDefaultPort = (
        "tgdb.channel.defaultPort",
        "defaultPort",
        "8222",
        "The default port specifier"
    )

    ChannelDefaultProtocol = (
        "tgdb.channel.defaultProtocol",
        "defaultProtocol",
        "ssl",
        "The default protocol"
    )

    ChannelSendSize = (
        "tgdb.channel.sendSize",
        "sendSize",
        "122",
        "TCP send packet size in KBs"
    )

    ChannelRecvSize = (
        "tgdb.channel.recvSize",
        "recvSize",
        "128",
        "TCP recv packet size in KB"
    )

    ChannelPingInterval = (
        "tgdb.channel.pingInterval",
        "pingInterval",
        "30",
        "Keep alive ping intervals"
    )

    ChannelConnectTimeout = ( # 5 sec timeout.
        "tgdb.channel.connectTimeout",
        "connectTimeout",
        "5.000",
        "Timeout for connection to establish, before it gives up and tries the ftUrls if specified"
    )

    ChannelFTHosts = (
        "tgdb.channel.ftHosts",
        "ftHosts",
        None,
        "Alternate fault tolerant list of &lt;host:port&gt; pair seperated by comma. "
    )

    ChannelFTRetryIntervalSeconds = (
        "tgdb.channel.ftRetryIntervalSeconds",
        "ftRetryIntervalSeconds",
        "10",
        "The connect retry interval to ftHosts"
    )

    ChannelFTRetryCount = (
        "tgdb.channel.ftRetryCount",
        "ftRetryCount",
        "3",
        "The number of times ro retry "
    )

    ChannelDefaultUserID = (
        "tgdb.channel.defaultUserID",
        "defaultUserID",
        None,
        "The default user Id for the connection"
    )

    ChannelProtocol = (
        "tgdb.channel.protocol",
        "protocol",
        None,
        "The protocol for the channel"
    )

    ChannelHost = (
        "tgdb.channel.host",
        "host",
        None,
        "The host for the channel"
    )

    ChannelPort = (
        "tgdb.channel.port",
        "port",
        None,
        "The port for the channel"
    )

    ChannelUserID = (
        "tgdb.channel.userID",
        "userID",
        None,
        "The user id for the connection if it is not specified in the API. See the rules for picking the user name"
    )

    ChannelPassword = (
        "tgdb.channel.password",
        "password",
        None,
        "The password for the username"
    )

    ChannelClientId = (
        "tgdb.channel.clientId",
        "clientId",
        "tgdb.python-api.client",
        "The client id to be used for the connection"
    )

    ConnectionDatabaseName = (
        "tgdb.connection.dbName",
        "dbName",
        None,
        "The database name the client is connecting to. It is used as part of verification for ssl channels"
    )

    ConnectionSpecifiedRoles = (
        "tgdb.connection.specifiedRoles",
        "roles",
        None,
        "The role name(s) that the user wants to log in as."
    )

    ConnectionPoolUseDedicatedChannelPerConnection = (
        "tgdb.connectionpool.useDedicatedChannelPerConnection",
        "useDedicatedChannelPerConnection",
        "false",
        "A boolean value indicating either to multiplex mulitple connections on a single tcp socket or use dedicate " +
        "socket per connection. A true value consumes resource but provides good performance. Also check " +
        "the max number of connections"
    )

    ConnectionPoolDefaultPoolSize = (
        "tgdb.connectionpool.defaultPoolSize",
        "defaultPoolSize",
        "10",
        "The default connection pool size to use when creating a ConnectionPool"
    )

    # 0 = mean immediate, Integer Max for indefinite.
    ConnectionReserveTimeoutSeconds = (
        "tgdb.connectionpool.connectionReserveTimeoutSeconds",
        "connectionpoolReserveTimeoutSeconds",
        "10",
        "A timeout parameter indicating how long to wait before getting a connection from the pool"
    )

    # Default Value is None which corresponds to no timeout
    ConnectionOperationTimeoutSeconds = (
        "tgdb.connection.operationTimeoutSeconds",
        "connectionOperationTimeoutSeconds",
        "3600",
        "A timeout parameter indicating how long to wait for a operation before giving up. "
        "Some queries are long running, and may override this behavior."
    )

    ConnectionIdleTimeoutSeconds = (
        "tgdb.connection.idleTimeoutSeconds",
        "connectionIdleTimeoutSeconds",
        None,
        "An idle teimeout parameter requested to server, before the server disconnects. This may or may not be honored "
        "by the server."
    )

    ConnectionDateFormat = (
        "tgdb.connection.dateFormat",
        "dateFormat",
        "YYYY-MM-DD",
        "Date format for this connection"
    )

    ConnectionTimeFormat = (
        "tgdb.connection.timeFormat",
        "timeFormat",
        "HH:mm:ss",
        "Date format for this connection"
    )
    ConnectionTimeStampFormat = (
        "tgdb.connection.timeStampFormat",
        "timeStampFormat",
        "YYYY-MM-DD HH:mm:ss.zzz",
        "Timestamp format for this connection"
    )

    ConnectionLocale = (
        "tgdb.connection.locale",
        "locale",
        "en_US",
        "Locale for this connection"
    )

    ConnectionDefaultQueryLanguage = (
        "tgdb.connection.defaultQueryLanguage",
        "queryLanguage",
        "gremlin",
        "Default query language format for this connection"
    )

    TlsLibraryLocation = (
        "tgdb.tls.library.location",
        "tlsLibraryLocation",
        None,
        "Transport level security library location. Should be relative or absolute path to the SSL library."
    )

    # TSL Parameters
    TlsProviderName = (
        "tgdb.tls.provider.name",
        "tlsProviderName",
        "SunJSSE", # The default is the Sun JSSE.
        "Transport level Security provider. Work with your InfoSec team to change this value"
    )

    TlsProviderClassName = (
        "tgdb.tls.provider.className",
        "tlsProviderClassName",
        "com.sun.net.ssl.internal.ssl.Provider", # The default is the Sun JSSE.
        "The underlying Provider implementation. Work with your InfoSec team to change this value."
    )

    TlsProviderConfigFile = (
        "tgdb.tls.provider.configFile",
        "tlsProviderConfigFile",
        None,
        "Some providers require extra configuration parameters, and it can be passed as a file"
    )

    TlsProtocol = (
        "tgdb.tls.protocol",
        "tlsProtocol",
        "TLSv1.2",
        "tlsProtocol version. The system only supports 1.2+"
    )

    TlsCipherSuites = ( # Use the Default Cipher Suites
        "tgdb.tls.cipherSuites",
        "cipherSuites",
        None,
        "A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher " +
        "list and Openssl list that supports 1.2 protocol "
    )

    TlsVerifyDatabaseName = (
        "tgdb.tls.verifyDBName",
        "verifyDBName",
        "false",
        "Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL."
    )

    TlsExpectedHostName = (
        "tgdb.tls.expectedHostName",
        "expectedHostName",
        None,
        "The expected hostName for the certificate. This is for future use"
    )

    TlsTrustedCertificates = (
        "tgdb.tls.trustedCertificates",
        "trustedCertificates",
        None,
        "The list of trusted Certificates"
    )

    KeyStorePassword = (
        "tgdb.security.keyStorePassword",
        "keyStorePassword",
        None,
        "The Keystore for the password"
    )

    EnableConnectionTrace = (
        "tgdb.connection.enableTrace",
        "enableTrace",
        "false",
        "The flag for debugging purpose, to enable the commit trace"
    )

    ConnectionTraceDir = (
        "tgdb.connection.enableTraceDir",
        "enableTraceDir",
        ".",
        "The base directory to hold commit trace log"
    )

    InvalidName = (None, None, None, None)

    def __init__(self, pn, an, dv, desc):
        self.__pn__ = pn
        self.__an__ = an
        self.__dv__ = dv
        self.__desc__ = desc

    @property
    def propertyname(self) -> str:
        """
        The full property name. Should be reverse-domain-notation style.

        :return: The full property name.
        :rtype: str
        """
        return self.__pn__

    @property
    def aliasname(self) -> str:
        """
        The aliasname. This should be a shortened version of the \
        :meth:`propertyname<tgdb.utils.ConfigName.propertyname>` attribute and should be unique across all \
        :meth:`aliasname<tgdb.utils.ConfigName.aliasname>`.

        :return: The shortened, unique name for this configuration property.
        :rtype: str
        """
        return self.__an__

    @property
    def defaultvalue(self) -> typing.Optional[str]:
        """
        The default value for this configuration property.

        :return: The default value.
        :rtype: typing.Optional[str]
        """
        return self.__dv__

    @property
    def description(self) -> str:
        """
        A description of this configuration property.

        :return: The description.
        :rtype: str
        """
        return self.__desc__

    @classmethod
    def fromName(cls, name: str):
        """
        Converts a name into a :class:`ConfigName<tgdb.utils.ConfigName>` object.

        :param name: The name to use as a lookup.
        :type name: str
        :return: The :class:`ConfigName<tgdb.utils.ConfigName>` object that matches the name.
        :rtype: tgdb.utils.ConfigName
        """
        for cn in ConfigName:
            if cn == cls.InvalidName:
                continue
            if name.casefold() == cn.propertyname.casefold():
                return cn
            if name.casefold() == cn.aliasname.casefold():
                return cn
        return ConfigName.InvalidName

    @classmethod
    def asMap(cls) -> typing.Dict[str, str]:
        cnmap: typing.Dict[str, str] = dict()
        for cn in ConfigName:
            cnmap[cn.propertyname] = cn.defaultvalue
            cnmap[cn.aliasname] = cn.defaultvalue
        return cnmap


class TGProperties(typing.MutableMapping):
    """ Intended for accessing configuration data."""

    def __init__(self, map: typing.Optional[typing.Union[typing.Dict, typing.MutableMapping]] = None):
        """ Initializes a new property object

        :param map: Is a map to base this configuration property based off of.
        :type map: Optional[Union[Dict, TGProperties]]
        """
        self.__map: typing.Dict[str, typing.Any]
        if map is None:
            self.__map = {}
        elif isinstance(map, dict):
            self.__map = dict(map)
        elif isinstance(map, TGProperties):
            self.__map = dict(map.__map)
        else:
            raise TypeError('Bad argument type. Expect dict or TGProperties')

    def __setitem__(self, k: typing.Union[str, ConfigName], v: typing.Any) -> None:
        """ Sets a particular property in this TGProperties instance

        :param k: The key for this property
        :type k: Either a string or a ConfigName
        :param v: The value to set for this property
        :type v: Anything
        :returns: Nothing
        :rtype: None
        """

        if isinstance(k, ConfigName):
            self.__map[k.propertyname] = v
            self.__map[k.aliasname] = v
        elif isinstance(k, str):
            self.__map[k] = v
        else:
            raise TypeError('Bad argument type for key. Expect str or ConfigName')

    def __delitem__(self, v) -> None:
        raise AttributeError('Readonly properties - cannot delete value')

    def __getitem__(self, k: typing.Union[str, ConfigName]) -> typing.Optional[str]:
        """ Gets a particular property from this TGProperties instance

        :param k: The key for this property
        :type k: typing.Union[str, tgdb.utils.ConfigName]
        :returns: the value for that property, or None if that property doesn't exist.
        :rtype: typing.Optional[str]
        """
        if isinstance(k, ConfigName):
            pv = self.__map.get(k.propertyname, None)
            av = self.__map.get(k.aliasname, None)
            return pv if pv is not None else av
        elif isinstance(k, str):
            return self.__map[k]
        else:
            raise TypeError('Bad argument type for key. Expect str or ConfigName, not {0} which is of type {1}'\
                            .format(str(k), str(type(k))))

    def __len__(self) -> int:
        """
        Gets the number of set elements in this properties.

        :return: The number of set element.
        :rtype: int
        """
        return len(self.__map)

    def __iter__(self):
        return iter(self.__map)

    def get(self, k: typing.Union[str, ConfigName], v: typing.Any = None):
        """
        Gets a particular property from this TGProperties instance, with a fallback if it doesn't exist

        :param k: The key for this property
        :type k: typing.Union[str, ConfigName]
        :param v: A value for this property if one doesn't already exist
        :type v: typing.Any
        :returns: the value for that property, or None if that property doesn't exist.
        :rtype: typing.Any
        """
        try:
            ret = self[k]
            if v is None and isinstance(k, ConfigName):
                v = k.defaultvalue
            if ret is None:
                ret = v
            return ret
        except KeyError:
            return v

    def getDef(self, k: ConfigName):
        """
        Gets the default value of k if not set in the properties.

        :param k: The property configuration file.
        :type: ConfigName
        :return: The configuration value for that property, or that configuration name's default value.
        :rtype: typing.Optional[str]
        """
        return self.get(k, k.defaultvalue)




NullString = '0000'
Space = ' '
NewLine = '\r\n'


class HexUtils:
    """A class for all of the hex utilities used by this API."""

    @classmethod
    def formatHex(cls, buf: typing.Union[bytes, bytearray], buflen: int = 0, ll: int = 48, cl: int = 2,
                  prettyprint: bool = True) -> str:
        """ Prints bytes as a hex string

        :param buf: bytes to print as hex string
        :type buf: typing.Union[bytes, bytearray]
        :param buflen: how many bytes of the buffer to print, if less than 1, then print the whole thing (default: 0).
        :type buflen: int
        :param ll: how many bytes to print on a line (default: 48). Only used when pretty printing.
        :type ll: int
        :param cl: how many bytes to print before inserting a space (default: 2). Only used when pretty printing.
        :type cl: int
        :param prettyprint: whether to print the format in a more human-readable format (if True), or something machines
                    can more easily decipher (if False) (default: True).
        :type prettyprint: bool
        :returns: The formatted hexadecimal of the bytes-like.
        :rtype: str
        """
        if buf is None or len(buf) == 0:
            return NullString

        bnewline = False
        lineno = 1
        blen = len(buf) if buflen <= 0 else buflen
        buflist = []
        if prettyprint:
            buflist.append('Formatted Byte Array:{0}'.format(NewLine))
            buflist.append('{0:08x}{1}'.format(0, Space))
            for i in range(0, blen):
                if bnewline:
                    bnewline = False
                    buflist.append('{0}{1:08x}{2}'.format(NewLine, (lineno * ll), Space))
                buflist.append('{0:02x}'.format(buf[i]))

                if (i+1) % cl == 0:
                    buflist.append(Space)

                if (i+1) % ll == 0:
                    bnewline = True
                    lineno = lineno + 1
        else:
            for i in range(0, blen):
                buflist.append('{0:02x}'.format(buf[i]))

        return ''.join(buflist)


class Class:
    """The class-based utilities required by this API."""

    @classmethod
    def forName(cls, clsname: str):
        """
        Find the object by name or dotted path, importing as necessary.

        :param clsname: name of class to discover.
        :type clsname: str
        :returns: The module, class, or object needed.
        """
        # clsname = 'tgdbapi.{0}'.format(clsname)
        import importlib

        parts = clsname.split('.')
        module, n = None, 0
        while n < len(parts):
            modname = '.'.join(parts[:n+1])
            try:
                module = importlib.import_module(modname)
                n = n+1
            except ModuleNotFoundError:
                break

        object = module if module is not None else __builtins__

        for part in parts[n:]:
            try:
                object = getattr(object, part)
            except AttributeError:
                return None
        return object

