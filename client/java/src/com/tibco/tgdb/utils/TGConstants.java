package com.tibco.tgdb.utils;

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

 * File name :TGConstants
 * Created on: 1/3/15
 * Created by: suresh

 * SVN Id: $Id: TGConstants.java 622 2016-03-19 20:51:12Z ssubrama $
 */

import java.util.ArrayList;
import java.util.List;

/**
 * Global Constants
 */
public interface TGConstants {

    public static final List<?> EmptyList = new ArrayList<>();
    public static byte[] EmptyByteArray = new byte[0];
    public static String EmptyString = "";


    static final long   U64_NULL               = 0xffffffffffffffffL;
    static final byte   U64PACKED_NULL        = (byte)0xf0;

}
