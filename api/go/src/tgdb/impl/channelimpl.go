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
 * File Name: channelimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: channelimpl.go 4585 2020-10-28 18:42:11Z nimish $
 */

package impl

import (
	"bytes"
	"crypto"
	"encoding/binary"

	//	"crypto/aes"
//	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"fmt"
//	"golang.org/x/crypto/blowfish"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"tgdb"
	"time"
)

// Default settings - A fallback option - to construct a valid connection URL to TGDB server
const (
	gDefaultHost string = "localhost"
	gDefaultPort int    = 8222
)

type LinkUrl struct {
	ftUrls   []tgdb.TGChannelUrl
	urlHost  string
	isIPv6   bool
	urlPort  int
	protocol tgdb.TGProtocol
	urlStr   string
	urlProps *SortedProperties // This is always sorted
	urlUser  string
}

func DefaultLinkUrl() *LinkUrl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(LinkUrl{})

	// Instantiate a single instance collection of all environment variable settings
	tgEnv := NewTGEnvironment()
	newChannelUrl := LinkUrl{
		ftUrls:   make([]tgdb.TGChannelUrl, 0),
		urlHost:  tgEnv.GetChannelDefaultHost(),
		isIPv6:   false,
		urlPort:  tgEnv.GetChannelDefaultPort(),
		protocol: tgdb.ProtocolTCP,
		urlProps: &SortedProperties{},
		urlUser:  tgEnv.GetChannelDefaultUser(),
	}

	// Get default FT hosts as configured as environment settings
	defaultFtUrls := tgEnv.GetChannelFTHosts()
	if len(defaultFtUrls) > 0 {
		ftUrls := strings.Split(defaultFtUrls, ",")
		for _, ftUrl := range ftUrls {
			if ftUrl != "" {
				newLinkUrl := NewLinkUrl(ftUrl)
				newChannelUrl.ftUrls = append(newChannelUrl.ftUrls, newLinkUrl)
			}
		}
	}
	return &newChannelUrl
}

func NewLinkUrl(sUrl string) *LinkUrl {
	newChannelUrl := DefaultLinkUrl()
	sUrl = strings.TrimSpace(sUrl)
	proto, host, port, ip6Flag, user, ftUrls, props, err := parseUrlComponents(sUrl)
	if err != nil {
		return nil
	}
	newChannelUrl.urlStr = sUrl
	newChannelUrl.isIPv6 = ip6Flag
	newChannelUrl.urlHost = host
	newChannelUrl.urlPort = port
	newChannelUrl.protocol = proto
	newChannelUrl.ftUrls = ftUrls
	if props != nil {
		newChannelUrl.urlProps = props.(*SortedProperties)
	}
	newChannelUrl.urlUser = user
	newChannelUrl.urlProps.AddProperty(GetConfigFromKey(ChannelUserID).GetName(), user)
	return newChannelUrl
}

func NewLinkUrlWithComponents(proto tgdb.TGProtocol, host string, port int) *LinkUrl {
	newChannelUrl := DefaultLinkUrl()
	newChannelUrl.urlHost = host
	newChannelUrl.urlPort = port
	newChannelUrl.protocol = proto
	return newChannelUrl
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGChannelUrl
/////////////////////////////////////////////////////////////////

func ParseChannelUrl(cUrl string) *LinkUrl {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:ParseChannelUrl w/ URL string as '%s'", cUrl))
	if cUrl == "" {
		return DefaultLinkUrl()
	}

	//logger.Log(fmt.Sprintf("Returning LinkUrl:ParseChannelUrl w/ URL string as '%s'", cUrl))
	return NewLinkUrl(cUrl)
}

/////////////////////////////////////////////////////////////////
// Private functions for TGChannelUrl
/////////////////////////////////////////////////////////////////

// "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
// "tcp://scott@foo.bar.com:8700"
// "tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}"
// "http://[2001:db8:1f70::999:de8:7648:6e8]:100/{userID=Admin}"
func parseUrlComponents(sUrl string) (tgdb.TGProtocol, string, int, bool, string, []tgdb.TGChannelUrl, tgdb.TGProperties, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseUrlComponents w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	protocol := tgdb.ProtocolTCP
	var user, host string
	var port int
	var ip6Flag bool
	ftUrls := make([]tgdb.TGChannelUrl, 0)

	if len(sUrl) == 0 {
		errMsg := fmt.Sprintf("Returning LinkUrl:parseUrlComponents as invalid length for input URL string '%s'", sUrl)
		err := GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	protocol, urlSubstr, err := parseProtocol(sUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Returning LinkUrl:parseUrlComponents - error in parsing protocol '%+v", err.Error()))
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	if len(urlSubstr) == 0 {
		errMsg := fmt.Sprintf("Returning LinkUrl:parseUrlComponents as incorrect Host/Port specified in input URL string '%s'", sUrl)
		err := GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	user, host, port, ip6Flag, propSubstr, err := parseUserAndHostAndPort(urlSubstr)
	if err != nil {
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	if len(propSubstr) == 0 {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning LinkUrl:parseUrlComponents as there are no properties specified - hence no FTUrls mentioned as part of property set"))
		}
		return protocol, host, port, ip6Flag, user, nil, nil, nil
	}
	props, ftUrls, err3 := parseProperties(propSubstr)
	if err3 != nil {
		logger.Error(fmt.Sprintf("Returning LinkUrl:parseUrlComponents - error in parsing properties '%+v", err.Error()))
		return protocol, host, port, ip6Flag, user, ftUrls, nil, err3
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning LinkUrl:parseUrlComponents w/ '%+v' '%+v' '%+v' '%+v' '%+v' '%+v' '%+v'", protocol, host, port, ip6Flag, user, ftUrls, props))
	}
	return protocol, host, port, ip6Flag, user, ftUrls, props, nil
}

func parseProtocol(sUrl string) (tgdb.TGProtocol, string, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseProtocol w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	var protocolStr, urlSubstring string
	protocol := tgdb.ProtocolTCP // Default
	if len(sUrl) == 0 {
		errMsg := "Unable to parse protocol from the channel URL string"
		return protocol, "", GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Expected format of sUrl is one of the following
	// (a) http://scott@foo.bar.com:8700/... 	                    <== protocol://User/DNS/Port/...
	// (b) https://scott@foo.bar.com/...	                        <== protocol://User/DNS/...
	// (c) tcp://scott@10.20.30.40:8123/...	                        <== protocol://User/IPv4/Port/...
	// (d) ssl://scott@10.20.30.40/...	                            <== protocol://User/IPv4//...

	idx := strings.IndexAny(sUrl, "://")
	if idx > 0 {
		strComponents := strings.Split(sUrl, "://")
		protocolStr = strings.ToLower(strComponents[0])
		urlSubstring = strComponents[1]
		switch protocolStr {
		case tgdb.ProtocolHTTP.String():
			protocol = tgdb.ProtocolHTTP
		case tgdb.ProtocolHTTPS.String():
			protocol = tgdb.ProtocolHTTPS
		case tgdb.ProtocolSSL.String():
			protocol = tgdb.ProtocolSSL
		case tgdb.ProtocolTCP.String():
			fallthrough
		default:
			protocol = tgdb.ProtocolTCP
		}
	} else {
		protocol = tgdb.ProtocolTCP
		urlSubstring = sUrl
	}

	//logger.Log(fmt.Sprintf("Returning LinkUrl:parseProtocol w/ URL string as '%+v' '%+v'", protocol, urlSubstring))
	return protocol, urlSubstring, nil
}

func parseUserAndHostAndPort(sUrl string) (string, string, int, bool, string, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseUserAndHostAndPort w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	var urlSubstring, user, host string
	var port int
	var ip6Flag bool
	var hostStr, hostPortStr, userHostPortStr string

	if len(sUrl) == 0 {
		errMsg := "Unable to parse user/host/port from the channel URL string"
		return user, host, port, ip6Flag, "", GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Expected format of sUrl is one of the following
	// (a) scott@foo.bar.com:8700/... 	                        <== User/DNS/Port/...
	// (b) scott@foo.bar.com/...	                            <== User/DNS/Port/...
	// (c) scott@10.20.30.40:8123/...	                        <== User/IPv4/Port/...
	// (d) scott@10.20.30.40/...	                            <== User/IPv4//...
	// (e) scott@[1234:567;2345:678::2342:453:2341]:3456/...	<== User/IPv6/Port/...
	// (f) scott@[1234:567:2345:678::2342:453:2341]/...	        <== User/IPv6//...
	// (g) foo.bar.com:8700/... 	                            <== DNS/Port/...
	// (h) foo.bar.com/...	                                    <== DNS/Port/...
	// (i) 10.20.30.40:8123/...	                                <== IPv4/Port/...
	// (j) 10.20.30.40/...	                                    <== IPv4/...
	// (k) 1234:567;2345:678::2342:453:2341]:3456/...	        <== IPv6/Port/...
	// (l) 1234:567:2345:678::2342:453:2341]/...	            <== IPv6//...

	//logger.Debug(fmt.Sprintf("ParseUserAndHostAndPort has incoming URL string as: '%+v'", sUrl))
	if strings.Contains(sUrl, "/") {
		strComponents := strings.Split(sUrl, "/")
		urlSubstring = strComponents[1] // Remaining string that may have properties in N1'V1,N2=V2,... format
		userHostPortStr = strComponents[0]
	} else {
		// At this point, there is no trailing '/' indicating it is probably one of the ftUrls or no properties
		userHostPortStr = sUrl
	}

	//logger.Debug(fmt.Sprintf("ParseUserAndHostAndPort userHostPortStr string as: '%+v'", userHostPortStr))
	// The following if block takes care of formats a through f - to extract user as part of host-port string
	if strings.Contains(userHostPortStr, "@") {
		userHostPortComps := strings.Split(userHostPortStr, "@")
		user = userHostPortComps[0]
		hostPortStr = userHostPortComps[1]
		if len(hostPortStr) == 0 {
			// This is the case for 'scott@/...'
			errMsg := "Invalid or missing host Name and/or port in the channel URL string"
			return user, host, port, ip6Flag, "", GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
		}
	} else {
		// At this point, user is still NOT specified as shown in formats g through l above
		hostPortStr = userHostPortStr
	}
	//logger.Debug(fmt.Sprintf("ParseUserAndHostAndPort hostPortStr string as: '%+v'", hostPortStr))

	// At this point, it is either host:port or only host w/o port - Hence separate out port first
	if strings.Contains(hostPortStr, ":") {
		hostPortComps := strings.Split(hostPortStr, ":")
		hostStr = hostPortComps[0]
		portstr := hostPortComps[1]
		if len(portstr) == 0 {
			port = gDefaultPort
		} else {
			port1, err := strconv.Atoi(portstr)
			if err != nil {
				errMsg := "Unable to parse user/host/port from the channel URL string"
				return user, host, port, ip6Flag, "", GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
			}
			port = port1
		}
	} else {
		// There is no port specified - so assigned it to the default port
		port = gDefaultPort
		hostStr = hostPortStr
	}

	// First determine whether the host is in Ipv6 format a.k.a. [xxxx:xxxx:xxxx:xxxx:xxxx]
	if strings.HasPrefix(hostStr, "[") {
		// This implies the host is in IPv6 format
		host = hostStr[1 : len(hostStr)-1] // Strip '[' and ']' from both ends
		if len(host) == 0 {
			host = gDefaultHost
		}
		ip6Flag = true
	} else {
		// This implies the host is in either IPv4 or DNS format
		hostPortComps := strings.Split(hostStr, ":")
		host = hostPortComps[0]
		if len(host) == 0 {
			host = gDefaultHost
		}
	}

	//logger.Log(fmt.Sprintf("Returning LinkUrl:parseUserAndHostAndPort returning '%s', '%s', '%d', '%+v' and remainder as: '%+v'", user, host, port, ip6Flag, urlSubstring))
	return user, host, port, ip6Flag, urlSubstring, nil
}

func parseProperties(sUrl string) (tgdb.TGProperties, []tgdb.TGChannelUrl, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseProperties w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	if len(sUrl) == 0 {
		errMsg := "Unable to parse properties from the channel URL string"
		return nil, nil, GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Expected format of sUrl is one of the following
	// {userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}

	// Intentionally do not set the following components to default values
	props := NewSortedProperties()
	ftUrls := make([]tgdb.TGChannelUrl, 0)

	if !strings.HasPrefix(sUrl, "{") {
		errMsg := "Unable to parse properties from the channel URL string"
		return nil, nil, GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	sUrl = sUrl[1 : len(sUrl)-1] // Strip '{' and '}' from both ends
	nvPairs := strings.Split(sUrl, ";")
	if len(nvPairs) == 0 {
		errMsg := "Malformed URL property specification - Must begin with { and end with }. All key=value must be separated with ;"
		return nil, nil, GetErrorByType(TGErrorProtocolNotSupported, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// At this point, nvPairs is expected to have each entry as 'Name=value' string
	for _, nvPair := range nvPairs {
		kvp := strings.Split(nvPair, "=")
		props.AddProperty(kvp[0], kvp[1])
	}
	//logger.Debug(fmt.Sprintf("Url Properties parsed from input URL are: '%+v'", props))

	// Once all the properties are scanned and parsed, check whether any of the supplied properties was ftHosts
	configName := PreDefinedConfigurations[ChannelFTHosts]
	if DoesPropertyExist(props, configName.aliasName) || DoesPropertyExist(props, configName.configPropName) {
		//logger.Debug(fmt.Sprintf("Url Properties does have an existing property named / aliased as 'ftHosts'"))
		cn := GetConfigFromName("ftHosts")
		if cn == nil || cn.GetName() == "" {
			ftUrls = nil
		} else {
			ftUrlStr := props.GetProperty(cn, "")
			ftUs := ftUrlStr
			//logger.Debug(fmt.Sprintf("ftus: '%+v'", ftus))
			if len(ftUs) == 0 {
				ftUrls = nil
			} else {
				fts := strings.Split(ftUs, ",")
				for _, fUrl := range fts {
					if fUrl != "" {
						//logger.Debug(fmt.Sprintf("ParseProperties trying to create new LinkURL for FtUrl '%+v'", fUrl))
						newLinkUrl := NewLinkUrl(fUrl)
						//logger.Debug(fmt.Sprintf("ParseProperties created new LinkURL '%+v' for FtUrl '%+v'", newLinkUrl, fUrl))
						ftUrls = append(ftUrls, newLinkUrl)
					}
				}
			}
		}
	}

	//logger.Log(fmt.Sprintf("Returning LinkUrl:ParseProperties returning Props as '%+v', and ftUrls as: '%+v'", props, ftUrls))
	return props, ftUrls, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannelUrl
/////////////////////////////////////////////////////////////////

// GetFTUrls gets the Fault Tolerant URLs
func (obj *LinkUrl) GetFTUrls() []tgdb.TGChannelUrl {
	return obj.ftUrls
}

// GetHost gets the host part of the URL
func (obj *LinkUrl) GetHost() string {
	return obj.urlHost
}

// GetPort gets the port on which it is connected
func (obj *LinkUrl) GetPort() int {
	return obj.urlPort
}

// GetProperties gets the URL Properties
func (obj *LinkUrl) GetProperties() tgdb.TGProperties {
	return obj.urlProps
}

// GetProtocol gets the protocol used as part of the URL
func (obj *LinkUrl) GetProtocol() tgdb.TGProtocol {
	return obj.protocol
}

// GetUrlAsString gets the string form of the URL
func (obj *LinkUrl) GetUrlAsString() string {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:GetUrlAsString"))
	if len(obj.urlStr) == 0 {
		if obj.isIPv6 {
			obj.urlStr = fmt.Sprintf("%s://%s@[%s]:%d", obj.protocol.String(), obj.urlUser, strings.ToLower(obj.urlHost), obj.urlPort)
		} else {
			obj.urlStr = fmt.Sprintf("%s://%s@%s:%d", obj.protocol.String(), obj.urlUser, strings.ToLower(obj.urlHost), obj.urlPort)
		}
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning LinkUrl:GetUrlAsString w/ URL string as '%s'", obj.urlStr))
	}
	return obj.urlStr
}

// GetUser gets the user associated with the URL
func (obj *LinkUrl) GetUser() string {
	if obj.urlUser != "" {
		return obj.urlUser
	}
	user := ""
	userConfig := GetConfigFromKey(ChannelUserID)
	if userConfig != nil {
		user = obj.urlProps.GetProperty(userConfig, "")
		if user == "" {
			env := NewTGEnvironment()
			user = env.GetChannelUser()
			if user == "" {
				user = env.GetChannelDefaultUser()
			}
		}
	}
	return user
}

func (obj *LinkUrl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("LinkUrl:{")
	buffer.WriteString(fmt.Sprintf("FtUrls: %+v", obj.ftUrls))
	buffer.WriteString(fmt.Sprintf(", UrlHost: %s", obj.urlHost))
	buffer.WriteString(fmt.Sprintf(", IsIPv6: %+v", obj.isIPv6))
	buffer.WriteString(fmt.Sprintf(", UrlPort: %d", obj.urlPort))
	buffer.WriteString(fmt.Sprintf(", Protocol: %+v", obj.protocol))
	buffer.WriteString(fmt.Sprintf(", UrlStr: %s", obj.urlStr))
	buffer.WriteString(fmt.Sprintf(", UrlProps: %+v", obj.urlProps))
	buffer.WriteString(fmt.Sprintf(", UrlUser: %s", obj.urlUser))
	buffer.WriteString("}")
	return buffer.String()
}


// ======= Exception channel Type =======
type ExceptionChannelType int

const (
	RethrowException ExceptionChannelType = iota
	RetryOperation
	Disconnected
)

type ExceptionHandleResult struct {
	ExceptionType    ExceptionChannelType // types.TGExceptionType
	ExceptionMessage string
}

func (proType ExceptionChannelType) ChannelException() *ExceptionHandleResult {
	// Use a buffer for efficient string concatenation
	var exceptionResult ExceptionHandleResult

	if proType&RethrowException == RethrowException {
		exceptionResult = ExceptionHandleResult{
			ExceptionType:    TGErrorGeneralException,
			ExceptionMessage: "TGDB-CHANNEL-FAIL:Failed to reconnect",
		}
	}
	if proType&RetryOperation == RetryOperation {
		exceptionResult = ExceptionHandleResult{
			ExceptionType:    TGErrorRetryIOException,
			ExceptionMessage: "TGDB-CHANNEL-RETRY:Channel Reconnected, Retry Operation",
		}
	}
	if proType&Disconnected == Disconnected {
		exceptionResult = ExceptionHandleResult{
			ExceptionType:    TGErrorChannelDisconnected,
			ExceptionMessage: "TGDB-CHANNEL-FAIL:Failed to reconnect",
		}
	}
	return &exceptionResult
}

var ConnectionsToChannel int32

type AbstractChannel struct {
	authToken         int64
	channelLinkState  tgdb.LinkState
	channelProperties *SortedProperties
	channelUrl        *LinkUrl
	clientId          string
	connectionIndex   int
	cryptographer     tgdb.TGDataCryptoGrapher
	inboxAddress      string
	needsPing         bool
	numOfConnections  int32
	lastActiveTime    time.Time
	primaryUrl        *LinkUrl
	reader            *ChannelReader
	requestId         int64
	responses         map[int64]tgdb.TGChannelResponse
	sessionId         int64
	exceptionLock     sync.Mutex    // reentrant-lock for synchronizing sending/receiving messages over the wire
	exceptionCond     *sync.Cond    // Condition for lock
	sendLock          sync.Mutex    // reentrant-lock for synchronizing sending/receiving messages over the wire
	tracer            tgdb.TGTracer // Used for tracing the information flow during the execution
	user			string
	pw				[]byte
}

func DefaultAbstractChannel() *AbstractChannel {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AbstractChannel{})

	newChannel := AbstractChannel{
		authToken:        -1,
		connectionIndex:  0,
		needsPing:        false,
		numOfConnections: 0,
		lastActiveTime:   time.Now(),
		channelLinkState: tgdb.LinkNotConnected,
		channelUrl:       DefaultLinkUrl(),
		primaryUrl:       DefaultLinkUrl(),
		responses:        make(map[int64]tgdb.TGChannelResponse, 0),
		sessionId:        -1,
	}
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	return &newChannel
}

func NewAbstractChannel(linkUrl *LinkUrl, props *SortedProperties) *AbstractChannel {
	newChannel := DefaultAbstractChannel()
	newChannel.channelUrl = linkUrl
	newChannel.primaryUrl = linkUrl
	newChannel.channelProperties = props
	// TODO: Uncomment the following two lines to test Trace functionality
	//enableTraceFlag := newChannel.channelProperties.GetProperty(utils.GetConfigFromKey(utils.EnableConnectionTrace), "true")
	//if enableTraceFlag == "true" {
	enableTraceFlag := newChannel.channelProperties.GetPropertyAsBoolean(GetConfigFromKey(EnableConnectionTrace))
	if enableTraceFlag {
		traceDir := newChannel.channelProperties.GetProperty(GetConfigFromKey(ConnectionTraceDir), ".")
		clientId := newChannel.channelProperties.GetProperty(GetConfigFromKey(ChannelClientId), "")
		newChannel.tracer = NewChannelTracer(clientId, traceDir)
		//newChannel.tracer.Start()
	}
	return newChannel
}

/////////////////////////////////////////////////////////////////
// Private functions for TGChannel / Derived Channels
/////////////////////////////////////////////////////////////////

func getChannelClientProtocolVersion() uint16 {
	return GetProtocolVersion()
}

func getServerProtocolVersion() uint16 {
	return 0
}

func isChannelClosing(obj tgdb.TGChannel) bool {
	if obj.GetLinkState() == tgdb.LinkClosing {
		return true
	}
	return false
}

func isChannelClosed(obj tgdb.TGChannel) bool {
	if obj.GetLinkState() == tgdb.LinkClosing || obj.GetLinkState() == tgdb.LinkClosed || obj.GetLinkState() == tgdb.LinkTerminated {
		return true
	}
	return false
}

func isChannelConnected(obj tgdb.TGChannel) bool {
	if obj.GetLinkState() == tgdb.LinkConnected {
		return true
	}
	return false
}

func (obj *AbstractChannel) DoAuthenticate() tgdb.TGError {
	return nil
}

func (obj *AbstractChannel) SetAuthToken (token int64) {
	obj.authToken = token
}


func (obj *AbstractChannel) channelToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("AbstractChannel:{")
	buffer.WriteString(fmt.Sprintf("AuthToken: %d", obj.authToken))
	//buffer.WriteString(fmt.Sprintf(", ChannelProperties: %+v", obj.ChannelProperties))
	buffer.WriteString(fmt.Sprintf(", ClientId: %s", obj.clientId))
	buffer.WriteString(fmt.Sprintf(", ConnectionIndex: %d", obj.connectionIndex))
	//buffer.WriteString(fmt.Sprintf(", DataCryptoGrapher: %+v", obj.cryptoGrapher))
	buffer.WriteString(fmt.Sprintf(", InboxAddress: %s", obj.inboxAddress))
	buffer.WriteString(fmt.Sprintf(", NeedsPing: %+v", obj.needsPing))
	buffer.WriteString(fmt.Sprintf(", NumOfConnections: %d", obj.numOfConnections))
	buffer.WriteString(fmt.Sprintf(", LastActiveTime: %+v", obj.lastActiveTime))
	buffer.WriteString(fmt.Sprintf(", LinkState: %s", obj.channelLinkState.String()))
	buffer.WriteString(fmt.Sprintf(", ChannelUrl: %s", obj.channelUrl.String()))
	buffer.WriteString(fmt.Sprintf(", PrimaryUrl: %s", obj.primaryUrl.String()))
	buffer.WriteString(fmt.Sprintf(", RequestId: %d", obj.requestId))
	buffer.WriteString(fmt.Sprintf(", Responses: %+v", obj.responses))
	buffer.WriteString(fmt.Sprintf(", SessionId: %d", obj.sessionId))
	//buffer.WriteString(fmt.Sprintf(", Reader: %s", obj.GetReader().String()))
	//buffer.WriteString(fmt.Sprintf(", Tracer: %s", obj.GetTracer().String()))
	buffer.WriteString(fmt.Sprintf(", ExceptionCond: %+v", obj.exceptionCond))
	buffer.WriteString("}")
	return buffer.String()
}

func (obj *AbstractChannel) GetChannelPassword() []byte {
	if len(obj.pw) > 0 {
		return obj.pw
	}
	pwd := ""
	if len(obj.channelProperties.GetAllProperties()) > 0 {
		pwd = obj.channelProperties.GetProperty(GetConfigFromKey(ChannelPassword), "")
	}
	return []byte(pwd)
}

func (obj *AbstractChannel) SetChannelPassword(pword []byte) {
	obj.pw = pword
}

func (obj *AbstractChannel) getDatabaseName() string {
	dbName := ""
	if len(obj.channelProperties.GetAllProperties()) > 0 {
		dbName = obj.channelProperties.GetProperty(GetConfigFromKey(ConnectionDatabaseName), "")
	}
	return dbName
}

func (obj *AbstractChannel) GetChannelUserName() string {
	if len(obj.user) > 0 {
		return obj.user
	}
	user := ""
	if len(obj.channelProperties.GetAllProperties()) > 0 {
		user = obj.channelProperties.GetProperty(GetConfigFromKey(ChannelUserID), "")
	}
	return user
}

func (obj *AbstractChannel) SetChannelUserName(uname string) {
	obj.user = uname
}

func (obj *AbstractChannel) isChannelPingable() bool {
	return obj.needsPing
}

func (obj *AbstractChannel) setChannelAuthToken(authToken int64) {
	obj.authToken = authToken
}

func (obj *AbstractChannel) setChannelClientId(clientId string) {
	obj.clientId = clientId
}

func (obj *AbstractChannel) setChannelInboxAddr(addr string) {
	obj.inboxAddress = addr
}

func (obj *AbstractChannel) setChannelSessionId(sessionId int64) {
	obj.sessionId = sessionId
}

// SetDataCryptoGrapher sets the data cryptographer
func (obj *AbstractChannel) setDataCryptoGrapher(crypto tgdb.TGDataCryptoGrapher) {
	obj.cryptographer = crypto
}

func (obj *AbstractChannel) setNoOfConnections(num int32) {
	obj.numOfConnections = num
}

/////////////////////////////////////////////////////////////////
// Helper (Quite Involved) functions for AbstractChannel
/////////////////////////////////////////////////////////////////

func channelConnect(obj tgdb.TGChannel) tgdb.TGError {
	//logger.Log(fmt.Sprintf("Entering AbstractChannel:channelConnect"))
	if isChannelConnected(obj) {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("AbstractChannel:channelConnect channel is already connected"))
		}
		obj.SetNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
		return nil
	}
	if isChannelClosed(obj) || obj.GetLinkState() == tgdb.LinkNotConnected {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelConnect about to channelTryRepeatConnect for object '%+v'", obj.String()))
		}
		err := channelTryRepeatConnect(obj, false)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelConnect channelTryRepeatConnect failed w/ '%+v'", err.Error()))
			return err
		}
		obj.SetChannelLinkState(tgdb.LinkConnected)
		obj.SetNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelConnect successfully established socket connection and now has '%d' number of connections", obj.GetNoOfConnections()))
		}
	} else {
		logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelConnect channelTryRepeatConnect - connect called on an invalid state := '%s'", obj.GetLinkState().String()))
		errMsg := fmt.Sprintf("Connect called on an invalid state := '%s'", obj.GetLinkState().String())
		return NewTGGeneralExceptionWithMsg(errMsg)
	}
	//logger.Log(fmt.Sprintf("Returning AbstractChannel:channelConnect having '%d' number of connections", obj.GetNoOfConnections()))
	return nil
}

func channelDisConnect(obj tgdb.TGChannel) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractChannel:channelDisConnect"))
	}
	obj.ChannelLock()
	defer obj.ChannelUnlock()

	if !isChannelConnected(obj) {
		logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelDisConnect channel is already disconnected"))
		return nil
	}

	if obj.GetNoOfConnections() == 0 {
		logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelDisConnect calling disconnect more than number of connects"))
		return nil
	}
	obj.SetNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, -1))
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelDisConnect"))
	}
	return nil
}

