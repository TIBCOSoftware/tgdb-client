package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"time"
)

const (
	multiTxnUrl = "tcp://scott@localhost:8222"
	//multiTxnUrl1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//multiTxnUrl2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//multiTxnUrl3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//multiTxnUrl4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//multiTxnUrl5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	multiTxnPwd         = "scott"
	//prefetchMetaData = false
)

//var basicNodeType, rateNodeType, testNodeType types.TGNodeType
//var john, smith, kelly types.TGNode
//var brother, wife types.TGEdge

func MultiTxn0(conn types.TGConnection, gof types.TGGraphObjectFactory) types.TGNode {
	fmt.Println(">>>>>>> Entering MultiTxn0 - Insert Simple Node(John) of basicnode with a few properties.<<<<<<<")
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.GetGraphMetadata <<<<<<<")
		return nil
	}

	basicNodeType, err := gmd.GetNodeType("basicnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.GetNodeType('basicnode') <<<<<<<")
		return nil
	}
	if basicNodeType != nil {
		fmt.Printf(">>>>>>> 'basicNodeType' is found with %d attributes <<<<<<<\n", len(basicNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'basicNodeType' is not found from meta data fetch <<<<<<<")
		return nil
	}

	// Node (j1)
	j1, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during gof.CreateNode(j1) <<<<<<<")
		return nil
	}
	_ = j1.SetOrCreateAttribute("name", "john")
	_ = j1.SetOrCreateAttribute("age", 30)
	//_ = j1.SetOrCreateAttribute("nickname", "美麗")
	timeStampFormat := time.RFC3339
	tsom, err1 := time.Parse(timeStampFormat, "2016-10-25T15:09:30Z")
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during time.Parse(timeStampFormat, 2016-10-31T21:32:12) <<<<<<<")
		return nil
	}
	_ = j1.SetOrCreateAttribute("createtm", tsom)
	//_ = j1.SetOrCreateAttribute("networth", 2378989.567)
	_ = j1.SetOrCreateAttribute("flag", 'D')
	_ = j1.SetOrCreateAttribute("desc", "Hi TIBCO Team!\n\nThe second stop on the TIBCO NOW Global Tour is just days away. We saw extreme value from Singapore and the excitement and now we do it again in Berlin this time with 545 registered attendees! We have reached and exceeded our target and will be closing registration before we run into any capacity issues. We are very excited about this event and to see what is coming with some game changing product updates shown for the first time at TIBCO NOW Berlin. (There will be a Sharpen the Saw on this Friday)\n\n")
	err = conn.InsertEntity(j1)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.InsertEntity(j1) <<<<<<<")
		return nil
	}

	// Node (j2)
	j2, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during gof.CreateNode(j2) <<<<<<<")
		return nil
	}
	_ = j2.SetOrCreateAttribute("name", "jane")
	_ = j2.SetOrCreateAttribute("age", 30)
	//_ = j1.SetOrCreateAttribute("nickname", "美麗")
	tsom, err1 = time.Parse(timeStampFormat, "2016-10-25T15:09:30Z")
	if err1 != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during time.Parse(timeStampFormat, 2016-10-31T21:32:12) <<<<<<<")
		return nil
	}
	_ = j2.SetOrCreateAttribute("createtm", tsom)
	//_ = j2.SetOrCreateAttribute("networth", 2378989.567)
	_ = j2.SetOrCreateAttribute("flag", 'D')
	_ = j2.SetOrCreateAttribute("desc", "Hi TIBCO Team!\n\nThe second stop on the TIBCO NOW Global Tour is just days away. We saw extreme value from Singapore and the excitement and now we do it again in Berlin this time with 545 registered attendees! We have reached and exceeded our target and will be closing registration before we run into any capacity issues. We are very excited about this event and to see what is coming with some game changing product updates shown for the first time at TIBCO NOW Berlin. (There will be a Sharpen the Saw on this Friday)\n\n")
	err = conn.InsertEntity(j2)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.InsertEntity(j1) <<<<<<<")
		return nil
	}

	for i:=0; i<1000; i++ {
		// Edge (Brother)
		edge, err := gof.CreateEdgeWithDirection(j1, j2, types.DirectionTypeDirected)
		if err != nil {
			fmt.Println(">>>>>>> Returning from MultiTxn0 - error during gof.CreateEdgeWithDirection(j1, j2) <<<<<<<")
			return nil
		}
		_ = edge.SetOrCreateAttribute("name", "spouse")
		_ = edge.SetOrCreateAttribute("desc", "This is test...")
		err = conn.InsertEntity(edge)
		if err != nil {
			fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.InsertEntity(brother) <<<<<<<")
			return nil
		}
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn0 - error during conn.Commit() <<<<<<<")
		return nil
	}

	//john = j1
	fmt.Println(">>>>>>> Returning MultiTxn0 w/ NO ERROR!!!<<<<<<<")
	return j1
}

