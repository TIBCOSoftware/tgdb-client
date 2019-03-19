package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

func insertTransaction(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering insertTransaction: Insert Few Simple Nodes with individual properties <<<<<<<")

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println("Returning from insertTransaction - error during conn.GetGraphMetadata")
		return
	}

	//testNodeType, err := gmd.GetNodeType("basicnode")
	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println("Returning from insertTransaction - error during gmd.GetNodeType('testnode')")
		return
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'testnode' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'testnode' is not found from meta data fetch <<<<<<<")
		return
	}

	// Node # 1
	node1, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node1) <<<<<<<")
		return
	}

	_ = node1.SetOrCreateAttribute("name", "Bruce###################Wayne")
	_ = node1.SetOrCreateAttribute("multiple", 7)
	_ = node1.SetOrCreateAttribute("rate", 5.5)
	_ = node1.SetOrCreateAttribute("nickname", "超級英雄")
	_ = node1.SetOrCreateAttribute("level", 4.0)
	_ = node1.SetOrCreateAttribute("age", 38)
	//d, _ := strconv.ParseFloat("23423573223423.89813", 64)
	//d, _ := model.NewFromString("23423573223423.89813")
	//d, _ := model.NewFromString("234235732234235723590735124523.89813275891735070")
	//_ = node1.SetOrCreateAttribute("networth", d)
	err = conn.InsertEntity(node1)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node1) <<<<<<<")
		return
	}

	// Node # 11
	node11, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node11) <<<<<<<")
		return
	}
	_ = node11.SetOrCreateAttribute("name", "Peter Parker")
	_ = node11.SetOrCreateAttribute("multiple", 7)
	_ = node11.SetOrCreateAttribute("rate", 3.3)
	_ = node11.SetOrCreateAttribute("nickname", "超級英雄")
	_ = node11.SetOrCreateAttribute("level", 4.0)
	_ = node11.SetOrCreateAttribute("age", 24)
	err = conn.InsertEntity(node11)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node11) <<<<<<<")
		return
	}

	// Node # 12
	node12, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node12) <<<<<<<")
		return
	}
	_ = node12.SetOrCreateAttribute("name", "Clark Kent")
	_ = node12.SetOrCreateAttribute("multiple", 7)
	_ = node12.SetOrCreateAttribute("rate", 6.6)
	_ = node12.SetOrCreateAttribute("nickname", "超級英雄")
	_ = node12.SetOrCreateAttribute("level", 4.0)
	_ = node12.SetOrCreateAttribute("age", 32)
	err = conn.InsertEntity(node12)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node12) <<<<<<<")
		return
	}

	// Node # 13
	node13, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node13) <<<<<<<")
		return
	}
	_ = node13.SetOrCreateAttribute("name", "James Logan Howlett")
	_ = node13.SetOrCreateAttribute("multiple", 7)
	_ = node13.SetOrCreateAttribute("rate", 4.4)
	_ = node13.SetOrCreateAttribute("nickname", "超級英雄")
	_ = node13.SetOrCreateAttribute("level", 4.0)
	_ = node13.SetOrCreateAttribute("age", 40)
	err = conn.InsertEntity(node13)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node13) <<<<<<<")
		return
	}

	// Node # 2
	node2, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node2) <<<<<<<")
		return
	}
	_ = node2.SetOrCreateAttribute("name", "Mary Jane Watson")
	_ = node2.SetOrCreateAttribute("multiple", 14)
	_ = node2.SetOrCreateAttribute("rate", 6.3)
	_ = node2.SetOrCreateAttribute("nickname", "美麗")
	_ = node2.SetOrCreateAttribute("level", 3.0)
	_ = node2.SetOrCreateAttribute("age", 22)
	err = conn.InsertEntity(node2)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node2) <<<<<<<")
		return
	}

	// Node # 21
	node21, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node21) <<<<<<<")
		return
	}
	_ = node21.SetOrCreateAttribute("name", "Lois Lane")
	_ = node21.SetOrCreateAttribute("multiple", 14)
	_ = node21.SetOrCreateAttribute("rate", 6.4)
	_ = node21.SetOrCreateAttribute("nickname", "美麗")
	_ = node21.SetOrCreateAttribute("level", 3.0)
	_ = node21.SetOrCreateAttribute("age", 30)
	err = conn.InsertEntity(node21)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node21) <<<<<<<")
		return
	}

	// Node # 22
	node22, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node22) <<<<<<<")
		return
	}
	_ = node22.SetOrCreateAttribute("name", "Jean Grey")
	_ = node22.SetOrCreateAttribute("multiple", 14)
	_ = node22.SetOrCreateAttribute("rate", 6.2)
	_ = node22.SetOrCreateAttribute("nickname", "美麗")
	_ = node22.SetOrCreateAttribute("level", 3.0)
	_ = node22.SetOrCreateAttribute("age", 30)
	err = conn.InsertEntity(node22)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node22) <<<<<<<")
		return
	}

	// Node # 23
	node23, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node23) <<<<<<<")
		return
	}
	_ = node23.SetOrCreateAttribute("name", "Selina Kyle")
	_ = node23.SetOrCreateAttribute("multiple", 14)
	_ = node23.SetOrCreateAttribute("rate", 6.5)
	_ = node23.SetOrCreateAttribute("nickname", "Criminal")
	_ = node23.SetOrCreateAttribute("level", 3.0)
	_ = node23.SetOrCreateAttribute("age", 30)
	err = conn.InsertEntity(node23)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node23) <<<<<<<")
		return
	}

	// Node # 24
	node24, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node24) <<<<<<<")
		return
	}
	_ = node24.SetOrCreateAttribute("name", "Harley Quinn")
	_ = node24.SetOrCreateAttribute("multiple", 14)
	_ = node24.SetOrCreateAttribute("rate", 5.8)
	_ = node24.SetOrCreateAttribute("nickname", "Criminal")
	_ = node24.SetOrCreateAttribute("level", 3.0)
	_ = node24.SetOrCreateAttribute("age", 30)
	err = conn.InsertEntity(node24)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node24) <<<<<<<")
		return
	}

	// Node # 3
	node3, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node3) <<<<<<<")
		return
	}
	_ = node3.SetOrCreateAttribute("name", "Lex Luthor")
	_ = node3.SetOrCreateAttribute("multiple", 52)
	_ = node3.SetOrCreateAttribute("rate", 7.1)
	_ = node3.SetOrCreateAttribute("nickname", "壞人")
	_ = node3.SetOrCreateAttribute("level", 11.0)
	_ = node3.SetOrCreateAttribute("age", 26)
	err = conn.InsertEntity(node3)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node3) <<<<<<<")
		return
	}

	// Node # 31
	node31, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node31) <<<<<<<")
		return
	}
	_ = node31.SetOrCreateAttribute("name", "Harvey Dent")
	_ = node31.SetOrCreateAttribute("multiple", 52)
	_ = node31.SetOrCreateAttribute("rate", 4.6)
	_ = node31.SetOrCreateAttribute("nickname", "壞人")
	_ = node31.SetOrCreateAttribute("level", 11.0)
	_ = node31.SetOrCreateAttribute("age", 40)
	err = conn.InsertEntity(node31)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node31) <<<<<<<")
		return
	}

	// Node # 32
	node32, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node32) <<<<<<<")
		return
	}
	_ = node32.SetOrCreateAttribute("name", "Victor Creed")
	_ = node32.SetOrCreateAttribute("multiple", 52)
	_ = node32.SetOrCreateAttribute("rate", 6.4)
	_ = node32.SetOrCreateAttribute("nickname", "壞人")
	_ = node32.SetOrCreateAttribute("level", 11.0)
	_ = node32.SetOrCreateAttribute("age", 40)
	err = conn.InsertEntity(node32)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node32) <<<<<<<")
		return
	}

	// Node # 33
	node33, err := gof.CreateNodeInGraph(testNodeType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(node33) <<<<<<<")
		return
	}
	_ = node33.SetOrCreateAttribute("name", "Norman Osborn")
	_ = node33.SetOrCreateAttribute("multiple", 52)
	_ = node33.SetOrCreateAttribute("rate", 6.3)
	_ = node33.SetOrCreateAttribute("nickname", "壞人")
	_ = node33.SetOrCreateAttribute("level", 11.0)
	_ = node33.SetOrCreateAttribute("age", 50)
	err = conn.InsertEntity(node33)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(node33) <<<<<<<")
		return
	}

	// Edge # 1
	edge1, err := gof.CreateEdgeWithDirection(node1, node31, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(edge1) <<<<<<<")
		return
	}
	_ = edge1.SetOrCreateAttribute("name", "Nemesis")
	err = conn.InsertEntity(edge1)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(edge1) <<<<<<<")
		return
	}

	// Edge # 2
	edge2, err := gof.CreateEdgeWithDirection(node1, node23, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(edge12 <<<<<<<")
		return
	}
	_ = edge2.SetOrCreateAttribute("name", "Frenemy")
	err = conn.InsertEntity(edge2)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(edge2) <<<<<<<")
		return
	}

	// Edge # 3
	edge3, err := gof.CreateEdgeWithDirection(node1, node24, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(edge3) <<<<<<<")
		return
	}
	_ = edge3.SetOrCreateAttribute("name", "Nemesis")
	err = conn.InsertEntity(edge3)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(edge1) <<<<<<<")
		return
	}

	// Edge # 4
	edge4, err := gof.CreateEdgeWithDirection(node1, node12, types.DirectionTypeBiDirectional)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during gof.CreateNode(edge4) <<<<<<<")
		return
	}
	_ = edge4.SetOrCreateAttribute("name", "Teammate")
	err = conn.InsertEntity(edge4)
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.InsertEntity(edge4) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from insertTransaction - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning upsertTransaction: Successful w/ NO ERRORS !!! <<<<<<<")
}

