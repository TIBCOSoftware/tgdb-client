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
 * File Name: adminimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: adminimpl.go 4588 2020-10-28 23:56:58Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"tgdb"
	"time"
)

func extractCacheInfoFromInputStream(is tgdb.TGInputStream) (*CacheStatisticsImpl, tgdb.TGError) {
	dataCacheMaxEntries, err := is.(*ProtocolDataInputStream).ReadInt() // data cache max entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMaxEntries from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMaxEntries as '%+v'", dataCacheMaxEntries))
	}

	dataCacheEntries, err := is.(*ProtocolDataInputStream).ReadInt() // data cache entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheEntries from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheEntries as '%+v'", dataCacheEntries))
	}

	dataCacheHits, err := is.(*ProtocolDataInputStream).ReadLong() // data cache hits
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheHits from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheHits as '%+v'", dataCacheHits))
	}

	dataCacheMisses, err := is.(*ProtocolDataInputStream).ReadLong() // data cache misses
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMisses from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMisses as '%+v'", dataCacheMisses))
	}

	dataCacheMaxMemory, err := is.(*ProtocolDataInputStream).ReadLong() // data cache max memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataCacheMaxMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataCacheMaxMemory as '%+v'", dataCacheMaxMemory))
	}

	indexCacheMaxEntries, err := is.(*ProtocolDataInputStream).ReadInt() // index cache max entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMaxEntries from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMaxEntries as '%+v'", indexCacheMaxEntries))
	}

	indexCacheEntries, err := is.(*ProtocolDataInputStream).ReadInt() // index cache entries
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheEntries from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheEntries as '%+v'", indexCacheEntries))
	}

	indexCacheHits, err := is.(*ProtocolDataInputStream).ReadLong() // index cache hits
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheHits from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheHits as '%+v'", indexCacheHits))
	}

	indexCacheMisses, err := is.(*ProtocolDataInputStream).ReadLong() // index cache misses
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMisses from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMisses as '%+v'", indexCacheMisses))
	}

	indexCacheMaxMemory, err := is.(*ProtocolDataInputStream).ReadLong() // index cache max memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCacheMaxMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCacheMaxMemory as '%+v'", indexCacheMaxMemory))
	}

	cacheInfoImpl := NewCacheStatisticsImpl(dataCacheMaxEntries, dataCacheEntries, dataCacheHits,
		dataCacheMisses, dataCacheMaxMemory, indexCacheMaxEntries, indexCacheEntries, indexCacheHits,
		indexCacheMisses, indexCacheMaxMemory)
	return cacheInfoImpl, nil
}

func extractConnectionListFromInputStream(is tgdb.TGInputStream) ([]tgdb.TGConnectionInfo, tgdb.TGError) {
	connList := make([]tgdb.TGConnectionInfo, 0)
	connCount, err := is.(*ProtocolDataInputStream).ReadLong() // connection count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading connCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read connCount as '%+v'", connCount))
	}

	for i := 0; i < int(connCount); i++ {
		listenerName, err := is.(*ProtocolDataInputStream).ReadUTF() // Listener Name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading listenerName from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read listenerName as '%+v'", listenerName))
		}

		clientId, err := is.(*ProtocolDataInputStream).ReadUTF() // client Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading clientId from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read clientId as '%+v'", clientId))
		}

		sessionId, err := is.(*ProtocolDataInputStream).ReadLong() // session Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading sessionId from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read sessionId as '%+v'", sessionId))
		}

		userName, err := is.(*ProtocolDataInputStream).ReadUTF() // user Name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userName from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userName as '%+v'", userName))
		}

		remoteAddr, err := is.(*ProtocolDataInputStream).ReadUTF() // remote address
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading remoteAddr from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read remoteAddr as '%+v'", remoteAddr))
		}

		tStamp, err := is.(*ProtocolDataInputStream).ReadLong() // timestamp
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading tStamp from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read tStamp as '%+v'", tStamp))
		}

		createTS := time.Now().Unix() - (tStamp/1000000000)
		connectionInfo := NewConnectionInfoImpl(listenerName, clientId, sessionId, userName, remoteAddr, createTS)
		connList = append(connList, connectionInfo)
	}
	return connList, nil
}

func extractDatabaseInfoFromInputStream(is tgdb.TGInputStream) (*DatabaseStatisticsImpl, tgdb.TGError) {
	dbSize, err := is.(*ProtocolDataInputStream).ReadLong() // database size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dbSize from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dbSize as '%+v'", dbSize))
	}

	numDataSegments, err := is.(*ProtocolDataInputStream).ReadInt() // no of data segments
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numDataSegments from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read numDataSegments as '%+v'", numDataSegments))
	}
	dataSize, err := is.(*ProtocolDataInputStream).ReadLong() // data size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataSize from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataSize as '%+v'", dataSize))
	}
	dataUsed, err := is.(*ProtocolDataInputStream).ReadLong() // data used
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataUsed from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataUsed as '%+v'", dataUsed))
	}
	dataFree, err := is.(*ProtocolDataInputStream).ReadLong() // data free
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataFree from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataFree as '%+v'", dataFree))
	}
	dataBlockSize, err := is.(*ProtocolDataInputStream).ReadInt() // data block size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading dataBlockSize from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read dataBlockSize as '%+v'", dataBlockSize))
	}
	numIndexSegments, err := is.(*ProtocolDataInputStream).ReadInt() // no of index segments
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numIndexSegments from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read numIndexSegments as '%+v'", numIndexSegments))
	}
	indexSize, err := is.(*ProtocolDataInputStream).ReadLong() // index size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexSize from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexSize as '%+v'", indexSize))
	}
	indexUsed, err := is.(*ProtocolDataInputStream).ReadLong() // index used
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexUsed from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexUsed as '%+v'", indexUsed))
	}
	indexFree, err := is.(*ProtocolDataInputStream).ReadLong() // index free
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexFree from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexFree as '%+v'", indexFree))
	}
	indexBlockSize, err := is.(*ProtocolDataInputStream).ReadInt() // index block size
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexBlockSize from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexBlockSize as '%+v'", indexBlockSize))
	}
	databaseInfoImpl := NewDatabaseStatisticsImpl(dbSize, numDataSegments, dataSize, dataUsed, dataFree,
		dataBlockSize, numIndexSegments, indexSize, indexUsed, indexFree, indexBlockSize)
	return databaseInfoImpl, nil
}

func extractDescriptorListFromInputStream(is tgdb.TGInputStream) ([]tgdb.TGAttributeDescriptor, tgdb.TGError) {
	attrDescList := make([]tgdb.TGAttributeDescriptor, 0)
	descCount, err := is.(*ProtocolDataInputStream).ReadInt() // attribute descriptor count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading descCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read descCount as '%+v'", descCount))
	}
	for i := 0; i < descCount; i++ {
		aType, err := is.(*ProtocolDataInputStream).ReadByte() // type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading aType from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read aType as '%+v'", aType))
		}
		attrId, err := is.(*ProtocolDataInputStream).ReadInt() // Attribute Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrId from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrId as '%+v'", attrId))
		}
		attrName, err := is.(*ProtocolDataInputStream).ReadUTF() // attribute Name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrName from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrName as '%+v'", attrName))
		}
		attrType, err := is.(*ProtocolDataInputStream).ReadByte() // attribute type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading AttrType from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read AttrType as '%+v'", attrType))
		}
		isArrayFlag, err := is.(*ProtocolDataInputStream).ReadBoolean() // IsArray flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isArrayFlag from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isArrayFlag as '%+v'", isArrayFlag))
		}
		isEncryptedFlag, err := is.(*ProtocolDataInputStream).ReadBoolean() // Is_Encrypted flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isEncryptedFlag from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isEncryptedFlag as '%+v'", isEncryptedFlag))
		}
		attrDesc := NewAttributeDescriptorAsArray(attrName, GetAttributeTypeFromId(int(attrType)).GetTypeId(), isArrayFlag)
		attrDesc.SetAttributeId(int64(attrId))
		attrDesc.SetIsEncrypted(isEncryptedFlag)

		if attrType == AttributeTypeNumber {
			precision, err := is.(*ProtocolDataInputStream).ReadShort() // Precision
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading Precision from message buffer"))
				return nil, err
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read Precision as '%+v'", precision))
			}
			scale, err := is.(*ProtocolDataInputStream).ReadShort() // Scale
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading Scale from message buffer"))
				return nil, err
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read Scale as '%+v'", scale))
			}
			attrDesc.SetPrecision(precision)
			attrDesc.SetScale(scale)
		}
		attrDescList = append(attrDescList, attrDesc)
	}
	return attrDescList, nil
}

func extractIndexListFromInputStream(is tgdb.TGInputStream) ([]tgdb.TGIndexInfo, tgdb.TGError) {
	indexList := make([]tgdb.TGIndexInfo, 0)
	indexCount, err := is.(*ProtocolDataInputStream).ReadInt() // index count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexCount as '%+v'", indexCount))
	}
	for i := 0; i < indexCount; i++ {
		indexType, err := is.(*ProtocolDataInputStream).ReadByte() // Index Type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexType from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexType as '%+v'", indexType))
		}
		indexId, err := is.(*ProtocolDataInputStream).ReadInt() // Index Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexId from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexId as '%+v'", indexId))
		}
		indexName, err := is.(*ProtocolDataInputStream).ReadUTF() // Index Name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexName from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read indexName as '%+v'", indexName))
		}
		isUniqueFlag, err := is.(*ProtocolDataInputStream).ReadBoolean() // isUnique flag
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading isUniqueFlag from message buffer"))
			return nil, err
		}

		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read isUniqueFlag as '%+v'", isUniqueFlag))
		}
		attrCount, err := is.(*ProtocolDataInputStream).ReadInt() // attribute count
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrCount from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrCount as '%+v'", attrCount))
		}
		attributes := make([]string, 0)
		for j := 0; j < attrCount; j++ {
			attrName, err := is.(*ProtocolDataInputStream).ReadUTF() // attribute Name
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrName from message buffer"))
				return nil, err
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read attrName as '%+v'", attrName))
			}
			attributes = append(attributes, attrName)
		}

		nodeCount, err := is.(*ProtocolDataInputStream).ReadInt() // node count
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nodeCount from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nodeCount as '%+v'", nodeCount))
		}
		nodes := make([]string, 0)
		for k := 0; k < nodeCount; k++ {
			nodeName, err := is.(*ProtocolDataInputStream).ReadUTF() // node Name
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nodeName from message buffer"))
				return nil, err
			}
			if logger.IsDebug() {
				logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nodeName as '%+v'", nodeName))
			}
			nodes = append(nodes, nodeName)
		}

		blocksize, err := is.(*ProtocolDataInputStream).ReadInt() // block size
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading blocksize from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read blocksize as '%+v'", blocksize))
		}

		numEntries, err := is.(*ProtocolDataInputStream).ReadLong() // num of Entries
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numEntries from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read numEntries as '%+v'", numEntries))
		}

		status, err := is.(*ProtocolDataInputStream).ReadBytes() // status
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading status from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read status as '%+v'", status))
		}

		indexInfo := NewIndexInfoImpl(indexId, indexName, indexType, isUniqueFlag, attributes, nodes, numEntries, string(status))
		indexList = append(indexList, indexInfo)
	}
	return indexList, nil
}

