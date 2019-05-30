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

 * File name : BlobUnitTest.${EXT}
 * Created on: 10/17/2018
 * Created by: suresh
 * SVN Id: $Id: BlobUnitTest.java 3148 2019-04-26 00:35:38Z sbangar $
 */

package com.tibco.tgdb.test;

import com.tibco.tgdb.model.TGKey;
import com.tibco.tgdb.model.TGNode;

import java.io.FileInputStream;
import java.io.FileOutputStream;

public class BlobUnitTest extends ConnectionUnitTest{

    public BlobUnitTest(String[] args) {
        super(args);
    }

    protected void testcase1() throws Exception {
/*
        FileInputStream img = new FileInputStream("/Users/suresh/Desktop/GraphDb/Roadmap_Presentations/TIBCO_community_graphdatabase_hero_370x275.png");
        FileInputStream doc = new FileInputStream("/Users/suresh/Desktop/GraphDb/Roadmap_Presentations/Graph-DB-Product-Page.docx");
        TGNode n1 = gof.createNode(basicNodeType);
        n1.setAttribute("name", "john");
        n1.setAttribute("image", img);
        n1.setAttribute("doc", doc);
        conn.insertEntity(n1);
        conn.commit();
*/
        TGKey key = gof.createCompositeKey("basicnode");
        key.setAttribute("name", "john");
        TGNode n2 = (TGNode) conn.getEntity(key, null);
        byte[] buf = n2.getAttribute("image").getAsBytes();
        FileOutputStream fos = new FileOutputStream("/Users/suresh/Desktop/GraphDb/Roadmap_Presentations/Retrievd_TIBCO_community_graphdatabase_hero_370x275.png");
        fos.write(buf);
        fos.flush();
        fos.close();
//        TGBlob blob1 = n2.getAttribute("image").getAsBlob();
//        blob1.free();
    }

    public static void main(String[] args)  {
        try {
            BlobUnitTest blobUnitTest = new BlobUnitTest(args);
            blobUnitTest.connect();
            blobUnitTest.testcase1();
        }
        catch (Exception e) {
            e.printStackTrace();;
        }
    }
}
