package utils

import (
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGConfigName.go
 * Created on: Sep 23, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	ChannelDefaultHost = iota
	ChannelDefaultPort
	ChannelDefaultProtocol
	ChannelSendSize
	ChannelRecvSize
	ChannelPingInterval
	ChannelConnectTimeout
	ChannelFTHosts
	ChannelFTRetryIntervalSeconds
	ChannelFTRetryCount
	ChannelDefaultUserID
	ChannelUserID
	ChannelPassword
	ChannelClientId
	ConnectionDatabaseName
	ConnectionPoolUseDedicatedChannelPerConnection
	ConnectionPoolDefaultPoolSize
	ConnectionReserveTimeoutSeconds
	ConnectionOperationTimeoutSeconds
	TlsProviderName
	TlsProviderClassName
	TlsProviderConfigFile
	TlsProtocol
	TlsCipherSuites
	TlsVerifyDatabaseName
	TlsExpectedHostName
	TlsTrustedCertificates
	KeyStorePassword
	EnableConnectionTrace
	ConnectionTraceDir
	InvalidName
)

type ConfigName struct {
	configPropName string
	aliasName      string
	defaultValue   string
	description    string
}

var PreDefinedConfigurations = map[int]ConfigName{
	ChannelDefaultHost:                             {configPropName: "tgdb.channel.defaultHost", aliasName: "defaultHost", defaultValue: "localhost", description: "The default host specifier"},
	ChannelDefaultPort:                             {configPropName: "tgdb.channel.defaultPort", aliasName: "defaultPort", defaultValue: "8700", description: "The default port specifier"},
	ChannelDefaultProtocol:                         {configPropName: "tgdb.channel.defaultProtocol", aliasName: "defaultProtocol", defaultValue: "tcp", description: "The default protocol"},
	ChannelSendSize:                                {configPropName: "tgdb.channel.sendSize", aliasName: "sendSize", defaultValue: "122", description: "TCP send packet size in KBs"},
	ChannelRecvSize:                                {configPropName: "tgdb.channel.recvSize", aliasName: "recvSize", defaultValue: "128", description: "TCP recv packet size in KB"},
	ChannelPingInterval:                            {configPropName: "tgdb.channel.pingInterval", aliasName: "pingInterval", defaultValue: "30", description: "Keep alive ping intervals"},
	ChannelConnectTimeout:                          {configPropName: "tgdb.channel.connectTimeout", aliasName: "connectTimeout", defaultValue: "1000", description: "Timeout for connection to establish, before it gives up and tries the ftUrls if specified"}, //1 sec timeout
	ChannelFTHosts:                                 {configPropName: "tgdb.channel.ftHosts", aliasName: "ftHosts", defaultValue: "", description: "Alternate fault tolerant list of &lt;host:port&gt; pair separated by comma"},
	ChannelFTRetryIntervalSeconds:                  {configPropName: "tgdb.channel.ftRetryIntervalSeconds", aliasName: "ftRetryIntervalSeconds", defaultValue: "10", description: "The connect retry interval to ftHosts"},
	ChannelFTRetryCount:                            {configPropName: "tgdb.channel.ftRetryCount", aliasName: "ftRetryCount", defaultValue: "3", description: "The number of times ro retry"},
	ChannelDefaultUserID:                           {configPropName: "tgdb.channel.defaultUserID", aliasName: "defaultUserID", defaultValue: "", description: "The default user id for the connection"},
	ChannelUserID:                                  {configPropName: "tgdb.channel.userID", aliasName: "userID", defaultValue: "", description: "The user id for the connection if it is not specified in the API. See the rules for picking the user name"},
	ChannelPassword:                                {configPropName: "tgdb.channel.password", aliasName: "password", defaultValue: "", description: "The password for the username"},
	ChannelClientId:                                {configPropName: "tgdb.channel.clientId", aliasName: "clientId", defaultValue: "tgdb.go-api.client", description: "The client id to be used for the connection"},
	ConnectionDatabaseName:                         {configPropName: "tgdb.connection.dbName", aliasName: "dbName", defaultValue: "", description: "The database name the client is connecting to. It is used as part of verification for ssl channels"},
	ConnectionPoolUseDedicatedChannelPerConnection: {configPropName: "tgdb.connectionpool.useDedicatedChannelPerConnection", aliasName: "useDedicatedChannelPerConnection", defaultValue: "false", description: ""},
	ConnectionPoolDefaultPoolSize:                  {configPropName: "tgdb.connectionpool.defaultPoolSize", aliasName: "defaultPoolSize", defaultValue: "10", description: "The default connection pool size to use when creating a ConnectionPool"},
	//0 = mean immediate, Integer Max for indefinite
	ConnectionReserveTimeoutSeconds:                {configPropName: "tgdb.connectionpool.connectionReserveTimeoutSeconds", aliasName: "connectionReserveTimeoutSeconds", defaultValue: "10", description: "A timeout parameter indicating how long to wait before getting a connection from the pool"},
	//Represented in ms. Default Value is 10sec
	ConnectionOperationTimeoutSeconds:              {configPropName: "tgdb.connection.operationTimeoutSeconds", aliasName: "connectionOperationTimeoutSeconds", defaultValue: "10", description: "A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior"},
	// TODO: Ask TGDB Engineering Team
	TlsProviderName:                                {configPropName: "tgdb.tls.provider.name", aliasName: "tlsProviderName", defaultValue: "SunJSSE", description: "Transport level Security provider. Work with your InfoSec team to change this value"},
	// TODO: Ask TGDB Engineering Team - The default is the Sun JSSE. One can specify the tibco wrapper class for FIPS
	TlsProviderClassName:                           {configPropName: "tgdb.tls.provider.className", aliasName: "tlsProviderClassName", defaultValue: "com.sun.net.ssl.internal.ssl.Provider", description: "The underlying Provider implementation. Work with your InfoSec team to change this value"},
	TlsProviderConfigFile:                          {configPropName: "tgdb.tls.provider.configFile", aliasName: "tlsProviderConfigFile", defaultValue: "", description: "Some providers require extra configuration paramters, and it can be passed as a file"},
	TlsProtocol:                                    {configPropName: "tgdb.tls.protocol", aliasName: "tlsProtocol", defaultValue: "TLSv1.2", description: "TLSProtocol version. The system only supports 1.2+"},
	//Use the Default Cipher Suites
	TlsCipherSuites:                                {configPropName: "tgdb.tls.cipherSuites", aliasName: "cipherSuites", defaultValue: "", description: "A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher list and Openssl list that supports 1.2 protocol"},
	TlsVerifyDatabaseName:                          {configPropName: "tgdb.tls.verifyDBName", aliasName: "verifyDBName", defaultValue: "false", description: "Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL"},
	TlsExpectedHostName:                            {configPropName: "tgdb.tls.expectedHostName", aliasName: "expectedHostName", defaultValue: "", description: "The expected hostName for the certificate. This is for future use"},
	TlsTrustedCertificates:                         {configPropName: "tgdb.tls.trustedCertificates", aliasName: "trustedCertificates", defaultValue: "", description: "The list of trusted Certificates"},
	KeyStorePassword:                               {configPropName: "tgdb.security.keyStorePassword", aliasName: "keyStorePassword", defaultValue: "", description: "The Keystore for the password"},
	EnableConnectionTrace:                          {configPropName: "tgdb.connection.enableTrace", aliasName: "enableTrace", defaultValue: "false", description: "The flag for debugging purpose, to enable the commit trace"},
	ConnectionTraceDir:                             {configPropName: "tgdb.connection.enableTraceDir", aliasName: "enableTraceDir", defaultValue: ".", description: "The base directory to hold commit trace log"},
	InvalidName:                                    {configPropName: "", aliasName: "", defaultValue: "", description: ""},
}

