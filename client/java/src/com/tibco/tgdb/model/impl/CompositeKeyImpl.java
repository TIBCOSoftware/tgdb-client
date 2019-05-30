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
 * <p/>
 * File name : CompositeKeyImpl.${EXT}
 * Created on: 3/19/16
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: CompositeKeyImpl.java 3145 2019-04-26 00:26:51Z nimish $
 */


package com.tibco.tgdb.model.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.model.impl.attribute.AbstractAttribute;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;

import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.Map;

public class CompositeKeyImpl implements TGKey {

    private TGGraphMetadata graphMetadata;
    private String typeName;
    Map<String, TGAttribute> attributes = new LinkedHashMap<String, TGAttribute>();

    //FIXME: Not sure desc name is needed or not
    public CompositeKeyImpl(TGGraphMetadata graphMetadata, String typeName) throws TGException {
        this.graphMetadata = graphMetadata;
        this.typeName = typeName;

        /* Not require to have a desc name
        if (typeName != null) return;

        if (graphMetadata.getNodeType(typeName) == null)
            throw new TGException(String.format("Invalid NodeType specified :%s", typeName));
        */
    }

    @Override
    public void setAttribute(String name, Object value) throws TGException {

        TGAttribute attr = null;

        if (value == null) throw new TGException(String.format("Value is null"));

        attr = attributes.get(name);
        if (attr == null) {
            TGAttributeDescriptor attrDesc = null;
            attrDesc = graphMetadata.getAttributeDescriptor(name);
            if (attrDesc == null) {
                attrDesc = ((GraphMetadataImpl)graphMetadata).createAttributeDescriptor(name, value.getClass());
            }
            attr = AbstractAttribute.createAttribute(null, attrDesc, value);
        }
        try {
            attr.setValue(value);
            attributes.put(name, attr);
        }
        catch (Exception e) {
            throw TGException.buildException("Can't set value to attribute", null, e);
        }
    }

    @Override
    public void writeExternal(TGOutputStream os) throws TGException, IOException {
        if (this.typeName != null) {
            os.writeBoolean(true); //TypeName exists
            os.writeUTF(this.typeName);
        }
        else {
            os.writeBoolean(false);
        }
        os.writeShort(attributes.size());
        for (TGAttribute attr : attributes.values()) {
        	//Null value is not allowed and therefore no need to include isNull flag
            attr.writeExternal(os);
        }
    }

    @Override
    public void readExternal(TGInputStream is) throws TGException, IOException {
        throw new TGException("Not Supported operation");
    }
}
