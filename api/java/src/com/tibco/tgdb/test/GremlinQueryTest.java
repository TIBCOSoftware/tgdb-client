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
 * File name : GremlinQueryTest.${EXT}
 * Created on: 11/07/2018
 * Created by: vincent
 * SVN Id: $Id: GremlinQueryTest.java 4651 2020-11-04 22:11:46Z vchung $
 */

package com.tibco.tgdb.test;

import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.*;
import static com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE.TYPE_PATH;
import static java.lang.System.out;
import static com.tibco.tgdb.test.GremlinQueryTest.Strategy.*;

import java.util.*;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.gremlin.DefaultGraphTraversal;
import com.tibco.tgdb.model.TGAttributeType;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.query.TGResultDataDescriptor;
import com.tibco.tgdb.query.TGResultSetMetaData;
import com.tibco.tgdb.query.impl.ResultSetMetaData;
import com.tibco.tgdb.utils.ResultSetUtils;
import org.apache.tinkerpop.gremlin.process.traversal.Order;
import org.apache.tinkerpop.gremlin.process.traversal.P;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.structure.Edge;
import org.apache.tinkerpop.gremlin.structure.util.empty.EmptyGraph;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.gremlin.GraphTraversalSource;
import com.tibco.tgdb.gremlin.__;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGAttribute;
import com.tibco.tgdb.query.TGResultSet;
import com.tibco.tgdb.query.TGResultDataDescriptor.DATA_TYPE;

//This test uses data from OpenFlight and Caliper data
public class GremlinQueryTest {
    private String url = "tcp://scott@localhost:8222/{dbname=gqt}";
    private String passwd = "scott";
    private TGLogger logger;
    private TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
    private boolean testAll = true;
    private boolean testFlight = true;
    private boolean testCaliper = true;
    private boolean showResultMeta = false;
    private boolean testValue = false;
    private boolean testCondition = false;
    private boolean testSpecialStep = false;
    private boolean testTraversal = false;
    private boolean testNegCase = false;
    private boolean testAggr = false;
    private boolean testEdge = false;
    private boolean testException = false;
    private boolean testRSMeta = false;
    private boolean testGroup = false;
    private boolean testRepeat = false;
    private boolean testSP = false;
    private boolean testByteCode = false;
    private boolean testDedup = false;
    private boolean testOrder = false;
    private boolean testOrderDedup = false;
    private boolean testVids = false;
    private boolean testACL = false;
    private GraphTraversalSource g;
    private TGConnection conn;
    private Validator listValueValidator;
    private Validator resultSizeValidator;
    private Validator fixedSizeValidator;
    private Validator mapValueValidator;
    private Validator resultAttrSizeValidator;
    private Processor resultListProcessor;
    private Processor countInResultProcessor;
    private Processor mapInResultProcessor;
    private Processor listInResultProcessor;
    private Processor attrInResultProcessor;
    private Processor entityPathInResultProcessor;
    private Processor resultCountProcessor;
    //private StringExtractor<TGEntity> entityTypeExtractor;

    enum Strategy {
        Invalid,
        VerifyListValue,
        VerifyListSize,
    };

    class ResultType {
        DATA_TYPE type;
        TGAttributeType scalarType;
        DATA_TYPE keyType;
        TGAttributeType keyScalarType;
        DATA_TYPE valueType;
        TGAttributeType valueScalarType;
        ResultType[] containedType;
        String annot;
    }

    //Whatever info needed for validator and processor
    class TestContext {
        List resultList;
        Object expectedValue;
        Object returnedValue;
        Validator validator;
        Processor processor;
        Strategy strategy;

        TestContext() {

        }

        TestContext(Object expectedValue) {
            this.expectedValue = expectedValue;
        }

        TestContext(List resultList, Object expectedValue) {
            this.resultList = resultList;
            this.expectedValue = expectedValue;
        }
    }

    class NVPair<String,V> {
        String name;
        V value;

        NVPair(String name, V value) {
            this.name = name;
            this.value = value;
        }
    }

    /*
    Similar to toString method.
    An entity can have multiple ways to show its information
    when display the query results
    private interface StringExtractor<E> {
    String extract(E e);
    }
    */

    //Too many lambdas in a file causes Intellij overload
    /*
    Using lambda to execution a gremlin query
    private interface Query {
        List execute();
    }
    */

    private interface Validator {
        //boolean validate(List valueList, List<Integer> retValues, List<Integer> expectedValues);
        boolean validate(TestContext ctx);
    }

    private interface Processor {
        boolean process(TestContext ctx);
    }

    private String getStringValue(Iterator<String> argIter) {
        if (argIter.hasNext()) {
            String s = argIter.next();
            return s;
        }
        return null;
    }

    private String getStringValue(Iterator<String> argIter, String defaultValue) {
        String s = getStringValue(argIter);
        if (s == null) {
            return defaultValue;
        } else {
            return s;
        }
    }

    private int getIntValue(Iterator<String> argIter, int defaultValue) {
        String s = getStringValue(argIter);
        if (s != null) {
            try {
                int i = Integer.parseInt(s);
                return i;
            } catch (NumberFormatException e) {
                out.printf("Invalid number : %s\n", s);
            }
        }
        return defaultValue;
    }

    private boolean getBoolValue(Iterator<String> argIter, boolean defaultValue) {
        String s = getStringValue(argIter);
        if (s != null) {
            boolean b = Boolean.parseBoolean(s);
            return b;
        }
        return defaultValue;
    }

    private void getArgs(String[] args) {
        List<String> argList = Arrays.asList(args);
        Iterator<String> argIter = argList.iterator();
        while (argIter.hasNext()) {
            String s = argIter.next();
            out.printf("Arg : \"%s\"\n", s);
            if (s.equalsIgnoreCase("-url")) {
                url = getStringValue(argIter, "tcp://scott@localhost:8222/{dbname=gqt}");
            } else if (s.equalsIgnoreCase("-password") || s.equalsIgnoreCase("-pw")) {
                passwd = getStringValue(argIter, "scott");
            } else if (s.equalsIgnoreCase("-loglevel") || s.equalsIgnoreCase("-ll")) {
                String ll = getStringValue(argIter, "Debug");
                try {
                    logLevel = TGLogger.TGLevel.valueOf(ll);
                } catch (IllegalArgumentException e) {
                    out.printf("Invalid log level value '%s'...ignored\n", ll);
                }
            } else if (s.equalsIgnoreCase("-flightonly") || s.equalsIgnoreCase("-fo")) {
                //testFlight = getBoolValue(argIter, true);
                testFlight = true;
            } else if (s.equalsIgnoreCase("-caliperonly") || s.equalsIgnoreCase("-co")) {
                testCaliper = true;
            } else if (s.equalsIgnoreCase("-showresultmeta") || s.equalsIgnoreCase("-srm")) {
                showResultMeta = true;
            } else if (s.equalsIgnoreCase("-testall") || s.equalsIgnoreCase("-ta")) {
                testAll = getBoolValue(argIter, true);
            } else if (s.equalsIgnoreCase("-testvalue") || s.equalsIgnoreCase("-tv")) {
                 testValue = true;
            } else if (s.equalsIgnoreCase("-testcond") || s.equalsIgnoreCase("-tc")) {
                testCondition = true;
            } else if (s.equalsIgnoreCase("-testspec") || s.equalsIgnoreCase("-tspc")) {
                testSpecialStep = true;
            } else if (s.equalsIgnoreCase("-testtrv") || s.equalsIgnoreCase("-ttv")) {
                testTraversal = true;
            } else if (s.equalsIgnoreCase("-testneg") || s.equalsIgnoreCase("-tneg")) {
                testNegCase = true;
            } else if (s.equalsIgnoreCase("-testaggr") || s.equalsIgnoreCase("-tagg")) {
                testAggr = true;
            } else if (s.equalsIgnoreCase("-testedge") || s.equalsIgnoreCase("-te")) {
                testEdge = true;
            } else if (s.equalsIgnoreCase("-testexcp") || s.equalsIgnoreCase("-tex")) {
                testException = true;
            } else if (s.equalsIgnoreCase("-testrsmeta") || s.equalsIgnoreCase("-trsm")) {
                testRSMeta = true;
            } else if (s.equalsIgnoreCase("-testgrp") || s.equalsIgnoreCase("-tg")) {
                testGroup = true;
            } else if (s.equalsIgnoreCase("-testrep") || s.equalsIgnoreCase("-tr")) {
                testRepeat = true;
            } else if (s.equalsIgnoreCase("-testsp") || s.equalsIgnoreCase("-tsp")) {
                testSP = true;
            } else if (s.equalsIgnoreCase("-testbc") || s.equalsIgnoreCase("-tbc")) {
                testByteCode = true;
            } else if (s.equalsIgnoreCase("-testdedup") || s.equalsIgnoreCase("-tdd")) {
                testDedup = true;
            } else if (s.equalsIgnoreCase("-testorder") || s.equalsIgnoreCase("-to")) {
                testOrder = true;
            } else if (s.equalsIgnoreCase("-testorderddup") || s.equalsIgnoreCase("-tod")) {
                testOrderDedup = true;
            } else if (s.equalsIgnoreCase("-testvid") || s.equalsIgnoreCase("-tvid")) {
                testVids = true;
            } else if (s.equalsIgnoreCase("-testacl") || s.equalsIgnoreCase("-tacl")) {
                testACL = true;
            } else {
                out.printf("Skip argument %s\n", s);
            }
        }
        if (testValue || testCondition || testSpecialStep || testTraversal || testNegCase || testAggr ||
            testEdge || testException || testRSMeta || testGroup || testRepeat || testSP || testByteCode ||
            testDedup || testVids || testOrder || testOrderDedup || testACL) {
            testAll = false;
        }
    }

