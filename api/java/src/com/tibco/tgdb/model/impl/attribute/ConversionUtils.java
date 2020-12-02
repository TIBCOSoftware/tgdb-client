/**
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
 * File name : ConversionUtils.${EXT}
 * Created on: 06/22/2018
 * Created by: suresh
 * SVN Id: $Id: ConversionUtils.java 3996 2020-05-16 01:09:49Z vchung $
 */

package com.tibco.tgdb.model.impl.attribute;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.utils.TGEnvironment;

import static com.tibco.tgdb.exception.TGException.TGExceptionType.TypeConversionError;

public class ConversionUtils {

    //Do not change the id, if the order changes
    public enum TGAttributeValueType {
        Invalid(-1),
        Boolean(0),
        Byte(1),
        Char(2),
        Short(3),
        Integer(4),
        Long(5),
        Float(6),
        Double(7),
        BigDecimal(8),
        String(9),
        Calendar(10),
        ByteArray(11),
        CharArray(12),
        ShortArray(13),
        IntArray(14),
        LongArray(15),
        FloatArray(16),
        DoubleArray(17),
        InputStream(18),
        Externalizable(19),
        TGSerializable(20);
        
        public int id;
        TGAttributeValueType(int id) {
            this.id = id;
        }
        
        public static TGAttributeValueType fromId(int id)
        {
            if (id < 0) return  TGAttributeValueType.Invalid;
            for (TGAttributeValueType type:TGAttributeValueType.values()) {
                if (type.id == id) return type;
            }
            return TGAttributeValueType.Invalid;
        }
    }

    public static byte[] inputStream2ByteArray(InputStream in) throws TGException
    {
        ByteArrayOutputStream bos = new ByteArrayOutputStream();
        try {
            int cnt = 0;
            while ((cnt = in.available()) > 0) {
                byte buf[] = new byte[cnt];
                int readcnt = in.read(buf, 0, cnt);
                bos.write(buf, 0, readcnt);
            }
            return bos.toByteArray();
        }

        catch (IOException ioe)
        {
            throw new TGException(ioe);
        }

    }

    public static byte[] bigDecimal2ByteArray(BigDecimal bd) throws TGException
    {
        try {
            ByteArrayOutputStream bos = new ByteArrayOutputStream();
            DataOutputStream dos = new DataOutputStream(bos);
            byte buf[] = bd.unscaledValue().toByteArray();
            dos.writeInt(bd.scale());
            dos.writeInt(buf.length);
            dos.write(buf,0,buf.length);
            return bos.toByteArray();
        }
        catch (IOException ioe) {
            throw new TGException(ioe);
        }
    }

    public static BigDecimal byteArrayToBigDecimal(byte[] buf) throws TGException
    {
        try {
            ByteArrayInputStream bis = new ByteArrayInputStream(buf);
            DataInputStream dis = new DataInputStream(bis);
            int scale = dis.readInt();
            byte bd[] = new byte[dis.readInt()];
            dis.readFully(bd);
            BigInteger bi = new BigInteger(bd);
            return new BigDecimal(bi, scale);
        }
        catch (Exception e) {
            throw new TGException(e);
        }
    }

    public static Calendar string2Calendar(String s) throws TGException
    {
        try {
            //SS:TODO Get the Datetime format from the environment.
            //SS:Should we take per connection?
            Calendar cal = Calendar.getInstance();
            SimpleDateFormat sdf = new SimpleDateFormat(TGEnvironment.getInstance().getDefaultDateTimeFormat());
            cal.setTime(sdf.parse(s));
            return cal;
        }
        catch (ParseException pe) {
            throw new TGException(pe);
        }
    }

    public static String calendar2String(Calendar cal) throws TGException
    {
        try {
            //SS:TODO Get the Datetime format from the Connection
            SimpleDateFormat dtf = new SimpleDateFormat(TGEnvironment.getInstance().getDefaultDateTimeFormat());
            return dtf.format(cal.getTime());
        }
        catch (Exception e) {
            throw new TGException(e);
        }
    }

    public static Calendar long2Calendar(long l) throws TGException
    {
        Calendar cal = Calendar.getInstance();
        cal.setTimeInMillis(l);
        return cal;
    }

