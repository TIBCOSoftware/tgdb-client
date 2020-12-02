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
 *  File name : JavaAdminTest.java
 *  Created on: 6/10/20
 *  Created by: sbangar@tibco.com
 *  Version   : 3.0.0
 *  Since     : 3.0.0
 *  SVN Id    : : $
 *
 */

 package com.tibco.tgdb.test;

 import com.tibco.tgdb.connection.TGConnection;

 import static java.lang.System.out;
 import com.tibco.tgdb.admin.TGAdminConnection;
 import com.tibco.tgdb.connection.TGConnectionFactory;
 import com.tibco.tgdb.exception.TGException;
 import com.tibco.tgdb.model.*;


 public class JavaAdminTest {

     String url = "tcp://127.0.0.1:8222/{dbName=demodb;}";
     String user = "admin";
     String pwd = "admin";
     protected TGAdminConnection adminConn = null;

     private void setup() throws Exception {
     	adminConn = (TGAdminConnection)TGConnectionFactory.getInstance().createAdminConnection(url, user, pwd, null);
     	adminConn.connect();
     }
     private void cleanup() throws Exception {
     	adminConn.disconnect();
         out.println("Disconnected.");
     }
     private void testCreateUser() throws Exception{
     	String username = "Sneha";
     	String password = "testadmin";
     	String[] roles = {"basicrole","operator","user","sysadm"};
     	String invalidRole = "invalidRole";

     	//create new user with basicrole
     	try {
     		adminConn.createUser(username, password,roles[0]);
     		System.out.println("Successfully created user :" + username);

     	}
     	catch(TGException e) {
     		System.out.println("Server Error code: "+ e.getServerErrorCode()+ " Error Message: "+ e.getMessage());
     	}

     	//try to create already existing user
     	try {
     		adminConn.createUser(username, password);
     		System.out.println("Successfully created user :" + username);
     	}
     	catch(TGException e) {
     		System.out.println("Server Error code: "+ e.getServerErrorCode()+ " Error Message: "+ e.getMessage());
     	}

     	username = "tibco";
     	//create user with multiple roles
     	try {
     		adminConn.createUser(username, password,roles[0],roles[2]);
     		System.out.println("Successfully created user :" + username);
     	}
     	catch(TGException e) {
     		System.out.println("Server Error code: "+ e.getServerErrorCode()+ " Error Message: "+ e.getMessage());
     	}

     	username = "seattle";
     	//create user without any roles
     	try {
     		adminConn.createUser(username, password);
     		System.out.println("Successfully created user :" + username);
     	}
     	catch(TGException e) {
     		System.out.println("Server Error code: "+ e.getServerErrorCode()+ " Error Message: "+ e.getMessage());
     	}

     	username = "graphdb";
     	//create user with invalid role
     	try {
     		adminConn.createUser(username, password,invalidRole);
     		System.out.println("Successfully created user :" + username);
     	}
     	catch(TGException e) {
     		System.out.println("Server Error code: "+ e.getServerErrorCode()+ " Error Message: "+ e.getMessage());
     	}

     }
     public static void main(String[] args) throws Exception {
     	JavaAdminTest javaAdminTest = new JavaAdminTest();
     	javaAdminTest.setup();
     	javaAdminTest.testCreateUser();
     	javaAdminTest.cleanup();
     }

 }