    private void setup() throws Exception {
        TGLogger.TGLevel logLevel = TGLogger.TGLevel.Debug;
        int i = 1;

        out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
        logger = TGLogManager.getInstance().getLogger();
        logger.setLevel(logLevel);

        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        //Following two lines function the same.  traversal() returns GraphTraversalSource
        //EmptyGraph.instance().traversal().withRemote(conn);
        //GraphTraversalSource g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        //GraphTraversalSource g = (GraphTraversalSource) EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        //Pass in TGConnection instead of RemoteConnection from Gremlin
        //We may look into supporting RemoteConnection
        g = EmptyGraph.instance().traversal(GraphTraversalSource.class).withRemote(conn);
        GraphTraversal t = g.V();


        //entityTypeExtractor = e -> e.getEntityType().getName();

        //The following lambdas can also be implemented as method and use method reference to invoke them.
        //The lambda expression can access the member variable directly
        //Simple count of the number item in the list match the expected value
        //The first item in 'e' is the number of value element in the list
        //resultSizeValidator = (l, r, e) -> {
        resultSizeValidator = c -> {
            List<Integer> e = (List<Integer>) c.expectedValue;
            boolean pass = (c.resultList.size() == e.get(1));
            out.printf("%d value(s) returned, expected %d\n", c.resultList.size(), e.get(1));
            return pass;
        };

        //Check numnber of entries and value of each entry
        //Watch out for ordering of the result which can cause verification to report erroneous failure
        listValueValidator = c -> {
            boolean pass = true;
            List<Integer> r = (List<Integer>) c.returnedValue;
            List<Integer> e = (List<Integer>) c.expectedValue;
            //The first value is the size of the list
            if (r.size() != e.size()) {
                out.printf("Comparison list size mismatch : result list size : %d, expected list size : %d\n", r.size(), e.size());
                pass = false;
            } if (!r.get(0).equals(e.get(0))) {
                //This should not happen if the previous condition check is false
                out.printf("Wrong result list size : %d, expected : %d\n", r.get(0), e.get(0));
                pass = false;
            } else {
                for (int idx=1; idx<e.size(); idx++) {
                    if (!r.get(idx).equals(e.get(idx))) {
                        out.printf("Item %d value : %d does not match expected value : %d\n", idx, r.get(idx), e.get(idx));
                        pass = false;
                        break;
                    }
                }
            }
            return pass;
        };

        mapValueValidator = c -> {
            boolean pass = true;
            List e = (List) c.expectedValue;
            if (c.resultList.size() != (int) e.get(0)) {
                out.printf("Wrong result list size : %d, expected : %d\n", c.resultList.size(), e.get(0));
                pass = false;
            } else {
                for (int idx=0; idx<c.resultList.size(); idx++) {
                    Map<String, Object> rMap = (Map<String, Object>) c.resultList.get(idx);
                    Map<String, Object> eMap = (Map<String, Object>) e.get(idx + 1);
                    if (rMap.size() != eMap.size()) {
                        out.printf("Item %d map size : %d does not match expected size : %d\n", idx + 1, rMap.size(), eMap.size());
                        pass = false;
                        break;
                    }
                    for (Map.Entry<String, Object> entry : eMap.entrySet()) {
                        Object rEntry = rMap.get(entry.getKey());
                        if (rEntry == null) {
                            out.printf("Item %d key value : '%s' not found\n", idx + 1, entry.getKey());
                            return false;
                        } else {
                            if (c.strategy == VerifyListSize) {
                                if (!(rEntry instanceof List)) {
                                    out.printf("Item %d value is not a list\n", idx + 1);
                                    return false;
                                } else {
                                    if (((Collection) rEntry).size() != (int) entry.getValue()) {
                                        out.printf("Item %d with key : '%s' has list size : %d but expected list size is : %d\n",
                                                idx + 1, entry.getKey(), ((Collection) rEntry).size(), entry.getValue());
                                        return false;
                                    }
                                }
                            } else {
                                //FIXME: Need to handle more types
                                if ((long) rEntry != (int) entry.getValue()) {
                                    out.printf("Item %d with key : '%s' and value : %d does not match the expected value : %d\n",
                                            idx + 1, entry.getKey(), rEntry, entry.getValue());
                                    return false;
                                }
                            }
                        }
                    }
                }
            }
            return pass;
        };

        //Check numnber of entities and numbner of attributes in each entity
        //Watch out for ordering of the result which can cause verification to report erroneous failure
        resultAttrSizeValidator = c -> {
            boolean pass = true;
            List<Integer> e = (List<Integer>) c.expectedValue;
            List<Integer> r = (List<Integer>) c.returnedValue;
            if (c.resultList.size() != e.get(0)) {
                out.printf("Wrong result list size : %d, expected : %d\n", c.resultList.size(), e.get(0));
                pass = false;
            } else {
                for (int idx=1; idx<e.size(); idx++) {
                    if (!r.get(idx).equals(e.get(idx))) {
                        out.printf("Item %d attribute count %d does not match expected count %d\n",
                                idx, r.get(idx), e.get(idx));
                        pass = false;
                        break;
                    }
                }
            }
            return pass;
        };

        //Check number of maps and numnber of entries in each map
        //It assumes each map has the same number of entries
        fixedSizeValidator = c -> {
            boolean pass = true;
            List<Integer> r = (List<Integer>)c.returnedValue;
            List<Integer> e = (List<Integer>)c.expectedValue;
            if (!r.get(0).equals(e.get(0))) {
                out.printf("Wrong result list size : %d, expected : %d\n", r.get(0), e.get(0));
                pass = false;
            } else {
                for (int idx=1; idx<r.size(); idx++) {
                    if (!r.get(idx).equals(e.get(1))) {
                        out.printf("Item %d value : %d does not match expected count : %d\n",
                                idx, r.get(idx), e.get(1));
                        pass = false;
                        break;
                    }
                }
            }
            return pass;
        };

        /*
        Same as the implementation below
        processResultList = l -> {
        IntStream.range(0, l.size()).forEach(idx -> out.printf("%d %s\n", idx+1, l.get(idx)));
        List<Integer> retValues = Arrays.asList(l.size());
        return retValues;
        };
        */

        resultListProcessor = c -> {
            int idx = 1;
            for (Object value : c.resultList) {
                out.printf("%d %s\n", idx++, value);
            }
            c.returnedValue = Collections.singletonList(c.resultList.size());
            return true;
        };

        mapInResultProcessor = c -> {
            List<Integer> retValues = new ArrayList<>();
            retValues.add(c.resultList.size());
            int idx = 1;
            for (Object value : c.resultList) {
                out.printf("%d %s\n", idx++, value);
                retValues.add(((Map) value).size());
            }
            c.returnedValue = retValues;
            return true;
        };

        //FIXME: Would be nice to combine this with mapInResultProcessor
        listInResultProcessor = c -> {
            List<Integer> retValues = new ArrayList<>();
            retValues.add(c.resultList.size());
            int idx = 1;
            for (Object value : c.resultList) {
                out.printf("%d %s\n", idx++, value);
                retValues.add(((List) value).size());
            }
            c.returnedValue = retValues;
            return true;
        };

        countInResultProcessor = c -> {
            List<Integer> retValues = new ArrayList<>();
            retValues.add(c.resultList.size());
            for (Object value : c.resultList) {
                out.println(value);
                retValues.add(((Long)value).intValue());
            }
            c.returnedValue = retValues;
            return true;
        };

        attrInResultProcessor = c -> {
            List<Integer> retValues = new ArrayList<>();
            retValues.add(c.resultList.size());
            int idx = 1;
            List<TGEntity> el = (List<TGEntity>) c.resultList;
            for (TGEntity ent : el) {
                Collection<TGAttribute> attrs = ent.getAttributes();
                if (idx > 1) {
                    out.println();
                }
                if (attrs.size() == 0) {
                    out.printf("%d has no attribute\n", idx);
                } else {
                    int ai = 1;
                    for (TGAttribute attr : attrs) {
                        out.printf("%d-%d Attr name : %s, value : %s\n", idx, ai++, attr.getAttributeDescriptor().getName(),
                                attr.getValue().toString());
                    }
                }
                retValues.add(attrs.size());
                idx++;
            }
            c.returnedValue = retValues;
            return true;
        };

        //May allow attribute to be specified instead of only the type name
        //output should be n [type, type, type]
        entityPathInResultProcessor = c -> {
            List<Integer> retValues = new ArrayList<>();
            retValues.add(c.resultList.size());
            int idx = 1;
            List<List<TGEntity>>pl = (List<List<TGEntity>>) c.resultList;
            for (List<TGEntity> p : pl) {
                StringBuilder sb = new StringBuilder("[");
                int ei = 0;
                for(TGEntity ent : p) {
                    if (ei++ > 0) {
                        sb.append(", ");
                    }
                    //FIXME: No need to use this complex approach to extract an entity value
                    //appendValueToStringBuilder(sb, ent, entityTypeExtractor);
                    sb.append(ent.getEntityType().getName());
                }
                sb.append("]");
                out.printf("%d %s\n", idx, sb.toString());
                sb.setLength(0);
                idx++;
            }
            c.returnedValue = retValues;
            return true;
        };

        resultCountProcessor = c -> {
            c.returnedValue = Collections.singletonList(c.resultList.size());
            return true;
        };
    }

    private void cleanup() throws Exception {
        conn.disconnect();
        out.println("Disconnected.");
    }

    /*
    private <E> void appendValueToStringBuilder(StringBuilder sb, E e, StringExtractor extractor) {
    sb.append(extractor.extract(e));
    }
    */