    public static Double string2Double(String s) throws TGException
    {
        try {
            return Double.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a double", s),
                    TypeConversionError, nfe);
        }
    }

    public static Float string2Float(String s) throws TGException
    {
        try {
            return Float.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a float", s),
                    TypeConversionError, nfe);
        }
    }

    public static Integer string2Integer(String s) throws TGException
    {
        try {
            return Integer.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a integer", s),
                    TypeConversionError, nfe);
        }
    }

    public static Long string2Long(String s) throws TGException
    {
        try {
            return Long.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a long", s),
                    TypeConversionError, nfe);
        }
    }

    public static Short string2Short(String s) throws TGException
    {
        try {
            return Short.valueOf(s);
        }
        catch (NumberFormatException nfe) {
            throw TGException.buildException(String.format("Cannot convert String:%s to a short", s),
                    TypeConversionError, nfe);
        }
    }

    public static Character string2Character(String s) throws TGException
    {
        if (s.length() == 1)  {
            return (char)s.charAt(0);
        }
        else { //let see if it is unicode point.
            try {
                int v = Integer.parseInt(s, 16); //The value is expected to be in "FFFF"
                return (char) v;
            }catch (NumberFormatException nfe) {
                throw TGException.buildException(String.format("Cannot convert String:%s to a character", s),
                        TypeConversionError, nfe);
            }
        }
    }

    /**
     * Try serializing the value to a byte array. This is used by Blob/Clob/and AEAD Data types
     * @param value
     * @return
     * @throws TGException
     * @throws IOException
     */
    public static byte[] toByteArray(Object value, TGAttributeType type) throws TGException, IOException
    {
        if (value == null) throw new TGException("Null value passed");
        ByteArrayOutputStream bos = new ByteArrayOutputStream(128);
        DataOutputStream dos = new DataOutputStream(bos);

        switch (type) {
            case Boolean:
                dos.writeBoolean((boolean)value);
                break;

            case Byte:
                dos.writeByte((int)value);
                break;

            case Char:
                dos.writeChar((char)value);
                break;

            case Short:
                dos.writeShort((Short)value);
                break;

            case Integer:
                dos.writeInt((int)value);
                break;

            case Long:
                dos.writeLong((long)value);
                break;

            case Float:
                dos.writeFloat((float)value);
                break;

            case Double:
                dos.writeDouble((double)value);
                break;

            case Number:
                byte[] numbuf = ConversionUtils.bigDecimal2ByteArray((BigDecimal)value);
                dos.writeInt(numbuf.length);
                dos.write(numbuf);
                break;

            case String:
                dos.writeUTF((String)value);
                break;

            case Time:
            case Date:
            case TimeStamp:
                String s = ConversionUtils.calendar2String((Calendar)value);
                dos.writeUTF(s);
                break;

            case Blob:
                byte[] buf = (byte[]) value;
                dos.writeInt(buf.length);
                dos.write(buf);
                break;

            case Clob:
                char[] cbuf = (char[]) value;
                dos.writeInt(cbuf.length);
                for (char c : cbuf) dos.writeChar(c);
                break;

            default:
                throw new TGTypeCoercionNotSupported(TGAttributeType.Blob, value.getClass().getSimpleName());
        }
        dos.flush();
        return bos.toByteArray();
    }

    public static Object fromByteArray(byte buf[], TGAttributeType type) throws TGException, IOException
    {
        if ((buf == null) || (buf.length == 0)) throw new TGException("Null value passed");
        ByteArrayInputStream bis = new ByteArrayInputStream(buf);
        DataInputStream dis = new DataInputStream(bis);

        switch (type) {
            case Boolean :
                return dis.readBoolean();
            case Byte :
                return dis.readByte();
            case Char :
                return dis.readChar();
            case Short :
                return dis.readShort();
            case Integer :
                return dis.readInt();
            case Long :
                return dis.readLong();
            case Float :
                return dis.readFloat();
            case Double :
                return dis.readDouble();
            case Number : {
                int len = dis.readInt();
                byte buf1[] = new byte[len];
                dis.readFully(buf1);
                return ConversionUtils.byteArrayToBigDecimal(buf1);
            }
            case String :
                return dis.readUTF();
            case Time :
            case Date:
            case TimeStamp:
                {
                    String s = dis.readUTF();
                    return ConversionUtils.string2Calendar(s);
                }
            case Blob : {
                int len = dis.readInt();
                byte buf1[] = new byte[len];
                dis.read(buf1, 0, len);
                return buf1;
            }
            case Clob :
            {
                int len = dis.readInt();
                char cb[] = new char[len];
                for (int i=0; i<len; i++) cb[i] = dis.readChar();
                return cb;
            }
            default:
                throw new TGTypeCoercionNotSupported(TGAttributeType.Blob, byte[].class.getSimpleName());
        }

    }



}
