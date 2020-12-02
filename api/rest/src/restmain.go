/*
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
 * File name: restmain.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id$
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tgdb"
	"tgdb/factory"
	"tgdb/impl"
	"tgdbrest"
	"time"
)

var topURLBase string
var authenticateURLBase string
var disconnectURLBase string
var hostPort string
var pingURLBase string
var transactionURLBase string
var queryURLBase string
var traverseURLBase string
var adminURLBase string
var metadataURLBase string
var entityURLBase string
var importExportURLBase string

var connPool tgdb.TGConnectionPool


var mapUSERNAME_AND_REMOTE_IP2TOKEN = make(map[string] tgdbrest.TGDBRestOAuthMetadataForConnection)

var TGDB_REST_PATH = "/TGDB/OData.svc/"
var TGDB_REST_FULL_PATH = "https://localhost:3001/TGDB/OData.svc/"
var TGDB_REST_METADATA_PATH = "/TGDB/OData.svc/$metadata"
var TGDB_REST_ATTR_DESC_PATH = "/TGDB/OData.svc/AttributeDescriptors"
var TGDB_REST_NODETYPES_PATH = "/TGDB/OData.svc/NodeTypes"
var TGDB_REST_EDGETYPES_PATH = "/TGDB/OData.svc/EdgeTypes"
var TGDB_REST_CONNECTIONS_PATH = "/TGDB/OData.svc/Connection"
var TGDB_REST_TYPES_PATH = "/TGDB/OData.svc/Types"
var TGDB_REST_TYPEDETAILS_PATH = "/TGDB/OData.svc/TypeDetails"
var hostPort4OData = ""
var dbURL4OData = ""
var TGDB_REST_API_USERNAME = "api"
var TGDB_REST_API_PASSWORD = "api"
var TGDB_REST_DB_CONNECTION_POOL_SIZE = 5
var HTTP_PROTOCOL = "http"

var engineNamePtr *string
var logDirPtr *string
var logLevelPtr *string
var logFileSizePtr *int
var logFileCountPtr *int
var logToConsolePtr *bool


var banner = "********************************************************************************\n" +
	"TIBCO(R) Graph Database {0}.\n" +
	"Copyright (c) 2016-2020 TIBCO Software Inc. All rights reserved.\n" +
	"{1} Enabled.\n" +
	"Please read the accompanying License and ReadMe documents;\n" +
	"your use of the software constitutes your acceptance of the terms contained in these documents.\n" +
	"********************************************************************************"

const (
	TGDB_API_ACCESS_KEY_ID = "TGDBAccessKeyId"
	TGDB_AUTH_TOKEN = "TGDBAuthToken"
)

var logger = impl.DefaultTGLogManager().GetLogger()

var mapAttributeTypeId2AttributeTypeName = map[int]string {
	impl.AttributeTypeInvalid:	"invalidType",
	impl.AttributeTypeBoolean:	"bool",
	impl.AttributeTypeByte:		"byte",
	impl.AttributeTypeChar:		"char",
	impl.AttributeTypeShort:     "short",
	impl.AttributeTypeInteger:   "integer",
	impl.AttributeTypeLong:      "long",
	impl.AttributeTypeFloat:     "float32",
	impl.AttributeTypeDouble:    "float64",
	impl.AttributeTypeNumber:    "number",
	impl.AttributeTypeString:    "string",
	impl.AttributeTypeDate:      "date",
	impl.AttributeTypeTime:      "time",
	impl.AttributeTypeTimeStamp: "time.Time",
	impl.AttributeTypeBlob:      "blob",
	impl.AttributeTypeClob:      "clob",
}

func createDisplayStringForVersionBanner (version *impl.TGClientVersion) string {
	//TIBCO(R) Graph Database 3.0.0 Build(0) Revision(4143) Enterprise Edition.
	return fmt.Sprintf("%d.%d.%d Build(%d) Revision(%d) %s Edition",
		version.GetMajor(), version.GetMinor(), version.GetUpdate(), version.GetBuildNo(), version.GetBuildRevision(), fromEditionTypetoEditionTypeString(version.GetEdition()))
}

func fromBuildTypetoBuildTypeString (buildType byte) string {
	switch buildType {
	case impl.BuildTypeBeta:
		{
			return "Beta"
		}
	case impl.BuildTypeEngineering:
		{
			return "Engineering-Build"
		}
	case impl.BuildTypeProduction:
		{
			return "Production"
		}
	default:
		{
			return "Beta"
		}
	}

}

func fromEditionTypetoEditionTypeString (edition byte) string {
	switch edition {
		case impl.EditionEvaluation:
			{
				return "Evaluation"
			}
		case impl.EditionCommunity:
			{
				return "Community"
			}
		case impl.EditionEnterprise:
			{
				return "Enterprise"
			}
		case impl.EditionDeveloper:
			{
				return "Developer"
			}
		default:
			{
				return "Evaluation"
			}
	}
}

func main() {

	version := impl.GetClientVersion()
	strVersion := createDisplayStringForVersionBanner (version)
	resultBanner := strings.Replace(banner, "{0}", strVersion, 1)
	resultBanner = strings.Replace(resultBanner, "{1}", fromBuildTypetoBuildTypeString(version.GetBuildType()), 1)
	fmt.Println(resultBanner)

	hostPortPtr := flag.String("listen", "localhost:9500", "Specify Server Listen Address")
	dbURLPtr := flag.String("dburl", "tcp://scott@localhost:8222/{dbName=demodb}", "Specify TGDB Database URL")
	randomName := "tgdb-rest-" + strconv.Itoa(os.Getpid())
	engineNamePtr = flag.String("name", randomName, "Specify Engine Name")
	logLevelPtr = flag.String("loglevel", "Info", "Specify Log Level")
	logDirPtr = flag.String("logdir", "NOT SET", "Specify Log Directory")

	logFileCountPtr = flag.Int("logfilecount", 10, "Specify Log File Count")
	logFileSizePtr = flag.Int("logfilesize", 10, "Specify Max Log File Size")
	logToConsolePtr = flag.Bool("logtoconsole", true, "Specify the messages on console")


	flag.Usage = func() {
		usageString := "usage: tgdb-rest [--listen <host:port>] [--dburl db_url] [--name name] [--loglevel Error|Warning|Info|Debug] [--logdir log_directory_path] [--logtoconsole true|false]\n\n"
		fmt.Fprintf(os.Stdout, usageString)
		fmt.Fprintf(os.Stdout, "optional arguments:\n")
		fmt.Fprintf(os.Stdout, "  --listen <host:port>  The host & port for server to listen to         (default \"localhost:9500\").\n")
		fmt.Fprintf(os.Stdout, "  --dburl db_name       The URL to connect to TIBCO Graph Database      (default \"tcp://scott@localhost:8222/{dbName=demodb}\").\n")
		fmt.Fprintf(os.Stdout, "  --name name           The name of the tgdb-rest server instance       (default \"tgdb-rest-<pid>\").\n")
		fmt.Fprintf(os.Stdout, "  --loglevel            The loglevel to set (Error|Warning|Info|Debug)  (default \"Warning\").\n")
		fmt.Fprintf(os.Stdout, "  --logdir              The directory to store logs                     (default \"Current Working Directory)\"\n")
		//fmt.Fprintf(os.Stdout, "  --logfilesize         The max filesize of each log file in MB         (default 10 MB)\n")
		//fmt.Fprintf(os.Stdout, "  --logfilecount        The max number of log files (rollover after threashold) (default 10)\n")
		fmt.Fprintf(os.Stdout, "  --logtoconsole        The boolean value to log messages on console    (default true)\n")
		fmt.Fprintf(os.Stdout, "  -h, --help            Show this help message.\n")
	}

	flag.Parse()
	//slogger.Log("Engine-Name: " + *engineNamePtr)
	hostPort4OData = *hostPortPtr
	dbURL4OData = *dbURLPtr

	*logFileSizePtr = *logFileSizePtr * 1000000
	//*logFileSizePtr = *logFileSizePtr * 1000

	logLvl := strings.ToLower(*logLevelPtr)
	if strings.Compare(logLvl, "error") == 0 {
		logger.SetLogLevel(tgdb.ErrorLog)
	} else if strings.Compare(logLvl, "warning") == 0 {
		logger.SetLogLevel(tgdb.WarningLog)
	} else if strings.Compare(logLvl, "info") == 0 {
		logger.SetLogLevel(tgdb.InfoLog)
	} else if strings.Compare(logLvl, "debug") == 0 {
		logger.SetLogLevel(tgdb.DebugLog)
	} else {
		logger.SetLogLevel(tgdb.WarningLog)
	}

	if *logToConsolePtr {
		logger.SetLogWriter(os.Stdout)
	} else {
		var logBaseDr string
		var err1 error
		if strings.Compare(*logDirPtr, "NOT SET") == 0 {
			logBaseDr, err1 = os.Getwd()
			if err1 != nil {
				fmt.Println("Error: " + err1.Error())
				return
			}
			logBaseDr = logBaseDr + string(os.PathSeparator) + "log"
		} else {
			logBaseDr = *logDirPtr
		}

		if len(logBaseDr) > 0 {
			if string(logBaseDr[len(logBaseDr)-1]) == string(os.PathSeparator) {
				logBaseDr = logBaseDr[0 : len(logBaseDr)-1]
			}
		}

		error := logger.SetLogBaseDir(logBaseDr)
		if error != nil {
			fmt.Println("Error: " + error.Error())
			return
		}

		error = logger.SetFileNameBase(*engineNamePtr)
		if error != nil {
			fmt.Println("Error: " + error.Error())
			return
		}
		logger.SetFileCount(*logFileCountPtr)
		logger.SetFileSize(*logFileSizePtr)

		logger.Info("Log Directory:  " + logBaseDr)

		logger.Info("Log File Size:  " + strconv.Itoa(*logFileSizePtr) + " Bytes")

		logger.Info("Log File Count: " + strconv.Itoa(*logFileCountPtr))
	}
	logger.SetLogPrefix("")

	logger.Info("LogLevel:       " + logLvl)
	logger.Info("TIBCO Graph Database REST Server Starting...")

	initializeSetupData ();
	registerConnectDisconnectURL()
	registerMetadataURL ()
	registerQueryURL ()
	registerTransactionURL ()

	registerODataURL()
	registerVizFileServURL()

	err := initTGDBConnectionPool ()

	if err != nil {
		logger.Error("Error starting the tgdb-rest server: " + err.Error())
		return
	}
	logDBSpecificInfo ()

	logger.Info("TIBCO Graph Database REST Server Running At: " + hostPort4OData)
	error := http.ListenAndServe(hostPort4OData, nil)
	if error != nil {
		logger.Error("TIBCO Graph Database REST Server failed at: " + hostPort4OData)
		logger.Error("Error Message: " + error.Error())
	}
}

func logDBSpecificInfo() {
	logger.Info("Connected to TIBCO Graph Database URL: " + dbURL4OData)
	logger.Info("Connection Pool Size: " + strconv.Itoa(TGDB_REST_DB_CONNECTION_POOL_SIZE))
}

func registerTransactionURL() {
	http.HandleFunc(transactionURLBase, transactionURLHandler)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + transactionURLBase)
}

func transactionURLHandler(w http.ResponseWriter, r *http.Request) {
	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	verb := headers["Verb"]
	if len(verb) < 1 || strings.Compare(strings.ToLower(verb), "createnode") == 0 {
		tgdbrest.RESTTransaction(connection, w, r, headers, body)
	}
	resetConnectionWithPrevToken(connection, prevToken)

}

func initTGDBConnectionPool() tgdb.TGError {
	connFactory := impl.NewTGConnectionFactory()
	var err tgdb.TGError

	connPool, err = connFactory.CreateConnectionPoolWithType(dbURL4OData, TGDB_REST_API_USERNAME, TGDB_REST_API_PASSWORD, TGDB_REST_DB_CONNECTION_POOL_SIZE, nil, tgdb.TypeAdmin)
	if (err != nil) {
		logger.Error(err.Error())
		return err
	}

	err = connPool.Connect();
	if (err != nil) {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func registerTopURL() {
	// register the topURLBase
	http.HandleFunc(topURLBase, topURLHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + topURLBase)
}

func registerVizFileServURL() {
	// register the topURLBase
	http.HandleFunc(topURLBase + "viz/images/", vizFileServHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + topURLBase + "viz/images/")

	http.HandleFunc(topURLBase + "viz/scripts/", vizFileServHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + topURLBase + "viz/scripts/")

	http.HandleFunc(topURLBase + "viz/css/", vizFileServHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + topURLBase + "viz/css/")

	http.HandleFunc(topURLBase + "viz/spotfire/", vizFileServHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + topURLBase + "viz/spotfire/")
}



func registerConnectDisconnectURL() {
	// register the authentication URL for authentication
	http.HandleFunc(authenticateURLBase, authenticationURLHandler)
	logger.Info("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + authenticateURLBase)

	//logger.Info ("Registering: " + disconnectURLBase)
	//http.HandleFunc(disconnectURLBase, disconnectURLHandler)
}

func registerMetadataURL() {

	// register the metadata endpoint URL for NodeTypes
	http.HandleFunc(metadataURLBase + "NodeTypes/", metadataURLHandler4NodeTypes)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + metadataURLBase + "NodeTypes/")

	// register the metadata endpoint URL for EdgeTypes
	http.HandleFunc(metadataURLBase + "EdgeTypes/", metadataURLHandler4EdgeTypes)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + metadataURLBase + "EdgeTypes/")

	// register the metadata endpoint URL for AttributeDescriptors
	http.HandleFunc(metadataURLBase + "AttributeDescriptors/", metadataURLHandler4AttributeDescriptors)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + metadataURLBase + "AttributeDescriptors/")

	http.HandleFunc(metadataURLBase + "Users/", metadataURLHandler4Users)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + metadataURLBase + "Users/")

	http.HandleFunc(metadataURLBase + "Connections/", metadataURLHandler4Connections)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + metadataURLBase + "Connections/")
}

func registerQueryURL() {
	// register the Query endpoint URL
	http.HandleFunc(queryURLBase, queryURLHandler)
	logger.Info ("Registered REST URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + queryURLBase)
}

/*
func registerAdminURL() {
	logger.Log ("Registering: " + adminURLBase + "Info/")
	http.HandleFunc(adminURLBase + "Info/", adminURLHandler4Info)
}
*/