func getTransaction(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering getTransaction: Unique Get Operation <<<<<<<")

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println("Returning from getTransaction - error during conn.GetGraphMetadata")
		return
	}

	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println("Returning from getTransaction - error during gmd.GetNodeType('testnode')")
		return
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'testnode' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'testnode' is not found from meta data fetch <<<<<<<")
		return
	}

	depth := 5
	resultCount := 100
	edgeLimit := 0

	// Key
	key, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(key) <<<<<<<")
		return
	}

	option := query.NewQueryOption()
	_ = option.SetPreFetchSize(resultCount)
	_ = option.SetTraversalDepth(depth)
	_ = option.SetEdgeLimit(edgeLimit)

	_ = key.SetOrCreateAttribute("name", "John Doe")
	entity, err := conn.GetEntity(key, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:name-John Doe) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'John Doe' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'John Doe' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key.SetOrCreateAttribute("name", "Bruce Wayne")
	entity, err = conn.GetEntity(key, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:name-Bruce Wayne) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Bruce Wayne' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Bruce Wayne' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key.SetOrCreateAttribute("name", "Peter Parker")
	entity, err = conn.GetEntity(key, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:name-Peter Parker) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Peter Parker' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Peter Parker' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key.SetOrCreateAttribute("name", "Mary Jane Watson")
	entity, err = conn.GetEntity(key, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:name-Mary Jane Watson) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Mary Jane Watson' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Mary Jane Watson' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key.SetOrCreateAttribute("name", "Super Jane")
	entity, err = conn.GetEntity(key, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:name-Super Jane) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Super Jane' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Super Jane' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	// Key1
	key1, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(key1) <<<<<<<")
		return
	}

	_ = key1.SetOrCreateAttribute("nickname", "Stupid")
	entity, err = conn.GetEntity(key1, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:nickname-Stupid) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Stupid' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Stupid' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key1.SetOrCreateAttribute("name", "Jean Grey")
	_ = key1.SetOrCreateAttribute("nickname", "美麗")
	entity, err = conn.GetEntity(key1, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:nickname-美麗) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for '美麗' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for '美麗' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key1.SetOrCreateAttribute("name", "Harley Quinn")
	_ = key1.SetOrCreateAttribute("nickname", "Criminal")
	entity, err = conn.GetEntity(key1, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:nickname-Criminal) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Criminal' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Criminal' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_ = key1.SetOrCreateAttribute("name", "James Logan Howlett")
	_ = key1.SetOrCreateAttribute("nickname", "超級英雄")
	entity, err = conn.GetEntity(key1, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:nickname-超級英雄) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for '超級英雄' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for '超級英雄' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	// Key2
	key2, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(key2) <<<<<<<")
		return
	}

	_ = key2.SetOrCreateAttribute("name", "Lex Luthor")
	_ = key2.SetOrCreateAttribute("nickname", "壞人")
	_ = key2.SetOrCreateAttribute("level", 11.0)
	_ = key2.SetOrCreateAttribute("rate", 7.1)
	entity, err = conn.GetEntity(key2, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:level|rate=11.0|7.1) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for '壞人'(11.0)(7.1)' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for '壞人'(11.0)(7.1)' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	// Key3
	key3, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(key3) <<<<<<<")
		return
	}

	_ = key3.SetOrCreateAttribute("multiple", 52)
	entity, err = conn.GetEntity(key3, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.GetEntity(key-attr:multiple-52) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'multiple(52)' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'multiple(52)' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning getTransaction: Successful w/ NO ERRORS !!! <<<<<<<")
}

