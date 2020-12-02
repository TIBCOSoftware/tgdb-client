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
 * File name : BlobAttribute.${EXT}
 * Created on: 06/04/2018
 * Created by: suresh
 * SVN Id: $Id: BlobAttribute.java 3631 2019-12-11 01:12:03Z ssubrama $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.impl.GraphMetadataImpl;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.io.InputStream;
import java.math.BigDecimal;
import java.nio.ByteBuffer;
import java.util.concurrent.atomic.AtomicLong;

class BlobAttribute extends AbstractAttribute  {

    private long entityId;
    private boolean isCached;
    static AtomicLong gUniqueId = new AtomicLong(0);

    BlobAttribute(TGAttributeDescriptor desc) {
        super(desc);
        this.entityId = gUniqueId.decrementAndGet();
    }

    @Override
    public void setValue(Object value) throws TGException {
        if (value == null) 
        {
        	this.value = value;
        	setModified();
        	return;
        }

        if (value instanceof byte[]) {
            this.value = value;
        }
        else if (value instanceof String) {
            this.value = ((String)value).getBytes();
        }
        else if (value instanceof BigDecimal) {
            this.value = ConversionUtils.bigDecimal2ByteArray(BigDecimal.class.cast(value));
        }
        else if (value instanceof ByteBuffer) {
            this.value = ((ByteBuffer)value).array();
        }
        else if (value instanceof InputStream) {
            this.value = ConversionUtils.inputStream2ByteArray((InputStream) value);
        }
        else {
            throw new TGTypeCoercionNotSupported(TGAttributeType.Blob, value.getClass().getSimpleName());
            //SS:TODO - list what is supported.
        }
        setModified();
    }


    public byte[] getAsBytes() throws TGException {
        if (entityId < 0) return (byte[]) this.value;
        if (isCached) return (byte[])this.value;
        GraphMetadataImpl gmi = (GraphMetadataImpl) owner.getGraphMetadata();
        ConnectionImpl conn = gmi.getConnection();
        if (this.desc.isEncrypted()) {
            this.value = conn.decryptEntity(entityId);
        }
        else {
            this.value = conn.getLargeObjectAsBytes(entityId);
        }
        isCached = true;
        return (byte[]) this.value;

    }

    public ByteBuffer getAsByteBuffer() throws TGException
    {
        byte[] buf = getAsBytes();
        return ByteBuffer.wrap(buf);
    }

    @Override
    public void readValue(TGInputStream is) throws TGException, IOException {
        this.entityId = is.readLong();
        this.isCached = false;
    }

    @Override
    public void writeValue(TGOutputStream os) throws TGException, IOException {
        os.writeLong(this.entityId);
        if (this.value == null) {
            os.writeBoolean(false);
        }
        else {
            os.writeBoolean(true);
            os.writeBytes(byte[].class.cast(value));
        }
    }

}
