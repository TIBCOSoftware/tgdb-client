package admin

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"time"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: AdminHelper.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/////////////////////////////////////////////////////////////////
// Available functions for AdminHelper
/////////////////////////////////////////////////////////////////

func extractCacheInfoFromInputStream(is types.TGInputStream) (*CacheStatisticsImpl, types.TGError) {
	dataCacheMaxEntries, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // data cache max entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMaxEntries from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMaxEntries as '%+v'", dataCacheMaxEntries))

	dataCacheEntries, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // data cache entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheEntries from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheEntries as '%+v'", dataCacheEntries))

	dataCacheHits, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data cache hits
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheHits from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheHits as '%+v'", dataCacheHits))

	dataCacheMisses, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data cache misses
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMisses from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMisses as '%+v'", dataCacheMisses))

	dataCacheMaxMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data cache max memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMaxMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMaxMemory as '%+v'", dataCacheMaxMemory))

	indexCacheMaxEntries, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // index cache max entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMaxEntries from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMaxEntries as '%+v'", indexCacheMaxEntries))

	indexCacheEntries, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // index cache entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheEntries from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheEntries as '%+v'", indexCacheEntries))

	indexCacheHits, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index cache hits
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheHits from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheHits as '%+v'", indexCacheHits))

	indexCacheMisses, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index cache misses
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMisses from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMisses as '%+v'", indexCacheMisses))

	indexCacheMaxMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index cache max memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMaxMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMaxMemory as '%+v'", indexCacheMaxMemory))

	cacheInfoImpl := NewCacheStatisticsImpl(dataCacheMaxEntries, dataCacheEntries, dataCacheHits,
		dataCacheMisses, dataCacheMaxMemory, indexCacheMaxEntries, indexCacheEntries, indexCacheHits,
		indexCacheMisses, indexCacheMaxMemory)
	return cacheInfoImpl, nil
}

func extractConnectionListFromInputStream(is types.TGInputStream) ([]TGConnectionInfo, types.TGError) {
	connList := make([]TGConnectionInfo, 0)
	connCount, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // connection count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading connCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read connCount as '%+v'", connCount))

	for i := 0; i < int(connCount); i++ {
		listenerName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // Listener name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading listenerName from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read listenerName as '%+v'", listenerName))

		clientId, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // client Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading clientId from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read clientId as '%+v'", clientId))

		sessionId, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // session Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading sessionId from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read sessionId as '%+v'", sessionId))

		userName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // user name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userName from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userName as '%+v'", userName))

		remoteAddr, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // remote address
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading remoteAddr from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read remoteAddr as '%+v'", remoteAddr))

		tStamp, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // timestamp
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading tStamp from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read tStamp as '%+v'", tStamp))

		createTS := time.Now().Unix() - (tStamp/1000000000)
		connectionInfo := NewConnectionInfoImpl(listenerName, clientId, sessionId, userName, remoteAddr, createTS)
		connList = append(connList, connectionInfo)
	}
	return connList, nil
}

func extractDatabaseInfoFromInputStream(is types.TGInputStream) (*DatabaseStatisticsImpl, types.TGError) {
	dbSize, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // database size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dbSize from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dbSize as '%+v'", dbSize))

	numDataSegments, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // no of data segments
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numDataSegments from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read numDataSegments as '%+v'", numDataSegments))

	dataSize, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataSize from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataSize as '%+v'", dataSize))

	dataUsed, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data used
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataUsed from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataUsed as '%+v'", dataUsed))

	dataFree, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // data free
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataFree from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataFree as '%+v'", dataFree))

	dataBlockSize, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // data block size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataBlockSize from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataBlockSize as '%+v'", dataBlockSize))

	numIndexSegments, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // no of index segments
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numIndexSegments from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read numIndexSegments as '%+v'", numIndexSegments))

	indexSize, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexSize from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexSize as '%+v'", indexSize))

	indexUsed, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index used
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexUsed from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexUsed as '%+v'", indexUsed))

	indexFree, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // index free
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexFree from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexFree as '%+v'", indexFree))

	indexBlockSize, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // index block size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexBlockSize from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexBlockSize as '%+v'", indexBlockSize))

	databaseInfoImpl := NewDatabaseStatisticsImpl(dbSize, numDataSegments, dataSize, dataUsed, dataFree,
		dataBlockSize, numIndexSegments, indexSize, indexUsed, indexFree, indexBlockSize)
	return databaseInfoImpl, nil
}

