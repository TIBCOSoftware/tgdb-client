package main

import (
	"crypto/x509"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/samples"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Hello, Welcome to the Playground!")

	/**
		Start TGDB Server with default initdb.conf and tgdb.conf
	*/
	samples.SimpleConnection()

	/**
		Restart TGDB Server, w/ default initdb.conf and tgdb.conf
	*/
	//samples.SimpleTransaction()

	/**
		Restart TGDB Server, w/ default initdb.conf and tgdb.conf, if you are executing this test suite in a standalone mode
		and uncomment QueryTransactions::createTestData(qConn) before executing the following in a standalone mode.
		Otherwise, let it continue after executing above test suite - SimpleTransaction()
	*/
	//samples.QueryTransactions()

	/**
		Restart TGDB Server, w/ default initdb.conf and tgdb.conf
	*/
	//samples.MultiTransactionTest()

	/**
		Restart TGDB Server, execute metadatascript.tql to create necessary nodes and attribute metadata via
		'${TGDB_HOME}/bin/tgdb-admin --file metadatascript.tql' before executing the following.
	*/
	//samples.MetadataTest()

	/**
		Restart TGDB Server with ancestry-initdb.conf and ancestry-tgdb.conf before executing the following.
	//*/
	//samples.AncestryTest()

	/**
		Restart TGDB Server with stemcell-initdb.conf and stemcell-tgdb.conf before executing the following.
		Make sure you also have data files available.
	*/
	//samples.StemCellMain()

	/**
		// TODO: Revisit later - Test it with data files in the proper format, if-n-when available
		Restart TGDB Server with stock-initdb.conf and stock-tgdb.conf before executing the following.
		Make sure you also have data files available.
	*/
	//samples.TradingDataMain()

	//testSSLConnectivity()

	fmt.Println("Play Over!!!")
}

func testSSLConnectivity() {
	//certPool := x509.SystemCertPool()

	sysTrustFile := fmt.Sprintf("%s%slib%ssecurity%scacerts", os.Getenv("JRE_HOME"), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))
	fmt.Printf("======> About to ReadFile '%+v'\n", sysTrustFile)
	rootPEM, err := ioutil.ReadFile(sysTrustFile)
	if err != nil {
		fmt.Printf("ERROR: Failed to read root certificate authority: %s\n", sysTrustFile)
	}
	//fmt.Printf("======> Got the certificate data as: '%+v'\n", string(rootPEM))

	certs, err := x509.ParseCertificateRequest(rootPEM)
	if err != nil {
		fmt.Printf("ERROR: Failed to parse root certificates: %s\n", sysTrustFile)
	}
	//fmt.Printf("======> Parsed %d certificates from system root\n", len(certs))
	fmt.Printf("======> Parsed certificate request '%+v' from system root\n", certs)

	//for _, cert := range certs {
	//	fmt.Printf("======> Adding the certificate data as: '%+v'\n", cert)
	//	certPool.AddCert(cert)
	//	//ok := certPool.AddCert(cert)
	//	//if !ok {
	//	//	fmt.Printf("ERROR: Failed to parse root certificate: %+v\n", ok)
	//	//	return
	//	//}
	//}
	//
	//config := &tls.Config{
	//	ClientCAs:          certPool,
	//	//RootCAs:            certPool,
	//	InsecureSkipVerify: false,
	//}
	//
	//host := "localhost"
	//port := 8223
	//serverAddr := fmt.Sprintf("%s:%d", host, port)
	//sslConn, cErr := tls.Dial("ssl", serverAddr, config)
	//if cErr != nil {
	//	fmt.Printf("ERROR: CreateSocket Failed to connect to the server at '%s' w/ error: '%+v'\n", serverAddr, cErr)
	//} else {
	//	fmt.Printf("======> Successfully created SSL connection for '%s' as '%+v'", serverAddr, sslConn)
	//}
}