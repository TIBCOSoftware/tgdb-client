/**
 * Copyright (c) 2018 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : BlobAttribute.${EXT}
 * Created on: 6/4/18
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
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
        this.value = conn.getLargeObjectAsBytes(entityId);
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
