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
 * File name : AttributeDescriptorImpl.${EXT}
 * Created on: 1/23/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: AttributeDescriptorImpl.java 2344 2018-06-11 23:21:45Z ssubrama $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGAttributeDescriptor;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGSystemObject;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.util.concurrent.atomic.AtomicInteger;

public class AttributeDescriptorImpl implements TGAttributeDescriptor {

    static AtomicInteger gLocalAttributeId = new AtomicInteger(0);
    static TGLogger gLogger = TGLogManager.getInstance().getLogger();

	private String name;
	private TGAttributeType type;
	private boolean isArray;
	private int attributeId;
    private short scale;
    private short precision;

    private AttributeDescriptorImpl() {}

    //Used during meta data fetch
    AttributeDescriptorImpl(int id) {
        this.type = TGAttributeType.Invalid;
        this.attributeId = id;
        this.isArray = false;
    }

	public AttributeDescriptorImpl (String name, TGAttributeType type) {
		this.name = name;
		this.type = type;
        this.isArray = false;
		//Purposely make it to a negative number
		this.attributeId = gLocalAttributeId.decrementAndGet();
        if (type == TGAttributeType.Number) {
            scale = 5;
            precision = 20;
        }
	}

	public AttributeDescriptorImpl (String name, TGAttributeType type, boolean isArray) {
		this.name = name;
		this.type = type;
        this.isArray = isArray;
		//Purposely make it to a negative number
		this.attributeId = gLocalAttributeId.decrementAndGet();
        if (type == TGAttributeType.Number) {
            scale = 5;
            precision = 20;
        }
	}

	//TODO:  To be used when created from server side data
	public AttributeDescriptorImpl (String name, TGAttributeType type, boolean isArray, int attributeId) {
		this.name = name;
		this.type = type;
		this.isArray = isArray;
		this.attributeId = attributeId;
	}

	@Override
	public TGSystemObject.TGSystemType getSystemType() {
		return TGSystemObject.TGSystemType.AttributeDescriptor;
	}

    @Override
    public int getAttributeId() {
        return attributeId;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public TGAttributeType getType() {
        return type;
    }

    @Override
    public boolean isArray() {
        return isArray;
    }

    @Override
    public short getPrecision() { return precision; }

    @Override
    public short getScale() { return scale; }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException
    {
        try {
        	os.writeByte(TGSystemType.AttributeDescriptor.type());  // sysobject desc attribute descriptor
            os.writeInt(attributeId);
            os.writeUTF(name);
            os.writeByte(type.typeId());
            os.writeBoolean(isArray);
            if (type == TGAttributeType.Number) {
                os.writeShort(precision);
                os.writeShort(scale);
            }
        } catch (IOException ioe) {
            gLogger.log(TGLevel.Warning, "Failed to write attribute description for : %s", name);
            throw ioe;
        }
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
    	int sysObjectType = is.readByte(); // read the sysobject desc field which should be 0 for attribute descriptor
    	TGSystemType stype = TGSystemType.fromValue(sysObjectType);
    	if (stype != TGSystemType.AttributeDescriptor) {
    		gLogger.log(TGLevel.Warning, "Attribute descriptor has invalid input stream value : %d", sysObjectType);
    		//FIXME: Throw exception is needed
    	}
        this.attributeId = is.readInt();
        this.name = is.readUTF();
        this.type = TGAttributeType.fromTypeId(is.readByte());
        this.isArray = is.readBoolean();
        if (type == TGAttributeType.Number) {
            precision = is.readShort();
            scale = is.readShort();
        }
    }
    
    public void setAttributeId(int id) {
    	this.attributeId = id;
    }

    public void setPrecision(short precision)
    {
        if (type == TGAttributeType.Number) {
            this.precision = precision;
        }
    }

    public void setScale(short scale)
    {
        if (type == TGAttributeType.Number) {
            this.scale = scale;
        }
    }
}