func extractMemoryInfoFromInputStream(is tgdb.TGInputStream) (*ServerMemoryInfoImpl, tgdb.TGError) {
	freeProcessMemory, err := is.(*ProtocolDataInputStream).ReadLong() // free process memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading freeProcessMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read freeProcessMemory as '%+v'", freeProcessMemory))
	}
	memProcessUsagePct, err := is.(*ProtocolDataInputStream).ReadInt() // process memory usage percentage
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memProcessUsagePct from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memProcessUsagePct as '%+v'", memProcessUsagePct))
	}
	maxProcessMemory, err := is.(*ProtocolDataInputStream).ReadLong() // max process memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxProcessMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxProcessMemory as '%+v'", maxProcessMemory))
	}
	usedProcessMemory := maxProcessMemory - freeProcessMemory

	fileLocation, err := is.(*ProtocolDataInputStream).ReadUTF() // shared file location
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading fileLocation from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read fileLocation as '%+v'", fileLocation))
	}
	freeSharedMemory, err := is.(*ProtocolDataInputStream).ReadLong() // free shared memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading freeSharedMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read freeSharedMemory as '%+v'", freeSharedMemory))
	}
	memSharedUsagePct, err := is.(*ProtocolDataInputStream).ReadInt() // shared memory usage percentage
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memSharedUsagePct from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memSharedUsagePct as '%+v'", memSharedUsagePct))
	}
	maxSharedMemory, err := is.(*ProtocolDataInputStream).ReadLong() // max shared memory
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxSharedMemory from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxSharedMemory as '%+v'", maxSharedMemory))
	}
	usedSharedMemory := maxSharedMemory - freeSharedMemory

	processMemory := NewMemoryInfoImpl(freeProcessMemory, maxProcessMemory, usedProcessMemory, "")
	sharedMemory := NewMemoryInfoImpl(freeSharedMemory, maxSharedMemory, usedSharedMemory, fileLocation)

	memoryInfo := NewServerMemoryInfoImpl(processMemory, sharedMemory)
	return memoryInfo, nil
}

func extractNetListenersInfoFromInputStream(is tgdb.TGInputStream) ([]tgdb.TGNetListenerInfo, tgdb.TGError) {
	listenerList := make([]tgdb.TGNetListenerInfo, 0)
	listenerCount, err := is.(*ProtocolDataInputStream).ReadLong() // listener count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading listenerCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read listenerCount as '%+v'", listenerCount))
	}

	for i := 0; i < int(listenerCount); i++ {
		bufLen, err := is.(*ProtocolDataInputStream).ReadInt() // buffer length
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bufLen from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read bufLen as '%+v'", bufLen))
		}

		nameBytes := make([]byte, bufLen)
		_, err = is.(*ProtocolDataInputStream).ReadIntoBuffer(nameBytes) // Name bytes
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading nameBytes from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read nameBytes as '%+v'", nameBytes))
		}
		listenerName := string(nameBytes)

		currConn, err := is.(*ProtocolDataInputStream).ReadInt() // current connection Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading currConn from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read currConn as '%+v'", currConn))
		}

		maxConn, err := is.(*ProtocolDataInputStream).ReadInt() // maximum connections
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading maxConn from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read maxConn as '%+v'", maxConn))
		}

		bufLen4Port, err := is.(*ProtocolDataInputStream).ReadInt() // buffer length for port
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bufLen4Port from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read bufLen4Port as '%+v'", bufLen4Port))
		}

		portBytes := make([]byte, bufLen4Port)
		_, err = is.(*ProtocolDataInputStream).ReadIntoBuffer(portBytes) // port bytes
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading portBytes from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read portBytes as '%+v'", portBytes))
		}
		portNumber := string(portBytes)

		listenerInfo := NewNetListenerInfoImpl(currConn, maxConn, listenerName, portNumber)
		listenerList = append(listenerList, listenerInfo)
	}
	return listenerList, nil
}

func extractServerInfoFromInputStream(is tgdb.TGInputStream) (*ServerInfoImpl, tgdb.TGError) {
	serverStatus, err  := extractServerStatusFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading serverStatus from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read serverStatus as '%+v'", serverStatus))
	}

	memoryInfo, err  := extractMemoryInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading memoryInfo from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read memoryInfo as '%+v'", memoryInfo))
	}

	netListenerInfo, err  := extractNetListenersInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading netListenerInfo from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read netListenerInfo as '%+v'", netListenerInfo))
	}

	transactionInfo, err  := extractTransactionsInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading transactionInfo from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read transactionInfo as '%+v'", transactionInfo))
	}

	cacheInfo, err  := extractCacheInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading cacheInfo from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read cacheInfo as '%+v'", cacheInfo))
	}

	databaseInfo, err  := extractDatabaseInfoFromInputStream(is)
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading databaseInfo from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read databaseInfo as '%+v'", databaseInfo))
	}

	serverInfo := NewServerInfoImpl(cacheInfo, databaseInfo, memoryInfo, netListenerInfo, serverStatus, transactionInfo)
	return serverInfo, nil
}

func extractServerStatusFromInputStream(is tgdb.TGInputStream) (*ServerStatusImpl, tgdb.TGError) {
	name, err := is.(*ProtocolDataInputStream).ReadUTF() // Name
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading Name from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read Name as '%+v'", name))
	}

	status, err := is.(*ProtocolDataInputStream).ReadByte() // status
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading status from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read status as '%+v'", status))
	}

	processId, err := is.(*ProtocolDataInputStream).ReadInt() // process Id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading processId from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read processId as '%+v'", processId))
	}

	duration, err := is.(*ProtocolDataInputStream).ReadLong() // duration
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading duration from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read duration as '%+v'", duration))
	}

	sVersion, err := is.(*ProtocolDataInputStream).ReadLong() // version
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading sVersion from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read sVersion as '%+v'", sVersion))
	}

	sVersionInfo := NewTGServerVersion(sVersion)
	uptime := time.Duration(time.Now().Unix() - (duration/1000000000))
	serverInfo := NewServerStatusImpl(name, sVersionInfo, string(processId), tgdb.ServerStates(status), uptime)

	return serverInfo, nil
}

func extractTransactionsInfoFromInputStream(is tgdb.TGInputStream) (*TransactionStatisticsImpl, tgdb.TGError) {
	txnProcessorsCount, err := is.(*ProtocolDataInputStream).ReadShort() // txn processor count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnProcessorsCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnProcessorsCount as '%+v'", txnProcessorsCount))
	}

	txnProcessedCount, err := is.(*ProtocolDataInputStream).ReadLong() // txn processed count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnProcessoedCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnProcessedCount as '%+v'", txnProcessedCount))
	}

	txnSuccessCount, err := is.(*ProtocolDataInputStream).ReadLong() // txn success count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnSuccessCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnSuccessCount as '%+v'", txnSuccessCount))
	}

	avgProcessTime, err := is.(*ProtocolDataInputStream).ReadDouble() // average processing time
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading avgProcessTime from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read avgProcessTime as '%+v'", avgProcessTime))
	}

	pendingTxnCount, err := is.(*ProtocolDataInputStream).ReadLong() // pending txn count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading pendingTxnCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read pendingTxnCount as '%+v'", pendingTxnCount))
	}

	txnLogQueueDepth, err := is.(*ProtocolDataInputStream).ReadInt() // txn logger queue depth
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading txnLogQueueDepth from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read txnLogQueueDepth as '%+v'", txnLogQueueDepth))
	}

	transactionsInfo := NewTransactionStatisticsImpl(avgProcessTime, pendingTxnCount, txnLogQueueDepth, int64(txnProcessorsCount), txnProcessedCount, txnSuccessCount)
	return transactionsInfo, nil
}

func extractUserListFromInputStream(is tgdb.TGInputStream) ([]tgdb.TGUserInfo, tgdb.TGError) {
	userList := make([]tgdb.TGUserInfo, 0)
	userCount, err := is.(*ProtocolDataInputStream).ReadInt() // user count
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userCount from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userCount as '%+v'", userCount))
	}

	for i := 0; i < userCount; i++ {
		userType, err := is.(*ProtocolDataInputStream).ReadByte() // user Type
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userType from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userType as '%+v'", userType))
		}

		userId, err := is.(*ProtocolDataInputStream).ReadInt() // user Id
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userId from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userId as '%+v'", userId))
		}

		userName, err := is.(*ProtocolDataInputStream).ReadUTF() // user Name
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userName from message buffer"))
			return nil, err
		}
		if logger.IsDebug() {
					logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read userName as '%+v'", userName))
		}

		switch tgdb.TGSystemType(userType) {
		case tgdb.SystemTypePrincipal:
			_, err := is.(*ProtocolDataInputStream).ReadBytes()
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bytes from message buffer"))
				return nil, err
			}
			numRoles, err := is.(*ProtocolDataInputStream).ReadInt() //number of roles
			if err != nil {
				logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading numRoles from message buffer"))
				return nil, err
			}

			roleIds := make([]int, numRoles)
			for j := 0; j < numRoles; j++ {
				roleId, err := is.(*ProtocolDataInputStream).ReadInt() //roleID
				if err != nil {
					logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading roleID from message buffer"))
					return nil, err
				}
				roleIds = append(roleIds, roleId)
			}
		default:
		}
		userInfo := NewUserInfoImpl(userId, userName, userType)
		userList = append(userList, userInfo)
	}
	return userList, nil
}