func initializeSetupData() {

	topURLBase = "/TGDB/"

	authenticateURLBase = topURLBase + "Authenticate" + "/"
	disconnectURLBase = topURLBase + "Disconnect" + "/"
	pingURLBase = topURLBase + "Ping" + "/"
	transactionURLBase = topURLBase + "Transaction" + "/"
	queryURLBase = topURLBase + "Query" + "/"
	traverseURLBase = topURLBase + "Traverse" + "/"
	adminURLBase = topURLBase + "Admin" + "/"
	metadataURLBase = topURLBase + "Metadata" + "/"
	entityURLBase = topURLBase + "Entity" + "/"
	importExportURLBase = topURLBase + "ImportExport" + "/"
}

func vizFileServHandler (w http.ResponseWriter, r *http.Request) {

	if strings.Compare(strings.ToLower("Options"), strings.ToLower(r.Method)) == 0 {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		w.WriteHeader(204)
		return
	}

	substring := r.URL.Path[5:]
	substring = ".." + substring

	dat, err := ioutil.ReadFile(substring)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Resource Not Found: " + r.URL.Path + ".\n"))
		return
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Write(dat)
	}
}

func topURLHandler (w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w,"Request Method is: %s", r.Method);

	_, _, _, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	fmt.Fprintf(w, createURLsForAllEndpoints(topURLBase))
}

