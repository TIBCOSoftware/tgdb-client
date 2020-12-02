/** * Copyright 2019 TIBCO Software Inc. All rights reserved.  * * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except * in compliance with the License.  * A copy of the License is included in the distribution package with this file.  * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name : ResultSetMetaData.${EXT}
 * Created on: 05/01/2016
 * Created by: chung
 * SVN Id: $Id: ResultSetMetaData.java 3149 2019-04-26 00:45:37Z sbangar $
 */

package com.tibco.tgdb.query.impl;

import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.query.TGResultSetMetaData;
import com.tibco.tgdb.query.TGResultDataDescriptor;
import com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE;


import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.*;

public class ResultSetMetaData implements TGResultSetMetaData {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private DATA_TYPE resultType = TYPE_UNKNOWN;
    private String annot = null;
    private TGResultDataDescriptor ddesc = null;

    ResultSetMetaData() {
        resultType = TYPE_UNKNOWN;
    }

    ResultSetMetaData(String annot) {
        resultType = TYPE_UNKNOWN;
        this.annot = annot;
    }

    void initialize(TGGraphMetadata gmd) {
        ddesc = constructDataDescriptor(gmd, annot);
        resultType = ddesc.getDataType();
    }

    //How to avoid RESULT_TYPE prefix??
    public DATA_TYPE getResultType() {
        return resultType;
    }

    public TGResultDataDescriptor getResultDataDescriptor() {
        return ddesc;
    }

    //FIXME: Work in progress.
    private TGResultDataDescriptor constructDataDescriptor(TGGraphMetadata gmd, String annot) {
        TGResultDataDescriptor desc = null;
        int length = annot.length();
        if (length == 0) {
            return null;
        }
        char[] chars = annot.toCharArray();
        switch(chars[0]) {
            case 'V':
                desc = new ResultDataDescriptor(TYPE_NODE);
                break;
            case 'E':
                desc = new ResultDataDescriptor(TYPE_EDGE);
                break;
            case 'P':
                desc = constructPathDataDescriptor(gmd, annot);
                break;
            case 'S':
                desc = new ResultDataDescriptor(TYPE_SCALAR);
                break;
            case '[':
                desc = new ResultDataDescriptor(TYPE_LIST);
                TGResultDataDescriptor nDesc = constructDataDescriptor(gmd, annot.substring(1));
                TGResultDataDescriptor[] aDesc = new TGResultDataDescriptor[1];
                aDesc[0] = nDesc;
                ((ResultDataDescriptor) desc).setContainedDescriptors(aDesc);
                ((ResultDataDescriptor) desc).setIsArray(true);
                break;
            case '{':
                desc = new ResultDataDescriptor(TYPE_MAP);
                TGResultDataDescriptor keyDesc = new ResultDataDescriptor(TYPE_SCALAR, 's');
                TGResultDataDescriptor valueDesc = new ResultDataDescriptor(TYPE_SCALAR);
                ((ResultDataDescriptor) desc).setKeyDescriptor(keyDesc);
                ((ResultDataDescriptor) desc).setValueDescriptor(valueDesc);
                ((ResultDataDescriptor) desc).setIsMap(true);
                break;
            case '(':
                desc = new ResultDataDescriptor(TYPE_TUPLE);
                break;
            case '?':
            case 'b':
            case 'c':
            case 'h':
            case 'i':
            case 'l':
            case 'f':
            case 'd':
            case 's':
            case 'D':
            case 'e':
            case 'a':
            case 'n':
                desc = new ResultDataDescriptor(TYPE_SCALAR, chars[0]);
                break;
            default:
                gLogger.log(TGLogger.TGLevel.Info, "Failed to initialize data descriptor with type annotation : %c(%s)", chars[0], annot);
                break;
        }
        return desc;
    }

    //FIXME:  Need to be fortified
    TGResultDataDescriptor constructPathDataDescriptor(TGGraphMetadata gmd, String annot) {
        TGResultDataDescriptor desc = new ResultDataDescriptor(TYPE_PATH);
        String[] data = annot.substring(2, annot.length() - 1).split(",");
        TGResultDataDescriptor[] edesc = new TGResultDataDescriptor[data.length];
        for (int i=0; i<data.length; i++) {
            edesc[i] = constructDataDescriptor(gmd, data[i]);
        }
        ((ResultDataDescriptor) desc).setContainedDescriptors(edesc);
        return desc;
    }

    public String getAnnot() {
        return annot;
    }
}