func channelHandleException(obj tgdb.TGChannel, ex tgdb.TGError, bReconnect bool) *ExceptionHandleResult {
	bReconnect = false
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractChannel:channelHandleException w/ Error: '%+v' and Reconnect Flag: '%+v'", ex, bReconnect))
	}
	obj.ExceptionLock()
	defer func() {
		if bReconnect {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelHandleException about to obj.exceptionCond.Broadcast()"))
			}
			obj.GetExceptionCondition().Broadcast()
		}
		obj.ExceptionUnlock()
	} ()

	if ex.GetErrorType() != TGErrorIOException {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AbstractChannel:channelHandleException w/ RethrowException"))
		}
		return RethrowException.ChannelException()
	}

	connectionOpTimeout := obj.GetProperties().GetPropertyAsInt(GetConfigFromKey(ConnectionOperationTimeoutSeconds))

	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractChannel:channelHandleException Infinite Loop"))
		}
		if bReconnect {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Returning AbstractChannel:channelHandleException Infinite Loop"))
			}
			break
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelHandleException Infinite Loop about to obj.exceptionCond.Wait()"))
		}
		//obj.GetExceptionCondition().Wait()
		time.Sleep(time.Duration(connectionOpTimeout) * time.Second)
		//obj.GetExceptionCondition().Broadcast()
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelHandleException Infinite Loop about to check isChannelConnected()"))
		}

		if isChannelConnected(obj) {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Returning AbstractChannel:channelHandleException Infinite Loop w/ RetryOperation as channel is connected"))
			}
			return RetryOperation.ChannelException()
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelHandleException Infinite Loop about to check IsClosed()"))
		}
		if obj.IsClosed() {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Returning AbstractChannel:channelHandleException Infinite Loop w/ DisconnectedException as channel is closed"))
			}
			return Disconnected.ChannelException()
		}
	} // End of Infinite Loop

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelHandleException about to obj.channelReconnect()"))
	}
	if channelReconnect(obj) {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Returning AbstractChannel:channelHandleException w/ RetryOperation as failure in channelReconnect()"))
		}
		return RetryOperation.ChannelException()
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelHandleException w/ DisconnectedException for input exception: '%+v' and Reconnect Flag: '%+v'", ex, bReconnect))
	}
	return Disconnected.ChannelException()
}

// channelProcessMessage processes a message received on the channel. This is called from the ChannelReader.
func channelProcessMessage(obj tgdb.TGChannel, msg tgdb.TGMessage) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelProcessMessage"))
	}
	reqId := msg.GetRequestId()
	channelResponseMap := obj.GetResponses()
	channelResponse := channelResponseMap[reqId]

	if channelResponse == nil {
		errMsg := fmt.Sprintf("AbstractChannel:channelProcessMessage - Received no response message for corresponding request :%d", reqId)
		logger.Error(fmt.Sprintf("ERROR: Returning %s", errMsg))
		//return exception.GetErrorByType(types.TGErrorGeneralException, types.TGDB_CHANNEL_ERROR, errMsg, "")
		return nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AbstractChannel:channelProcessMessage about to channelResponse.SetReply() w/ MSG"))
	}
	channelResponse.SetReply(msg)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelProcessMessage"))
	}
	return nil
}

func channelReconnect(obj tgdb.TGChannel) bool {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractChannel:channelReconnect"))
	}
	cn1 := GetConfigFromKey(ChannelFTHosts)
	ftHosts := obj.GetProperties().GetProperty(cn1, "")
	if len(ftHosts) <= 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning AbstractChannel:channelReconnect - There are no FT host URLs configured for this channel"))
		return false
	}

	// This is needed here to avoid a FD leak
	// Execute Derived channel's method - Ignore the Error Handling
	_ = obj.CloseSocket()

	oldUrl := obj.GetChannelURL()
	cn := GetConfigFromKey(ChannelFTRetryIntervalSeconds)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelReconnect config for ChannelFTRetryIntervalSeconds is '%+v", cn))
	}
	connectInterval := obj.GetProperties().GetPropertyAsInt(cn)
	cn = GetConfigFromKey(ChannelFTRetryCount)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelReconnect config for ChannelFTRetryCount is '%+v", cn))
	}
	retryCount := obj.GetProperties().GetPropertyAsInt(cn)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelReconnect Retrying to reconnnect %d times at interval of %d seconds to FTUrls.", retryCount, connectInterval))
	}

	obj.SetChannelLinkState(tgdb.LinkReconnecting)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AbstractChannel:channelReconnect about to obj.channelTryRepeatConnect()"))
	}
	err := channelTryRepeatConnect(obj, true)
	if err != nil {
		obj.SetChannelURL(oldUrl.(*LinkUrl))
		obj.SetChannelLinkState(tgdb.LinkClosed)
		logger.Error(fmt.Sprintf("ERROR: Returning AbstractChannel:channelReconnect - failed to reconnect w/ error: /%+v'", err.Error()))
		return false
	}
	obj.SetChannelLinkState(tgdb.LinkConnected)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractChannel:channelReconnect w/ NO Errors"))
	}
	return true
}

func channelRequestReply(obj tgdb.TGChannel, request tgdb.TGMessage) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelRequestReply"))
	}
	var respMessage tgdb.TGMessage

	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractChannel:channelRequestReply Infinite Loop"))
		}
		resp, err := func() (tgdb.TGMessage, tgdb.TGError) {
			obj.ChannelLock()
			defer obj.ChannelUnlock()

			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelRequestReply Infinite Loop about to obj.Send()"))
			}
			// Execute Derived channel's method
			err := obj.Send(request)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelRequestReply obj.Send failed w/ '%+v'", err.Error()))
				ehResult := channelHandleException(obj, err, true)
				if ehResult.ExceptionType == RethrowException {
					logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelRequestReply - Failed to send message"))
					if err.GetErrorType() == TGErrorGeneralException {
						return nil, NewTGGeneralExceptionWithMsg(err.Error())
					}
					errMsg := fmt.Sprintf("AbstractChannel:channelRequestReply - %s w/ error: %s", TGDB_SEND_ERROR, err.Error())
					return nil, BuildException(tgdb.TGTransactionStatus(err.GetErrorType()), errMsg)
				} else if ehResult.ExceptionType == Disconnected {
					logger.Error(fmt.Sprint("Returning AbstractChannel:channelRequestReply - channel got disconnected"))
					return nil, NewTGChannelDisconnected(err.GetErrorCode(), err.GetErrorType(), err.GetErrorMsg(), err.GetErrorDetails())
				} else {
					// TODO: Revisit later - Should we not throw an error?
					logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelRequestReply in obj.Send retrying to send message on url: '%s'", obj.GetChannelURL().GetUrlAsString()))
					//continue
					return nil, nil
				}
			}

			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelRequestReply Infinite Loop about to obj.ReadWireMsg()"))
			}
			// Execute Derived channel's method
			msg, err := obj.ReadWireMsg()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelRequestReply obj.ReadWireMsg failed w/ '%+v'", err.Error()))
				ehResult := channelHandleException(obj, err, true)
				if ehResult.ExceptionType == RethrowException {
					logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelRequestReply - Failed to read message"))
					if err.GetErrorType() == TGErrorGeneralException {
						return nil, NewTGGeneralExceptionWithMsg(err.Error())
					}
					errMsg := fmt.Sprintf("AbstractChannel:channelRequestReply - %s w/ error: %s", TGDB_SEND_ERROR, err.Error())
					return nil, BuildException(tgdb.TGTransactionStatus(err.GetErrorType()), errMsg)
				} else if ehResult.ExceptionType == Disconnected {
					logger.Error(fmt.Sprint("Returning AbstractChannel:channelRequestReply - channel got disconnected"))
					return nil, NewTGChannelDisconnected(err.GetErrorCode(), err.GetErrorType(), err.GetErrorMsg(), err.GetErrorDetails())
				} else {
					// TODO: Revisit later - Should we not throw an error?
					logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelRequestReply in obj.ReadWireMsg retrying to send message on url: '%s'", obj.GetChannelURL().GetUrlAsString()))
					//continue
					return nil, nil
				}
			}

			//obj.ChannelUnlock()
			return msg, nil
		} ()
		if resp == nil && err == nil {
			continue
		} else if err != nil {
			if err.GetErrorType() == TGSuccess {
				respMessage = nil
				break
			}
			return nil, err
		} else {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelRequestReply Breaking Loop successfully w/ msgResponse: '%+v'", resp))
			}
			respMessage = resp
			break
		}
	} // End of Infinite Loop

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelRequestReply w/ %+v", respMessage))
	}
	return respMessage, nil
}

func channelSendMessage(obj tgdb.TGChannel, msg tgdb.TGMessage, resendFlag bool) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractChannel:channelSendMessage w/ Message type: '%+v'", msg.GetVerbId()))
	}
	var error tgdb.TGError
	var resendMode tgdb.ResendMode
	if resendFlag {
		resendMode = tgdb.ModeReconnectAndResend
	} else {
		resendMode = tgdb.ModeReconnectAndRaiseException
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelSendMessage using '%s'", resendMode.String()))
	}

	var count int
	count = 0

	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Entering AbstractChannel:channelSendMessage Infinite Loop"))
		}
		contFlag, err := func() (bool, tgdb.TGError) {
			obj.ChannelLock()
			defer obj.ChannelUnlock()

			if !isChannelConnected(obj) {
				logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelSendMessage - channel is closed"))
				errMsg := fmt.Sprint("AbstractChannel:channelSendMessage - channel is closed")
				return false, GetErrorByType(TGErrorGeneralException, TGDB_CHANNEL_ERROR, errMsg, "")
			}
			obj.ChannelLock()

			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelSendMessage Infinite Loop about to obj.Send()"))
			}
			// Execute Derived channel's message communication mechanism
			err := obj.Send(msg)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelSendMessage obj.Send failed w/ '%+v'", err.Error()))
				ehResult := channelHandleException(obj, err, false)
				if ehResult.ExceptionType == RethrowException {
					logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelSendMessage - Failed to send message"))
					if err.GetErrorType() == TGErrorGeneralException {
						return false, NewTGGeneralExceptionWithMsg(err.Error())
					}
					errMsg := fmt.Sprintf("AbstractChannel:channelSendMessage - %s w/ error: %s", TGDB_SEND_ERROR, err.Error())
					return false, BuildException(tgdb.TGTransactionStatus(err.GetErrorType()), errMsg)
				} else if ehResult.ExceptionType == Disconnected {
					logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelSendMessage - channel got disconnected"))
					return false, NewTGChannelDisconnected(err.GetErrorCode(), err.GetErrorType(), err.GetErrorMsg(), err.GetErrorDetails())
				} else {
					// TODO: Revisit later - Should we not throw an error?
					logger.Warning(fmt.Sprintf("WARNING: AbstractChannel:channelSendMessage Retrying to send message on url: '%s'", obj.GetChannelURL().GetUrlAsString()))
					//continue
					//return true, nil
					count++
					if count > 5 {
						return false, nil
					}

				}
			}
			return false, nil
		} ()
		if contFlag {
			continue
		} else {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelSendMessage Breaking Loop"))
			}
			error = err
			break
		}
	} // End of Infinite Loop
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelSendMessage"))
	}
	return error
}

