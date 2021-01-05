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
 *  File name :attr.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: attrimpl.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates all the Attribute and its related classes
 """

from ctypes import c_float
from decimal import *
from datetime import *
from typing import *
from abc import *
import tgdb.pdu as tgpdu
from tgdb.impl.atomics import *
from tgdb.log import *
import tgdb.exception as tgexception
import tgdb.model as tgmodel

class AttributeDescriptorImpl(tgmodel.TGAttributeDescriptor):

    globalId: AtomicReference = AtomicReference('i', -1)

    def __init__(self, attrname=None, attrtype=tgmodel.TGAttributeType.Invalid):
        self.__id__ = AttributeDescriptorImpl.globalId.decrement()
        self.__name__ = attrname
        self.__type__ = attrtype
        self.__isArray__: bool = False
        self.__isEncrypted__: bool = False
        self.__precision__: int = 0
        self.__scale__: int = 0


    @property
    def id(self):
        return self.__id__

    @property
    def type(self):
        return self.__type__

    @property
    def isArray(self):
        return self.__isArray__

    @property
    def isEncrypted(self):
        return self.__isEncrypted__

    @property
    def precision(self):
        return self.__precision__

    @property
    def scale(self):
        return self.__scale__

    @property
    def systemType(self):
        return tgmodel.TGSystemType.AttributeDescriptor

    @property
    def name(self):
        return self.__name__

    def writeExternal(self, os):
        raise tgexception.TGException("Attribute Descriptor are readonly")

    def readExternal(self, input: tgpdu.TGInputStream):
        systype = input.readByte()
        if systype != tgmodel.TGSystemType.AttributeDescriptor.value:
            gLogger.log(TGLevel.Warning, "Attribute descriptor has invalid input stream value : {0}", systype)
        self.__id__ = input.readInt()
        self.__name__ = input.readUTF()
        self.__type__ = tgmodel.TGAttributeType.fromTypeId(input.readByte())
        self.__isArray__ = input.readBoolean()
        self.__isEncrypted__ = input.readBoolean()
        if self.__type__ == tgmodel.TGAttributeType.Number:
            self.__precision__ = input.readShort()
            self.__scale__ = input.readShort()


class AbstractAttribute(tgmodel.TGAttribute):

    def __init__(self, desc: tgmodel.TGAttributeDescriptor = None, owner: tgmodel.TGEntity = None):
        self.__desc__ = desc
        self.__owner__ = owner
        self._val: Any = None
        self.__isModified__ = False

    @property
    def descriptor(self) -> tgmodel.TGAttributeDescriptor:
        return self.__desc__

    @property
    def isNull(self) -> bool:
        return self._val is None

    @property
    def isModified(self) -> bool:
        return self.__isModified__

    def modified(self):
        self.__isModified__ = True

    def resetModified(self):
        self.__isModified__ = False

    @property
    def value(self) -> object:
        return self._val

    @property
    def owner(self) -> tgmodel.TGEntity:
        return self.__owner__

    @owner.setter
    def owner(self, value):
        self.__owner__ = value

    def __int__(self):
        return self.longValue

    def __str__(self):
        return self.stringValue

    def __bool__(self):
        return self.booleanValue

    def __repr__(self):
        return self.stringValue

    def __float__(self):
        return self.doubleValue

    def writeExternal(self, os: tgpdu.TGOutputStream):
        os.writeInt(self.descriptor.id)
        os.writeBoolean(self.isNull)
        if self.isNull:
            return
        if self.descriptor.isEncrypted:
            self.writeEncrypted(os)
        else:
            self.writeValue(os)

        return

    def readExternal(self, instream: tgpdu.TGInputStream):
        isnull = instream.readBoolean()
        if isnull:
            self._val = None
            return
        if (self.descriptor.isEncrypted and ((self.descriptor.type != tgmodel.TGAttributeType.Blob)
                                             or (self.descriptor.type != tgmodel.TGAttributeType.Clob))):
            self.readEncrypted(instream)
            return
        self.readValue(instream)

    def writeEncrypted(self, os: tgpdu.TGOutputStream):
        import tgdb.impl.pduimpl as pduimpl
        tmp_os = pduimpl.ProtocolDataOutputStream()
        self.writeValue(tmp_os)
        os.writeEncrypted(bytes(tmp_os.buffer))

    def readEncrypted(self, instream: tgpdu.TGInputStream):
        import tgdb.impl.pduimpl as pduimpl
        data = instream.readDecrypt()
        tmp_is = pduimpl.ProtocolDataInputStream(bytearray(data))
        self.readValue(tmp_is)

    @abstractmethod
    def writeValue(self, os: tgpdu.TGOutputStream): pass

    @abstractmethod
    def readValue(self, instream: tgpdu.TGInputStream): pass

    @classmethod
    def createAttribute(cls, owner: tgmodel.TGEntity, desc: tgmodel.TGAttributeDescriptor, _value=None) -> tgmodel.TGAttribute:

        aa: tgmodel.TGAttribute = None
        if desc.type == tgmodel.TGAttributeType.Boolean:
            aa = BooleanAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Byte:
            aa = ByteAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Char:
            aa = CharAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Short:
            aa = ShortAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Integer:
            aa = IntAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Long:
            aa = LongAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Float:
            aa = FloatAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Double:
            aa = DoubleAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.String:
            aa = StringAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Number:
            aa = NumberAttribute(desc, owner)
        elif desc.type in (tgmodel.TGAttributeType.Date, tgmodel.TGAttributeType.Time, tgmodel.TGAttributeType.TimeStamp):
            aa = TimestampAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Blob:
            aa = BlobAttribute(desc, owner)
        elif desc.type == tgmodel.TGAttributeType.Clob:
            aa = ClobAttribute(desc, owner)
        else:
            raise tgexception.TGException("Invalid descriptor:{0} encountered".format(desc))

        aa.owner = owner
        if _value is not None:
            aa._setValue(_value)

        return aa

    @classmethod
    def createFromStream(cls, owner: tgmodel.TGEntity, instream: tgpdu.TGInputStream) -> tgmodel.TGAttribute:

        aid = instream.readInt()
        if aid not in owner.graphMetadata.types:
            raise tgexception.TGException('Invalid attributeId:{0} encountered while deserialized'.format(aid))
        desc = owner.graphMetadata.types[aid]
        if desc is None or not isinstance(desc, tgmodel.TGAttributeDescriptor):
            raise tgexception.TGException('Invalid attributeId:{0} encountered while deserialized'.format(aid))

        aa: tgmodel.TGAttribute = AbstractAttribute.createAttribute(owner, desc)
        aa.readExternal(instream)

        return aa


class BooleanAttribute(AbstractAttribute):

    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        self._val = bool(_value)
        self.modified()

    @property
    def booleanValue(self) -> bool:
        return bool(self._val)

    @property
    def byteValue(self):
        return (0, 1)[self.booleanValue]

    @property
    def charValue(self):
        return ('N', 'Y')[self.booleanValue]

    @property
    def shortValue(self):
        return self.byteValue

    @property
    def intValue(self):
        return self.byteValue

    @property
    def longValue(self):
        return self.byteValue

    @property
    def stringValue(self) -> str:
        return ('False', 'True')[self.booleanValue]

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeBoolean(self._val)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readBoolean()


class ByteAttribute(AbstractAttribute):

    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value: int):

        try:
            abyte = int(_value) & 0xFF
            if self._val == abyte:
                return
            self._val = abyte
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.intValue]

    @property
    def byteValue(self):
        return int(self._val)

    @property
    def charValue(self):
        return chr(self.byteValue)

    @property
    def shortValue(self):
        return self.byteValue

    @property
    def intValue(self) -> int:
        return self.byteValue

    @property
    def longValue(self):
        return self.byteValue

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.byteValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeByte(self.byteValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = int(instream.readByte())


class CharAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):

        try:
            achar = chr(_value)
            if self._val == achar:
                return
            self._val = achar
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.shortValue > 0]

    @property
    def byteValue(self):
        return self.shortValue & 0xFF

    @property
    def charValue(self):
        return chr(self._val)

    @property
    def shortValue(self):
        return ord(self.charValue)

    @property
    def intValue(self) -> int:
        return self.shortValue

    @property
    def longValue(self):
        return self.shortValue

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.charValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeChar(self.charValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readChar()


class ShortAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):

        try:
            ashort = int(_value) & 0xFFFF
            if self._val == ashort:
                return
            self._val = ashort
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.shortValue > 0]

    @property
    def byteValue(self):
        return self.shortValue & 0xFF

    @property
    def charValue(self):
        return chr(self.shortValue)

    @property
    def shortValue(self):
        return int(self._val)

    @property
    def intValue(self) -> int:
        return self.shortValue

    @property
    def longValue(self):
        return self.shortValue

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.shortValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeShort(self.shortValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readShort()


class IntAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            aint = int(_value) & 0xFFFFFFFF
            if self._val == aint:
                return
            self._val = aint
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.intValue > 0]

    @property
    def byteValue(self):
        return self.intValue & 0xFF

    @property
    def charValue(self):
        return chr(self._val)

    @property
    def shortValue(self):
        return self.intValue & 0xFFFF

    @property
    def intValue(self) -> int:
        return int(self._val)

    @property
    def longValue(self):
        return self.intValue

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.intValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeInt(self.intValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readInt() & 0xFFFFFFFF


class LongAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            aint = int(_value) & 0xFFFFFFFFFFFFFFFF
            if self._val == aint:
                return
            self._val = aint
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.intValue > 0]

    @property
    def byteValue(self):
        return self.intValue & 0xFF

    @property
    def charValue(self):
        return chr(self._val)

    @property
    def shortValue(self):
        return self.intValue & 0xFFFF

    @property
    def intValue(self) -> int:
        return self.longValue & 0xFFFFFFFF

    @property
    def longValue(self):
        return int(self._val)

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.longValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeLong(self.longValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readLong()


class FloatAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            afloat = c_float(_value)
            if self._val == afloat:
                return
            self._val = afloat
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.intValue > 0]

    @property
    def byteValue(self):
        return self.intValue & 0xFF

    @property
    def shortValue(self):
        return self.intValue & 0xFFFF

    @property
    def intValue(self) -> int:
        return self.longValue & 0xFFFFFFFF

    @property
    def longValue(self):
        return int(self.floatValue)

    @property
    def floatValue(self):
        return c_float(self._val)

    @property
    def doubleValue(self):
        return self.floatValue

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.floatValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeFloat(self.floatValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readFloat()


class DoubleAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            afloat = float(_value)
            if self._val == afloat:
                return
            self._val = afloat
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.intValue > 0]

    @property
    def byteValue(self):
        return self.intValue & 0xFF

    @property
    def intValue(self) -> int:
        return self.longValue & 0xFFFFFFFF

    @property
    def longValue(self):
        return int(self.doubleValue)

    @property
    def floatValue(self):
        return self.doubleValue

    @property
    def doubleValue(self):
        return float(self._val)

    @property
    def stringValue(self) -> str:
        return "{0}".format(self.doubleValue)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeDouble(self.doubleValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readDouble()


class StringAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            if _value is None:
                self._val = None
                self.modified()
                return

            astr = str(_value)
            if self._val == astr:
                return
            self._val = astr
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))

    @property
    def booleanValue(self) -> bool:
        return (True, False)[self.stringValue is None or self.stringValue == '']

    @property
    def intValue(self) -> int:
        return self.longValue & 0xFFFFFFFF

    @property
    def longValue(self):
        try:
            aint = int(self.stringValue)
            return aint
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(self.stringValue))

    @property
    def floatValue(self):
        try:
            afloat = c_float(self.stringValue)
            return afloat
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(self.stringValue))

    @property
    def doubleValue(self):
        try:
            afloat = float(self.stringValue)
            return afloat
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(self.stringValue))

    @property
    def stringValue(self) -> str:
        return str(self._val)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeUTF(self.stringValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._val = instream.readUTF()


class NumberAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            if _value is None:
                self._val = None
                self.modified()
                return

            adec = Decimal(_value)
            adec = adec.scaleb(0, Context(prec=self.descriptor.precision, rounding=ROUND_HALF_UP))
            if self._val == adec:
                return
            self._val = adec
            self.modified()
            return
        except (ValueError, OverflowError, TypeError, DecimalException):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(_value))


    @property
    def intValue(self) -> int:
        return self.longValue & 0xFFFFFFFF

    @property
    def longValue(self):
        try:
            aDec:Decimal = self.value
            aint:int = int(aDec)
            return aint
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(self.stringValue))

    @property
    def floatValue(self):
        return self.doubleValue

    @property
    def doubleValue(self):
        try:
            aDec: Decimal = self.value
            afloat: float = float(aDec)
            return afloat
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherece value:{0} to integer".format(self.stringValue))

    @property
    def stringValue(self) -> str:
        return str(self._val)

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeShort(self.descriptor.precision)
        os.writeShort(self.descriptor.scale)
        os.writeUTF(self.stringValue)

    def readValue(self, instream: tgpdu.TGInputStream):
        prec = instream.readShort()
        scale = instream.readShort()
        dstr = instream.readUTF()
        adec = Decimal(dstr)
        adec = adec.scaleb(0, Context(prec, rounding=ROUND_HALF_UP))
        self._val = adec


"""
Python datetime only supports Gregorian Calendar.
TGDB support Julian calendar and can store from -8192 ... +8192 yrs. (BC and AD)
At some point need to replace the datetime using more sophisticated time
https://www.astropy.org 
"""


class TimestampAttribute(AbstractAttribute):
    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)

    def _setValue(self, _value):
        try:
            if _value is None:
               self._val = None
               self.modified()
               return

            adt = None
            if isinstance(_value, str):
                adt = datetime.fromisoformat(_value)
            elif isinstance(_value, float):
                adt = datetime.utcfromtimestamp(_value) # Always a UTC timestamp ???
            elif isinstance(_value, Tuple):
                t: list = [0, 0, 0, 0, 0, 0, 0]
                _v: Tuple = tuple(_value)
                for i in range(0, len(_v)):
                    t[i] = _v[i]
                adt = datetime(t[0], t[1], t[2], t[3], t[4], t[5], t[6])
            elif isinstance(_value, date):
                adt = datetime.fromisoformat(_value.isoformat())
            elif isinstance(_value, datetime):
                adt = datetime.fromisoformat(_value.isoformat())
            else:
                raise tgexception.TGException("invalid type {0} specified".format(type(_value)))

            if self._val == adt:
                return
            self._val = adt
            self.modified()
            return
        except (ValueError, OverflowError, TypeError):
            raise tgexception.TGException("Cannot coherce value:{0} to integer".format(_value))

    @property
    def stringValue(self) -> str:
        return str(self._val)

    @property
    def dateValue(self) -> date:
        if self._val is None:
            return date.min

        adt: datetime = self._val
        return adt.date()

    @property
    def timeValue(self) -> time:
        if self._val is None:
            return time.min

        adt: datetime = self._val
        return adt.time()

    @property
    def timestampValue(self) -> datetime:
        if self._val is None:
            return datetime.min

        adt: datetime = self._val
        return adt

    def writeValue(self, os: tgpdu.TGOutputStream):
        adt: datetime = self._val
        t = (adt.year, adt.month, adt.day, adt.hour, adt.minute, adt.second, adt.microsecond // 1000)
        os.writeBoolean(True), # only AD is supported
        os.writeShort(t[0])
        os.writeByte(t[1])
        os.writeByte(t[2])
        os.writeByte(t[3])
        os.writeByte(t[4])
        os.writeByte(t[5])
        os.writeShort(t[6])
        os.writeByte(-1) # TGNoZone = -1
        return

    def readValue(self, instream: tgpdu.TGInputStream):
        era = instream.readBoolean()
        yr  = instream.readShort()
        mon = instream.readByte()
        day = instream.readByte()
        hr  = instream.readByte()
        min = instream.readByte()
        sec = instream.readByte()
        msec = instream.readShort()
        tz  = instream.readByte()
        if tz > -1:
            instream.readShort()
        adt: datetime = datetime(yr, mon, day, hr, min, sec, msec, None)
        self._val = adt
        return


class BlobAttribute(AbstractAttribute):

    gUniqueId = AtomicReference('i', 0)

    def __init__(self, desc: tgmodel.TGAttributeDescriptor, owner: object = None):
        super().__init__(desc, owner)
        self._entityId = BlobAttribute.gUniqueId.decrement()
        self._cached = False

    def _setValue(self, value):
        if self._val != value:
            if value is None:
                self._val = value
            elif isinstance(value, bytes):
                self._val = value
            elif isinstance(value, str):
                self._val = bytes(value, encoding="utf-8")
            elif isinstance(value, bytearray):
                self._val = bytes(value)
            else:
                raise tgexception.TGTypeCoercionNotSupported("Could not convert {0} of type {1} to bytes."
                                                             .format(str(value), str(type(value))), totype=bytes)
            self.modified()

    def writeValue(self, os: tgpdu.TGOutputStream):
        os.writeLong(self._entityId)
        if self._val is None:
            os.writeBoolean(False)
        else:
            os.writeBoolean(True)
            os.writeBytes(self._val)

    def readValue(self, instream: tgpdu.TGInputStream):
        self._entityId = instream.readLong()
        self._cached = instream.readBoolean()
        if self._cached:
            self._val = instream.readBytes()

    @property
    def value(self) -> bytes:
        if self._entityId >= 0 or not self._cached:
            conn = self.owner.graphMetadata.graphObjectFactory.connection
            self._val = conn.getLargeObjectAsBytes(self._entityId, self.descriptor.isEncrypted)
        return self._val

    def asBytes(self):
        return self._val


class ClobAttribute(BlobAttribute):

    @property
    def value(self) -> str:
        return super().value.decode('utf-8')

    def asChars(self):
        return self.value