func metadataURLHandler4NodeTypes (w http.ResponseWriter, r *http.Request) {

	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}
	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}

	tgdbrest.MetadataQueryForNodeTypes(connection, w, r, headers, bodyMap);
	resetConnectionWithPrevToken(connection, prevToken)
}

func resetConnectionWithPrevToken(conn tgdb.TGConnection,  prevToken int64) tgdb.TGError {
	conn.GetChannel().SetAuthToken(prevToken)
	logger.Debug("ConnectionPool Release " + strconv.FormatInt(conn.GetConnectionId(), 10))
	_, err := connPool.ReleaseConnection(conn)

	if err != nil {
		logger.Error("error:" + err.Error())
		return err
	}
	return nil
}

func setConnectionWithUserToken(nToken int64) (tgdb.TGConnection, int64, tgdb.TGError) {
	conn, err := connPool.Get()
	logger.Debug ("ConnectionPool Get " + strconv.FormatInt(conn.GetConnectionId(), 10))

	if err != nil {
		logger.Error("error:" + err.Error())
		return nil, -1, err
	}

	channel := conn.GetChannel()
	prevToken := channel.GetAuthToken()
	channel.SetAuthToken(nToken)
	return conn, prevToken, nil
}

func metadataURLHandler4EdgeTypes (w http.ResponseWriter, r *http.Request) {

	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}
	tgdbrest.MetadataQueryForEdgeTypes(connection, w, r, headers, bodyMap);
	resetConnectionWithPrevToken(connection, prevToken)
}

