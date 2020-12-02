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
 * File name : ByteAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: ByteAttribute.java 3881 2020-04-16 22:40:29Z nimish $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

class ByteAttribute extends AbstractAttribute {

    ByteAttribute(TGAttributeDescriptor desc) {
        super(desc);
    }

    @Override
    public void setValue(Object value) throws TGException {
    	if (value == null) 
        {
    		this.value = value;
    		setModified();
        	return;
        }

        if (value instanceof Byte) {
            setByte((Byte)value);
        }
        else if (value instanceof Boolean) {
            setByte((byte)(((Boolean)value) ? 1 : 0));
        }
        else if (value instanceof Number) {
            setByte(((Number)value).byteValue());
        }
        else if (value instanceof String) {
            setByte(new Byte((String)value));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Byte, value.getClass().getSimpleName());
        }
        setModified();
    }

    void setByte(Byte b) {
        if (!isNull() && (this.value.equals(b))) return;
        this.value = b;
        setModified();
    }
    @Override
    public boolean getAsBoolean() {
        byte b = (byte) this.value;
        return b == 0;

    }

    @Override
    public byte getAsByte() {
        return (byte) this.value;
    }

    @Override
    public char getAsChar() {

        byte b = (byte) this.value;
        return (char) b;

    }

    @Override
    public short getAsShort() {
        Byte b = (Byte) this.value;
        return b.shortValue();
    }

    @Override
    public int getAsInt() {
        Byte b = (Byte) this.value;
        return b.intValue();
    }

    @Override
    public long getAsLong() {
        Byte b = (Byte) this.value;
        return b.longValue();
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readByte();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeByte((byte)this.value);
    }
}
