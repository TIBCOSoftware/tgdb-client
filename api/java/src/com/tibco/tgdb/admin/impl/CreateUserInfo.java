/**
 * Copyright 2020 TIBCO Software Inc. All rights reserved.
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
 *  File name : CreateUserInfo.java
 *  Created on: 6/04/20
 *  Created by: sbangar@tibco.com
 *  Version   : 3.0.0
 *  Since     : 3.0.0
 *  SVN Id    : : $
 *
 */

package com.tibco.tgdb.admin.impl;

import java.util.ArrayList;
import java.util.List;


public class CreateUserInfo {
	protected String name;
    protected String passwd;
    protected List<String> roles = new ArrayList<String>();

    public CreateUserInfo(String _name,String _passwd, String ..._roles) {
        this.name = _name;
        this.passwd = _passwd;
        if(_roles != null) {
            for(String roleName :_roles)
                roles.add(roleName);
        }
    }

    public String getName() {
        return name;
    }
    public String getPasswd() {
        return passwd;
    }
    public List<String> getRoles(){
        return roles;
    }
}