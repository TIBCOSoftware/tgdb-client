package channel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"strconv"
	"strings"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
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
 * File name: TGLinkUrl.go
 * Created on: Dec 01, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// Default settings - A fallback option - to construct a valid connection URL to TGDB server
const (
	gDefaultHost string = "localhost"
	gDefaultPort int    = 8222
)

type LinkUrl struct {
	ftUrls   []types.TGChannelUrl
	urlHost  string
	isIPv6   bool
	urlPort  int
	protocol types.TGProtocol
	urlStr   string
	urlProps *utils.SortedProperties // This is always sorted
	urlUser  string
}

func DefaultLinkUrl() *LinkUrl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(LinkUrl{})

	// Instantiate a single instance collection of all environment variable settings
	tgEnv := utils.NewTGEnvironment()
	newChannelUrl := LinkUrl{
		ftUrls:   make([]types.TGChannelUrl, 0),
		urlHost:  tgEnv.GetChannelDefaultHost(),
		isIPv6:   false,
		urlPort:  tgEnv.GetChannelDefaultPort(),
		protocol: types.ProtocolTCP,
		urlProps: &utils.SortedProperties{},
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
		newChannelUrl.urlProps = props.(*utils.SortedProperties)
	}
	newChannelUrl.urlUser = user
	//logger.Log(fmt.Sprintf("Configuration name for ChannelUserID is '%+v'", utils.GetConfigFromKey(utils.ChannelUserID)))
	newChannelUrl.urlProps.AddProperty(utils.GetConfigFromKey(utils.ChannelUserID).GetName(), user)
	return newChannelUrl
}