type AdminRequestMessage struct {
	*AbstractProtocolMessage
	command    AdminCommand
	logDetails *ServerLogDetails
}

func DefaultAdminRequestMessage() *AdminRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AdminRequestMessage{})

	newMsg := AdminRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
		command:                 AdminCommandInvalid,
	}
	newMsg.HeaderSetVerbId(VerbAdminRequest)
	newMsg.HeaderSetMessageByteBufLength(int(reflect.TypeOf(newMsg).Size()))
	newMsg.SetUpdatableFlag(true)
	return &newMsg
}

// Create New Message Instance
func NewAdminRequestMessage(authToken, sessionId int64) *AdminRequestMessage {
	newMsg := DefaultAdminRequestMessage()
	newMsg.HeaderSetAuthToken(authToken)
	newMsg.HeaderSetSessionId(sessionId)
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *AdminRequestMessage) GetCommand() AdminCommand {
	return msg.command
}

func (msg *AdminRequestMessage) GetLogLevel() *ServerLogDetails {
	return msg.logDetails
}

func (msg *AdminRequestMessage) SetCommand(cmd AdminCommand) {
	msg.command = cmd
}

func (msg *AdminRequestMessage) SetLogLevel(connId *ServerLogDetails) {
	msg.logDetails = connId
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AdminRequestMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminRequest:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminRequest:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminRequest:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminRequest:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("AdminRequest::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AdminRequestMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminRequest:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminRequest:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminRequest:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminRequest:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("AdminRequest::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AdminRequestMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AdminRequestMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AdminRequestMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AdminRequestMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AdminRequestMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AdminRequestMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AdminRequestMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AdminRequestMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AdminRequestMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AdminRequestMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AdminRequestMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AdminRequestMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AdminRequestMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AdminRequestMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AdminRequestMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AdminRequestMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.GetVerbId()).GetName())
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AdminRequestMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AdminRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AdminRequest:{")
	buffer.WriteString(fmt.Sprintf("Command: %d", msg.command))
	buffer.WriteString(fmt.Sprintf(", ServerLogDetails: %+v", msg.logDetails))
	strArray := []string{buffer.String(), msg.APMMessageToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AdminRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AdminRequestMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AdminRequestMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AdminRequestMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AdminRequestMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	//startPos := os.GetPosition()
	//logger.Log(fmt.Sprintf("Entering AdminRequest:WritePayload at output buffer position: '%d'", startPos))
	//os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the commit buffer length
	//os.(*iostream.ProtocolDataOutputStream).WriteInt(0) // This is for the checksum for the commit buffer to be added later.  Currently not used

	dataLen := 0
	checkSum := 0

	switch msg.command {
	case AdminCommandCreateUser:
	case AdminCommandCreateAttrDesc:
	case AdminCommandCreateIndex:
	case AdminCommandCreateNodeType:
	case AdminCommandCreateEdgeType:
	case AdminCommandShowUsers:
		fallthrough
	case AdminCommandShowAttrDescs:
		fallthrough
	case AdminCommandShowIndices:
		fallthrough
	case AdminCommandShowTypes:
		fallthrough
	case AdminCommandShowInfo:
		fallthrough
	case AdminCommandShowConnections:
		os.(*ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*ProtocolDataOutputStream).WriteInt(int(msg.command))
	case AdminCommandDescribe:
	case AdminCommandSetLogLevel:
		os.(*ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*ProtocolDataOutputStream).WriteInt(int(msg.command))
		os.(*ProtocolDataOutputStream).WriteShort(int(msg.logDetails.GetLogLevel()))
		os.(*ProtocolDataOutputStream).WriteLong(int64(msg.logDetails.GetLogComponent()))
	case AdminCommandStopServer:
		fallthrough
	case AdminCommandCheckpointServer:
		os.(*ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*ProtocolDataOutputStream).WriteInt(int(msg.command))
	case AdminCommandDisconnectClient:
	case AdminCommandKillConnection:
		os.(*ProtocolDataOutputStream).WriteInt(dataLen)
		os.(*ProtocolDataOutputStream).WriteInt(checkSum)
		os.(*ProtocolDataOutputStream).WriteInt(int(msg.command))
		os.(*ProtocolDataOutputStream).WriteLong(msg.GetSessionId())
		os.(*ProtocolDataOutputStream).WriteBoolean(true)
	default:
	}

	//currPos := os.GetPosition()
	//length := currPos - startPos
	//_, err := os.(*iostream.ProtocolDataOutputStream).WriteIntAt(startPos, length)
	//if err != nil {
	//	return err
	//}
	//logger.Log(fmt.Sprintf("Returning AdminRequest::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AdminRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.HeaderGetMessageByteBufLength(), msg.HeaderGetVerbId(),
		msg.HeaderGetSequenceNo(), msg.HeaderGetTimestamp(), msg.HeaderGetRequestId(),
		msg.HeaderGetDataOffset(), msg.HeaderGetAuthToken(), msg.HeaderGetSessionId(), msg.GetIsUpdatable(),
		msg.command, msg.logDetails)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminRequest:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AdminRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	var bLen, vId int
	var seq, tStamp, reqId, token, sId int64
	var offset int16
	var uFlag bool
	_, err := fmt.Fscanln(b, &bLen, &vId, &seq, &tStamp, &reqId, &offset, &token, &sId, &uFlag, &msg.command, &msg.logDetails)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminRequest:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	msg.HeaderSetMessageByteBufLength(bLen)
	msg.HeaderSetVerbId(vId)
	msg.HeaderSetSequenceNo(seq)
	msg.HeaderSetTimestamp(tStamp)
	msg.HeaderSetRequestId(reqId)
	msg.HeaderSetAuthToken(token)
	msg.HeaderSetSessionId(sId)
	msg.HeaderSetDataOffset(offset)
	msg.SetIsUpdatable(uFlag)
	return nil
}



type AdminResponseMessage struct {
	*AbstractProtocolMessage
	attrDescriptors []tgdb.TGAttributeDescriptor
	connections     []tgdb.TGConnectionInfo
	indices         []tgdb.TGIndexInfo
	serverInfo      *ServerInfoImpl
	users           []tgdb.TGUserInfo
}

func DefaultAdminResponseMessage() *AdminResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(AdminResponseMessage{})

	newMsg := AdminResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.HeaderSetVerbId(VerbAdminResponse)
	newMsg.HeaderSetMessageByteBufLength(int(reflect.TypeOf(newMsg).Size()))
	return &newMsg
}

// Create New Message Instance
func NewAdminResponseMessage(authToken, sessionId int64) *AdminResponseMessage {
	newMsg := DefaultAdminResponseMessage()
	newMsg.HeaderSetAuthToken(authToken)
	newMsg.HeaderSetSessionId(sessionId)
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for AdminResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *AdminResponseMessage) GetDescriptorList() []tgdb.TGAttributeDescriptor {
	return msg.attrDescriptors
}

func (msg *AdminResponseMessage) GetConnectionList() []tgdb.TGConnectionInfo {
	return msg.connections
}

func (msg *AdminResponseMessage) GetIndexList() []tgdb.TGIndexInfo {
	return msg.indices
}

func (msg *AdminResponseMessage) GetServerInfo() *ServerInfoImpl {
	return msg.serverInfo
}

func (msg *AdminResponseMessage) GetUserList() []tgdb.TGUserInfo {
	return msg.users
}

func (msg *AdminResponseMessage) SetDescriptorList(list []tgdb.TGAttributeDescriptor) {
	msg.attrDescriptors = list
}

func (msg *AdminResponseMessage) SetConnectionList(list []tgdb.TGConnectionInfo) {
	msg.connections = list
}

func (msg *AdminResponseMessage) SetIndexList(list []tgdb.TGIndexInfo) {
	msg.indices = list
}

func (msg *AdminResponseMessage) SetServerInfo(sInfo *ServerInfoImpl) {
	msg.serverInfo = sInfo
}

func (msg *AdminResponseMessage) SetUserList(list []tgdb.TGUserInfo) {
	msg.users = list
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *AdminResponseMessage) FromBytes(buffer []byte) (tgdb.TGMessage, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminResponse:FromBytes"))
	}
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, CreateExceptionByType(TGErrorInvalidMessageLength)
	}

	is := NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponse:FromBytes read bufLen as '%+v'", bufLen))
	}
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, GetErrorByType(TGErrorInvalidMessageLength, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminResponse:FromBytes - about to APMReadHeader"))
	}
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminResponse:FromBytes - about to ReadPayload"))
	}
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("AdminResponse::FromBytes resulted in '%+v'", msg))
	}
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *AdminResponseMessage) ToBytes() ([]byte, int, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminResponse:ToBytes"))
	}
	os := DefaultProtocolDataOutputStream()

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminResponse:ToBytes - about to APMWriteHeader"))
	}
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Inside AdminResponse:ToBytes - about to WritePayload"))
	}
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, GetErrorByType(TGErrorIOException, INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponse:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("AdminResponse::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	}
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *AdminResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *AdminResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *AdminResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *AdminResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *AdminResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *AdminResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *AdminResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *AdminResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *AdminResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *AdminResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *AdminResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *AdminResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *AdminResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *AdminResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *AdminResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *AdminResponseMessage) SetTimestamp(timestamp int64) tgdb.TGError {
	if !(msg.GetIsUpdatable() || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.GetVerbId()).GetName())
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *AdminResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *AdminResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AdminResponse:{")
	buffer.WriteString(fmt.Sprintf("AttributeDescriptors: %+v", msg.attrDescriptors))
	buffer.WriteString(fmt.Sprintf(", Connections: %+v", msg.connections))
	buffer.WriteString(fmt.Sprintf(", Indices: %+v", msg.indices))
	buffer.WriteString(fmt.Sprintf(", ServerInfo: %+v", msg.serverInfo))
	buffer.WriteString(fmt.Sprintf(", Users: %+v", msg.users))
	strArray := []string{buffer.String(), msg.APMMessageToString() + "}"}
	msgStr := strings.Join(strArray, ", ")
	return msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *AdminResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) tgdb.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *AdminResponseMessage) ReadHeader(is tgdb.TGInputStream) tgdb.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header Attributes to output stream