func channelSendRequest(obj tgdb.TGChannel, msg tgdb.TGMessage, channelResponse tgdb.TGChannelResponse, resendFlag bool) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering AbstractChannel:channelSendRequest w/ Message type: '%+v' ChannelResponse: '%+v'", msg.GetVerbId(), channelResponse))
	}
	reqId := channelResponse.GetRequestId()
	msg.SetRequestId(reqId)

	var respMessage tgdb.TGMessage
	var resendMode tgdb.ResendMode
	if resendFlag {
		resendMode = tgdb.ModeReconnectAndResend
	} else {
		resendMode = tgdb.ModeReconnectAndRaiseException
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelSendRequest using '%s'", resendMode.String()))
	}

	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelSendRequest Infinite Loop"))
		}
		resp, err := func() (tgdb.TGMessage, tgdb.TGError)  {
			obj.ChannelLock()
			defer obj.ChannelUnlock()

			if !isChannelConnected(obj) {
				errMsg := fmt.Sprintf("AbstractChannel:channelSendRequest - channel is closed")
				logger.Error(fmt.Sprintf("ERROR: Returning %s", errMsg))
				return nil, GetErrorByType(TGErrorGeneralException, TGDB_CHANNEL_ERROR, errMsg, "")
			}
			// TODO: Uncomment once Trace functionality is implemented and tested
			//if obj.GetTracer() != nil {
			//	obj.GetTracer().Trace(msg)
			//}
			//obj.ChannelLock()
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelSendRequest about to set channel response '%+v' in map '%+v'", channelResponse, obj.GetResponses()))
			}
			obj.SetResponse(reqId, channelResponse)

			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelSendRequest Infinite Loop about to obj.Send()"))
			}
			// Execute Derived channel's message communication mechanism
			err := obj.Send(msg)
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelSendRequest after obj.Send()"))
			}
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: AbstractChannel:channelSendRequest obj.Send failed w/ '%+v'", err.Error()))
				//obj.ChannelUnlock()
				ehResult := channelHandleException(obj, err, false)
				if ehResult.ExceptionType == RethrowException {
					logger.Error(fmt.Sprint("ERROR: Returning AbstractChannel:channelSendRequest - Failed to send message"))
					if err.GetErrorType() == TGErrorGeneralException {
						return nil, NewTGGeneralExceptionWithMsg(err.Error())
					}
					errMsg := fmt.Sprintf("AbstractChannel:channelSendRequest - %s w/ error: %s", TGDB_SEND_ERROR, err.Error())
					return nil, BuildException(tgdb.TGTransactionStatus(err.GetErrorType()), errMsg)
				} else if ehResult.ExceptionType == Disconnected {
					logger.Error(fmt.Sprint("Returning AbstractChannel:channelSendRequest - channel got disconnected"))
					return nil, NewTGChannelDisconnected(err.GetErrorCode(), err.GetErrorType(), err.GetErrorMsg(), err.GetErrorDetails())
				} else {
					// TODO: Revisit later - Should we not throw an error?
					logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelSendRequest Infinite Loop retrying to send message on url: '%s'", obj.GetChannelURL().GetUrlAsString()))
					//continue
					return nil, nil
				}
			}
			if !channelResponse.IsBlocking() {
				//obj.ChannelUnlock()
				logger.Warning(fmt.Sprint("WARNING: Returning AbstractChannel:channelSendRequest as channel response is NOT blocking"))
				//return nil, nil
				return nil, NewTGSuccessWithMsg("WARNING: Returning AbstractChannel:channelSendRequest as channel response is NOT blocking")
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside AbstractChannel:channelSendRequest Infinite Loop about to channelResponse.Await()"))
			}
			channelResponse.Await(channelResponse.(*BlockingChannelResponse))
			delete(obj.GetResponses(), reqId)
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelSendRequest Infinite Loop about to channelResponse.GetReply()"))
			}
			msgResponse := channelResponse.GetReply()

			if msgResponse != nil && msgResponse.GetVerbId() == VerbExceptionMessage {
				//obj.ChannelUnlock()
				exMsg := msgResponse.(*ExceptionMessage)
				if exMsg.GetExceptionType() == TGErrorRetryIOException {
					//continue
					return nil, nil
				}
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractChannel:channelSendRequest Breaking Loop for VerbExceptionMessage w/ msgRespbnse: '%+v'", msgResponse.String()))
				return nil, NewTGGeneralExceptionWithMsg(exMsg.GetExceptionMsg())
			}
			//obj.ChannelUnlock()
			return msgResponse, nil
		} ()
		if resp == nil && err == nil {
			continue
		} else if err != nil {
			if err.GetErrorType() == TGSuccess {
				respMessage = nil
				break
			}
			return nil, err
		} else {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelSendRequest Breaking Loop successfully w/ msgResponse: '%+v'", resp))
			}
			respMessage = resp
			break
		}
	} // End of Infinite Loop
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelSendRequest w/ %+v", respMessage))
	}
	return respMessage, nil
}

func channelStart(obj tgdb.TGChannel) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelStart"))
	}
	if !isChannelConnected(obj) {
		errMsg := fmt.Sprint("AbstractChannel:channelStart - channel is not connected")
		logger.Error(fmt.Sprintf("ERROR: Returning %s", errMsg))
		return GetErrorByType(TGErrorGeneralException, TGDB_CHANNEL_ERROR, errMsg, "")
	}
	obj.EnablePing()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStart about to start channel Reader"))
	}
	go obj.GetReader().Start()
	// TODO: Uncomment once Trace functionality is implemented and tested
	//logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStart about to start channel Tracer"))
	//go obj.GetTracer().Start()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractChannel:channelStart"))
	}
	return nil
}

func channelStop(obj tgdb.TGChannel, bForcefully bool) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelStop"))
	}
	obj.ChannelLock()
	defer func() {
		if isChannelClosing(obj) {
			obj.SetChannelLinkState(tgdb.LinkClosed)
		}
		// Execute Derived channel's method - Ignore Error Handling
		obj.ChannelUnlock()
	} ()

	if !isChannelConnected(obj) {
		logger.Warning(fmt.Sprint("WARNING: Returning AbstractChannel:channelStop as channel is already disconnected"))
		return
	}

	if bForcefully || obj.GetNoOfConnections() == 0 {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStop stopping channel"))
		}
		obj.DisablePing()
		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStop about to stop channel Reader"))
		}
		obj.GetReader().Stop()
		// TODO: Uncomment once Trace functionality is implemented and tested
		//logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStop about to stop channel Tracer"))
		//obj.GetTracer().Stop()

		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStop about to CreateMessageForVerb()"))
		}
		// Send the disconnect request. sendRequest will not receive a channel response since the channel will be disconnected.
		msgRequest, err := CreateMessageForVerb(VerbDisconnectChannelRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Inside AbstractChannel:channelStop VerbDisconnectChannelRequest CreateMessageForVerb failed with '%s'", err.Error()))
			// Execute Derived channel's method - Ignore Error Handling
			_ = obj.CloseSocket()
			return
		}
		// Execute Derived channel's method
		err = obj.Send(msgRequest)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Inside AbstractChannel:channelStop VerbDisconnectChannelRequest send failed with '%s'", err.Error()))
			// Execute Derived channel's method - Ignore Error Handling
			_ = obj.CloseSocket()
			return
		}
		obj.SetChannelLinkState(tgdb.LinkClosing)

		if logger.IsDebug() {
			logger.Debug(fmt.Sprint("Inside AbstractChannel:channelStop about to CloseSocket()"))
		}
		// Execute Derived channel's method - Ignore Error Handling
		_ = obj.CloseSocket()
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractChannel:channelStop"))
	}
	return
}

// channelTerminated closes the socket channel. This is called from the ChannelReader.
func channelTerminated(obj tgdb.TGChannel, killMsg string) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelTerminated"))
	}
	obj.ExceptionLock()
	defer obj.ExceptionUnlock()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTerminated about to terminate session/channel with '%s'", killMsg))
	}

	obj.SetChannelLinkState(tgdb.LinkTerminated)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AbstractChannel:channelTerminated about to CloseSocket()"))
	}
	// Execute Derived channel's method - Ignore Error Handling
	_ = obj.CloseSocket()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning AbstractChannel:channelTerminated w/ '%s'", killMsg))
	}
	return
}

func channelTryRepeatConnect(obj tgdb.TGChannel, sleepOnFirstInvocation bool) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AbstractChannel:channelTryRepeatConnect"))
	}
	cn := GetConfigFromKey(ChannelFTRetryIntervalSeconds)
	//logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect config for ChannelFTRetryIntervalSeconds is '%+v", cn))
	connectInterval := obj.GetProperties().GetPropertyAsInt(cn)
	cn = GetConfigFromKey(ChannelFTRetryCount)
	//logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect config for ChannelFTRetryCount is '%+v", cn))
	retryCount := obj.GetProperties().GetPropertyAsInt(cn)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect Trying to connnect %d times at interval of %d seconds to FTUrls", retryCount, connectInterval))
	}

	reconnected := false
	ftUrls := obj.GetPrimaryURL().GetFTUrls()
	urlCount := len(ftUrls)
	index := obj.GetConnectionIndex()
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect current object's primary url '%s' has FTUrls as '%+v'", obj.GetPrimaryURL().GetUrlAsString(), ftUrls))
	}

	for {
		if urlCount > 0 {
			url := ftUrls[index]
			obj.SetChannelURL(url.(*LinkUrl))
		}
		// From here onwards, object's primary Attributes will be used such as PrimaryUrl, LinkUrl etc.
		urlStr := obj.GetPrimaryURL().GetUrlAsString()
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect Infinite Loop to create a socket for URL: '%s'", urlStr))
		}

		for i := 0; i < retryCount; i++ {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect Attempt:%d to connect to URL:%s", i, urlStr))
			}
			if sleepOnFirstInvocation {
				time.Sleep(time.Duration(connectInterval) * time.Second)
				sleepOnFirstInvocation = false
			}

			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect about to CreateSocket() on attempt:%d to URL:%s", i, urlStr))
			}
			// Execute Derived channel's method
			err := obj.CreateSocket()
			if err != nil {
				logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelTryRepeatConnect about to CloseSocket() on attempt:%d to URL:%s w/ '%+v'", i, urlStr, err.Error()))
				// Execute Derived channel's method - Ignore Error Handling
				_ = obj.CloseSocket()
				continue
			}

			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect about to OnConnect() on attempt:%d to URL:%s", i, urlStr))
			}
			// Execute Derived channel's method
			err = obj.OnConnect()
			if err != nil {
				logger.Warning(fmt.Sprintf("WARNING: Inside AbstractChannel:channelTryRepeatConnect Failed to execute channel specific OnConnect w/ '%+v'", err.Error()))
				// Execute Derived channel's method - Ignore Error Handling
				_ = obj.CloseSocket()
				continue
			}
			obj.SetConnectionIndex(index)
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect successfully created socket and executed OnConnect() on attempt:%d to URL:%s", i, urlStr))
			}
			reconnected = true
			break
		} // End of for loop for Retry Attempts

		if urlCount > 0 {
			index = (index + 1) % urlCount
		} else {
			index += 1
		}

		if index != obj.GetConnectionIndex() || reconnected {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AbstractChannel:channelTryRepeatConnect breaking from Infinite Loop"))
			}
			break
		}
	} // End of Outer Infinite For loop

	if !reconnected {
		errMsg := fmt.Sprintf("AbstractChannel:channelTryRepeatConnect %s - failed %d attempts to connect to TGDB Server.", "TGDB-CONNECT-ERR", retryCount)
		logger.Error(fmt.Sprintf("ERROR: Returning '%s'", errMsg))
		return NewTGConnectionTimeoutWithMsg(errMsg)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AbstractChannel:channelTryRepeatConnect w/ NO error after successfully creating socket and executing OnConnect()"))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// ChannelLock locks the communication channel between TGDB client and server
func (obj *AbstractChannel) ChannelLock() {
	obj.sendLock.Lock()
}

// ChannelUnlock unlocks the communication channel between TGDB client and server
func (obj *AbstractChannel) ChannelUnlock() {
	obj.sendLock.Unlock()
}

// Connect connects the underlying channel using the URL end point
func (obj *AbstractChannel) Connect() tgdb.TGError {
	return channelConnect(obj)
}

// DisablePing disables the pinging ability to the channel
func (obj *AbstractChannel) DisablePing() {
	obj.needsPing = false
}

// Disconnect disconnects the channel from its URL end point
func (obj *AbstractChannel) Disconnect() tgdb.TGError {
	return channelDisConnect(obj)
}

// EnablePing enables the pinging ability to the channel
func (obj *AbstractChannel) EnablePing() {
	obj.needsPing = true
}

// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
func (obj *AbstractChannel) ExceptionLock() {
	obj.exceptionLock.Lock()
}

// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
func (obj *AbstractChannel) ExceptionUnlock() {
	obj.exceptionLock.Unlock()
}

// GetAuthToken gets Authorization Token
func (obj *AbstractChannel) GetAuthToken() int64 {
	return obj.authToken
}

// GetClientId gets Client Name
func (obj *AbstractChannel) GetClientId() string {
	return obj.clientId
}

// GetChannelURL gets the channel URL
func (obj *AbstractChannel) GetChannelURL() tgdb.TGChannelUrl {
	return obj.channelUrl
}

// GetConnectionIndex gets the Connection Index
func (obj *AbstractChannel) GetConnectionIndex() int {
	return obj.connectionIndex
}

// GetDataCryptoGrapher gets the data cryptographer handle
func (obj *AbstractChannel) GetDataCryptoGrapher() tgdb.TGDataCryptoGrapher {
	return obj.cryptographer
}

// GetExceptionCondition gets the Exception Condition
func (obj *AbstractChannel) GetExceptionCondition() *sync.Cond {
	return obj.exceptionCond
}

// GetLinkState gets the Link/channel State
func (obj *AbstractChannel) GetLinkState() tgdb.LinkState {
	return obj.channelLinkState
}

// GetNoOfConnections gets number of connections this channel has
func (obj *AbstractChannel) GetNoOfConnections() int32 {
	return obj.numOfConnections
}

// GetPrimaryURL gets the Primary URL
func (obj *AbstractChannel) GetPrimaryURL() tgdb.TGChannelUrl {
	return obj.primaryUrl
}

// GetProperties gets the channel Properties
func (obj *AbstractChannel) GetProperties() tgdb.TGProperties {
	return obj.channelProperties
}

// GetReader gets the channel Reader
func (obj *AbstractChannel) GetReader() tgdb.TGChannelReader {
	return obj.reader
}

// GetResponses gets the channel Response Map
func (obj *AbstractChannel) GetResponses() map[int64]tgdb.TGChannelResponse {
	return obj.responses
}

// GetSessionId gets Session id
func (obj *AbstractChannel) GetSessionId() int64 {
	return obj.sessionId
}

// GetTracer gets the channel Tracer
func (obj *AbstractChannel) GetTracer() tgdb.TGTracer {
	return obj.tracer
}

// IsChannelPingable checks whether the channel is pingable or not
func (obj *AbstractChannel) IsChannelPingable() bool {
	return obj.needsPing
}

// IsClosed checks whether channel is open or closed
func (obj *AbstractChannel) IsClosed() bool {
	return isChannelClosed(obj)
}

// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
func (obj *AbstractChannel) SendMessage(msg tgdb.TGMessage) tgdb.TGError {
	return channelSendMessage(obj, msg, true)
}

// SendRequest sends a Message, waits for a response in the message format, and blocks the thread till it gets the response
func (obj *AbstractChannel) SendRequest(msg tgdb.TGMessage, response tgdb.TGChannelResponse) (tgdb.TGMessage, tgdb.TGError) {
	return channelSendRequest(obj, msg, response, true)
}

// SetChannelLinkState sets the Link/channel State
func (obj *AbstractChannel) SetChannelLinkState(state tgdb.LinkState) {
	obj.channelLinkState = state
}

// SetChannelURL sets the channel URL
func (obj *AbstractChannel) SetChannelURL(url tgdb.TGChannelUrl) {
	obj.channelUrl = url.(*LinkUrl)
}

// SetConnectionIndex sets the connection index
func (obj *AbstractChannel) SetConnectionIndex(index int) {
	obj.connectionIndex = index
}

// SetNoOfConnections sets number of connections
func (obj *AbstractChannel) SetNoOfConnections(count int32) {
	obj.numOfConnections = count
}

// SetResponse sets the ChannelResponse Map
func (obj *AbstractChannel) SetResponse(reqId int64, response tgdb.TGChannelResponse) {
	obj.responses[reqId] = response
}

// Start starts the channel so that it can send and receive messages
func (obj *AbstractChannel) Start() tgdb.TGError {
	return channelStart(obj)
}

// Stop stops the channel forcefully or gracefully
func (obj *AbstractChannel) Stop(bForcefully bool) {
	channelStop(obj, bForcefully)
}

// CreateSocket creates a network socket to transfer the messages in the byte format
func (obj *AbstractChannel) CreateSocket() tgdb.TGError {
	logger.Error(fmt.Sprintf("####### ======> ERROR: Entering AbstractChannel:CreateSocket"))
	// No-op for Now! This needs to be implemented by derived channels (TCP/SSL/HTTP)
	return nil
}

// CloseSocket closes the network socket
func (obj *AbstractChannel) CloseSocket() tgdb.TGError {
	logger.Error(fmt.Sprintf("####### ======> ERROR: Entering AbstractChannel:CloseSocket"))
	// No-op for Now! This needs to be implemented by derived channels (TCP/SSL/HTTP)
	return nil
}

// OnConnect executes functional logic after successfully establishing the connection to the server
func (obj *AbstractChannel) OnConnect() tgdb.TGError {
	logger.Error(fmt.Sprintf("####### ======> ERROR: Entering AbstractChannel:OnConnect"))
	// No-op for Now! This needs to be implemented by derived channels (TCP/SSL/HTTP)
	return nil
}

// ReadWireMsg read the message from the network in the byte format
func (obj *AbstractChannel) ReadWireMsg() (tgdb.TGMessage, tgdb.TGError) {
	logger.Error(fmt.Sprintf("####### ======> ERROR: Entering AbstractChannel:ReadWireMsg"))
	// No-op for Now! This needs to be implemented by derived channels (TCP/SSL/HTTP)
	return nil, nil
}

// Send Message to the server, compress and or encrypt.
// Hence it is abstraction, that the channel knows about it.
// @param msg       The message that needs to be sent to the server
func (obj *AbstractChannel) Send(msg tgdb.TGMessage) tgdb.TGError {
	logger.Error(fmt.Sprintf("####### ======> ERROR: Entering AbstractChannel:Send w/ Message as '%s'", msg.String()))
	// No-op for Now! This needs to be implemented by derived channels (TCP/SSL/HTTP)
	return nil
}

func (obj *AbstractChannel) String() string {
	return obj.channelToString()
}

type ChannelTracer struct {
	msgQueue  *SimpleQueue
	msgTracer *ChannelMessageTracer
	clientId  string
	isRunning bool
}

func DefaultChannelTracer() *ChannelTracer {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelTracer{})

	newChannelTracer := ChannelTracer{
		msgQueue:  NewSimpleQueue(),
		clientId:  "",
		isRunning: false,
	}

	return &newChannelTracer
}

func NewChannelTracer(client, traceDir string) *ChannelTracer {
	newChannelTracer := DefaultChannelTracer()
	newChannelTracer.clientId = client
	msgTracer := NewChannelMessageTracer(newChannelTracer.msgQueue, client, traceDir)
	newChannelTracer.msgTracer = msgTracer
	return newChannelTracer
}