func metadataURLHandler4AttributeDescriptors (w http.ResponseWriter, r *http.Request) {
	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}


	verb := headers["Verb"]
	if len(verb) < 1 || strings.Compare(strings.ToLower(verb), "get") == 0 {
		tgdbrest.MetadataQueryForAttributeDescriptors(connection, w, r, headers, bodyMap)
	} else if strings.Compare(strings.ToLower(verb), "create") == 0 {
		tgdbrest.MetadataCreateForAttributeDescriptors(connection, w, r, headers, bodyMap)
	}
	resetConnectionWithPrevToken(connection, prevToken)
}

func metadataURLHandler4Users (w http.ResponseWriter, r *http.Request) {
	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}


	verb := headers["Verb"]
	if len(verb) < 1 || strings.Compare(strings.ToLower(verb), "get") == 0 {
		tgdbrest.MetadataQueryForUsers(connection, w, r, headers, bodyMap)
	}

	resetConnectionWithPrevToken(connection, prevToken)

}

func metadataURLHandler4Connections (w http.ResponseWriter, r *http.Request) {
	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}

	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}

	verb := headers["Verb"]
	if len(verb) < 1 || strings.Compare(strings.ToLower(verb), "get") == 0 {
		tgdbrest.MetadataQueryForConnections(connection, w, r, headers, bodyMap)
	}

	resetConnectionWithPrevToken(connection, prevToken)

}


func queryURLHandler (w http.ResponseWriter, r *http.Request) {

	t_start := time.Now()

	if strings.Compare(strings.ToLower("Options"), strings.ToLower(r.Method)) == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(204)
		return
	}

	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}
	connection, prevToken, err := setConnectionWithUserToken (nToken)
	if err != nil {
		return
	}

	bodyMap := make(map[string]string)

	for k, v := range body {
		bodyMap[k] = v.(string)
	}

	tgdbrest.Query(connection, w, r, headers, bodyMap)
	t_end := time.Now()
	diff := t_end.Sub(t_start)

	queryStr, ok := bodyMap["GremlinQuery"]
	if ok {
		logger.Debug("Query String: " + queryStr)
		logger.Debug("Query Execution Time: " + diff.String())
	}

	resetConnectionWithPrevToken(connection, prevToken)
}



func createURLsForAllEndpoints (topURLBase string) string {

	var resultString bytes.Buffer

	restEndpointsInfo := tgdbrest.RESTEndpointsInfo{
		Vendor:       "TIBCO Software, Inc.",
		}

	restEndpointsInfo.EndpointInfo[0].URL = hostPort + pingURLBase
	restEndpointsInfo.EndpointInfo[0].Description = "This is a Ping Endpoint"

	restEndpointsInfo.EndpointInfo[1].URL = hostPort + transactionURLBase
	restEndpointsInfo.EndpointInfo[1].Description = "This is a Transaction Endpoint"

	restEndpointsInfo.EndpointInfo[2].URL = hostPort + queryURLBase
	restEndpointsInfo.EndpointInfo[2].Description = "This is a Query Endpoint"

	restEndpointsInfo.EndpointInfo[3].URL = hostPort + traverseURLBase
	restEndpointsInfo.EndpointInfo[3].Description = "This is a Traverse Endpoint"

	restEndpointsInfo.EndpointInfo[4].URL = hostPort + adminURLBase
	restEndpointsInfo.EndpointInfo[4].Description = "This is an Admin Endpoint"

	restEndpointsInfo.EndpointInfo[5].URL = hostPort + metadataURLBase
	restEndpointsInfo.EndpointInfo[5].Description = "This is a Metadata Endpoint"

	restEndpointsInfo.EndpointInfo[6].URL = hostPort + entityURLBase
	restEndpointsInfo.EndpointInfo[6].Description = "This is an Entity Endpoint"

	restEndpointsInfo.EndpointInfo[7].URL = hostPort + importExportURLBase
	restEndpointsInfo.EndpointInfo[7].Description = "This is an Import Export Endpoint"

	restEndpointsInfo.EndpointInfo[8].URL = hostPort + authenticateURLBase
	restEndpointsInfo.EndpointInfo[8].Description = "This is Authentication Endpoint"

	restEndpointsInfo.EndpointInfo[9].URL = hostPort + disconnectURLBase
	restEndpointsInfo.EndpointInfo[9].Description = "This is Disconnect Endpoint"

	b, err := json.MarshalIndent(restEndpointsInfo, "", "\t")
	if err != nil {
		logger.Error("error:" + err.Error())
	}
	resultString.WriteString(string(b))
	return resultString.String()
}

