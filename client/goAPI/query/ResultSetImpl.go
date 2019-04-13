package query

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGResultSet.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ResultSet struct {
	conn       types.TGConnection
	currPos    int
	isOpen     bool
	resultId   int
	resultList []types.TGEntity
}

func DefaultResultSet() *ResultSet {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ResultSet{})

	newResults := ResultSet{
		currPos:    -1,
		isOpen:     true,
		resultId:   -1,
		resultList: make([]types.TGEntity, 0),
	}
	return &newResults
}

// Make sure that the ResultSet implements the TGResultSet interface
var _ types.TGResultSet = (*ResultSet)(nil)

func NewResultSet(conn types.TGConnection, resultId int) *ResultSet {
	newResults := DefaultResultSet()
	newResults.conn = conn
	newResults.resultId = resultId
	return newResults
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGResultSet
/////////////////////////////////////////////////////////////////

func (obj *ResultSet) GetConnection() types.TGConnection {
	return obj.conn
}

func (obj *ResultSet) GetCurrentPos() int {
	return obj.currPos
}

func (obj *ResultSet) GetIsOpen() bool {
	return obj.isOpen
}

func (obj *ResultSet) GetResultId() int {
	return obj.resultId
}

func (obj *ResultSet) GetResults() []types.TGEntity {
	return obj.resultList
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGResultSet
/////////////////////////////////////////////////////////////////

// AddEntityToResultSet adds another entity to the result set
func (obj *ResultSet) AddEntityToResultSet(entity types.TGEntity) types.TGResultSet {
	obj.resultList = append(obj.resultList, entity)
	//obj.currPos++
	return obj
}

// Close closes the result set
func (obj *ResultSet) Close() types.TGResultSet {
	obj.isOpen = false
	return obj
}

// Count returns nos of entities returned by the query. The result set has a cursor which prefetches
// "n" rows as per the query constraint. If the nos of entities returned by the query is less
// than prefetch count, then all are returned.
func (obj *ResultSet) Count() int {
	if obj.isOpen == false {
		return 0
	}
	return len(obj.resultList)
}

// First returns the first entity in the result set
func (obj *ResultSet) First() types.TGEntity {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.resultList) == 0 {
		return nil
	}
	return obj.resultList[0]
}

// Last returns the last Entity in the result set
func (obj *ResultSet) Last() types.TGEntity {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.resultList) == 0 {
		return nil
	}
	return obj.resultList[len(obj.resultList)-1]
}

// GetAt gets the entity at the position.
func (obj *ResultSet) GetAt(position int) types.TGEntity {
	if obj.isOpen == false {
		return nil
	}
	if position >= 0 && position < len(obj.resultList) {
		return obj.resultList[position]
	}
	return nil
}

// GetExceptions gets the Exceptions in the result set
func (obj *ResultSet) GetExceptions() []types.TGError {
	// TODO: Revisit later - No-Op for Now
	return nil
}

// GetPosition gets the Current cursor position. A result set upon creation is set to the position 0.
func (obj *ResultSet) GetPosition() int {
	if obj.isOpen == false {
		return 0
	}
	return obj.currPos
}

// HasExceptions checks whether the result set has any exceptions
func (obj *ResultSet) HasExceptions() bool {
	// TODO: Revisit later - No-Op for Now
	return false
}

// HasNext Check whether there is next entry in result set
func (obj *ResultSet) HasNext() bool {
	if obj.isOpen == false {
		return false
	}
	if len(obj.resultList) == 0 {
		return false
	} else if obj.currPos < (len(obj.resultList)-1) {
		return true
	}
	return false
}

// Next returns the next entity w.r.t to the current cursor position in the result set
func (obj *ResultSet) Next() types.TGEntity {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.resultList) == 0 {
		return nil
	} else if obj.currPos < (len(obj.resultList)-1) {
		obj.currPos++
		return obj.resultList[obj.currPos]
	}
	return nil
}

// Skip skips a number of position
func (obj *ResultSet) Prev() types.TGEntity {
	if obj.isOpen == false {
		return nil
	}
	if obj.currPos > 0 {
		obj.currPos--
		return obj.resultList[obj.currPos]
	}
	return nil
}

// Skip skips a number of position
func (obj *ResultSet) Skip(position int) types.TGResultSet {
	if obj.isOpen == false {
		return obj
	}
	newPos := obj.currPos + position
	if newPos >= 0 && newPos < len(obj.resultList) {
		obj.currPos = newPos
	}
	return obj
}

func (obj *ResultSet) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ResultSet:{")
	buffer.WriteString(fmt.Sprintf("Connection: %+v", obj.conn))
	buffer.WriteString(fmt.Sprintf(", currPos: %+v", obj.currPos))
	buffer.WriteString(fmt.Sprintf(", isOpen: %+v", obj.isOpen))
	buffer.WriteString(fmt.Sprintf(", ResultId: %+v", obj.resultId))
	buffer.WriteString(fmt.Sprintf(", ResultList: %+v", obj.resultList))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ResultSet) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.conn, obj.currPos, obj.isOpen, obj.resultId, obj.resultList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ResultSet:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *ResultSet) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.conn, &obj.currPos, &obj.isOpen, &obj.resultId, &obj.resultList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ResultSet:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