/////////////////////////////////////////////////////////////////
// Private functions for ChannelTracer
/////////////////////////////////////////////////////////////////

func (obj *ChannelTracer) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelTracer:{")
	buffer.WriteString(fmt.Sprintf("ClientId: %+v", obj.clientId))
	buffer.WriteString(fmt.Sprintf(", MsgQueue: %d", obj.msgQueue))
	buffer.WriteString(fmt.Sprintf(", MsgTracer: %d", obj.msgTracer))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for TGTracer
/////////////////////////////////////////////////////////////////

// Start starts the channel tracer
func (obj *ChannelTracer) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelTracer:Start ..."))
	if !obj.isRunning {
		obj.msgTracer.Start()
		obj.isRunning = true
	}
	//logger.Log(fmt.Sprint("Returning ChannelTracer:Start ..."))
}

// Stop stops the channel tracer
func (obj *ChannelTracer) Stop() {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering ChannelTracer:Stop ..."))
	}
	if obj.isRunning {
		// Finish / Flush any remaining processing
		obj.msgTracer.Stop()
		obj.isRunning = false
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning ChannelTracer:Stop ..."))
	}
}

// Trace traces the path the message has taken
func (obj *ChannelTracer) Trace(msg tgdb.TGMessage) {
	//logger.Log(fmt.Sprint("Entering ChannelTracer:Trace"))
	obj.msgQueue.Enqueue(msg)
	//logger.Log(fmt.Sprintf("Returning ChannelTracer:Trace ..."))
}


type ChannelMessageTracer struct {
	currentSuffix int
	traceFile     *os.File
	isRunning     bool
	msgQueue      *SimpleQueue
	traceFileName string
}

const MaxFileSize int64 = 1 << 20

func DefaultChannelMessageTracer() *ChannelMessageTracer {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelMessageTracer{})

	newChannelMessageTracer := ChannelMessageTracer{
		currentSuffix: 0,
		isRunning:     false,
		msgQueue:      NewSimpleQueue(),
		traceFileName: "",
	}

	return &newChannelMessageTracer
}

func NewChannelMessageTracer(queue *SimpleQueue, client, traceDir string) *ChannelMessageTracer {
	newChannelMessageTracer := DefaultChannelMessageTracer()
	newChannelMessageTracer.msgQueue = queue
	newChannelMessageTracer.traceFileName = filepath.FromSlash(fmt.Sprint(traceDir, "/", client, ".trace"))
	newChannelMessageTracer.createTraceFile(0, newChannelMessageTracer.currentSuffix)
	return newChannelMessageTracer
}

/////////////////////////////////////////////////////////////////
// Private functions for ChannelMessageTracer
/////////////////////////////////////////////////////////////////

// exists returns whether the given file or directory exists, and whether it is a directory or not
func exists(path string) (bool, bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, false, err
	}
	if os.IsNotExist(err) {
		return false, false, nil
	}
	if info.IsDir() {
		return true, true, nil
	}
	return true, false, err
}

// isFileReadyForRollover checks if the file needs to be rolled over with incremented suffix
func (obj *ChannelMessageTracer) isFileReadyForRollover(suffix, msgBufLen int) bool {
	if obj == nil {
		return false
	}

	fileWithSuffix := fmt.Sprintf("%s%d", obj.traceFileName, suffix)
	fileInfo, err := os.Stat(fileWithSuffix)
	if err != nil {
		return false
	}

	if fileInfo.Size()+int64(msgBufLen) < MaxFileSize {
		return false
	}
	return true
}

// createTraceFile creates a new trace file with new suffix
func (obj *ChannelMessageTracer) createTraceFile(oldSuffix, newSuffix int) bool {
	if obj == nil {
		return false
	}

	// First check the existence of file with old suffix and close it for writing
	traceFileWithOldSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, oldSuffix)
	oFlag, _, err := exists(traceFileWithOldSuffix)
	if err != nil {
		return false
	}

	if !oFlag {
		return false
	} else {
		_ = obj.traceFile.Sync()  // Flush
		_ = obj.traceFile.Close() // Close FD
		obj.traceFile = nil
	}

	// Next check the existence of file with new suffix and open it for writing
	traceFileWithNewSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, newSuffix)
	nFlag, _, err := exists(traceFileWithNewSuffix)
	if err != nil {
		return false
	}

	if nFlag {
		return false
	}

	fp, err := os.OpenFile(traceFileWithNewSuffix, syscall.O_RDWR, 0644)
	if err != nil {
		return false
	}
	obj.traceFile = fp
	obj.currentSuffix = newSuffix
	return true
}

// extractAndTraceMessage reads a message from the message queue and traces it in the trace file
func (obj *ChannelMessageTracer) extractAndTraceMessage() {
	if obj == nil {
		return
	}

	//logger.Log(fmt.Sprintf("Entering ChannelMessageTracer:extractAndTraceMessage w/ Message Tracer object: '%+v'", obj.String()))
	for {
		if !obj.isRunning {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Breaking ChannelMessageTracer:extractAndTraceMessage loop since message tracer is not running '%+v'", obj.isRunning))
			}
			break
		}

		// At this point, the trace file with suffix is expected to be ready for writing contents in it
		msg := obj.msgQueue.Dequeue()
		if msg != nil {
			msgBuf, msgLen, err := msg.(tgdb.TGMessage).ToBytes()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in msg.ToBytes() w/ '%+v'", err.Error()))
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			// Check if the file is about to exceed the limit after appending current message buffer
			if obj.isFileReadyForRollover(obj.currentSuffix, msgLen) {
				// Create a new file with incremented suffix
				cFlag := obj.createTraceFile(obj.currentSuffix, obj.currentSuffix+1)
				if !cFlag {
					logger.Error(fmt.Sprint("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in obj.createTraceFile()"))
					time.Sleep(1000 * time.Millisecond)
					continue
				}
			}

			_, err1 := obj.traceFile.Write(msgBuf)
			if err1 != nil {
				logger.Error(fmt.Sprintf("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in obj.traceFile.Write() w/ '%+v'", err1.Error()))
				time.Sleep(1000 * time.Millisecond)
				continue
			}
		} else {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprint("Inside ChannelMessageTracer:extractAndTraceMessage - No pending messages in the queue"))
			}
			time.Sleep(1000 * time.Millisecond)
			continue
		}
	} // End of Infinite Loop
	obj.isRunning = false
	//logger.Log(fmt.Sprintf("Returning ChannelMessageTracer:extractAndTraceMessage w/ Message Tracer object: '%+v'", obj.String()))
}

func (obj *ChannelMessageTracer) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelMessageTracer:{")
	buffer.WriteString(fmt.Sprintf("TraceFileName: %+v", obj.traceFileName))
	buffer.WriteString(fmt.Sprintf(", MsgQueue: %d", obj.msgQueue))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for TGTracer
/////////////////////////////////////////////////////////////////

// Start starts the channel message tracer
func (obj *ChannelMessageTracer) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelMessageTracer:Start ..."))
	if !obj.isRunning {
		obj.isRunning = true
		// Start reading and processing messages from the wire
		obj.extractAndTraceMessage()
	}
	//logger.Log(fmt.Sprint("Returning ChannelMessageTracer:Start ..."))
}

// Stop stops the channel message tracer
func (obj *ChannelMessageTracer) Stop() {
	//logger.Log(fmt.Sprint("Entering ChannelMessageTracer:Stop ..."))
	if obj.isRunning {
		// Finish / Flush any remaining processing
		traceFileWithSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, obj.currentSuffix)
		flag, _, _ := exists(traceFileWithSuffix)
		if flag {
			_ = obj.traceFile.Sync()  // Flush
			_ = obj.traceFile.Close() // Close FD
			obj.traceFile = nil
		}
		obj.isRunning = false
	}
	//logger.Log(fmt.Sprint("Returning ChannelMessageTracer:Stop ..."))
}


type BlockingChannelResponse struct {
	status    tgdb.ChannelResponseStatus
	requestId int64
	timeout   int64
	reply     tgdb.TGMessage
	lock      sync.Mutex // reentrant-lock for synchronizing sending/receiving messages over the wire
	cond      *sync.Cond // Condition for lock
}

func DefaultBlockingChannelResponse(reqId int64) *BlockingChannelResponse {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BlockingChannelResponse{})

	newBlockingChannelResponse := BlockingChannelResponse{
		requestId: reqId,
		status:    tgdb.Waiting,
		timeout:   -1,
	}
	newBlockingChannelResponse.cond = sync.NewCond(&newBlockingChannelResponse.lock) // Condition for lock

	return &newBlockingChannelResponse
}

func NewBlockingChannelResponse(reqId, rTimeout int64) *BlockingChannelResponse {
	newBlockingChannelResponse := DefaultBlockingChannelResponse(reqId)
	newBlockingChannelResponse.timeout = rTimeout
	return newBlockingChannelResponse
}

/////////////////////////////////////////////////////////////////
// Private functions for BlockingChannelResponse
/////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////
// Implement functions for TGChannelResponse
/////////////////////////////////////////////////////////////////

// Await waits (loops) till the channel response receives reply message from the server
func (obj *BlockingChannelResponse) Await(tester tgdb.StatusTester) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering BlockingChannelResponse:Await - %d", obj.status))
	}
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	//defer func() {
	//	obj.cond.Signal()
	//	obj.lock.Unlock()
	//} ()

	//go func() {
	count := 0
	for {
		//obj.cond.Wait()
		time.Sleep(time.Duration(obj.timeout) * time.Millisecond)
		//obj.cond.Signal()
		// Terminating Condition for this Infinite Loop is:
		// 	(a) Break if the channel response object status is NOT WAITING - Status is set via SetReply()/Signal() execution
		//logger.Log(fmt.Sprintf("Inside BlockingChannelResponse:Await Loop) - abou to Test obj.status - %d", obj.status))
		if !tester.Test(obj.status) {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Breaking out from BlockingChannelResponse:Await w/ contents as '%+v'", obj.String()))
			}
			break
		}
		// TODO: Remove this block once testing is over
		count++
		if (count%10000) == 0 {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside BlockingChannelResponse:Await(%d) ... BlockingChannelResponse - %d", count, obj.status))
			}
		}
	}
	//}()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning BlockingChannelResponse:Await ..."))
	}
}

// GetCallback gets a Callback object
func (obj *BlockingChannelResponse) GetCallback() tgdb.Callback {
	// Not applicable / available for BlockingChannelResponse
	return nil
}

// GetReply gets Reply object
func (obj *BlockingChannelResponse) GetReply() tgdb.TGMessage {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	return obj.reply
}

// GetRequestId gets Request id
func (obj *BlockingChannelResponse) GetRequestId() int64 {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	return obj.requestId
}

// GetStatus gets Status
func (obj *BlockingChannelResponse) GetStatus() tgdb.ChannelResponseStatus {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	return obj.status
}

// IsBlocking checks whether this channel response is blocking or not
func (obj *BlockingChannelResponse) IsBlocking() bool {
	return true
}

// Reset resets the state of channel response and initializes everything
func (obj *BlockingChannelResponse) Reset() {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:Reset ..."))
	obj.status = tgdb.Waiting
	obj.reply = nil
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:Reset ..."))
}

// SetReply sets the reply message received from the server
func (obj *BlockingChannelResponse) SetReply(msg tgdb.TGMessage) {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:SetReply ..."))
	obj.reply = msg
	obj.status = tgdb.Ok
	obj.cond.Broadcast()
	//logger.Log(fmt.Sprintf("Returning BlockingChannelResponse:SetReply %d", obj.status))
}

// SetRequestId sets Request id
func (obj *BlockingChannelResponse) SetRequestId(reqId int64) {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:SetRequestId ..."))
	obj.requestId = reqId
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:SetRequestId ..."))
}

// Signal lets other listeners of channel response know the status of this channel response
func (obj *BlockingChannelResponse) Signal(cStatus tgdb.ChannelResponseStatus) {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:Signal ..."))
	obj.status = cStatus
	obj.cond.Broadcast()
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:Signal ..."))
}

func (obj *BlockingChannelResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BlockingChannelResponse:{")
	buffer.WriteString(fmt.Sprintf("Status: %d", obj.status))
	buffer.WriteString(fmt.Sprintf(", RequestId: %d", obj.requestId))
	buffer.WriteString(fmt.Sprintf(", Timeout: %d", obj.timeout))
	buffer.WriteString(fmt.Sprintf(", Cond: %+v", obj.cond))
	buffer.WriteString(fmt.Sprintf(", Reply: %+v", obj.reply))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for StatusTester
/////////////////////////////////////////////////////////////////

// Test checks whether the channel response is in WAIT mode or not
func (obj *BlockingChannelResponse) Test(status tgdb.ChannelResponseStatus) bool {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	if obj.status == tgdb.Waiting {
		return true
	}
	return false
}


var gReaders int64

type ChannelReader struct {
	channel   tgdb.TGChannel
	isRunning bool
	name      string
	readerNum int64
}

func DefaultChannelReader() *ChannelReader {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelReader{})

	newChannelReader := ChannelReader{
		isRunning: false,
		readerNum: atomic.AddInt64(&gReaders, 1),
	}

	return &newChannelReader
}

func NewChannelReader(rChannel tgdb.TGChannel) *ChannelReader {
	newChannelReader := DefaultChannelReader()
	newChannelReader.channel = rChannel
	newChannelReader.name = fmt.Sprintf("TGLinkReader@[%s-%d]", rChannel.GetClientId(), newChannelReader.readerNum)
	return newChannelReader
}

/////////////////////////////////////////////////////////////////
// Helper functions for ChannelReader
/////////////////////////////////////////////////////////////////

// readAndProcessLoop reads a message from the network and processes it
func (obj *ChannelReader) readAndProcessLoop() {
	if obj == nil {
		return
	}
	//logger.Log(fmt.Sprintf("Entering ChannelReader:readAndProcessLoop w/ Reader object: '%+v'", obj.String()))
	for {
		// Terminating Conditions for this Infinite Loop are:
		// 	(a) Break if the channel reader is NOT RUNNING
		// 	(b) Break if the channel reader / GO Routine is INTERRUPTED
		// 	(c) Break if the channel is CLOSED
		// 	(d) Break if the request on the wire is to DISCONNECT from the SERVER
		// 	(e) Break - in case of ERROR - if the exceptionResult is NOT RetryOperation
		// Looping Conditions for this Infinite Loop are:
		// 	(a) Continue if the message on the wire is EMPTY or cannot be READ
		// 	(b) Continue if the message on the wire is PING (HeartBeat) message w/o Processing the message
		// 	(c) Continue if the message on the wire is ANYTHING else after Processing the message
		// 	(d) Continue - in case of ERROR - if the exceptionResult is ANYTHING OTHER THAN RetryOperation after setting the reply on channelResponse

		if !obj.isRunning {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("WARNING: Breaking ChannelReader:readAndProcessLoop loop since reader is not running '%+v'", obj.isRunning))
			}
			break
		}

		if obj.channel.IsClosed() {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("WARNING: Breaking ChannelReader:readAndProcessLoop loop since channel is closed"))
			}
			break
		}

		// Execute Derived Channel's method
		msg, err := obj.channel.ReadWireMsg()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop obj.channel.ReadWireMsg failed w/ '%+v'", err.Error()))
			if !obj.isRunning {
				logger.Info(fmt.Sprintf("INFO: Breaking ChannelReader:readAndProcessLoop reader is not running (2) '%+v'", obj.isRunning))
				break
			}
			exceptionResult := channelHandleException(obj.channel, err, true)
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop obj.channel.ReadWireMsg failed - exceptionResult '%+v'", exceptionResult))
			for _, resp := range obj.channel.GetResponses() {
				if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside ChannelReader:readAndProcessLoop about to channelResponse.SetReply() w/ new EXCEPTION MSG"))
				}
				resp.SetReply(NewExceptionMessageWithType(int(exceptionResult.ExceptionType), exceptionResult.ExceptionMessage))
			}
			if exceptionResult.ExceptionType != RetryOperation {
				logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop since Reader thread returned w/o Retrying due to error - exceptionResult '%+v'", exceptionResult))
				break
			}
			//logger.Error(fmt.Sprintf("INFO: Breaking ChannelReader:readAndProcessLoop loop - Read Wire Message resulted in error: '%+v'", err))
			//break
		}

		if obj.channel.IsClosed() {
			logger.Warning(fmt.Sprintf("WARNING: Breaking ChannelReader:readAndProcessLoop loop since channel is closed"))
			break
		}

		if msg == nil {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop - Read Message Again since MSG is NIL"))
			}
			continue
		}

		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop - Read Message of type '%+v'", msg.GetVerbId()))
		}
		if msg.GetVerbId() == VerbPingMessage {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop Trying to Read Message Again since MSG is PingMessage"))
			}
			continue
		}

		// Server Requested to disconnect
		if msg.GetVerbId() == VerbSessionForcefullyTerminated {
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Breaking ChannelReader:readAndProcessLoop loop w/ Forceful Termination Message is '%+v'", msg.String()))
			}
			channelTerminated(obj.channel, msg.(*SessionForcefullyTerminatedMessage).GetKillString())
			obj.isRunning = false
			break
		}

		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop Processing Message of type '%+v'", msg.GetVerbId()))
		}
		err = channelProcessMessage(obj.channel, msg)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop channelProcessMessage() failed w/ '%+v'", err.Error()))
			if !obj.isRunning {
				logger.Info(fmt.Sprintf("INFO: Inside ChannelReader:readAndProcessLoop reader is not running (3) '%+v'", obj.isRunning))
				break
			}
			exceptionResult := channelHandleException(obj.channel, err, true)
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop channelProcessMessage() failed - exceptionResult (2) '%+v'", exceptionResult))
			for _, resp := range obj.channel.GetResponses() {
				if logger.IsDebug() {
					logger.Debug(fmt.Sprint("Inside ChannelReader:readAndProcessLoop about to channelResponse.SetReply() w/ new EXCEPTION MSG (2)"))
				}
				resp.SetReply(NewExceptionMessageWithType(int(exceptionResult.ExceptionType), exceptionResult.ExceptionMessage))
			}
			if exceptionResult.ExceptionType != RetryOperation {
				logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop since Reader thread returned w/o Retrying due to error (2) - exceptionResult '%+v'", exceptionResult))
				break
			}
			//logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop - ProcessMessage resulted in error: '%+v'", err))
			//break
		}

		if obj.channel.IsClosed() {
			logger.Warning(fmt.Sprintf("WARNING: Breaking ChannelReader:readAndProcessLoop loop since channel is closed"))
			break
		}
	} // End of Infinite Loop
	obj.isRunning = false
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning ChannelReader:readAndProcessLoop w/ Reader object: '%+v'", obj.String()))
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> ChannelReader
/////////////////////////////////////////////////////////////////

// Start starts the channel reader
func (obj *ChannelReader) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelReader:Start ..."))
	if !obj.isRunning {
		obj.isRunning = true
		// Start reading and processing messages from the wire
		//obj.readAndProcessLoop()
		go obj.readAndProcessLoop()
	}
	//logger.Log(fmt.Sprint("Returning ChannelReader:Start ..."))
}

// Stop stops the channel reader
func (obj *ChannelReader) Stop() {
	//logger.Log(fmt.Sprint("Entering ChannelReader:Stop ..."))
	if obj.isRunning {
		// Finish / Flush any remaining processing
		//obj.readAndProcessLoop()
		go obj.readAndProcessLoop()
		obj.isRunning = false
	}
	//logger.Log(fmt.Sprint("Returning ChannelReader:Stop ..."))
}

