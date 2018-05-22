/**
 * Copyright (c) 2017 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : BasicAssignmentTest.${EXT}
 * Created on: 9/6/17
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
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
