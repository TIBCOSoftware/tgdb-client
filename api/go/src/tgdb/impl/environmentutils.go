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
 * File Name: environmentutils.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: environmentutils.go 3626 2019-12-09 19:35:03Z nimish $
 */

package impl

import (
	"os"
	"tgdb"
	"strconv"
	"strings"
)

type TGEnvironment struct {
	envMap map[ConfigName]string
	//rwMutex sync.RWMutex // Intentionally kept Private - rw-lock for synchronizing read-n-update of env configuration
}

// Global Instance
//var gInstance *TGEnvironment
//var once sync.Once

func defaultTGEnvironment() *TGEnvironment {
	gInstance := &TGEnvironment{}
	//once.Do(func() {	// Commented for unit tests
	envSuperSet := make(map[ConfigName]string, 0)
	// Load new environment set from pre-defined system configurations
	for _, config := range PreDefinedConfigurations {
		envSuperSet[config] = config.GetDefaultValue()
	}
	// Append new environment set from O/S level environment settings
	sysEnv := os.Environ()
	for _, kvPair := range sysEnv {
		s := strings.Split(kvPair, "=")
		newConfigKey := NewConfigName(s[0], s[0], s[1])
		envSuperSet[*newConfigKey] = s[1]
	}
	gInstance.envMap = envSuperSet
	//})
	return gInstance
}

func NewTGEnvironment() *TGEnvironment {
	return defaultTGEnvironment()
}

/////////////////////////////////////////////////////////////////
// Helper functions for TGEnvironment
/////////////////////////////////////////////////////////////////

func GetConfig(name string) *ConfigName {
	//logger.Log(fmt.Sprintf("GetConfig received lookup Name as '%+v'\n", Name))
	cn := GetConfigFromName(name)
	//logger.Log(fmt.Sprintf("GetConfig found config as '%+v'\n", cn))
	if cn == nil || cn.GetName() == "" {
		return nil
	}
	return cn
}

/////////////////////////////////////////////////////////////////
// Implement functions for TGEnvironment
/////////////////////////////////////////////////////////////////

// TODO: Revisit later - for more testing
func (obj *TGEnvironment) GetAsSortedProperties() tgdb.TGProperties {
	sp := NewSortedProperties()
	for cn := range obj.envMap {
		name := cn.GetName()
		value := cn.GetDefaultValue()
		if name == "" {
			continue
		}
		sp.AddProperty(name, value)
	}
	return sp
}

func (obj *TGEnvironment) GetChannelClientId() string {
	cn := GetConfigFromKey(ChannelClientId)
	if cn == nil || cn.GetName() == "" {
		return "tgdb-client"
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelClientId Value is '%+v'\n", value))
	if value == "" {
		return "tgdb-client"
	}
	return value
}

func (obj *TGEnvironment) GetChannelConnectTimeout() int {
	cn := GetConfigFromKey(ChannelConnectTimeout)
	if cn == nil || cn.GetName() == "" {
		return -1
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelConnectTimeout Value is '%+v'\n", value))
	if value == "" {
		return -1
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetChannelDefaultHost() string {
	cn := GetConfigFromKey(ChannelDefaultHost)
	if cn == nil || cn.GetName() == "" {
		return "localhost"
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelDefaultHost Value is '%+v'\n", value))
	if value == "" {
		return "localhost"
	}
	return value
}

func (obj *TGEnvironment) GetChannelDefaultPort() int {
	cn := GetConfigFromKey(ChannelDefaultPort)
	if cn == nil || cn.GetName() == "" {
		return -1
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelDefaultPort Value is '%+v'\n", value))
	if value == "" {
		return -1
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetChannelDefaultUser() string {
	cn := GetConfigFromKey(ChannelDefaultUserID)
	if cn == nil || cn.GetName() == "" {
		return ""
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelDefaultUser Value is '%+v'\n", value))
	if value == "" {
		return ""
	}
	return value
}

func (obj *TGEnvironment) GetChannelFTHosts() string {
	cn := GetConfigFromKey(ChannelFTHosts)
	if cn == nil || cn.GetName() == "" {
		return ""
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelFTHosts Value is '%+v'\n", value))
	if value == "" {
		return ""
	}
	return value
}

func (obj *TGEnvironment) GetChannelPingInterval() int {
	cn := GetConfigFromKey(ChannelPingInterval)
	if cn == nil || cn.GetName() == "" {
		return -1
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelPingInterval Value is '%+v'\n", value))
	if value == "" {
		return -1
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetChannelSendSize() int {
	cn := GetConfigFromKey(ChannelSendSize)
	if cn == nil || cn.GetName() == "" {
		return 122
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelSendSize Value is '%+v'\n", value))
	if value == "" {
		return 122
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetChannelReceiveSize() int {
	cn := GetConfigFromKey(ChannelRecvSize)
	if cn == nil || cn.GetName() == "" {
		return 128
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelReceiveSize Value is '%+v'\n", value))
	if value == "" {
		return 128
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetChannelUser() string {
	cn := GetConfigFromKey(ChannelUserID)
	if cn == nil || cn.GetName() == "" {
		return ""
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetChannelUser Value is '%+v'\n", value))
	if value == "" {
		return ""
	}
	return value
}

func (obj *TGEnvironment) GetConnectionPoolDefaultPoolSize() int {
	cn := GetConfigFromKey(ConnectionPoolDefaultPoolSize)
	if cn == nil || cn.GetName() == "" {
		return -1
	}
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetConnectionPoolDefaultPoolSize Value is '%+v'\n", value))
	if value == "" {
		return -1
	}
	v, _ := strconv.Atoi(value)
	return v
}

func (obj *TGEnvironment) GetDefaultDateTimeFormat() string {
	return "01-02-2006 15:04:05" // In Go language, 01 ==> mm, 02 ==> dd, 15|03 ==> hh, 04 ==> MM, 05 ==> ss, 06 ==> yy
}

func (obj *TGEnvironment) GetEnvironmentProperty(name string) interface{} {
	cn := GetConfig(name)
	if cn == nil || cn.GetName() == "" {
		return ""
	}
	//logger.Log(fmt.Sprintf("GetEnvironmentProperty has properties as '%+v' - '%+v'", len(obj.TGEnv), obj.TGEnv))
	value := obj.envMap[*cn]
	//logger.Log(fmt.Sprintf("GetEnvironmentProperty Value is '%+v'\n", value))
	return value
}

func (obj *TGEnvironment) SetEnvironmentProperty(name string, value string) {
	// Check whether the configuration is already loaded in environment variable set
	cn := GetConfig(name)
	if cn == nil || cn.GetName() == "" {
		return
	}
	// Set only if the configuration is already present
	obj.envMap[*cn] = value
}