func authenticationURLHandler (w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
	}

	var err tgdb.TGError

	conn, err := connPool.Get()
	if err != nil {
		logger.Error("Returning from SimpleConnectAndDisconnect - error during conn.Connect: " + err.Error())
		return
	}

	defer connPool.ReleaseConnection(conn)

	prevToken := conn.GetChannel().GetAuthToken()

	absChannel, ok := conn.GetChannel().(*impl.TCPChannel)
	if ok {
		absChannel.SetChannelUserName(user)
		absChannel.SetChannelPassword([]byte(pass))
	}

	err = absChannel.DoAuthenticateForRESTConsumer()
	if err != nil {
		// Handle the error here
	}
	userToken := absChannel.GetAuthToken()

	absChannel.SetAuthToken(prevToken)

	var authResponse tgdbrest.TGDBRestAuthenticateResponse
	authResponse.Token = strconv.FormatInt(userToken, 10)
	b, error := json.MarshalIndent(authResponse, "", "\t")
	if error != nil {
		logger.Error("error:" + err.Error())
	}
	fmt.Fprintf(w, string(b))
}


//func disconnectURLHandler (w http.ResponseWriter, r *http.Request) {
//
//	_, _, nToken, bResult := isAuthenticRequest(w, r)
//	if !bResult {
//		return
//	}
//
//	connection := mapToken2Connection[nToken]
//	connection.Disconnect()
//	delete(mapToken2Connection, nToken)
//	var disconnectResponse tgdbrest.DisconnectResponse
//	disconnectResponse.Description = "Client Disconnected Successfully."
//
//	b, err := json.MarshalIndent(disconnectResponse, "", "\t")
//	if err != nil {
//		logger.Error("error:" + err.Error())
//	}
//	fmt.Fprintf(w, string(b))
//}



