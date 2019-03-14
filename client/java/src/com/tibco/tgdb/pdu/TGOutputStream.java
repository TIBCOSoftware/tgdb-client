package com.tibco.tgdb.pdu;

import com.tibco.tgdb.exception.TGException;

import java.io.DataOutput;
import java.io.IOException;
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
 * File name :TGOutputStream
 * Created on: 12/17/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGOutputStream.java 2343 2018-06-08 18:33:46Z ssubrama $
 */

public interface TGOutputStream extends  DataOutput {

    /**
     * Skip N bytes. Allocate if necessary
     * @param n
     */
    void skip(int n);

    /**
     * Get the current write position
     * @return
     */
    public int getPosition();

    /**
     * Get the underlying Buffer
     * @return
     */
    public byte[] getBuffer();

    /**
     * Get the total write length
     * @return
     */
    public int getLength();

    /**
     * Convenient WriteBytes, which writes the len, and the byte[] buf
     * @param buf
     * @throws IOException
     */
    public void writeBytes(byte[] buf) throws IOException ;

    /**
     * write a long value as varying length into the buffer.
     * @param value
     */
    public void writeVarLong(long value) throws IOException;

    /**
     * Write UTFString
      * @param str
     * @return
     * @throws UTFDataFormatException
     */
    public int writeUTFString(String str) throws UTFDataFormatException;

    /**
     * Write boolean at a given position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeBooleanAt(int pos, boolean value) throws TGException;

    /**
     * Write a byte at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeByteAt(int pos, int value) throws TGException;

    /**
     * Write Short at a position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeShortAt(int pos, int value) throws TGException;

    /**
     * Write a Java Char at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeCharAt(int pos, int value) throws TGException;

    /**
     * Write Integer at the position.Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeIntAt(int pos, int value) throws TGException;

    /**
     * Write Long at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeLongAt(int pos, long value) throws TGException;

    /**
     * Write Float at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeFloatAt(int pos, float value) throws TGException;

    /**
     * Write Double at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param value
     * @return
     * @throws TGException
     */
    public int writeDoubleAt(int pos, double value) throws TGException;


    /**
     * Write String at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param s
     * @return
     * @throws TGException
     */
    public int writeBytesAt(int pos, String s) throws TGException;

    /**
     * write Chars at the position. Buffer should have sufficient space to write the content.
     * @param pos
     * @param s
     * @return
     * @throws TGException
     */
    public int writeCharsAt(int pos, String s) throws TGException;


}