func MultiTxn1(conn types.TGConnection, gof types.TGGraphObjectFactory, basicNodeType types.TGNodeType) types.TGNode {
	fmt.Println(">>>>>>> Entering MultiTxn1 <<<<<<<")

	// Node
	node, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1 - error during gof.CreateNode(node) <<<<<<<")
		return nil
	}
	_ = node.SetOrCreateAttribute("name", "john")
	_ = node.SetOrCreateAttribute("age", 30)
	_ = node.SetOrCreateAttribute("desc", "美麗")
	env := utils.NewTGEnvironment()
	createTm := time.Date(2016, time.October, 01, 15, 05, 30, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
	_ = node.SetOrCreateAttribute("createtm", createTm)
	//_ = node.SetOrCreateAttribute("networth", 2378989.567)
	_ = node.SetOrCreateAttribute("flag", 'D')
	err = conn.InsertEntity(node)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1 - error during conn.InsertEntity(node) <<<<<<<")
		return nil
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1 - error during conn.Commit() <<<<<<<")
		return nil
	}

	//john = node
	fmt.Println(">>>>>>> Returning MultiTxn1 w/ NO ERROR!!!<<<<<<<")
	return node
}

func MultiTxn1_1(conn types.TGConnection, gof types.TGGraphObjectFactory) types.TGNode {
	fmt.Println(">>>>>>> Entering MultiTxn1_1: Get the Entity that we inserted.<<<<<<<")

	// Key
	key, err := gof.CreateCompositeKey("basicnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1_1 - error during gof.CreateCompositeKey(basicnode) <<<<<<<")
		return nil
	}
	_ = key.SetOrCreateAttribute("name", "john2")
	entity, err := conn.GetEntity(key, query.DefaultQueryOption())

	node := entity.(*model.Node)

	fmt.Println(">>>>>>> Returning MultiTxn1_1 w/ NO ERROR!!!<<<<<<<")
	return node
}

func MultiTxn1_2(conn types.TGConnection, gof types.TGGraphObjectFactory, basicNodeType types.TGNodeType) {
	fmt.Println(">>>>>>> Entering MultiTxn1_2: Again insert John. This should raise Unique Key Constraint violation.<<<<<<<")

	// Node
	node, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1_2 - error during gof.CreateNode(node) <<<<<<<")
		return
	}
	_ = node.SetOrCreateAttribute("name", "john")
	_ = node.SetOrCreateAttribute("age", 30)
	_ = node.SetOrCreateAttribute("desc", "美麗")
	err = conn.InsertEntity(node)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn1_2 - error during conn.InsertEntity() <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Printf(">>>>>>> Returning from MultiTxn1_2 - Expected exception: %s\n <<<<<<<", err.Error())
		return
	}

	fmt.Println(">>>>>>> Returning MultiTxn1_2 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn2(conn types.TGConnection, gof types.TGGraphObjectFactory, john types.TGNode) {
	fmt.Println(">>>>>>> Entering MultiTxn2: Update Node John's attribute. <<<<<<<")

	// Node (John)
	_ = john.SetOrCreateAttribute("age", 35)
	err := conn.UpdateEntity(john)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn2 - error during conn.UpdateEntity(node) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn2 - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning MultiTxn2 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn3(conn types.TGConnection, gof types.TGGraphObjectFactory, basicNodeType types.TGNodeType) (types.TGNode, types.TGEdge) {
	fmt.Println(">>>>>>> Entering MultiTxn3 - Insert 2 nodes, and set a relation between them.<<<<<<<")

	// Node (Smith)
	smith, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during gof.CreateNode(node) <<<<<<<")
		return nil, nil
	}
	_ = smith.SetOrCreateAttribute("name", "smith")
	_ = smith.SetOrCreateAttribute("age", 30)
	_ = smith.SetOrCreateAttribute("desc", "will")
	err = conn.InsertEntity(smith)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during conn.InsertEntity(smith) <<<<<<<")
		return nil, nil
	}

	// Node (Kelly)
	kelly, err := gof.CreateNodeInGraph(basicNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during gof.CreateNode(node) <<<<<<<")
		return nil, nil
	}
	_ = kelly.SetOrCreateAttribute("name", "kelly")
	_ = kelly.SetOrCreateAttribute("age", 28)
	_ = kelly.SetOrCreateAttribute("desc", "Ki")
	err = conn.InsertEntity(kelly)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during conn.InsertEntity(kelly) <<<<<<<")
		return nil, nil
	}

	// Edge (Brother)
	brother, err := gof.CreateEdgeWithDirection(smith, kelly, types.DirectionTypeDirected)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during gof.CreateEdgeWithDirection(smith, kelly) <<<<<<<")
		return nil, nil
	}
	_ = brother.SetOrCreateAttribute("name", "Sister")
	err = conn.InsertEntity(brother)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during conn.InsertEntity(brother) <<<<<<<")
		return nil, nil
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3 - error during conn.Commit() <<<<<<<")
		return nil, nil
	}

	fmt.Println(">>>>>>> Returning MultiTxn3 w/ NO ERROR!!!<<<<<<<")
	return kelly, brother
}