    //FIXME:  Display more error messages
    private boolean validateResultMetaData(ResultType expectedType, TGResultDataDescriptor rdd) {
        String ident = "";
        if (rdd == null) {
            return false;
        }
        if (expectedType.type != rdd.getDataType()) {
            return false;
        }
        if (rdd.getDataType() == TYPE_SCALAR) {
            if (!expectedType.scalarType.equals(rdd.getScalarType())) {
                return false;
            }
        } else if (rdd.getDataType() == TYPE_MAP) {
            TGResultDataDescriptor kdd = rdd.getKeyDescriptor();
            TGResultDataDescriptor vdd = rdd.getValueDescriptor();
            if (kdd == null || vdd == null) {
                return false;
            }
            if (expectedType.keyType != kdd.getDataType() ||
                    !expectedType.keyScalarType.equals(kdd.getScalarType())) {
                return  false;
            } else if (expectedType.valueType != vdd.getDataType() ||
                    !expectedType.valueScalarType.equals(vdd.getScalarType())) {
                return false;
            }
            //May be able to combine list and path check together
        } else if (rdd.getDataType() == TYPE_LIST) {
            TGResultDataDescriptor[] ndd = rdd.getContainedDescriptors();
            if (ndd == null) {
                return false;
            } else if (expectedType.containedType.length != ndd.length) {
                return false;
            }
            for (int i=0; i<ndd.length; i++) {
                boolean pass = validateResultMetaData(expectedType.containedType[i], ndd[i]);
                if (pass == false) {
                    return pass;
                }
            }
        } else if (rdd.getDataType() == TYPE_PATH) {
            TGResultDataDescriptor[] ndd = rdd.getContainedDescriptors();
            if (ndd == null) {
                return false;
            } else if (expectedType.containedType.length != ndd.length) {
                return false;
            }
            for (int i=0; i<ndd.length; i++) {
                boolean pass = validateResultMetaData(expectedType.containedType[i], ndd[i]);
                if (pass == false) {
                    return pass;
                }
            }
        }
        return true;
    }

