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
 * File name : FloatAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: FloatAttribute.java 3631 2019-12-11 01:12:03Z ssubrama $
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

class FloatAttribute extends AbstractAttribute {

    FloatAttribute(TGAttributeDescriptor desc)
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

        if (value instanceof Float) {
            setFloat((Float) value);
        }
        else if (value instanceof Number) {
            setFloat(((Number)value).floatValue());
        }
        else if (value instanceof String) {
            setFloat(ConversionUtils.string2Float(String.class.cast(value)));
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Float, value.getClass().getSimpleName());
        }
    }

    void setFloat(Float f) {
        if (!isNull() && value.equals(f)) return;
        value = f;
        setModified();
    }

    @Override
    public boolean getAsBoolean() {
        if (isNull()) return false;
        float d = (Float) this.value;
        return d > 0;
    }


    @Override
    public short getAsShort() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.shortValue();
    }

    @Override
    public int getAsInt() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.intValue();
    }

    @Override
    public long getAsLong() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.longValue();
    }

    @Override
    public float getAsFloat() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.floatValue();
    }

    @Override
    public double getAsDouble() {
        if (isNull()) return 0;
        Float d = (Float) this.value;
        return d.doubleValue();
    }

    public BigDecimal getAsNumber() {
        if (isNull()) return null;
        return new BigDecimal(this.value.toString());
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readFloat();

    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeFloat((Float)this.value);
    }
}