// Make sure that the ConfigName implements the TGConfigName interface
var _ types.TGConfigName = (*ConfigName)(nil)

func NewConfigName(name, alias string, value string) *ConfigName {
	existingConfig := GetConfigFromName(name)
	if existingConfig.configPropName != "" && existingConfig.aliasName != "" {
		return existingConfig
	}
	return &ConfigName{configPropName: name, aliasName: alias, defaultValue: value}
}

/////////////////////////////////////////////////////////////////
// Helper Public functions for TGConfigName
/////////////////////////////////////////////////////////////////

// GetConfigFromKey returns the TGConfigName given its full qualified string form or its alias name.
func GetConfigFromKey(key int) *ConfigName {
	if config, ok := PreDefinedConfigurations[key]; ok {
		return &config
	}
	invalid := PreDefinedConfigurations[InvalidName]
	return &invalid
}

// GetConfigFromKey returns the TGConfigName for specified name
func GetConfigFromName(name string) *ConfigName {
	for _, config := range PreDefinedConfigurations {
		if strings.ToLower(config.configPropName) == strings.ToLower(name) {
			return &config
		}
		if (config.aliasName != "") && (strings.ToLower(config.aliasName) == strings.ToLower(name)) {
			return &config
		}
	}
	invalid := PreDefinedConfigurations[InvalidName]
	return &invalid
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConfigName
/////////////////////////////////////////////////////////////////

// GetAlias gets configuration Alias
func (c *ConfigName) GetAlias() string {
	return c.aliasName
}

// GetDefaultValue gets configuration Default Value
func (c *ConfigName) GetDefaultValue() string {
	return c.defaultValue
}

// GetName gets configuration name
func (c *ConfigName) GetName() string {
	return c.configPropName
}

// GetDesc gets configuration description
func (c *ConfigName) GetDesc() string {
	return c.description
}

// SetDesc sets configuration description
func (c *ConfigName) SetDesc(desc string) {
	c.description = desc
}
