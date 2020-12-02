/** * Copyright 2019 TIBCO Software Inc. All rights reserved.  * * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except * in compliance with the License.  * A copy of the License is included in the distribution package with this file.  * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name : ResultDataDescriptor.${EXT}
 * Created on: 05/01/2016
 * Created by: chung
 * SVN Id: $Id: ResultDataDescriptor.java 3149 2019-04-26 00:45:37Z sbangar $
 */

package com.tibco.tgdb.query.impl;

import static com.tibco.tgdb.model.TGAttributeType.*;
import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.*;

import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGSystemObject;
import com.tibco.tgdb.query.TGResultDataDescriptor;


public class ResultDataDescriptor implements TGResultDataDescriptor {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private DATA_TYPE dataType = TYPE_NODE;
    private String annot = null;
    private boolean isMap = false;
    private boolean isArray = false;
    private boolean hasType = false;
    private int containedSize = 0;
    private TGAttributeType scalarType = Invalid;

    private TGSystemObject sysObject;
    private TGResultDataDescriptor keyDesc = null;
    private TGResultDataDescriptor valueDesc = null;
    private TGResultDataDescriptor[] containedDesc = null;

    ResultDataDescriptor() {
        dataType = TYPE_NODE;
    }

    ResultDataDescriptor(DATA_TYPE type) {
        dataType = type;
    }

    ResultDataDescriptor(DATA_TYPE type, char scalarAnnot) {
        dataType = type;
        hasType = true;
        switch (scalarAnnot) {
            case '?':
                scalarType = Boolean;
                break;
            case 'b':
                scalarType = Byte;
                break;
            case 'c':
                scalarType = Char;
                break;
            case 'h':
                scalarType = Short;
                break;
            case 'i':
                scalarType = Integer;
                break;
            case 'l':
                scalarType = Long;
                break;
            case 'f':
                scalarType = Float;
                break;
            case 'd':
                scalarType = Double;
                break;
            case 's':
                scalarType = String;
                break;
            case 'D':
                scalarType = Date;
                break;
            case 'e':
                scalarType = Time;
                break;
            case 'a':
                scalarType = TimeStamp;
                break;
            case 'n':
                scalarType = Number;
                break;
            default:
                hasType = false;
                break;
        }
    }

    public DATA_TYPE getDataType() {
        return dataType;
    }

    public int getContainedDataSize() {
        return containedSize;
    }

    public boolean isMap() {
        return isMap;
    }

    public boolean isArray() {
        return isArray;
    }

    public boolean hasConcreteType() {
        return hasType;
    }

    public TGAttributeType getScalarType() {
        return scalarType;
    }

    public TGSystemObject getSystemObject() {
        return sysObject;
    }

    public TGResultDataDescriptor getKeyDescriptor() {
        return keyDesc;
    }

    public TGResultDataDescriptor getValueDescriptor() {
        return valueDesc;
    }

    public TGResultDataDescriptor[] getContainedDescriptors() {
        return containedDesc;
    }

    public TGResultDataDescriptor getContainedDescriptor(int position) {
        return null;
    }

    void setContainedSize(int size) {
        containedSize = size;
    }

    void setIsArray(boolean isArray) {
        this.isArray = isArray;
    }

    void setIsMap(boolean isMap) {
        this.isMap = isMap;
    }

    void setHasConcreteType(boolean hasType) {
        this.hasType = hasType;
    }

    void setSystemObject(TGSystemObject sysObject) {
        this.sysObject = sysObject;
    }

    void setKeyDescriptor(TGResultDataDescriptor key) {
        keyDesc = key;
    }

    void setValueDescriptor(TGResultDataDescriptor value) {
        valueDesc = value;
    }

    void setContainedDescriptors(TGResultDataDescriptor[] desc) {
        containedDesc = desc;
        containedSize = desc.length;
    }
}