func (obj *ChannelReader) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelReader:{")
	buffer.WriteString(fmt.Sprintf("Name: %+v", obj.name))
	buffer.WriteString(fmt.Sprintf(", IsRunning: %+v", obj.isRunning))
	buffer.WriteString(fmt.Sprintf(", ReaderNum: %d", obj.readerNum))
	buffer.WriteString(fmt.Sprintf(", Channel: %s", obj.channel.String()))
	buffer.WriteString("}")
	return buffer.String()
}


type DataCryptoGrapher struct {
	sessionId      int64
	remoteCert     *x509.Certificate
	pubKey         crypto.PublicKey
	//algoParameters *pkix.AlgorithmIdentifier
}

func DefaultDataCryptoGrapher() *DataCryptoGrapher {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DataCryptoGrapher{})

	newChannelUrl := DataCryptoGrapher{
		sessionId: 0,
	}

	return &newChannelUrl
}

func NewDataCryptoGrapher(sessionId int64, serverCertBytes []byte) (*DataCryptoGrapher, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering NewDataCryptoGrapher() w/ serverCertBytes as '%+v'", serverCertBytes))
	}
	newCryptoGrapher := DefaultDataCryptoGrapher()
	newCryptoGrapher.sessionId = sessionId

	// TODO: Uncomment once DataCryptoGrapher is implemented
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - about to x509.ParseCertificate(()"))
	//cert, err := x509.ParseCertificate(serverCertBytes)
	//if err != nil {
	//	errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
	//	return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
	//}
	//newCryptoGrapher.remoteCert = cert
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - parsed certificate as'%+v'", cert))
	//
	//logger.Log(fmt.Sprint("Inside NewDataCryptoGrapher - about to x509.ParsePKIXPublicKey()"))
	//pubKey, err := x509.ParsePKIXPublicKey(serverCertBytes)
	//if err != nil {
	//	errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse PUBLIC KEY from the certificate buffer")
	//	return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err.Error())
	//}
	//newCryptoGrapher.pubKey = pubKey
	//logger.Log(fmt.Sprintf("Inside NewDataCryptoGrapher - parsed public key as'%+v'", pubKey))

	/**
	algoParams, err1 := getAlgorithmParameters(newCryptoGrapher.pubKey)
	if err1 != nil {
		errMsg := fmt.Sprint("NewDataCryptoGrapher -- Unable to parse CERTIFICATE from the certificate buffer")
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, err1.Error())
	}
	newCryptoGrapher.algoParameters = algoParams
	*/
	return newCryptoGrapher, nil
}

/////////////////////////////////////////////////////////////////
// Helper functions for DataCryptoGrapher
/////////////////////////////////////////////////////////////////

func (obj *DataCryptoGrapher) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DataCryptoGrapher:{")
	buffer.WriteString(fmt.Sprintf("SessionId: %d", obj.sessionId))
	buffer.WriteString(fmt.Sprintf(", RemoteCert: %+v", obj.remoteCert))
	buffer.WriteString(fmt.Sprintf(", PubKey: %+v", obj.pubKey))
	//buffer.WriteString(fmt.Sprintf(", AlgoParameters: %d", obj.algoParameters))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Private functions for DataCryptoGrapher
/////////////////////////////////////////////////////////////////

func getAlgorithmParameters(pubKey crypto.PublicKey) (*pkix.AlgorithmIdentifier, tgdb.TGError) {
	// TODO: Uncomment once DataCryptoGrapher is implemented
	/**
	if (publicKey == null) return null;

	if (publicKey instanceof DSAPublicKey) {
		AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
		DSAPublicKey dsakey = (DSAPublicKey) publicKey;
		DSAParameterSpec dsaParams = (DSAParameterSpec) dsakey.getParams();
		algparams.init(dsaParams);
		return algparams;
	}

	if (publicKey instanceof ECPublicKey) {
		AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
		ECPublicKey eckey = (ECPublicKey) publicKey;
		ECParameterSpec ecParams = (ECParameterSpec) eckey.getParams();
		algparams.init(ecParams);
		return algparams;
	}

	if (publicKey instanceof DHPublicKey) {
		AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
		DHPublicKey dhkey = (DHPublicKey) publicKey;
		DHParameterSpec dhParams = (DHParameterSpec) dhkey.getParams();
		algparams.init(dhParams);
		return algparams;
	}

	return null;  //RSA doesn't have
	*/
	// No-op for Now!
	return nil, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGDataCryptoGrapher
/////////////////////////////////////////////////////////////////

// Decrypt decrypts the buffer
//func (obj *DataCryptoGrapher) Decrypt(encBuffer []byte) ([]byte, types.TGError) {
func (obj *DataCryptoGrapher) Decrypt(is tgdb.TGInputStream) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering DataCryptoGrapher:Decrypt()"))
	}
	out := DefaultProtocolDataOutputStream()
	buf := make([]byte, 0)

	rand, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading rand from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", rand))
	}

	len, err := is.(*ProtocolDataInputStream).ReadLong()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading len from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", len))
	}

	cnt := len / 8
	rem := len % 8

	for i:=0; i<int(cnt); i++ {
		val, err := is.(*ProtocolDataInputStream).ReadLong()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading val from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", val))
		}

		org := val ^ rand
		_ = out.WriteLongAsBytes(org)
	}

	for i:=0; i<int(rem); i++ {
		val, err := is.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning DataCryptoGrapher:Decrypt w/ Error in reading val from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside DataCryptoGrapher:Decrypt read resultId as '%+v'", val))
		}

		out.WriteByte(int(val))
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Returning DataCryptoGrapher:Decrypt() w/ decrypted buffer as '%+v'", buf))
	}
	return out.ToByteArray()
}

// Encrypt encrypts the buffer
func (obj *DataCryptoGrapher) Encrypt(rawBuf []byte) ([]byte, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Entering DataCryptoGrapher:Encrypt() w/ raw buffer as '%+v'", rawBuf))
	}
	// TODO: Uncomment once DataCryptoGrapher is implemented
	/**
	try {
		Cipher cipher = Cipher.getInstance(publicKey.getAlgorithm());
		cipher.init(Cipher.ENCRYPT_MODE, publicKey, algparams);
		return cipher.doFinal(data);
	}
	catch (Exception e) {
		throw new TGException(e);
	}

	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
	

	block, err := blowfish.NewCipher(obj.remoteCert.RawSubjectPublicKeyInfo)
	if err != nil {
		return nil, GetErrorByType(TGErrorSecurityException, INTERNAL_SERVER_ERROR, err.Error(), "")
	}
	encryptedBuf := make([]byte, aes.BlockSize+len(rawBuf))
	iv := encryptedBuf[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedBuf[aes.BlockSize:], rawBuf)
	fmt.Printf("%x\n", encryptedBuf)
	
	block := blowfish.NewCipher(obj.remoteCert.RawSubjectPublicKeyInfo)
	encryptedBuf := make([]byte, aes.BlockSize+len(decBuffer))
	iv := encryptedBuf[:aes.BlockSize]

	algo := obj.remoteCert.PublicKeyAlgorithm
	switch algo {
	case x509.RSA:
	case x509.DSA:
		block, err := des.NewCipher(rawBuf)
		if err != nil {
			return nil, exception.GetErrorByType(types.TGErrorSecurityException, types.INTERNAL_SERVER_ERROR, err.Error(), "")
		}
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(encryptedBuf[aes.BlockSize:], rawBuf)
	case x509.ECDSA:
	}
	*/
	//logger.Log(fmt.Sprintf("Returning DataCryptoGrapher:Decrypt() w/ encrypted buffer as '%+v'", encryptedBuf))
	return nil, nil
}


const (
	dataBufferSize = 32 * 1024 // 32 KB
)

type TCPChannel struct {
	*AbstractChannel
	shutdownLock   sync.RWMutex // rw-lock for synchronizing read-n-update of env configuration
	isSocketClosed bool         // indicate if the connection is already closed
	msgCh          chan tgdb.TGMessage
	socket         *net.TCPConn
	input          *ProtocolDataInputStream
	output         *ProtocolDataOutputStream
}

func DefaultTCPChannel() *TCPChannel {
	newChannel := TCPChannel{
		AbstractChannel: DefaultAbstractChannel(),
		msgCh:           make(chan tgdb.TGMessage),
		isSocketClosed:  false,
	}
	buff := make([]byte, 0)
	newChannel.input = NewProtocolDataInputStream(buff)
	newChannel.output = NewProtocolDataOutputStream(0)
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	return &newChannel
}

func NewTCPChannel(linkUrl *LinkUrl, props *SortedProperties) *TCPChannel {
	newChannel := TCPChannel{
		AbstractChannel: NewAbstractChannel(linkUrl, props),
		msgCh:           make(chan tgdb.TGMessage),
		isSocketClosed:  false,
	}
	buff := make([]byte, 0)
	newChannel.input = NewProtocolDataInputStream(buff)
	newChannel.output = NewProtocolDataOutputStream(0)
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	return &newChannel
}

/////////////////////////////////////////////////////////////////
// Private functions for TCPChannel
/////////////////////////////////////////////////////////////////


func (obj *TCPChannel) DoAuthenticateForRESTConsumer() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:DoAuthenticateForRESTConsumer"))
	}
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := CreateMessageForVerb(VerbAuthenticateRequest)

	channelResponse := NewBlockingChannelResponse(msgRequest.GetRequestId(), -1)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: TCPChannel::DoAuthenticateForRESTConsumer failed w/ '%+v'", err.Error()))
		return err
	}

	msgRequest.(*AuthenticateRequestMessage).SetClientId(obj.clientId)
	msgRequest.(*AuthenticateRequestMessage).SetInboxAddr(obj.inboxAddress)
	msgRequest.(*AuthenticateRequestMessage).SetUserName(obj.GetChannelUserName())
	msgRequest.(*AuthenticateRequestMessage).SetPassword(obj.GetChannelPassword())
	msgRequest.(*AuthenticateRequestMessage).SetDatabaseName(obj.getDatabaseName())
	msgResponse, err := channelSendRequest(obj, msgRequest, channelResponse, true)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::DoAuthenticateForRESTConsumer channelSendRequest failed w/ '%+v'", err.Error()))
		return err
	}
	if ! msgResponse.(*AuthenticateResponseMessage).IsSuccess() {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::DoAuthenticateForRESTConsumer msgResponse.(*AuthenticateResponseMessage).IsSuccess() failed"))
		return NewTGBadAuthenticationWithRealm(INTERNAL_SERVER_ERROR, TGErrorBadAuthentication, "Bad username/password combination", "", "tgdb")
	}

	obj.setChannelAuthToken(msgResponse.GetAuthToken())
	obj.setChannelSessionId(msgResponse.GetSessionId())

	cryptoDataGrapher, err := NewDataCryptoGrapher(msgResponse.GetSessionId(), msgResponse.(*AuthenticateResponseMessage).GetServerCertBuffer())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::DoAuthenticateForRESTConsumer NewDataCryptoGrapher failed w/ '%+v'", err.Error()))
		return err
	}
	obj.setDataCryptoGrapher(cryptoDataGrapher)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:DoAuthenticateForRESTConsumer Successfully authenticated for user: '%s'", obj.GetChannelUserName()))
	}
	return nil
}


func (obj *TCPChannel) DoAuthenticate() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:doAuthenticate"))
	}
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := CreateMessageForVerb(VerbAuthenticateRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: TCPChannel::doAuthenticate CreateMessageForVerb(VerbAuthenticateRequest) failed w/ '%+v'", err.Error()))
		return err
	}

	msgRequest.(*AuthenticateRequestMessage).SetClientId(obj.clientId)
	msgRequest.(*AuthenticateRequestMessage).SetInboxAddr(obj.inboxAddress)
	msgRequest.(*AuthenticateRequestMessage).SetUserName(obj.GetChannelUserName())
	msgRequest.(*AuthenticateRequestMessage).SetPassword(obj.GetChannelPassword())
	msgRequest.(*AuthenticateRequestMessage).SetDatabaseName(obj.getDatabaseName())

	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:doAuthenticate about to request reply for request '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:doAuthenticate received reply as '%+v'", msgResponse.String()))
	if ! msgResponse.(*AuthenticateResponseMessage).IsSuccess() {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate msgResponse.(*AuthenticateResponseMessage).IsSuccess() failed"))
		return NewTGBadAuthenticationWithRealm(INTERNAL_SERVER_ERROR, TGErrorBadAuthentication, "Bad username/password combination", "", "tgdb")
	}

	obj.setChannelAuthToken(msgResponse.GetAuthToken())
	obj.setChannelSessionId(msgResponse.GetSessionId())

	cryptoDataGrapher, err := NewDataCryptoGrapher(msgResponse.GetSessionId(), msgResponse.(*AuthenticateResponseMessage).GetServerCertBuffer())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate NewDataCryptoGrapher failed w/ '%+v'", err.Error()))
		return err
	}
	obj.setDataCryptoGrapher(cryptoDataGrapher)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:doAuthenticate Successfully authenticated for user: '%s'", obj.GetChannelUserName()))
	}
	return nil
}

func (obj *TCPChannel) SetAuthToken(token int64) {
	obj.setChannelAuthToken(token)
}

func (obj *TCPChannel) performHandshake(sslMode bool) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:performHandshake"))
	}
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := CreateMessageForVerb(VerbHandShakeRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake CreateMessageForVerb(VerbAuthenticateRequest) failed w/ '%+v'", err.Error()))
		return err
	}

	msgRequest.(*HandShakeRequestMessage).SetRequestType(InitiateRequest)

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:performHandshake about to request reply for InitiateRequest '%+v'", msgRequest.String()))
	}
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::doAuthenticate channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:performHandshake received reply as '%+v'", msgResponse.String()))
	}
	if msgResponse.GetVerbId() != VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*SessionForcefullyTerminatedMessage).GetKillString()
			return NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	response := msgResponse.(*HandShakeResponseMessage)
	if response.GetResponseStatus() != ResponseAcceptChallenge {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake response.GetResponseStatus() is NOT ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	// Validate the version specific information on the response object
	serverVersion := response.GetChallenge()
	clientVersion := GetClientVersion()
	err = obj.validateHandshakeResponseVersion(serverVersion, clientVersion)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake validateHandshakeResponseVersion failed w/ '%+v'", err.Error()))
		return err
	}

	challenge := clientVersion.GetVersionAsLong()

	// Ignore Error Handling
	_ = msgRequest.(*HandShakeRequestMessage).UpdateSequenceAndTimeStamp(-1)
	msgRequest.(*HandShakeRequestMessage).SetRequestType(ChallengeAccepted)
	msgRequest.(*HandShakeRequestMessage).SetSslMode(sslMode)
	msgRequest.(*HandShakeRequestMessage).SetChallenge(challenge)

	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:performHandshake about to request reply for ChallengeAccepted '%+v'", msgRequest.String()))
	msgResponse, err = channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake channelRequestReply failed w/ '%+v'", err.Error()))
		return err
	}
	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:performHandshake received reply (2) as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*SessionForcefullyTerminatedMessage).GetKillString()
			return NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	response = msgResponse.(*HandShakeResponseMessage)
	if response.GetResponseStatus() != ResponseProceedWithAuthentication {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::performHandshake response.GetResponseStatus() is NOT ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel::performHandshake Handshake w/ Remote Server is successful."))
	}
	return nil
}

func (obj *TCPChannel) setSocket(newSocket *net.TCPConn) tgdb.TGError {
	obj.socket = newSocket
	err := obj.socket.SetNoDelay(true)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel:setSocket Failed to set NoDelay flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set NoDelay flag to true")
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	}

	err = obj.socket.SetLinger(0) // <= 0 means Do not linger
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel:setSocket Failed to set NoLinger flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set NoLinger flag to true")
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	}

	buff := make([]byte, dataBufferSize)
	obj.input = NewProtocolDataInputStream(buff)
	obj.input.BufLen = 0
	obj.output = NewProtocolDataOutputStream(dataBufferSize)
	//clientId = properties.get(ConfigName.ChannelClientId.getName());
	//if (clientId == null) {
	//	clientId = properties.get(ConfigName.ChannelClientId.getAlias());
	//	if (clientId == null) {
	//		clientId = TGEnvironment.getInstance().getChannelClientId();
	//	}
	//}
	clientId := obj.GetProperties().GetProperty(GetConfigFromKey(ChannelClientId), "")
	obj.setChannelClientId(clientId)
	obj.setChannelInboxAddr(obj.socket.RemoteAddr().String()) //SS:TODO: Is this correct
	return nil
}

func (obj *TCPChannel) setBuffers(newSocket *net.TCPConn) tgdb.TGError {
	sendSize := obj.channelProperties.GetPropertyAsInt(GetConfigFromKey(ChannelSendSize))
	if sendSize > 0 {
		err := newSocket.SetWriteBuffer(sendSize*1024)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::setBuffers newSocket.SetWriteBuffer failed w/ '%+v'", err.Error()))
			errMsg := fmt.Sprintf("TCPChannel:setBuffers unable to set write buffer limit to '%d'", sendSize*1024)
			return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
		}
	}
	receiveSize := obj.channelProperties.GetPropertyAsInt(GetConfigFromKey(ChannelRecvSize))
	if receiveSize > 0 {
		err := newSocket.SetReadBuffer(receiveSize*1024)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::setBuffers SetReadBuffer failed w/ '%+v'", err.Error()))
			errMsg := fmt.Sprintf("TCPChannel:setBuffers unable to set read buffer limit to '%d'", receiveSize*1024)
			return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
		}
	}
	return nil
}

func (obj *TCPChannel) tryRead() (tgdb.TGMessage, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("======> Entering TCPChannel:tryRead"))
	n, err := obj.input.Available()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::tryRead obj.input.Available() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::tryRead there is no data available to be read"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if n <= 0 {
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:tryRead about to request message '%d' bytes from the wire", n))
	}
	return obj.ReadWireMsg()
}

func (obj *TCPChannel) validateHandshakeResponseVersion(sVersion int64, cVersion *TGClientVersion) tgdb.TGError {
	serverVersion := NewTGServerVersion(sVersion)
	sStrVer := serverVersion.GetVersionString()

	cStrVer := cVersion.GetVersionString()

	if 	serverVersion.GetMajor() == cVersion.GetMajor() &&
		serverVersion.GetMinor() == cVersion.GetMinor() &&
		serverVersion.GetUpdate() == cVersion.GetUpdate() {
		return nil
	}

	errMsg := fmt.Sprintf("======> Inside SSLChannel:validateHandshakeResponseVersion - Version mismatch between client(%s) & server(%s)", cStrVer, sStrVer)
	if logger.IsDebug() {
		logger.Debug(errMsg)
	}
	return GetErrorByType(TGErrorVersionMismatchException, "", errMsg, "")
}

func (obj *TCPChannel) writeLoop(done chan bool) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Entering TCPChannel:writeLoop"))
	}
	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("======> Inside TCPChannel:writeLoop entering infinite loop"))
		}
		select { // Non-blocking channel operation
		case msg, ok := <-obj.msgCh: // Retrieve the message from the channel
			if !ok {
				//if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
				//	logMessage("TCPChannel::writeLoop unable to retrieve msg from the channel);
				//}
				logger.Error(fmt.Sprint("ERROR: Returning TCPChannel:writeLoop unable to retrieve message from obj.msgCh"))
				return
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("======> Inside TCPChannel:writeLoop retrieved message from obj.msgCh as '%+v'", msg.String()))
			}

			err := obj.writeToWire(msg)
			if err != nil {
				// TODO: Revisit later - Do something
				logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::writeLoop unable to obj.writeToWire w/ '%+v'", err.Error()))
				return
			}

			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("======> Inside TCPChannel:writeLoop successfully wrote message '%+v' on the socket", msg.String()))
			}
			break
		default:
			// TODO: Revisit later - Do something
		}
	}	// End of Infinite Loop
	// Send an acknowledgement of completion to the parent thread
	done <- true
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Returning TCPChannel:writeLoop"))
	}
}