func extractDescriptorListFromInputStream(is types.TGInputStream) ([]types.TGAttributeDescriptor, types.TGError) {
	attrDescList := make([]types.TGAttributeDescriptor, 0)
	descCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // attribute descriptor count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading descCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read descCount as '%+v'", descCount))

	for i := 0; i < descCount; i++ {
		aType, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading aType from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read aType as '%+v'", aType))

		attrId, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // Attribute Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrId from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrId as '%+v'", attrId))

		attrName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // attribute name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrName from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrName as '%+v'", attrName))

		attrType, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // attribute type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrType from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrType as '%+v'", attrType))

		isArrayFlag, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean() // isArray flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isArrayFlag from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isArrayFlag as '%+v'", isArrayFlag))

		isEncryptedFlag, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean() // isEncrypted flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isEncryptedFlag from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isEncryptedFlag as '%+v'", isEncryptedFlag))

		attrDesc := model.NewAttributeDescriptorAsArray(attrName, types.GetAttributeTypeFromId(int(attrType)).GetTypeId(), isArrayFlag)
		attrDesc.SetAttributeId(int64(attrId))
		attrDesc.SetIsEncrypted(isEncryptedFlag)

		if attrType == types.AttributeTypeNumber {
			precision, err := is.(*iostream.ProtocolDataInputStream).ReadShort() // precision
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading precision from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read precision as '%+v'", precision))

			scale, err := is.(*iostream.ProtocolDataInputStream).ReadShort() // scale
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading scale from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read scale as '%+v'", scale))
			attrDesc.SetPrecision(precision)
			attrDesc.SetScale(scale)
		}
		attrDescList = append(attrDescList, attrDesc)
	}
	return attrDescList, nil
}

func extractIndexListFromInputStream(is types.TGInputStream) ([]TGIndexInfo, types.TGError) {
	indexList := make([]TGIndexInfo, 0)
	indexCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // index count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCount as '%+v'", indexCount))

	for i := 0; i < indexCount; i++ {
		indexType, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // Index Type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexType from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexType as '%+v'", indexType))

		indexId, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // Index Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexId from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexId as '%+v'", indexId))

		indexName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // Index name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexName from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexName as '%+v'", indexName))

		isUniqueFlag, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean() // isUnique flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isUniqueFlag from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isUniqueFlag as '%+v'", isUniqueFlag))

		attrCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // attribute count
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrCount from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrCount as '%+v'", attrCount))

		attributes := make([]string, 0)
		for j := 0; j < attrCount; j++ {
			attrName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // attribute name
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrName from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrName as '%+v'", attrName))
			attributes = append(attributes, attrName)
		}

		nodeCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // node count
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nodeCount from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nodeCount as '%+v'", nodeCount))

		nodes := make([]string, 0)
		for k := 0; k < nodeCount; k++ {
			nodeName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // node name
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nodeName from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nodeName as '%+v'", nodeName))
			nodes = append(nodes, nodeName)
		}

		blocksize, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // block size
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading blocksize from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read blocksize as '%+v'", blocksize))

		indexInfo := NewIndexInfoImpl(indexId, indexName, indexType, isUniqueFlag, attributes, nodes)
		indexList = append(indexList, indexInfo)
	}
	return indexList, nil
}

func extractMemoryInfoFromInputStream(is types.TGInputStream) (*ServerMemoryInfoImpl, types.TGError) {
	freeProcessMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // free process memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading freeProcessMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read freeProcessMemory as '%+v'", freeProcessMemory))

	memProcessUsagePct, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // process memory usage percentage
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memProcessUsagePct from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memProcessUsagePct as '%+v'", memProcessUsagePct))

	maxProcessMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // max process memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxProcessMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxProcessMemory as '%+v'", maxProcessMemory))
	usedProcessMemory := maxProcessMemory - freeProcessMemory

	fileLocation, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // shared file location
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading fileLocation from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read fileLocation as '%+v'", fileLocation))

	freeSharedMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // free shared memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading freeSharedMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read freeSharedMemory as '%+v'", freeSharedMemory))

	memSharedUsagePct, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // shared memory usage percentage
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memSharedUsagePct from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memSharedUsagePct as '%+v'", memSharedUsagePct))

	maxSharedMemory, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // max shared memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxSharedMemory from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxSharedMemory as '%+v'", maxSharedMemory))
	usedSharedMemory := maxSharedMemory - freeSharedMemory

	processMemory := NewMemoryInfoImpl(freeProcessMemory, maxProcessMemory, usedProcessMemory, "")
	sharedMemory := NewMemoryInfoImpl(freeSharedMemory, maxSharedMemory, usedSharedMemory, fileLocation)

	memoryInfo := NewServerMemoryInfoImpl(processMemory, sharedMemory)
	return memoryInfo, nil
}

func extractNetListenersInfoFromInputStream(is types.TGInputStream) ([]TGNetListenerInfo, types.TGError) {
	listenerList := make([]TGNetListenerInfo, 0)
	listenerCount, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // listener count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading listenerCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read listenerCount as '%+v'", listenerCount))

	for i := 0; i < int(listenerCount); i++ {
		bufLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buffer length
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bufLen from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read bufLen as '%+v'", bufLen))

		nameBytes := make([]byte, bufLen)
		_, err = is.(*iostream.ProtocolDataInputStream).ReadIntoBuffer(nameBytes) // name bytes
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nameBytes from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nameBytes as '%+v'", nameBytes))
		listenerName := string(nameBytes)

		currConn, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // current connection Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading currConn from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read currConn as '%+v'", currConn))

		maxConn, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // maximum connections
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxConn from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxConn as '%+v'", maxConn))

		bufLen4Port, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buffer length for port
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bufLen4Port from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read bufLen4Port as '%+v'", bufLen4Port))

		portBytes := make([]byte, bufLen4Port)
		_, err = is.(*iostream.ProtocolDataInputStream).ReadIntoBuffer(portBytes) // port bytes
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading portBytes from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read portBytes as '%+v'", portBytes))
		portNumber := string(portBytes)

		listenerInfo := NewNetListenerInfoImpl(currConn, maxConn, listenerName, portNumber)
		listenerList = append(listenerList, listenerInfo)
	}
	return listenerList, nil
}

