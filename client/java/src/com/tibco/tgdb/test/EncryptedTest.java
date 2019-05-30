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
 * File name : EncryptedTest.${EXT}
 * Created on: 04/24/2019
 * Created by: suresh
 * SVN Id: $Id: EncryptedTest.java 3148 2019-04-26 00:35:38Z sbangar $
 */

package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.model.*;
import com.tibco.tgdb.query.TGQueryOption;

import java.math.BigDecimal;
import java.util.Calendar;
import java.util.TimeZone;


public class EncryptedTest {
    public String url = "tcp://admin@localhost:8222";
    public String passwd = "admin";
    TGGraphObjectFactory gof;
    TGGraphMetadata gmd;
    TGConnection conn = null;

    TGNodeType encryptedNode, unencryptedNode;

    public void connect() throws TGException {
        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);
        conn.connect();

        gof = conn.getGraphObjectFactory();
        if (gof == null) throw new TGException("gof is null");
        gmd = conn.getGraphMetadata(true);
        if (gmd == null) throw new TGException("gmd is null");

        encryptedNode = gmd.getNodeType("encryptedNode");
        if (encryptedNode == null) throw new TGException("Node desc encryptedNode not found");


        System.out.println("Leaving connect()...");
    }

    public void run() throws Exception {
        System.out.println("Run test.");
        TGNode node = gof.createNode(encryptedNode);

        node.setAttribute("stringPkey", "ThisIsAPkey");
        node.setAttribute("boolAttr", false);
//        node.setAttribute("byteAttr", (byte) 0xba);
        node.setAttribute("charAttr", '*');
        node.setAttribute("shortAttr", (short) 6385);
        node.setAttribute("intAttr", 73741825);
        node.setAttribute("longAttr", (long) 1342177281);
        node.setAttribute("floatAttr", (float) 2.23);
        node.setAttribute("doubleAttr", 2336.32424);
        node.setAttribute("numberAttr", new BigDecimal("234235732234235723590735124523.89813275891735070"));
        node.setAttribute("stringAttr", "betterStringKey");
        node.setAttribute("dateAttr", new Calendar
                .Builder()
                .setDate(2016, 10, 31)
                .build());
        node.setAttribute("timeAttr", new Calendar
                .Builder()
                .setTimeOfDay(21, 32, 12, 845)
                .setTimeZone(TimeZone.getDefault())
                .build());
        node.setAttribute("timestampAttr", new Calendar
                .Builder()
                .setDate(2016, 10, 25)
                .setTimeOfDay(8,9,30,999)
                .build());

        conn.insertEntity(node);
        conn.commit();

        TGKey key = gof.createCompositeKey("encryptedNode");
        key.setAttribute("stringPkey", "ThisIsAPkey");
        TGQueryOption option = TGQueryOption.createQueryOption();

        TGEntity entity = conn.getEntity(key, option);

        if (entity instanceof TGNode) {
            System.out.println("Found node. stringPkey = " + entity.getAttribute("stringPkey").getAsString());
            System.out.println("boolAttr = " + entity.getAttribute("boolAttr").getAsString());
//            System.out.println("byteAttr = " + entity.getAttribute("byteAttr").getAsString());
            System.out.println("charAttr = " + entity.getAttribute("charAttr").getAsString());
            System.out.println("shortAttr = " + entity.getAttribute("shortAttr").getAsString());
            System.out.println("intAttr = " + entity.getAttribute("intAttr").getAsString());
            System.out.println("longAttr = " + entity.getAttribute("longAttr").getAsString());
            System.out.println("floatAttr = " + entity.getAttribute("floatAttr").getAsString());
            System.out.println("doubleAttr = " + entity.getAttribute("doubleAttr").getAsString());
            System.out.println("numberAttr = " + entity.getAttribute("numberAttr").getAsString());
            System.out.println("stringAttr = " + entity.getAttribute("stringAttr").getAsString());
            System.out.println("dateAttr = " + entity.getAttribute("dateAttr").getAsString());
            System.out.println("timeAttr = " + entity.getAttribute("timeAttr").getAsString());
            System.out.println("timestampAttr = " + entity.getAttribute("timestampAttr").getAsString());
        }
    }

    public static void main(String[] args) throws Exception {
        try {
            EncryptedTest test = new EncryptedTest();
            test.connect();
            test.run();
        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }
}