func (obj *TCPChannel) writeToWire(msg tgdb.TGMessage) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:writeToWire w/ Msg: '%+v'", msg.String()))
	}
	obj.DisablePing()
	msgBytes, bufLen, err := msg.ToBytes()
	if err != nil {
		errMsg := fmt.Sprintf("TCPChannel::writeToWire unable to convert message into byte format")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, err.Error()))
		//return exception.GetErrorByType(types.TGErrorIOException, "TGErrorProtocolNotSupported", errMsg, err.GetErrorMsg())
		return err
	}

	// Clear timeout deadlines set at the time of creation of the socket
	sErr := obj.socket.SetDeadline(time.Time{})
	if sErr != nil {
		errMsg := fmt.Sprintf("TCPChannel::writeToWire - unable to clear the deadline over TCP socket")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	//// Reset timeout deadlines starting from NOW!!!
	//timeout := NewTGEnvironment().GetChannelConnectTimeout()
	//sErr = obj.socket.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	//if sErr != nil {
	//	errMsg := fmt.Sprintf("TCPChannel::writeToWire - unable to reset the deadline over TCP socket")
	//	logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
	//	return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	//}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:writeToWire about to write message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	}
	// Put the data packet on the socket for network transmission
	_, sErr = obj.socket.Write(msgBytes[0:bufLen])
	if sErr != nil {
		errMsg := fmt.Sprintf("TCPChannel::writeToWire - unable to send message bytes over TCP socket")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return GetErrorByType(TGErrorIOException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:writeToWire successfully wrote message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	}
	return nil
}

func intToBytes(value int, bytes []byte, offset int) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:intToBytes w/ value as '%d', byteArray as '%+v' and offset '%d'", value, bytes, offset))
	}
	for i := 0; i < 4; i++ {
		bytes[offset+i] = byte((value >> uint(8*(3-i))) & 0xff)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:intToBytes w/ byteArray as '%+v'", bytes))
	}
}

/////////////////////////////////////////////////////////////////
// Helper functions for TCPChannel
/////////////////////////////////////////////////////////////////

func (obj *TCPChannel) GetIsClosed() bool {
	return obj.isSocketClosed
}

func (obj *TCPChannel) SetIsClosed(flag bool) {
	obj.isSocketClosed = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// ChannelLock locks the communication channel between TGDB client and server
func (obj *TCPChannel) ChannelLock() {
	obj.sendLock.Lock()
}

// ChannelUnlock unlocks the communication channel between TGDB client and server
func (obj *TCPChannel) ChannelUnlock() {
	obj.sendLock.Unlock()
}

// Connect connects the underlying channel using the URL end point
func (obj *TCPChannel) Connect() tgdb.TGError {
	return channelConnect(obj)
}

// DisablePing disables the pinging ability to the channel
func (obj *TCPChannel) DisablePing() {
	obj.needsPing = false
}

// Disconnect disconnects the channel from its URL end point
func (obj *TCPChannel) Disconnect() tgdb.TGError {
	return channelDisConnect(obj)
}

// EnablePing enables the pinging ability to the channel
func (obj *TCPChannel) EnablePing() {
	obj.needsPing = true
}

// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
func (obj *TCPChannel) ExceptionLock() {
	obj.exceptionLock.Lock()
}

// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
func (obj *TCPChannel) ExceptionUnlock() {
	obj.exceptionLock.Unlock()
}

// GetAuthToken gets Authorization Token
func (obj *TCPChannel) GetAuthToken() int64 {
	return obj.authToken
}

// GetClientId gets Client Name
func (obj *TCPChannel) GetClientId() string {
	return obj.clientId
}

// GetChannelURL gets the channel URL
func (obj *TCPChannel) GetChannelURL() tgdb.TGChannelUrl {
	return obj.channelUrl
}

// GetConnectionIndex gets the Connection Index
func (obj *TCPChannel) GetConnectionIndex() int {
	return obj.connectionIndex
}

// GetExceptionCondition gets the Exception Condition
func (obj *TCPChannel) GetExceptionCondition() *sync.Cond {
	return obj.exceptionCond
}

// GetLinkState gets the Link/channel State
func (obj *TCPChannel) GetLinkState() tgdb.LinkState {
	return obj.channelLinkState
}

// GetNoOfConnections gets number of connections this channel has
func (obj *TCPChannel) GetNoOfConnections() int32 {
	return obj.numOfConnections
}

// GetPrimaryURL gets the Primary URL
func (obj *TCPChannel) GetPrimaryURL() tgdb.TGChannelUrl {
	return obj.primaryUrl
}

// GetProperties gets the channel Properties
func (obj *TCPChannel) GetProperties() tgdb.TGProperties {
	return obj.channelProperties
}

// GetReader gets the channel Reader
func (obj *TCPChannel) GetReader() tgdb.TGChannelReader {
	return obj.reader
}

// GetResponses gets the channel Response Map
func (obj *TCPChannel) GetResponses() map[int64]tgdb.TGChannelResponse {
	return obj.responses
}

// GetSessionId gets Session id
func (obj *TCPChannel) GetSessionId() int64 {
	return obj.sessionId
}

// GetTracer gets the channel Tracer
func (obj *TCPChannel) GetTracer() tgdb.TGTracer {
	return obj.tracer
}

// IsChannelPingable checks whether the channel is pingable or not
func (obj *TCPChannel) IsChannelPingable() bool {
	return obj.needsPing
}

// IsClosed checks whether channel is open or closed
func (obj *TCPChannel) IsClosed() bool {
	return isChannelClosed(obj)
}

// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
func (obj *TCPChannel) SendMessage(msg tgdb.TGMessage) tgdb.TGError {
	return channelSendMessage(obj, msg, true)
}

// SendRequest sends a Message, waits for a response in the message format, and blocks the thread till it gets the response
func (obj *TCPChannel) SendRequest(msg tgdb.TGMessage, response tgdb.TGChannelResponse) (tgdb.TGMessage, tgdb.TGError) {
	return channelSendRequest(obj, msg, response, true)
}

// SetChannelLinkState sets the Link/channel State
func (obj *TCPChannel) SetChannelLinkState(state tgdb.LinkState) {
	obj.channelLinkState = state
}

// SetChannelURL sets the channel URL
func (obj *TCPChannel) SetChannelURL(url tgdb.TGChannelUrl) {
	obj.channelUrl = url.(*LinkUrl)
}

// SetConnectionIndex sets the connection index
func (obj *TCPChannel) SetConnectionIndex(index int) {
	obj.connectionIndex = index
}

// SetNoOfConnections sets number of connections
func (obj *TCPChannel) SetNoOfConnections(count int32) {
	obj.numOfConnections = count
}

// SetResponse sets the ChannelResponse Map
func (obj *TCPChannel) SetResponse(reqId int64, response tgdb.TGChannelResponse) {
	obj.responses[reqId] = response
}

// Start starts the channel so that it can send and receive messages
func (obj *TCPChannel) Start() tgdb.TGError {
	return channelStart(obj)
}

// Stop stops the channel forcefully or gracefully
func (obj *TCPChannel) Stop(bForcefully bool) {
	channelStop(obj, bForcefully)
}

// CreateSocket creates the physical link socket
func (obj *TCPChannel) CreateSocket() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:CreateSocket"))
	}
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	obj.SetChannelLinkState(tgdb.LinkNotConnected)
	host := obj.channelUrl.urlHost
	port := obj.channelUrl.urlPort
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:CreateSocket attempting to resolve address for '%s'", serverAddr))

	tcpAddr, tErr := net.ResolveTCPAddr(tgdb.ProtocolTCP.String(), serverAddr)
	if tErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket net.ResolveTCPAddr failed w/ '%+v'", tErr.Error()))
		errMsg := fmt.Sprintf("TCPChannel:CreateSocket unable to resolve channel address '%s'", serverAddr)
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, tErr.Error())
	}
	//logger.Debug(fmt.Sprintf("======> Inside TCPChannel:CreateSocket resolved TCP address for '%s' as '%+v'", serverAddr, tcpAddr))

	tcpConn, cErr := net.DialTCP(tgdb.ProtocolTCP.String(), nil, tcpAddr)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to connect to the server at '%s' w/ '%+v'", serverAddr, cErr.Error()))
		failureMessage := fmt.Sprintf("Failed to connect to the server at '%s'" + serverAddr)
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:CreateSocket created TCP connection for '%s' as '%+v'", serverAddr, tcpConn))
	}

	//timeout := NewTGEnvironment().GetChannelConnectTimeout()
	//dErr := tcpConn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	//if dErr != nil {
	//	logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set deadline of '%+v' seconds on the connection to the server w/ '%+v'", time.Duration(timeout) * time.Second, dErr.Error()))
	//	failureMessage := fmt.Sprintf("Failed to set the timeout '%d' on socket", timeout)
	//	return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, dErr.Error())
	//}

	err := tcpConn.SetKeepAlive(true)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set keep alive flag to true w/ '%+v'", err.Error()))
		failureMessage := fmt.Sprint("Failed to set keep alive flag to true")
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}

	// Set Read / Write Buffer Size on the socket
	tcErr := obj.setBuffers(tcpConn)
	if tcErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set buffers w/ '%+v'", tcErr.Error()))
		return tcErr
	}
	tcErr = obj.setSocket(tcpConn)
	if tcErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CreateSocket Failed to set socket value to the object w/ '%+v'", tcErr.Error()))
		return tcErr
	}
	obj.SetIsClosed(false)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:CreateSocket w/ TCP Connection as '%+v'", *obj.socket))
	}
	return nil
}

// CloseSocket closes the socket
func (obj *TCPChannel) CloseSocket() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:CloseSocket w/ socket: '%+v'", obj.socket))
	}
	obj.shutdownLock.Lock()
	defer func() {
		obj.SetIsClosed(true)
		obj.shutdownLock.Unlock()
		obj.socket = nil
		obj.input = nil
		obj.output = nil
	} ()

	if obj.socket != nil {
		cErr := obj.socket.Close()
		if cErr != nil {
			failureMessage := "Failed to close the socket to the server"
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::CloseSocket %s w/ '%+v'", failureMessage, cErr.Error()))
			//return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
		}
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:CloseSocket for socket: '%+v'", obj.socket))
	}
	return nil
}

// OnConnect executes all the channel specific activities
func (obj *TCPChannel) OnConnect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:OnConnect about to tryRead w/ socket: '%+v'", obj.socket))
	}
	msg, err := obj.tryRead()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.tryRead() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect there is no data available to be read"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}
	if msg != nil {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("======> Inside TCPChannel:OnConnect tryRead() read Message as '%+v'", msg.String()))
		}
	}

	if msg != nil && msg.GetVerbId() == VerbSessionForcefullyTerminated {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:OnConnect since Message is of Forceful Termination Type"))
		return NewTGChannelDisconnectedWithMsg(msg.(*SessionForcefullyTerminatedMessage).GetKillString())
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:OnConnect about to performHandshake"))
	}
	err = obj.performHandshake(false)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.performHandshake() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect error in performing handshake with server"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:OnConnect about to doAuthenticate"))
	}
	err = obj.DoAuthenticate()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::OnConnect obj.doAuthenticate() failed w/ '%+v'", err.Error()))
		errMsg := "TCPChannel::OnConnect error in authentication with server"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:OnConnect w/ socket: '%+v'", obj.socket))
	}
	return nil
}

// ReadWireMsg reads the message from the wire in the form of byte stream
func (obj *TCPChannel) ReadWireMsg() (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:ReadWireMsg  w/ socket: '%+v'", obj.socket))
	}
	obj.input.BufLen = dataBufferSize
	in := obj.input
	if in == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:ReadWireMsg since obj.input is NIL"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}

	obj.DisablePing()
	if obj.GetIsClosed() {
		logger.Warning(fmt.Sprint("WARNING: Returning TCPChannel:ReadWireMsg since TCP channel is Closed"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}
	obj.lastActiveTime = time.Now()

	totalBuffer := make ([]byte, 4)
	count, err11 := obj.socket.Read(totalBuffer)
	var totalBytesOnSocket uint32
	if err11 != nil || count <= 0 {
		errMsg := "TCPChannel::ReadWireMsg obj.socket.Read failed"
		logger.Error(fmt.Sprintf("ERROR: Returning %s", errMsg))
		return nil, GetErrorByType(TGErrorIOException, "", errMsg, "")
	} else {
		totalBytesOnSocket = binary.BigEndian.Uint32(totalBuffer)
	}
	// Read the data on the socket
	buff := make([]byte, totalBytesOnSocket-4)
	nInterim := 0
	for ;nInterim < int(totalBytesOnSocket-4); {
		nCurrent, sErr := obj.socket.Read(buff[nInterim: totalBytesOnSocket-4])
		if sErr != nil || nCurrent <= 0 {
			errMsg := "TCPChannel::ReadWireMsg obj.socket.Read failed"
			logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
			return nil, GetErrorByType(TGErrorIOException, "", errMsg, sErr.Error())
		}
		nInterim = nInterim + nCurrent
	}

	slice := buff[0:nInterim]
	totalBuffer = append(totalBuffer, slice...)

	in.Buf = make ([]byte, totalBytesOnSocket)
	copy(in.Buf, totalBuffer[:totalBytesOnSocket])
	in.BufLen = int(totalBytesOnSocket)

	msg, err := CreateMessageFromBuffer(in.Buf, 0, int(totalBytesOnSocket))
	if err != nil {
		errMsg := "TCPChannel::ReadWireMsg - unable to create a message from the input stream bytes"
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, err.Error()))
		//return nil, exception.GetErrorByType(types.TGErrorIOException, err.GetErrorCode(), err.GetErrorMsg(), errMsg)
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside TCPChannel:ReadWireMsg Created Message from buffer as '%+v'", msg.String()))
	}

	if msg.GetVerbId() == VerbExceptionMessage {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == VerbExceptionMessage"))
		errMsg := msg.(*ExceptionMessage).GetExceptionMsg()
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if msg.GetVerbId() == VerbHandShakeResponse {
		if msg.(*HandShakeResponseMessage).GetResponseStatus() == ResponseChallengeFailed {
			errMsg := msg.(*HandShakeResponseMessage).GetErrorMessage()
			logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == VerbHandShakeResponse w/ '%+v'", errMsg))
			return nil, GetErrorByType(TGErrorVersionMismatchException, "", errMsg, "")
		}
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:ReadWireMsg w/ Socket '%+v' and Message as '%+v'", obj.socket, msg.String()))
	}
	return msg, nil
}

// Send sends the message to the server, compress and or encrypt.
// Hence it is abstraction, that the channel knows about it.
// @param msg       The message that needs to be sent to the server
func (obj *TCPChannel) Send(msg tgdb.TGMessage) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering TCPChannel:Send w/ Socket '%+v' and Message as '%+v'", obj.socket, msg.String()))
	}
	if obj.output == nil || obj.GetIsClosed() {
		logger.Error(fmt.Sprint("ERROR: Returning TCPChannel::Send as the channel is closed"))
		errMsg := fmt.Sprintf("TCPChannel:Send - unable to send message to server as the channel is closed")
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	}

	err := obj.writeToWire(msg)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:Send w/ error '%+v'", err))
	}
	return err

	//// Wait for success notification from the GO routine
	//done := make(chan bool, 1)
	//// Push the message on to the channel (FIFO pipe)
	//obj.msgCh <- msg
	//
	//// TODO: Revisit later - for performance and optimization
	//// Execute sending each message content in another thread/go-routine
	////go obj.writeLoop(done)
	//// This is a common function called by both SendMessage (non-blocking) and SendRequest (blocking)
	//// Hence this should be handled called in another thread/go-routing in SendMessage ONLY
	//obj.writeLoop(done)
	//<-done
	//logger.Log(fmt.Sprintf("======> Exiting TCPChannel:Send w/ Message as '%+v'", msg.String()))
	//return nil
}

func (obj *TCPChannel) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TCPChannel:{")
	buffer.WriteString(fmt.Sprintf("IsSocketClosed: %+v", obj.isSocketClosed))
	buffer.WriteString(fmt.Sprintf(", MsgCh: %+v", obj.msgCh))
	buffer.WriteString(fmt.Sprintf(", Socket: %+v", obj.socket))
	//buffer.WriteString(fmt.Sprintf(", Input: %+v", obj.input.String()))
	//buffer.WriteString(fmt.Sprintf(", Output: %+v", obj.output.String()))
	strArray := []string{buffer.String(), obj.channelToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}


type SSLChannel struct {
	*AbstractChannel
	shutdownLock   sync.RWMutex // rw-lock for synchronizing read-n-update of env configuration
	isSocketClosed bool         // indicate if the connection is already closed
	msgCh          chan tgdb.TGMessage
	socket         *tls.Conn
	tlsConfig      *tls.Config
	input          *ProtocolDataInputStream
	output         *ProtocolDataOutputStream
}

func DefaultSSLChannel() *SSLChannel {
	newChannel := SSLChannel{
		AbstractChannel: DefaultAbstractChannel(),
		msgCh:           make(chan tgdb.TGMessage),
		isSocketClosed:  false,
	}
	buff := make([]byte, 0)
	newChannel.input = NewProtocolDataInputStream(buff)
	newChannel.output = NewProtocolDataOutputStream(0)
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	return &newChannel
}

func NewSSLChannel(linkUrl *LinkUrl, props *SortedProperties) (*SSLChannel, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("======> Entering SSLChannel:NewSSLChannel w/ linkUrl: '%s'", linkUrl.String()))
	newChannel := SSLChannel{
		AbstractChannel: NewAbstractChannel(linkUrl, props),
		msgCh:           make(chan tgdb.TGMessage),
		isSocketClosed:  false,
	}
	buff := make([]byte, 0)
	newChannel.input = NewProtocolDataInputStream(buff)
	newChannel.output = NewProtocolDataOutputStream(0)
	newChannel.exceptionCond = sync.NewCond(&newChannel.exceptionLock) // Condition for lock
	newChannel.reader = NewChannelReader(&newChannel)
	config, err := initTLSConfig(props)
	if err != nil {
		return nil, err
	}
	newChannel.tlsConfig = config
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:NewSSLChannel w/ TLSConfig: '%+v'", config))
	}
	return &newChannel, nil
}

/////////////////////////////////////////////////////////////////
// Private functions for SSLChannel
/////////////////////////////////////////////////////////////////

