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
 * File name : CharAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: CharAttribute.java 3154 2019-04-26 18:31:55Z sbangar $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

class CharAttribute extends AbstractAttribute {

    CharAttribute (TGAttributeDescriptor desc)
    {
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
        
        if (value instanceof Character) {
            setChar((Character)value);
        }
        else if (value instanceof Byte) {
            byte b = (Byte) value;
            setChar((char) b);
        }
        else if (value instanceof Number) {
            Number n = (Number) value;
            setChar((char) n.shortValue());
        }
        else if (value instanceof String) { //Unicode String value
            String s = value.toString();
            setChar(ConversionUtils.string2Character(s));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Char, value.getClass().getSimpleName());
        }

    }

    void setChar(Character c) {
        if (!isNull() && this.value.equals(c)) return;
        this.value = c;
        setModified();
    }

    @Override
    public char getAsChar() {
        return (char) this.value;
    }

    @Override
    public byte getAsByte() {
        return (byte)Character.getNumericValue((char)this.value);

    }

    @Override
    public short getAsShort() {
        return (short)Character.getNumericValue((char)this.value);
    }

    @Override
    public int getAsInt() {
        return Character.getNumericValue((char)this.value);
    }

    @Override
    public long getAsLong() {
        return Character.getNumericValue((char)this.value);
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readChar();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeChar((char)this.value);
    }
}
