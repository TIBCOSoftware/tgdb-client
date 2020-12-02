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
 * File name : ShortAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: ShortAttribute.java 3631 2019-12-11 01:12:03Z ssubrama $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.math.BigDecimal;

class ShortAttribute extends AbstractAttribute {

    ShortAttribute (TGAttributeDescriptor desc) {
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
        if (value instanceof Short) {
            setShort((Short.class.cast(value)));
        }
        else if (value instanceof Number) {
            setShort(((Number)value).shortValue());
        }
        else if (value instanceof String) {
            setShort(ConversionUtils.string2Short(String.class.cast(value)));
        }
        else if (value instanceof Character) {
            char c = (char) value;
            setShort((short)Character.getNumericValue(c));
        }
        else if (value instanceof Boolean) {
            boolean b = (boolean) value;
            setShort((short)(b ? 1 : 0));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Short, value.getClass().getSimpleName());
        }
    }

    private void setShort(Short s) {
        if (!isNull() && this.value.equals(s)) return;
        this.value = s;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        if (isNull()) return false;
        long d = (Short) this.value;
        return d > 0;
    }

    @Override
    public short getAsShort() {
        if (isNull()) return 0;
        Short d = (Short) this.value;
        return d.shortValue();
    }

    @Override
    public int getAsInt() {
        if (isNull()) return 0;
        Short d = (Short) this.value;
        return d.intValue();
    }

    @Override
    public long getAsLong() {
        if (isNull()) return 0;
        Short d = (Short) this.value;
        return d;
    }

    @Override
    public float getAsFloat() {
        if (isNull()) return 0;
        Short d = (Short) this.value;
        return d.floatValue();
    }

    @Override
    public double getAsDouble() {
        if (isNull()) return 0;
        Short d = (Short) this.value;
        return d.doubleValue();
    }

    public BigDecimal getAsNumber() {
        if (isNull()) return null;
        return new BigDecimal(this.value.toString());
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        if (this.desc.isEncrypted()) {
            readDecrypted(is);
            return;
        }
        this.value = is.readShort();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeShort((Short)this.value);
    }
}
