package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.utils.TGConstants;

import java.io.EOFException;
import java.io.IOException;
import java.io.InputStream;
import java.io.UTFDataFormatException;
import java.util.Map;

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
 * File name :ProtocolDataInputStream
 * Created on: 1/31/15
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: ProtocolDataInputStream.java 771 2016-05-05 11:40:52Z vchung $
 */
public class ProtocolDataInputStream extends InputStream implements TGInputStream{

    byte[]          buffer  = null;
    int             buflen  = 0;
    int             curpos  = 0;
    String          encoding = null;
    transient private char[] charbuf = null;
    transient private int mark = 0;
    Map             referenceMap = null;



    public ProtocolDataInputStream(byte[] buf)
    {
        this.buffer = buf;
        buflen = this.buffer.length;

    }

    @Override
    public void readFully(byte[] ba) throws IOException {

        if (ba == null) throw new NullPointerException("input argument to readfully is null.");
        readFully(ba, curpos, ba.length);

    }

    @Override
    public void readFully(byte[] ba, int readcurpos, int readlen) throws IOException {

        if (readlen <= 0) return;
        if (curpos+readlen > buflen) throw new EOFException();
        System.arraycopy(buffer,curpos,ba,readcurpos,readlen);
        curpos += readlen;

    }

    @Override
    public int skipBytes(int n) throws IOException {
        if (n < 0) return 0;
        int avail = buflen - curpos;
        if (avail <= n) {
            curpos = buflen;
            return avail;
        }
        curpos += n;
        return n;
    }

    @Override
    public boolean readBoolean() throws IOException {
        if (curpos >= buflen) throw new EOFException();
        return ((buffer[curpos++] == (byte)0) ? false : true);
    }

    @Override
    public byte readByte() throws IOException {
        if (curpos >= buflen) throw new EOFException();
        return buffer[curpos++];
    }

    @Override
    public int readUnsignedByte() throws IOException {
        if (curpos >= buflen) throw new EOFException();
        return (int)(buffer[curpos++]&0xff);
    }

    @Override
    public short readShort() throws IOException {
        if (curpos+2 > buflen) throw new EOFException();
        return (short)(((buffer[curpos++] << 8)) + (buffer[curpos++]&0xff));
    }

    @Override
    public int readUnsignedShort() throws IOException {
        if (curpos+2 > buflen) throw new EOFException();
        int i = (((buffer[curpos++] << 8)) + (buffer[curpos++]&0x00ff));
        return (i & 0x0000ffff);
    }

    @Override
    public char readChar() throws IOException {
        if (curpos+2 > buflen) throw new EOFException();
        return (char)(((buffer[curpos++] << 8)) + (buffer[curpos++]&0xff));
    }

    @Override
    public int readInt() throws IOException {

        if (curpos+4 > buflen) throw new EOFException();
        return ( ((buffer[curpos++] << 24)) +
                ((buffer[curpos++] << 16)&0x00ff0000) +
                ((buffer[curpos++] << 8)&0x0000ff00) +
                ((buffer[curpos++])&0xff));
    }

    @Override
    public long readLong() throws IOException {
        if (curpos+8 > buflen) throw new EOFException();
        return ((long)(readInt())<<32) + (readInt()&0xffffffffL);
    }

    @Override
    public float readFloat() throws IOException {
        return Float.intBitsToFloat(readInt());
    }

    @Override
    public double readDouble() throws IOException {
        return Double.longBitsToDouble(readLong());
    }

    @Override
    public String readLine() throws IOException {
        return null;
    }

    @Override
    public String readUTF() throws IOException {
        int start = curpos;
        try
        {
            int utflen = readUnsignedShort();
            return readUTFString(utflen);
        }
        catch(UTFDataFormatException utfe)
        {
            curpos = start;
            throw utfe;
        }
        catch(EOFException eofe)
        {
            curpos = start;
            throw eofe;
        }
    }

    @Override
    public int read() throws IOException {
        return (curpos < buflen) ? (buffer[curpos++] & 0xff) : -1;
    }