func MultiTxn3_1(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering MultiTxn3_1: Get the Entity that we inserted. <<<<<<<")

	// Key
	key, err := gof.CreateCompositeKey("basicnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3_1 - error during gof.CreateCompositeKey(basicnode) <<<<<<<")
		return
	}
	_ = key.SetOrCreateAttribute("name", "smith")
	entity, err := conn.GetEntity(key, query.DefaultQueryOption())
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn3_1 - error during conn.GetEntity(node) <<<<<<<")
		return
	}

	fmt.Printf(">>>>>>> Returning MultiTxn3_1 w/ entity: '%+v'<<<<<<<\n", entity)
}

func MultiTxn4(conn types.TGConnection, gof types.TGGraphObjectFactory, john, kelly types.TGNode) types.TGEdge {
	fmt.Println(">>>>>>> Entering MultiTxn4: Add an edge between 2 existing nodes - In case between john and kelly <<<<<<<")

	// Edge (Wife)
	wife, err := gof.CreateEdgeWithDirection(kelly, john, types.DirectionTypeDirected)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn4 - error during gof.CreateNode(wife) <<<<<<<")
		return nil
	}
	_ = wife.SetOrCreateAttribute("name", "wife")
	err = conn.InsertEntity(wife)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn4 - error during conn.InsertEntity(node) <<<<<<<")
		return nil
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn4 - error during conn.Commit() <<<<<<<")
		return nil
	}

	fmt.Println(">>>>>>> Returning MultiTxn4 w/ NO ERROR!!!<<<<<<<")
	return wife
}

func MultiTxn4_0(conn types.TGConnection, gof types.TGGraphObjectFactory) types.TGNode {
	fmt.Println(">>>>>>> Entering MultiTxn4_0: Get the Entity that we inserted <<<<<<<")

	// Key
	key, err := gof.CreateCompositeKey("basicnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn4_0 - error during gof.CreateCompositeKey(basicnode) <<<<<<<")
		return nil
	}
	_ = key.SetOrCreateAttribute("name", "kelly")

	entity, err := conn.GetEntity(key, query.DefaultQueryOption())
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn4_0 - error during conn.GetEntity(key) <<<<<<<")
		return nil
	}

	node := entity.(*model.Node)

	fmt.Printf(">>>>>>> Returning MultiTxn4_0 w/ entity: '%+v'<<<<<<<\n", entity)
	return node
}