func updateTransaction(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering updateTransaction: A Sample Update <<<<<<<")

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println("Returning from updateTransaction - error during conn.GetGraphMetadata")
		return
	}

	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println("Returning from updateTransaction - error during gmd.GetNodeType('testnode')")
		return
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'testnode' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'testnode' is not found from meta data fetch <<<<<<<")
		return
	}

	// Node # 13
	keyNode13, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(keyNode13) <<<<<<<")
		return
	}
	_ = keyNode13.SetOrCreateAttribute("name", "James Logan Howlett")
	_ = keyNode13.SetOrCreateAttribute("multiple", 7)
	_ = keyNode13.SetOrCreateAttribute("rate", 4.4)
	_ = keyNode13.SetOrCreateAttribute("nickname", "超級英雄")
	_ = keyNode13.SetOrCreateAttribute("level", 4.0)
	_ = keyNode13.SetOrCreateAttribute("age", 40)
	node13, err := conn.GetEntity(keyNode13, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.GetEntity(keyNode13) <<<<<<<")
		return
	}
	if node13 == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'James Logan Howlett(keyNode13)' returned NOTHING <<<<<<<")
		return
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'James Logan Howlett(keyNode13)' returned '%+v' <<<<<<<\n", node13)
	}
	_ = node13.SetOrCreateAttribute("nickname", "超級")
	_ = node13.SetOrCreateAttribute("age", 43)
	err = conn.UpdateEntity(node13)
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.UpdateEntity(node13) <<<<<<<")
		return
	}

	// Node # 32
	keyNode32, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(keyNode32) <<<<<<<")
		return
	}
	_ = keyNode32.SetOrCreateAttribute("name", "Victor Creed")
	_ = keyNode32.SetOrCreateAttribute("multiple", 52)
	_ = keyNode32.SetOrCreateAttribute("rate", 6.4)
	_ = keyNode32.SetOrCreateAttribute("nickname", "壞人")
	_ = keyNode32.SetOrCreateAttribute("level", 11.0)
	_ = keyNode32.SetOrCreateAttribute("age", 40)
	node32, err := conn.GetEntity(keyNode13, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.GetEntity(keyNode32) <<<<<<<")
		return
	}
	if node32 == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Victor Creed(keyNode32)' returned NOTHING <<<<<<<")
		return
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Victor Creed(keyNode13)' returned '%+v' <<<<<<<\n", node13)
	}
	_ = node32.SetOrCreateAttribute("nickname", "害怕")
	_ = node32.SetOrCreateAttribute("level", 11.0)
	_ = node32.SetOrCreateAttribute("age", 38)
	err = conn.UpdateEntity(node32)
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.UpdateEntity(node32) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.Commit() <<<<<<<")
		return
	}

	// Key2
	key2, err := gof.CreateCompositeKey("testnode")
	if err != nil {
		fmt.Println(">>>>>>> Returning from getTransaction - error during gof.CreateCompositeKey(key2) <<<<<<<")
		return
	}

	_ = key2.SetOrCreateAttribute("name", "Victor Creed")
	_ = key2.SetOrCreateAttribute("nickname", "害怕")
	_ = key2.SetOrCreateAttribute("level", 11.0)
	_ = key2.SetOrCreateAttribute("rate", 6.4)
	entity, err := conn.GetEntity(key2, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.GetEntity(key-attr:nil) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'Sabertooth(11.0)(6.4)' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'Sabertooth(11.0)(6.4)' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from updateTransaction - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning updateTransaction: Successful w/ NO ERRORS !!! <<<<<<<")
}

