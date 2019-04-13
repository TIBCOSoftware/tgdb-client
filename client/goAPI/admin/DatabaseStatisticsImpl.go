package admin

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
 * File name: DatabaseStatisticsImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
var _ TGDatabaseStatistics = (*DatabaseStatisticsImpl)(nil)

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
