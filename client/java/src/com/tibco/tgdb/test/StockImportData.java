/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * File name : StockImportData.${EXT}
 * Created on: 1/13/15
 * Created by: chung 
 * <p/>
 * SVN Id: $Id: ConnectionTest1.java 748 2016-04-25 17:10:38Z vchung $
 */


package com.tibco.tgdb.test;

import java.io.BufferedReader;
import java.io.FileReader;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Date;
import java.util.Iterator;
import java.util.List;
import java.util.Stack;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.connection.TGConnectionFactory;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGEdge;
import com.tibco.tgdb.model.TGGraphMetadata;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.model.TGNode;
import com.tibco.tgdb.model.TGNodeType;

public class StockImportData {
    static String format = "yyyy-MM-dd";

	public String url = "tcp://scott@localhost:8222";
    public String passwd = "scott";
    public TGLogger.TGLevel logLevel = TGLogger.TGLevel.Info;
    public int edgeFetchCount = -1;
    public int nodeFetchCount = -1;
    public int nodeCommitCount = 1000;
    public int edgeCommitCount = 10000;
    
    TGConnection conn = null;
    TGGraphObjectFactory gof;

    TGNodeType yearPriceType = null;
    TGNodeType quarterPriceType = null;
    TGNodeType monthPriceType = null;
    TGNodeType weekPriceType = null;
    TGNodeType dayPriceType = null;
    TGNodeType hourPriceType = null;
    TGNodeType minPriceType = null;

    TGNodeType stockType = null;

    String importFileName = null;
    String stockName = null;
    String companyName = null;
    boolean treatDoubleAsString = false;
    boolean noYearToDayEdge = false;

    void StockImportData() {
    }