func initTLSConfig(props *SortedProperties) (*tls.Config, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Entering SSLChannel:initTLSConfig"))
	}

	// Load System certificate
	rootCertPool, err := x509.SystemCertPool()
	if err != nil {
		errMsg := fmt.Sprint("ERROR: Returning SSLChannel::initTLSConfig Failed to read system certificate pool")
		logger.Error(errMsg)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	}

	//sysTrustFile := fmt.Sprintf("%s%slib%ssecurity%scacerts", os.Getenv("JRE_HOME"), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:initTLSConfig about to ReadFile '%+v'", sysTrustFile))
	//pem, err := ioutil.ReadFile(sysTrustFile)
	//if err != nil {
	//	errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Failed to read client certificate authority: %s", sysTrustFile)
	//	logger.Error(errMsg)
	//	return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//}

	certPool := x509.NewCertPool()
	clientCertificates := make([]tls.Certificate, 0)

	//logger.Debug(fmt.Sprint("======> Inside SSLChannel:initTLSConfig about to add system certificate to certificate pool"))
	//if !certPool.AppendCertsFromPEM(pem) {
	//	errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Can't parse client certificate data from '%s", sysTrustFile)
	//	logger.Error(errMsg)
	//	//return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	//}

	// Load the user defined certificates.
	trustedCerts := props.GetProperty(GetConfigFromKey(TlsTrustedCertificates), "")
	if trustedCerts == "" {
		errMsg := fmt.Sprint("WARNING: Returning SSLChannel::initTLSConfig There are no user defined certificates")
		logger.Warning(errMsg)
		//return nil, nil
	} else {
		userCertificateFilePaths := strings.Split(trustedCerts, ",")
		for _, userCertFile := range userCertificateFilePaths {
			userCertData, err := ioutil.ReadFile(userCertFile)
			if err != nil {
				errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Failed to read user certificate file: %s", userCertFile)
				logger.Error(errMsg)
				return nil, GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
			}
			if !certPool.AppendCertsFromPEM(userCertData) {
				errMsg := fmt.Sprintf("ERROR: Returning SSLChannel::initTLSConfig Can't parse client certificate data from '%s", userCertFile)
				logger.Error(errMsg)
				return nil, GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
			}
			clientCertificates = append(clientCertificates)
		}
	}

	tlsConfig := &tls.Config{
		Certificates:       clientCertificates,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: false,
		Rand:               rand.Reader,
		RootCAs:            rootCertPool,
	}
	// TODO: This list may change if GO Language supports more or different - Ref. https://golang.org/pkg/crypto/tls/
	// TODO: Revisit later to find out if there is any API in GO to get this as a list.
	suites := []uint16{
		// TLS 1.0 - 1.2 cipher suites.
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		// TLS 1.3 cipher suites.
		//tls.TLS_AES_128_GCM_SHA256,
		//tls.TLS_AES_256_GCM_SHA384,
		//tls.TLS_CHACHA20_POLY1305_SHA256,
		// TLS_FALLBACK_SCSV isn't a standard cipher suite but an indicator
		// that the client is doing version fallback. See RFC 7507.
		tls.TLS_FALLBACK_SCSV,
	}
	supportedSuites := FilterSuitesById(suites)
	tlsConfig.CipherSuites = supportedSuites

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Returning SSLChannel:initTLSConfig"))
	}
	return tlsConfig, nil
}

//func (obj *SSLChannel) channelConnect() types.TGError {
//	//logger.Log(fmt.Sprint("======> Entering SSLChannel:channelConnect"))
//	if isChannelConnected(obj) {
//		logger.Log(fmt.Sprintf("======> SSLChannel::channelConnect channel is already connected"))
//		obj.setNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
//		return nil
//	}
//	if isChannelClosed(obj) || obj.GetLinkState() == types.LinkNotConnected {
//		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:channelConnect about to channelTryRepeatConnect for object '%+v'", obj.String()))
//		err := channelTryRepeatConnect(obj, false)
//		if err != nil {
//			logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::channelConnect channelTryRepeatConnect failed"))
//			return err
//		}
//		obj.SetChannelLinkState(types.LinkConnected)
//		obj.setNoOfConnections(atomic.AddInt32(&ConnectionsToChannel, 1))
//		logger.Log(fmt.Sprintf("======> Returning SSLChannel:channelConnect successfully established socket connection and now has '%d' number of connections", obj.NumOfConnections))
//	} else {
//		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::channelConnect channelTryRepeatConnect - connect called on an invalid state := '%s'", obj.GetLinkState().String()))
//		errMsg := fmt.Sprintf("======> Connect called on an invalid state := '%s'", obj.GetLinkState().String())
//		return exception.NewTGGeneralExceptionWithMsg(errMsg)
//	}
//	logger.Log(fmt.Sprintf("======> Returning SSLChannel:channelConnect having '%d' number of connections", obj.GetNoOfConnections()))
//	return nil
//}

func (obj *SSLChannel) doAuthenticate() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Entering SSLChannel:doAuthenticate"))
	}
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := CreateMessageForVerb(VerbAuthenticateRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: SSLChannel::doAuthenticate pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed"))
		return err
	}

	msgRequest.(*AuthenticateRequestMessage).SetClientId(obj.clientId)
	msgRequest.(*AuthenticateRequestMessage).SetInboxAddr(obj.inboxAddress)
	msgRequest.(*AuthenticateRequestMessage).SetUserName(obj.GetChannelUserName())
	msgRequest.(*AuthenticateRequestMessage).SetPassword(obj.GetChannelPassword())

	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:doAuthenticate about to request reply for request '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::doAuthenticate channelRequestReply failed"))
		return err
	}
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:doAuthenticate received reply as '%+v'", msgResponse.String()))
	if !msgResponse.(*AuthenticateResponseMessage).IsSuccess() {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::doAuthenticate msgResponse.(*pdu.AuthenticateResponseMessage).IsSuccess() failed"))
		return NewTGBadAuthenticationWithRealm(INTERNAL_SERVER_ERROR, TGErrorBadAuthentication, "Bad username/password combination", "", "tgdb")
	}

	obj.setChannelAuthToken(msgResponse.GetAuthToken())
	obj.setChannelSessionId(msgResponse.GetSessionId())

	cryptoDataGrapher, err := NewDataCryptoGrapher(msgResponse.GetSessionId(), msgResponse.(*AuthenticateResponseMessage).GetServerCertBuffer())
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::doAuthenticate NewDataCryptoGrapher failed w/ '%+v'", err.Error()))
		return err
	}
	obj.setDataCryptoGrapher(cryptoDataGrapher)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:doAuthenticate Successfully authenticated for user: '%s'", obj.GetChannelUserName()))
	}
	return nil
}

func (obj *SSLChannel) performHandshake(sslMode bool) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Entering SSLChannel:performHandshake"))
	}
	// Use Message Factory method to create appropriate message structure (class) based on input type
	msgRequest, err := CreateMessageForVerb(VerbHandShakeRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake pdu.CreateMessageForVerb(pdu.VerbAuthenticateRequest) failed"))
		return err
	}

	msgRequest.(*HandShakeRequestMessage).SetRequestType(InitiateRequest)

	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:performHandshake about to request reply for InitiateRequest '%+v'", msgRequest.String()))
	msgResponse, err := channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::doAuthenticate channelRequestReply failed"))
		return err
	}
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:performHandshake received reply as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*SessionForcefullyTerminatedMessage).GetKillString()
			return NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	response := msgResponse.(*HandShakeResponseMessage)
	if response.GetResponseStatus() != ResponseAcceptChallenge {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	// Validate the version specific information on the response object
	serverVersion := response.GetChallenge()
	clientVersion := GetClientVersion()
	err = obj.validateHandshakeResponseVersion(serverVersion, clientVersion)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::performHandshake validateHandshakeResponseVersion failed w/ '%+v'", err.Error()))
		return err
	}

	challenge := clientVersion.GetVersionAsLong()

	// Ignore Error Handling
	_ = msgRequest.(*HandShakeRequestMessage).UpdateSequenceAndTimeStamp(-1)
	msgRequest.(*HandShakeRequestMessage).SetRequestType(ChallengeAccepted)
	msgRequest.(*HandShakeRequestMessage).SetSslMode(sslMode)
	msgRequest.(*HandShakeRequestMessage).SetChallenge(challenge)

	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:performHandshake about to request reply for ChallengeAccepted '%+v'", msgRequest.String()))
	msgResponse, err = channelRequestReply(obj, msgRequest)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake channelRequestReply failed"))
		return err
	}
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:performHandshake received reply (2) as '%+v'", msgResponse.String()))
	if msgResponse.GetVerbId() != VerbHandShakeResponse {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake HandshakeResponse message response NOT received"))
		if msgResponse.GetVerbId() == VerbSessionForcefullyTerminated {
			errMsg := msgResponse.(*SessionForcefullyTerminatedMessage).GetKillString()
			return NewTGChannelDisconnectedWithMsg(errMsg)
		}
		errMsg := fmt.Sprintf("Expecting a HandshakeResponse message, and received: '%d'. Cannot connect to the server at: '%s'", msgResponse.GetVerbId(), obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}

	response = msgResponse.(*HandShakeResponseMessage)
	if response.GetResponseStatus() != ResponseProceedWithAuthentication {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::performHandshake response.GetResponseStatus() is NOT pdu.ResponseAcceptChallenge"))
		errMsg := fmt.Sprintf("'%s': Handshake Failed. Cannot connect to the server at: '%s'", TGDB_HNDSHKRESP_ERROR, obj.channelUrl.GetUrlAsString())
		return NewTGGeneralException(TGDB_HNDSHKRESP_ERROR, TGErrorGeneralException, errMsg, "")
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel::performHandshake Handshake w/ Remote Server is successful."))
	}
	return nil
}

func (obj *SSLChannel) setSocket(newSocket *tls.Conn) tgdb.TGError {
	obj.socket = newSocket
	//err := obj.socket.SetNoDelay(true)
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:setSocket Failed to set NoDelay flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set NoDelay flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	//}
	//
	//err = obj.socket.SetLinger(0) // <= 0 means Do not linger
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:setSocket Failed to set NoLinger flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set NoLinger flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, err.Error())
	//}

	buff := make([]byte, dataBufferSize)
	obj.input = NewProtocolDataInputStream(buff)
	obj.input.BufLen = 0
	obj.output = NewProtocolDataOutputStream(dataBufferSize)
	clientId := obj.GetProperties().GetProperty(GetConfigFromKey(ChannelClientId), "")
	obj.setChannelClientId(clientId)
	obj.setChannelInboxAddr(obj.socket.RemoteAddr().String()) //SS:TODO: Is this correct
	return nil
}

func (obj *SSLChannel) setBuffers(newSocket *tls.Conn) tgdb.TGError {
	//sendSize := obj.ChannelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelSendSize))
	//if sendSize > 0 {
	//	err := newSocket.SetWriteBuffer(sendSize*1024)
	//	if err != nil {
	//		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::setBuffers newSocket.SetWriteBuffer failed"))
	//		errMsg := fmt.Sprintf("SSLChannel:setBuffers unable to set write buffer limit to '%d'", sendSize*1024)
	//		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//	}
	//}
	//receiveSize := obj.ChannelProperties.GetPropertyAsInt(utils.GetConfigFromKey(utils.ChannelRecvSize))
	//if receiveSize > 0 {
	//	err := newSocket.SetReadBuffer(receiveSize*1024)
	//	if err != nil {
	//		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::setBuffers SetReadBuffer failed"))
	//		errMsg := fmt.Sprintf("SSLChannel:setBuffers unable to set read buffer limit to '%d'", receiveSize*1024)
	//		return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, err.Error())
	//	}
	//}
	return nil
}

func (obj *SSLChannel) tryRead() (tgdb.TGMessage, tgdb.TGError) {
	//logger.Log(fmt.Sprint("======> Entering SSLChannel:tryRead"))
	n, err := obj.input.Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::tryRead obj.input.Available() failed"))
		errMsg := "SSLChannel::tryRead there is no data available to be read"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if n <= 0 {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel::tryRead as there are no bytes to read from the wire"))
		return nil, nil
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:tryRead about to request message '%d' bytes from the wire", n))
	}
	return obj.ReadWireMsg()
}

func (obj *SSLChannel) validateHandshakeResponseVersion(sVersion int64, cVersion *TGClientVersion) tgdb.TGError {
	serverVersion := NewTGServerVersion(sVersion)
	sStrVer := serverVersion.GetVersionString()

	cStrVer := cVersion.GetVersionString()

	if 	serverVersion.GetMajor() == cVersion.GetMajor() &&
		serverVersion.GetMinor() == cVersion.GetMinor() &&
		serverVersion.GetUpdate() == cVersion.GetUpdate() {
		return nil
	}

	errMsg := fmt.Sprintf("======> Inside SSLChannel:validateHandshakeResponseVersion - Version mismatch between client(%s) & server(%s)", cStrVer, sStrVer)
	if logger.IsDebug() {
		logger.Debug(errMsg)
	}
	return GetErrorByType(TGErrorVersionMismatchException, "", errMsg, "")
}

func (obj *SSLChannel) writeLoop(done chan bool) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Entering SSLChannel:writeLoop"))
	}
	for {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("======> Inside SSLChannel:writeLoop entering infinite loop"))
		}
		select { // Non-blocking channel operation
		case msg, ok := <-obj.msgCh: // Retrieve the message from the channel
			if !ok {
				//if (gLogger.isEnabled(TGLogger.TGLevel.DebugWire)) {
				//	logMessage("SSLChannel::writeLoop unable to retrieve msg from the channel);
				//}
				logger.Error(fmt.Sprint("ERROR: Returning SSLChannel:writeLoop unable to retrieve message from obj.msgCh"))
				return
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("======> Inside SSLChannel:writeLoop retrieved message from obj.msgCh as '%+v'", msg.String()))
			}

			err := obj.writeToWire(msg)
			if err != nil {
				// TODO: Revisit later - Do something
				logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::writeLoop unable to obj.writeToWire"))
				return
			}

			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("======> Inside SSLChannel:writeLoop successfully wrote message '%+v' on the socket", msg.String()))
			}
			break
		default:
			// TODO: Revisit later - Do something
		}
	} // End of Infinite Loop
	// Send an acknowledgement of completion to the parent thread
	done <- true
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("======> Returning SSLChannel:writeLoop"))
	}
}

func (obj *SSLChannel) writeToWire(msg tgdb.TGMessage) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:writeToWire w/ Msg: '%+v'", msg.String()))
	}
	obj.DisablePing()
	msgBytes, bufLen, err := msg.ToBytes()
	if err != nil {
		errMsg := fmt.Sprintf("SSLChannel::writeToWire - unable to convert message into byte format")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, err.Error()))
		//return exception.GetErrorByType(types.TGErrorIOException, "TGErrorProtocolNotSupported", errMsg, err.GetErrorMsg())
		return err
	}

	// Clear timeout deadlines set at the time of creation of the socket
	sErr := obj.socket.SetDeadline(time.Time{})
	if sErr != nil {
		errMsg := fmt.Sprintf("SSLChannel::writeToWire - unable to clear the deadline over SSL socket")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	// Reset timeout deadlines starting from NOW!!!
	timeout := NewTGEnvironment().GetChannelConnectTimeout()

	//sErr = obj.socket.SetDeadline(time.Time{})
	sErr = obj.socket.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if sErr != nil {
		errMsg := fmt.Sprintf("SSLChannel::writeToWire - unable to reset the deadline over SSL socket")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:writeToWire about to write message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	}
	// Put the data packet on the socket for network transmission
	_, sErr = obj.socket.Write(msgBytes[0:bufLen])
	if sErr != nil {
		errMsg := fmt.Sprintf("SSLChannel::writeToWire - unable to send message bytes over SSL socket")
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return GetErrorByType(TGErrorIOException, "TGErrorProtocolNotSupported", errMsg, sErr.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:writeToWire successfully wrote message bytes on the socket as '%+v'", msgBytes[0:bufLen]))
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Helper functions for SSLChannel
/////////////////////////////////////////////////////////////////

func (obj *SSLChannel) GetIsClosed() bool {
	return obj.isSocketClosed
}

func (obj *SSLChannel) SetIsClosed(flag bool) {
	obj.isSocketClosed = flag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// ChannelLock locks the communication channel between TGDB client and server
func (obj *SSLChannel) ChannelLock() {
	obj.sendLock.Lock()
}

// ChannelUnlock unlocks the communication channel between TGDB client and server
func (obj *SSLChannel) ChannelUnlock() {
	obj.sendLock.Unlock()
}

// Connect connects the underlying channel using the URL end point
func (obj *SSLChannel) Connect() tgdb.TGError {
	return channelConnect(obj)
}

// DisablePing disables the pinging ability to the channel
func (obj *SSLChannel) DisablePing() {
	obj.needsPing = false
}

// Disconnect disconnects the channel from its URL end point
func (obj *SSLChannel) Disconnect() tgdb.TGError {
	return channelDisConnect(obj)
}

// EnablePing enables the pinging ability to the channel
func (obj *SSLChannel) EnablePing() {
	obj.needsPing = true
}

// ExceptionLock locks the communication channel between TGDB client and server in case of business exceptions
func (obj *SSLChannel) ExceptionLock() {
	obj.exceptionLock.Lock()
}

// ExceptionUnlock unlocks the communication channel between TGDB client and server in case of business exceptions
func (obj *SSLChannel) ExceptionUnlock() {
	obj.exceptionLock.Unlock()
}

// GetAuthToken gets Authorization Token
func (obj *SSLChannel) GetAuthToken() int64 {
	return obj.authToken
}

// GetClientId gets Client Name
func (obj *SSLChannel) GetClientId() string {
	return obj.clientId
}

// GetChannelURL gets the channel URL
func (obj *SSLChannel) GetChannelURL() tgdb.TGChannelUrl {
	return obj.channelUrl
}

// GetConnectionIndex gets the Connection Index
func (obj *SSLChannel) GetConnectionIndex() int {
	return obj.connectionIndex
}

// GetExceptionCondition gets the Exception Condition
func (obj *SSLChannel) GetExceptionCondition() *sync.Cond {
	return obj.exceptionCond
}

// GetLinkState gets the Link/channel State
func (obj *SSLChannel) GetLinkState() tgdb.LinkState {
	return obj.channelLinkState
}

// GetNoOfConnections gets number of connections this channel has
func (obj *SSLChannel) GetNoOfConnections() int32 {
	return obj.numOfConnections
}

// GetPrimaryURL gets the Primary URL
func (obj *SSLChannel) GetPrimaryURL() tgdb.TGChannelUrl {
	return obj.primaryUrl
}

// GetProperties gets the channel Properties
func (obj *SSLChannel) GetProperties() tgdb.TGProperties {
	return obj.channelProperties
}

// GetReader gets the channel Reader
func (obj *SSLChannel) GetReader() tgdb.TGChannelReader {
	return obj.reader
}

// GetResponses gets the channel Response Map
func (obj *SSLChannel) GetResponses() map[int64]tgdb.TGChannelResponse {
	return obj.responses
}

// GetSessionId gets Session id
func (obj *SSLChannel) GetSessionId() int64 {
	return obj.sessionId
}

// GetTracer gets the channel Tracer
func (obj *SSLChannel) GetTracer() tgdb.TGTracer {
	return obj.tracer
}

// IsChannelPingable checks whether the channel is pingable or not
func (obj *SSLChannel) IsChannelPingable() bool {
	return obj.needsPing
}

// IsClosed checks whether channel is open or closed
func (obj *SSLChannel) IsClosed() bool {
	return isChannelClosed(obj)
}

// SendMessage sends a Message on this channel, and returns immediately - An Asynchronous or Non-Blocking operation
func (obj *SSLChannel) SendMessage(msg tgdb.TGMessage) tgdb.TGError {
	return channelSendMessage(obj, msg, true)
}

// SendRequest sends a Message, waits for a response in the message format, and blocks the thread till it gets the response
func (obj *SSLChannel) SendRequest(msg tgdb.TGMessage, response tgdb.TGChannelResponse) (tgdb.TGMessage, tgdb.TGError) {
	return channelSendRequest(obj, msg, response, true)
}

// SetChannelLinkState sets the Link/channel State
func (obj *SSLChannel) SetChannelLinkState(state tgdb.LinkState) {
	obj.channelLinkState = state
}

// SetChannelURL sets the channel URL
func (obj *SSLChannel) SetChannelURL(url tgdb.TGChannelUrl) {
	obj.channelUrl = url.(*LinkUrl)
}

// SetConnectionIndex sets the connection index
func (obj *SSLChannel) SetConnectionIndex(index int) {
	obj.connectionIndex = index
}

// SetNoOfConnections sets number of connections
func (obj *SSLChannel) SetNoOfConnections(count int32) {
	obj.numOfConnections = count
}

// SetResponse sets the ChannelResponse Map
func (obj *SSLChannel) SetResponse(reqId int64, response tgdb.TGChannelResponse) {
	obj.responses[reqId] = response
}

// Start starts the channel so that it can send and receive messages
func (obj *SSLChannel) Start() tgdb.TGError {
	return channelStart(obj)
}

// Stop stops the channel forcefully or gracefully
func (obj *SSLChannel) Stop(bForcefully bool) {
	channelStop(obj, bForcefully)
}

// CreateSocket creates the physical link socket
func (obj *SSLChannel) CreateSocket() tgdb.TGError {
	/**
			super.createSocket();
	        try {
	            sslSocket = (SSLSocket)sslSocketFactory.createSocket(this.socket, getHost(), getPort(), true);
	            String[] suites = sslSocket.getEnabledCipherSuites();
	            supportedSuites = TGCipherSuite.filterSuites(suites);
	            sslSocket.setEnabledCipherSuites(supportedSuites);
	            sslSocket.setEnabledProtocols(TLSProtocols);
	            sslSocket.setUseClientMode(true);
	        }
	        catch (IOException ioe) {
	            throw new TGChannelDisconnectedException(ioe);
	        }
	*/
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:CreateSocket"))
	}
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	obj.SetChannelLinkState(tgdb.LinkNotConnected)
	host := obj.channelUrl.urlHost
	port := obj.channelUrl.urlPort
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:CreateSocket attempting to resolve address for '%s'", serverAddr))
	}

	//tcpAddr, tErr := net.ResolveSSLAddr(types.ProtocolSSL.String(), serverAddr)
	//if tErr != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket net.ResolveSSLAddr failed"))
	//	errMsg := fmt.Sprintf("SSLChannel:CreateSocket unable to resolve channel address '%s'", serverAddr)
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, tErr.Error())
	//}
	////logger.Debug(fmt.Sprintf("======> Inside SSLChannel:CreateSocket resolved SSL address for '%s' as '%+v'", serverAddr, tcpAddr))
	//
	sslConn, cErr := tls.Dial(tgdb.ProtocolSSL.String(), serverAddr, obj.tlsConfig)
	if cErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::CreateSocket Failed to connect to the server at '%s' w/ '%+v'", serverAddr, cErr.Error()))
		failureMessage := fmt.Sprintf("Failed to connect to the server at '%s'", serverAddr)
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:CreateSocket created SSL connection for '%s' as '%+v'", serverAddr, sslConn))
	}

	timeout := NewTGEnvironment().GetChannelConnectTimeout()
	dErr := sslConn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if dErr != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::CreateSocket Failed to set deadline of '%+v' seconds on the connection to the server", time.Duration(timeout)*time.Second))
		failureMessage := fmt.Sprintf("Failed to set the timeout '%d' on socket", timeout)
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, dErr.Error())
	}

	//err := sslConn.SetKeepAlive(true)
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set keep alive flag to true"))
	//	failureMessage := fmt.Sprint("Failed to set keep alive flag to true")
	//	return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
	//}

	// Set Read / Write Buffer Size on the socket
	tcErr := obj.setBuffers(sslConn)
	if tcErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set buffers"))
		return tcErr
	}
	tcErr = obj.setSocket(sslConn)
	if tcErr != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::CreateSocket Failed to set socket value to the object"))
		return tcErr
	}
	obj.SetIsClosed(false)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:CreateSocket w/ SSL Connection as '%+v'", *obj.socket))
	}
	return nil
}

