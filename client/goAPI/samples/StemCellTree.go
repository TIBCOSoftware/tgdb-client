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
)

var (
	stemTestUrl = "tcp://scott@localhost:8222"
	//stemTestUrl1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//stemTestUrl2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//stemTestUrl3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//stemTestUrl4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//stemTestUrl5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	stemTestPwd = "scott"
	//prefetchMetaData = false
	mESCMap             = make(map[string]types.TGNode, 0)
	hESCMap             = make(map[string]types.TGNode, 0)
	treatDoubleAsString = false
	edgeFetchCount      = -1
	nodeFetchCount      = -1
	nodeCommitCount     = 100
	edgeCommitCount     = 100
	//hEscFile            = "./hESC_mESC.csv"
	//hEscNetworkFile     = "./hESC_comp_network_025.dat"
	//mEscNetworkFile     = "./mESC_comp_network_025.dat"
	hEscFile            = "/Users/achavan/Downloads/genomes/hESC_mESC.csv"
	hEscNetworkFile     = "/Users/achavan/Downloads/genomes/hESC_22k_edge.dat"
	mEscNetworkFile     = "/Users/achavan/Downloads/genomes/mESC_31k_edge.dat"
)

func stemCellUsage() {
	// custom usage (help) output here if needed
	fmt.Println("")
	fmt.Println("Application Flags:")
	flag.PrintDefaults()
	fmt.Println("")
}

func readFileLineByLine(fileName string) (map[int]string, error) {
	fmt.Printf(">>>>>>> Entering readFileLineByLine from '%s'<<<<<<<\n", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(">>>>>>> Returning from readFileLineByLine - error during os.Open(fileName) <<<<<<<")
		return nil, err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	m := make(map[int]string)
	lineCount := 0
	for {
		var buffer bytes.Buffer

		var l []byte
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				//fmt.Println(">>>>>>> Inside readFileLineByLine - breaking since it reached the end of line <<<<<<<")
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				fmt.Println(">>>>>>> Inside readFileLineByLine - breaking since we reached EOF <<<<<<<")
				break
			}
		}

		if err == io.EOF {
			fmt.Println(">>>>>>> Inside readFileLineByLine - breaking since we reached EOF <<<<<<<")
			break
		}

		line := buffer.String()
		m[lineCount] = line
		lineCount++
	}

	if err != io.EOF {
		fmt.Println(">>>>>>> Returning readFileLineByLine - error while reading the line <<<<<<<")
		return nil, err
	}

	fmt.Println(">>>>>>> Returning readFileLineByLine with lineMap <<<<<<<")
	return m, nil
}

func importNodes(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering ImportNodes <<<<<<<")

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	hESCNodeType, err := gmd.GetNodeType("hESC")
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during conn.GetNodeType('hESC') <<<<<<<")
		return
	}

	mESCNodeType, err := gmd.GetNodeType("mESC")
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during conn.GetAttributeDescriptor('mESC') <<<<<<<")
		return
	}

	lineMap, err1 := readFileLineByLine(hEscFile)
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during readFileLineByLine(hEscFile) <<<<<<<")
		return
	}

	count := 0
	fmt.Println(">>>>>>> Entering into the ImportNodes loop to process lineMap <<<<<<<")
	for lineNo, line := range lineMap {
		fieldArray := strings.Split(line, ";")
		fmt.Printf(">>>>>>> Entering into the loop to process lineMap[%d] w/ tokens: '%+v' <<<<<<<\n", lineNo, fieldArray)
		//if (!arr[0].matches("^\\d+$")) {
		if strings.EqualFold(fieldArray[0], "^\\d+$") {
			continue
		}

		// Node
		hESCNode, err := gof.CreateNodeInGraph(hESCNodeType)
		if err != nil {
			fmt.Printf(">>>>>>> Continuing within ImportNodes - error during gof.CreateNodeInGraph(hESCNodeType) at line#'%d'<<<<<<<\n", lineNo)
			continue
		}
		_ = hESCNode.SetOrCreateAttribute("symbol", fieldArray[1])
		_ = hESCNode.SetOrCreateAttribute("name", fieldArray[4])
		hESCMap[fieldArray[0]] = hESCNode
		err = conn.InsertEntity(hESCNode)
		if err != nil {
			fmt.Printf(">>>>>>> Continuing within ImportNodes - error during conn.InsertEntity(hESCNode[%s]) at line#'%d' <<<<<<<\n", fieldArray[4], lineNo)
			continue
		}
		count++

		// Node
		mESCNode, err := gof.CreateNodeInGraph(mESCNodeType)
		if err != nil {
			fmt.Printf(">>>>>>> Continuing within ImportNodes - error during gof.CreateNodeInGraph(mESCNodeType) at line#'%d' <<<<<<<\n", lineNo)
			continue
		}
		_ = mESCNode.SetOrCreateAttribute("symbol", fieldArray[3])
		_ = mESCNode.SetOrCreateAttribute("name", fieldArray[4])
		mESCMap[fieldArray[2]] = mESCNode
		err = conn.InsertEntity(mESCNode)
		if err != nil {
			fmt.Printf(">>>>>>> Continuing within ImportNodes - error during conn.InsertEntity(mESCNode[%s]) at line#'%d' <<<<<<<\n", fieldArray[4], lineNo)
			continue
		}
		count++

		if (count % nodeCommitCount) == 0 {
			_, err = conn.Commit()
			if err != nil {
				fmt.Println(">>>>>>> Breaking from ImportNodes - (count%nodeCommitCount) == 0 <<<<<<<")
				break
			}
		}
		if count == nodeFetchCount {
			fmt.Println(">>>>>>> Breaking from ImportNodes - count == nodeFetchCount <<<<<<<")
			break
		}
	} // End of line loop
	fmt.Println(">>>>>>> Finished looping into the ImportNodes loop to process lineMap <<<<<<<")

	// Last commit for hESC, mESC nodes
	if (count % nodeCommitCount) != 0 {
		_, err = conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Breaking from ImportNodes - (count%nodeCommitCount) != 0 <<<<<<<")
			return
		}
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportNodes - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Printf(">>>>>>> Returning from ImportNodes - Finished processing %d nodes <<<<<<<\n", count)
}