    public final String readUTFString(int utflen) throws EOFException, UTFDataFormatException {
        if (curpos+utflen > buflen) throw new EOFException();
        int start = curpos;
        int c, cval, char2, char3;
        int strlen=0, lastpos=curpos+utflen;

        if (charbuf == null || charbuf.length < utflen)
            charbuf = new char[utflen==0 ? 16 : utflen];

        while(curpos < lastpos) {
            c = (int) buffer[curpos++] & 0xff;
            cval = c >> 4;
            if (cval <= 7)   // 0xxxxxxx
            {
                charbuf[strlen++] = (char)c;
            }
            else
            if (cval == 12 || cval == 13)   // 110x xxxx   10xx xxxx
            {
                if (curpos+1 > lastpos)
                {
                    curpos = start;
                    throw new UTFDataFormatException();
                }
                char2 = (int) buffer[curpos++];
                if ((char2 & 0xC0) != 0x80)
                {
                    curpos = start;
                    throw new UTFDataFormatException();
                }
                charbuf[strlen++] = (char)(((c & 0x1F) << 6) | (char2 & 0x3F));
            }
            else
            if (cval == 14)   // 1110 xxxx  10xx xxxx  10xx xxxx
            {
                if (curpos+2 > lastpos)
                {
                    curpos = start;
                    throw new UTFDataFormatException();
                }
                char2 = (int)buffer[curpos++];
                char3 = (int)buffer[curpos++];
                if (((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80))
                {
                    curpos = start;
                    throw new UTFDataFormatException();
                }
                charbuf[strlen++] = (char)(((c & 0x0F) << 12) |
                        ((char2 & 0x3F) << 6)  |
                        ((char3 & 0x3F) << 0));
            }
            else    // 10xx xxxx,  1111 xxxx
            {
                curpos = start;
                throw new UTFDataFormatException();
            }
        }
        if (curpos != lastpos)
        {
            curpos = start;
            throw new UTFDataFormatException();
        }
        return new String(charbuf,0,strlen);
    }



    //--------------- readVarLong ------------------------------

    public long readVarLong() throws EOFException
    {
        if (curpos >= buflen) throw new EOFException();

        byte lenByte = buffer[curpos];

        if (lenByte == TGConstants.U64PACKED_NULL)
        {
            curpos++;
            return TGConstants.U64_NULL;
        }

        if ((lenByte & ((byte)0x80)) == 0)
        {
            curpos++;
            return (long)lenByte;
        }

        if ((lenByte & ((byte)0x40)) == 0)
        {
            if (curpos+2 > buflen) throw new EOFException();
            short s = (short)(((buffer[curpos++] << 8)) + (buffer[curpos++]&0xff));
            return (long)(s & 0x3fff);
        }

        if ((lenByte & ((byte)0x20)) == 0)
        {
            if (curpos+4 > buflen) throw new EOFException();
            int ival = (((buffer[curpos++] << 24)) +
                    ((buffer[curpos++] << 16)&0x00ff0000) +
                    ((buffer[curpos++] << 8)&0x0000ff00) +
                    ((buffer[curpos++])&0xff));
            return (long)(ival & 0x1fffffff);
        }

        int count = (int)(lenByte & ((byte)0x0f));
        curpos++;

        if (curpos+count > buflen) throw new EOFException();

        long val = 0;
        for (int i=0; i<count; i++)
        {
            val <<= 8;
            val |= (buffer[curpos++] & 0x00ff);
        }

        return val;

    }



    @Override
    public long skip(long n) throws IOException {
        return skipBytes((int)n);
    }

    @Override
    public synchronized void mark(int readlimit) {
        mark = readlimit;
    }

    @Override
    public synchronized void reset()  {
        curpos = mark;
    }

    @Override
    public boolean markSupported() {
        return true;
    }

    @Override
    public int read(byte[] b) throws IOException {
        return read(b, 0, b.length);
    }



    @Override
    public int read(byte[] b, int off, int len) throws IOException {
        if (b == null) {
            throw new NullPointerException();
        }
        else
        if ((off < 0) || (off > b.length) || (len < 0) ||
                ((off + len) > b.length) || ((off + len) < 0)) {
            throw new IndexOutOfBoundsException();
        }

        if (curpos >= buflen)
            return -1;

        if (curpos + len > buflen)
            len = buflen - curpos;

        if (len <= 0)
            return 0;

        System.arraycopy(buffer, curpos, b, off, len);
        curpos += len;
        return len;
    }

    @Override
    public int available() throws IOException {
        return buflen - curpos;
    }

    public byte[] readBytes() throws IOException {
        int len = readInt();
        if (len == 0) return TGConstants.EmptyByteArray;
        if (len == -1) throw new IOException("Read data corrupt");
        byte[] buf = new byte[len];
        read(buf);
        return buf;
    }


    public long getPosition() {
        return curpos;
    }

    /**
     * Atomic call.
     * @param position
     * @return
     */

    public long setPosition(long position) {
        int oldPos = curpos;
        curpos = (int) position;

        return oldPos;
    }
    
    public void setReferenceMap(Map map) {
    	referenceMap = map;
    }
    
    public Map getReferenceMap() {
    	return referenceMap;
    }
}
