package com.tibco.tgdb.test;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.query.TGQuery;
import com.tibco.tgdb.query.TGQueryOption;

public class QueryTest {

    public static void main(String[] args) throws Exception {
        String url = "tcp://scott@localhost:8222";
        String passwd = "scott";

        TGConnection conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();
        System.out.println("Start test query");
        TGQuery Query1 = conn.createQuery("testquery < X '5ef';");
        Query1.execute();
        conn.executeQuery("testquery < X '5ef';", TGQueryOption.DEFAULT_QUERY_OPTION);
        Query1.close();
        conn.executeQuery("testquery < X '5ef';", TGQueryOption.DEFAULT_QUERY_OPTION);

        conn.disconnect();
        System.out.println("Disconnected.");
    }
}