func extractServerInfoFromInputStream(is types.TGInputStream) (*ServerInfoImpl, types.TGError) {
	serverStatus, err  := extractServerStatusFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading serverStatus from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read serverStatus as '%+v'", serverStatus))

	memoryInfo, err  := extractMemoryInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memoryInfo from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memoryInfo as '%+v'", memoryInfo))

	netListenerInfo, err  := extractNetListenersInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading netListenerInfo from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read netListenerInfo as '%+v'", netListenerInfo))

	transactionInfo, err  := extractTransactionsInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading transactionInfo from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read transactionInfo as '%+v'", transactionInfo))

	cacheInfo, err  := extractCacheInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading cacheInfo from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read cacheInfo as '%+v'", cacheInfo))

	databaseInfo, err  := extractDatabaseInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading databaseInfo from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read databaseInfo as '%+v'", databaseInfo))

	serverInfo := NewServerInfoImpl(cacheInfo, databaseInfo, memoryInfo, netListenerInfo, serverStatus, transactionInfo)
	return serverInfo, nil
}

func extractServerStatusFromInputStream(is types.TGInputStream) (*ServerStatusImpl, types.TGError) {
	name, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // name
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading name from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read name as '%+v'", name))

	status, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // status
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading status from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read status as '%+v'", status))

	processId, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // process Id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading processId from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read processId as '%+v'", processId))

	duration, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // duration
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading duration from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read duration as '%+v'", duration))

	sVersion, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // version
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading sVersion from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read sVersion as '%+v'", sVersion))

	sVersionInfo := utils.NewTGServerVersion(sVersion)
	uptime := time.Duration(time.Now().Unix() - (duration/1000000000))
	serverInfo := NewServerStatusImpl(name, sVersionInfo, string(processId), ServerStates(status), uptime)

	return serverInfo, nil
}

func extractTransactionsInfoFromInputStream(is types.TGInputStream) (*TransactionStatisticsImpl, types.TGError) {
	txnProcessorsCount, err := is.(*iostream.ProtocolDataInputStream).ReadShort() // txn processor count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnProcessorsCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnProcessorsCount as '%+v'", txnProcessorsCount))

	txnProcessedCount, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // txn processed count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnProcessoedCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnProcessedCount as '%+v'", txnProcessedCount))

	txnSuccessCount, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // txn success count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnSuccessCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnSuccessCount as '%+v'", txnSuccessCount))

	avgProcessTime, err := is.(*iostream.ProtocolDataInputStream).ReadDouble() // average processing time
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading avgProcessTime from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read avgProcessTime as '%+v'", avgProcessTime))

	pendingTxnCount, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // pending txn count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading pendingTxnCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read pendingTxnCount as '%+v'", pendingTxnCount))

	txnLogQueueDepth, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // txn logger queue depth
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnLogQueueDepth from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnLogQueueDepth as '%+v'", txnLogQueueDepth))

	transactionsInfo := NewTransactionStatisticsImpl(avgProcessTime, pendingTxnCount, txnLogQueueDepth, int64(txnProcessorsCount), txnProcessedCount, txnSuccessCount)
	return transactionsInfo, nil
}

func extractUserListFromInputStream(is types.TGInputStream) ([]TGUserInfo, types.TGError) {
	userList := make([]TGUserInfo, 0)
	userCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // user count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userCount from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userCount as '%+v'", userCount))

	for i := 0; i < userCount; i++ {
		userType, err := is.(*iostream.ProtocolDataInputStream).ReadByte() // user Type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userType from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userType as '%+v'", userType))

		userId, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // user Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userId from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userId as '%+v'", userId))

		userName, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // user name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userName from message buffer"))
			return nil, err
		}
		logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userName as '%+v'", userName))

		switch types.TGSystemType(userType) {
		case types.SystemTypePrincipal:
			bufLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buffer length
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bufLen from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read bufLen as '%+v'", bufLen))

			buf := make([]byte, bufLen)
			_, err = is.(*iostream.ProtocolDataInputStream).ReadIntoBuffer(buf) // user role
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading principal info from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprint("Inside AdminResponseMessage:ReadPayload read principal info"))

			userRole, err := is.(*iostream.ProtocolDataInputStream).ReadUTF() // user role
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userRole from message buffer"))
				return nil, err
			}
			logger.Log(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userRole as '%+v'", userRole))
		default:
		}
		userInfo := NewUserInfoImpl(userId, userName, userType)
		userList = append(userList, userInfo)
	}
	return userList, nil
}
