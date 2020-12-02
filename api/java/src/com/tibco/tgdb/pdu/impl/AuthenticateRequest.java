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
 * File name :AuthenticateRequest
 * Created on: 12/24/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: AuthenticateRequest.java 4062 2020-06-04 17:48:37Z sbangar $
 */

package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;
import com.tibco.tgdb.utils.TGConstants;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class AuthenticateRequest extends AbstractProtocolMessage {


    String clientId;
    String inboxAddr;
    String userName;
    String dbName;
    List<String> roles = null;
    byte[] password = TGConstants.EmptyByteArray;

    @Override
    public VerbId getVerbId() {
        return VerbId.AuthenticateRequest;
    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
        if ((dbName == null) || dbName.length() == 0) {
            os.writeBoolean(true);
        }
        else {
            os.writeBoolean(false); //No clientId
            os.writeUTF(dbName);
        }
        if ((clientId == null) || clientId.length() == 0) {
            os.writeBoolean(true);
        }
        else {
            os.writeBoolean(false); //No clientId
            os.writeUTF(clientId);
        }
        if ((inboxAddr == null) || inboxAddr.length() == 0) {
            os.writeBoolean(true);
        }
        else {
            os.writeBoolean(false);
            os.writeUTF(inboxAddr);
        }

        if(roles == null)
             os.writeInt(-1);//use all roles
        else{
            os.writeInt(roles.size());
            for(String roleName : roles)
                os.writeUTF(roleName);
        }

        if ((userName == null) || userName.length() == 0) {
            os.writeBoolean(true);
        }
        else {
            os.writeBoolean(false);
            os.writeUTF(userName);  //Can't be null.
        }
        os.writeBytes(password);
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
        //For Testing purpose only.
        boolean bIsClientId = is.readBoolean();
        if (!bIsClientId) {
            this.clientId = is.readUTF();
        }
        this.inboxAddr = is.readUTF();
        this.userName = is.readUTF();
        this.password = is.readBytes();
    }


    public String getClientId() { return clientId;}
    public void setClientId(String s) { this.clientId = s;}

    public String getInboxAddr() { return inboxAddr;}
    public void setInboxAddr(String s) { this.inboxAddr = s;}

    public String getUserName() { return userName;}
    public void setUserName(String s) { this.userName = s;}

    public byte[] getPassword() { return password;}
    public void setPassword(byte[] b) { this.password = b;}

    public void setDatabaseName(String dbName) { this.dbName = dbName;}
    public String getDatabaseName() { return this.dbName; }

    public void setRoles(List<String> roleList) { this.roles = roleList; }
    public List<String> getRoles() { return this.roles; }

}