func (msg *AdminResponseMessage) WriteHeader(os tgdb.TGOutputStream) tgdb.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload Attributes
func (msg *AdminResponseMessage) ReadPayload(is tgdb.TGInputStream) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering AdminResponseMessage:ReadPayload"))
	}
	//bLen, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // buf length
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading bLen from message buffer"))
	//	return err
	//}
	//logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read buf length as '%+v'", bLen))
	//
	//checkSum, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // checksum
	//if err != nil {
	//	logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading checkSum from message buffer"))
	//	return err
	//}
	//logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read checkSum as '%+v'", checkSum))

	resultId, err := is.(*ProtocolDataInputStream).ReadInt() // result Id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading resultId from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read resultId as '%+v'", resultId))
	}

	command, err := is.(*ProtocolDataInputStream).ReadInt() // command
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading command from message buffer"))
		return err
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside AdminResponseMessage:ReadPayload read command as '%+v'", command))
	}

	switch AdminCommand(command) {
	case AdminCommandCreateUser:
	case AdminCommandCreateAttrDesc:
	case AdminCommandCreateIndex:
	case AdminCommandCreateNodeType:
	case AdminCommandCreateEdgeType:
	case AdminCommandShowUsers:
		userList, err := extractUserListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading userList from message buffer"))
			return err
		}
		msg.SetUserList(userList)
	case AdminCommandShowAttrDescs:
		attrDescList, err := extractDescriptorListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading attrDescList from message buffer"))
			return err
		}
		msg.SetDescriptorList(attrDescList)
	case AdminCommandShowIndices:
		indexList, err := extractIndexListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading indexList from message buffer"))
			return err
		}
		msg.SetIndexList(indexList)
	case AdminCommandShowTypes:
	case AdminCommandShowInfo:
		serverInfo, err := extractServerInfoFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading serverInfo from message buffer"))
			return err
		}
		msg.SetServerInfo(serverInfo)
	case AdminCommandShowConnections:
		connList, err := extractConnectionListFromInputStream(is)
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning AdminResponseMessage:ReadPayload w/ Error in reading connList from message buffer"))
			return err
		}
		msg.SetConnectionList(connList)
	case AdminCommandDescribe:
	case AdminCommandSetLogLevel:
	case AdminCommandStopServer:
	case AdminCommandCheckpointServer:
	case AdminCommandDisconnectClient:
	case AdminCommandKillConnection:
	default:
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning AdminResponseMessage:ReadPayload"))
	}
	return nil
}

// WritePayload exports the values of the message specific payload Attributes to output stream
func (msg *AdminResponseMessage) WritePayload(os tgdb.TGOutputStream) tgdb.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *AdminResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.HeaderGetMessageByteBufLength(), msg.HeaderGetVerbId(),
		msg.HeaderGetSequenceNo(), msg.HeaderGetTimestamp(), msg.HeaderGetRequestId(),
		msg.HeaderGetDataOffset(), msg.HeaderGetAuthToken(), msg.HeaderGetSessionId(),
		msg.GetIsUpdatable(), msg.attrDescriptors, msg.connections, msg.indices, msg.serverInfo, msg.users)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminResponse:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *AdminResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	var bLen, vId int
	var seq, tStamp, reqId, token, sId int64
	var offset int16
	var uFlag bool
	_, err := fmt.Fscanln(b, &bLen, &vId, &seq, &tStamp, &reqId, &offset, &token, &sId, &uFlag,
		&msg.attrDescriptors, &msg.connections, &msg.indices, &msg.serverInfo, &msg.users)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning AdminResponse:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	msg.HeaderSetMessageByteBufLength(bLen)
	msg.HeaderSetVerbId(vId)
	msg.HeaderSetSequenceNo(seq)
	msg.HeaderSetTimestamp(tStamp)
	msg.HeaderSetRequestId(reqId)
	msg.HeaderSetAuthToken(token)
	msg.HeaderSetSessionId(sId)
	msg.HeaderSetDataOffset(offset)
	msg.SetIsUpdatable(uFlag)
	return nil
}

//var logger = logging.DefaultTGLogManager().GetLogger()

type CacheStatisticsImpl struct {
	dataCacheMaxEntries  int
	dataCacheEntries     int
	dataCacheHits        int64
	dataCacheMisses      int64
	dataCacheMaxMemory   int64
	indexCacheMaxEntries int
	indexCacheEntries    int
	indexCacheHits       int64
	indexCacheMisses     int64
	indexCacheMaxMemory  int64
}

// Make sure that the CacheStatisticsImpl implements the TGCacheStatistics interface
var _ tgdb.TGCacheStatistics = (*CacheStatisticsImpl)(nil)

func DefaultCacheStatisticsImpl() *CacheStatisticsImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(CacheStatisticsImpl{})

	return &CacheStatisticsImpl{}
}

