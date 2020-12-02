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
 * File name : BooleanAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: BooleanAttribute.java 3154 2019-04-26 18:31:55Z sbangar $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

class BooleanAttribute extends AbstractAttribute {

    BooleanAttribute(TGAttributeDescriptor type) {
        super(type);
        this.value = false;
    }

    @Override
    public void setValue(Object value) throws TGException {
    	if (value == null) 
        {
    		this.value = value;
    		setModified();
        	return;
        }
    	
        if (value instanceof Boolean) {
            setBoolean((Boolean)value);
        }
        else if (value instanceof Number) {
            int v = (int)value;
            setBoolean(v > 0);
        }
        else if (value instanceof String) {
            setBoolean(Boolean.valueOf(value.toString()));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Boolean, value.getClass().getSimpleName());
        }
    }

    void setBoolean(Boolean b) {
        setModified();
    	if (!isNull() && (this.value.equals(b))) return;
        this.value = b;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        return Boolean.class.cast(this.value);
    }

    @Override
    public byte getAsByte() {
        Boolean b = (Boolean) this.value;
        return (byte) (b ? 1 : 0);
    }

    @Override
    public char getAsChar() {
        Boolean b = (Boolean) this.value;
        return (char) (b ? 1 : 0);
    }

    @Override
    public short getAsShort() {
        Boolean b = (Boolean) this.value;
        return (short) (b ? 1 : 0);
    }

    @Override
    public int getAsInt() {
        Boolean b = (Boolean) this.value;
        return (int) (b ? 1 : 0);
    }

    @Override
    public long getAsLong() {
        Boolean b = (Boolean) this.value;
        return (long) (b ? 1L : 0L);
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeBoolean(Boolean.class.cast(value));
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        value = is.readBoolean();
    }
}
