package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.utils.TGConstants;

import java.io.IOException;
import java.io.OutputStream;
import java.io.UTFDataFormatException;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * <p/>
 * File name :ProtocolDataOutputStream
 * Created on: 1/31/15
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: ProtocolDataOutputStream.java 583 2016-03-15 02:02:39Z vchung $
 */

public class ProtocolDataOutputStream extends OutputStream implements TGOutputStream {


    byte[]          buffer  = null;
    int             buflen  = 0;
    int count = 0;
    transient private char[] charbuf = null;



    public ProtocolDataOutputStream() {
        buffer = new byte[256];
        buflen = buffer.length;
    }

    public ProtocolDataOutputStream(int len) {
        buffer = new byte[len];
        buflen = buffer.length;
    }


    public void skip(int n) {
        if (count + n > buflen) ensure(n);
        count += n;
    }

    /* only used in write mode */
    final void ensure(int len) {
        if (count + len <= buflen) return;
        int newlen=0;
        if (len > 100000)
            newlen = count + len + 2048;
        else newlen = (count + len) * 2;
        byte[] b = new byte[newlen];
        System.arraycopy(buffer, 0, b, 0, count);
        buffer = b;
        buflen = newlen;
    }


    public int getPosition() {
        return count;
    }


    public final byte[] getBuffer() {
        return buffer;
    }

    public final int getLength() {
        return count;
    }


    @Override
    public void writeBoolean(boolean value) throws IOException {
        if (count >= buflen) ensure(1);
        buffer[count++] = value ? (byte) 1 : (byte) 0;
    }

    @Override
    public void writeByte(int value) throws IOException {
        if (count >= buflen) ensure(1);
        buffer[count++] = (byte) value;
    }

    @Override
    public void writeShort(int value) throws IOException {
        if (count + 2 > buflen) ensure(2);
        buffer[count++] = (byte) (value >>> 8);
        buffer[count++] = (byte) (value);
    }

    @Override
    public void writeChar(int value) throws IOException {
        if (count + 2 > buflen) ensure(2);
        buffer[count++] = (byte) (value >>> 8);
        buffer[count++] = (byte) (value);
    }

    @Override
    public void writeInt(int value) throws IOException {

        if (count + 4 > buflen) ensure(4);

        buffer[count++] = (byte) (value >> 24);
        buffer[count++] = (byte) (value >> 16);
        buffer[count++] = (byte) (value >> 8);
        buffer[count++] = (byte) (value);

    }

    @Override
    public void writeLong(long value) throws IOException {

        if (count + 8 > buflen) ensure(8);

        // splitting in into two ints, is much faster
        // than shifiting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
        // also much much faster than old code.
        int a = (int)(value>>32);
        int b = (int)value;

        buffer[count] = (byte) (a >> 24);
        buffer[count + 1] = (byte) (a >> 16);
        buffer[count + 2] = (byte) (a >> 8);
        buffer[count + 3] = (byte) a;
        buffer[count + 4] = (byte) (b >> 24);
        buffer[count + 5] = (byte) (b >> 16);
        buffer[count + 6] = (byte) (b >> 8);
        buffer[count + 7] = (byte) b;
        count += 8;

    }

    @Override
    public void writeFloat(float value) throws IOException {
        writeInt(Float.floatToIntBits(value));
    }

    @Override
    public void writeDouble(double value) throws IOException {
        writeLong(Double.doubleToLongBits(value));
    }

    @Override
    public void writeBytes(String s) throws IOException {
        byte[] buf = s != null ? s.getBytes() : TGConstants.EmptyByteArray;
        writeBytes(buf);
    }


    public void writeBytes(byte[] buf) throws IOException {
        int len = buf.length;
        ensure(len+4);
        writeInt(len);
        for (int i = 0 ; i < len ; i++) {
            int v = buf[i];
            writeByte(v);
        }
    }



    @Override
    public void writeChars(String s) throws IOException {
        int len = s.length();
        ensure(len);
        for (int i = 0 ; i < len ; i++) {
            int v = s.charAt(i);
            writeChar(v);
        }
    }

    @Override
    public void writeUTF(String str) throws IOException {

        int start = count, len = 0;

        if (count + 2 > buflen) ensure(2 + str.length() * 3);
        count += 2;

        try {
            len = writeUTFString(str);
        }
        catch(UTFDataFormatException e) {
            count = start;
            throw e;
        }

        /* now write length */
        buffer[start]   = (byte)((len >>> 8) & 0xFF);
        buffer[start+1] = (byte)(len & 0xFF);

    }

    @Override
    public void write(int value) throws IOException {
        if (count >= buflen) ensure(1);
        buffer[count++] = (byte) value;
    }

    @Override
    public void write(byte[] value, int writepos, int writelen) throws IOException {
        if (value == null) {
            throw new NullPointerException();
        } else if ((writepos < 0) || (writepos > value.length) || (writelen < 0) ||
                ((writepos + writelen) > value.length) || ((writepos + writelen) < 0)) {
            throw new IndexOutOfBoundsException();
        }
        else
        if (writelen == 0) {
            return;
        }
        ensure(writelen);
        System.arraycopy(value, writepos, buffer, count, writelen);
        count += writelen;
    }