func NewCacheStatisticsImpl(_dataCacheMaxEntries, _dataCacheEntries int,
	_dataCacheHits, _dataCacheMisses, _dataCacheMaxMemory int64,
	_indexCacheMaxEntries, _indexCacheEntries int,
	_indexCacheHits, _indexCacheMisses, _indexCacheMaxMemory int64) *CacheStatisticsImpl {
	newCacheStatistics := DefaultCacheStatisticsImpl()
	newCacheStatistics.dataCacheMaxEntries = _dataCacheMaxEntries
	newCacheStatistics.dataCacheEntries = _dataCacheEntries
	newCacheStatistics.dataCacheHits = _dataCacheHits
	newCacheStatistics.dataCacheMisses = _dataCacheMisses
	newCacheStatistics.dataCacheMaxMemory = _dataCacheMaxMemory
	newCacheStatistics.indexCacheMaxEntries = _indexCacheMaxEntries
	newCacheStatistics.indexCacheEntries = _indexCacheEntries
	newCacheStatistics.indexCacheHits = _indexCacheHits
	newCacheStatistics.indexCacheMisses = _indexCacheMisses
	newCacheStatistics.indexCacheMaxMemory = _indexCacheMaxMemory
	return newCacheStatistics
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGCacheStatisticsImpl
/////////////////////////////////////////////////////////////////

func (obj *CacheStatisticsImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CacheStatisticsImpl:{")
	buffer.WriteString(fmt.Sprintf("DataCacheMaxEntries: '%d'", obj.dataCacheMaxEntries))
	buffer.WriteString(fmt.Sprintf(", DataCacheEntries: '%d'", obj.dataCacheEntries))
	buffer.WriteString(fmt.Sprintf(", DataCacheHits: '%d'", obj.dataCacheHits))
	buffer.WriteString(fmt.Sprintf(", DataCacheMisses: '%d'", obj.dataCacheMisses))
	buffer.WriteString(fmt.Sprintf(", DataCacheMaxMemory: '%d'", obj.dataCacheMaxMemory))
	buffer.WriteString(fmt.Sprintf(", IndexCacheMaxEntries: '%d'", obj.indexCacheMaxEntries))
	buffer.WriteString(fmt.Sprintf(", IndexCacheEntries: '%d'", obj.indexCacheEntries))
	buffer.WriteString(fmt.Sprintf(", IndexCacheHits: '%d'", obj.indexCacheHits))
	buffer.WriteString(fmt.Sprintf(", IndexCacheMisses: '%d'", obj.indexCacheMisses))
	buffer.WriteString(fmt.Sprintf(", IndexCacheMaxMemory: '%d'", obj.indexCacheMaxMemory))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGCacheStatistics
/////////////////////////////////////////////////////////////////

// GetDataCacheEntries returns the data-cache entries
func (obj *CacheStatisticsImpl) GetDataCacheEntries() int {
	return obj.dataCacheEntries
}

// GetDataCacheHits returns the data-cache hits
func (obj *CacheStatisticsImpl) GetDataCacheHits() int64 {
	return obj.dataCacheHits
}

// GetDataCacheMisses returns the data-cache misses
func (obj *CacheStatisticsImpl) GetDataCacheMisses() int64 {
	return obj.dataCacheMisses
}

// GetDataCacheMaxEntries returns the data-cache max entries
func (obj *CacheStatisticsImpl) GetDataCacheMaxEntries() int {
	return obj.dataCacheMaxEntries
}

// GetDataCacheMaxMemory returns the data-cache max memory
func (obj *CacheStatisticsImpl) GetDataCacheMaxMemory() int64 {
	return obj.dataCacheMaxMemory
}

// GetIndexCacheEntries returns the index-cache entries
func (obj *CacheStatisticsImpl) GetIndexCacheEntries() int {
	return obj.indexCacheEntries
}

// GetIndexCacheHits returns the index-cache hits
func (obj *CacheStatisticsImpl) GetIndexCacheHits() int64 {
	return obj.indexCacheHits
}

// GetIndexCacheMisses returns the index-cache misses
func (obj *CacheStatisticsImpl) GetIndexCacheMisses() int64 {
	return obj.indexCacheMisses
}

// GetIndexCacheMaxMemory returns the index-cache max memory
func (obj *CacheStatisticsImpl) GetIndexCacheMaxMemory() int64 {
	return obj.indexCacheMaxMemory
}

// GetIndexCacheMaxEntries returns the index-cache max entries
func (obj *CacheStatisticsImpl) GetIndexCacheMaxEntries() int {
	return obj.indexCacheMaxEntries
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *CacheStatisticsImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.dataCacheMaxEntries, obj.dataCacheEntries, obj.dataCacheHits,
		obj.dataCacheMisses, obj.dataCacheMaxMemory, obj.indexCacheMaxEntries, obj.indexCacheEntries,
		obj.indexCacheHits, obj.indexCacheMisses, obj.indexCacheMaxMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CacheStatisticsImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *CacheStatisticsImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.dataCacheMaxEntries, &obj.dataCacheEntries, &obj.dataCacheHits,
		&obj.dataCacheMisses, &obj.dataCacheMaxMemory, &obj.indexCacheMaxEntries, &obj.indexCacheEntries,
		&obj.indexCacheHits, &obj.indexCacheMisses, &obj.indexCacheMaxMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning CacheStatisticsImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ConnectionInfoImpl struct {
	ListenerName         string
	ClientID             string
	SessionID            int64
	UserName             string
	RemoteAddress        string
	CreatedTimeInSeconds int64
}

// Make sure that the ConnectionInfoImpl implements the TGConnectionInfo interface
var _ tgdb.TGConnectionInfo = (*ConnectionInfoImpl)(nil)

func DefaultConnectionInfoImpl() *ConnectionInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ConnectionInfoImpl{})

	return &ConnectionInfoImpl{}
}

func NewConnectionInfoImpl(_listnerName, _clientID string, _sessionID int64,
	_userName, _remoteAddress string, _createdTimeInSeconds int64) *ConnectionInfoImpl {
	newConnectionInfo := DefaultConnectionInfoImpl()
	newConnectionInfo.ListenerName = _listnerName
	newConnectionInfo.ClientID = _clientID
	newConnectionInfo.SessionID = _sessionID
	newConnectionInfo.UserName = _userName
	newConnectionInfo.RemoteAddress = _remoteAddress
	newConnectionInfo.CreatedTimeInSeconds = _createdTimeInSeconds
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGConnectionInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ConnectionInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("ListenerName: '%s'", obj.ListenerName))
	buffer.WriteString(fmt.Sprintf(", ClientID: '%s'", obj.ClientID))
	buffer.WriteString(fmt.Sprintf(", SessionID: '%d'", obj.SessionID))
	buffer.WriteString(fmt.Sprintf(", UserName: '%s'", obj.UserName))
	buffer.WriteString(fmt.Sprintf(", RemoteAddress: '%s'", obj.RemoteAddress))
	buffer.WriteString(fmt.Sprintf(", CreatedTimeInSeconds: '%d'", obj.CreatedTimeInSeconds))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGConnectionInfo
/////////////////////////////////////////////////////////////////

// GetClientID returns a client ID of listener
func (obj *ConnectionInfoImpl) GetClientID() string {
	return obj.ClientID
}

// GetCreatedTimeInSeconds returns a time when the listener was created
func (obj *ConnectionInfoImpl) GetCreatedTimeInSeconds() int64 {
	return obj.CreatedTimeInSeconds
}

// GetListenerName returns a Name of a particular listener
func (obj *ConnectionInfoImpl) GetListenerName() string {
	return obj.ListenerName
}

// GetRemoteAddress returns a remote address of listener
func (obj *ConnectionInfoImpl) GetRemoteAddress() string {
	return obj.RemoteAddress
}

// GetSessionID returns a session ID of listener
func (obj *ConnectionInfoImpl) GetSessionID() int64 {
	return obj.SessionID
}

// GetUserName returns a user-Name associated with listener
func (obj *ConnectionInfoImpl) GetUserName() string {
	return obj.UserName
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.ListenerName, obj.ClientID, obj.SessionID,
		obj.UserName, obj.RemoteAddress, obj.CreatedTimeInSeconds)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ConnectionInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ConnectionInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.ListenerName, &obj.ClientID, &obj.SessionID,
		&obj.UserName, &obj.RemoteAddress, &obj.CreatedTimeInSeconds)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ConnectionInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type DatabaseStatisticsImpl struct {
	blockSize        int
	dataBlockSize    int
	dataFree         int64
	dataSize         int64
	dataUsed         int64
	dbSize           int64
	indexFree        int64
	indexSize        int64
	indexUsed        int64
	numDataSegments  int
	numIndexSegments int
}

// Make sure that the DatabaseStatisticsImpl implements the TGDatabaseStatistics interface
var _ tgdb.TGDatabaseStatistics = (*DatabaseStatisticsImpl)(nil)

func DefaultDatabaseStatisticsImpl() *DatabaseStatisticsImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(DatabaseStatisticsImpl{})

	return &DatabaseStatisticsImpl{}
}

