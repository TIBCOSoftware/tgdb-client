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
 * File name : LongAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: LongAttribute.java 3154 2019-04-26 18:31:55Z sbangar $
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

class LongAttribute extends AbstractAttribute {

    LongAttribute(TGAttributeDescriptor desc)
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
        if (value instanceof Long) {
            setLong((Long) value);
        }
        else if (value instanceof Number) {
            setLong(((Number)value).longValue());
        }
        else if (value instanceof String) {
            setLong(ConversionUtils.string2Long(String.class.cast(value)));
        }
        else if (value instanceof Character) {
            char c = (char) value;
            setLong((long)Character.getNumericValue(c));
        }
        else if (value instanceof Boolean) {
            boolean b = (boolean) value;
            setLong(b ? 1L : 0L);
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Long, value.getClass().getSimpleName());
        }
    }

    private void setLong(Long l)
    {
        if (!isNull() && this.value.equals(l)) return;
        this.value = l;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        if (isNull()) return false;
        long d = (Long) this.value;
        return d > 0;
    }

    @Override
    public short getAsShort() {
        if (isNull()) return 0;
        Long d = (Long) this.value;
        return d.shortValue();
    }

    @Override
    public int getAsInt() {
        if (isNull()) return 0;
        Long d = (Long) this.value;
        return d.intValue();
    }

    @Override
    public long getAsLong() {
        if (isNull()) return 0;
        Long d = (Long) this.value;
        return d;
    }

    @Override
    public float getAsFloat() {
        if (isNull()) return 0;
        Long d = (Long) this.value;
        return d.floatValue();
    }

    @Override
    public double getAsDouble() {
        if (isNull()) return 0;
        Long d = (Long) this.value;
        return d.doubleValue();
    }

    public BigDecimal getAsNumber() {
        if (isNull()) return null;
        return new BigDecimal(this.value.toString());
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readLong();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeLong((Long)this.value);

    }
}
