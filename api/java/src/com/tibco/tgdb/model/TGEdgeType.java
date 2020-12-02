package com.tibco.tgdb.model;

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

 * File name : TGEdgeType.java
 * Created on: 01/22/2015
 * Created by: suresh
 * SVN Id: $Id: TGEdgeType.java 3142 2019-04-26 00:15:06Z nimish $
 */

public interface TGEdgeType extends TGEntityType {

    TGEdge.DirectionType getDirectionType();
    
    TGNodeType getFromNodeType();
    
    TGNodeType getToNodeType();
    
}
