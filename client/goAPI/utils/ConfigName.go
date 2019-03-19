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
	PropName     string
	AliasName    string
	DefaultValue string
	Desc         string
}

var PreDefinedConfigurations = map[int]ConfigName{
	ChannelDefaultHost:                             {PropName: "tgdb.channel.defaultHost", AliasName: "defaultHost", DefaultValue: "localhost", Desc: "The default host specifier"},
	ChannelDefaultPort:                             {PropName: "tgdb.channel.defaultPort", AliasName: "defaultPort", DefaultValue: "8700", Desc: "The default port specifier"},
	ChannelDefaultProtocol:                         {PropName: "tgdb.channel.defaultProtocol", AliasName: "defaultProtocol", DefaultValue: "tcp", Desc: "The default protocol"},
	ChannelSendSize:                                {PropName: "tgdb.channel.sendSize", AliasName: "sendSize", DefaultValue: "122", Desc: "TCP send packet size in KBs"},
	ChannelRecvSize:                                {PropName: "tgdb.channel.recvSize", AliasName: "recvSize", DefaultValue: "128", Desc: "TCP recv packet size in KB"},
	ChannelPingInterval:                            {PropName: "tgdb.channel.pingInterval", AliasName: "pingInterval", DefaultValue: "30", Desc: "Keep alive ping intervals"},
	ChannelConnectTimeout:                          {PropName: "tgdb.channel.connectTimeout", AliasName: "connectTimeout", DefaultValue: "1000", Desc: "Timeout for connection to establish, before it gives up and tries the ftUrls if specified"}, //1 sec timeout
	ChannelFTHosts:                                 {PropName: "tgdb.channel.ftHosts", AliasName: "ftHosts", DefaultValue: "", Desc: "Alternate fault tolerant list of &lt;host:port&gt; pair separated by comma"},
	ChannelFTRetryIntervalSeconds:                  {PropName: "tgdb.channel.ftRetryIntervalSeconds", AliasName: "ftRetryIntervalSeconds", DefaultValue: "10", Desc: "The connect retry interval to ftHosts"},
	ChannelFTRetryCount:                            {PropName: "tgdb.channel.ftRetryCount", AliasName: "ftRetryCount", DefaultValue: "3", Desc: "The number of times ro retry"},
	ChannelDefaultUserID:                           {PropName: "tgdb.channel.defaultUserID", AliasName: "defaultUserID", DefaultValue: "", Desc: "The default user id for the connection"},
	ChannelUserID:                                  {PropName: "tgdb.channel.userID", AliasName: "userID", DefaultValue: "", Desc: "The user id for the connection if it is not specified in the API. See the rules for picking the user name"},
	ChannelPassword:                                {PropName: "tgdb.channel.password", AliasName: "password", DefaultValue: "", Desc: "The password for the username"},
	ChannelClientId:                                {PropName: "tgdb.channel.clientId", AliasName: "clientId", DefaultValue: "tgdb.go-api.client", Desc: "The client id to be used for the connection"},
	ConnectionDatabaseName:                         {PropName: "tgdb.connection.dbName", AliasName: "dbName", DefaultValue: "", Desc: "The database name the client is connecting to. It is used as part of verification for ssl channels"},
	ConnectionPoolUseDedicatedChannelPerConnection: {PropName: "tgdb.connectionpool.useDedicatedChannelPerConnection", AliasName: "useDedicatedChannelPerConnection", DefaultValue: "false", Desc: ""},
	ConnectionPoolDefaultPoolSize:                  {PropName: "tgdb.connectionpool.defaultPoolSize", AliasName: "defaultPoolSize", DefaultValue: "10", Desc: "The default connection pool size to use when creating a ConnectionPool"},
	//0 = mean immediate, Integer Max for indefinite
	ConnectionReserveTimeoutSeconds:                {PropName: "tgdb.connectionpool.connectionReserveTimeoutSeconds", AliasName: "connectionReserveTimeoutSeconds", DefaultValue: "10", Desc: "A timeout parameter indicating how long to wait before getting a connection from the pool"},
	//Represented in ms. Default Value is 10sec
	ConnectionOperationTimeoutSeconds:              {PropName: "tgdb.connection.operationTimeoutSeconds", AliasName: "connectionOperationTimeoutSeconds", DefaultValue: "10", Desc: "A timeout parameter indicating how long to wait for a operation before giving up. Some queries are long running, and may override this behavior"},
	// TODO: Ask TGDB Engineering Team
	TlsProviderName:                                {PropName: "tgdb.tls.provider.name", AliasName: "tlsProviderName", DefaultValue: "SunJSSE", Desc: "Transport level Security provider. Work with your InfoSec team to change this value"},
	// TODO: Ask TGDB Engineering Team - The default is the Sun JSSE. One can specify the tibco wrapper class for FIPS
	TlsProviderClassName:                           {PropName: "tgdb.tls.provider.className", AliasName: "tlsProviderClassName", DefaultValue: "com.sun.net.ssl.internal.ssl.Provider", Desc: "The underlying Provider implementation. Work with your InfoSec team to change this value"},
	TlsProviderConfigFile:                          {PropName: "tgdb.tls.provider.configFile", AliasName: "tlsProviderConfigFile", DefaultValue: "", Desc: "Some providers require extra configuration paramters, and it can be passed as a file"},
	TlsProtocol:                                    {PropName: "tgdb.tls.protocol", AliasName: "tlsProtocol", DefaultValue: "TLSv1.2", Desc: "TLSProtocol version. The system only supports 1.2+"},
	//Use the Default Cipher Suites
	TlsCipherSuites:                                {PropName: "tgdb.tls.cipherSuites", AliasName: "cipherSuites", DefaultValue: "", Desc: "A list cipher suites that the InfoSec team has cleared. The default list is a common list of JSSE's cipher list and Openssl list that supports 1.2 protocol"},
	TlsVerifyDatabaseName:                          {PropName: "tgdb.tls.verifyDBName", AliasName: "verifyDBName", DefaultValue: "false", Desc: "Verify the Database name in the certificate. TGDB provides self signed certificate for easy-to-use SSL"},
	TlsExpectedHostName:                            {PropName: "tgdb.tls.expectedHostName", AliasName: "expectedHostName", DefaultValue: "", Desc: "The expected hostName for the certificate. This is for future use"},
	TlsTrustedCertificates:                         {PropName: "tgdb.tls.trustedCertificates", AliasName: "trustedCertificates", DefaultValue: "", Desc: "The list of trusted Certificates"},
	KeyStorePassword:                               {PropName: "tgdb.security.keyStorePassword", AliasName: "keyStorePassword", DefaultValue: "", Desc: "The Keystore for the password"},
	EnableConnectionTrace:                          {PropName: "tgdb.connection.enableTrace", AliasName: "enableTrace", DefaultValue: "false", Desc: "The flag for debugging purpose, to enable the commit trace"},
	ConnectionTraceDir:                             {PropName: "tgdb.connection.enableTraceDir", AliasName: "enableTraceDir", DefaultValue: ".", Desc: "The base directory to hold commit trace log"},
	InvalidName:                                    {PropName: "", AliasName: "", DefaultValue: "", Desc: ""},
}

// Make sure that the ConfigName implements the TGConfigName interface
var _ types.TGConfigName = (*ConfigName)(nil)

func NewConfigName(name, alias string, value string) *ConfigName {
	existingConfig := GetConfigFromName(name)
	if existingConfig.PropName != "" && existingConfig.AliasName != "" {
		return existingConfig
	}
	return &ConfigName{PropName: name, AliasName: alias, DefaultValue: value}
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
		if strings.ToLower(config.PropName) == strings.ToLower(name) {
			return &config
		}
		if (config.AliasName != "") && (strings.ToLower(config.AliasName) == strings.ToLower(name)) {
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
	return c.AliasName
}

// GetDefaultValue gets configuration Default Value
func (c *ConfigName) GetDefaultValue() string {
	return c.DefaultValue
}

// GetName gets configuration name
func (c *ConfigName) GetName() string {
	return c.PropName
}