func isAuthenticRequest (w http.ResponseWriter, r *http.Request) (map[string] string, map[string] interface{}, int64, bool) {

	var restNodeTypesRequest tgdbrest.TGDBRestRequest

	err := json.NewDecoder(r.Body).Decode(&restNodeTypesRequest)
	if err != nil {
		logger.Error("Error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil, -1, false
	}

	api_auth_token := restNodeTypesRequest.Headers["Token"]

	if api_auth_token == "" {
		var tgdbError tgdbrest.TGDBRESTError
		tgdbError.ErrorMessage = "Please specify valid " + TGDB_AUTH_TOKEN + " as a query parameter in the request."

		b, err := json.MarshalIndent(tgdbError, "", "\t")
		if err != nil {
			logger.Error("error:" + err.Error())
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fmt.Fprintf(w, string(b))
		return nil, nil, -1, false
	}

	nToken, err := strconv.ParseInt(api_auth_token, 10, 64)
	if err != nil {
		var tgdbError tgdbrest.TGDBRESTError
		tgdbError.ErrorMessage = "Please specify valid " + TGDB_AUTH_TOKEN + " as a query parameter in the request."

		b, err := json.MarshalIndent(tgdbError, "", "\t")
		if err != nil {
			logger.Error("error:" + err.Error())
		}
		fmt.Fprintf(w, string(b))
		return nil, nil, -1, false
	}
	return restNodeTypesRequest.Headers, restNodeTypesRequest.Body, nToken, true
}

/*
func adminURLHandler4Info (w http.ResponseWriter, r *http.Request) {
	headers, body, nToken, bResult := isAuthenticRequest(w, r)
	if !bResult {
		return
	}
	connection := mapToken2Connection[nToken]
	tgdbrest.AdminURL4Info(connection, w, r, headers, body);

}
*/


func registerODataURL () {
	http.HandleFunc(TGDB_REST_PATH, oDataURLHandler)
	logger.Info("Registered OData URL: " + HTTP_PROTOCOL + "://" + hostPort4OData + TGDB_REST_PATH)
}

func oDataURLHandler (w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
	}

	conn, err := connPool.Get()

	if err != nil {
		logger.Error("Error during Connect to TGDB Server: " + err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	prevToken := conn.GetChannel().GetAuthToken()

	absChannel, ok := conn.GetChannel().(*impl.TCPChannel)
	if ok {
		absChannel.SetChannelUserName(user)
		absChannel.SetChannelPassword([]byte(pass))
	}

	err = absChannel.DoAuthenticateForRESTConsumer()
	if err != nil {
		// Handle the error here
	}
	userToken := absChannel.GetAuthToken()
	absChannel.SetAuthToken(prevToken)

	var tgdbRestOAuthMetadata tgdbrest.TGDBRestOAuthMetadataForConnection
	//var currentConnection tgdb.TGConnection
	
	tgdbRestOAuthMetadata = tgdbrest.TGDBRestOAuthMetadataForConnection{
		Token:               strconv.FormatInt(userToken, 10),
		UserName:            user,
		ConnectionTimestamp: time.Now().String(),
	}

	if strings.Compare(r.URL.Path, TGDB_REST_PATH) == 0 {
		fmt.Fprintf(w, readSvcCollectionContent())
	}
	if strings.Compare(r.URL.Path, TGDB_REST_METADATA_PATH) == 0 {
		fmt.Fprintf(w, readSvcMetadataCollectionContent())
	}
	if strings.Compare(r.URL.Path, TGDB_REST_CONNECTIONS_PATH) == 0 {
		fmt.Fprintf(w, readSvcConnections(tgdbRestOAuthMetadata))
	}
	if strings.Compare(r.URL.Path, TGDB_REST_TYPES_PATH) == 0 {
		content := readSvcTypes(conn)
		if content == "" {
			//TODO: Handle error with more specific information
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error.\n"))
			return
		}
		fmt.Fprintf(w, content)
	}
	if strings.Compare(r.URL.Path, TGDB_REST_TYPEDETAILS_PATH) == 0 {
		//content := readSvcTypeDetails(currentConnection)
		content := readSvcTypeDetails(conn)
		if content == "" {
			//TODO: Handle error with more specific information
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error.\n"))
			return
		}
		fmt.Fprintf(w, content)
	}
	resetConnectionWithPrevToken(conn, prevToken)
}

func readSvcCollectionContent () string {

	//dat, err := ioutil.ReadFile("G:\\workspace\\svn\\sgdb\\3.0\\api\\rest\\src\\odata_collection.xml")

	serviceContents := "<service xmlns=\"http://www.w3.org/2007/app\" xmlns:atom=\"http://www.w3.org/2005/Atom\" xml:base=\"https://localhost:3001/TGDB/OData.svc/\">\n" +
		"\t<workspace>\n" +
		"\t\t<atom:title>Default</atom:title>\n" +
		"\t\t<collection href=\"Connection\">\n" +
		"\t\t\t<atom:title>Connection</atom:title>\n" +
		"\t\t</collection>\n" +

		"\t\t<collection href=\"Types\">\n" +
		"\t\t\t<atom:title>Types</atom:title>\n" +
		"\t\t</collection>\n" +

		"\t\t<collection href=\"TypeDetails\">\n" +
		"\t\t\t<atom:title>TypeDetails</atom:title>\n" +
		"\t\t</collection>\n" +

		"\t</workspace>\n" +
		"</service>"

	return serviceContents
}

func readSvcMetadataCollectionContent () string {

	//dat, err := ioutil.ReadFile("G:\\workspace\\svn\\sgdb\\3.0\\api\\rest\\src\\spotfireset\\tgdb_odata_metadata_collection.xml")

	serviceMetadata := "<edmx:Edmx xmlns:edmx=\"http://schemas.microsoft.com/ado/2007/06/edmx\" Version=\"1.0\">\n" +
		"\t<edmx:DataServices xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" m:DataServiceVersion=\"3.0\" m:MaxDataServiceVersion=\"3.0\">\n" +
		"\t\t<Schema xmlns=\"http://schemas.microsoft.com/ado/2009/11/edm\" Namespace=\"TGDBSpotfire\">\n" +
		"\t\t\t<EntityType Name=\"VirtualConnection\">\n" +
		"\t\t\t\t<Key>\n" +
		"\t\t\t\t\t<PropertyRef Name=\"ConnectionURL\"/>\n" +
		"\t\t\t\t</Key>\n" +
		"\t\t\t\t<Property Name=\"ConnectionURL\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"Token\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"UserName\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"ConnectionTimestamp\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t</EntityType>\n" +


		"\t\t\t<EntityType Name=\"Type\">\n" +
		"\t\t\t\t<Key>\n" +
		"\t\t\t\t\t<PropertyRef Name=\"Name\"/>\n" +
		"\t\t\t\t</Key>\n" +
		"\t\t\t\t<Property Name=\"Name\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"Type\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"SysId\" Type=\"Edm.Int32\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"EntryCount\" Type=\"Edm.Int32\" Nullable=\"false\"/>\n" +
		"\t\t\t</EntityType>\n" +

		"\t\t\t<EntityType Name=\"TypeDetail\">\n" +
		"\t\t\t\t<Key>\n" +
		"\t\t\t\t\t<PropertyRef Name=\"Name\"/>\n" +
		"\t\t\t\t</Key>\n" +
		"\t\t\t\t<Property Name=\"Name\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"SysId\" Type=\"Edm.Int32\" Nullable=\"false\"/>\n" +
		"\t\t\t\t<Property Name=\"Type\" Type=\"Edm.String\" Nullable=\"false\"/>\n" +
		"\t\t\t</EntityType>\n" +


		"\t\t\t<EntityContainer Name=\"DemoService\" m:IsDefaultEntityContainer=\"true\">\n" +
		"\t\t\t\t<EntitySet Name=\"Connection\" EntityType=\"TGDBSpotfire.VirtualConnection\"/>\n" +
		"\t\t\t\t<EntitySet Name=\"Types\" EntityType=\"TGDBSpotfire.Type\"/>\n" +
		"\t\t\t\t<EntitySet Name=\"TypeDetails\" EntityType=\"TGDBSpotfire.TypeDetail\"/>\n" +
		"\t\t\t</EntityContainer>\n" +

		"\t\t</Schema>\n" +
		"\t</edmx:DataServices>\n" +
		"</edmx:Edmx>"

	return serviceMetadata
}

/*
func readSvcAdvertisementsContent () string {
	dat, err := ioutil.ReadFile("G:\\workspace\\svn\\sgdb\\3.0\\api\\rest\\src\\Advertisements.xml")
	check(err)
	//fmt.Print(string(dat))
	return string(dat)
}
*/


func readSvcConnections (tgdbRESTOAuthMetadata tgdbrest.TGDBRestOAuthMetadataForConnection) string {
	connectionURL := "http://" + hostPort4OData + "/TGDB/"

	connectionsContent := "<feed xml:base=\"https://services.odata.org/V3/OData/OData.svc/\" xmlns=\"http://www.w3.org/2005/Atom\" xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" xmlns:georss=\"http://www.georss.org/georss\" xmlns:gml=\"http://www.opengis.net/gml\">\n" +
		"\t<id>https://services.odata.org/V3/OData/OData.svc/Connection</id>\n" +
		"\t<title type=\"text\">Connection</title>\n" +
		"\t<updated>2020-06-17T21:06:39Z</updated>\n" +
		"\t<link rel=\"self\" title=\"Connection\" href=\"Connection\" />\n" +
		"\t<entry>\n" +
		"\t\t<id>https://services.odata.org/V3/OData/OData.svc/Connection</id>\n" +
		"\t\t<category term=\"TGDBSpotfire.VirtualConnection\" scheme=\"http://schemas.microsoft.com/ado/2007/08/dataservices/scheme\" />\n" +
		"\t\t<content type=\"*/*\" src=\"Connection/$value\" />\n" +
		"\t\t<m:properties>\n" +
		"\t\t\t<d:ConnectionURL m:type=\"Edm.String\">" + connectionURL + "</d:ConnectionURL>\n" +
		"\t\t\t<d:Token m:type=\"Edm.String\">" + tgdbRESTOAuthMetadata.Token + "</d:Token>\n" +
		"\t\t\t<d:UserName m:type=\"Edm.String\">" + tgdbRESTOAuthMetadata.UserName + "</d:UserName>\n" +
		"\t\t\t<d:ConnectionTimestamp m:type=\"Edm.String\">" + tgdbRESTOAuthMetadata.ConnectionTimestamp + "</d:ConnectionTimestamp>\n" +
		"\t\t</m:properties>\n" +
		"\t</entry>\n" +
		"</feed>"

	return connectionsContent
}

func readSvcTypeDetails (conn tgdb.TGConnection) string {
	entries := formTGDBRestOAuthMetadataForTypeDetails(conn)
	if entries == nil {
		return ""
	}
	return createODataStringForTypeDetails (entries)
}

func readSvcTypes (conn tgdb.TGConnection) string {
	entries := formTGDBRestOAuthMetadataForTypes(conn)
	if entries == nil {
		return ""
	}
	return createODataStringForTypes (entries)
}


func createODataStringForTypes(types []tgdbrest.TGDBRestOAuthMetadataForTypes) string {

	typesContent := "<feed xml:base=\"https://services.odata.org/V3/OData/OData.svc/\" xmlns=\"http://www.w3.org/2005/Atom\" xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" xmlns:georss=\"http://www.georss.org/georss\" xmlns:gml=\"http://www.opengis.net/gml\">\n" +
		"\t<id>https://services.odata.org/V3/OData/OData.svc/Types</id>\n" +
		"\t<title type=\"text\">Types</title>\n" +
		"\t<updated>2020-06-17T21:06:39Z</updated>\n" +
		"\t<link rel=\"self\" title=\"Types\" href=\"Types\" />\n"

	for i := 0; i < len(types); i++ {
		typesContent = typesContent + createODataStringForSingleType (types[i])
	}

	typesContent = typesContent + "</feed>"

	return typesContent
}

func createODataStringForSingleType(entityType tgdbrest.TGDBRestOAuthMetadataForTypes) string {
	resultContent := "\t<entry>\n" +
		"\t\t<id>https://services.odata.org/V3/OData/OData.svc/Types</id>\n" +
		"\t\t<category term=\"TGDBSpotfire.Type\" scheme=\"http://schemas.microsoft.com/ado/2007/08/dataservices/scheme\" />\n" +
		"\t\t<content type=\"*/*\" src=\"Type/$value\" />\n" +
		"\t\t<m:properties>\n" +
		"\t\t\t<d:Name m:type=\"Edm.String\">" + entityType.Name + "</d:Name>\n" +
		"\t\t\t<d:Type m:type=\"Edm.String\">" + entityType.Type + "</d:Type>\n" +
		"\t\t\t<d:SysId m:type=\"Edm.Int32\">" + strconv.Itoa(entityType.SysId) + "</d:SysId>\n" +
		"\t\t\t<d:EntryCount m:type=\"Edm.Int32\">" + strconv.Itoa(int(entityType.EntryCount)) + "</d:EntryCount>\n" +
		"\t\t</m:properties>\n" +
		"\t</entry>\n"

	return resultContent
}

func createODataStringForSingleTypeDetail(entityType tgdbrest.TGDBRestOAuthMetadataForTypeDetails) string {
	resultContent := "\t<entry>\n" +
		"\t\t<id>https://services.odata.org/V3/OData/OData.svc/TypeDetails</id>\n" +
		"\t\t<category term=\"TGDBSpotfire.TypeDetail\" scheme=\"http://schemas.microsoft.com/ado/2007/08/dataservices/scheme\" />\n" +
		"\t\t<content type=\"*/*\" src=\"TypeDetail/$value\" />\n" +
		"\t\t<m:properties>\n" +
		"\t\t\t<d:Name m:type=\"Edm.String\">" + entityType.Name + "</d:Name>\n" +
		"\t\t\t<d:SysId m:type=\"Edm.Int32\">" + strconv.Itoa(entityType.SysId) + "</d:SysId>\n" +
		"\t\t\t<d:Type m:type=\"Edm.String\">" + entityType.Type + "</d:Type>\n" +
		"\t\t</m:properties>\n" +
		"\t</entry>\n"

	return resultContent
}


func formTGDBRestOAuthMetadataForTypes (conn tgdb.TGConnection) []tgdbrest.TGDBRestOAuthMetadataForTypes {

	a := make([]tgdbrest.TGDBRestOAuthMetadataForTypes, 0)

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphObjectFactory")
		//return err
		//handle error here
		return nil
	}
	if gof == nil {
		// Handle error here
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - Graph Object Factory is null")
		//return
		return nil
	}

	//if prefetchMetaData {
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphMetadata")
		//return
		// handle error here
		return nil
	}


	// Handle NodeTypes
	nodeTypes, err := gmd.GetNodeTypes()
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	for i := 0; i < len(nodeTypes); i++ {
		var odataMetadataForTypes tgdbrest.TGDBRestOAuthMetadataForTypes

		odataMetadataForTypes.Name = nodeTypes[i].GetName()
		odataMetadataForTypes.Type = "N"
		odataMetadataForTypes.SysId = nodeTypes[i].GetEntityTypeId()
		odataMetadataForTypes.EntryCount = int32(nodeTypes[i].(*impl.NodeType).GetNumEntries())

		a = append (a, odataMetadataForTypes)
	}


	// Handle EdgeTypes
	edgeTypes, err := gmd.GetEdgeTypes()
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	for i := 0; i < len(edgeTypes); i++ {
		var odataMetadataForTypes tgdbrest.TGDBRestOAuthMetadataForTypes

		odataMetadataForTypes.Name = edgeTypes[i].GetName()
		odataMetadataForTypes.Type = "E"
		odataMetadataForTypes.SysId = edgeTypes[i].GetEntityTypeId()
		odataMetadataForTypes.EntryCount = int32(edgeTypes[i].(*impl.EdgeType).GetNumEntries())

		a = append (a, odataMetadataForTypes)
	}


	//fmt.Println(nodeTypes)
	return a
}

func formTGDBRestOAuthMetadataForTypeDetails (conn tgdb.TGConnection) []tgdbrest.TGDBRestOAuthMetadataForTypeDetails {

	a := make([]tgdbrest.TGDBRestOAuthMetadataForTypeDetails, 0)

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphObjectFactory")
		//return err
		//handle error here
		return nil
	}
	if gof == nil {
		// Handle error here
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - Graph Object Factory is null")
		//return
		return nil
	}

	//if prefetchMetaData {
	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		//fmt.Println("Returning from SimpleConnectAndGetServerMetadata - error during conn.GetGraphMetadata")
		//return
		// handle error here
		return nil
	}

	// Handle NodeTypes
	nodeTypes, err := gmd.GetNodeTypes()
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	// Handle NodeTypes
	for i := 0; i < len(nodeTypes); i++ {
		attrDescs := nodeTypes[i].GetAttributeDescriptors()
		for j := 0; j < len (attrDescs); j++ {
			var odataMetadataForTypeDetails tgdbrest.TGDBRestOAuthMetadataForTypeDetails
			odataMetadataForTypeDetails.Name = attrDescs[j].GetName()
			//fmt.Println(odataMetadataForTypeDetails.Name)
			odataMetadataForTypeDetails.SysId = nodeTypes[i].GetEntityTypeId()
			odataMetadataForTypeDetails.Type = mapAttributeTypeId2AttributeTypeName[attrDescs[j].GetAttrType()]//strconv.Itoa(attrDescs[j].GetAttrType())
			a = append (a, odataMetadataForTypeDetails)
		}
	}

	// Handle EdgeTypes
	edgeTypes, err := gmd.GetEdgeTypes()
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	for i := 0; i < len(edgeTypes); i++ {
		attrDescs := edgeTypes[i].GetAttributeDescriptors();
		for j := 0; j < len(attrDescs); j++ {
			var odataMetadataForTypeDetails tgdbrest.TGDBRestOAuthMetadataForTypeDetails
			odataMetadataForTypeDetails.Name = attrDescs[j].GetName()
			odataMetadataForTypeDetails.SysId = edgeTypes[i].GetEntityTypeId()
			odataMetadataForTypeDetails.Type = mapAttributeTypeId2AttributeTypeName[attrDescs[j].GetAttrType()]//strconv.Itoa(attrDescs[j].GetAttrType())
			a = append(a, odataMetadataForTypeDetails)
		}
	}
	return a
}



