package channel

import (
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"sync"
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
 * File name: TGChannelFactory.go
 * Created on: Dec 01, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

var logger = logging.DefaultTGLogManager().GetLogger()

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

func (obj *TGChannelFactory) createChannelWithProperties(urlPath, userName, password string, props map[string]string) (types.TGChannel, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering TGChannelFactory:createChannelWithProperties w/ URL: '%s' User: '%s', Pwd: '%s'", urlPath, userName, password))
	if len(urlPath) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - urlPath is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid URL specified as '%s'", urlPath)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	if len(userName) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - userName is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid user specified as '%s'", userName)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	if len(password) == 0 {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:createChannelWithProperties - password is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties Invalid password specified as '%s'", password)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	properties := utils.NewSortedProperties()
	if props != nil {
		for k, v := range props {
			properties.AddProperty(k, v)
		}
	}
	channelUrl := ParseChannelUrl(urlPath)
	if channelUrl != nil {
		urlProps := channelUrl.GetProperties().(*utils.SortedProperties)
		for _, kvp := range urlProps.GetAllProperties() {
			properties.AddProperty(kvp.KeyName, kvp.KeyValue)
		}
	}
	err1 := utils.SetUserAndPassword(properties, userName, password)
	if err1 != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGChannelFactory:createChannelWithProperties - unable to set user and password in the property set w/ Error: '%+v'", err1.Error()))
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithProperties unable to set user '%s' and password in the property set", userName)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorGeneralException", errMsg, err1.Error())
	}
	return obj.CreateChannelWithUrlProperties(channelUrl, properties)
}

func (obj *TGChannelFactory) CreateChannelWithUrlProperties(channelUrl types.TGChannelUrl, props *utils.SortedProperties) (types.TGChannel, types.TGError) {
	//logger.Log(fmt.Sprintf("Entering TGChannelFactory:CreateChannelWithUrlProperties w/ ChannelURL: '%+v' and Properties: '%+v'", channelUrl, props))
	if channelUrl == nil {
		logger.Error(fmt.Sprint("ERROR: Returning TGChannelFactory:CreateChannelWithUrlProperties - channelUrl is EMPTY"))
		errMsg := fmt.Sprintf("TGChannelFactory:CreateChannelWithUrlProperties Invalid URL specified as '%+v'", channelUrl)
		return nil, exception.GetErrorByType(types.TGErrorGeneralException, "TGErrorGeneralException", errMsg, "")
	}
	channelProtocol := channelUrl.GetProtocol()
	switch channelProtocol {
	case types.ProtocolTCP:
		return NewTCPChannel(channelUrl.(*LinkUrl), props), nil
	case types.ProtocolSSL:
		return NewSSLChannel(channelUrl.(*LinkUrl), props)
	case types.ProtocolHTTP:
		fallthrough
		//return NewHTTPChannel(channelUrl.(*LinkUrl), props), nil
	case types.ProtocolHTTPS:
		fallthrough
		//return NewHTTPSChannel(channelUrl.(*LinkUrl), props), nil
	default:
		errMsg := fmt.Sprintf("TGChannelFactory:createChannelWithUrlProperties protocol '%s' not supported", channelProtocol.String())
		return nil, exception.GetErrorByType(types.TGErrorProtocolNotSupported, "TGErrorProtocolNotSupported", errMsg, "")
	}
	return nil, nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGChannel
/////////////////////////////////////////////////////////////////

// Create a channel on the URL specified using the userName and password.
// A URL is represented as a string of the form
//         <protocol>://[user@]['['ipv6']'] | ipv4 [:][port][/]'{' name:value;... '}'
// @param urlPath A url string.
// @param userName The userName for the channel. The userId provided overrides all other userIds that can be inferred.
//         The rules for overriding are in this order
//         a. The argument 'userId' is the highest priority. If Null then
//         b. The user@url is considered. If that is Null
//         c. the "userID=value" from the URL string is considered.
//         d. If all of them is Null, then the default User associated to the installation will be taken.
// @param password An encrypted password associated with the userName
// @return a Channel
func (obj *TGChannelFactory) CreateChannel(urlPath, userName, password string) (types.TGChannel, types.TGError) {
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
func (obj *TGChannelFactory) CreateChannelWithProperties(urlPath, userName, password string, props map[string]string) (types.TGChannel, types.TGError) {
	return obj.createChannelWithProperties(urlPath, userName, password, props)
}