func NewLinkUrlWithComponents(proto types.TGProtocol, host string, port int) *LinkUrl {
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
func parseUrlComponents(sUrl string) (types.TGProtocol, string, int, bool, string, []types.TGChannelUrl, types.TGProperties, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseUrlComponents w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	protocol := types.ProtocolTCP
	var user, host string
	var port int
	var ip6Flag bool
	ftUrls := make([]types.TGChannelUrl, 0)

	if len(sUrl) == 0 {
		errMsg := fmt.Sprintf("Returning LinkUrl:parseUrlComponents as invalid length for input URL string '%s'", sUrl)
		err := exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	protocol, urlSubstr, err := parseProtocol(sUrl)
	if err != nil {
		logger.Log(fmt.Sprintf("Returning LinkUrl:parseUrlComponents - error in parsing protocol '%+v", err.Error()))
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	if len(urlSubstr) == 0 {
		errMsg := fmt.Sprintf("Returning LinkUrl:parseUrlComponents as incorrect Host/Port specified in input URL string '%s'", sUrl)
		err := exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	user, host, port, ip6Flag, propSubstr, err := parseUserAndHostAndPort(urlSubstr)
	if err != nil {
		return protocol, host, port, ip6Flag, user, nil, nil, err
	}
	if len(propSubstr) == 0 {
		logger.Log(fmt.Sprint("Returning LinkUrl:parseUrlComponents as there are no properties specified - hence no FTUrls mentioned as part of property set"))
		return protocol, host, port, ip6Flag, user, nil, nil, nil
	}
	props, ftUrls, err3 := parseProperties(propSubstr)
	if err3 != nil {
		logger.Log(fmt.Sprintf("Returning LinkUrl:parseUrlComponents - error in parsing properties '%+v", err.Error()))
		return protocol, host, port, ip6Flag, user, ftUrls, nil, err3
	}
	logger.Log(fmt.Sprintf("Returning LinkUrl:parseUrlComponents w/ '%+v' '%+v' '%+v' '%+v' '%+v' '%+v' '%+v'", protocol, host, port, ip6Flag, user, ftUrls, props))
	return protocol, host, port, ip6Flag, user, ftUrls, props, nil
}

func parseProtocol(sUrl string) (types.TGProtocol, string, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseProtocol w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	var protocolStr, urlSubstring string
	protocol := types.ProtocolTCP // Default
	if len(sUrl) == 0 {
		errMsg := "Unable to parse protocol from the channel URL string"
		return protocol, "", exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Expected format of sUrl is one of the following
	// (a) http://scott@foo.bar.com:8700/... 	                    <== protocol://User/DNS/Port/...
	// (b) https://scott@foo.bar.com/...	                        <== protocol://User/DNS/...
	// (c) tcp://scott@10.20.30.40:8123/...	                        <== protocol://User/IPv4/Port/...
	// (d) ssl://scott@10.20.30.40/...	                            <== protocol://User/IPv4//...

	//logger.Log(fmt.Sprintf("ParseProtocol has incoming URL string as: '%+v'", sUrl))
	idx := strings.IndexAny(sUrl, "://")
	if idx > 0 {
		strComponents := strings.Split(sUrl, "://")
		protocolStr = strings.ToLower(strComponents[0])
		urlSubstring = strComponents[1]
		switch protocolStr {
		case types.ProtocolHTTP.String():
			protocol = types.ProtocolHTTP
		case types.ProtocolHTTPS.String():
			protocol = types.ProtocolHTTPS
		case types.ProtocolSSL.String():
			protocol = types.ProtocolSSL
		case types.ProtocolTCP.String():
			fallthrough
		default:
			protocol = types.ProtocolTCP
		}
	} else {
		protocol = types.ProtocolTCP
		urlSubstring = sUrl
	}

	//logger.Log(fmt.Sprintf("Returning LinkUrl:parseProtocol w/ URL string as '%+v' '%+v'", protocol, urlSubstring))
	return protocol, urlSubstring, nil
}

func parseUserAndHostAndPort(sUrl string) (string, string, int, bool, string, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseUserAndHostAndPort w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	var urlSubstring, user, host string
	var port int
	var ip6Flag bool
	var hostStr, hostPortStr, userHostPortStr string

	if len(sUrl) == 0 {
		errMsg := "Unable to parse user/host/port from the channel URL string"
		return user, host, port, ip6Flag, "", exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
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

	//logger.Log(fmt.Sprintf("ParseUserAndHostAndPort has incoming URL string as: '%+v'", sUrl))
	if strings.Contains(sUrl, "/") {
		strComponents := strings.Split(sUrl, "/")
		urlSubstring = strComponents[1] // Remaining string that may have properties in N1'V1,N2=V2,... format
		userHostPortStr = strComponents[0]
	} else {
		// At this point, there is no trailing '/' indicating it is probably one of the ftUrls or no properties
		userHostPortStr = sUrl
	}

	//logger.Log(fmt.Sprintf("ParseUserAndHostAndPort userHostPortStr string as: '%+v'", userHostPortStr))
	// The following if block takes care of formats a through f - to extract user as part of host-port string
	if strings.Contains(userHostPortStr, "@") {
		userHostPortComps := strings.Split(userHostPortStr, "@")
		user = userHostPortComps[0]
		hostPortStr = userHostPortComps[1]
		if len(hostPortStr) == 0 {
			// This is the case for 'scott@/...'
			errMsg := "Invalid or missing host name and/or port in the channel URL string"
			return user, host, port, ip6Flag, "", exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
		}
	} else {
		// At this point, user is still NOT specified as shown in formats g through l above
		hostPortStr = userHostPortStr
	}
	//logger.Log(fmt.Sprintf("ParseUserAndHostAndPort hostPortStr string as: '%+v'", hostPortStr))

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
				return user, host, port, ip6Flag, "", exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
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

func parseProperties(sUrl string) (types.TGProperties, []types.TGChannelUrl, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering LinkUrl:parseProperties w/ URL string as '%s'", sUrl))
	// Intentionally do not set the following components to default values
	if len(sUrl) == 0 {
		errMsg := "Unable to parse properties from the channel URL string"
		return nil, nil, exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// Expected format of sUrl is one of the following
	// {userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}

	// Intentionally do not set the following components to default values
	props := utils.NewSortedProperties()
	ftUrls := make([]types.TGChannelUrl, 0)

	if !strings.HasPrefix(sUrl, "{") {
		errMsg := "Unable to parse properties from the channel URL string"
		return nil, nil, exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	//logger.Log(fmt.Sprintf("ParseProperties has incoming URL string as: '%+v'", sUrl))
	sUrl = sUrl[1 : len(sUrl)-1] // Strip '{' and '}' from both ends
	nvPairs := strings.Split(sUrl, ";")
	if len(nvPairs) == 0 {
		errMsg := "Malformed URL property specification - Must begin with { and end with }. All key=value must be separated with ;"
		return nil, nil, exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	// At this point, nvPairs is expected to have each entry as 'name=value' string
	for _, nvPair := range nvPairs {
		kvp := strings.Split(nvPair, "=")
		props.AddProperty(kvp[0], kvp[1])
	}
	//logger.Log(fmt.Sprintf("Url Properties parsed from input URL are: '%+v'", props))

	// Once all the properties are scanned and parsed, check whether any of the supplied properties was ftHosts
	if utils.DoesPropertyExist(props, "ftHosts") {
		//logger.Log(fmt.Sprintf("Url Properties does have an existing property named / aliased as 'ftHosts'"))
		cn := utils.GetConfigFromName("ftHosts")
		if cn == nil || cn.GetName() == "" {
			ftUrls = nil
		} else {
			ftUrlStr := props.GetProperty(cn, "")
			ftUs := ftUrlStr
			//logger.Log(fmt.Sprintf("ftus: '%+v'", ftus))
			if len(ftUs) == 0 {
				ftUrls = nil
			} else {
				fts := strings.Split(ftUs, ",")
				for _, fUrl := range fts {
					if fUrl != "" {
						//logger.Log(fmt.Sprintf("ParseProperties trying to create new LinkURL for FtUrl '%+v'", fUrl))
						newLinkUrl := NewLinkUrl(fUrl)
						//logger.Log(fmt.Sprintf("ParseProperties created new LinkURL '%+v' for FtUrl '%+v'", newLinkUrl, fUrl))
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
func (obj *LinkUrl) GetFTUrls() []types.TGChannelUrl {
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
func (obj *LinkUrl) GetProperties() types.TGProperties {
	return obj.urlProps
}

// GetProtocol gets the protocol used as part of the URL
func (obj *LinkUrl) GetProtocol() types.TGProtocol {
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
	logger.Log(fmt.Sprintf("Returning LinkUrl:GetUrlAsString w/ URL string as '%s'", obj.urlStr))
	return obj.urlStr
}

// GetUser gets the user associated with the URL
func (obj *LinkUrl) GetUser() string {
	if obj.urlUser != "" {
		return obj.urlUser
	}
	user := ""
	userConfig := utils.GetConfigFromKey(utils.ChannelUserID)
	if userConfig != nil {
		user = obj.urlProps.GetProperty(userConfig, "")
		if user == "" {
			env := utils.NewTGEnvironment()
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
