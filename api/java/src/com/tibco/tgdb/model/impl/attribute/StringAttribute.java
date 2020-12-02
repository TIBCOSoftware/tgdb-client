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
 * File name : StringAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: StringAttribute.java 3882 2020-04-17 00:27:45Z nimish $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.impl.GraphMetadataImpl;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGProperties;

import java.io.IOException;
import java.math.BigDecimal;
import java.nio.CharBuffer;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Locale;

class StringAttribute extends AbstractAttribute {
    final static int MAX_STRING_ATTR_LENGTH = 1000 - Short.BYTES - 1;

    StringAttribute(TGAttributeDescriptor desc)
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
        String s = value.toString();
        int len = utfLength(s);
        if (len > MAX_STRING_ATTR_LENGTH ) {
            throw new TGException(String.format("UTF length of String exceed the maximum string length supported for a String Attribute (%d > %d)",
                    len, MAX_STRING_ATTR_LENGTH));
        }
        this.value = s;
        setModified();
    }

    @Override
    public byte[] getAsBytes() {
        String str = (String) value;
        return str.getBytes();
    }

    @Override
    public char[] getAsChars() {
        String str = (String) value;
        CharBuffer cb = CharBuffer.wrap(str);
        return cb.array();

    }

    @Override
    public short getAsShort() {
        return 0;
    }

    @Override
    public int getAsInt() {
        return 0;
    }

    @Override
    public long getAsLong() {
        return 0;
    }

    @Override
    public float getAsFloat() {
        if (this.value == null) return Float.NaN;
        return Float.parseFloat((String.class.cast(this.value)));

    }

    @Override
    public double getAsDouble() {
        if (this.value == null) return Double.NaN;
        return Double.parseDouble((String.class.cast(this.value)));
    }

    @Override
    public String getAsString() {
        if (this.value == null) return null;
        return this.value.toString();
    }

    @Override
    public Calendar getAsDate() {
        return getAsCalendar(TGAttributeType.Date);
    }

    @Override
    public Calendar getAsTime() {
        return getAsCalendar(TGAttributeType.Time);
    }

    @Override
    public Calendar getAsTimestamp() {
        return getAsCalendar(TGAttributeType.TimeStamp);
    }

    private Calendar getAsCalendar(TGAttributeType type)
    {
        if (this.value == null) return null;
        GraphMetadataImpl gmi = (GraphMetadataImpl) owner.getGraphMetadata();
        ConnectionImpl conn = gmi.getConnection();
        TGProperties<String, String> properties = conn.getProperties();
        String format = null;
        String locale = properties.getProperty(ConfigName.ConnectionLocale);
        try {
            switch (type) {
                case Date:
                    format = properties.getProperty(ConfigName.ConnectionDateFormat);
                    break;

                case Time:
                    format = properties.getProperty(ConfigName.ConnectionTimeFormat);
                    break;

                case TimeStamp:
                    format = properties.getProperty(ConfigName.ConnectionTimeStampFormat);
                    break;

                default:
                    return null;
            }
            Calendar cal = Calendar.getInstance();
            SimpleDateFormat sdf = new SimpleDateFormat(format, Locale.forLanguageTag(locale));
            cal.setTime(sdf.parse(String.class.cast(this.value)));
            return cal;
        }
        catch (Exception e) { }
        return null;

    }

    @Override
    public BigDecimal getAsNumber() {
        if (this.value == null) return null;
        return new BigDecimal(String.class.cast(this.value));
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.value = is.readUTF();
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeUTF((String)this.value);
    }

    private int utfLength(String str) {
        int strlen = str.length();
        char c;
        int utflen = 0;
        for (int i = 0; i < strlen; i++) {
            c = str.charAt(i);
            if ((c >= 0x0001) && (c <= 0x007F)) {
                utflen++;
            } else if (c > 0x07FF) {
                utflen += 3;
            } else {
                utflen += 2;
            }
        }
        return utflen;
    }
}