    /**
     * Test runner use for validating result against a list of values
     * @param testName        Name of the test
     * @param targetValues    The expected values used for validation
     * @param processor    Code to process and display the result list
     * @param validator    Code to Validate the result
     * @param traversal    DefaultGraphTraversal to be executed
     * @throws Exception
     */
    //private void test(String testName, List<Integer> targetValues, Processor processor, ValueValidator validator, Query query) throws Exception {
    private void execTest(String testName, List<Integer> targetValues, Processor processor, Validator validator, DefaultGraphTraversal traversal) throws Exception {
        out.println(testName);
        out.println(traversal.getBytecode());
        List resultList = traversal.toList();
        TestContext ctx = new TestContext(resultList, targetValues);
        processor.process(ctx);
        boolean pass = validator.validate(ctx);
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    /**
     * Test runner use for validating result against a fixed value
     * @param testName		Name of the test
     * @param targetValue 	The expected value used for validation
     * @param processor 	Code to process and display the result list
     * @param validator 	Code to Validate the result
     * @param traversal 	DefaultGraphTraversal to be executed
     * @throws Exception
     */
    private void execTest(String testName, int targetValue, Processor processor, Validator validator, DefaultGraphTraversal traversal) throws Exception {
        List<Integer> targetValues = new ArrayList<>();
        targetValues.add(1);
        targetValues.add(targetValue);
        execTest(testName, targetValues, processor, validator, traversal);
    }

    /**
     * Test runner use for validating result against a list of values of a string query
     * @param testName		Name of the test
     * @param targetValues 	The expected values used for validation
     * @param processor 	Code to process and display the result list
     * @param validator 	Code to Validate the result
     * @param query 		The query to be executed
     * @throws Exception
     */
    private void execTest(String testName, List<Integer> targetValues, Processor processor, Validator validator, String query) throws Exception {
        boolean pass = false;
        out.println(testName);
        out.println(query);
        try {
            TGResultSet<Object> resultSet = conn.executeQuery(query, null);
            List valueList = new ArrayList(resultSet.toCollection());
            TestContext ctx = new TestContext(valueList, targetValues);
            processor.process(ctx);
            pass = validator.validate(ctx);
        } catch (TGException e) {
            out.printf("Unexpected exception : %s(%s)\n", e.getMessage(), e.getExceptionType());
        }
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    /**
     * Test runner use for validating result against a list of values of a string query
     * @param testName		Name of the test
     * @param targetValue 	The expected value used for validation
     * @param processor 	Code to process and display the result list
     * @param validator 	Code to Validate the result
     * @param query 		The query to be executed
     * @throws Exception
     */
    private void execTest(String testName, int targetValue, Processor processor, Validator validator, String query) throws Exception {
        List<Integer> targetValues = new ArrayList<>();
        targetValues.add(1);
        targetValues.add(targetValue);
        execTest(testName, targetValues, processor, validator, query);
    }

    /**
     * Basic test runner for custom display and validation
     * @param testName        Name of the test
     * @param ctx        Code responsible for display and checking the results
     * @param traversal  DefaultGraphTraversal to be executed
     * @throws Exception
     */
    private void execTest(String testName, TestContext ctx, DefaultGraphTraversal traversal) throws Exception {
        out.println(testName);
        out.println(traversal.getBytecode());
        //The lambda expression can access the member variable directly
        //List resultList = query.execute();
        ctx.resultList = traversal.toList();
        ctx.processor.process(ctx);
        boolean pass = ctx.validator.validate(ctx);
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    /**
     * Test runner use for validating result against a list of values
     * @param testName		Name of the test
     * @param expectedType  Expected result set meta data annotation
     * @param query 	    Gremlin query string
     * @throws Exception
     */
    private void execTest(String testName, ResultType expectedType, String query) throws Exception {
        boolean pass = false;
        out.println(testName);
        out.println(query);
        try {
            TGResultSet<Object> resultSet = conn.executeQuery(query, null);
            TGResultSetMetaData rsmd = resultSet.getMetaData();
            ResultSetUtils.printRSMetaData(resultSet);
            int idx = 1;
            while (resultSet.hasNext()) {
                Object value = resultSet.next();
                out.printf("%d %s\n", idx++, value);
            }
            if ((expectedType == null && rsmd == null)) {
                pass = true;
            } else if ((expectedType != null && rsmd == null)) {
                out.println("Result set has no meta data - unexpected");
            } else if (expectedType == null && rsmd != null) {
                out.println("Result set has meta data - unexpected");
            } else {
                String annot = ((ResultSetMetaData) rsmd).getAnnot();
                out.println("Result type annotation : " + annot);
                if (expectedType.annot != null && !expectedType.annot.equals(annot)) {
                    out.printf("Annotation string mismatch - expected : %s, returned : %s\n", expectedType.annot, annot);
                } else {
                    pass = validateResultMetaData(expectedType, rsmd.getResultDataDescriptor());
                }
            }

        } catch (TGException e) {
            out.printf("Unexpected exception : %s(%s)\n", e.getMessage(), e.getExceptionType());
        }
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    /**
     * Test runner use for validating result against a list of values of a string query
     * @param testName		Name of the test
     * @param ctx		    Test context
     * @param query 		The query to be executed
     * @throws Exception
     */
    private void execTest(String testName, TestContext ctx, String query) throws Exception {
        boolean pass = false;
        out.println(testName);
        out.println(query);
        try {
            TGResultSet<Object> resultSet = conn.executeQuery(query, null);
            List resultList = new ArrayList(resultSet.toCollection());
            ctx.resultList = resultList;
            ctx.processor.process(ctx);
            pass = ctx.validator.validate(ctx);
        } catch (TGException e) {
            out.printf("Unexpected exception : %s(%s)\n", e.getMessage(), e.getExceptionType());
        }
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    //Need to add check for specific exception
    private void testException(String testName, String query) throws Exception {
        boolean pass = false;
        out.println(testName);
        try {
            conn.executeQuery(query, null);
            out.printf("Expected exception not thrown\n");
        } catch (TGException e) {
            out.printf("Got expected exception : %s(%s)\n", e.getMessage(), e.getExceptionType());
            pass = true;
        }
        out.println(testName + " ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    private void testValues() throws Exception {
        if (testAll == false && testValue == false) {
            out.println("Skip Values test");
            return;
        }
        out.println("Values test start");
        //This should return a list of primitive values
        execTest("Test values", 4, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").values());

        execTest("Test values count", 4, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").values().count());

        /*
        execTest("Test values fold", 1, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").values().fold());
        */

        //returns a list with a single element of a map of key/value
        execTest("Test valueMap all", 4, mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap());

        //returns a list with a single element of a map of key/value
        execTest("Test valueMap select", 2, mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap("itemname", "oops", "itemid"));

        //returns a list with a single element of a map of key/value
        execTest("Test valueMap count", 1, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap().count());

        //returns a list with a single element of a map of key/value
        /*
        execTest("Test valueMap fold", 1, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").valueMap().fold());
                */

        execTest("Test cdibatch values", 1, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdibatch", "batchid", "17012LXFP342").values("batchid", "oops"));
        out.println("Values test end");
        out.println();
    }

    private void testConditions() throws Exception {
        if (testAll == false && testCondition == false) {
            out.println("Skip Conditions test");
            return;
        }
        out.println("Conditions test start");
        //return TGNode directly
        execTest("Test has step only", Arrays.asList(1,4), attrInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44"));

        execTest("Test and-and condition", Arrays.asList(1,4), attrInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").has("groupid", 2200));

        execTest("Test conditions", Arrays.asList(3,4,2,2), attrInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", P.eq("172CBAFEVZ08").or(P.eq("172CBAFFPU57")).or(P.eq("172CBAFGLK39"))).order().by("cdiid"));

        execTest("Test conditions valueMap", Arrays.asList(3,2,1,1), mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", P.eq("172CBAFEVZ08").or(P.eq("172CBAFFPU57")).or(P.eq("172CBAFGLK39"))).order().by("itemid").valueMap("itemid", "itemname"));

        /*FIXME: Disable this test 4/16/2019.  Takes up too much memory. Revisit later.
         out.println("Test between conditions valueMap");
         * Can add a non-unique index for groupid to test out the index and
         * at the same time reduce memory usage due to scanning all the entities in the db.
         */
        /* Disable again 8/27/2020 because because of memory issue */
        execTest("Test between conditions valueMap", 2724, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "groupid", P.between(700, 701)).valueMap("itemid", "itemname", "groupid"));
         /**/

        //FIXME:  Need to rework index check logic on the server side.
        //This query should do a unique get instead of the
        //range up prefetch on the server side.
        /* disable for now to speed up the test 4/17/2019
          out.println("Test and/or steps");
//        List valueList = g.V().hasLabel("cdi").and(
        valueList = g.V().hasLabel("cdi").and(
        		__.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57")),
        		__.has("groupid",P.gte(1700))).valueMap().toList();
        for (Object value : valueList) {
        out.println(value);
        }
        out.println("Test and/or steps ended");
        out.println();
        */

        execTest("Test and/or with empty or steps", Arrays.asList(3,2,1,4), mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(
                __.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57"))).or().
                hasLabel("cdibatch").and(__.has("batchid","17012LXFP342")).valueMap());

        execTest("Test or after V", Arrays.asList(3,2,1,4), mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().or(__.hasLabel("cdi").and(
                __.or(__.has("cdiid","172CBAFEVZ08"),__.has("cdiid", "172CBAFFPU57"))),
                __.hasLabel("cdibatch").and(__.has("batchid","17012LXFP342"))).valueMap());

        execTest("V limit", 2, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("JFK")).
                or(P.eq("HKG").or(P.eq("YYZ")))).limit(2));
        out.println("Conditions test end");
        out.println();
    }

    private void testSpecialSteps() throws Exception {
        if (testAll == false && testSpecialStep == false) {
            out.println("Skip Special Step test");
            return;
        }
        out.println("Special Step test start");
        out.println("Test pagerank");
        List valueList = g.V().pageRank().valueMap("@pagerank", "itemname").toList();
        for (Object value : valueList) {
            out.println(value);
        }
        out.println("Test pagerank - not working currently");
        out.println("Test pagerank ended");
        out.println();
        out.println("Special Step test end");
        out.println();
    }

    private void testTraversals() throws Exception {
        if (testAll == false && testTraversal == false) {
            out.println("Skip Traversal test");
            return;
        }
        out.println("Traversal test start");
        execTest("Traversal 1", Arrays.asList(4,2,2,2,2), mapInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CBAFEVZ08").outE("produces").has("quantity", P.gt(5)).
                inV().has("groupid", P.gt(1000)).out("produces").valueMap());

        execTest("Traversal 2", Arrays.asList(280,4), mapInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXEI727")).outE("contains").simplePath().
                inV().has("groupid", P.neq(1000)).limit(10).out("produces").valueMap());

        execTest("Traversal 3", Arrays.asList(28,4), mapInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXEI727")).outE("contains").inV().valueMap());

        execTest("Traversal 4", Arrays.asList(28,1), mapInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").valueMap("quantity"));

        execTest("Traversal 5", Arrays.asList(28,5), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().path().by("cdiid").by("quantity"));

        /*
        execTest("Traversal 6 multiple folds", 1, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CBAFEVZ08").outE("produces").has("quantity", P.gt(5)).
                inV().has("groupid", P.gt(1000)).out("produces").values().fold().fold().fold());
                        */

        execTest("Traversal 7.0", 32, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                        outE("produces").inV().outE("produces").inV().valueMap());

        //FIXME: Not working. -- Need to investigate -- similar behavior using gremlin console also
        /*
        valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18").or(__.has("cdiid", "172CBAFEVZ08")).
        or(__.has("cdiid", "172CBAFFPU57"))).
        outE("produces").inV().outE("produces").inV().path().by("cdiid").by("quantity").toList();
        */
        execTest("Traversal 7", Arrays.asList(32,5), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").inV().path().by("cdiid").by("quantity"));

        //FIXME:Not supporting by just entity in the path - not sure why
        /*
        out.println("Traversal 7.1");
        valueList = g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
        or(P.eq("172CBAFFPU57")))).
        outE("produces").inV().outE("produces").inV().path().toList();
        for (Object value : valueList) {
        out.println(value);
        }
        out.println("Traversal 7.1 ended");
        out.println();
        */

        execTest("Traversal 7.2", 32, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").inV().path().count());

        execTest("Traversal 7.3", Arrays.asList(32,5), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").inV().simplePath().path().by("cdiid").by("quantity"));

        execTest("Traversal 7.4", Arrays.asList(32,5), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").inV().simplePath().path());

        execTest("Traversal 7.5", 32, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").inV().simplePath().path().count());

        //Not supporting by just entity in the path
        execTest("Traversal Flight data 1", Arrays.asList(64629,5), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").inV().outE("routeType").inV().path().by("iataCode"));

        //FIXME: Check on this test
        /*
        out.println("Traversal Flight data 2");
        valueList = g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
        outE("routeType").inV().outE("routeType").inV().outE("routeType").inV().
        outE("routeType").inV().has("iataCode", "CDG").path().count().toList();
        for (Object value : valueList) {
        out.println(value);
        }
        out.println("Traversal Flight data 2 ended");
        out.println();
        */

        execTest("Traversal Flight data 2.0", Arrays.asList(92607,7), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").inV().
                outE("routeType").inV().
                outE("routeType").inV().
                simplePath().
                has("iataCode", "CDG").path().by("iataCode"));

        execTest("Traversal Flight data 2.0 - count", 2213, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
//				outE("routeType").inV().
                simplePath().has("iataCode", "CDG").path().count());

        execTest("Traversal Flight data 2.1", Arrays.asList(102779,6), mapInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").inV().
                outE("routeType").inV().
                outE("routeType").inV().
                has("iataCode", "CDG").valueMap());

        execTest("Traversal Flight data 2.2", Arrays.asList(64629,3), listInResultProcessor, fixedSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").
                out("routeType").
                out("routeType").
                //has("iataCode", "CDG").
                path().by("iataCode"));

        //Gremlin string tests
        execTest("Traversal string 1", Arrays.asList(28,5), entityPathInResultProcessor, fixedSizeValidator,
                "gremlin://g.V().hasLabel('cdi').has('cdiid', '172CDIXDQC18').outE('produces').inV().outE('produces').inV().path();");

        execTest("Out edges from SFO", 249, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").has("iataCode", "SFO").
                outE("routeType").values("iataCode"));

        execTest("In edges to SFO", 250, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").has("iataCode", "SFO").
                inE("routeType").values("iataCode"));

        execTest("Flight from SFO to RNO", Arrays.asList(1,3), listInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").has("iataCode", "SFO").
                outE("routeType").inV().has("airportType", "iataCode", "RNO").path().by("iataCode"));
        out.println("Traversal test end");
        out.println();
    }

    private void testNegativeCases() throws Exception {
        if (testAll == false && testNegCase == false) {
            out.println("Skip Negative Case test");
            return;
        }
        out.println("Negative Case test start");
        gremlinQueryIllegalSequence1();
        gremlinQueryIllegalSequence2();
        out.println("Negative Case test end");
        out.println();
    }

    //Aggregation tests
    private void testAggregations() throws Exception {
        if (testAll == false && testAggr == false) {
            out.println("Skip Aggregation test");
            return;
        }
        out.println("Aggregation test start");
        execTest("Aggregation raw data", 28, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().outE().values("quantity"));

        execTest("Sum", 2230, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().outE().values("quantity").sum());

        execTest("Max", 80, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().outE().values("quantity").max());

        execTest("Min", 79, countInResultProcessor, listValueValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().outE().values("quantity").min());

        //Test 'mean' step
        boolean pass = false;
        out.println("Mean");
        List valueList = g.V().hasLabel("cdi").and(__.has("cdiid", "172CDIXDQC18")).outE("produces").inV().
                outE("produces").inV().outE().values("quantity").mean().toList();
        if (valueList.size() == 1) {
            double val = (double) valueList.get(0);
            if (val > 79.63 && val < 79.65) {
                pass = true;
            }
        }
        for (Object value : valueList) {
            out.println(value);
        }
        out.println("Mean ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
        out.println("Aggregation test end");
        out.println();
    }

    private void testEdges() throws Exception {
        if (testAll == false && testEdge == false) {
            out.println("Skip Edge test");
            return;
        }
        out.println("Edge test start");
        //Edge tests
        execTest("Get all edges", 1996020, resultCountProcessor, resultSizeValidator, (DefaultGraphTraversal) g.E());

        execTest("Get 5 edges", 5, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal) g.E().limit(5));

        execTest("Get all 'produces' edges", 489165, resultCountProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.E().hasLabel("produces").has("quantity", P.gt(30)));

        execTest("Get all 'UA' edges", 1089, resultCountProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.E().hasLabel("routeType").has("iataCode", "UA").valueMap());

        execTest("Traverse from 'UA' edges to nodes to 'SW' edges", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.E().hasLabel("routeType").has("iataCode", "UA").
                inV().outE("routeType").has("iataCode", "SW").path().by("iataCode"));

        execTest("Traverse from 'UA' edges to nodes", 1089, resultCountProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.E().hasLabel("routeType").has("iataCode", "UA").inV().path().by("iataCode"));
        out.println("Edge test end");
        out.println();
    }

    private void testExceptionHandling() throws Exception {
        if (testAll == false && testException == false) {
            out.println("Skip Exception Handling test");
            return;
        }
        out.println("Exception Handling test start");
        testException("Traversal bad string 1 - 'limitt'",
            "gremlin://g.V().hasLabel('cdi').limitt(10).has('cdiid', '172CDIXDQC18').outE('produces').inV().outE('produces').inV().path();");
        testException("Traversal bad string 2 - 'unknown'",
            "gremlin://g.V().has('airportType', 'iataCode', 'SFO').unknown();");
        out.println("Exception Handling test end");
        out.println();
    }

    private void testResultSetMetaData() throws Exception {
        if (testAll == false && testRSMeta == false) {
            out.println("Skip RS Meta Data test");
            return;
        }
        out.println("RS meta data test start");

        ResultType rt = new ResultType();

        rt.annot = "{s,S}";
        rt.type = TYPE_MAP;
        rt.keyType = TYPE_SCALAR;
        rt.keyScalarType = TGAttributeType.String;
        rt.valueType = TYPE_SCALAR;
        rt.valueScalarType = TGAttributeType.Invalid;
        execTest( "Value map resultset meta data", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().valueMap().limit(5);");

        rt = new ResultType();
        rt.annot = "P(V,E,V,V)";
        rt.type = TYPE_PATH;
        rt.containedType = new ResultType[4];
        for (int i=0; i<4; i++) {
            rt.containedType[i] = new ResultType();
            rt.containedType[i].type = TYPE_NODE;
        }
        rt.containedType[1].type = TYPE_EDGE;
        execTest( "Path resultset meta data", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().out().has('iataCode', 'CDG').path().limit(5);");

        //FIXME: Support union values - use O or | symbol
        rt = new ResultType();
        rt.annot = "P(s,s,s,s,s)"; //"P[s]", "P([s])", "P(d,[s])", "P(d,[(s,V)]) - for repeating pattern", how to specify nullable element;
        rt.type = TYPE_PATH;
        rt.containedType = new ResultType[5];
        for (int i=0; i<5; i++) {
            rt.containedType[i] = new ResultType();
            rt.containedType[i].type = TYPE_SCALAR;
            rt.containedType[i].scalarType = TGAttributeType.String;
        }
        execTest( "Path(String) resultset meta data", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().outE().inV().path().by('iataCode').limit(5);");

        rt = new ResultType();
        rt.annot = "[P(s,s,s,s,s)]";
        rt.type = TYPE_LIST;
        rt.containedType = new ResultType[1];
        rt.containedType[0] = new ResultType();
        rt.containedType[0].type = TYPE_PATH;
        rt.containedType[0].containedType = new ResultType[5];
        for (int i=0; i<5; i++) {
            rt.containedType[0].containedType[i] = new ResultType();
            rt.containedType[0].containedType[i].type = TYPE_SCALAR;
            rt.containedType[0].containedType[i].scalarType = TGAttributeType.String;
        }
        /*
        execTest( "List of paths resultset meta data", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().outE().inV().path().by('iataCode').limit(5).fold();");
        */

        rt = new ResultType();
        rt.annot = "l";
        rt.type = TYPE_SCALAR;
        rt.scalarType = TGAttributeType.Long;
        execTest( "Long value resultset meta data 1", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().out().path().count();");

        rt = new ResultType();
        rt.annot = "l";
        rt.type = TYPE_SCALAR;
        rt.scalarType = TGAttributeType.Long;
        execTest( "Long value resultset meta data 2", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFOOO" + "').outE().inV().out().path().count();");

        rt = null;
        execTest( "Empty resultset meta data 1", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFOOO" + "').outE().inV().out().path();");

        execTest( "Empty resultset meta data 2", rt,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFOOO" + "');");
        out.println("RS meta data test end");
        out.println();
    }

    //FIXME: Need to add data file initialization support
    //Test name as section name
    //[Data format type:Name of the test]
    //value record each line
    //e.g. [MapIntValue:Group by Test 1]
    //Mercury:1516
    //Venus:3760
    //Earth:3959
    private void setupGroupTestValues1(Map<String, Integer> valueMap) throws Exception {
        valueMap.clear();
        valueMap.put("Hong Kong", 6);
        valueMap.put("United States", 148);
        valueMap.put("Philippines", 1);
        valueMap.put("Japan", 6);
        valueMap.put("United Kingdom", 10);
        valueMap.put("United Arab Emirates", 2);
        valueMap.put("Switzerland", 2);
        valueMap.put("New Zealand", 3);
        valueMap.put("Canada", 14);
        valueMap.put("South Korea", 7);
        valueMap.put("Netherlands", 2);
        valueMap.put("El Salvador", 3);
        valueMap.put("Ireland", 1);
        valueMap.put("China", 7);
        valueMap.put("Taiwan", 5);
        valueMap.put("Denmark", 1);
        valueMap.put("Mexico", 20);
        valueMap.put("France", 4);
        valueMap.put("Australia", 2);
        valueMap.put("Germany", 5);
    }

    private void setupGroupTestValues2(Map<String, Integer> valueMap) throws Exception {
        valueMap.clear();
        valueMap.put("Ghana", 1); valueMap.put("Congo (Brazzaville)", 4); valueMap.put("Bahrain", 1); valueMap.put("India", 5);
        valueMap.put("Canada", 44); valueMap.put("Turkey", 8); valueMap.put("Belgium", 8); valueMap.put("Taiwan", 7);
        valueMap.put("Finland", 7); valueMap.put("Trinidad and Tobago", 3); valueMap.put("Netherlands Antilles", 6);
        valueMap.put("South Africa", 6); valueMap.put("Bermuda", 4); valueMap.put("Georgia", 2);
        valueMap.put("Central African Republic", 1); valueMap.put("Jamaica", 6); valueMap.put("Peru", 5);
        valueMap.put("Germany", 41); valueMap.put("Puerto Rico", 7); valueMap.put("Hong Kong", 11);
        valueMap.put("Guinea", 1); valueMap.put("United States", 393); valueMap.put("Chad", 1); valueMap.put("Aruba", 2);
        valueMap.put("Madagascar", 2); valueMap.put("Thailand", 3); valueMap.put("Costa Rica", 5); valueMap.put("Sweden", 7);
        valueMap.put("Vietnam", 6); valueMap.put("Poland", 5); valueMap.put("Jordan", 5); valueMap.put("Nigeria", 5);
        valueMap.put("Kuwait", 1); valueMap.put("Bulgaria", 2); valueMap.put("Tunisia", 8); valueMap.put("Croatia", 4);
        valueMap.put("Sri Lanka", 1); valueMap.put("United Kingdom", 60); valueMap.put("United Arab Emirates", 11);
        valueMap.put("Kenya", 2); valueMap.put("Switzerland", 12); valueMap.put("Spain", 43);
        valueMap.put("Lebanon", 2); valueMap.put("Djibouti", 1); valueMap.put("Venezuela", 2); valueMap.put("Liberia", 2);
        valueMap.put("Azerbaijan", 2); valueMap.put("Cuba", 1); valueMap.put("Czech Republic", 6); valueMap.put("Saint Lucia", 1);
        valueMap.put("Burkina Faso", 1); valueMap.put("Mauritania", 1); valueMap.put("Israel", 5); valueMap.put("Australia", 2);
        valueMap.put("Cameroon", 4); valueMap.put("Cyprus", 2); valueMap.put("Malaysia", 3); valueMap.put("Iceland", 4);
        valueMap.put("Oman", 1); valueMap.put("Armenia", 3); valueMap.put("Gabon", 1); valueMap.put("Austria", 6);
        valueMap.put("South Korea", 14); valueMap.put("El Salvador", 6); valueMap.put("Luxembourg", 2); valueMap.put("Brazil", 12);
        valueMap.put("Turks and Caicos Islands", 2); valueMap.put("Algeria", 9); valueMap.put("Jersey", 1);
        valueMap.put("Slovenia", 2); valueMap.put("Antigua and Barbuda", 2); valueMap.put("Ecuador", 2); valueMap.put("Colombia", 8);
        valueMap.put("Hungary", 2); valueMap.put("Japan", 19); valueMap.put("Belarus", 1); valueMap.put("Mauritius", 2);
        valueMap.put("New Zealand", 3); valueMap.put("Senegal", 3); valueMap.put("Honduras", 1); valueMap.put("Italy", 53);
        valueMap.put("Ethiopia", 1); valueMap.put("Haiti", 4); valueMap.put("Singapore", 2); valueMap.put("Egypt", 6);
        valueMap.put("Russia", 8); valueMap.put("Malta", 2); valueMap.put("Saudi Arabia", 6); valueMap.put("Cape Verde", 4);
        valueMap.put("Netherlands", 7); valueMap.put("Pakistan", 3); valueMap.put("Ireland", 15);
        valueMap.put("China", 22); valueMap.put("Martinique", 1); valueMap.put("Lithuania", 1); valueMap.put("France", 59);
        valueMap.put("Serbia", 2); valueMap.put("Reunion", 1); valueMap.put("Romania", 2); valueMap.put("Togo", 1);
        valueMap.put("Niger", 3); valueMap.put("Philippines", 1); valueMap.put("Cote d'Ivoire", 2); valueMap.put("Uzbekistan", 2);
        valueMap.put("Congo (Kinshasa)", 1); valueMap.put("Barbados", 1); valueMap.put("Norway", 7);
        valueMap.put("Dominican Republic", 14); valueMap.put("Denmark", 8); valueMap.put("Mexico", 33); valueMap.put("Montenegro", 2);
        valueMap.put("Benin", 1); valueMap.put("Angola", 1); valueMap.put("Portugal", 5); valueMap.put("Bahamas", 2);
        valueMap.put("Grenada", 2); valueMap.put("Greece", 7); valueMap.put("Cayman Islands", 2); valueMap.put("Latvia", 3);
        valueMap.put("Morocco", 10); valueMap.put("Mali", 1); valueMap.put("Panama", 3); valueMap.put("Guadeloupe", 1);
        valueMap.put("Guatemala", 1); valueMap.put("Guyana", 1); valueMap.put("Chile", 4); valueMap.put("Argentina", 5);
        valueMap.put("Virgin Islands", 3); valueMap.put("Ukraine", 3);
    }

    //FIXME: Server does not support T.label yet.
    //FIXME: Need to verify caliper related test results
    private void testGroupSteps() throws Exception {
        if (testAll == false && testGroup == false) {
            out.println("Skip Group test");
            return;
        }
        out.println("Group test start");
        Map<String, Integer> expectedValues = new HashMap<>();
        setupGroupTestValues1(expectedValues);

        TestContext ctx = new TestContext();
        ctx.processor = resultListProcessor;
        ctx.validator = mapValueValidator;
        List evList = new ArrayList();
        evList.add(1);
        evList.add(expectedValues);
        ctx.expectedValue = evList;
        execTest("Group count string test 1", ctx,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().groupCount().by('country');");

        execTest("Group count string test 2", ctx,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().values('country').groupCount();");

        //FIXME: The result meta data returned is not correct. It's a server issue.
        //I believe server can describe a map of list
        ctx.strategy = VerifyListSize;
        execTest("Group string test 1", ctx,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().group().by('country');");

        //FIXME: The result meta data returned is not correct. It's a server issue.
        //I believe server can describe a map of list
        execTest("Group string test 2", ctx,
        "gremlin://g.V().has('airportType', 'iataCode', '" + "SFO" + "').outE().inV().group().by('country').by('iataCode');");

        expectedValues.clear();
        expectedValues.put("United States", 2);
        expectedValues.put("France", 1);
        ctx.strategy = VerifyListValue;
        execTest("Group count test 1", ctx, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("JFK")).or(P.eq("CDG"))).values("country").groupCount());

        setupGroupTestValues2(expectedValues);
        execTest("Group count test 2", ctx, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("JFK")).or(P.eq("CDG"))).outE().inV().groupCount().by("country"));

        expectedValues.clear();
        expectedValues.put("77", 1); expectedValues.put("23", 1); expectedValues.put("57", 1); expectedValues.put("59", 1);
        expectedValues.put("79", 10); expectedValues.put("80", 18);
        execTest("Group count test 3", ctx, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").values("quantity").groupCount());

        execTest("Group count test 3_1", 32, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).or(P.eq("172CBAFFPU57")))).
                outE("produces").inV().outE("produces").values("quantity"));

        setupGroupTestValues1(expectedValues);
        ctx.strategy = VerifyListSize;
        execTest("Group test 1", ctx, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").outE().inV().group().by("country"));

        execTest("Group test 2", ctx, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").outE().inV().group().by("country").by("iataCode"));

        expectedValues.clear();
        expectedValues.put("77", 1); expectedValues.put("23", 1); expectedValues.put("57", 1); expectedValues.put("59", 1);
        expectedValues.put("79", 10); expectedValues.put("80", 18);
        execTest("Group test 3", ctx, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).outE("produces").inV().outE("produces").group().by("quantity"));

        execTest("Group test 4", ctx, (DefaultGraphTraversal)
            g.V().hasLabel("cdi").and(__.has("cdiid", P.eq("172CDIXDQC18").or(P.eq("172CBAFEVZ08")).
                or(P.eq("172CBAFFPU57")))).outE("produces").inV().outE("produces").group().by("quantity").by("prodid"));

        //groupCount by entity is not supported and therefore the results should be zero
        execTest("Group test 5", 0, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("JFK")).or(P.eq("CDG"))).groupCount());
        out.println("Group test end");
        out.println();
    }

    private void testRepeats() throws Exception {
        if (testAll == false && testRepeat == false) {
            out.println("Skip Repeat test");
            return;
        }
        out.println("Repeat test start");
        //returns a list of 9 values 'JFK'
        execTest("Repeat step with end condition", 9, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").repeat(__.outE("routeType").has("iataCode", "UA").inV()).
                times(2).has("iataCode", "JFK"));

        execTest("Repeat step with end condition and emit", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").repeat(__.outE("routeType").has("iataCode", "UA").inV()).emit().
                times(2).has("iataCode", "JFK"));

        execTest("Repeat step with end condition and path", 9, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").repeat(__.outE("routeType").has("iataCode", "UA").inV()).
                times(2).has("iataCode", "JFK").simplePath().path());

        execTest("Repeat step with end condition and path", 9, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").repeat(__.outE("routeType").has("iataCode", "UA").inV()).
                times(2).has("iataCode", "JFK").simplePath().path().by("iataCode"));

        execTest("Two hopes traversal and path 1", 9, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().has("iataCode", "JFK").simplePath().path().by("iataCode"));

        execTest("Repeat step with end condition and emit and path 1", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").repeat(__.outE("routeType").has("iataCode", "UA").inV()).emit().
                times(2).has("iataCode", "JFK").simplePath().path().by("iataCode"));

        execTest("Repeat step with end condition and emit and path 2", 38, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("AUS"))).repeat(__.outE("routeType").has("iataCode", "UA").inV()).emit().
                times(2).has("iataCode", P.eq("JFK").or(P.eq("YYZ"))).simplePath().path().by("iataCode"));

        execTest("Two hopes traversal and path 2",36, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("AUS"))).outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().has("iataCode", P.eq("JFK").or(P.eq("YYZ"))).simplePath().path().by("iataCode"));

        execTest("One hope traversal and path",2, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("AUS"))).outE("routeType").has("iataCode", "UA").inV().
                has("iataCode", P.eq("JFK").or(P.eq("YYZ"))).simplePath().path().by("iataCode"));

        execTest("Repeat step with end condition and emit and path string query 1", 1358, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', 'SFO').repeat(bothE('routeType').has('iataCode', 'UA').bothV()).emit()." +
                        "times(3).has('iataCode', 'JFK').simplePath().path().by('iataCode');");

        execTest("Repeat step with end condition and emit and path string query 2", 6648, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', eq('SFO').or(eq('AUS'))).repeat(bothE('routeType').has('iataCode', 'UA').bothV()).emit()." +
                        "times(3).has('iataCode', eq('JFK').or(eq('YYZ'))).simplePath().path().by('iataCode');");

        execTest("Repeat step with end condition and emit and path string query 3", 175, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', 'SFO').repeat(outE('routeType').has('iataCode', 'UA').inV()).emit()." +
                        "times(3).has('iataCode', 'JFK').simplePath().path().by('iataCode');");

        execTest("Repeat step with end condition and emit and path string query 4", 175, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', 'SFO').repeat(bothE('routeType').has('iataCode', 'UA').inV()).emit()." +
                        "times(3).has('iataCode', 'JFK').simplePath().path().by('iataCode');");

        execTest("Repeat step with end condition and emit and path string query 5", 175, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', 'SFO').repeat(outE('routeType').has('iataCode', 'UA').bothV()).emit()." +
                        "times(3).has('iataCode', 'JFK').simplePath().path().by('iataCode');");

        execTest("Two hopes traversal and path 1", 1320, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", "SFO").bothE("routeType").has("iataCode", "UA").bothV().
                bothE("routeType").has("iataCode", "UA").bothV().
                bothE("routeType").has("iataCode", "UA").bothV().
                has("iataCode", "JFK").simplePath().path().by("iataCode"));

        execTest("Repeat step with end condition and emit and path string query 2", 6500, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', eq('SFO').or(eq('AUS'))).repeat(bothE('routeType').has('iataCode', 'UA').bothV())." +
                        "times(3).has('iataCode', eq('JFK').or(eq('YYZ'))).simplePath().path().by('iataCode');");

        /*

        execTest("Repeat step with end condition and emit and path string query 2", 6648, resultListProcessor, resultSizeValidator,
                "gremlin://g.V().has('airportType', 'iataCode', 'SFO').emit().repeat(bothE('routeType').has('iataCode', 'UA').bothV())." +
                        "times(2).has('iataCode', eq('JFK').or(eq('YYZ'))).simplePath().path().by('iataCode');");
         */

        out.println("Repeat test end");
        out.println();
    }

    private void testByteCodeQuery() throws Exception {
        if (testAll == false && testByteCode == false) {
            out.println("Skip Bytecode test");
            return;
        }
        out.println("Bytecode test start");
        Map<String, Integer> expectedValues = new HashMap<>();
        TestContext ctx = new TestContext();
        ctx.processor = resultListProcessor;
        ctx.validator = mapValueValidator;
        ctx.strategy = VerifyListSize;
        List evList = new ArrayList();
        evList.add(1);
        evList.add(expectedValues);
        ctx.expectedValue = evList;
        expectedValues.put("77", 1); expectedValues.put("23", 1); expectedValues.put("57", 1); expectedValues.put("59", 1);
        expectedValues.put("79", 10); expectedValues.put("80", 18);
        execTest("Byte code test 1(should have same results as group test 3)", ctx,
            "gbc://[[], [V(), hasLabel(cdi), and([[], [has(cdiid, or(eq(172CDIXDQC18), eq(172CBAFEVZ08), eq(172CBAFFPU57)))]]), outE(produces), inV(), outE(produces), group(), by(quantity)]]");
        out.println("Bytecode test end");
        out.println();
    }

    private void testSP() throws Exception {
        if (testAll == false && testSP == false) {
            out.println("Skip SP test");
            return;
        }
        out.println("SP test start");
        boolean pass = true;
        try {
            out.println("Stored Proc Test 1");
            TGResultSet<TGEntity> resultSet = conn.executeQuery("gremlin://g.execSP('getAirports', 'United States', 'Seattle');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            int idx = 1;
            while (resultSet.hasNext()) {
                TGEntity entity = (TGEntity) resultSet.next();
                out.printf("%d %s\n", idx++, entity.getAttribute("iataCode").getAsString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("Stored Proc Test 1 ends - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();

        pass = true;
        try {
            out.println("Stored Proc Test 2");
            TGResultSet<TGEntity> resultSet = conn.executeQuery("gremlin://g.execSP('getAirports', 'United States', 'Seattle').values();", null);
            ResultSetUtils.printRSMetaData(resultSet);
            int idx = 1;
            while (resultSet.hasNext()) {
                Object value = resultSet.next();
                out.printf("%d %s\n", idx++, value);
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("Stored Proc Test 2 ends - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();

        pass = true;
        try {
            out.println("Stored Proc Test 3");
            TGResultSet<TGEntity> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('getAirports', 'United States', 'Seattle').values();", null);
            ResultSetUtils.printRSMetaData(resultSet);
            int idx = 1;
            while (resultSet.hasNext()) {
                Object value = resultSet.next();
                out.printf("%d %s\n", idx++, value);
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("Stored Proc Test 3 ends - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();

        pass = true;
        try {
            out.println("Stored Proc Test 4");
            TGResultSet<TGEntity> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('getAirports', 'United States', 'Seattle').out();", null);
            ResultSetUtils.printRSMetaData(resultSet);
            int idx = 1;
            while (resultSet.hasNext()) {
                TGEntity entity = (TGEntity) resultSet.next();
                out.printf("%d %s\n", idx++, entity.getAttribute("iataCode").getAsString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("Stored Proc Test 4 ends - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();

        out.println("SP tests based on api/python/src/tgdb/test/storedproc/testproc.py start");
        pass = true;
        try {
            out.println("retD start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retD');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retD end");
        out.println();

        pass = true;
        try {
            out.println("retN start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retN');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retN end");
        out.println();

        pass = true;
        try {
            out.println("retT start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retT');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retT end");
        out.println();

        pass = true;
        try {
            out.println("retT2D start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retT2D');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retT2D end");
        out.println();

        pass = true;
        try {
            out.println("retTL start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retTL');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retTL end");
        out.println();


        pass = true;
        try {
            out.println("retNestedT start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retNestedT');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retNestedT end");
        out.println();

        pass = true;
        try {
            out.println("retEmptyL start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retEmptyL');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            int i = 0;
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
            if (i == 0) {
                out.println("retEmptyL returns nothing as expected");
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retEmptyL end");
        out.println();

        pass = true;
        try {
            out.println("retLN start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLN');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLN end");
        out.println();

        pass = true;
        try {
            out.println("retLLN start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLLN');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLLN end");
        out.println();

        pass = true;
        try {
            out.println("retLTNl start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLTNl');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLTNl end");
        out.println();

        pass = true;
        try {
            out.println("retLP start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLP');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLP end");
        out.println();

        pass = true;
        try {
            out.println("retLPLs start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLPLs');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLPLs end");
        out.println();

        pass = true;
        try {
            out.println("retLPLsV start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLPLsV');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLPLsV end");
        out.println();

        pass = true;
        try {
            out.println("retLPdL start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retLPdL');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retLPdL end");
        out.println();

        pass = true;
        try {
            out.println("retML start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retML');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retML end");
        out.println();

        pass = true;
        try {
            out.println("retMLP start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retMLP');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retMLP end");
        out.println();

        pass = true;
        try {
            out.println("retMLPL start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SEA').execSP('retMLPL');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("retMLPL end");
        out.println("SP tests based on api/python/src/tgdb/test/storedproc/testproc.py end");
        out.println();

        out.println("SP tests using builtin sp start");
        pass = true;
        try {
            out.println("pageRank start");
            TGResultSet<Object> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'country', 'Germany').execSP('pageRank');", null);
            ResultSetUtils.printRSMetaData(resultSet);
            while (resultSet.hasNext()) {
                Object ent = (Object) resultSet.next();
                out.printf("Value : %s\n", ent.toString());
            }
        }
        catch (TGException e) {
            out.printf("Unexpected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
            pass = false;
        }
        out.println("pageRank end");
        out.println("SP tests using builtin sp end");
        out.println("SP test end");
        out.println();
    }

    private void testDedup() throws Exception {
        if (testAll == false && testDedup == false) {
            out.println("Skip Dedup test");
            return;
        }
        out.println("Dedup test start");
        //This should return a list of primitive values
        /*
        execTest("Test values dedup", 4, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().values("iataCode").dedup());

         */

        execTest("Test out without dedup", 249, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).out("routeType").values("iataCode"));

        execTest("Test out with dedup", 104, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).out("routeType").dedup().values("iataCode"));

        execTest("Test out with dedup - string query", 104, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('airportType').and(__.has('iataCode', P.eq('SFO'))).out('routeType').dedup().values('iataCode');");

        execTest("Test outE without dedup", 249, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).outE("routeType").values("iataCode"));

        execTest("Test outE with dedup", 42, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).outE("routeType").dedup().by("iataCode").values("iataCode"));

        execTest("Test outE with dedup - string query", 42, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('airportType').and(__.has('iataCode', P.eq('SFO'))).outE('routeType').dedup().by('iataCode').values('iataCode');");

        execTest("Test out without dedup by int value", 3, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").out("produces").valueMap());

        execTest("Test outE with dedup by int value but attribute not exists", 0, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("cdi", "cdiid", "172CDIXEAY44").out("produces").dedup().by("groupid").valueMap());

        execTest("Test outE without dedup", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().valueMap());

        execTest("Test outE with dedup by int value", 3, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().dedup().by("groupid").valueMap());

        execTest("Test outE with dedup by string value", 4, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().dedup().by("itemname").valueMap());

        execTest("Test outE with dedup by string value - string query", 4, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('cdibatch').and(__.has('batchid', '17012LXFP342')).outE('contains').inV().dedup().by('itemname').valueMap();");

        execTest("Test string value", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("itemname"));

        execTest("Test string value dedup", 4, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("itemname").dedup());

        execTest("Test int value", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("groupid"));

        execTest("Test int value dedup", 3, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("groupid").dedup());

        execTest("Test int value dedup - string query", 3, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('cdibatch').and(has('batchid', '17012LXFP342')).outE('contains').inV().values('groupid').dedup();");

        //negative case
        execTest("Test values dedup for more than one value - should fail", 0, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().values("iataCode", "city").dedup());

        execTest("Test nodes dedup by values - second 'by' is not supported - should fail", 0, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("airportType").and(__.has("iataCode", P.eq("SFO"))).
                outE("routeType").has("iataCode", "UA").inV().
                outE("routeType").has("iataCode", "UA").inV().dedup().by("iataCode").by("city"));
        out.println("Dedup test end");
        out.println();
    }

    private void testOrderBy() throws Exception {
        if (testAll == false && testOrder == false) {
            out.println("Skip Order By test");
            return;
        }
        out.println("Order by test start");
        execTest("Entity order by string", 249, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().order().by("iataCode").values("iataCode"));

        execTest("Entity order by string with limit", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().order().by("iataCode", Order.decr).limit(10).values("iataCode"));

        execTest("Entity order by int", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().order().by("groupid").valueMap());

        execTest("Entity order by int with limit", 5, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().order().by("groupid", Order.decr).valueMap().limit(5));

        execTest("Value order by string", 249, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().values("iataCode").order());

        execTest("Value order by string with limit", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().values("iataCode").order().by(Order.decr).limit(10));

        execTest("Value order by int", 8, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("groupid").order());

        execTest("Value order by int with limit", 5, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().hasLabel("cdibatch").and(__.has("batchid", "17012LXFP342")).outE("contains").inV().values("groupid").order().by(Order.decr).limit(5));

        execTest("Entity order by string - string query", 249, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().order().by('iataCode').values('iataCode');");

        execTest("Entity order by string with limit - string query", 10, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().order().by('iataCode').limit(10).values('iataCode');");

        execTest("Entity order by int - string query", 8, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('cdibatch').and(has('batchid', '17012LXFP342')).outE('contains').inV().order().by('groupid').values('itemname');");

        execTest("Entity order by int - string query", 5, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('cdibatch').and(has('batchid', '17012LXFP342')).outE('contains').inV().order().by('groupid').limit(5).values('itemname');");

        execTest("Value order by string - string query", 249, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().values('iataCode').order();");

        execTest("Value order by string with limit - string query", 10, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().values('iataCode').order().limit(10);");
        out.println("Order by test end");
        out.println();
    }

    private void testOrderDedup() throws Exception {
        if (testAll == false && testOrderDedup == false) {
            out.println("Skip Order By and Dedup test");
            return;
        }
        out.println("Order by and Dedup test start");

        execTest("Value order by string and dedup with limit", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().values("iataCode").order().dedup().limit(10));

        execTest("Value order by string and dedup", 104, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().values("iataCode").order().dedup());

        execTest("Value order by string desc and dedup", 104, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().values("iataCode").order().by(Order.decr).dedup());

        execTest("Entity order by string and dedup by a different attribute", 101, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().dedup().by("city").order().by("iataCode").valueMap("city","iataCode"));

        execTest("Entity order by id and dedup by string", 104, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().dedup().order().by("iataCode").valueMap("city","iataCode"));

        execTest("Entity order by and dedup by same attribute", 10, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType","iataCode","SFO").out().dedup().by("iataCode").order().by("iataCode").limit(10).valueMap("city","iataCode"));

        execTest("Value order by string and dedup with limit - String", 10, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().values('iataCode').order().dedup().limit(10);");

        execTest("Value order by string and dedup - String", 104, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().values('iataCode').order().dedup();");

        execTest("Value order by string desc and dedup - String", 104, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().values('iataCode').order().by(decr).dedup();");

        execTest("Entity order by string and dedup by a different attribute - String", 101, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().dedup().by('city').order().by('iataCode').valueMap('city','iataCode');");

        execTest("Entity order by id and dedup by string - String", 104, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().dedup().order().by('iataCode').valueMap('city','iataCode');");

        execTest("Entity order by and dedup by same attribute - String", 10, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('airportType','iataCode','SFO').out().dedup().by('iataCode').order().by('iataCode').limit(10).valueMap('city','iataCode');");

        execTest("Entity order by int and dedup by string with limit", 5, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('cdibatch', 'batchid', eq('17012LXFP342').or(eq('17012LXED811'))).outE('contains').inV().order().by('groupid', decr).dedup().by('itemname').valueMap().limit(20);");

        execTest("Test V order by and dedup - string query", 237, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('airportType').dedup().by('country').order().by('country').valueMap('iataCode','country');");

        execTest("Test E order by and dedup - string query", 50, resultListProcessor, resultSizeValidator,
            "gremlin://g.E().hasLabel('routeType').dedup().by('name').order().by('name',asc).limit(50).valueMap();");

        out.println("Order by and Dedup test end");
        out.println();
    }

    private void testVStepWithIds() throws Exception {
        if (testAll == false && testVids == false) {
            out.println("Skip V step with ids test");
            return;
        }
        out.println("V step with ids test start");
        execTest("Test returning vertex ids", 3, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("SEA")).or(P.eq("JFK"))).values("@id"));

        List idList = g.V().has("airportType", "iataCode", P.eq("SFO").or(P.eq("SEA")).or(P.eq("JFK"))).values("@id").toList();

        execTest("Test vertex query with ids", 3, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V(idList.get(0), idList.get(1), idList.get(2)).valueMap());

        execTest("Test vertex traversal with ids", 1875, resultListProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V(idList.get(0), idList.get(1), idList.get(2)).inE().outV().inE().outV().has("iataCode", "YYZ").path().by("iataCode"));

        execTest("Test vertex query with ids and condition", 2, mapInResultProcessor, resultSizeValidator, (DefaultGraphTraversal)
            g.V(idList.get(0), idList.get(1), idList.get(2)).has("iataCode", P.neq("SFO")).valueMap());

        out.println("V step with ids test end");
        out.println("");
    }

    //use joe/joe user to run this test.
    private void testACL() throws Exception {
        if (testAll == false && testACL == false) {
            out.println("Skip ACL test");
            return;
        }
        out.println("ACL test start - need user 'joe' to run this test");
        //FIXME: Need to use a separate user to get these ids first.
        List idList = g.V().has("airlineType", "iataCode", P.eq("UA").or(P.eq("BA")).or(P.eq("AC"))).values("@id").toList();
        //Use joe/joe defined to gdtdb.conf to test this.  With joe/joe it should return nothing with error message.
        if (idList.size() == 3) {
            execTest("Test V query with ids - runtime acl check", 3, mapInResultProcessor, resultSizeValidator, (DefaultGraphTraversal)
                    g.V(idList.get(0), idList.get(1), idList.get(2)).has("iataCode", P.neq("SFO")).valueMap());
        }

        execTest("Test V query - validation time check", 0, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().hasLabel('airlineType').dedup().by('name').order().by('name',asc).limit(50).valueMap();");

        execTest("Test E query - validation time check", 0, resultListProcessor, resultSizeValidator,
            "gremlin://g.E().hasLabel('produces').dedup().by('workcenter').order().by('workcenter',asc).limit(50).valueMap('workcenter','quantity');");

        execTest("Test edge traversal - runtime time check", 0, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('cdi','cdiid','172CBAFDTZ85').outE().inV();");

        execTest("Test node traversal - runtime time check", 0, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('cdi','cdiid','172CBAFDTZ85').in('contains').has('batchid','17012LXEI727');");

        execTest("Test edge traversal - runtime time check", 0, resultListProcessor, resultSizeValidator,
            "gremlin://g.V().has('cdi','cdiid','172CBAFDTZ85').inE('contains').outV();");

        /* to be enabled with new data model
        execTest("Test edge traversal - validation time check", 0, resultListProcessor, resultSizeValidator,
        "g.V().has('workcenterType','workcenter','A528502').outE().inV().has('cdibatch','batchid','17012LXEW018');");

        execTest("Test edge traversal - runtime time check", 0, resultListProcessor, resultSizeValidator,
        "g.V().has('workcenterType','workcenter','A528502').outE().inV().has('machineType','machine','A528512');");
         */
        out.println("ACL test end");
        out.println("");
    }

    void gremlinQueryIllegalSequence1() throws Exception {
        boolean pass = true;
        out.println("Traversal with illegal sequences 1");
        List valueList = g.V().hasLabel("cdi").has("cdiid", "172CDIXDQC18").outE("produces").inV().outE("pduce").out().path().by("cdiid").by("quantity").toList();
        if (valueList.size() > 0) {
            out.println("Unexpected result returned");
            pass = false;
        }
        for (Object value : valueList) {
            out.println(value);
        }
        out.println("Traversal with illegal sequences 1 ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    void gremlinQueryIllegalSequence2() throws Exception {
        boolean pass = true;
        out.println("Traversal with illegal sequences 2");
        List<Edge> nodeList = g.V().has("airportType", "iataCode", "SFO").outE().outE().toList();

        if (nodeList.isEmpty()) {
            out.println("Traverse from 'SFO' node returns nothing - expected");
        } else {
            out.println("Traverse from 'SFO' node returns something - not expected");
            pass = false;
        }
        out.println("Traversal with illegal sequences 2 ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    void gremlinStringQueryUnknownStep() {
        boolean pass = true;
        try {
            out.println("String query with invalid step name 'unknown'");
            TGResultSet<TGEntity> resultSet = conn.executeQuery("gremlin://g.V().has('airportType', 'iataCode', 'SFO').unknown();", null);
            //ResultSetUtils.printRSMetaData(resultSet);
            if (resultSet.hasNext()) {
                out.println("Result set no empty - not expected");
            } else {
                out.println("Expected exception not thrown");
            }
            pass = false;
        }
        catch (TGException e) {
            out.printf("Expected exception(%s) for gremlin string query with error code : %s\n", e.getMessage(), e.getExceptionType());
        }
        out.println("String query with invalid step name  'unknown' ended - " + (pass ? "SUCCEED" : "FAILED"));
        out.println();
    }

    private void futureTestCases() throws Exception {
        //Test for client side validation by gremlib mincore. This is no validation in mincore.
        //valueList = g.V().has("cdi", "groupid", 1000).group().by(T.label).by(__.count()).by("foo").toList();
        //valueList = g.V().group().by(T.label).by(__.count()).by("foo");
        //Column.keys;
        /*
        * Followings are not active cases
        valueList = g.V().has("houseMemberType", "yearBorn", P.between(1906, 1936)).toList();

        valueList = g.V().has("airport","code",P.eq("SFO").or(P.eq("JFK"))).toList();
        valueList = g.V().hasLabel("airport").and(__.or(__.has("code","JFK"),__.has("code", "SFO")),__.has("runways",P.gt(1))).valueMap().toList();

        int startYear = 1846;
        int endYear = 1905;
        valueList = g.V().has("houseMemberType", "yearBorn", P.gt(startYear).and(P.lt(endYear))).toList();

        out.println("Predicate");
        valueList = g.V().has("airport","code",P.eq("SFO").or(P.eq("JFK"))).toList();
        out.println("Predicate ended");


        //Path traversal query
        valueList = g.V().
        		has("airport", "code", "AUS").
        		repeat(__.out().simplePath()).
        		emit().
        		times(5).
        		has("airport", "code", "AGR").
        		path().by("code").
        		limit(10).toList();
        for (Object value : valueList) {
        out.println(value);
        }

        valueList = g.V().
        		has("airport", "code", "AUS").
        		repeat(__.out()).
        		emit().
        		times(5).
        		has("airport", "code", "AGR").
        		path().by("code").
        		limit(10).toList();
        for (Object value : valueList) {
        out.println(value);
        }

        valueList = g.V().has("code","AUS").
        match(__.as("aus").values("runways").as("ausr"),
        __.as("aus").out("route").as("outa").values("runways").as("outr")
        .where("ausr",P.eq("outr"))).
        select("outa").valueMap().select("code","runways").toList();

        //Referencing system id
        g.V().hasId(6).out().has(T.id,P.lt(46)).path().by("code").toList();

        //another way to specify the vertex type/label
        valueList = g.V().where(__.label().is(P.eq("airport"))).count().toList();


        valueList = g.V().hasLabel("airport").has("city", P.between("Dal","Dat")).
        values("city").order().dedup().toList();

        valueList = g.V().has("airport","code","AUS").out().
        not(__.where(__.in("contains").has("code","US"))).
        valueMap("code","city").toList();

        valueList = g.V().has("airport","code","AUS").out().
        where(__.values("runways").is(P.gt(6)).or().values("runways").is(4)).
        valueMap("code","runways").toList();

        valueList = g.V().has("code","AUS").as("a").out().as("b").
        filter(__.select("a","b").by("runways").where("a",P.eq("b"))).
        valueMap("code","runways").toList();

        valueList = g.V().has("airport","city","London").as("a","r").in("contains").as("b").
        where("a",P.eq("b")).by("country").by("code").
        select("a","r","b").by("code").by("region").toList();


        //Bonaparte queries
        valueList = g.V().has("houseMemberType", "yearBorn", P.between(startYear, endYear)).toList();

        valueList = g.V().has("houseMemberType", "yearBorn", P.gt(startYear).and(P.lt(endYear))).toList();

        //In progress, model conv UI query
        valueList = g.V().has("project", "cost", P.gte(500)).in("project_school").in("school_state").has("getstate", "uuid",
        	    		"USA-California").toList();

        //select and airport if outgoing routes > 100 and outgoing route distance > 1000 and the destination airport has > 7 runways
        valueList = g.V().hasLabel("airport").where(
        		__.and(__.out("route").count().is(P.gt(100)),
        			__.outE("route").has("dist", P.gt(1000)).inV().has("runways",P.gt(7)))
        		).values("code").toList();

        //One way of using 'and' step
        valueList = g.V().where(__.outE("created").and().outE("knows")).values("name").toList();

        //select an edge to return
        valueList = g.V().has("code","MIA").outE().as("e").inV().has("code","DFW").select("e").toList();

        //select 10 airports and return result with a count of 10 and props of each of the 10 airports
        valueList = g.V().hasLabel("airport").limit(10).fold().
        		project("count","fields").
        		by(__.unfold().count()).
        		by(__.unfold().project("props").by(__.values("code", "desc")).fold()).toList();

        valueList = g.V().hasLabel("airport").and(__.has("region","US-TX"),__.has("longest",P.gte(12000))).values("code").toList();

        valueList = g.V().has("region","US-TX").has("longest",P.gte(12000)).toList();
        */
    }

    public static void main(String[] args) throws Exception {
        GremlinQueryTest gqt = new GremlinQueryTest();
        gqt.getArgs(args);
        gqt.setup();
        gqt.testValues();
        gqt.testConditions();
        //gqt.testSpecialSteps();
        gqt.testTraversals();
        gqt.testNegativeCases();
        gqt.testAggregations();
        gqt.testEdges();
        gqt.testExceptionHandling();
        gqt.testResultSetMetaData();
        gqt.testGroupSteps();
        gqt.testRepeats();
        gqt.testSP();
        gqt.testByteCodeQuery();
        gqt.testDedup();
        gqt.testOrderBy();
        gqt.testOrderDedup();
        gqt.testVStepWithIds();
        gqt.testACL();
        gqt.cleanup();
    }
}