func queryTransaction(conn types.TGConnection) {
	fmt.Println(">>>>>>> Entering queryTransaction: A Sample Query <<<<<<<")

	option := query.DefaultQueryOption()
	_ = option.SetPreFetchSize(-1)
	_ = option.SetTraversalDepth(-1)
	_ = option.SetEdgeLimit(-1)

	queryStr := "@nodetype = 'testnode' and ((nickname = '壞人' and level = 11.0) or (level = 3.0));"	// Expect 9 Entities back
	//queryStr := "@nodetype = 'testnode' and (name = 'Bruce Wayne');"	// Expect 1 Entity back
	rSet, err := conn.ExecuteQuery(queryStr, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from QueryTransaction - error during conn.ExecuteQuery) <<<<<<<")
		return
	}
	fmt.Printf(">>>>>>> Inside queryTransaction: ResultSet: '%+v' <<<<<<<\n", rSet)

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from queryTransaction - error during conn.Commit() <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Returning queryTransaction: Successful w/ NO ERRORS !!! <<<<<<<")
}

func SimpleGetTransaction() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleGetTransaction - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleGetTransaction - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleGetTransaction - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from SimpleGetTransaction - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside SimpleGetTransaction: Get Transaction: Read Inserted Nodes <<<<<<<")
	getTransaction(conn, gof)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleGetTransaction - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from SimpleGetTransaction - successfully disconnected. <<<<<<<")
}