func NewDatabaseStatisticsImpl(_dbSize int64, _numDataSegments int, _dataSize, _dataUsed, _dataFree int64,
	_dataBlockSize, _numIndexSegments int, _indexSize, _indexUsed, _indexFree int64, _blockSize int) *DatabaseStatisticsImpl {
	newConnectionInfo := DefaultDatabaseStatisticsImpl()
	newConnectionInfo.blockSize = _blockSize
	newConnectionInfo.dataBlockSize = _dataBlockSize
	newConnectionInfo.dataFree = _dataFree
	newConnectionInfo.dataSize = _dataSize
	newConnectionInfo.dataUsed = _dataUsed
	newConnectionInfo.dbSize = _dbSize
	newConnectionInfo.indexFree = _indexFree
	newConnectionInfo.indexSize = _indexSize
	newConnectionInfo.indexUsed = _indexUsed
	newConnectionInfo.numDataSegments = _numDataSegments
	newConnectionInfo.numIndexSegments = _numIndexSegments
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGDatabaseStatisticsImpl
/////////////////////////////////////////////////////////////////

func (obj *DatabaseStatisticsImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DatabaseStatisticsImpl:{")
	buffer.WriteString(fmt.Sprintf("BlockSize: '%d'", obj.blockSize))
	buffer.WriteString(fmt.Sprintf(", DataBlockSize: '%d'", obj.dataBlockSize))
	buffer.WriteString(fmt.Sprintf(", DataFree: '%d'", obj.dataFree))
	buffer.WriteString(fmt.Sprintf(", DataSize: '%d'", obj.dataSize))
	buffer.WriteString(fmt.Sprintf(", DataUsed: '%d'", obj.dataUsed))
	buffer.WriteString(fmt.Sprintf(", DbSize: '%d'", obj.dbSize))
	buffer.WriteString(fmt.Sprintf(", IndexFree: '%d'", obj.indexFree))
	buffer.WriteString(fmt.Sprintf(", IndexSize: '%d'", obj.indexSize))
	buffer.WriteString(fmt.Sprintf(", IndexUsed: '%d'", obj.indexUsed))
	buffer.WriteString(fmt.Sprintf(", NumDataSegments: '%d'", obj.numDataSegments))
	buffer.WriteString(fmt.Sprintf(", NumIndexSegments: '%d'", obj.numIndexSegments))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGDatabaseStatistics
/////////////////////////////////////////////////////////////////

// GetBlockSize returns the block size
func (obj *DatabaseStatisticsImpl) GetBlockSize() int {
	return obj.blockSize
}

// GetDataBlockSize returns the block size of data
func (obj *DatabaseStatisticsImpl) GetDataBlockSize() int {
	return obj.dataBlockSize
}

// GetDataFree returns the free data size
func (obj *DatabaseStatisticsImpl) GetDataFree() int64 {
	return obj.dataFree
}

// GetDataSize returns data size
func (obj *DatabaseStatisticsImpl) GetDataSize() int64 {
	return obj.dataSize
}

// GetDataUsed returns the size of data used
func (obj *DatabaseStatisticsImpl) GetDataUsed() int64 {
	return obj.dataUsed
}

// GetDbSize returns the size of database
func (obj *DatabaseStatisticsImpl) GetDbSize() int64 {
	return obj.dbSize
}

// GetIndexFree returns the free index size
func (obj *DatabaseStatisticsImpl) GetIndexFree() int64 {
	return obj.indexFree
}

// GetIndexSize returns the index size
func (obj *DatabaseStatisticsImpl) GetIndexSize() int64 {
	return obj.indexSize
}

// GetIndexUsed returns the size of index used
func (obj *DatabaseStatisticsImpl) GetIndexUsed() int64 {
	return obj.indexUsed
}

// GetNumDataSegments returns the number of data segments
func (obj *DatabaseStatisticsImpl) GetNumDataSegments() int {
	return obj.numDataSegments
}

// GetNumIndexSegments returns the number of index segments
func (obj *DatabaseStatisticsImpl) GetNumIndexSegments() int {
	return obj.numIndexSegments
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *DatabaseStatisticsImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.blockSize, obj.dataBlockSize, obj.dataFree, obj.dataSize, obj.dataUsed,
		obj.dbSize, obj.indexFree, obj.indexSize, obj.indexUsed, obj.numDataSegments, obj.numIndexSegments)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DatabaseStatisticsImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *DatabaseStatisticsImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.blockSize, &obj.dataBlockSize, &obj.dataFree, &obj.dataSize, &obj.dataUsed,
		&obj.dbSize, &obj.indexFree, &obj.indexSize, &obj.indexUsed, &obj.numDataSegments, &obj.numIndexSegments)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning DatabaseStatisticsImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type IndexInfoImpl struct {
	sysId      int
	indexType  byte
	name       string
	uniqueFlag bool
	attributes []string
	nodeTypes  []string
	numEntries int64
	status     string
}

// Make sure that the IndexInfoImpl implements the TGIndexInfo interface
var _ tgdb.TGIndexInfo = (*IndexInfoImpl)(nil)

func DefaultIndexInfoImpl() *IndexInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(IndexInfoImpl{})

	return &IndexInfoImpl{}
}

func NewIndexInfoImpl(sysId int, name string, indexType byte,
	isUnique bool, attributes, nodeTypes []string, entries int64, status string) *IndexInfoImpl {
	newConnectionInfo := DefaultIndexInfoImpl()
	newConnectionInfo.sysId = sysId
	newConnectionInfo.indexType = indexType
	newConnectionInfo.name = name
	newConnectionInfo.uniqueFlag = isUnique
	newConnectionInfo.attributes = attributes
	newConnectionInfo.nodeTypes = nodeTypes
	newConnectionInfo.numEntries = entries
	newConnectionInfo.status = status
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGIndexInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("IndexInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("SysId: '%d'", obj.sysId))
	buffer.WriteString(fmt.Sprintf(", IndexType: '%+v'", obj.indexType))
	buffer.WriteString(fmt.Sprintf(", Name: '%s'", obj.name))
	buffer.WriteString(fmt.Sprintf(", IsUnique: '%+v'", obj.uniqueFlag))
	buffer.WriteString(fmt.Sprintf(", Attributes: '%+v'", obj.attributes))
	buffer.WriteString(fmt.Sprintf(", NodeTypes: '%+v'", obj.nodeTypes))
	buffer.WriteString(fmt.Sprintf(", NumEntries: '%+v'", obj.numEntries))
	buffer.WriteString(fmt.Sprintf(", Status: '%+v'", obj.status))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGIndexInfo
/////////////////////////////////////////////////////////////////

// GetAttributes returns a collection of attribute names
func (obj *IndexInfoImpl) GetAttributeNames() []string {
	return obj.attributes
}

// GetName returns the index Name
func (obj *IndexInfoImpl) GetName() string {
	return obj.name
}

// GetNumEntries returns the number of entries for the index
func (obj *IndexInfoImpl) GetNumEntries() int64 {
	return obj.numEntries
}

// GetType returns the index type
func (obj *IndexInfoImpl) GetType() byte {
	return obj.indexType
}

// GetStatus returns the status of the index
func (obj *IndexInfoImpl) GetStatus() string {
	return obj.status
}

// GetSystemId returns the system ID
func (obj *IndexInfoImpl) GetSystemId() int {
	return obj.sysId
}

// GetNodeTypes returns a collection of node types
func (obj *IndexInfoImpl) GetNodeTypes() []string {
	return obj.nodeTypes
}

// IsUnique returns the information whether the index is unique
func (obj *IndexInfoImpl) IsUnique() bool {
	return obj.uniqueFlag
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.sysId, obj.indexType, obj.name, obj.uniqueFlag, obj.attributes, obj.nodeTypes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IndexInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *IndexInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.sysId, &obj.indexType, &obj.name, &obj.uniqueFlag, &obj.attributes, &obj.nodeTypes)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning IndexInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type MemoryInfoImpl struct {
	freeMemory               int64
	maxMemory                int64
	usedMemory               int64
	sharedMemoryFileLocation string
}

// Make sure that the MemoryInfoImpl implements the TGMemoryInfo interface
var _ tgdb.TGMemoryInfo = (*MemoryInfoImpl)(nil)

func DefaultMemoryInfoImpl() *MemoryInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(MemoryInfoImpl{})

	return &MemoryInfoImpl{}
}

func NewMemoryInfoImpl(_freeMemory, _maxMemory, _usedMemory int64, _sharedMemoryFileLocation string) *MemoryInfoImpl {
	newConnectionInfo := DefaultMemoryInfoImpl()
	newConnectionInfo.freeMemory = _freeMemory
	newConnectionInfo.maxMemory = _maxMemory
	newConnectionInfo.usedMemory = _usedMemory
	newConnectionInfo.sharedMemoryFileLocation = _sharedMemoryFileLocation
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGMemoryInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MemoryInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("FreeMemory: '%d'", obj.freeMemory))
	buffer.WriteString(fmt.Sprintf(", MaxMemory: '%d'", obj.maxMemory))
	buffer.WriteString(fmt.Sprintf(", UsedMemory: '%d'", obj.usedMemory))
	buffer.WriteString(fmt.Sprintf(", SharedMemoryFileLocation: '%s'", obj.sharedMemoryFileLocation))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMemoryInfo
/////////////////////////////////////////////////////////////////

// GetFreeMemory returns the free memory size from server
func (obj *MemoryInfoImpl) GetFreeMemory() int64 {
	return obj.freeMemory
}

// GetMaxMemory returns the max memory size from server
func (obj *MemoryInfoImpl) GetMaxMemory() int64 {
	return obj.maxMemory
}

// GetSharedMemoryFileLocation returns the shared memory file location
func (obj *MemoryInfoImpl) GetSharedMemoryFileLocation() string {
	return obj.sharedMemoryFileLocation
}

// GetUsedMemory returns the used memory size from server
func (obj *MemoryInfoImpl) GetUsedMemory() int64 {
	return obj.usedMemory
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.freeMemory, obj.maxMemory, obj.usedMemory, obj.sharedMemoryFileLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MemoryInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *MemoryInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.freeMemory, &obj.maxMemory, &obj.usedMemory, &obj.sharedMemoryFileLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning MemoryInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type NetListenerInfoImpl struct {
	currentConnections int
	maxConnections     int
	listenerName       string
	portNumber         string
}

// Make sure that the NetListenerInfoImpl implements the TGNetListenerInfo interface
var _ tgdb.TGNetListenerInfo = (*NetListenerInfoImpl)(nil)

func DefaultNetListenerInfoImpl() *NetListenerInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(NetListenerInfoImpl{})

	return &NetListenerInfoImpl{}
}

func NewNetListenerInfoImpl(_currentConnections, _maxConnections int, _listenerName, _portNumber string) *NetListenerInfoImpl {
	newConnectionInfo := DefaultNetListenerInfoImpl()
	newConnectionInfo.currentConnections = _currentConnections
	newConnectionInfo.maxConnections = _maxConnections
	newConnectionInfo.listenerName = _listenerName
	newConnectionInfo.portNumber = _portNumber
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGNetListenerInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("NetListenerInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("CurrentConnections: '%d'", obj.currentConnections))
	buffer.WriteString(fmt.Sprintf(", MaxConnections: '%d'", obj.maxConnections))
	buffer.WriteString(fmt.Sprintf(", ListenerName: '%s'", obj.listenerName))
	buffer.WriteString(fmt.Sprintf(", PortNumber: '%s'", obj.portNumber))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGNetListenerInfo
/////////////////////////////////////////////////////////////////

// GetCurrentConnections returns the count of current connections
func (obj *NetListenerInfoImpl) GetCurrentConnections() int {
	return obj.currentConnections
}

// GetMaxConnections returns the count of max connections
func (obj *NetListenerInfoImpl) GetMaxConnections() int {
	return obj.maxConnections
}

// GetListenerName returns the listener Name
func (obj *NetListenerInfoImpl) GetListenerName() string {
	return obj.listenerName
}

// GetPortNumber returns the port detail of this listener
func (obj *NetListenerInfoImpl) GetPortNumber() string {
	return obj.portNumber
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.currentConnections, obj.maxConnections, obj.listenerName, obj.portNumber)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NetListenerInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *NetListenerInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.currentConnections, &obj.maxConnections, &obj.listenerName, &obj.portNumber)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning NetListenerInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type ServerInfoImpl struct {
	cacheInfo        *CacheStatisticsImpl
	databaseInfo     *DatabaseStatisticsImpl
	memoryInfo       *ServerMemoryInfoImpl
	netListenersInfo []tgdb.TGNetListenerInfo
	serverStatusInfo *ServerStatusImpl
	transactionsInfo *TransactionStatisticsImpl
}

// Make sure that the ServerInfoImpl implements the TGServerInfo interface
var _ tgdb.TGServerInfo = (*ServerInfoImpl)(nil)

func DefaultServerInfoImpl() *ServerInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerInfoImpl{})

	return &ServerInfoImpl{}
}

func NewServerInfoImpl(_cacheInfo *CacheStatisticsImpl, _databaseInfo *DatabaseStatisticsImpl,
	_memoryInfo *ServerMemoryInfoImpl, _netListenersInfo []tgdb.TGNetListenerInfo,
	_serverStatusInfo *ServerStatusImpl, _transactionsInfo *TransactionStatisticsImpl) *ServerInfoImpl {
	newServerInfo := DefaultServerInfoImpl()
	newServerInfo.cacheInfo = _cacheInfo
	newServerInfo.databaseInfo = _databaseInfo
	newServerInfo.memoryInfo = _memoryInfo
	newServerInfo.netListenersInfo = _netListenersInfo
	newServerInfo.serverStatusInfo = _serverStatusInfo
	newServerInfo.transactionsInfo = _transactionsInfo
	return newServerInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("CacheInfo: '%+v'", obj.cacheInfo))
	buffer.WriteString(fmt.Sprintf(", DatabaseInfo: '%+v'", obj.databaseInfo))
	buffer.WriteString(fmt.Sprintf(", MemoryInfo: '%+v'", obj.memoryInfo))
	buffer.WriteString(fmt.Sprintf(", NetListenersInfo: '%+v'", obj.netListenersInfo))
	buffer.WriteString(fmt.Sprintf(", ServerStatusInfo: '%+v'", obj.serverStatusInfo))
	buffer.WriteString(fmt.Sprintf(", TransactionsInfo: '%+v'", obj.transactionsInfo))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerInfo
/////////////////////////////////////////////////////////////////

// GetCacheInfo returns cache statistics information from server
func (obj *ServerInfoImpl) GetCacheInfo() tgdb.TGCacheStatistics {
	return obj.cacheInfo
}

// GetDatabaseInfo returns database statistics information from server
func (obj *ServerInfoImpl) GetDatabaseInfo() tgdb.TGDatabaseStatistics {
	return obj.databaseInfo
}

// GetMemoryInfo returns object corresponding to specific memory type
func (obj *ServerInfoImpl) GetMemoryInfo(memType tgdb.MemType) tgdb.TGMemoryInfo {
	return obj.memoryInfo.GetServerMemoryInfo(memType)
}

// GetNetListenersInfo returns a collection of information on NetListeners
func (obj *ServerInfoImpl) GetNetListenersInfo() []tgdb.TGNetListenerInfo {
	return obj.netListenersInfo
}

// GetServerStatus returns the information on Server Status including Name, version etc.
func (obj *ServerInfoImpl) GetServerStatus() tgdb.TGServerStatus {
	return obj.serverStatusInfo
}

// GetTransactionsInfo returns transaction statistics from server including processed transaction count, successful transaction count, average processing time etc.
func (obj *ServerInfoImpl) GetTransactionsInfo() tgdb.TGTransactionStatistics {
	return obj.transactionsInfo
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.cacheInfo, obj.databaseInfo, obj.memoryInfo,
		obj.netListenersInfo, obj.serverStatusInfo, obj.transactionsInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.cacheInfo, &obj.databaseInfo, &obj.memoryInfo,
		&obj.netListenersInfo, &obj.serverStatusInfo, &obj.transactionsInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


// ======= Server Log Level Types =======
type TGLogLevel int

const (
	TGLL_Console TGLogLevel = -2
	TGLL_Invalid TGLogLevel = -1
	TGLL_Fatal   TGLogLevel = iota
	TGLL_Error
	TGLL_Warn
	TGLL_Info
	TGLL_User
	TGLL_Debug
	TGLL_DebugFine
	TGLL_DebugFiner
	TGLL_MaxLogLevel
)

// ======= Server Log Component Types =======
type TGLogComponent int64

const (
	TGLC_COMMON_COREMEMORY TGLogComponent = iota
	TGLC_COMMON_CORECOLLECTIONS
	TGLC_COMMON_COREPLATFORM
	TGLC_COMMON_CORESTRING
	TGLC_COMMON_UTILS
	TGLC_COMMON_GRAPH
	TGLC_COMMON_MODEL
	TGLC_COMMON_NET
	TGLC_COMMON_PDU
	TGLC_COMMON_SEC
	TGLC_COMMON_FILES
	TGLC_COMMON_RESV2

	//Server Components
	TGLC_SERVER_CDMP
	TGLC_SERVER_DB
	TGLC_SERVER_EXPIMP
	TGLC_SERVER_INDEX
	TGLC_SERVER_INDEXBTREE
	TGLC_SERVER_INDEXISAM
	TGLC_SERVER_QUERY
	TGLC_SERVER_QUERY_RESV1
	TGLC_SERVER_QUERY_RESV2
	TGLC_SERVER_TXN
	TGLC_SERVER_TXNLOG
	TGLC_SERVER_TXNWRITER
	TGLC_SERVER_STORAGE
	TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_SERVER_GRAPH
	TGLC_SERVER_MAIN
	TGLC_SERVER_RESV2
	TGLC_SERVER_RESV3
	TGLC_SERVER_RESV4

	//Security Components
	TGLC_SECURITY_DATA
	TGLC_SECURITY_NET
	TGLC_SECURITY_RESV1
	TGLC_SECURITY_RESV2

	TGLC_ADMIN_LANG
	TGLC_ADMIN_CMD
	TGLC_ADMIN_MAIN
	TGLC_ADMIN_AST
	TGLC_ADMIN_GREMLIN

	TGLC_CUDA_GRAPHMGR
	TGLC_CUDA_KERNELEXECUTIVE
	TGLC_CUDA_RESV1
)

const (
	// User Defined Components
	TGLC_LOG_COREALL     TGLogComponent = TGLC_COMMON_COREMEMORY | TGLC_COMMON_CORECOLLECTIONS | TGLC_COMMON_COREPLATFORM | TGLC_COMMON_CORESTRING
	TGLC_LOG_GRAPHALL    TGLogComponent = TGLC_COMMON_GRAPH | TGLC_SERVER_GRAPH
	TGLC_LOG_MODEL       TGLogComponent = TGLC_COMMON_MODEL
	TGLC_LOG_NET         TGLogComponent = TGLC_COMMON_NET
	TGLC_LOG_PDUALL      TGLogComponent = TGLC_COMMON_PDU | TGLC_SERVER_CDMP
	TGLC_LOG_SECALL      TGLogComponent = TGLC_COMMON_SEC | TGLC_SECURITY_DATA | TGLC_SECURITY_NET
	TGLC_LOG_CUDAALL     TGLogComponent = TGLC_LOG_GRAPHALL | TGLC_CUDA_GRAPHMGR | TGLC_CUDA_KERNELEXECUTIVE
	TGLC_LOG_TXNALL      TGLogComponent = TGLC_SERVER_TXN | TGLC_SERVER_TXNLOG | TGLC_SERVER_TXNWRITER
	TGLC_LOG_STORAGEALL  TGLogComponent = TGLC_SERVER_STORAGE | TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_LOG_PAGEMANAGER TGLogComponent = TGLC_SERVER_STORAGEPAGEMANAGER
	TGLC_LOG_ADMINALL    TGLogComponent = TGLC_ADMIN_LANG | TGLC_ADMIN_CMD | TGLC_ADMIN_MAIN | TGLC_ADMIN_AST | TGLC_ADMIN_GREMLIN
	TGLC_LOG_MAIN        TGLogComponent = TGLC_SERVER_MAIN | TGLC_ADMIN_MAIN
)

const (
	TGLC_LOG_GLOBAL TGLogComponent = 0xFFFFFFFFFFFFFFF
)

type ServerLogDetails struct {
	logLevel     TGLogLevel
	logComponent TGLogComponent
}

func DefaultServerLogDetails() *ServerLogDetails {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerLogDetails{})

	return &ServerLogDetails{}
}

func NewServerLogDetails(logLevel TGLogLevel, _logComponent TGLogComponent) *ServerLogDetails {
	newServerLogDetails := DefaultServerLogDetails()
	newServerLogDetails.logLevel = logLevel
	newServerLogDetails.logComponent = _logComponent
	return newServerLogDetails
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerLogDetails
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerLogDetails:{")
	buffer.WriteString(fmt.Sprintf("LogLevel: '%+v'", obj.logLevel))
	buffer.WriteString(fmt.Sprintf(", LogComponent: '%+v'", obj.logComponent))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> ServerLogDetails
/////////////////////////////////////////////////////////////////

// GetLogLevel returns the server log level
func (obj *ServerLogDetails) GetLogLevel() TGLogLevel {
	return obj.logLevel
}

// GetLogComponent returns the server log component
func (obj *ServerLogDetails) GetLogComponent() TGLogComponent {
	return obj.logComponent
}

// SetLogLevel sets the server log level
func (obj *ServerLogDetails) SetLogLevel(level TGLogLevel) {
	obj.logLevel = level
}

// SetLogComponent sets the server log component
func (obj *ServerLogDetails) SetLogComponent(comp TGLogComponent) {
	obj.logComponent = comp
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.logLevel, obj.logComponent)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerLogDetails:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerLogDetails) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.logLevel, &obj.logComponent)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerLogDetails:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ServerMemoryInfoImpl struct {
	processMemory *MemoryInfoImpl
	sharedMemory  *MemoryInfoImpl
}

// Make sure that the ServerMemoryInfoImpl implements the TGServerMemoryInfo interface
var _ tgdb.TGServerMemoryInfo = (*ServerMemoryInfoImpl)(nil)

func DefaultServerMemoryInfoImpl() *ServerMemoryInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerMemoryInfoImpl{})

	return &ServerMemoryInfoImpl{}
}

