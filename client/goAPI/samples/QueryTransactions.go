package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	qUrl = "tcp://scott@localhost:8222"
	//qUrl1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//qUrl2 = "tcp://scott@localhost:8222/{qConnectTimeout=30}";
	//qUrl3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//qUrl4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//qUrl5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	qPassword = "scott"
	//prefetchMetaData = false
)

func connect() (types.TGConnection, types.TGError) {
	qConnFactory := connection.NewTGConnectionFactory()
	qConn, err := qConnFactory.CreateConnection(qUrl, "", qPassword, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from connect - error during CreateConnection <<<<<<<")
		return nil, err
	}

	err = qConn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from connect - error during conn.Connect <<<<<<<")
		return nil, err
	}
	return qConn, nil
}

func disConnect(qConn types.TGConnection) {
	err := qConn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from disConnect - error during qConn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from disConnect - successfully disconnected. <<<<<<<")
}

func createNode(gof types.TGGraphObjectFactory, nodeType types.TGNodeType) (types.TGNode, types.TGError) {
	if nodeType != nil {
		return gof.CreateNodeInGraph(nodeType)
	} else {
		return gof.CreateNode()
	}
}

func createTestData(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering createTestData: Creating Test Data <<<<<<<")
	gof, err := qConn.GetGraphObjectFactory()
	if err != nil || gof == nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.GetGraphObjectFactory <<<<<<<")
		return err
	}

	gmd, err := qConn.GetGraphMetadata(true)
	if err != nil {
		errMsg := fmt.Sprint("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetGraphMetadata")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println("Returning from createTestData - error during gmd.GetNodeType('testnode')")
		return err
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'testnode' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		errMsg := fmt.Sprint(">>>>>>> 'testnode' is not found from meta data fetch <<<<<<<")
		return exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Node # 1
	node1, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node1) <<<<<<<")
		return err
	}
	_ = node1.SetOrCreateAttribute("name", "Bruce-Wayne")
	_ = node1.SetOrCreateAttribute("multiple", 7)
	_ = node1.SetOrCreateAttribute("rate", 5.5)
	_ = node1.SetOrCreateAttribute("nickname", "Superhero")
	_ = node1.SetOrCreateAttribute("level", 4.0)
	_ = node1.SetOrCreateAttribute("age", 40)
	err = qConn.InsertEntity(node1)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node1) <<<<<<<")
		return err
	}

	// Node # 11
	node11, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node11) <<<<<<<")
		return err
	}
	_ = node11.SetOrCreateAttribute("name", "Peter-Parker")
	_ = node11.SetOrCreateAttribute("multiple", 7)
	_ = node11.SetOrCreateAttribute("rate", 3.3)
	_ = node11.SetOrCreateAttribute("nickname", "Superhero")
	_ = node11.SetOrCreateAttribute("level", 4.0)
	_ = node11.SetOrCreateAttribute("age", 24)
	err = qConn.InsertEntity(node11)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node11) <<<<<<<")
		return err
	}

	// Node # 12
	node12, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node12) <<<<<<<")
		return err
	}
	_ = node12.SetOrCreateAttribute("name", "Clark-Kent")
	_ = node12.SetOrCreateAttribute("multiple", 7)
	_ = node12.SetOrCreateAttribute("rate", 6.6)
	_ = node12.SetOrCreateAttribute("nickname", "Superhero")
	_ = node12.SetOrCreateAttribute("level", 4.0)
	_ = node12.SetOrCreateAttribute("age", 32)
	err = qConn.InsertEntity(node12)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node12) <<<<<<<")
		return err
	}

	// Node # 13
	node13, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node13) <<<<<<<")
		return err
	}
	_ = node13.SetOrCreateAttribute("name", "James-Logan-Howlett")
	_ = node13.SetOrCreateAttribute("multiple", 7)
	_ = node13.SetOrCreateAttribute("rate", 4.4)
	_ = node13.SetOrCreateAttribute("nickname", "Superhero")
	_ = node13.SetOrCreateAttribute("level", 4.0)
	_ = node13.SetOrCreateAttribute("age", 40)
	err = qConn.InsertEntity(node13)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node13) <<<<<<<")
		return err
	}

	// Node # 14
	node14, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node14) <<<<<<<")
		return err
	}
	_ = node14.SetOrCreateAttribute("name", "Diana-Prince")
	_ = node14.SetOrCreateAttribute("multiple", 7)
	_ = node14.SetOrCreateAttribute("rate", 5.9)
	_ = node14.SetOrCreateAttribute("nickname", "Superheroine")
	_ = node14.SetOrCreateAttribute("level", 4.0)
	_ = node14.SetOrCreateAttribute("age", 35)
	err = qConn.InsertEntity(node14)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node14) <<<<<<<")
		return err
	}

	// Node # 15
	node15, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node15) <<<<<<<")
		return err
	}
	_ = node15.SetOrCreateAttribute("name", "Jean-Grey")
	_ = node15.SetOrCreateAttribute("multiple", 7)
	_ = node15.SetOrCreateAttribute("rate", 6.2)
	_ = node15.SetOrCreateAttribute("nickname", "Superheroine")
	_ = node15.SetOrCreateAttribute("level", 4.0)
	_ = node15.SetOrCreateAttribute("age", 50)
	err = qConn.InsertEntity(node15)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node15) <<<<<<<")
		return err
	}

	// Node # 2
	node2, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node2) <<<<<<<")
		return err
	}
	_ = node2.SetOrCreateAttribute("name", "Mary-Jane-Watson")
	_ = node2.SetOrCreateAttribute("multiple", 14)
	_ = node2.SetOrCreateAttribute("rate", 6.3)
	_ = node2.SetOrCreateAttribute("nickname", "Girlfriend")
	_ = node2.SetOrCreateAttribute("level", 3.0)
	_ = node2.SetOrCreateAttribute("age", 22)
	err = qConn.InsertEntity(node2)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node2) <<<<<<<")
		return err
	}

	// Node # 21
	node21, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node21) <<<<<<<")
		return err
	}
	_ = node21.SetOrCreateAttribute("name", "Lois-Lane")
	_ = node21.SetOrCreateAttribute("multiple", 14)
	_ = node21.SetOrCreateAttribute("rate", 6.4)
	_ = node21.SetOrCreateAttribute("nickname", "Girlfriend")
	_ = node21.SetOrCreateAttribute("level", 3.0)
	_ = node21.SetOrCreateAttribute("age", 30)
	err = qConn.InsertEntity(node21)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node21) <<<<<<<")
		return err
	}

	// Node # 22
	node22, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node22) <<<<<<<")
		return err
	}
	_ = node22.SetOrCreateAttribute("name", "Raven-Darkhölme")
	_ = node22.SetOrCreateAttribute("multiple", 14)
	_ = node22.SetOrCreateAttribute("rate", 7.2)
	_ = node22.SetOrCreateAttribute("nickname", "Criminal")
	_ = node22.SetOrCreateAttribute("level", 3.0)
	_ = node22.SetOrCreateAttribute("age", 36)
	err = qConn.InsertEntity(node22)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node22) <<<<<<<")
		return err
	}

	// Node # 23
	node23, err := createNode(gof, testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node23) <<<<<<<")
		return err
	}
	_ = node23.SetOrCreateAttribute("name", "Selina-Kyle")
	_ = node23.SetOrCreateAttribute("multiple", 14)
	_ = node23.SetOrCreateAttribute("rate", 6.5)
	_ = node23.SetOrCreateAttribute("nickname", "Criminal")
	_ = node23.SetOrCreateAttribute("level", 3.0)
	_ = node23.SetOrCreateAttribute("age", 30)
	err = qConn.InsertEntity(node23)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node23) <<<<<<<")
		return err
	}

	// Node # 24
	node24, err := gof.CreateNode()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node24) <<<<<<<")
		return err
	}
	_ = node24.SetOrCreateAttribute("name", "Harley-Quinn")
	_ = node24.SetOrCreateAttribute("multiple", 14)
	_ = node24.SetOrCreateAttribute("rate", 5.8)
	_ = node24.SetOrCreateAttribute("nickname", "Criminal")
	_ = node24.SetOrCreateAttribute("level", 3.0)
	_ = node24.SetOrCreateAttribute("age", 30)
	err = qConn.InsertEntity(node24)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node24) <<<<<<<")
		return err
	}

	// Node # 3
	node3, err := gof.CreateNode()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node3) <<<<<<<")
		return err
	}
	_ = node3.SetOrCreateAttribute("name", "Lex-Luthor")
	_ = node3.SetOrCreateAttribute("multiple", 52)
	_ = node3.SetOrCreateAttribute("rate", 7.1)
	_ = node3.SetOrCreateAttribute("nickname", "Bad guy")
	_ = node3.SetOrCreateAttribute("level", 11.0)
	_ = node3.SetOrCreateAttribute("age", 26)
	err = qConn.InsertEntity(node3)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node3) <<<<<<<")
		return err
	}

	// Node # 31
	node31, err := gof.CreateNode()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node31) <<<<<<<")
		return err
	}
	_ = node31.SetOrCreateAttribute("name", "Harvey-Dent")
	_ = node31.SetOrCreateAttribute("multiple", 52)
	_ = node31.SetOrCreateAttribute("rate", 4.6)
	_ = node31.SetOrCreateAttribute("nickname", "Bad guy")
	_ = node31.SetOrCreateAttribute("level", 11.0)
	_ = node31.SetOrCreateAttribute("age", 40)
	err = qConn.InsertEntity(node31)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node31) <<<<<<<")
		return err
	}

	// Node # 32
	node32, err := gof.CreateNode()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node32) <<<<<<<")
		return err
	}
	_ = node32.SetOrCreateAttribute("name", "Victor-Creed")
	_ = node32.SetOrCreateAttribute("multiple", 52)
	_ = node32.SetOrCreateAttribute("rate", 6.4)
	_ = node32.SetOrCreateAttribute("nickname", "Bad guy")
	_ = node32.SetOrCreateAttribute("level", 11.0)
	_ = node32.SetOrCreateAttribute("age", 50)
	err = qConn.InsertEntity(node32)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node32) <<<<<<<")
		return err
	}

	// Node # 33
	node33, err := gof.CreateNode()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(node33) <<<<<<<")
		return err
	}
	_ = node33.SetOrCreateAttribute("name", "Norman-Osborn")
	_ = node33.SetOrCreateAttribute("multiple", 52)
	_ = node33.SetOrCreateAttribute("rate", 6.3)
	_ = node33.SetOrCreateAttribute("nickname", "Bad guy")
	_ = node33.SetOrCreateAttribute("level", 11.0)
	_ = node33.SetOrCreateAttribute("age", 50)
	err = qConn.InsertEntity(node33)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(node33) <<<<<<<")
		return err
	}

	// Edge # 1
	edge1, err := gof.CreateEdgeWithDirection(node1, node31, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge1) <<<<<<<")
		return err
	}
	_ = edge1.SetOrCreateAttribute("name", "Nemesis")
	err = qConn.InsertEntity(edge1)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge1) <<<<<<<")
		return err
	}

	// Edge # 2
	edge2, err := gof.CreateEdgeWithDirection(node1, node23, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge2) <<<<<<<")
		return err
	}
	_ = edge2.SetOrCreateAttribute("name", "Frenemy")
	err = qConn.InsertEntity(edge2)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge2) <<<<<<<")
		return err
	}

	// Edge # 3
	edge3, err := gof.CreateEdgeWithDirection(node1, node24, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge3) <<<<<<<")
		return err
	}
	_ = edge3.SetOrCreateAttribute("name", "Nemesis")
	err = qConn.InsertEntity(edge3)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge3) <<<<<<<")
		return err
	}

	// Edge # 4
	edge4, err := gof.CreateEdgeWithDirection(node1, node12, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge4) <<<<<<<")
		return err
	}
	_ = edge4.SetOrCreateAttribute("name", "Teammate")
	err = qConn.InsertEntity(edge4)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge4) <<<<<<<")
		return err
	}

	// Edge # 5
	edge5, err := gof.CreateEdgeWithDirection(node13, node15, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge5) <<<<<<<")
		return err
	}
	_ = edge5.SetOrCreateAttribute("name", "Friend")
	err = qConn.InsertEntity(edge5)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge5) <<<<<<<")
		return err
	}

	// Edge # 6
	edge6, err := gof.CreateEdgeWithDirection(node13, node22, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge6) <<<<<<<")
		return err
	}
	_ = edge6.SetOrCreateAttribute("name", "Enemy")
	err = qConn.InsertEntity(edge6)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge6) <<<<<<<")
		return err
	}

	// Edge # 7
	edge7, err := gof.CreateEdgeWithDirection(node13, node32, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during gof.CreateNode(edge7) <<<<<<<")
		return err
	}
	_ = edge7.SetOrCreateAttribute("name", "Enemy")
	err = qConn.InsertEntity(edge7)
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.InsertEntity(edge7) <<<<<<<")
		return err
	}

	_, err = qConn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from createTestData - error during qConn.Commit() <<<<<<<")
		return err
	}

	fmt.Println(">>>>>>> Returning createTestData: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testPKey(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testPKey: A Unique Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use primary key
	queryStr := fmt.Sprint("@nodetype = 'testnode' and name = 'Lex-Luthor';")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testPKey - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testPKey: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testPKey: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testPartialUnique(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testPartialUnique: A Partial Unique Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use unique key testidx1
	queryStr := fmt.Sprint("@nodetype = 'testnode' and nickname = 'Bad guy' and level = 11.0;")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testPartialUnique - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testPartialUnique: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testPartialUnique: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testNonUnique(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testNonUnique: A Non Unique Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use nonUnique key testidx2
	queryStr := fmt.Sprint("@nodetype = 'testnode' and nickname = 'Superhero';")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testNonUnique - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testNonUnique: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testNonUnique: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testGreaterThan(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testGreaterThan: A Greater Than Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use testidx3 partially
	queryStr := fmt.Sprint("@nodetype = 'testnode' and age > 32;")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testGreaterThan - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testGreaterThan: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testGreaterThan: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testLessThan(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testLessThan: A Less Than Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use testidx3 partially
	queryStr := fmt.Sprint("@nodetype = 'testnode' and age < 28;")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testLessThan - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testLessThan: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testLessThan: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testRange(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testRange: A Range Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use testidx3
	queryStr := fmt.Sprint("@nodetype = 'testnode' and age > 28 and age < 33 and level > 2.9 and level < 4.1;")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testRange - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testRange: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testRange: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testComplex(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testComplex: A Complex Query <<<<<<<")

	//query, err := qConn.CreateQuery("testquery < X '5ef';")
	//rSet := query.Execute()

	option := query.DefaultQueryOption()

	queryStr1 := fmt.Sprint("@nodetype = 'testnode' and nickname = 'Criminal' and level = 3.0 and (rate = 5.8 or rate = 6.5);")
	queryStr2 := fmt.Sprint("@nodetype = 'testnode' and (nickname = 'Criminal' or nickname = 'foo') and level = 3.0 and (rate = 5.8 or rate = 6.5);")
	queryStr3 := fmt.Sprint("@nodetype = '@nodetype = 'testnode' and (nickname = 'Criminal' and level = 3.0 and rate = 5.8) or (nickname = '壞人' and level = 11.0 and rate = 6.3);")
	queryStr4 := fmt.Sprint("@nodetype = 'testnode' and ((nickname = 'Criminal' and level = 3.0 and rate = 5.8) or (nickname = '壞人' and level = 11.0 and rate = 6.3));")
	//queryStr5 := fmt.Sprint("@nodetype = 'testnode' and level < 5.2 and rate < 4.5;")
	//queryStr6 := fmt.Sprint("testquery < X '5ef';")
	rSet, err := qConn.ExecuteQuery(queryStr1, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testComplex - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	rSet, err = qConn.ExecuteQuery(queryStr2, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testComplex - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	rSet, err = qConn.ExecuteQuery(queryStr3, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testComplex - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	rSet, err = qConn.ExecuteQuery(queryStr4, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testComplex - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testComplex: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	//query.Close()

	fmt.Println(">>>>>>> Returning testComplex: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func testNodeTypeOnly(qConn types.TGConnection) types.TGError {
	fmt.Println(">>>>>>> Entering testNodeTypeOnly: A Node Desc Only Query <<<<<<<")

	option := query.DefaultQueryOption()

	// use testidx3
	queryStr := fmt.Sprint("@nodetype = 'testnode';")
	rSet, err := qConn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from testNodeTypeOnly - error during qConn.ExecuteQuery) <<<<<<<")
		return err
	}
	if rSet != nil {
		entityCount := 0
		for {
			if !rSet.HasNext() {
				break
			}
			entityCount++
			fmt.Printf(">>>>>>> Inside testNodeTypeOnly: ResultSet has entity(%d) as: '%+v' <<<<<<<\n", entityCount, rSet.Next())
		}
	}

	fmt.Println(">>>>>>> Returning testNodeTypeOnly: Successful w/ NO ERRORS !!! <<<<<<<")
	return nil
}

func QueryTransactions() {
	fmt.Println(">>>>>>> Entering QueryTransactions <<<<<<<")

	qConn, err := connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during connect() <<<<<<<")
		return
	}

	// Uncomment createTestData(qConn) if-n-when you are executing this test suite in standalone mode.
	// Otherwise use the data that has been created earlier via SimpleTransaction::SimpleInsertTransaction()
	// in file SimpleTransactions.go

	//err = createTestData(qConn)
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning QueryTransactions - error during createTestData() <<<<<<<")
	//	return
	//}

	err = testPKey(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testPKey() <<<<<<<")
		return
	}

	err = testPartialUnique(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testPartialUnique() <<<<<<<")
		return
	}

	err = testNonUnique(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testNonUnique() <<<<<<<")
		return
	}

	err = testGreaterThan(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testGreaterThan() <<<<<<<")
		return
	}

	err = testLessThan(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testLessThan() <<<<<<<")
		return
	}

	err = testRange(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testRange() <<<<<<<")
		return
	}

	err = testNodeTypeOnly(qConn)
	if err != nil {
		fmt.Println(">>>>>>> Returning QueryTransactions - error during testNodeTypeOnly() <<<<<<<")
		return
	}

	//err = testComplex(qConn)
	//if err != nil {
	//	fmt.Println(">>>>>>> Returning QueryTransactions - error during testComplex() <<<<<<<")
	//	return
	//}

	disConnect(qConn)
	fmt.Println(">>>>>>> Returning QueryTransactions w/ NO ERRORS!!! <<<<<<<")
}
