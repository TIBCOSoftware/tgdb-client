package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"time"
)

const (
	metaUrl = "tcp://scott@localhost:8222"
	//url1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//url2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//url3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//url4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//url5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	metaPassword = "scott"
	//metaPassword         = "admin"
	//prefetchMetaData = false
)

/**
 * Need to execute metadescript.tql first
 */
func Test1(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering Test1 <<<<<<<")
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	nodeOneAttrType, err := gmd.GetNodeType("nodeOneAttr")
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during conn.GetNodeType('nodeOneAttr') <<<<<<<")
		return
	}
	if nodeOneAttrType != nil {
		fmt.Printf(">>>>>>> 'nodeOneAttrType' is found with %d attributes <<<<<<<\n", len(nodeOneAttrType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'nodeOneAttrType' is not found from meta data fetch <<<<<<<")
		return
	}

	// Node
	node, err := gof.CreateNodeInGraph(nodeOneAttrType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during gof.CreateNode(node) <<<<<<<")
		return
	}
	_ = node.SetOrCreateAttribute("stringAttr", "StringKey")
	err = conn.InsertEntity(node)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during conn.InsertEntity(node) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during conn.Commit() <<<<<<<")
		return
	}

	// Key
	key, err := gof.CreateCompositeKey("nodeOneAttr")
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during gof.CreateCompositeKey(nodeOneAttr) <<<<<<<")
		return
	}
	_ = key.SetOrCreateAttribute("stringAttr", "StringKey")
	option := query.NewQueryOption()

	entity, err := conn.GetEntity(key, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test1 - error during conn.GetEntity(key-attr:stringAttr-StringKey) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> conn.GetEntity() for 'StringKey' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> conn.GetEntity() for 'StringKey' returned '%+v' <<<<<<<\n", entity)
		//printEntities(ent, printDepth, 0, "", new HashMap<Integer, TGEntity>());
	}

	fmt.Println(">>>>>>> Returning Test1 w/ NO ERROR!!!<<<<<<<")
}

func Test2(conn types.TGConnection, gof types.TGGraphObjectFactory) {
	fmt.Println(">>>>>>> Entering Test2 <<<<<<<")
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during conn.GetGraphMetadata <<<<<<<")
		return
	}

	nodeAllAttrsType, err := gmd.GetNodeType("nodeAllAttrs")
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during conn.GetNodeType('nodeTenAttrs') <<<<<<<")
		return
	}
	if nodeAllAttrsType != nil {
		fmt.Printf(">>>>>>> 'nodeAllAttrs' is found with %d attributes <<<<<<<\n", len(nodeAllAttrsType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'nodeAllAttrs' is not found from meta data fetch <<<<<<<")
		return
	}

	// Node
	node, err := gof.CreateNodeInGraph(nodeAllAttrsType)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during gof.CreateNode(node) <<<<<<<")
		return
	}
	_ = node.SetOrCreateAttribute("boolAttr", false)
	_ = node.SetOrCreateAttribute("byteAttr", byte(0xba))
	_ = node.SetOrCreateAttribute("charAttr", '*')
	_ = node.SetOrCreateAttribute("shortAttr", int16(6385))
	_ = node.SetOrCreateAttribute("intAttr", int(73741825))
	_ = node.SetOrCreateAttribute("longAttr", int64(1342177281))
	_ = node.SetOrCreateAttribute("floatAttr", float32(2.23))
	_ = node.SetOrCreateAttribute("doubleAttr", float64(2336.32424))
	n, _ := utils.NewTGDecimalFromString("234235732234235723590735124523.89813275891735070")
	_ = node.SetOrCreateAttribute("numberAttr", n)
	_ = node.SetOrCreateAttribute("stringAttr", "betterStringKey")

	env := utils.NewTGEnvironment()
	tod := time.Date(2016, time.October, 31, 0, 0, 0, 0, time.UTC).Format(env.GetDefaultDateTimeFormat())
	_ = node.SetOrCreateAttribute("dateAttr", tod)

	tom := time.Date(1970, 01, 01, 21, 32, 12, 845000, time.UTC).Format(env.GetDefaultDateTimeFormat())
	_ = node.SetOrCreateAttribute("timeAttr", tom)

	tsom := time.Date(2016, time.October, 25, 8, 9, 30, 999000, time.UTC).Format(env.GetDefaultDateTimeFormat())
	_ = node.SetOrCreateAttribute("timestampAttr", tsom)

	err = conn.InsertEntity(node)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during conn.InsertEntity(node) <<<<<<<")
		return
	}

	_, err = conn.Commit()
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during conn.Commit() <<<<<<<")
		return
	}
	fmt.Println(">>>>>>> Inside Test2 - Node w/ all 13 types of attributes successfully inserted <<<<<<<")

	// Key
	key, err := gof.CreateCompositeKey("nodeAllAttrs")
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during gof.CreateCompositeKey(nodeAllAttrs) <<<<<<<")
		return
	}
	_ = key.SetOrCreateAttribute("stringAttr", "betterStringKey")
	option := query.NewQueryOption()

	entity, err := conn.GetEntity(key, option)
	if err != nil {
		fmt.Println(">>>>>>> Returning from Test2 - error during conn.GetEntity(key-attr:stringAttr-betterStringKey) <<<<<<<")
		return
	}
	if entity == nil {
		fmt.Println(">>>>>>> Inside Test2 - conn.GetEntity() for 'betterStringKey' returned NOTHING <<<<<<<")
	} else {
		fmt.Printf(">>>>>>> Inside Test2 - conn.GetEntity() for 'betterStringKey' returned '%+v' <<<<<<<\n", entity)
	}

	fmt.Println(">>>>>>> Returning Test2 w/ NO ERROR!!!<<<<<<<")
}

func MetadataTest() {
	fmt.Println("Entering MetadataTest")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(metaUrl, "", metaPassword, nil)
	if err != nil {
		fmt.Println("Returning from MetadataTest - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from MetadataTest - error during conn.Connect")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println(">>>>>>> Returning from MetadataTest - error during conn.GetGraphObjectFactory <<<<<<<")
		return
	}
	if gof == nil {
		fmt.Println(">>>>>>> Returning from MetadataTest - Graph Object Factory is null <<<<<<<")
		return
	}

	fmt.Println(">>>>>>> Inside MetadataTest: About to Test1 <<<<<<<")
	Test1(conn, gof)
	fmt.Println(">>>>>>> Inside MetadataTest: About to Test2 <<<<<<<")
	Test2(conn, gof)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from MetadataTest - error during conn.Disconnect")
		return
	}
	fmt.Println("Returning from MetadataTest - successfully disconnected.")
}