func MultiTxn5(conn types.TGConnection, gof types.TGGraphObjectFactory, wife types.TGEdge) {
	fmt.Println(">>>>>>> Entering MultiTxn5: Update an existing Edge <<<<<<<")

	// Node (wife)
	if wife != nil {
		_ = wife.SetOrCreateAttribute("name", "wife")
		//_ = wife.SetOrCreateAttribute("dom", "10/2/2016")
		//env := utils.NewTGEnvironment()
		//dom := time.Date(2016, time.October, 10, 00, 00, 00, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
		//_ = wife.SetOrCreateAttribute("dom", dom)
		_ = wife.SetOrCreateAttribute("desc", "This is update...")
		err := conn.UpdateEntity(wife)
		if err != nil {
			fmt.Println(">>>>>>> Returning from MultiTxn5 - error during conn.UpdateEntity(wife) <<<<<<<")
			return
		}

		_, err = conn.Commit()
		if err != nil {
			fmt.Println(">>>>>>> Returning from MultiTxn5 - error during conn.Commit() <<<<<<<")
			return
		}
	}

	fmt.Println(">>>>>>> Returning MultiTxn5 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn6(conn types.TGConnection, gof types.TGGraphObjectFactory, john types.TGNode) {
	fmt.Println(">>>>>>> Entering MultiTxn6: Deleting Node john <<<<<<<")

	err := conn.DeleteEntity(john)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn6 - error during conn.DeleteEntity(john) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn6 - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning MultiTxn6 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn6_1(conn types.TGConnection, gof types.TGGraphObjectFactory, john types.TGNode) {
	fmt.Println(">>>>>>> Entering MultiTxn6_1: Updating Node john Again - Should throw mismatch of ERA. or deleted <<<<<<<")

	// Node
	_ = john.SetOrCreateAttribute("age", 40)
	_ = john.SetOrCreateAttribute("desc", "美麗")
	err := conn.UpdateEntity(john)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn6_1 - error during conn.UpdateEntity(john) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn6_1 - expected error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning MultiTxn6_1 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn7(conn types.TGConnection, gof types.TGGraphObjectFactory, brother types.TGEdge) {
	fmt.Println(">>>>>>> Entering MultiTxn7: Deleting Egde <<<<<<<")

	// Edge (Brother)
	err := conn.DeleteEntity(brother)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn7 - error during conn.DeleteEntity(brother) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn7 - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning MultiTxn7 w/ NO ERROR!!!<<<<<<<")
}

func MultiTxn8(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering MultiTxn8: Query w/ Filter from smith to kelly <<<<<<<")

	startFilter := "@nodetype = 'basicnode' and name = 'smith';"
	traverserFilter := "@degree = 1;"
	endFilter := "@nodetype = 'basicnode' and name = 'kelly';"

	result, err := conn.ExecuteQueryWithFilter(startFilter, "", traverserFilter, endFilter, query.DefaultQueryOption())
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn8 - error during conn.ExecuteQueryWithFilter(startFilter, traverserFilter, endFilter) <<<<<<<")
		return
	}

	fmt.Printf(">>>>>>> Returning MultiTxn8 w/ results: '%+v'<<<<<<<\n", result)
}

func MultiTransactionTest() {
	fmt.Println("Entering MultiTransactionTest")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(multiTxnUrl, "", multiTxnPwd, nil)
	if err != nil {
		fmt.Println("Returning from MultiTransactionTest - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from MultiTransactionTest - error during conn.Connect")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTransactionTest - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from MultiTransactionTest - Graph Object Factory is null <<<<<<<")
		return
	}

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn8 - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	basicNodeType, err := gmd.GetNodeType("basicnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn8 - error during conn.GetNodeType('basicnode') <<<<<<<")
		return
	}
	if basicNodeType != nil {
		fmt.Printf(">>>>>>> 'basicNodeType' is found with %d attributes <<<<<<<\n", len(basicNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'basicNodeType' is not found from meta data fetch <<<<<<<")
		return
	}

	rateNodeType, err := gmd.GetNodeType("ratenode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn8 - error during conn.GetNodeType('ratenode') <<<<<<<")
		return
	}
	if basicNodeType != nil {
		fmt.Printf(">>>>>>> 'rateNodeType' is found with %d attributes <<<<<<<\n", len(rateNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'rateNodeType' is not found from meta data fetch <<<<<<<")
		return
	}

	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from MultiTxn8 - error during conn.GetNodeType('testnode') <<<<<<<")
		return
	}
	if basicNodeType != nil {
		fmt.Printf(">>>>>>> 'testNodeType' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'testNodeType' is not found from meta data fetch <<<<<<<")
		return
	}

	//fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn0 <<<<<<<")
	//MultiTxn0(conn, gof)
	//fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn1_1 <<<<<<<")
	//MultiTxn1_1(conn, gof)

	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn1 <<<<<<<")
	john := MultiTxn1(conn, gof, basicNodeType)
	//fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn1_2 <<<<<<<")
	//MultiTxn1_2(conn, gof)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn2 <<<<<<<")
	MultiTxn2(conn, gof, john)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn3 <<<<<<<")
	//kelly, brother := MultiTxn3(conn, gof, basicNodeType)
	kelly, _ := MultiTxn3(conn, gof, basicNodeType)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn3_1 <<<<<<<")
	MultiTxn3_1(conn, gof)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn4 <<<<<<<")
	wife := MultiTxn4(conn, gof, john, kelly)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn4_0 <<<<<<<")
	kelly = MultiTxn4_0(conn, gof)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn8 <<<<<<<")
	MultiTxn8(conn, gof)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn5 <<<<<<<")
	MultiTxn5(conn, gof, wife)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn6 <<<<<<<")
	MultiTxn6(conn, gof, john)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn1 <<<<<<<")
	john = MultiTxn1(conn, gof, basicNodeType)
	fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn6_1 <<<<<<<")
	MultiTxn6_1(conn, gof, john)
	//fmt.Println(">>>>>>> Inside MultiTransactionTest: About to MultiTxn7 <<<<<<<")
	//MultiTxn7(conn, gof, brother)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from MultiTransactionTest - error during conn.Disconnect")
		return
	}
	fmt.Println("Returning from MultiTransactionTest - successfully disconnected.")
}
