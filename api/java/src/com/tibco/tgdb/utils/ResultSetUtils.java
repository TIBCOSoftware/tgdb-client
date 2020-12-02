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

 * File name : ResultSetUtils.java
 * Created on: 04/19/2020
 * Created by: chung
 * SVN Id: $Id: ResultSetUtils.java 3137 2019-04-25 23:52:32Z sbangar $
 */

package com.tibco.tgdb.utils;

import com.tibco.tgdb.query.TGResultDataDescriptor;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGResultSetMetaData;

import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.*;
import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.TYPE_PATH;

public class ResultSetUtils {
    public static void printRSMetaData(TGResultSet rs) {
        TGResultSetMetaData rsmd = rs.getMetaData();
        if (rsmd == null) {
            System.out.println("Result has no meta data");
            return;
        }
        TGResultDataDescriptor.DATA_TYPE dt = rsmd.getResultType();
        System.out.printf("Result data type : %s\n", dt);
        TGResultDataDescriptor rdd = rsmd.getResultDataDescriptor();
        printMetaData(rdd, "");
    }

    private static void printMetaData(TGResultDataDescriptor rdd, String ident) {
        System.out.printf("%sType : %s, Scalar type : %s, isArray : %b, isMap : %b\n",
                ident,
                rdd.getDataType(), rdd.getDataType() == TYPE_SCALAR ? rdd.getScalarType() : "n\\a",
                rdd.isArray(), rdd.isMap());
        if (rdd.getDataType() == TYPE_MAP) {
            System.out.printf("%skey type : %s, key scalar type : %s, value type : %s, value scalar type : %s\n",
                    ident,
                    rdd.getKeyDescriptor().getDataType(), rdd.getKeyDescriptor().getScalarType(),
                    rdd.getValueDescriptor().getDataType(), rdd.getValueDescriptor().getScalarType());
        }
        if (rdd.getDataType() == TYPE_LIST) {
            TGResultDataDescriptor[] ndd = rdd.getContainedDescriptors();
            for (int i=0; i<ndd.length; i++) {
                printMetaData(ndd[i], ident + "   ");
            }
        }
        if (rdd.getDataType() == TYPE_PATH) {
            TGResultDataDescriptor[] ndd = rdd.getContainedDescriptors();
            for (int i=0; i<ndd.length; i++) {
                printMetaData(ndd[i], ident + "   ");
            }
        }
    }
}

