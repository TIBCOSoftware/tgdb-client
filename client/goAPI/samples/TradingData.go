package samples

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	tradingDataUrl = "tcp://scott@localhost:8222"
	//tradingDataUrl1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//tradingDataUrl2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//tradingDataUrl3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//tradingDataUrl4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//tradingDataUrl5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	tradingDataPwd = "scott"
	//prefetchMetaData = false
	format          = "2006-01-02"
	companyName     = "Coke"
	stockName       = "KO"
	doubleAsString  = false
	edgeCount       = -1
	nodeCount       = -1
	nodeCommCount   = 1000
	edgeCommCount   = 1000
	importFileName  = "/Users/achavan/Downloads/stocks/KO.csv"
	noYearToDayEdge = false
)

// Stack is a basic LIFO stack that re-sizes as needed.
type Stack struct {
	lines []string
	count int
}

// NewStack returns a new stack.
func NewStack() *Stack {
	return &Stack{}
}

// Push adds a node to the stack.
func (s *Stack) Push(n string) {
	s.lines = append(s.lines[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() string {
	if s.count == 0 {
		return ""
	}
	s.count--
	return s.lines[s.count]
}

var priceDataStack *Stack

func tradingDataUsage() {
	// custom usage (help) output here if needed
	fmt.Println("")
	fmt.Println("Application Flags:")
	flag.PrintDefaults()
	fmt.Println("")
}

func prepareDataStackForProcessing(fileName string) (*Stack, error) {
	fmt.Printf(">>>>>>> Entering prepareDataStackForProcessing from data file: '%s' <<<<<<<\n", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(">>>>>>> Returning from prepareDataStackForProcessing - error during os.Open(fileName) <<<<<<<")
		return nil, err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	priceDataStack = NewStack()
	lineMap := make(map[int]string)
	lineCount := 1
	for {
		var buffer bytes.Buffer

		var l []byte
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				//fmt.Println(">>>>>>> Inside prepareDataStackForProcessing - breaking since it reached the end of line <<<<<<<")
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				fmt.Println(">>>>>>> Inside prepareDataStackForProcessing - breaking since we reached EOF <<<<<<<")
				break
			}
		}

		if err == io.EOF {
			fmt.Println(">>>>>>> Inside prepareDataStackForProcessing - breaking since we reached EOF <<<<<<<")
			break
		}

		line := buffer.String()
		lineMap[lineCount] = line
		lineCount++
	}

	if err != io.EOF {
		fmt.Println(">>>>>>> Returning prepareDataStackForProcessing - error while reading the line <<<<<<<")
		return nil, err
	}

	// Prepare list to reverse the order of processing
	for lineNo, line := range lineMap {
		if line == "" {
			continue
		}
		//Read the first line
		if lineNo == 1 {
		//	stockInfo := strings.Split(line, ",")
		//	companyName = stockInfo[0]
		//	stockName = stockInfo[1]
		//}
		////Skip the title line;
		//if lineNo == 2 {
			continue
		}
		priceDataStack.Push(line)
	}

	fmt.Printf(">>>>>>> Returning prepareDataStackForProcessing with priceDataStack: '%+v' <<<<<<<\n", priceDataStack)
	return priceDataStack, nil
}

func importStockValues(conn types.TGConnection, gof types.TGGraphObjectFactory, priceDataStack *Stack) {
	fmt.Println(">>>>>>> Entering importStockValues <<<<<<<")

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	// Nodes
	//yearPriceType, err := gmd.GetNodeType("yearpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetNodeType('yearpricetype') <<<<<<<")
	//	return
	//}
	//
	//quarterPriceType, err := gmd.GetNodeType("quarterpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetAttributeDescriptor('quarterpricetype') <<<<<<<")
	//	return
	//}
	//
	//monthPriceType, err := gmd.GetNodeType("monthpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetNodeType('monthpricetype') <<<<<<<")
	//	return
	//}
	//
	//weekPriceType, err := gmd.GetNodeType("weekpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetAttributeDescriptor('weekpricetype') <<<<<<<")
	//	return
	//}

	dayPriceType, err := gmd.GetNodeType("daypricetype")
	if err != nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetAttributeDescriptor('daypricetype') <<<<<<<")
		return
	}

	//hourPriceType, err := gmd.GetNodeType("hourpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetNodeType('hourpricetype') <<<<<<<")
	//	return
	//}
	//
	//minPriceType, err := gmd.GetNodeType("minpricetype")
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetAttributeDescriptor('minpricetype') <<<<<<<")
	//	return
	//}

	stockType, err := gmd.GetNodeType("stocktype")
	if err != nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetNodeType('stocktype') <<<<<<<")
		return
	}


	// Node
	stockNode, err := gof.CreateNodeInGraph(stockType)
	if err != nil {
		fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.CreateNodeInGraph(stockType) <<<<<<<")
		return
	}
	_ = stockNode.SetOrCreateAttribute("name", stockName)
	//_ = stockNode.SetOrCreateAttribute("companyname", companyName)
	err = conn.InsertEntity(stockNode)
	if err != nil {
		fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.InsertEntity(stockNode) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Inside tradingDataTest - error during conn.Commit() <<<<<<<")
		return
	}

	//yhprice := 0.0
	//ylprice := 1000000.0
	//ycprice := 0.0
	//var yvol int64
	//
	//mhprice := 0.0
	//mlprice := 1000000.0
	//mcprice := 0.0
	//var mvol int64
	//
	//whprice := 0.0
	//wlprice := 1000000.0
	//wcprice := 0.0
	//var wvol int64
	//
	//currYear := 0
	//currMonth := 0
	//currWeek := 0
	var /*currStkYearNode, currStkMonthNode, currStkWeekNode, */currStkDayNode types.TGNode
	var /*prevStkYearNode, prevStkMonthNode, prevStkWeekNode, */prevStkDayNode types.TGNode

	dayNodeList := make([]types.TGNode, 0)

	//currStkDayNode, err = gof.CreateNodeInGraph(dayPriceType)
	//if err != nil {
	//	fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.CreateNodeInGraph(dayPriceType) <<<<<<<")
	//	return
	//}
	//
	//currStkWeekNode, err = gof.CreateNodeInGraph(weekPriceType)
	//if err != nil {
	//	fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.CreateNodeInGraph(weekPriceType) <<<<<<<")
	//	return
	//}
	//
	//currStkMonthNode, err = gof.CreateNodeInGraph(monthPriceType)
	//if err != nil {
	//	fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.CreateNodeInGraph(monthPriceType) <<<<<<<")
	//	return
	//}
	//
	//currStkYearNode, err = gof.CreateNodeInGraph(yearPriceType)
	//if err != nil {
	//	fmt.Println(">>>>>>> Continuing within tradingDataTest - error during conn.CreateNodeInGraph(yearPriceType) <<<<<<<")
	//	return
	//}

	lineCount := 1

	fmt.Println(">>>>>>> Entering into the loop to process priceDataStack <<<<<<<")
	for {
		line := priceDataStack.Pop()
		if line == "" {
			break
		}
		fmt.Printf(">>>>>>> Inside the loop within tradingDataTest processing line#'%d':%s <<<<<<<\n", priceDataStack.count+1, line)

		// Process the line.
		stockVal := strings.Split(line, ",")
		if len(stockVal) < 7 {
			continue
		}

		dateStr := stockVal[0]
		datepart := strings.Split(dateStr, "-")
		if len(datepart) < 3 {
			continue
		}

		tod, err := time.Parse(format, dateStr)
		if err != nil {
			fmt.Printf(">>>>>>> Inside the loop within tradingDataTest - error during time.Parse(%s, %s) <<<<<<<\n", format, dateStr)
			return
		}
		fmt.Printf(">>>>>>> Inside the loop within tradingDataTest - stock transaction date: '%+v' <<<<<<<\n", tod)
		//year, month, dayOfMonth := tod.Date()
		//dayOfYear := tod.YearDay()
		//dayOfWeek := tod.Weekday()
		//_, weekOfYear := tod.ISOWeek()

		oprice, _ := strconv.ParseFloat(stockVal[1], 64)
		hprice, _ := strconv.ParseFloat(stockVal[2], 64)
		lprice, _ := strconv.ParseFloat(stockVal[3], 64)
		cprice, _ := strconv.ParseFloat(stockVal[4], 64)
		//adprice, _ := strconv.ParseFloat(stockVal[5], 64)
		//vol, _ := strconv.ParseInt(stockVal[5], 10, 64)
		vol, _ := strconv.ParseInt(stockVal[6], 10, 64)

		prevStkDayNode = currStkDayNode

		// Node
		currStkDayNode, err = gof.CreateNodeInGraph(dayPriceType)
		if err != nil {
			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.CreateNodeInGraph(dayPriceType) <<<<<<<")
			return
		}
		_ = currStkDayNode.SetOrCreateAttribute("name", stockName+"-"+dateStr)
		if doubleAsString == true {
			_ = currStkDayNode.SetOrCreateAttribute("openprice", stockVal[1])
			_ = currStkDayNode.SetOrCreateAttribute("highprice", stockVal[2])
			_ = currStkDayNode.SetOrCreateAttribute("lowprice", stockVal[3])
			_ = currStkDayNode.SetOrCreateAttribute("closeprice", stockVal[4])
		} else {
			_ = currStkDayNode.SetOrCreateAttribute("openprice", oprice)
			_ = currStkDayNode.SetOrCreateAttribute("highprice", hprice)
			_ = currStkDayNode.SetOrCreateAttribute("lowprice", lprice)
			_ = currStkDayNode.SetOrCreateAttribute("closeprice", cprice)
		}
		_ = currStkDayNode.SetOrCreateAttribute("tradevolume", vol)
		_ = currStkDayNode.SetOrCreateAttribute("datestring", dateStr)
		_ = currStkDayNode.SetOrCreateAttribute("pricedate", tod.Unix())

		err = conn.InsertEntity(currStkDayNode)
		if err != nil {
			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(currStkDayNode) <<<<<<<")
			return
		}

		// to be used later for hours and minutes
		dayNodeList = append(dayNodeList, currStkDayNode)

		if prevStkDayNode != nil {
			// Edge
			nextDayEdge, err := gof.CreateEdgeWithDirection(prevStkDayNode, currStkDayNode, types.DirectionTypeBiDirectional)
			if err != nil {
				fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(prevStkDayNode, currStkDayNode) <<<<<<<")
				return
			}
			_ = nextDayEdge.SetOrCreateAttribute("name", "NextDay")
			err = conn.InsertEntity(nextDayEdge)
			if err != nil {
				fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(nextDayEdge) <<<<<<<")
				return
			}
		}

		//if year != currYear {
		//	if currStkYearNode != nil {
		//		if doubleAsString == true {
		//			_ = currStkYearNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(yhprice, 'f', 0, 64))
		//			_ = currStkYearNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(ylprice, 'f', 0, 64))
		//			_ = currStkYearNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(ycprice, 'f', 0, 64))
		//		} else {
		//			_ = currStkYearNode.SetOrCreateAttribute("highprice", yhprice)
		//			_ = currStkYearNode.SetOrCreateAttribute("lowprice", ylprice)
		//			_ = currStkYearNode.SetOrCreateAttribute("closeprice", ycprice)
		//		}
		//		_ = currStkYearNode.SetOrCreateAttribute("tradevolume", yvol)
		//	}
		//
		//	prevStkYearNode = currStkYearNode
		//
		//	// Node
		//	currStkYearNode, err = gof.CreateNodeInGraph(yearPriceType)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.CreateNodeInGraph(yearPriceType) <<<<<<<")
		//		return
		//	}
		//	_ = currStkYearNode.SetOrCreateAttribute("name", stockName+"-"+strconv.Itoa(year))
		//	if doubleAsString == true {
		//		_ = currStkYearNode.SetOrCreateAttribute("openprice", strconv.FormatFloat(oprice, 'f', 0, 64))
		//	} else {
		//		_ = currStkYearNode.SetOrCreateAttribute("openprice", oprice)
		//	}
		//	err = conn.InsertEntity(currStkYearNode)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(currStkYearNode) <<<<<<<")
		//		return
		//	}
		//
		//	if prevStkDayNode != nil {
		//		// Edge
		//		nextYearEdge, err := gof.CreateEdgeWithDirection(prevStkYearNode, currStkYearNode, types.DirectionTypeBiDirectional)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(prevStkYearNode, currStkYearNod) <<<<<<<")
		//			return
		//		}
		//		_ = nextYearEdge.SetOrCreateAttribute("name", "NextYear")
		//		err = conn.InsertEntity(nextYearEdge)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(nextYearEdge) <<<<<<<")
		//			return
		//		}
		//	}
		//
		//	// Edge
		//	edge, err := gof.CreateEdgeWithDirection(stockNode, currStkYearNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(prevStkYearNode, currStkYearNod) <<<<<<<")
		//		return
		//	}
		//	_ = edge.SetOrCreateAttribute("name", "YearPrice")
		//	_ = edge.SetOrCreateAttribute("year", year)
		//	err = conn.InsertEntity(edge)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(nextYearEdge) <<<<<<<")
		//		return
		//	}
		//
		//	currYear = year
		//	yhprice = 0.0
		//	ylprice = 1000000.0
		//	yvol = 0
		//}
		//
		//if lprice < ylprice {
		//	ylprice = lprice
		//}
		//if hprice > yhprice {
		//	yhprice = hprice
		//}
		//ycprice = cprice
		//yvol += vol
		//
		//if noYearToDayEdge == false {
		//	// Edge
		//	year2dayEdge, err := gof.CreateEdgeWithDirection(currStkYearNode, currStkDayNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(currStkYearNode, currStkDayNode) <<<<<<<")
		//		return
		//	}
		//	_ = year2dayEdge.SetOrCreateAttribute("name", "YearToDay")
		//	_ = year2dayEdge.SetOrCreateAttribute("dayofyear", dayOfYear)
		//	err = conn.InsertEntity(year2dayEdge)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(year2dayEdge) <<<<<<<")
		//		return
		//	}
		//}
		//
		//if month != time.Month(currMonth) {
		//	if currStkMonthNode != nil {
		//		if doubleAsString == true {
		//			_ = currStkMonthNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(mhprice, 'f', 0, 64))
		//			_ = currStkMonthNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(mlprice, 'f', 0, 64))
		//			_ = currStkMonthNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(mcprice, 'f', 0, 64))
		//		} else {
		//			_ = currStkMonthNode.SetOrCreateAttribute("highprice", mhprice)
		//			_ = currStkMonthNode.SetOrCreateAttribute("lowprice", mlprice)
		//			_ = currStkMonthNode.SetOrCreateAttribute("closeprice", mcprice)
		//		}
		//		_ = currStkMonthNode.SetOrCreateAttribute("tradevolume", mvol)
		//	}
		//	prevStkMonthNode = currStkMonthNode
		//
		//	// Node
		//	currStkMonthNode, err = gof.CreateNodeInGraph(monthPriceType)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.CreateNodeInGraph(monthPriceType) <<<<<<<")
		//		return
		//	}
		//	_ = currStkMonthNode.SetOrCreateAttribute("name", stockName+"-"+strconv.Itoa(currYear)+"-"+strconv.Itoa(int(month)))
		//	if doubleAsString == true {
		//		_ = currStkMonthNode.SetOrCreateAttribute("openprice", strconv.FormatFloat(oprice, 'f', 0, 64))
		//	} else {
		//		_ = currStkMonthNode.SetOrCreateAttribute("openprice", oprice)
		//	}
		//	err = conn.InsertEntity(currStkMonthNode)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(currStkMonthNode) <<<<<<<")
		//		return
		//	}
		//
		//	if prevStkMonthNode != nil {
		//		// Edge
		//		nextMonthEdge, err := gof.CreateEdgeWithDirection(prevStkMonthNode, currStkMonthNode, types.DirectionTypeBiDirectional)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(prevStkMonthNode, currStkMonthNode) <<<<<<<")
		//			return
		//		}
		//		_ = nextMonthEdge.SetOrCreateAttribute("name", "NextMonth")
		//		err = conn.InsertEntity(nextMonthEdge)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(nextMonthEdge) <<<<<<<")
		//			return
		//		}
		//	}
		//
		//	// Edge
		//	edge, err := gof.CreateEdgeWithDirection(stockNode, currStkMonthNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(stockNode, currStkMonthNode) <<<<<<<")
		//		return
		//	}
		//	_ = edge.SetOrCreateAttribute("name", "MonthPrice")
		//	_ = edge.SetOrCreateAttribute("year", currYear)
		//	_ = edge.SetOrCreateAttribute("month", month)
		//	err = conn.InsertEntity(edge)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(edge) <<<<<<<")
		//		return
		//	}
		//
		//	// Edge
		//	edge1, err := gof.CreateEdgeWithDirection(currStkYearNode, currStkMonthNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(stockNode, currStkMonthNode) <<<<<<<")
		//		return
		//	}
		//	_ = edge1.SetOrCreateAttribute("name", "YearToMonth")
		//	_ = edge1.SetOrCreateAttribute("month", month)
		//	err = conn.InsertEntity(edge1)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(edge1) <<<<<<<")
		//		return
		//	}
		//
		//	currMonth = int(month)
		//	mhprice = 0.0
		//	mlprice = 1000000.0
		//	mvol = 0
		//}
		//if lprice < mlprice {
		//	mlprice = lprice
		//}
		//if hprice > mhprice {
		//	mhprice = hprice
		//}
		//mcprice = cprice
		//mvol += vol
		//
		//// Edge
		//month2dayEdge, err := gof.CreateEdgeWithDirection(currStkMonthNode, currStkDayNode, types.DirectionTypeBiDirectional)
		//if err != nil {
		//	fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(currStkMonthNode, currStkDayNode) <<<<<<<")
		//	return
		//}
		//_ = month2dayEdge.SetOrCreateAttribute("name", "MonthToDay")
		//_ = month2dayEdge.SetOrCreateAttribute("dayofmonth", dayOfMonth)
		//err = conn.InsertEntity(month2dayEdge)
		//if err != nil {
		//	fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(edge1) <<<<<<<")
		//	return
		//}
		//
		//if weekOfYear != currWeek {
		//	if currStkWeekNode != nil {
		//		if doubleAsString == true {
		//			_ = currStkWeekNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(whprice, 'f', 0, 64))
		//			_ = currStkWeekNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(wlprice, 'f', 0, 64))
		//			_ = currStkWeekNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(wcprice, 'f', 0, 64))
		//		} else {
		//			_ = currStkWeekNode.SetOrCreateAttribute("highprice", whprice)
		//			_ = currStkWeekNode.SetOrCreateAttribute("lowprice", wlprice)
		//			_ = currStkWeekNode.SetOrCreateAttribute("closeprice", wcprice)
		//		}
		//		_ = currStkWeekNode.SetOrCreateAttribute("tradevolume", wvol)
		//	}
		//	prevStkWeekNode = currStkWeekNode
		//
		//	// Node
		//	currStkWeekNode, err = gof.CreateNodeInGraph(weekPriceType)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.CreateNodeInGraph(weekPriceType) <<<<<<<")
		//		return
		//	}
		//	if weekOfYear == 1 && currWeek == 52 && currMonth == 12 {
		//		_ = currStkWeekNode.SetOrCreateAttribute("name", stockName+"-"+strconv.Itoa(currYear+1)+"-"+strconv.Itoa(weekOfYear))
		//	} else {
		//		_ = currStkWeekNode.SetOrCreateAttribute("name", stockName+"-"+strconv.Itoa(currYear)+"-"+strconv.Itoa(weekOfYear))
		//	}
		//	if doubleAsString == true {
		//		_ = currStkWeekNode.SetOrCreateAttribute("openprice", strconv.FormatFloat(oprice, 'f', 0, 64))
		//	} else {
		//		_ = currStkWeekNode.SetOrCreateAttribute("openprice", oprice)
		//	}
		//	err = conn.InsertEntity(currStkWeekNode)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(currStkWeekNode) <<<<<<<")
		//		return
		//	}
		//
		//	if prevStkWeekNode != nil {
		//		// Edge
		//		nextWeekEdge, err := gof.CreateEdgeWithDirection(prevStkWeekNode, currStkWeekNode, types.DirectionTypeBiDirectional)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(prevStkWeekNode, currStkWeekNode) <<<<<<<")
		//			return
		//		}
		//		_ = nextWeekEdge.SetOrCreateAttribute("name", "NextWeek")
		//		err = conn.InsertEntity(nextWeekEdge)
		//		if err != nil {
		//			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(nextWeekEdge) <<<<<<<")
		//			return
		//		}
		//	}
		//
		//	// Edge
		//	edge, err := gof.CreateEdgeWithDirection(stockNode, currStkWeekNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(stockNode, currStkWeekNode) <<<<<<<")
		//		return
		//	}
		//	_ = edge.SetOrCreateAttribute("name", "WeekPrice")
		//	_ = edge.SetOrCreateAttribute("year", currYear)
		//	_ = edge.SetOrCreateAttribute("week", weekOfYear)
		//	err = conn.InsertEntity(edge)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(edge) <<<<<<<")
		//		return
		//	}
		//
		//	// Edge
		//	edge2, err := gof.CreateEdgeWithDirection(currStkYearNode, currStkWeekNode, types.DirectionTypeBiDirectional)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(currStkYearNode, currStkWeekNode) <<<<<<<")
		//		return
		//	}
		//	_ = edge2.SetOrCreateAttribute("name", "YearToWeek")
		//	_ = edge2.SetOrCreateAttribute("week", weekOfYear)
		//	err = conn.InsertEntity(edge2)
		//	if err != nil {
		//		fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(edge2) <<<<<<<")
		//		return
		//	}
		//
		//	currWeek = weekOfYear
		//	whprice = 0.0
		//	wlprice = 1000000.0
		//	wvol = 0
		//}
		//if lprice < wlprice {
		//	wlprice = lprice
		//}
		//if hprice > whprice {
		//	whprice = hprice
		//}
		//wcprice = cprice
		//wvol += vol
		//
		//// Edge
		//week2dayEdge, err := gof.CreateEdgeWithDirection(currStkWeekNode, currStkDayNode, types.DirectionTypeBiDirectional)
		//if err != nil {
		//	fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during gof.CreateEdgeWithDirection(currStkWeekNode, currStkDayNode) <<<<<<<")
		//	return
		//}
		//_ = week2dayEdge.SetOrCreateAttribute("name", "WeekToDay")
		//_ = week2dayEdge.SetOrCreateAttribute("dayofweek", dayOfWeek)
		//err = conn.InsertEntity(week2dayEdge)
		//if err != nil {
		//	fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.InsertEntity(week2dayEdge) <<<<<<<")
		//	return
		//}
		//
		_, err = conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Inside the loop within tradingDataTest - error during conn.Commit() <<<<<<<")
			return
		}

		lineCount++
	} // End of for loop
	fmt.Println(">>>>>>> Finished looping in the loop to process priceDataStack <<<<<<<")

	//if doubleAsString == true {
	//	_ = currStkYearNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(yhprice, 'f', 0, 64))
	//	_ = currStkYearNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(ylprice, 'f', 0, 64))
	//	_ = currStkYearNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(ycprice, 'f', 0, 64))
	//
	//	_ = currStkMonthNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(mhprice, 'f', 0, 64))
	//	_ = currStkMonthNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(mlprice, 'f', 0, 64))
	//	_ = currStkMonthNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(mcprice, 'f', 0, 64))
	//
	//	_ = currStkWeekNode.SetOrCreateAttribute("highprice", strconv.FormatFloat(whprice, 'f', 0, 64))
	//	_ = currStkWeekNode.SetOrCreateAttribute("lowprice", strconv.FormatFloat(wlprice, 'f', 0, 64))
	//	_ = currStkWeekNode.SetOrCreateAttribute("closeprice", strconv.FormatFloat(wcprice, 'f', 0, 64))
	//} else {
	//	_ = currStkYearNode.SetOrCreateAttribute("highprice", yhprice)
	//	_ = currStkYearNode.SetOrCreateAttribute("lowprice", ylprice)
	//	_ = currStkYearNode.SetOrCreateAttribute("closeprice", ycprice)
	//
	//	_ = currStkMonthNode.SetOrCreateAttribute("highprice", mhprice)
	//	_ = currStkMonthNode.SetOrCreateAttribute("lowprice", mlprice)
	//	_ = currStkMonthNode.SetOrCreateAttribute("closeprice", mcprice)
	//
	//	_ = currStkWeekNode.SetOrCreateAttribute("highprice", whprice)
	//	_ = currStkWeekNode.SetOrCreateAttribute("lowprice", wlprice)
	//	_ = currStkWeekNode.SetOrCreateAttribute("closeprice", wcprice)
	//}
	//_ = currStkYearNode.SetOrCreateAttribute("tradevolume", yvol)
	//_ = currStkMonthNode.SetOrCreateAttribute("tradevolume", mvol)
	//_ = currStkWeekNode.SetOrCreateAttribute("tradevolume", wvol)

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning from importStockValues - Finished processing <<<<<<<")
}