    String getStringValue(Iterator<String> argIter) {
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		return s;
    	}
    	return null;
    }
    
    String getStringValue(Iterator<String> argIter, String defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		return s;
    	}
    }

    int getIntValue(Iterator<String> argIter, int defaultValue) {
    	String s = getStringValue(argIter);
    	if (s == null) {
    		return defaultValue;
    	} else {
    		try {
    			int i = Integer.valueOf(s);
    			return i;
    		} catch (NumberFormatException e) {
    			System.out.printf("Invalid number : %s\n", s);
    		}
    		return defaultValue;
    	}
    }

    void getArgs(String[] args) {
    	List<String> argList = Arrays.asList(args);
    	Iterator<String> argIter = argList.iterator();
    	while (argIter.hasNext()) {
    		String s = argIter.next();
    		System.out.printf("Arg : \"%s\"\n", s);
    		if (s.equalsIgnoreCase("-url")) {
    			url = getStringValue(argIter, "tcp://scott@localhost:8222");
    		} else if (s.equalsIgnoreCase("-password") || s.equalsIgnoreCase("-pw")) {
    			passwd = getStringValue(argIter, "scott");
    		} else if (s.equalsIgnoreCase("-loglevel") || s.equalsIgnoreCase("-ll")) {
    			String ll = getStringValue(argIter, "Info");
    			try {
    				logLevel = TGLogger.TGLevel.valueOf(ll);
    			} catch(IllegalArgumentException e) {
    				System.out.printf("Invalid log level value '%s'...ignored\n", ll);
    	        }
            } else if (s.equalsIgnoreCase("-importfile") || s.equalsIgnoreCase("-f")) {
                   importFileName = getStringValue(argIter, "importfile");
            } else if (s.equalsIgnoreCase("-treatdoubleasstring") || s.equalsIgnoreCase("-dtos")) {
                   treatDoubleAsString = true;
            } else if (s.equalsIgnoreCase("-noyeartoday") || s.equalsIgnoreCase("-noy2d")) {
                   noYearToDayEdge = true;
            } else if (s.equalsIgnoreCase("-edgecount") || s.equalsIgnoreCase("-ec")) {
                   edgeFetchCount = getIntValue(argIter, edgeFetchCount);
            } else if (s.equalsIgnoreCase("-nodecount") || s.equalsIgnoreCase("-nc")) {
                   nodeFetchCount = getIntValue(argIter, nodeFetchCount);
               } else if (s.equalsIgnoreCase("-nodecommitcount") || s.equalsIgnoreCase("-ncc")) {
                   nodeCommitCount = getIntValue(argIter, nodeCommitCount);
               } else if (s.equalsIgnoreCase("-edgecommitcount") || s.equalsIgnoreCase("-ecc")) {
                   edgeCommitCount = getIntValue(argIter, edgeCommitCount);
            } else {
                System.out.printf("Skip argument %s\n", s);
            }
        }
    }

    void initMetaData() throws TGException {
        // Nodes
    	TGGraphMetadata gmd = conn.getGraphMetadata(true);
        
        yearPriceType = gmd.getNodeType("yearpricetype");
        quarterPriceType = gmd.getNodeType("quarterpricetype");
        monthPriceType = gmd.getNodeType("monthpricetype");
        weekPriceType = gmd.getNodeType("weekpricetype");
        dayPriceType = gmd.getNodeType("daypricetype");
        hourPriceType = gmd.getNodeType("hourpricetype");
        minPriceType = gmd.getNodeType("minpricetype");

        stockType = gmd.getNodeType("stocktype");
    }

    TGNode createDayRelatedNodes(String[] datepart, String[] stockval) {
    	return null;
        //System.out.printf("create day related nodes not implemented yet");
    }

    void importStockValues() throws Exception {
        String line;
        int currYear = 0;
        int currMonth = 0;
        int currWeek = 0;
        TGNode currStkYearNode = null;
        TGNode currStkMonthNode = null;
        TGNode currStkWeekNode = null;
        TGNode currStkDayNode = null;
        TGNode prevStkYearNode = null;
        TGNode prevStkMonthNode = null;
        TGNode prevStkWeekNode = null;
        TGNode prevStkDayNode = null;

        SimpleDateFormat df = new SimpleDateFormat(format);
        Calendar cal = Calendar.getInstance();

        List<TGNode> dayNodeList = new ArrayList<TGNode>();
        Stack<String> priceStack = new Stack<String>();

        // prepare list to reverse the order of processing
    	try (BufferedReader br = new BufferedReader(new FileReader(importFileName))) {
            //Read the first line
            line = br.readLine();
            if (line == null) {
                return;
            }
            String[] stockInfo = line.split(",");
            companyName = stockInfo[0];
            stockName = stockInfo[1];

            //skip the title line;
            line = br.readLine();

            while ((line = br.readLine()) != null) {
                priceStack.push(line);
            }
        }

        TGNode stockNode = gof.createNode(stockType);
        stockNode.setAttribute("companyname", companyName);
        stockNode.setAttribute("name", stockName);
        conn.insertEntity(stockNode);
        
        double yhprice = 0.0;
        double ylprice = 1000000.0;
        double ycprice = 0.0;
        long yvol = 0;

        double mhprice = 0.0;
        double mlprice = 1000000.0;
        double mcprice = 0.0;
        long mvol = 0;

        double whprice = 0.0;
        double wlprice = 1000000.0;
        double wcprice = 0.0;
        long wvol = 0;

        while(!priceStack.empty()) {

            line = priceStack.pop();

            // process the line.
            String[] stockval = line.split(",");
            if (stockval.length < 7) {
                continue;
            }

            String dateStr = stockval[0];
            String[] datepart = dateStr.split("-");
            if (datepart.length < 3) {
                continue;
            }

            Date date = df.parse(dateStr);
            cal.setTime(date);
 
            int year = cal.get(Calendar.YEAR);
            int month = cal.get(Calendar.MONTH) + 1;
            int dayOfWeek = cal.get(Calendar.DAY_OF_WEEK);
            int dayOfMonth = cal.get(Calendar.DAY_OF_MONTH);
            int dayOfYear = cal.get(Calendar.DAY_OF_YEAR);
            int weekOfYear = cal.get(Calendar.WEEK_OF_YEAR);

            double oprice = Double.valueOf(stockval[1]);
            double hprice = Double.valueOf(stockval[2]);
            double lprice = Double.valueOf(stockval[3]);
            double cprice = Double.valueOf(stockval[4]);
            long vol = Long.valueOf(stockval[5]);

            prevStkDayNode = currStkDayNode;

            currStkDayNode = gof.createNode(dayPriceType);
            currStkDayNode.setAttribute("name", stockName + "-" + dateStr);
            if (treatDoubleAsString == true) {
                currStkDayNode.setAttribute("openprice", String.valueOf(oprice));
                currStkDayNode.setAttribute("highprice", String.valueOf(hprice));
                currStkDayNode.setAttribute("lowprice", String.valueOf(lprice));
                currStkDayNode.setAttribute("closeprice", String.valueOf(cprice));
            } else {
                currStkDayNode.setAttribute("openprice", oprice);
                currStkDayNode.setAttribute("highprice", hprice);
                currStkDayNode.setAttribute("lowprice", lprice);
                currStkDayNode.setAttribute("closeprice", cprice);
            }
            currStkDayNode.setAttribute("tradevolume", vol);
            currStkDayNode.setAttribute("datestring", dateStr);
            currStkDayNode.setAttribute("pricedate", date.getTime());

            conn.insertEntity(currStkDayNode);
            // to be used later for hours and minutes
            dayNodeList.add(currStkDayNode);

            if (prevStkDayNode != null) {
                TGEdge nextDayEdge = gof.createEdge(prevStkDayNode, currStkDayNode, TGEdge.DirectionType.BiDirectional);
                nextDayEdge.setAttribute("name", "NextDay");
                conn.insertEntity(nextDayEdge);
            }

            if (year != currYear) {
                if (currStkYearNode != null) {
                    if (treatDoubleAsString == true) {
                        currStkYearNode.setAttribute("closeprice", String.valueOf(ycprice));
                        currStkYearNode.setAttribute("highprice", String.valueOf(yhprice));
                        currStkYearNode.setAttribute("lowprice", String.valueOf(ylprice));
                    } else {
                        currStkYearNode.setAttribute("closeprice", ycprice);
                        currStkYearNode.setAttribute("highprice", yhprice);
                        currStkYearNode.setAttribute("lowprice", ylprice);
                    }
                    currStkYearNode.setAttribute("tradevolume", yvol);
                }

                prevStkYearNode = currStkYearNode;
                currStkYearNode = gof.createNode(yearPriceType);
                currStkYearNode.setAttribute("name", stockName + "-" + String.valueOf(year));
                if (treatDoubleAsString == true) {
                    currStkYearNode.setAttribute("openprice", String.valueOf(oprice));
                } else {
                    currStkYearNode.setAttribute("openprice", oprice);
                }
                conn.insertEntity(currStkYearNode);

                if (prevStkDayNode != null) {
                	TGEdge nextYearEdge = gof.createEdge(prevStkYearNode, currStkYearNode, TGEdge.DirectionType.BiDirectional);
                	nextYearEdge.setAttribute("name", "NextYear");
                	conn.insertEntity(nextYearEdge);
            	}
                
                TGEdge edge = gof.createEdge(stockNode, currStkYearNode, TGEdge.DirectionType.BiDirectional);
                edge.setAttribute("name", "YearPrice");
                edge.setAttribute("year", year);
                conn.insertEntity(edge);

                currYear = year;
                yhprice = 0.0;
                ylprice = 1000000.0;
                yvol = 0;

            }
            if (lprice < ylprice) {
                ylprice = lprice;
            }
            if (hprice > yhprice) {
                yhprice = hprice;
            }
            ycprice = cprice;
            yvol += vol;
            
            if (noYearToDayEdge == false) {
            	TGEdge year2dayEdge = gof.createEdge(currStkYearNode, currStkDayNode, TGEdge.DirectionType.BiDirectional);
            	year2dayEdge.setAttribute("name", "YearToDay");
            	year2dayEdge.setAttribute("dayofyear", dayOfYear);
            	conn.insertEntity(year2dayEdge);
            }

            if (month != currMonth) {
                if (currStkMonthNode != null) {
                    if (treatDoubleAsString == true) {
                        currStkMonthNode.setAttribute("closeprice", String.valueOf(mcprice));
                        currStkMonthNode.setAttribute("highprice", String.valueOf(mhprice));
                        currStkMonthNode.setAttribute("lowprice", String.valueOf(mlprice));
                    } else {
                        currStkMonthNode.setAttribute("closeprice", mcprice);
                        currStkMonthNode.setAttribute("highprice", mhprice);
                        currStkMonthNode.setAttribute("lowprice", mlprice);
                    }
                    currStkMonthNode.setAttribute("tradevolume", mvol);
                }
                prevStkMonthNode = currStkMonthNode;
                currStkMonthNode = gof.createNode(monthPriceType);
                currStkMonthNode.setAttribute("name", stockName + "-" + String.valueOf(currYear) + "-" + String.valueOf(month));
                if (treatDoubleAsString == true) {
                    currStkMonthNode.setAttribute("openprice", String.valueOf(oprice));
                } else {
                    currStkMonthNode.setAttribute("openprice", oprice);
                }
                conn.insertEntity(currStkMonthNode);

                if (prevStkMonthNode != null) {
                	TGEdge nextMonthEdge = gof.createEdge(prevStkMonthNode, currStkMonthNode, TGEdge.DirectionType.BiDirectional);
                	nextMonthEdge.setAttribute("name", "NextMonth");
                	conn.insertEntity(nextMonthEdge);
            	}

                TGEdge edge = gof.createEdge(stockNode, currStkMonthNode, TGEdge.DirectionType.BiDirectional);
                edge.setAttribute("name", "MonthPrice");
                edge.setAttribute("year", currYear);
                edge.setAttribute("month", month);
                conn.insertEntity(edge);

	            edge = gof.createEdge(currStkYearNode, currStkMonthNode, TGEdge.DirectionType.BiDirectional);
                edge.setAttribute("name", "YearToMonth");
                edge.setAttribute("month", month);
                conn.insertEntity(edge);

                currMonth = month;
                mhprice = 0.0;
                mlprice = 1000000.0;
                mvol = 0;
            }
            if (lprice < mlprice) {
                mlprice = lprice;
            }
            if (hprice > mhprice) {
                mhprice = hprice;
            }
            mcprice = cprice;
            mvol += vol;

            TGEdge month2dayEdge = gof.createEdge(currStkMonthNode, currStkDayNode, TGEdge.DirectionType.BiDirectional);
            month2dayEdge.setAttribute("name", "MonthToDay");
            month2dayEdge.setAttribute("dayofmonth", dayOfMonth);
            conn.insertEntity(month2dayEdge);

            if (weekOfYear != currWeek) {
                if (currStkWeekNode != null) {
                    if (treatDoubleAsString == true) {
                        currStkWeekNode.setAttribute("closeprice", String.valueOf(wcprice));
                        currStkWeekNode.setAttribute("highprice", String.valueOf(whprice));
                        currStkWeekNode.setAttribute("lowprice", String.valueOf(wlprice));
                    } else {
                        currStkWeekNode.setAttribute("closeprice", wcprice);
                        currStkWeekNode.setAttribute("highprice", whprice);
                        currStkWeekNode.setAttribute("lowprice", wlprice);
                    }
                    currStkWeekNode.setAttribute("tradevolume", wvol);
                }
                prevStkWeekNode = currStkWeekNode;
                currStkWeekNode = gof.createNode(weekPriceType);
                if (weekOfYear == 1 && currWeek == 52 && currMonth == 12) {
                	currStkWeekNode.setAttribute("name", stockName + "-" + String.valueOf(currYear + 1) + "-" + String.valueOf(weekOfYear));
                } else {
                	currStkWeekNode.setAttribute("name", stockName + "-" + String.valueOf(currYear) + "-" + String.valueOf(weekOfYear));
                }
                if (treatDoubleAsString == true) {
                    currStkWeekNode.setAttribute("openprice", String.valueOf(oprice));
                } else {
                    currStkWeekNode.setAttribute("openprice", oprice);
                }
                conn.insertEntity(currStkWeekNode);

                if (prevStkWeekNode != null) {
                	TGEdge nextWeekEdge = gof.createEdge(prevStkWeekNode, currStkWeekNode, TGEdge.DirectionType.BiDirectional);
                	nextWeekEdge.setAttribute("name", "NextWeek");
                	conn.insertEntity(nextWeekEdge);
            	}

                TGEdge edge = gof.createEdge(stockNode, currStkWeekNode, TGEdge.DirectionType.BiDirectional);
                edge.setAttribute("name", "WeekPrice");
                edge.setAttribute("year", currYear);
                edge.setAttribute("week", weekOfYear);
                conn.insertEntity(edge);

	            edge = gof.createEdge(currStkYearNode, currStkWeekNode, TGEdge.DirectionType.BiDirectional);
                edge.setAttribute("name", "YearToWeek");
                edge.setAttribute("week", weekOfYear);
                conn.insertEntity(edge);

                currWeek = weekOfYear;
                whprice = 0.0;
                wlprice = 1000000.0;
                wvol = 0;
            }
            if (lprice < wlprice) {
                wlprice = lprice;
            }
            if (hprice > whprice) {
                whprice = hprice;
            }
            wcprice = cprice;
            wvol += vol;

            TGEdge week2dayEdge = gof.createEdge(currStkWeekNode, currStkDayNode, TGEdge.DirectionType.BiDirectional);
            week2dayEdge.setAttribute("name", "WeekToDay");
            week2dayEdge.setAttribute("dayofweek", dayOfWeek);
            conn.insertEntity(week2dayEdge);
        }
        if (treatDoubleAsString == true) {
            currStkYearNode.setAttribute("closeprice", String.valueOf(ycprice));
            currStkYearNode.setAttribute("highprice", String.valueOf(yhprice));
            currStkYearNode.setAttribute("lowprice", String.valueOf(ylprice));

            currStkMonthNode.setAttribute("closeprice", String.valueOf(mcprice));
            currStkMonthNode.setAttribute("highprice", String.valueOf(mhprice));
            currStkMonthNode.setAttribute("lowprice", String.valueOf(mlprice));

            currStkWeekNode.setAttribute("closeprice", String.valueOf(wcprice));
            currStkWeekNode.setAttribute("highprice", String.valueOf(whprice));
            currStkWeekNode.setAttribute("lowprice", String.valueOf(wlprice));

        } else {
            currStkYearNode.setAttribute("closeprice", ycprice);
            currStkYearNode.setAttribute("highprice", yhprice);
            currStkYearNode.setAttribute("lowprice", ylprice);

            currStkMonthNode.setAttribute("closeprice", mcprice);
            currStkMonthNode.setAttribute("highprice", mhprice);
            currStkMonthNode.setAttribute("lowprice", mlprice);

            currStkWeekNode.setAttribute("closeprice", wcprice);
            currStkWeekNode.setAttribute("highprice", whprice);
            currStkWeekNode.setAttribute("lowprice", wlprice);

        }
        currStkYearNode.setAttribute("tradevolume", yvol);
        currStkMonthNode.setAttribute("tradevolume", mvol);
        currStkWeekNode.setAttribute("tradevolume", wvol);
        conn.commit();
        //createDayRelatedNodes(currStkDayNode, cal);
    }

    void run() throws Exception {
    	System.out.printf("Using url : %s, password : %s, log level : %s\n", url, passwd, logLevel.toString());
    	System.out.printf(" max node count : %d, max edge count : %d, node commit count : %d, edge commit count : %d\n",
            nodeFetchCount, edgeFetchCount, nodeCommitCount, edgeCommitCount);
    	TGLogger logger = TGLogManager.getInstance().getLogger();
    	logger.setLevel(logLevel);

        conn = TGConnectionFactory.getInstance().createConnection(url, null, passwd, null);

        conn.connect();

        gof = conn.getGraphObjectFactory();
        
        if (gof == null) {
        	System.out.println("Graph Object Factory is null...exiting");
        	conn.disconnect();
        	return;
        }

        initMetaData();
        importStockValues();

        conn.disconnect();
        System.out.println("Connection test connection disconnected.");
    }

    public static void main(String[] args) throws Exception {
    	StockImportData importData = new StockImportData();
        importData.getArgs(args);
    	importData.run();
    }
}
