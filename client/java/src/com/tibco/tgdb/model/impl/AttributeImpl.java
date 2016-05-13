/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * <p/>
 * File name : AttributeImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: AttributeImpl.java 771 2016-05-05 11:40:52Z vchung $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;

public class AttributeImpl implements TGAttribute {

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    AbstractEntity owner;

    TGAttributeDescriptor type;
    //addEdge does not trigger change notification
    Object value;
    boolean isModified = false;

    public AttributeImpl(AbstractEntity owner) {
        this.owner = owner;
    }

    public AttributeImpl(AbstractEntity owner, TGAttributeDescriptor type, Object value) {
        this.owner = owner;
        this.value = value;
        this.type = type;
    }

    @Override
    public TGEntity getOwner() {
        return owner;
    }

    //FIXME:  Need to change this name to getAttributeDescriptor
    @Override
    public TGAttributeDescriptor getAttributeType() {
        return type;
    }

    @Override
    public boolean isNull() {
        return value==null;
    }

    @Override
    public Object getValue() {
        return value;
    }

    @Override
    public void setValue(Object value) throws TGException {
    	//FIXME: Need to match the type of the attribute descriptor
        this.value = value;
        isModified = true;
    }

    @Override
    public boolean getAsBoolean() {
        return Boolean.class.cast(value);
    }

    @Override
    public byte getAsByte() {
        return Byte.class.cast(value);
    }

    @Override
    public char getAsChar() {
        return Character.class.cast(value);
    }

    @Override
    public short getAsShort() {
        return Short.class.cast(value);
    }

    @Override
    public int getAsInt() {
        return Integer.class.cast(value);
    }

    @Override
    public long getAsLong() {
        return Long.class.cast(value);
    }

    @Override
    public float getAsFloat() {
        return Float.class.cast(value);
    }

    @Override
    public double getAsDouble() {
        return Double.class.cast(value);
    }

    @Override
    public String getAsString() {
        return String.class.cast(value);
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
       	int aid = type.getAttributeId();
        //null attribute is not allowed during entity creation
        os.writeInt(aid);
       	os.writeBoolean(isNull());
       	if (isNull()) {
       		return;
       	}
    	switch(type.getType()) {
    		case Boolean:
    			os.writeBoolean(Boolean.class.cast(value));
    			break;
    		case Byte:
    			os.writeByte(Byte.class.cast(value));
    			break;
    		case Char:
    			os.writeChar(Character.class.cast(value));
    			break;
    		case Short:
    			os.writeShort(Short.class.cast(value));
    			break;
    		case Int:
    			os.writeInt(Integer.class.cast(value));
    			break;
    		case Long:
    			os.writeLong(Long.class.cast(value));
    			break;
    		case Float:
    			os.writeFloat(Float.class.cast(value));
    			break;
    		case Double:
    			os.writeDouble(Double.class.cast(value));
    			break;
    		case String:
    			os.writeUTF(String.class.cast(value));
    			break;
    		default:
    			break;
    	}
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	int aid = is.readInt();
        TGAttributeDescriptor at = ((GraphMetadataImpl) owner.graphMetadata).getAttributeDescriptor(aid);
        this.type = at;
        if (at == null) {
        	//FIXME: retrieve entity type together with the entity?
        	gLogger.log(TGLevel.Warning, "cannot lookup attribute descriptor %d from graph meta data cache", aid);
        }
        if (is.readByte() == 1) {
        	value = null;
        	return;
        }
    	switch(type.getType()) {
    		case Boolean:
    			byte b = is.readByte();
    			value = Boolean.valueOf(b != 0);
    			break;
    		case Byte:
    			value = is.readByte();
    			break;
    		case Char:
    			value = is.readChar();
    			break;
    		case Short:
    			value = is.readShort();
    			break;
    		case Int:
    			value = is.readInt();
    			break;
    		case Long:
    			value = is.readLong();
    			break;
    		case Float:
    			value = is.readFloat();
    			break;
    		case Double:
    			value = is.readDouble();
    			break;
    		case String:
    			value = is.readUTF();
    			break;
    		default:
    			break;
    	}
    }
    
    void resetIsModified() {
    	this.isModified = false;
    }
    
    @Override
    public boolean isModified() {
    	return isModified;
    }
}