func NewServerMemoryInfoImpl(_processMemory, _sharedMemory *MemoryInfoImpl) *ServerMemoryInfoImpl {
	newConnectionInfo := DefaultServerMemoryInfoImpl()
	newConnectionInfo.processMemory = _processMemory
	newConnectionInfo.sharedMemory = _sharedMemory
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerMemoryInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerMemoryInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("ProcessMemory: '%+v'", obj.processMemory))
	buffer.WriteString(fmt.Sprintf(", SharedMemory: '%+v'", obj.sharedMemory))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerMemoryInfo
/////////////////////////////////////////////////////////////////

// GetProcessMemory returns server process memory
func (obj *ServerMemoryInfoImpl) GetProcessMemory() tgdb.TGMemoryInfo {
	return obj.processMemory
}

// GetSharedMemory returns server shared memory
func (obj *ServerMemoryInfoImpl) GetSharedMemory() tgdb.TGMemoryInfo {
	return obj.sharedMemory
}

// GetMemoryInfo returns the memory info for the specified type
func (obj *ServerMemoryInfoImpl) GetServerMemoryInfo(memType tgdb.MemType) tgdb.TGMemoryInfo {
	if memType == tgdb.MemoryProcess {
		return obj.GetProcessMemory()
	} else if memType == tgdb.MemoryShared {
		return obj.GetSharedMemory()
	} else {
		return nil
	}
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.processMemory, obj.sharedMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerMemoryInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerMemoryInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.processMemory, &obj.sharedMemory)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerMemoryInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type ServerStatusImpl struct {
	name      string
	processId string
	status    tgdb.ServerStates
	uptime    time.Duration
	version   *TGServerVersion
}

// Make sure that the ServerStatusImpl implements the TGServerStatus interface
var _ tgdb.TGServerStatus = (*ServerStatusImpl)(nil)

func DefaultServerStatusImpl() *ServerStatusImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ServerStatusImpl{})

	return &ServerStatusImpl{}
}

