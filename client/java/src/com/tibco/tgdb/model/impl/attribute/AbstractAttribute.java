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
 * File name : AbstractAttribute.${EXT}
 * Created on: 06/03/2018
 * Created by: suresh (suresh.subramani@tibco.com)
 * SVN Id: $Id: AbstractAttribute.java 3154 2019-04-26 18:31:55Z sbangar $
 */

package com.tibco.tgdb.model.impl.attribute;

import com.tibco.tgdb.connection.impl.ConnectionImpl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTypeCoercionNotSupported;
import com.tibco.tgdb.exception.TGTypeNotSupported;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.model.impl.AbstractEntity;
import com.tibco.tgdb.model.impl.GraphMetadataImpl;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.math.BigDecimal;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.util.Calendar;
import java.util.concurrent.atomic.AtomicLong;

public abstract class AbstractAttribute implements TGAttribute {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();
    static AtomicLong gUniqueId = new AtomicLong(0);

    protected AbstractEntity owner = null;
    protected TGAttributeDescriptor desc = null;
    protected boolean isModified = false;
    Object value;

    protected AbstractAttribute(TGAttributeDescriptor desc) {this.desc = desc;}
    protected AbstractAttribute(AbstractEntity owner) { this.owner = owner;}

     @Override
     public TGEntity getOwner() {
         return owner;
     }

     @Override
     public TGAttributeDescriptor getAttributeType() {
         return desc;
     }

     @Override
     public TGAttributeDescriptor getAttributeDescriptor() {
         return desc;
     }

     @Override
     public boolean isNull() {
         return value==null;
     }

     protected void setNull() {
        value = null;
        setModified();
     }

     @Override
     public Object getValue() {
         return value;
     }

     @Override
     public boolean isModified() {
         return isModified;
     }

     public void setModified() {
        this.isModified = true;
     }

     public void resetIsModified() {
         this.isModified = false;
     }

    public void writeExternal(TGOutputStream os) throws TGException, IOException {
        int aid = desc.getAttributeId();
        //null attribute is not allowed during entity creation
        os.writeInt(aid);
        os.writeBoolean(isNull());
        if (isNull()) {
            return;
        }
        if (desc.isEncrypted()) {
            writeEncrypted(os);
            return;
        }
        writeValue(os);
        return;
    }
    public void readExternal(TGInputStream is) throws TGException, IOException {
        //We have already read the AttributeId, so no need to read it.
        boolean isNull = is.readBoolean();
        if (isNull) {
            this.value = null;
            return;
        }

        if ((desc.isEncrypted()) &&
            (desc.getType() != TGAttributeType.Blob) &&
            (desc.getType() != TGAttributeType.Clob))
        {
            readDecrypted(is);
            return;
        }
        readValue(is);
    }

    public abstract void readValue(TGInputStream is) throws TGException, IOException;

    public abstract void writeValue(TGOutputStream os) throws TGException, IOException;



    public static <K extends AbstractAttribute> K readExternal(TGEntity owner, TGInputStream is) throws TGException, IOException {
        int aid = is.readInt();
        GraphMetadataImpl graphMetadata = (GraphMetadataImpl)((AbstractEntity) owner).getGraphMetadata();
        TGAttributeDescriptor at = graphMetadata.getAttributeDescriptor(aid);
        if (at == null) throw TGException.buildException(
                String.format("Invalid attributeId:%d encountered while deserialized", aid),
                "TGDB-CLIENT-READEXTERNAL", null);

        K aa = AbstractAttribute.createAttribute((AbstractEntity)owner, at, null);
        aa.readExternal(is);
        return aa;
    }
    
    public static <K extends AbstractAttribute> K createAttribute(AbstractEntity owner, 
                                                                  TGAttributeDescriptor desc, 
                                                                  Object value) throws TGException
    {
        TGAttributeType type = desc.getType();
        AbstractAttribute aa = null;

        switch (type) {
            case Boolean:
                aa = new BooleanAttribute(desc);
                break;

            case Byte:
                aa = new ByteAttribute(desc);
                break;

            case Char:
                aa = new CharAttribute(desc);
                break;

            case Short:
                aa = new ShortAttribute(desc);
                break;

            case Integer:
                aa = new IntegerAttribute(desc);
                break;

            case Long:
                aa = new LongAttribute(desc);
                break;

            case Float:
                aa = new FloatAttribute(desc);
                break;

            case Double:
                aa = new DoubleAttribute(desc);
                break;

            case Number:
                aa = new NumberAttribute(desc);
                break;

            case String:
                aa = new StringAttribute(desc);
                break;

            case Date:
            case Time:
            case TimeStamp:
                aa = new TimestampAttribute(desc);
                break;

            case Clob:
                aa = new ClobAttribute(desc);
                break;

            case Blob:
                aa = new BlobAttribute(desc);
                break;

            default:
                throw new TGTypeNotSupported(type);

        }

        aa.owner = owner;
        if (value != null)
        {
            aa.setValue(value);
        }
        
        return (K) aa;
    }

    protected void writeEncrypted(TGOutputStream os) throws TGException, IOException
    {
        byte[] blob = ConversionUtils.toByteArray(this.value, this.desc.getType());
        GraphMetadataImpl gmi = (GraphMetadataImpl) owner.getGraphMetadata();
        ConnectionImpl conn = gmi.getConnection();
        byte[] encrypted = conn.encryptEntity(blob);
        os.writeBytes(encrypted);
        return;
    }

    protected void readDecrypted(TGInputStream in) throws TGException, IOException
    {
        GraphMetadataImpl gmi = (GraphMetadataImpl) owner.getGraphMetadata();
        ConnectionImpl conn = gmi.getConnection();
        byte[] blob = conn.decryptBuffer(in);
        this.value = ConversionUtils.fromByteArray(blob, this.desc.getType());
        return;

        /*
        TGSystemObject.TGSystemType systemType = owner.getEntityType().getSystemType();
        switch (systemType) {
            case NodeType:
            {
                long entityId = in.readLong();
                blob = conn.decryptEntity(entityId);
                break;
            }
            case EdgeType:
            {
                byte[] encryptedBuf = in.readBytes();
                blob = conn.decryptBuffer(encryptedBuf);
                break;
            }
            default:
                throw new TGException(String.format("Decryption not supported for system types:%s", systemType.name()));
        }
        this.value = ConversionUtils.fromByteArray(blob, this.desc.getType());
        return;
        */
    }




    public static void main(String[] args) throws TGException {
        //BooleanAttribute battr = createAttribute(null);
    }

 }