func importHESCEdges(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering ImportHESCEdges <<<<<<<")

	// hESC Edges
	lineMap, err1 := readFileLineByLine(hEscNetworkFile)
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from ImportHESCEdges - error during readFileLineByLine(hEscFile) <<<<<<<")
		return
	}

	count := 0
	fmt.Println(">>>>>>> Entering into the ImportHESCEdges loop to process lineMap <<<<<<<")
	for lineNo, line := range lineMap {
		fieldArray := strings.Split(line, "\t")
		fromNode := hESCMap[fieldArray[0]]
		toNode := hESCMap[fieldArray[1]]
		infScore, err1 := strconv.ParseFloat(fieldArray[2], 64)
		if err1 != nil {
			fmt.Println(">>>>>>> Returning from ImportHESCEdges - error during readFileLineByLine(hEscFile) <<<<<<<")
			return
		}

		if fromNode != nil && toNode != nil {
			// Edge
			edge, err := gof.CreateEdgeWithDirection(fromNode, toNode, types.DirectionTypeBiDirectional)
			if err != nil {
				fmt.Printf(">>>>>>> Returning from ImportHESCEdges - error during gof.CreateEdgeWithDirection(fromNode, toNode) at line#'%d'<<<<<<<\n", lineNo)
				return
			}
			if treatDoubleAsString == true {
				_ = edge.SetOrCreateAttribute("infscore", strconv.FormatFloat(infScore, 'e', 0, 64))
			} else {
				_ = edge.SetOrCreateAttribute("infscore", infScore)
			}
			err = conn.InsertEntity(edge)
			if err != nil {
				fmt.Printf(">>>>>>> Returning from ImportHESCEdges - error during conn.InsertEntity(edge) at line#'%d' <<<<<<<\n", lineNo)
				return
			}
			count++
			if (count % nodeCommitCount) == 0 {
				_, err = conn.Commit()
				if err != nil {
					fmt.Println(">>>>>>> Breaking from ImportHESCEdges - (count%nodeCommitCount) == 0 <<<<<<<")
					break
				}
			}
		}
	} // End of line loop
	fmt.Println(">>>>>>> Finished looping into the ImportHESCEdges loop to process lineMap <<<<<<<")

	// Last commit for hESC nodes
	if (count % nodeCommitCount) != 0 {
		_, err := conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Breaking from ImportHESCEdges - (count%nodeCommitCount) != 0 <<<<<<<")
			return
		}
	}

	_, err := conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportHESCEdges - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Printf(">>>>>>> Returning from ImportHESCEdges - Finished processing %d hESC edges <<<<<<<\n", count)
}

