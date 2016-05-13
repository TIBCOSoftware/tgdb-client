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
 * File name : ByteArrayEntityId.${EXT}
 * Created on: 1/28/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: ByteArrayEntityId.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.model.impl;


import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGEntityId;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.util.Arrays;

public class ByteArrayEntityId implements TGEntityId {

    byte[] entityId;

    ByteArrayEntityId(long value) {
        entityId = new byte[16];

        int pos = writeLongAt(0, Long.MAX_VALUE);
        writeLongAt(pos, value);
        
    }
    /**
     * Make a copy.
     * @param buf
     * @throws TGException
     */
    public ByteArrayEntityId(byte[] buf) throws TGException {
        if (buf == null) throw new TGException("Cannot construct from null byte array");

        if (buf.length != 16) throw new TGException("Invalid buffer specified");

        entityId = Arrays.copyOf(buf, buf.length);
    }

    @Override
    public boolean equals(Object obj) {
        if (! (obj instanceof ByteArrayEntityId)) return false;

        byte[] b1 = this.entityId;
        byte[] b2 = ((ByteArrayEntityId)obj).entityId;
        return Arrays.equals(b1, b2);

    }

    /**
     * Return a copy
     * @return
     */
    @Override
    public byte[] toBytes() {
        return Arrays.copyOf(entityId, entityId.length);
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
        os.writeBytes(entityId);
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
        entityId = is.readBytes();
    }

    private int writeLongAt(int pos, long value) {


        // splitting in into two ints, is much faster
        // than shifiting bits for a long i.e (byte)val >> 56, (byte)val >> 48 ..
        // also much much faster than old code.
        int a = (int) (value >> 32);
        int b = (int) value;


        entityId[pos] = (byte) (a >> 24);
        entityId[pos + 1] = (byte) (a >> 16);
        entityId[pos + 2] = (byte) (a >> 8);
        entityId[pos + 3] = (byte) a;
        entityId[pos + 4] = (byte) (b >> 24);
        entityId[pos + 5] = (byte) (b >> 16);
        entityId[pos + 6] = (byte) (b >> 8);
        entityId[pos + 7] = (byte) b;

        return pos + 8;

    }

}