func NewServerStatusImpl(_name string, _version *TGServerVersion, _processId string, _status tgdb.ServerStates, _uptime time.Duration) *ServerStatusImpl {
	newServerStatus := DefaultServerStatusImpl()
	newServerStatus.name = _name
	newServerStatus.processId = _processId
	newServerStatus.status = _status
	newServerStatus.uptime = _uptime
	newServerStatus.version = _version
	return newServerStatus
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGServerStatusImpl
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ServerStatusImpl:{")
	buffer.WriteString(fmt.Sprintf("Name: '%s'", obj.name))
	buffer.WriteString(fmt.Sprintf(", ProcessId: '%s'", obj.processId))
	buffer.WriteString(fmt.Sprintf(", Status: '%+v'", obj.status))
	buffer.WriteString(fmt.Sprintf(", Uptime: '%+v'", obj.uptime))
	buffer.WriteString(fmt.Sprintf(", Version: '%+v'", obj.version))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGServerStatus
/////////////////////////////////////////////////////////////////

// GetName returns the Name of the server instance
func (obj *ServerStatusImpl) GetName() string {
	return obj.name
}

// GetProcessId returns the process ID of server
func (obj *ServerStatusImpl) GetProcessId() string {
	return obj.processId
}

// GetServerStatus returns the state information of server
func (obj *ServerStatusImpl) GetServerStatus() tgdb.ServerStates {
	return obj.status
}

// GetUptime returns the uptime information of server
func (obj *ServerStatusImpl) GetUptime() time.Duration {
	return obj.uptime
}

// GetServerVersion returns the server version information
func (obj *ServerStatusImpl) GetServerVersion() *TGServerVersion {
	return obj.version
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.name, obj.processId, obj.status, obj.uptime, obj.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerStatusImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ServerStatusImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.name, &obj.processId, &obj.status, &obj.uptime, &obj.version)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ServerStatusImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

// ======= Admin Command Types =======
type AdminCommand int

const (
	AdminCommandInvalid AdminCommand = iota
	AdminCommandCreateUser
	AdminCommandCreateRole
	AdminCommandCreateAttrDesc
	AdminCommandCreateIndex
	AdminCommandCreateNodeType
	AdminCommandCreateEdgeType
	AdminCommandShowUsers
	AdminCommandShowRoles
	AdminCommandShowAttrDescs
	AdminCommandShowIndices
	AdminCommandShowTypes
	AdminCommandShowInfo
	AdminCommandShowConnections
	AdminCommandDescribe
	AdminCommandSetLogLevel
	AdminCommandSetSPDirectory
	AdminCommandUpdateRole
	AdminCommandStopServer
	AdminCommandCheckpointServer
	AdminCommandDisconnectClient
	AdminCommandKillConnection
)

func (command AdminCommand) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer
	buffer.WriteString("")

	if command&AdminCommandInvalid == AdminCommandInvalid {
		buffer.WriteString("Admin Command Invalid")
	} else if command&AdminCommandCreateUser == AdminCommandCreateUser {
		buffer.WriteString("Admin Command Create User")
	} else if command&AdminCommandCreateAttrDesc == AdminCommandCreateAttrDesc {
		buffer.WriteString("Admin Command Create AttrDesc")
	} else if command&AdminCommandCreateIndex == AdminCommandCreateIndex {
		buffer.WriteString("Admin Command Create Index")
	} else if command&AdminCommandCreateNodeType == AdminCommandCreateNodeType {
		buffer.WriteString("Admin Command Create NodeType")
	} else if command&AdminCommandCreateEdgeType == AdminCommandCreateEdgeType {
		buffer.WriteString("Admin Command Create EdgeType")
	} else if command&AdminCommandShowUsers == AdminCommandShowUsers {
		buffer.WriteString("Admin Command Show Users")
	} else if command&AdminCommandShowAttrDescs == AdminCommandShowAttrDescs {
		buffer.WriteString("Admin Command Show AttrDesc")
	} else if command&AdminCommandShowIndices == AdminCommandShowIndices {
		buffer.WriteString("Admin Command Show Indices")
	} else if command&AdminCommandShowTypes == AdminCommandShowTypes {
		buffer.WriteString("Admin Command Show Types")
	} else if command&AdminCommandShowInfo == AdminCommandShowInfo {
		buffer.WriteString("Admin Command Show Info")
	} else if command&AdminCommandShowConnections == AdminCommandShowConnections {
		buffer.WriteString("Admin Command Show Connections")
	} else if command&AdminCommandDescribe == AdminCommandDescribe {
		buffer.WriteString("Admin Command Describe")
	} else if command&AdminCommandSetLogLevel == AdminCommandSetLogLevel {
		buffer.WriteString("Admin Command Set LogLevel")
	} else if command&AdminCommandStopServer == AdminCommandStopServer {
		buffer.WriteString("Admin Command Stop Server")
	} else if command&AdminCommandCheckpointServer == AdminCommandCheckpointServer {
		buffer.WriteString("Admin Command Checkpoint Server")
	} else if command&AdminCommandDisconnectClient == AdminCommandDisconnectClient {
		buffer.WriteString("Admin Command Disconnect Client")
	} else if command&AdminCommandKillConnection == AdminCommandKillConnection {
		buffer.WriteString("Admin Command Kill Connection")
	}
	return buffer.String()
}




type TransactionStatisticsImpl struct {
	averageProcessingTime       float64
	pendingTransactionsCount    int64
	transactionLoggerQueueDepth int
	transactionProcessorCount   int64
	transactionProcessedCount   int64
	transactionSuccessfulCount  int64
}

// Make sure that the TransactionStatisticsImpl implements the TGTransactionStatistics interface
var _ tgdb.TGTransactionStatistics = (*TransactionStatisticsImpl)(nil)

func DefaultTransactionStatisticsImpl() *TransactionStatisticsImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TransactionStatisticsImpl{})

	return &TransactionStatisticsImpl{}
}

func NewTransactionStatisticsImpl(_averageProcessingTime float64, _pendingTransactionsCount int64, _transactionLoggerQueueDepth int,
	_transactionProcessorCount, _transactionProcessedCount, _transactionSuccessfulCount int64) *TransactionStatisticsImpl {
	newConnectionInfo := DefaultTransactionStatisticsImpl()
	newConnectionInfo.averageProcessingTime = _averageProcessingTime
	newConnectionInfo.pendingTransactionsCount = _pendingTransactionsCount
	newConnectionInfo.transactionLoggerQueueDepth = _transactionLoggerQueueDepth
	newConnectionInfo.transactionProcessorCount = _transactionProcessorCount
	newConnectionInfo.transactionProcessedCount = _transactionProcessedCount
	newConnectionInfo.transactionSuccessfulCount = _transactionSuccessfulCount
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGTransactionStatisticsImpl
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TransactionStatisticsImpl:{")
	buffer.WriteString(fmt.Sprintf("AverageProcessingTime: '%+v'", obj.averageProcessingTime))
	buffer.WriteString(fmt.Sprintf(", PendingTransactionsCount: '%d'", obj.pendingTransactionsCount))
	buffer.WriteString(fmt.Sprintf(", TransactionLoggerQueueDepth: '%d'", obj.transactionLoggerQueueDepth))
	buffer.WriteString(fmt.Sprintf(", TransactionProcessorCount: '%d'", obj.transactionProcessorCount))
	buffer.WriteString(fmt.Sprintf(", TransactionProcessedCount: '%d'", obj.transactionProcessedCount))
	buffer.WriteString(fmt.Sprintf(", TransactionSuccessfulCount: '%d'", obj.transactionSuccessfulCount))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGTransactionStatistics
/////////////////////////////////////////////////////////////////

// GetAverageProcessingTime returns the average processing time for the transactions
func (obj *TransactionStatisticsImpl) GetAverageProcessingTime() float64 {
	return obj.averageProcessingTime
}

// GetPendingTransactionsCount returns the pending transactions count
func (obj *TransactionStatisticsImpl) GetPendingTransactionsCount() int64 {
	return obj.pendingTransactionsCount
}

// GetTransactionLoggerQueueDepth returns the queue depth of transactionLogger
func (obj *TransactionStatisticsImpl) GetTransactionLoggerQueueDepth() int {
	return obj.transactionLoggerQueueDepth
}

// GetTransactionProcessorsCount returns the transaction processors count
func (obj *TransactionStatisticsImpl) GetTransactionProcessorsCount() int64 {
	return obj.transactionProcessorCount
}

// GetTransactionProcessedCount returns the processed transaction count
func (obj *TransactionStatisticsImpl) GetTransactionProcessedCount() int64 {
	return obj.transactionProcessedCount
}

// GetTransactionSuccessfulCount returns the successful transactions count
func (obj *TransactionStatisticsImpl) GetTransactionSuccessfulCount() int64 {
	return obj.transactionSuccessfulCount
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.averageProcessingTime, obj.pendingTransactionsCount, obj.transactionLoggerQueueDepth,
		obj.transactionProcessorCount, obj.transactionProcessedCount, obj.transactionSuccessfulCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionStatisticsImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TransactionStatisticsImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.averageProcessingTime, &obj.pendingTransactionsCount, &obj.transactionLoggerQueueDepth,
		&obj.transactionProcessorCount, &obj.transactionProcessedCount, &obj.transactionSuccessfulCount)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TransactionStatisticsImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



type UserInfoImpl struct {
	UserId   int
	UserType byte
	UserName string
}

// Make sure that the UserInfoImpl implements the TGUserInfo interface
var _ tgdb.TGUserInfo = (*UserInfoImpl)(nil)

func DefaultUserInfoImpl() *UserInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(UserInfoImpl{})

	return &UserInfoImpl{}
}

func NewUserInfoImpl(_userId int, _userName string, _userType byte) *UserInfoImpl {
	newConnectionInfo := DefaultUserInfoImpl()
	newConnectionInfo.UserId = _userId
	newConnectionInfo.UserType = _userType
	newConnectionInfo.UserName = _userName
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGUserInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("UserInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("UserId: '%d'", obj.UserId))
	buffer.WriteString(fmt.Sprintf(", UserType: '%+v'", obj.UserType))
	buffer.WriteString(fmt.Sprintf(", UserName: '%s'", obj.UserName))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGUserInfo
/////////////////////////////////////////////////////////////////

// GetName returns the user Name
func (obj *UserInfoImpl) GetName() string {
	return obj.UserName
}

// GetSystemId returns the system ID for this user
func (obj *UserInfoImpl) GetSystemId() int {
	return obj.UserId
}

// GetType returns the user type
func (obj *UserInfoImpl) GetType() byte {
	return obj.UserType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.UserId, obj.UserType, obj.UserName)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning UserInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.UserId, &obj.UserType, &obj.UserName)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning UserInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
