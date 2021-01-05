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
 *		SVN Id: $Id: channel.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates channel interfaces
 """

import enum
from typing import *
import abc
from tgdb.utils import *
import tgdb.pdu as tgpdu
import tgdb.exception as tgexception

class Status(enum.Enum):
    """Connection status to the server."""
    Waiting = 0
    Ok = 1
    Pushed = 2
    Resend = 3
    Disconnected = 4,
    Closed = 5


class LinkState(enum.Enum):
    """Socket link state."""
    NotConnected = 0
    Connected = 1
    Closing = 2
    Closed = 3
    FailedOnSend = 4
    FailedOnRecv = 5
    FailedOnProcessing = 6
    Reconnecting = 7
    Terminated = 8

    @classmethod
    def fromId(cls, id):
        for ls in LinkState:
            if ls.value == id:
                return ls
        return LinkState.Terminated


class ResendMode(enum.Enum):
    """The resend mode when disconnected."""
    DontReconnectAndIgnore = 0
    ReconnectAndResend = 1
    ReconnectAndRaiseException = 2
    ReconnectAndIgnore = 3


class ProtocolType(enum.Enum):
    """Types of protocol."""
    Tcp = 0
    Ssl = 1


class TGChannelResponseWaiter(abc.ABC):
    """Keeps track of a single request/response from the server."""

    class WaiterStatus(enum.Enum):
        Waiting = 1
        Ok = 2
        Pushed = 3
        Resend = 4
        Disconnected = 5
        Closed = 6

    @property
    @abc.abstractmethod
    def isBlocking(self):
        """Is this waiter blocking?

        :returns: True if this waiter blocks the thread, else False.
        """

    @property
    @abc.abstractmethod
    def status(self) -> WaiterStatus:
        """Gets the status of this ChannelResponseWaiter"""

    @property
    @abc.abstractmethod
    def requestid(self):
        """The request id corresponding to the response that this waiter is for."""

    @abc.abstractmethod
    def awaiting(self, status: WaiterStatus):
        """Thread-safe checking of a status."""

    @property
    @abc.abstractmethod
    def reply(self):
        """Gets the thread-safe reply"""

    @reply.setter
    def reply(self, msg: tgpdu.TGMessage):
        """Thread safe reply setter"""


class TGChannelStopMethod(enum.Enum):
    """Way that the channel was stopped."""

    Graceful = 0
    ClientForceful = 1
    RemoteKill = 2


class ExceptionHandleResult(enum.Enum):
    """How to handle exceptions."""

    RethrowException = 0
    RetryOperation = 1
    Disconnected = 2


class TGSocket(abc.ABC):
    """Represents a single lower-level socket."""

    @property
    @abc.abstractmethod
    def handle(self):
        """Gets the lower-level socket."""

    @property
    @abc.abstractmethod
    def inboxaddr(self):
        """Gets the client-side address."""

    @abc.abstractmethod
    def connect(self):
        """Connect to the server."""

    @abc.abstractmethod
    def close(self):
        """Disconnect from the server."""

    @abc.abstractmethod
    def send(self, msg: tgpdu.TGMessage) -> int:
        """Send a message, returning the number of bytes sent."""

    @abc.abstractmethod
    def recvMsg(self) -> tgpdu.TGMessage:
        """Reads a message from."""

    def tryRead(self) -> tgpdu.TGMessage:
        """Try to read a message.

        :returns: None if no message was available, otherwise the message that was available.
        """


class TGChannel(abc.ABC):
    """Represents a request-reply session.

    Handles all initialization, tear down, and request-reply.
    """

    @property
    @abc.abstractmethod
    def linkstate(self) -> LinkState:
        """Gets whether the channel is connected, or any reason that it is not."""

    @property
    @abc.abstractmethod
    def properties(self) -> TGProperties:
        """Gets the properties for this channel."""

    @property
    @abc.abstractmethod
    def outboxaddr(self):
        """Gets the server-side address."""

    @property
    @abc.abstractmethod
    def protocolversions(self) -> tuple:
        """Gets the protocol version for both the client and server."""

    @property
    @abc.abstractmethod
    def authtoken(self) -> int:
        """Gets the authentication token that the server gave to us."""

    @property
    @abc.abstractmethod
    def sessionid(self) -> int:
        """Gets the session id for this session."""

    @property
    @abc.abstractmethod
    def clientid(self) -> str:
        """Gets this client's identifier."""

    @abc.abstractmethod
    def createSocket(self) -> TGSocket:
        """Creates and returns a default socket."""

    @abc.abstractmethod
    def connect(self):
        """Connect this channel to the server."""

    @abc.abstractmethod
    def start(self):
        """Starts the thread responsible for reading."""

    @abc.abstractmethod
    def disconnect(self):
        """Disconnect from the server."""

    @abc.abstractmethod
    def stop(self, stopmethod=TGChannelStopMethod.Graceful, msg=None):
        """Stop the listening thread."""

    @abc.abstractmethod
    def send(self, msg: tgpdu.TGMessage, response: TGChannelResponseWaiter = None) -> tgpdu.TGMessage:
        """Sends a message and waits for a response."""

    @abc.abstractmethod
    def readMessage(self) -> tgpdu.TGMessage:
        """Reads the message on the socket."""

    @property
    @abc.abstractmethod
    def waiters(self) -> Dict[int, TGChannelResponseWaiter]:
        """The waiters currently active.

        :returns: A dictionary with keys representing the request identifier and values of that channel response waiter.
        """

    @abc.abstractmethod
    def handleException(self, e: Exception) -> ExceptionHandleResult:
        """ Handles the exception.

        :param e: The exception to handle.
        :return: Returns how the exception was handled.
        """

    @property
    def isClosed(self):
        """Whether the channel is closed.

        :returns: True if the channel is closed.
        """
        return self.linkstate in (LinkState.Closed, LinkState.Closing, LinkState.Terminated)

    def inError(self):
        """Whether the channel is in an error."""
        return self.linkstate in (LinkState.FailedOnSend, LinkState.FailedOnRecv, LinkState.FailedOnProcessing)

    @property
    def isConnected(self):
        """Whether the channel is connected."""
        return self.linkstate == LinkState.Connected


    @classmethod
    def createChannel(cls, urlpath: str, username=None, password=None, dbName: str = None, props: TGProperties=None):
        """Creates a channel instance.

        :param urlPath: The URL for the database to connect to. When running on localhost, it should look like
            tcp://localhost:8222 for an insecure connection.
        :param username: The user's name of the database to use to connect to the database.
        :param password: The user's password for the database server.
        :param dbName: The database name to connect to.
        :param props: The properties to set up the connection with, contains any additional properties.

        :returns: The connection object requested.
        """
        klazz = None
        chprops = TGProperties(ConfigName.asMap())
        if props is not None:
            chprops.update(props)
        url = TGChannelUrl.parseUrl(urlpath)
        chprops.update(url.properties)

        if username is not None:
            chprops[ConfigName.ChannelUserID] = username

        if dbName is not None:
            chprops[ConfigName.ConnectionDatabaseName] = dbName

        if password is not None:
            chprops[ConfigName.ChannelPassword] = password

        if url.protocol == ProtocolType.Tcp:
            klazz = Class.forName('tgdb.impl.channelimpl.TcpChannel')
        elif url.protocol == ProtocolType.Ssl:
            klazz = Class.forName('tgdb.impl.channelimpl.SslChannel')
        else:
            raise tgexception.TGException("Invalid Url specified")
        return klazz(url, chprops)


class TGChannelUrl(abc.ABC):
    """Represents a parsed URL."""

    @property
    @abc.abstractmethod
    def protocol(self) -> ProtocolType:
        """Gets the protol corresponding with this URL"""

    @property
    @abc.abstractmethod
    def host(self) -> str:
        """The host for this URL."""

    @property
    @abc.abstractmethod
    def user(self) -> str:
        """The user for this URL."""

    @property
    @abc.abstractmethod
    def port(self) -> int:
        """The pot for this URL."""

    @property
    @abc.abstractmethod
    def properties(self) -> dict:
        """Any properties specified."""

    @property
    @abc.abstractmethod
    def url(self) -> str:
        """The full URL string."""

    @classmethod
    def parseUrl(cls, url):
        """Parse the URL and return a corresponding instance."""
        import tgdb.impl.channelimpl as tgchannelimpl
        return tgchannelimpl.LinkUrl.parse(url)
        








