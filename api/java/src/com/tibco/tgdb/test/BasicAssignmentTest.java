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

 * File name : BasicAssignmentTest.${EXT}
 * Created on: 09/06/2017
 * Created by: suresh
 * SVN Id: $Id: BasicAssignmentTest.java 3148 2019-04-26 00:35:38Z sbangar $
 */

package com.tibco.tgdb.test;

public class BasicAssignmentTest
{

    public static void main(String[] args) {
        try {
            TEntry.assignmentTest1();
        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }

    static class TEntry {
        String left;
        String right;
        String parent;
        boolean color;
        Object value;
        Object key;

        TEntry(String left, String right, String parent, boolean color, Object key, Object value)
        {
            this.left = left;
            this.right = right;
            this.parent = parent;
            this.color = color;
            this.key = key;
            this.value = value;
        }
        public static void assignmentTest1() throws Exception
        {
            TEntry p = new TEntry("p-left", "p-right", "p-parent", false, "p-key", "p-value");
            TEntry q = new TEntry("q-left", "q-right", "q-parent", false, "q-key", "q-value");
            p.key = q.key;
            p.value = q.value;
            p = q;
        }



    }
}
