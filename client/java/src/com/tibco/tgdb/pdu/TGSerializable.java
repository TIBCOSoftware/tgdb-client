package com.tibco.tgdb.pdu;

import com.tibco.tgdb.exception.TGException;

import java.io.IOException;


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
 * File name :TGSerializable
 * Created on: 1/31/15
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGSerializable.java 583 2016-03-15 02:02:39Z vchung $
 */
public interface TGSerializable {

    void writeExternal(TGOutputStream os) throws TGException, IOException;

    void readExternal(TGInputStream is) throws TGException, IOException;
}
