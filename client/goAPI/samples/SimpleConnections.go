package samples

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	url = "tcp://scott@localhost:8222"
	//url0 = "tcp://scott@localhost:8222"
	//url1 = "tcp://scott@[fe80::1c15:49f2:b621:7ced%en0:8222]";
	//url2 = "tcp://scott@localhost:8222/{connectTimeout=30}";
	//url3 = "tcp://scott@localhost:8222/{dbName=mod;verifyDBName=true}";
	//url4 = "ssl://scott@localhost:8223/{dbName=mod;verifyDBName=true}";
	//url5 = "ssl://scott@localhost:8223/{ftHosts=192.168.1.15:8222;ftRetryCount=5;ftRetryIntervalSeconds=30;dbName=mod;verifyDBName=true}";
	password         = "scott"
	//prefetchMetaData = false
)

func SimpleConnectAndValidateBootstrappedEntities() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.Connect")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetGraphObjectFactory")
		return
	}
	if gof == nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - Graph Object Factory is null")
		return
	}

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetGraphMetadata")
		return
	}

	attrDesc, err := gmd.GetAttributeDescriptor("factor")
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetAttributeDescriptor('factor')")
		return
	}
	if attrDesc != nil {
		fmt.Printf(">>>>>>> 'factor' has id %d and data desc: %s <<<<<<<\n", attrDesc.GetAttributeId(), types.GetAttributeTypeFromId(attrDesc.GetAttrType()).GetTypeName())
	} else {
		fmt.Println(">>>>>>> 'factor' is not found from meta data fetch <<<<<<<")
		return
	}

	attrDesc, err = gmd.GetAttributeDescriptor("level")
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetAttributeDescriptor('level')")
		return
	}
	if attrDesc != nil {
		fmt.Printf(">>>>>>> 'level' has id %d and data desc: %s <<<<<<<\n", attrDesc.GetAttributeId(), types.GetAttributeTypeFromId(attrDesc.GetAttrType()).GetTypeName())
	} else {
		fmt.Println(">>>>>>> 'level' is not found from meta data fetch <<<<<<<")
		return
	}

	testNodeType, err := gmd.GetNodeType("testnode")
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetNodeType('testnode')")
		return
	}
	if testNodeType != nil {
		fmt.Printf(">>>>>>> 'testnode' is found with %d attributes <<<<<<<\n", len(testNodeType.GetAttributeDescriptors()))
	} else {
		fmt.Println(">>>>>>> 'testnode' is not found from meta data fetch <<<<<<<")
		return
	}

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.Disconnect")
		return
	}
	fmt.Println("Returning from SimpleConnectAndValidateBootstrappedEntities - successfully disconnected.")
}

func SimpleConnectAndGetServerMetadata() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.Connect")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphObjectFactory")
		return
	}
	if gof == nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - Graph Object Factory is null")
		return
	}

	//if prefetchMetaData {
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphMetadata")
		return
	}
	fmt.Printf("Inside from SimpleConnectAndGetServerMetadata - read meta data as '%+v'", gmd)

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.Disconnect")
		return
	}
	fmt.Println("Returning from SimpleConnectAndGetServerMetadata - successfully disconnected.")
}

func SimpleConnectAndDisconnect() {
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(url, "", password, nil)
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndDisconnect - error during CreateConnection")
		return
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndDisconnect - error during conn.Connect")
		return
	}

	err = conn.Disconnect()
	if err != nil {
		fmt.Println("Returning from SimpleConnectAndDisconnect - error during conn.Disconnect")
		return
	}
	fmt.Println("Returning from SimpleConnectAndDisconnect - successfully disconnected.")
}

func SimpleConnection() {
	fmt.Println(">>>>>>> Inside SimpleConnection: About to Test SimpleConnectAndDisconnect() <<<<<<<")
	SimpleConnectAndDisconnect()
	fmt.Println(">>>>>>> Inside SimpleConnection: About to Test SimpleConnectAndGetServerMetadata() <<<<<<<")
	SimpleConnectAndGetServerMetadata()
	fmt.Println(">>>>>>> Inside SimpleConnection: About to Test SimpleConnectAndValidateBootstrappedEntities() <<<<<<<")
	SimpleConnectAndValidateBootstrappedEntities()
	fmt.Println(">>>>>>> Returning from SimpleConnection. <<<<<<<")
}