func importMESCEdges(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering ImportMESCEdges <<<<<<<")

	// mESC Edges
	lineMap, err1 := readFileLineByLine(mEscNetworkFile)
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from ImportMESCEdges - error during readFileLineByLine(hEscFile) <<<<<<<")
		return
	}

	count := 0
	fmt.Println(">>>>>>> Entering looping into the ImportMESCEdges loop to process lineMap <<<<<<<")
	for lineNo, line := range lineMap {
		fieldArray := strings.Split(line, "\t")
		fromNode := mESCMap[fieldArray[0]]
		toNode := mESCMap[fieldArray[1]]
		infScore, err1 := strconv.ParseFloat(fieldArray[2], 64)
		if err1 != nil {
			fmt.Println(">>>>>>> Returning from ImportMESCEdges - error during readFileLineByLine(hEscFile) <<<<<<<")
			return
		}

		if fromNode != nil && toNode != nil {
			// Edge
			edge, err := gof.CreateEdgeWithDirection(fromNode, toNode, types.DirectionTypeBiDirectional)
			if err != nil {
				fmt.Printf(">>>>>>> Returning from ImportMESCEdges - error during gof.CreateEdgeWithDirection(fromNode, toNode) at line#'%d'<<<<<<<\n", lineNo)
				return
			}
			if treatDoubleAsString == true {
				_ = edge.SetOrCreateAttribute("infscore", strconv.FormatFloat(infScore, 'f', 0, 64))
			} else {
				_ = edge.SetOrCreateAttribute("infscore", infScore)
			}
			err = conn.InsertEntity(edge)
			if err != nil {
				fmt.Println(">>>>>>> Returning from ImportMESCEdges - error during conn.InsertEntity(edge) <<<<<<<")
				return
			}
			count++
			if (count % nodeCommitCount) == 0 {
				_, err = conn.Commit()
				if err != nil {
					fmt.Println(">>>>>>> Breaking from ImportMESCEdges - (count%nodeCommitCount) == 0 <<<<<<<")
					break
				}
			}
		}
	} // End of line loop
	fmt.Println(">>>>>>> Finished looping into the ImportMESCEdges loop to process lineMap <<<<<<<")

	// Last commit for mESC nodes
	if (count % nodeCommitCount) != 0 {
		_, err := conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Breaking from ImportMESCEdges - (count%nodeCommitCount) != 0 <<<<<<<")
			return
		}
	}

	_, err := conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from ImportMESCEdges - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Printf(">>>>>>> Returning from ImportMESCEdges - Finished processing %d mESC edges <<<<<<<\n", count)
}

func stemCellTest() {
	fmt.Println(">>>>>>> Entering StemCellTest <<<<<<<")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(stemTestUrl, "", stemTestPwd, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from StemCellTest - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from StemCellTest - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from StemCellTest - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from StemCellTest - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside StemCellTest: About to ImportNodes <<<<<<<")
	importNodes(conn, gof)
	fmt.Println(">>>>>>> Inside StemCellTest: About to ImportHESCEdges <<<<<<<")
	importHESCEdges(conn, gof)
	fmt.Println(">>>>>>> Inside StemCellTest: About to ImportMESCEdges <<<<<<<")
	importMESCEdges(conn, gof)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from StemCellTest - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from StemCellTest - successfully disconnected. <<<<<<<")
}

func StemCellMain() {

	// Args parse
	//flag.StringVar(&stemTestUrl, "url", "tcp://scott@localhost:8222", "Please specify TGDB server url.")
	//flag.StringVar(&stemTestPwd, "password", "scott", "Please specify password to connect to TGDB server.")
	//flag.IntVar(&edgeFetchCount, "edgecount", -1, "Please specify number of edges to get for a node, -1 represents NO LIMIT.")
	//flag.IntVar(&edgeCommitCount, "edgecommitcount", -1, "Please specify number of edges to commit for a node. Default is 1000.")
	//flag.IntVar(&nodeFetchCount, "nodecount", -1, "Please specify number of nodes to get , -1 represents NO LIMIT.")
	//flag.IntVar(&nodeCommitCount, "nodecommitcount", -1, "Please specify number of nodes to commit. Default is 1000.")
	//flag.BoolVar(&treatDoubleAsString, "treatdoubleasstring", false, "Please specify whether to treat Double as a String or not.")
	//flag.StringVar(&hEscFile, "hESC", "./hESC_mESC.csv", "Please specify hESC data file name. Default is hESC_mESC.csv.")
	//flag.StringVar(&hEscNetworkFile, "hESCNET", "./hESC_comp_network_025.dat", "Please specify hESC edge data file name. Default is hESC_comp_network_025.dat.")
	//flag.StringVar(&mEscNetworkFile, "mESCNET", "./mESC_comp_network_025.dat", "Please specify mESC edge data file name. Default is mESC_comp_network_025.dat.")
	//
	//flag.Parse()
	//
	//// assign custom usage function (will be shown by default if -h or --help flag is passed)
	//flag.Usage = stemCellUsage
	//
	//// if no flags print usage (not default behaviour)
	//if len(os.Args) == 1 {
	//	stemCellUsage()
	//}

	fmt.Println(">>>>>>> Starting Stem Cell Data Import ... <<<<<<<")
	stemCellTest()
	fmt.Println(">>>>>>> Stem Cell Data Imported Successfully!!! <<<<<<<")
}