func SimpleUpdateTransaction() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside SimpleUpdateTransaction: Update Transaction: Update few Nodes <<<<<<<")
	updateTransaction(conn, gof)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from SimpleUpdateTransaction - successfully disconnected. <<<<<<<")
}

func SimpleQueryTransaction() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside SimpleQueryTransaction: Query Transaction: Query Inserted Nodes <<<<<<<")
	queryTransaction(conn)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from SimpleQueryTransaction - successfully disconnected. <<<<<<<")
}

func SimpleInsertTransaction() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - error during CreateConnection <<<<<<<")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - error during conn.Connect <<<<<<<")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside SimpleInsertTransaction: Insert Few Simple Nodes with individual properties <<<<<<<")
	insertTransaction(conn, gof)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - error during conn.Disconnect <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Returning from SimpleInsertTransaction - successfully disconnected. <<<<<<<")
}

func SimpleTransaction() {
	fmt.Println(">>>>>>> Inside SimpleTransaction: About to Test SimpleInsertTransaction() <<<<<<<")
	SimpleInsertTransaction()
	fmt.Println(">>>>>>> Inside SimpleTransaction: About to Test SimpleGetTransaction() <<<<<<<")
	SimpleGetTransaction()
	fmt.Println(">>>>>>> Inside SimpleTransaction: About to Test SimpleUpdateTransaction() <<<<<<<")
	SimpleUpdateTransaction()
	fmt.Println(">>>>>>> Inside SimpleTransaction: About to Test SimpleQueryTransaction() <<<<<<<")
	SimpleQueryTransaction()
	fmt.Println(">>>>>>> Returning from SimpleTransaction. <<<<<<<")
}
