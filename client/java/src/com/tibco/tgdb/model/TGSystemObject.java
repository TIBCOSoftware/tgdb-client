package com.tibco.tgdb.model;

import com.tibco.tgdb.pdu.TGSerializable;

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
 *
 * File name : TGSystemObject.java
 * Created by: suresh
 * SVN Id: $Id: TGSystemObject.java 622 2016-03-19 20:51:12Z ssubrama $
 */

public interface TGSystemObject extends TGSerializable{

    enum TGSystemType {
    	InvalidType(-1),
    	AttributeDescriptor(0),
    	NodeType(1),
        EdgeType(2),
        Index(3),
        Prinicapl(4),
        Role(5),
        Sequence(6),
        MaxSysObjectTypes(7);

        private final int type;

    	TGSystemType(int type) {
    		this.type = type;
    	}
    	
    	public int type() {
    		return type;
    	}

    	public static TGSystemType fromValue(int type) {
        	for (TGSystemType st : TGSystemType.values()) {
            	if (type == st.type) return st;
        	}
        	return InvalidType;
    	}
    }

    /**
     * Get the system desc enum
     * @return the system desc of the object
     */
    TGSystemType getSystemType();
    
    /**
     * Get the desc name.
     * @return the name of the object
     */
    String getName();

}