    /**
     * write a long value as varying length into the buffer.
     * @param value
     */
    public final void writeVarLong(long value)
    {
        if (value == TGConstants.U64_NULL)
        {
            if (count >= buflen) ensure(1);
            buffer[count++] = TGConstants.U64PACKED_NULL;
            return;
        }

        if (value < 0)
            throw new InternalError("Can not pack negative long value");

        if (value <= 0x7f) {
            if (count >= buflen) ensure(1);
            buffer[count++] = (byte) value;
            return;
        }

        if (value <= 0x3fff) {
            if (count + 2 > buflen) ensure(2);
            value |= 0x00008000;
            buffer[count++] = (byte) (value >>> 8);
            buffer[count++] = (byte) (value);
            return;
        }

        if (value <= 0x1fffffff) {
            if (count + 4 > buflen) ensure(4);
            value |= 0xC0000000;
            count += 4;
            for (int i=1; i<=4; i++,value>>>=8)
                buffer[count - i] = (byte) value;
            return;
        }

        /* we may need up to 9 bytes */
        if (count + 9 > buflen) ensure(9);

        /* calculate the number of non-zero bytes */
        long mask = 0xff00000000000000L;
        int  count = 8;
        for (int i=0; i<8; i++) {
            if ((value & mask) != 0L)
                break;
            count--;
            mask >>>= 8;
        }

        byte b = (byte)(count | 0xE0);
        buffer[count++] = b;

        count += count;
        for (int i=1; i<=count; i++,value>>>=8)
            buffer[count - i] = (byte) value;
    }

    /*
     *
     */
    public final int writeUTFString(String str) throws UTFDataFormatException {
        int i = 0, c, strlen = str.length(), start = count;

        if ((count + 3 * strlen) > buflen) ensure(strlen * 3);

        if (charbuf == null || (charbuf.length < strlen))
            charbuf = new char[strlen != 0 ? strlen : 16];

        str.getChars(0,strlen,charbuf,0);

        while(i<strlen) {
            c = charbuf[i++];
            if ((c >= 0x0001) && (c <= 0x007F)) {
                buffer[count++] = (byte) c;
            }
            else
            if (c > 0x07FF) {
                buffer[count++] = (byte) (0xE0 | ((c >> 12) & 0x0F));
                buffer[count++] = (byte) (0x80 | ((c >> 6) & 0x3F));
                buffer[count++] = (byte) (0x80 | ((c >> 0) & 0x3F));
            }
            else {
                buffer[count++] = (byte) (0xC0 | ((c >> 6) & 0x1F));
                buffer[count++] = (byte) (0x80 | ((c >> 0) & 0x3F));
            }
        }
        int len = count - start;
        if (len > 65535) {
            count = start;
            throw new UTFDataFormatException("String is too long");
        }
        return len;
    }


    public int writeBooleanAt(int pos, boolean value) throws TGException {
        if (pos < count) {
            buffer[pos++] = value ? (byte) 1 : (byte) 0;
            return pos;
        } else {
            throw new TGException("Invalid position specified :" + pos);
        }
    }


    public int writeByteAt(int pos, int value) throws TGException {
        if (pos < count) {
            buffer[pos++] = (byte) value;
            return pos;
        } else {
            throw new TGException("Invalid position specified :" + pos);
        }
    }


    public int writeShortAt(int pos, int value) throws TGException {
        if (pos + 2 < count) {
            buffer[pos++] = (byte) (value >>> 8);
            buffer[pos++] = (byte) (value);
            return pos;
        } else {
            throw new TGException("Invalid position specified :" + pos);
        }
    }


    public int writeCharAt(int pos, int value) throws TGException {
        if (pos + 2 < count) {
            buffer[pos++] = (byte) (value >>> 8);
            buffer[pos++] = (byte) (value);
            return pos;
        } else {
            throw new TGException("Invalid position specified :" + pos);
        }
    }


    public int writeIntAt(int pos, int value) throws TGException {

        if (pos + 4 < count) {

            buffer[pos++] = (byte) (value >> 24);
            buffer[pos++] = (byte) (value >> 16);
            buffer[pos++] = (byte) (value >> 8);
            buffer[pos++] = (byte) (value);
            return pos;
        } else {
            throw new TGException("Invalid position specified :" + pos);
        }

    }


    public int writeLongAt(int pos, long value) throws TGException {

        if (pos + 8 >= count) {
            throw new TGException("Invalid position specified :" + pos);
        }

        // splitting in into two ints, is much faster
        // than shifiting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
        // also much much faster than old code.
        int a = (int) (value >> 32);
        int b = (int) value;


        buffer[pos] = (byte) (a >> 24);
        buffer[pos + 1] = (byte) (a >> 16);
        buffer[pos + 2] = (byte) (a >> 8);
        buffer[pos + 3] = (byte) a;
        buffer[pos + 4] = (byte) (b >> 24);
        buffer[pos + 5] = (byte) (b >> 16);
        buffer[pos + 6] = (byte) (b >> 8);
        buffer[pos + 7] = (byte) b;

        return pos + 8;

    }


    public int writeFloatAt(int pos, float value) throws TGException {
        return writeIntAt(pos, Float.floatToIntBits(value));
    }


    public int writeDoubleAt(int pos, double value) throws TGException {
        return writeLongAt(pos, Double.doubleToLongBits(value));
    }


    public int writeBytesAt(int pos, String s) throws TGException {
        int len = s.length();
        if ((pos + len) >= count) throw new TGException("Invalid position specified :" + pos);

        for (int i = 0; i < len; i++) {
            int v = s.charAt(i);
            pos = writeByteAt(pos, v);
        }
        return pos;
    }


    public int writeCharsAt(int pos, String s) throws TGException {
        int len = s.length();
        if ((pos + len) >= count) throw new TGException("Invalid position specified :" + pos);
        for (int i = 0; i < len; i++) {
            int v = s.charAt(i);
            pos = writeCharAt(pos, v);
        }
        return pos;
    }

}