// CloseSocket closes the socket
func (obj *SSLChannel) CloseSocket() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:CloseSocket w/ socket: '%+v'", obj.socket))
	}
	obj.shutdownLock.Lock()
	defer func() {
		obj.SetIsClosed(true)
		obj.shutdownLock.Unlock()
		obj.socket = nil
		obj.input = nil
		obj.output = nil
	} ()

	if obj.socket != nil {
		cErr := obj.socket.Close()
		if cErr != nil {
			failureMessage := "Failed to close the socket to the server"
			logger.Error(fmt.Sprintf("ERROR: Returning SSLChannel::CloseSocket %s w/ '%+v'", failureMessage, cErr.Error()))
			//return exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorProtocolNotSupported", failureMessage, cErr.Error())
		}
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:CloseSocket for socket: '%+v'", obj.socket))
	}
	return nil
}

// OnConnect executes all the channel specific activities
func (obj *SSLChannel) OnConnect() tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:OnConnect about to tryRead"))
	}
	msg, err := obj.tryRead()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.tryRead() failed"))
		errMsg := "SSLChannel::OnConnect there is no data available to be read"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if msg != nil {
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("======> Inside SSLChannel:OnConnect tryRead() read Message as '%+v'", msg.String()))
		}
	}

	if msg != nil && msg.GetVerbId() == VerbSessionForcefullyTerminated {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:OnConnect since Message is of Forceful Termination Type"))
		return NewTGChannelDisconnectedWithMsg(msg.(*SessionForcefullyTerminatedMessage).GetKillString())
	}

	// Check Host Verifier
	err1 := obj.socket.Handshake()
	if err1 != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.socket.Handshake() failed"))
		errMsg := "SSLChannel::OnConnect obj.socket.Handshake() failed"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, err1.Error())
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:OnConnect about to performHandshake"))
	}
	err = obj.performHandshake(true)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.performHandshake() failed"))
		errMsg := "SSLChannel::OnConnect error in performing handshake with server"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:OnConnect about to doAuthenticate"))
	}
	err = obj.doAuthenticate()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::OnConnect obj.doAuthenticate() failed"))
		errMsg := "SSLChannel::OnConnect error in authentication with server"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:OnConnect"))
	}
	return nil
}

// ReadWireMsg reads the message from the wire in the form of byte stream
func (obj *SSLChannel) ReadWireMsg() (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:ReadWireMsg w/ SSLChannel as '%+v'", obj.String()))
	}
	obj.input.BufLen = dataBufferSize
	in := obj.input
	if in == nil {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:ReadWireMsg since obj.input is NIL"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}

	obj.DisablePing()
	if obj.GetIsClosed() {
		logger.Warning(fmt.Sprint("WARNING: Returning SSLChannel:ReadWireMsg since SSL channel is Closed"))
		// TODO: Revisit later - Should we not return an error?
		return nil, nil
	}
	obj.lastActiveTime = time.Now()

	// Read the data on the socket
	buff := make([]byte, dataBufferSize)
	n, sErr := obj.socket.Read(buff)
	if sErr != nil || n <= 0 {
		errMsg := "SSLChannel::ReadWireMsg obj.socket.Read failed"
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, sErr.Error()))
		return nil, GetErrorByType(TGErrorIOException, "", errMsg, sErr.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Read '%d' bytes from the wire in buff '%+v'", n, buff[:(2*n)]))
	}
	copy(in.Buf, buff[:n])
	in.BufLen = n
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Input Stream Buffer('%d') is '%+v'", in.BufLen, in.Buf[:(2*n)]))

	// Needed to avoid dirty data in the buffer when we handle the message
	buffer := make([]byte, n)
	//logger.Debug(fmt.Sprint("======> Inside SSLChannel:ReadWireMsg in.ReadFullyAtPos read msgBytes as '%+v'", msgBytes))
	copy(buffer, buff[:n])
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg copied into buffer as '%+v'", buffer))

	//intToBytes(size, msgBytes, 0)
	//bytesRead, _ := utils.FormatHex(msgBytes)
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg bytes read: '%s'", bytesRead))

	msg, err := CreateMessageFromBuffer(buffer, 0, n)
	if err != nil {
		errMsg := "SSLChannel::ReadWireMsg - unable to create a message from the input stream bytes"
		logger.Error(fmt.Sprintf("ERROR: Returning %s w/ '%+v'", errMsg, err.Error()))
		//return nil, exception.GetErrorByType(types.TGErrorIOException, "", errMsg, err.GetErrorMsg())
		return nil, err
	}
	//logger.Debug(fmt.Sprintf("======> Inside SSLChannel:ReadWireMsg Created Message from buffer as '%+v'", msg.String()))

	if msg.GetVerbId() == VerbExceptionMessage {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbExceptionMessage"))
		errMsg := msg.(*ExceptionMessage).GetExceptionMsg()
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}

	//if msg.GetVerbId() == pdu.VerbHandShakeResponse {
	//	if msg.GetResponseStatus() == pdu.ResponseChallengeFailed {
	//		errMsg := msg.GetErrorMessage()
	//		logger.Error(fmt.Sprintf("ERROR: Returning TCPChannel::ReadWireMsg msg.GetVerbId() == pdu.VerbHandShakeResponse w/ '%+v'", errMsg))
	//		return nil, exception.GetErrorByType(types.TGErrorVersionMismatchException, "", errMsg, "")
	//	}
	//}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning SSLChannel:ReadWireMsg w/ Message as '%+v'", msg.String()))
	}
	return msg, nil
}

// Send Message to the server, compress and/or encrypt.
// Hence it is abstraction, that the channel knows about it.
// @param msg       The message that needs to be sent to the server
func (obj *SSLChannel) Send(msg tgdb.TGMessage) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Entering SSLChannel:Send w/ Message as '%+v'", msg.String()))
	}
	if obj.output == nil || obj.GetIsClosed() {
		logger.Error(fmt.Sprint("ERROR: Returning SSLChannel::Send as the channel is closed"))
		errMsg := fmt.Sprintf("SSLChannel:Send - unable to send message to server as the channel is closed")
		return GetErrorByType(TGErrorGeneralException, "TGErrorProtocolNotSupported", errMsg, "")
	}

	err := obj.writeToWire(msg)
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("======> Returning TCPChannel:Send w/ error '%+v'", err))
	}
	return err
}

func (obj *SSLChannel) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("SSLChannel:{")
	buffer.WriteString(fmt.Sprintf("IsSocketClosed: %+v", obj.isSocketClosed))
	buffer.WriteString(fmt.Sprintf(", MsgCh: %+v", obj.msgCh))
	buffer.WriteString(fmt.Sprintf(", Socket: %+v", obj.socket))
	//buffer.WriteString(fmt.Sprintf("Input: %+v ", obj.input.String()))
	//buffer.WriteString(fmt.Sprintf("Output: %+v ", obj.output.String()))
	//buffer.WriteString("\n")
	strArray := []string{buffer.String(), obj.channelToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}


//var logger = logging.DefaultTGLogManager().GetLogger()

type TGChannelFactory struct {
}

var globalChannelFactory *TGChannelFactory
var cfOnce sync.Once

func newTGChannelFactory() *TGChannelFactory {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGChannelFactory{})

	cfOnce.Do(func() {
		globalChannelFactory = &TGChannelFactory{}
	})
	return globalChannelFactory
}

// Get an instance of the Channel Factory
func GetChannelFactoryInstance() *TGChannelFactory {
	return newTGChannelFactory()
}

/////////////////////////////////////////////////////////////////
// Private functions for TGChannelFactory
/////////////////////////////////////////////////////////////////

func (obj *TGChannelFactory) createChannelWithProperties(urlPath, userName, password string, props map[string]string) (tgdb.TGChannel, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering TGChannelFactory:createChannelWithProperties w/ URL: '%s' User: '%s', Pwd: '%s'", urlPath, userName, password))
	if len(urlPath) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - urlPath is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid URL specified as '%s'", urlPath)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	if len(userName) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - userName is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid user specified as '%s'", userName)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	if len(password) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - password is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid password specified as '%s'", password)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	properties := NewSortedProperties()
	if props != nil {
		for k, v := range props {
			properties.AddProperty(k, v)
		}
	}
	channelUrl := ParseChannelUrl(urlPath)
	if channelUrl != nil {
		urlProps := channelUrl.GetProperties().(*SortedProperties)
		for _, kvp := range urlProps.GetAllProperties() {
			properties.AddProperty(kvp.KeyName, kvp.KeyValue)
		}
	}
	err1 := SetUserAndPassword(properties, userName, password)
	if err1 != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGChannelFactory:createChannelWithProperties - unable to set user and password in the property set w/ Error: '%+v'", err1.Error()))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties unable to set user '%s' and password in the property set", userName)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorGeneralException", errMsg, err1.Error())
	}
	return obj.CreateChannelWithUrlProperties(channelUrl, properties)
}

func (obj *TGChannelFactory) CreateChannelWithUrlProperties(channelUrl tgdb.TGChannelUrl, props *SortedProperties) (tgdb.TGChannel, tgdb.TGError) {
	//logger.Log(fmt.Sprintf("Entering TGChannelFactory:CreateChannelWithUrlProperties w/ ChannelURL: '%+v' and Properties: '%+v'", channelUrl, props))
	if channelUrl == nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:CreateChannelWithUrlProperties - channelUrl is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:CreateChannelWithUrlProperties Invalid URL specified as '%+v'", channelUrl)
		return nil, GetErrorByType(TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	channelProtocol := channelUrl.GetProtocol()
	switch channelProtocol {
	case tgdb.ProtocolTCP:
		return NewTCPChannel(channelUrl.(*LinkUrl), props), nil
	case tgdb.ProtocolSSL:
		return NewSSLChannel(channelUrl.(*LinkUrl), props)
	case tgdb.ProtocolHTTP:
		fallthrough
		//return NewHTTPChannel(channelUrl.(*LinkUrl), props), nil
	case tgdb.ProtocolHTTPS:
		fallthrough
		//return NewHTTPSChannel(channelUrl.(*LinkUrl), props), nil
	default:
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithUrlProperties protocol '%s' not supported", channelProtocol.String())
		return nil, GetErrorByType(TGErrorProtocolNotSupported, "TGErrorProtocolNotSupported", errMsg, "")
	}
	return nil, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// Create a channel on the URL specified using the userName and password.
// A URL is represented as a string of the form
//         <protocol>://[user@]['['ipv6']'] | ipv4 [:][port][/]'{' Name:value;... '}'
// @param urlPath A url string.
// @param userName The userName for the channel. The userId provided overrides all other userIds that can be inferred.
//         The rules for overriding are in this order
//         a. The argument 'userId' is the highest priority. If Null then
//         b. The user@url is considered. If that is Null
//         c. the "userID=value" from the URL string is considered.
//         d. If all of them is Null, then the default User associated to the installation will be taken.
// @param password An encrypted password associated with the userName
// @return a Channel
func (obj *TGChannelFactory) CreateChannel(urlPath, userName, password string) (tgdb.TGChannel, tgdb.TGError) {
	props := make(map[string]string, 0)
	return obj.CreateChannelWithProperties(urlPath, userName, password, props)
}

// Create a channel on the URL specified using the user Name and password
// @param urlPath A url as a string form
// @param userName The userName for the channel. The userId provided overrides all other userIds that can be infered.
//               The rules for overriding are in this order
//               a. The argument 'userId' is the highest priority. If Null then
//               b. The user@url is considered. If that is Null
//               c. the "userID=value" from the URL string is considered.
//               d. The user retrieved from the Properties is considered
//               e. If all of them is Null, then the default User associated to the installation will be taken.
// @param password Encrypted password
// @param props A properties bag with Connection Properties. The URL inferred properties override this property bag.
// @return a connected channel
func (obj *TGChannelFactory) CreateChannelWithProperties(urlPath, userName, password string, props map[string]string) (tgdb.TGChannel, tgdb.TGError) {
	return obj.createChannelWithProperties(urlPath, userName, password, props)
}



type TGCipherSuite struct {
	suiteId     uint16
	opensslName string
	keyExch     string
	encryption  string
	bits        string
}

var PreDefinedCipherSuites = map[string]TGCipherSuite{
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         {tls.TLS_RSA_WITH_AES_128_CBC_SHA256, "AES128-SHA256", "RSA", "AES", "128"},
	"TLS_RSA_WITH_AES_256_CBC_SHA256":         {0x3d, "AES256-SHA256", "RSA", "AES", "256"},
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256":     {0x67, "DHE-RSA-AES128-SHA256", "DH", "AES", "128"},
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256":     {0x6b, "DHE-RSA-AES256-SHA256", "DH", "AES", "256"},
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         {tls.TLS_RSA_WITH_AES_128_GCM_SHA256, "AES128-GCM-SHA256", "RSA", "AESGCM", "128"},
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         {tls.TLS_RSA_WITH_AES_256_GCM_SHA384, "AES256-GCM-SHA384", "RSA", "AESGCM", "256"},
	"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256":     {0x9e, "DHE-RSA-AES128-GCM-SHA256", "DH", "AESGCM", "128"},
	"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384":     {0x9f, "DHE-RSA-AES256-GCM-SHA384", "DH", "AESGCM", "256"},
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": {tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, "ECDHE-ECDSA-AES128-SHA256", "ECDH", "AES", "128"},
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384": {0xc024, "ECDHE-ECDSA-AES256-SHA384", "ECDH", "AES", "256"},
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   {tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, "ECDHE-RSA-AES128-SHA256", "ECDH", "AES", "128"},
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384":   {0xc028, "ECDHE-RSA-AES256-SHA384", "ECDH", "AES", "256"},
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": {tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, "ECDHE-ECDSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128"},
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": {tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, "ECDHE-ECDSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256"},
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   {tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, "ECDHE-RSA-AES128-GCM-SHA256", "ECDH", "AESGCM", "128"},
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   {tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, "ECDHE-RSA-AES256-GCM-SHA384", "ECDH", "AESGCM", "256"},
	"TLS_INVALID_CIPHER":                      {0, "", "", "", ""},
}

func NewCipherSuite(id uint16, name, key, encr, bitSize string) *TGCipherSuite {
	return &TGCipherSuite{suiteId: id, opensslName: name, keyExch: key, encryption: encr, bits: bitSize}
}

/////////////////////////////////////////////////////////////////
// Helper Public functions for TGCipherSuite
/////////////////////////////////////////////////////////////////

// GetCipherSuite returns the TGCipherSuite given its full qualified string form or its alias Name.
func GetCipherSuite(nameOrAlias string) *TGCipherSuite {
	for name, suite := range PreDefinedCipherSuites {
		if 	strings.ToLower(name) == strings.ToLower(nameOrAlias) ||
			strings.ToLower(suite.opensslName) == strings.ToLower(nameOrAlias) {
			return &suite
		}
	}
	invalid := PreDefinedCipherSuites["TLS_INVALID_CIPHER"]
	return &invalid
}

// GetCipherSuiteById returns the TGCipherSuite given its ID.
func GetCipherSuiteById(id uint16) *TGCipherSuite {
	for _, suite := range PreDefinedCipherSuites {
		if 	suite.suiteId == id {
			return &suite
		}
	}
	invalid := PreDefinedCipherSuites["TLS_INVALID_CIPHER"]
	return &invalid
}

// FilterSuites returns CipherSuites that are supported by TGDB client
func FilterSuites(suites []string) []string {
	supportedSuites := make([]string, 0)
	for _, inputSuite := range suites {
		cs := GetCipherSuite(inputSuite)
		// Ignore "TLS_INVALID_CIPHER"
		if cs.suiteId != 0 {
			supportedSuites = append(supportedSuites, cs.opensslName)
		}
	}
	return supportedSuites
}

// FilterSuitesById returns CipherSuites that are supported by TGDB client
func FilterSuitesById(suites []uint16) []uint16 {
	supportedSuites := make([]uint16, 0)
	for _, inputSuite := range suites {
		cs := GetCipherSuiteById(inputSuite)
		// Ignore "TLS_INVALID_CIPHER"
		if cs.suiteId != 0 {
			supportedSuites = append(supportedSuites, cs.suiteId)
		}
	}
	return supportedSuites
}


