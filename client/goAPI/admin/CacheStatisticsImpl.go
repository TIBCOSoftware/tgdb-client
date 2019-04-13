package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
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
 * File name: CacheStatisticsImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

var logger = logging.DefaultTGLogManager().GetLogger()

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
var _ TGCacheStatistics = (*CacheStatisticsImpl)(nil)

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