func createODataStringForTypeDetails(types []tgdbrest.TGDBRestOAuthMetadataForTypeDetails) string {

	typesContent := "<feed xml:base=\"https://services.odata.org/V3/OData/OData.svc/\" xmlns=\"http://www.w3.org/2005/Atom\" xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" xmlns:georss=\"http://www.georss.org/georss\" xmlns:gml=\"http://www.opengis.net/gml\">\n" +
		"\t<id>https://services.odata.org/V3/OData/OData.svc/TypeDetails</id>\n" +
		"\t<title type=\"text\">TypeDetails</title>\n" +
		"\t<updated>2020-06-17T21:06:39Z</updated>\n" +
		"\t<link rel=\"self\" title=\"TypeDetails\" href=\"TypeDetails\" />\n"

	for i := 0; i < len(types); i++ {
		typesContent = typesContent + createODataStringForSingleTypeDetail (types[i])
	}

	typesContent = typesContent + "</feed>"

	return typesContent
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createConnectionForOAuth (db_url string, user string, password string) (*tgdbrest.TGDBRestOAuthMetadataForConnection, tgdb.TGError) {

	var odataMetadata tgdbrest.TGDBRestOAuthMetadataForConnection
	connFactory := factory.GetConnectionFactory()
	//conn, err := connFactory.CreateAdminConnection(url + "{dbName=" + db_name + "}", user, password, nil)

	var mapConnectionProps = make(map[string] string)
	mapConnectionProps["tgdb.channel.clientId"] = *engineNamePtr



	conn, err := connFactory.CreateAdminConnection(db_url, user, password, mapConnectionProps)
	odataMetadata.ConnectionTimestamp = time.Now().String()
	odataMetadata.UserName = user
	if err != nil {
		//logger.Debug("Returning from SimpleConnectAndDisconnect - error during CreateConnection")
		logger.Error("Error: " + err.Error())
		return nil, err
	}

	err = conn.Connect()
	if err != nil {
		//logger.Debug("Returning from SimpleConnectAndDisconnect - error during conn.Connect")
		logger.Error("Error: " + err.Error())
		return nil, err
	}

	token := conn.GetChannel().GetAuthToken()

	odataMetadata.Token = strconv.FormatInt(token, 10)
	return &odataMetadata, nil

}