func tradingDataTest() {
	fmt.Println(">>>>>>> Entering tradingDataTest")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(tradingDataUrl, "", tradingDataPwd, nil)
	if err != nil {
		fmt.Println("Returning from tradingDataTest - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from tradingDataTest - error during conn.Connect")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside tradingDataTest: About to importStockValues <<<<<<<")
	priceDataStack, err1 := prepareDataStackForProcessing(importFileName)
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from tradingDataTest - error during prepareDataStackForProcessing(importFileName) <<<<<<<")
		return
	}
	fmt.Printf(">>>>>>> Inside tradingDataTest: Gathered trading data as: '%+v' <<<<<<<\n", priceDataStack)

	fmt.Println(">>>>>>> Inside tradingDataTest: About to importStockValues <<<<<<<")
	importStockValues(conn, gof, priceDataStack)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from tradingDataTest - error during conn.Disconnect")
		return
	}
	fmt.Println(">>>>>>> Returning from tradingDataTest - successfully disconnected. <<<<<<<")
}

func TradingDataMain() {

	//// Args parse
	//flag.StringVar(&tradingDataUrl, "url", "tcp://scott@localhost:8222", "Please specify TGDB server url.")
	//flag.StringVar(&tradingDataPwd, "password", "scott", "Please specify password to connect to TGDB server.")
	//flag.StringVar(&importFileName, "importfile", "/Users/achavan/Downloads/stocks/KO.csv", "Please specify data file name.")
	//flag.IntVar(&edgeCount, "edgecount", -1, "Please specify number of edges to get , -1 represents NO LIMIT.")
	//flag.IntVar(&edgeCommCount, "edgecommitcount", -1, "Please specify number of edges to commit for a node. Default is 1000.")
	//flag.IntVar(&nodeCount, "nodecount", -1, "Please specify number of nodes to get , -1 represents NO LIMIT.")
	//flag.IntVar(&nodeCommCount, "nodecommitcount", -1, "Please specify number of nodes to commit. Default is 1000.")
	//flag.BoolVar(&doubleAsString, "treatdoubleasstring", false, "Please specify whether to treat Double as a String or not.")
	//flag.BoolVar(&noYearToDayEdge, "noyeartoday", false, "Please specify whether there is an edge relationship from year to day.")
	//
	//flag.Parse()
	//
	//// assign custom usage function (will be shown by default if -h or --help flag is passed)
	//flag.Usage = tradingDataUsage
	//
	//// if no flags print usage (not default behaviour)
	//if len(os.Args) == 1 {
	//	tradingDataUsage()
	//}

	fmt.Println(">>>>>>> Starting Trading Data Import ... <<<<<<<")
	tradingDataTest()
	fmt.Println(">>>>>>> Trading Data Imported Successfully!!! <<<<<<<")
}
