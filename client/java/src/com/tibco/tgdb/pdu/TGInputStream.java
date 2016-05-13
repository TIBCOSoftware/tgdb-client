package com.tibco.tgdb.pdu;

import java.io.DataInput;
import java.io.EOFException;
import java.io.IOException;
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
 * File name :TGInputStream
 * Created on: 12/21/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGInputStream.java 771 2016-05-05 11:40:52Z vchung $
 */
public interface TGInputStream extends DataInput
{

    /**
     * Read the current byte
     * @return
     * @throws IOException
     */
    public int read() throws IOException;


    /**
     * Read a Variable long field
     * @return
     * @throws EOFException
     */
    public long readVarLong() throws EOFException ;

    /**
     * Skip N bytes
     * @param n
     * @return
     * @throws IOException
     */
    public long skip(long n) throws IOException;


    /**
     * Mark the current Position
     * @param readlimit
     */
    public void mark(int readlimit);


    /**
     * Reset back to the old position
     */
    public  void reset();


    /**
     * Is Mark supported
     * @return
     */
    public boolean markSupported();


    /**
     *
     * @param b Buffer. buf cannot be null
     * @return
     * @throws IOException
     */
    public int read(byte[] b) throws IOException;


    /**
     * Similar to readFully.
     * @param b
     * @param off
     * @param len
     * @return
     * @throws IOException
     */
    public int read(byte[] b, int off, int len) throws IOException;

    /**
     * How much data is available on the stream
     * @return
     * @throws IOException
     */
    public int available() throws IOException;

    /**
     * Read an encoded byte array. writeBytes encodes the length, and the byte[]. This is equivalent to do a readInt, and read(byte[])
     * @return
     * @throws IOException
     */
    public byte[] readBytes() throws IOException;


    /**
     * Get the current positon of read
     * @return
     */
    public long getPosition();

    /**
     * Set the position of reading.
     * Atomic call.
     * @param position
     * @return
     */

    public long setPosition(long position);
    
    /**
     * Add a user maintained map for reference data 
     * @param map
     */
    public void setReferenceMap(Map map);
    
    /**
     * Return a user maintained reference map
     * @return
     */
    public Map getReferenceMap();
}
