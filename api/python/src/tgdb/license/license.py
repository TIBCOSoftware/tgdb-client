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
 * File name : license.py
 * Created on: 11/01/2019
 * Created by: Derek
 * SVN Id: $Id$
 *
 * Gets the license acceptance file.
 *
 * Not generally supposed to be used by client code.
"""

import os.path as ospath
import getpass
import time

import tgdb.version as tgvers
from tgdb.impl.pduimpl import ProtocolDataOutputStream
from tgdb.impl.pduimpl import ProtocolDataInputStream
import tgdb.log as tglog


def _getLicenseAcceptedPath():
    vers = tgvers.TGVersion()
    return ospath.expanduser('~/.tgdb/.license.accepted.' + str(vers.major))


def _getLicAcceptData():
    username = getpass.getuser()
    timestamp = time.time_ns()

    return username, timestamp


def writeLicense(licenseAcceptFlag: bool):

    outstream = ProtocolDataOutputStream()

    uname, timestamp = _getLicAcceptData()
    unamelen = len(uname)
    for i in range(32):
        if i < unamelen:
            outstream.writeByte(ord(uname[i]))
        else:
            outstream.writeByte(0)
    outstream.writeInt(unamelen, littleEndian=True)

    outstream.writeInt(0)               # Padding

    outstream.writeUnsignedLong(timestamp, True)
    outstream.writeByte(ord('y') if not licenseAcceptFlag else 0)
    outstream.writeBoolean(licenseAcceptFlag)

    outstream.writeShort(0)             # Padding
    outstream.writeInt(0)               # Padding

    version: tgvers.TGVersion = tgvers.TGVersion.getVersion()

    version.writeExternal(outstream)

    pos = outstream.position

    outstream.writeLong(0)

    outstream.writeHash64At(pos, 0, pos, littleEndian=True)

    # Create the directory if it does not already exist.
    if not ospath.exists(ospath.dirname(_getLicenseAcceptedPath())):
        import errno
        import os
        try:
            os.makedirs(ospath.dirname(_getLicenseAcceptedPath()))
        except OSError as err:
            raise err

    with open(_getLicenseAcceptedPath(), "wb") as file:
        file.write(outstream.buffer)


def readLicense() -> bool:
    ret: bool

    try:
        array: bytes
        with open(_getLicenseAcceptedPath(), "rb") as file:
            array = file.read()

        instream = ProtocolDataInputStream(bytearray(array))

        uname = ["\x00"] * 32

        for i in range(32):
            uname[i] = chr(instream.readByte())

        length = instream.readInt(True)
        if length > 32:
            tglog.gLogger.log(tglog.TGLevel.Error, 'Length stated is larger than expected: %d instead of a max of 32',
                              length)
            length = 32

        username = ""

        for i in range(length):
            username += uname[i]

        _ = instream.readInt()          # Padding to align the output properly

        timestamp = instream.readUnsignedLong(True)

        _ = instream.readByte()
        _ = instream.readBoolean()

        _ = instream.readShort()        # Padding
        _ = instream.readInt()          # Padding

        recv_version = tgvers.TGVersion.readExternal(instream)

        cur_version = tgvers.TGVersion.getVersion()

        act_uname, act_ts = _getLicAcceptData()

        recv_cs = instream.readUnsignedLong(True)
        act_cs = instream.readHash64From(0, instream.position - 8)

        ret = (act_uname[:length] == username[:length] and cur_version == recv_version and recv_cs == act_cs and
               act_ts >= timestamp)

    except (FileNotFoundError, IOError) as e:
        tglog.gLogger.log(tglog.TGLevel.Error, "Error: %s", str(e))
        ret = False

    return ret


if __name__ == '__main__':
    writeLicense(False)
    print(readLicense())