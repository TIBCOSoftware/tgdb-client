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
 *
 * File name : CreateNode.${EXT}
 * Created on: 09/13/2017
 * Created by: suresh
 * SVN Id: $Id: CreateNode.java 3881 2020-04-16 22:40:29Z nimish $
 */

package com.tibco.tgdb.test;


import java.util.Calendar;
import java.util.GregorianCalendar;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class CreateNode {

    public static void jira250(String[] args) throws Exception {

        String url = "tcp://127.0.0.1:8222";
        String user = "scott";
        String pwd = "scott";
        TGConnection conn = null;
        try {
            conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
            conn.connect();

            TGGraphObjectFactory gof = conn.getGraphObjectFactory();
            if (gof == null) {
                throw new Exception("Graph object not found");
            }
            TGGraphMetadata gmd = conn.getGraphMetadata(true);
            TGNodeType basicnode = gmd.getNodeType("basicnode");
            if (basicnode == null)
                throw new Exception("Node desc not found");

            TGNode basic1 = gof.createNode(basicnode);

            String pkey = "Murray";

            basic1.setAttribute("name", pkey);
            basic1.setAttribute("networth", new java.math.BigDecimal("10"));
            basic1.setAttribute("address", "Palo Alto CA");
            conn.insertEntity(basic1);

            conn.commit();
            System.out.println("Entity created");

            basic1.setAttribute("address", null);
            basic1.setAttribute("networth", null);
            conn.updateEntity(basic1);
            conn.commit();
            System.out.println("Entity updated");

            conn.getGraphMetadata(true);
            TGKey key = gof.createCompositeKey("basicnode");

            key.setAttribute("name", pkey);
            TGEntity entity = conn.getEntity(key, null);
            if (entity != null) {
                System.out.println("Entity retrieved :");
                System.out.println("Name = " + entity.getAttribute("name").getValue());
                System.out.println("Address = " + entity.getAttribute("address").getValue());
                System.out.println("Networth = " + entity.getAttribute("networth").getValue());
            }
        } finally {
            if (conn != null)
                conn.disconnect();
        }
    }

    public static void jira289(String[] args) throws Exception {

        String url = "tcp://127.0.0.1:8222";
        String user = "scott";
        String pwd = "scott";
        TGConnection conn = null;
        try {
            conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
            conn.connect();

            TGGraphObjectFactory gof = conn.getGraphObjectFactory();
            if (gof == null) {
                throw new Exception("Graph object not found");
            }
            TGGraphMetadata gmd = conn.getGraphMetadata(true);
            TGNodeType basicnode = gmd.getNodeType("basicnode");
            if (basicnode == null)
                throw new Exception("Node desc not found");

            TGNode basic1 = gof.createNode(basicnode);

            // 998 bytes
            //String longStr = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
            // 999 bytes
            String longStr = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

            System.out.println("String Attr has " + longStr.getBytes("UTF-8").length + " bytes");

            String pkey = "Mike1";

            basic1.setAttribute("name", pkey);
            basic1.setAttribute("desc", longStr);

            conn.insertEntity(basic1);
            conn.commit();
            System.out.println("Entity created");

            TGKey key = gof.createCompositeKey("basicnode");
            key.setAttribute("name", pkey);
            TGEntity entity = conn.getEntity(key, null);
            if (entity != null) {
                System.out.println("Entity found : " + entity.getAttribute("desc").getValue());
            } else {
                System.out.println("Entity not found");
            }


            entity.setAttribute("desc", longStr);
            conn.updateEntity(entity);
            conn.commit();
            System.out.println("Entity updated");

            TGEntity entity2 = conn.getEntity(key, null);
            if (entity2 != null) {
                System.out.println("Entity found : " + entity2.getAttribute("desc").getValue());
            } else {
                System.out.println("Entity not found");
            }

        } finally {
            if (conn != null)
                conn.disconnect();
        }
    }

    public static void jira181(String[] args) throws Exception {

        String url = "tcp://127.0.0.1:8222";
        String user = "scott";
        String pwd = "scott";
        TGConnection conn = null;
        try {
            conn = TGConnectionFactory.getInstance().createConnection(url, user, pwd, null);
            conn.connect();

            TGGraphObjectFactory gof = conn.getGraphObjectFactory();
            if (gof == null) {
                throw new Exception("Graph object not found");
            }
            TGGraphMetadata gmd = conn.getGraphMetadata(true);
            TGNodeType basicnode = gmd.getNodeType("basicnode");
            if (basicnode == null)
                throw new Exception("Node desc not found");

            TGNode basic1 = gof.createNode(basicnode);
            String pkey = "Mike1";
            basic1.setAttribute("name", pkey);
            basic1.setAttribute("ratedate", new Calendar
                    .Builder()
                    .setDate(4176, 1, 1)
                    .set(Calendar.ERA,GregorianCalendar.BC)
                    .build());

            conn.insertEntity(basic1);
            conn.commit();
            System.out.println("Entity created");

            conn.getGraphMetadata(true);
            TGKey key = gof.createCompositeKey("basicnode");

            key.setAttribute("name", pkey);
            TGEntity entity = conn.getEntity(key, null);
            if (entity != null) {
                System.out.println("RateDate = " + entity.getAttribute("ratedate").getValue());
            }

        } finally {
            if (conn != null)
                conn.disconnect();
        }
    }

    public static void main(String[] args) throws Exception {
        try {
            jira289(args);
        }
        catch (Exception e) {
            e.printStackTrace();
        }
        try {
            jira181(args);
        }
        catch (Exception e)
        {
            e.printStackTrace();
        }
    }
}