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
 * File Name: queryimpl.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: queryimpl.go 4271 2020-08-19 22:20:20Z nimish $
 */

package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"tgdb"
	"time"
)

type ResultSet struct {
	conn       tgdb.TGConnection
	currPos    int
	isOpen     bool
	resultId   int
	ResultList []interface{}
	MetaData	tgdb.TGResultSetMetaData
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
		ResultList: make([]interface{}, 0),
		MetaData: nil,
	}
	return &newResults
}

// Make sure that the ResultSet implements the TGResultSet interface
var _ tgdb.TGResultSet = (*ResultSet)(nil)

func NewResultSet(conn tgdb.TGConnection, resultId int) *ResultSet {
	newResults := DefaultResultSet()
	newResults.conn = conn
	newResults.resultId = resultId
	return newResults
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGResultSet
/////////////////////////////////////////////////////////////////

func (obj *ResultSet) GetConnection() tgdb.TGConnection {
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

func (obj *ResultSet) GetResults() []interface{} {
	return obj.ResultList
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGResultSet
/////////////////////////////////////////////////////////////////

// AddEntityToResultSet adds another entity to the result set
func (obj *ResultSet) AddEntityToResultSet(entity tgdb.TGEntity) tgdb.TGResultSet {
	obj.ResultList = append(obj.ResultList, entity)
	//obj.currPos++
	return obj
}

// Close closes the result set
func (obj *ResultSet) Close() tgdb.TGResultSet {
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
	return len(obj.ResultList)
}

// First returns the first entity in the result set
func (obj *ResultSet) First() interface{} {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.ResultList) == 0 {
		return nil
	}
	return obj.ResultList[0]
}

// Last returns the last Entity in the result set
func (obj *ResultSet) Last() interface{} {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.ResultList) == 0 {
		return nil
	}
	return obj.ResultList[len(obj.ResultList)-1]
}

// GetAt gets the entity at the position.
func (obj *ResultSet) GetAt(position int) interface{} {
	if obj.isOpen == false {
		return nil
	}
	if position >= 0 && position < len(obj.ResultList) {
		return obj.ResultList[position]
	}
	return nil
}

// GetExceptions gets the Exceptions in the result set
func (obj *ResultSet) GetExceptions() []tgdb.TGError {
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
	if len(obj.ResultList) == 0 {
		return false
	} else if obj.currPos < (len(obj.ResultList)-1) {
		return true
	}
	return false
}

// Next returns the next entity w.r.t to the current cursor position in the result set
func (obj *ResultSet) Next() interface{} {
	if obj.isOpen == false {
		return nil
	}
	if len(obj.ResultList) == 0 {
		return nil
	} else if obj.currPos < (len(obj.ResultList)-1) {
		obj.currPos++
		return obj.ResultList[obj.currPos]
	}
	return nil
}

// Skip skips a number of position
func (obj *ResultSet) Prev() interface{} {
	if obj.isOpen == false {
		return nil
	}
	if obj.currPos > 0 {
		obj.currPos--
		return obj.ResultList[obj.currPos]
	}
	return nil
}

// Skip skips a number of position
func (obj *ResultSet) Skip(position int) tgdb.TGResultSet {
	if obj.isOpen == false {
		return obj
	}
	newPos := obj.currPos + position
	if newPos >= 0 && newPos < len(obj.ResultList) {
		obj.currPos = newPos
	}
	return obj
}

// ToCollection converts the result set into a collection
func (obj *ResultSet) ToCollection() []interface{} {
	return obj.ResultList
}

func (obj *ResultSet) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ResultSet:{")
	buffer.WriteString(fmt.Sprintf("Connection: %+v", obj.conn))
	buffer.WriteString(fmt.Sprintf(", currPos: %+v", obj.currPos))
	buffer.WriteString(fmt.Sprintf(", isOpen: %+v", obj.isOpen))
	buffer.WriteString(fmt.Sprintf(", ResultId: %+v", obj.resultId))
	buffer.WriteString(fmt.Sprintf(", ResultList: %+v", obj.ResultList))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *ResultSet) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.conn, obj.currPos, obj.isOpen, obj.resultId, obj.ResultList)
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
	_, err := fmt.Fscanln(b, &obj.conn, &obj.currPos, &obj.isOpen, &obj.resultId, &obj.ResultList)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning ResultSet:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}

func (obj *ResultSet) GetMetadata() tgdb.TGResultSetMetaData {
	return obj.MetaData
}


// ======= Various Element Types for Gremlin Results =======
type ElementType int

const (
	ElementTypeInvalid ElementType = iota
	ElementTypeList
	ElementTypeAttr
	ElementTypeAttrValue
	ElementTypeAttrValueTransient
	ElementTypeEntity
	ElementTypeMap
)

func (elementType ElementType) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if elementType&ElementTypeInvalid == ElementTypeInvalid {
		buffer.WriteString("ElementTypeInvalid")
	} else if elementType&ElementTypeList == ElementTypeList {
		buffer.WriteString("ElementTypeList")
	} else if elementType&ElementTypeAttr == ElementTypeAttr {
		buffer.WriteString("ElementTypeAttr")
	} else if elementType&ElementTypeAttrValue == ElementTypeAttrValue {
		buffer.WriteString("ElementTypeAttrValue")
	} else if elementType&ElementTypeAttrValueTransient == ElementTypeAttrValueTransient {
		buffer.WriteString("ElementTypeAttrValueTransient")
	} else if elementType&ElementTypeEntity == ElementTypeEntity {
		buffer.WriteString("ElementTypeEntity")
	} else if elementType&ElementTypeMap == ElementTypeMap {
		buffer.WriteString("ElementTypeMap")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Helper functions for Gremlin Result
/////////////////////////////////////////////////////////////////

func FillCollection(entityStream tgdb.TGInputStream, gof tgdb.TGGraphObjectFactory, col []interface{}) ([]interface{}, tgdb.TGError) {
	//logger.Log(fmt.Sprint("Entering GremlinResult:FillCollection"))
	eleType, err := entityStream.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:FillCollection - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:FillCollection - unable to read element type in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted eleType: '%+v'", eleType))
	}

	if ElementType(eleType) == ElementTypeList {
		//return ConstructList(entityStream, gof, col)
		col, err := ConstructList(entityStream, gof, col)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:FillCollection - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "GremlinResult:FillCollection - unable to read element type in the response stream"
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		return col, nil
	} else {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:FillCollection - Invalid gremlin response collection type : %+v", ElementType(eleType)))
		errMsg := fmt.Sprintf("GremlinResult:FillCollection - Invalid gremlin response collection type : %+v", ElementType(eleType))
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, "")
	}
}

//func ConstructList(entityStream tgdb.TGInputStream, gof tgdb.TGGraphObjectFactory, col []interface{}) tgdb.TGError {
func ConstructList(entityStream tgdb.TGInputStream, gof tgdb.TGGraphObjectFactory, col []interface{}) ([]interface{}, tgdb.TGError) {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GremlinResult:ConstructList"))
	}
	size, err := entityStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to read size in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructList unable to read size in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted size: '%+v'", size))
	}

	eleType, err := entityStream.(*ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructList unable to read element type in the response stream"
		return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:ConstructList extracted eleType: '%+v'", eleType))
	}

	var dummyNode tgdb.TGNode
	if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
		node1, err := gof.CreateNode()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
			errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
			return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		dummyNode = node1
	}

	for i:=0; i<size; i++ {
		if ElementType(eleType) == ElementTypeEntity {
			entityType, err := entityStream.(*ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "GremlinResult:ConstructList - unable to read entity type in the entity stream"
				return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := tgdb.TGEntityKind(entityType)
			switch kindId {
			case tgdb.EntityKindNode:
				node, nErr := gof.CreateNode()
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = node.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to node.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructList - unable to node.ReadExternal in the entity stream"
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				col = append(col, node)
			case tgdb.EntityKindEdge:
				edge, nErr := gof.CreateEntity(tgdb.EntityKindEdge)
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = edge.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to edge.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructList - unable to edge.ReadExternal in the entity stream"
					return nil, GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				col = append(col, edge)
			case tgdb.EntityKindGraph:
				fallthrough
			case tgdb.EntityKindInvalid:
				fallthrough
			default:
				logger.Error(fmt.Sprint("ERROR: Returning GremlinResult:ConstructList - Invalid entity kind from gremlin response stream"))
				break
			}
		} else if ElementType(eleType) == ElementTypeList {
			colElem := make([]interface{}, 0)
			colElem, _ = ConstructList(entityStream, gof, colElem)
			col = append(col, colElem)
		} else if ElementType(eleType) == ElementTypeMap {
			mapElem := make(map[string]interface{}, 0)
			_ = ConstructMap(entityStream, gof, mapElem)
			col = append(col, mapElem)
		} else if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			attr, err := ReadExternalForEntity(dummyNode, entityStream)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read attr w/ Error: '%+v'", err.Error()))
				return nil, err
			}
			if ElementType(eleType) == ElementTypeAttr {
				col = append(col, attr)
			} else {
				col = append(col, attr.GetValue())
			}
		} else {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - Invalid element type '%+v' from Gremlin response stream", eleType))
		}
	}	// End of for loop

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning GremlinResult:ConstructList"))
	}
	return col, nil
}

func ConstructMap(entityStream tgdb.TGInputStream, gof tgdb.TGGraphObjectFactory, colMap map[string]interface{}) tgdb.TGError {
	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Entering GremlinResult:ConstructMap"))
	}
	size, err := entityStream.(*ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to read size in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructMap unable to read size in the response stream"
		return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted size: '%+v'", size))
	}

	for i:=0; i<size; i++ {
		key, err := entityStream.(*ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to read key in the response stream w/ error: '%s'", err.Error()))
			errMsg := "GremlinResult:ConstructMap - unable to read key in the entity stream"
			return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted key: '%+v'", key))
		}

		eleType, err := entityStream.(*ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "GremlinResult:ConstructMap unable to read element type in the response stream"
			return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
		}
		if logger.IsDebug() {
			logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:ConstructMap extracted eleType: '%+v'", eleType))
		}

		var dummyNode tgdb.TGNode
		if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			node1, err := gof.CreateNode()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
				errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
				return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			dummyNode = node1
		}

		if ElementType(eleType) == ElementTypeEntity {
			entityType, err := entityStream.(*ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "GremlinResult:ConstructMap - unable to read entity type in the entity stream"
				return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := tgdb.TGEntityKind(entityType)

			switch kindId {
			case tgdb.EntityKindNode:
				node, nErr := gof.CreateNode()
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructMap unable to create node for element type: '%+v'", ElementType(eleType))
					return GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = node.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to node.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructMap - unable to node.ReadExternal in the entity stream"
					return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				colMap[key] = node
			case tgdb.EntityKindEdge:
				edge, nErr := gof.CreateEntity(tgdb.EntityKindEdge)
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructMap unable to create node for element type: '%+v'", ElementType(eleType))
					return GetErrorByType(TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = edge.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to edge.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructMap - unable to edge.ReadExternal in the entity stream"
					return GetErrorByType(TGErrorGeneralException, "", errMsg, err.Error())
				}
				colMap[key] = edge
			case tgdb.EntityKindGraph:
				fallthrough
			case tgdb.EntityKindInvalid:
				fallthrough
			default:
				logger.Error(fmt.Sprint("ERROR: Returning GremlinResult:ConstructList - Invalid entity kind from gremlin response stream"))
				break
			}
		} else if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			attr, err := ReadExternalForEntity(dummyNode, entityStream)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read attr w/ Error: '%+v'", err.Error()))
				return err
			}
			colMap[key] = attr
			if ElementType(eleType) == ElementTypeAttr {
				colMap[key] = attr
			} else {
				colMap[key] = attr.GetValue()
			}
		} else if ElementType(eleType) == ElementTypeList {
			colElem := make([]interface{}, 0)
			colElem, _ = ConstructList(entityStream, gof, colElem)
			colMap[key] = colElem
		} else if ElementType(eleType) == ElementTypeMap {
			mapElem := make(map[string]interface{}, 0)
			_ = ConstructMap(entityStream, gof, mapElem)
			colMap[key] = mapElem
		}
	}	// End of for loop

	if logger.IsDebug() {
		logger.Debug(fmt.Sprint("Returning GremlinResult:ConstructMap"))
	}
	return nil
}


//var logger = logging.DefaultTGLogManager().GetLogger()

type TGQueryImpl struct {
	qryConnection tgdb.TGConnection
	qryHashId     int64
	qryOption     *TGQueryOptionImpl
	// These parameters are for tokenized query parameters specified as '?, ?, ...'
	qryParameters map[string]interface{}
}

func DefaultQuery() *TGQueryImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGQueryImpl{})

	newQuery := TGQueryImpl{
		qryHashId:     -1,
		qryParameters: make(map[string]interface{}, 0),
	}
	newQuery.qryOption = DefaultQueryOption()
	return &newQuery
}

// Make sure that the TGQueryImpl implements the TGQuery interface
var _ tgdb.TGQuery = (*TGQueryImpl)(nil)

func NewQuery(conn tgdb.TGConnection, queryHashId int64) *TGQueryImpl {
	newQuery := DefaultQuery()
	newQuery.qryConnection = conn
	newQuery.qryHashId = queryHashId
	return newQuery
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) GetQueryConnection() tgdb.TGConnection {
	return obj.qryConnection
}

func (obj *TGQueryImpl) GetQueryId() int64 {
	return obj.qryHashId
}

func (obj *TGQueryImpl) GetQueryOption() *TGQueryOptionImpl {
	return obj.qryOption
}

func (obj *TGQueryImpl) GetQueryParameters() map[string]interface{} {
	return obj.qryParameters
}

func (obj *TGQueryImpl) SetQueryId(qId int64) {
	obj.qryHashId = qId
}

func (obj *TGQueryImpl) SetQueryOption(queryOptions *TGQueryOptionImpl) {
	obj.qryOption = queryOptions
}

func (obj *TGQueryImpl) SetQueryParameters(params map[string]interface{}) {
	obj.qryParameters = params
}

/////////////////////////////////////////////////////////////////
// Private functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) setQueryParameter(name string, value interface{}) {
	obj.qryParameters[name] = value
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGQuery
/////////////////////////////////////////////////////////////////

// Close closes the Query
func (obj *TGQueryImpl) Close() {
	_, _ = obj.qryConnection.CloseQuery(obj.qryHashId)
}

// Execute executes the Query
func (obj *TGQueryImpl) Execute() tgdb.TGResultSet {
	rSet, _ := obj.qryConnection.ExecuteQueryWithId(obj.qryHashId, obj.qryOption)
	return rSet
}

// SetBoolean sets Boolean parameter
func (obj *TGQueryImpl) SetBoolean(name string, value bool) {
	obj.setQueryParameter(name, value)
}

// SetBytes sets Byte Parameter
func (obj *TGQueryImpl) SetBytes(name string, value []byte) {
	obj.setQueryParameter(name, value)
}

// SetChar sets Character Parameter
func (obj *TGQueryImpl) SetChar(name string, value string) {
	obj.setQueryParameter(name, value)
}

// SetDate sets Date Parameter
func (obj *TGQueryImpl) SetDate(name string, value time.Time) {
	obj.setQueryParameter(name, value)
}

// SetDouble sets Double Parameter
func (obj *TGQueryImpl) SetDouble(name string, value float64) {
	obj.setQueryParameter(name, value)
}

// SetFloat sets Float Parameter
func (obj *TGQueryImpl) SetFloat(name string, value float32) {
	obj.setQueryParameter(name, value)
}

// SetInt sets Integer Parameter
func (obj *TGQueryImpl) SetInt(name string, value int) {
	obj.setQueryParameter(name, value)
}

// SetLong sets Long Parameter
func (obj *TGQueryImpl) SetLong(name string, value int64) {
	obj.setQueryParameter(name, value)
}

// SetNull sets the parameter to null
func (obj *TGQueryImpl) SetNull(name string) {
	obj.setQueryParameter(name, nil)
}

// SetOption sets the Query Option
func (obj *TGQueryImpl) SetOption(options tgdb.TGQueryOption) {
	obj.qryOption = options.(*TGQueryOptionImpl)
}

// SetShort sets Short Parameter
func (obj *TGQueryImpl) SetShort(name string, value int16) {
	obj.setQueryParameter(name, value)
}

// SetString sets String Parameter
func (obj *TGQueryImpl) SetString(name string, value string) {
	obj.setQueryParameter(name, value)
}

func (obj *TGQueryImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TGQueryImpl:{")
	buffer.WriteString(fmt.Sprintf("QryConnection: %+v", obj.qryConnection))
	buffer.WriteString(fmt.Sprintf(", QryHashId: %+v", obj.qryHashId))
	buffer.WriteString(fmt.Sprintf(", QryOption: %+v", obj.qryOption.String()))
	buffer.WriteString(fmt.Sprintf(", QryParameters: %+v", obj.qryParameters))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.qryConnection, obj.qryHashId, obj.qryOption, obj.qryParameters)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.qryConnection, &obj.qryHashId, &obj.qryOption, &obj.qryParameters)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}



const (
	DefaultPrefetchSize    = 1000
	DefaultTraversalDepth  = 3
	DefaultEdgeLimit       = 0 // 0 ==> Unlimited
	DefaultOptionSortLimit = 0
	DefaultBatchSize       = 50

	OptionQueryBatchSize      = "batchsize"
	OptionQueryFetchSize      = "fetchsize"
	OptionQueryTraversalDepth = "traversaldepth"
	OptionQueryEdgeLimit      = "edgelimit"
	OptionQuerySortAttr       = "sortattrname"
	OptionQuerySortOrder      = "sortorder"			// 0 - asc, 1 - dsc
	OptionQuerySortLimit      = "sortresultlimit"
)

type TGQueryOptionImpl struct {
	optionProperties *SortedProperties
	mutable          bool
}

func DefaultQueryOption() *TGQueryOptionImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGQueryOptionImpl{})

	newQryOption := TGQueryOptionImpl{
		optionProperties: NewSortedProperties(),
		mutable:          false,
	}
	return &newQryOption
}

// Make sure that the TGQueryOptionImpl implements the TGQueryOption interface
var _ tgdb.TGQueryOption = (*TGQueryOptionImpl)(nil)

//func NewQueryOption(mutable bool) *TGQueryOptionImpl {
func NewQueryOption() *TGQueryOptionImpl {
	newQryOption := DefaultQueryOption()
	newQryOption.mutable = true
	newQryOption.preloadQueryOptions()
	return newQryOption
}

/////////////////////////////////////////////////////////////////
// Private functions from Interface ==> TGQueryOption
/////////////////////////////////////////////////////////////////

func (obj *TGQueryOptionImpl) preloadQueryOptions() {
	// Add property will either insert or update it using underlying TGProperties functionality
	//obj.AddProperty(OptionQueryBatchSize, -1)
	obj.AddProperty(OptionQueryFetchSize, strconv.Itoa(DefaultPrefetchSize))
	obj.AddProperty(OptionQueryTraversalDepth, strconv.Itoa(DefaultTraversalDepth))
	obj.AddProperty(OptionQueryEdgeLimit, strconv.Itoa(DefaultEdgeLimit))
	obj.AddProperty(OptionQuerySortAttr, "")
	obj.AddProperty(OptionQuerySortOrder, "0")	// Ascending Order
	obj.AddProperty(OptionQuerySortLimit, strconv.Itoa(DefaultOptionSortLimit))
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGQueryOption
/////////////////////////////////////////////////////////////////

func (obj *TGQueryOptionImpl) GetQueryOptionProperties() *SortedProperties {
	return obj.optionProperties
}

func (obj *TGQueryOptionImpl) GetIsMutable() bool {
	return obj.mutable
}

func (obj *TGQueryOptionImpl) SetUserAndPassword(user, pwd string) tgdb.TGError {
	return SetUserAndPassword(obj.optionProperties, user, pwd)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGQueryOption
/////////////////////////////////////////////////////////////////

// GetBatchSize gets the current value of the batch size
func (obj TGQueryOptionImpl) GetBatchSize() int {
	if DoesPropertyExist(obj.optionProperties, OptionQueryBatchSize) {
		cn := NewConfigName(OptionQueryBatchSize, OptionQueryBatchSize, strconv.Itoa(DefaultBatchSize))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultBatchSize)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultPrefetchSize
}

// SetBatchSize sets a limit on the batch. Default is 50
func (obj TGQueryOptionImpl) SetBatchSize(size int) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if size == 0 {
		size = -1
	}
	obj.SetProperty(OptionQueryBatchSize, strconv.Itoa(size))
	return nil
}

// GetPreFetchSize gets the current value of the pre-fetch size
func (obj *TGQueryOptionImpl) GetPreFetchSize() int {
	if DoesPropertyExist(obj.optionProperties, OptionQueryFetchSize) {
		cn := NewConfigName(OptionQueryFetchSize, OptionQueryFetchSize, strconv.Itoa(DefaultPrefetchSize))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultPrefetchSize)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultPrefetchSize
}

// SetPreFetchSize sets a limit on the number of entities(nodes and edges) return in a query. Default is 1000
func (obj *TGQueryOptionImpl) SetPreFetchSize(size int) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if size == 0 {
		size = -1
	}
	obj.SetProperty(OptionQueryFetchSize, strconv.Itoa(size))
	return nil
}

// GetTraversalDepth gets the current value of traversal depth
func (obj *TGQueryOptionImpl) GetTraversalDepth() int {
	if DoesPropertyExist(obj.optionProperties, OptionQueryTraversalDepth) {
		cn := NewConfigName(OptionQueryTraversalDepth, OptionQueryTraversalDepth, strconv.Itoa(DefaultTraversalDepth))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultTraversalDepth)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultTraversalDepth
}

// SetTraversalDepth sets the additional level of traversal from the query result set. Default is 3.
func (obj *TGQueryOptionImpl) SetTraversalDepth(depth int) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if depth == 0 {
		depth = -1
	}
	obj.SetProperty(OptionQueryTraversalDepth, strconv.Itoa(depth))
	return nil
}

// GetEdgeLimit gets the current value of edge limit
func (obj *TGQueryOptionImpl) GetEdgeLimit() int {
	if DoesPropertyExist(obj.optionProperties, OptionQueryEdgeLimit) {
		cn := NewConfigName(OptionQueryEdgeLimit, OptionQueryEdgeLimit, strconv.Itoa(DefaultEdgeLimit))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultEdgeLimit)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultEdgeLimit
}

// SetEdgeLimit sets the number of edges per node to be returned in a query.  Default is 0 which means unlimited.
func (obj *TGQueryOptionImpl) SetEdgeLimit(limit int) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if limit == 0 {
		limit = -1
	}
	obj.SetProperty(OptionQueryEdgeLimit, strconv.Itoa(limit))
	return nil
}

// GetSortAttrName gets sort attribute Name
func (obj *TGQueryOptionImpl) GetSortAttrName() string {
	if DoesPropertyExist(obj.optionProperties, OptionQuerySortAttr) {
		cn := NewConfigName(OptionQuerySortAttr, OptionQuerySortAttr, "")
		if s := obj.GetProperty(cn, ""); s != "" {
			return s
		}
	}
	return ""
}

// SetSortAttrName sets sort attribute Name
func (obj *TGQueryOptionImpl) SetSortAttrName(name string) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if name == "" || len(name) == 0 {
		errMsg := fmt.Sprintf("AttributeTypeInvalid attribute Name")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.SetProperty(OptionQuerySortAttr, name)
	return nil
}

// IsSortOrderDsc gets sort order desc
func (obj *TGQueryOptionImpl) IsSortOrderDsc() bool {
	if DoesPropertyExist(obj.optionProperties, OptionQuerySortOrder) {
		cn := NewConfigName(OptionQuerySortOrder, OptionQuerySortOrder, "false")
		if s := obj.GetProperty(cn, "false"); s != "" {
			v, _ := strconv.ParseBool(s)
			return v
		}
	}
	return false
}

// SetSortOrderDsc sets sort order desc
func (obj *TGQueryOptionImpl) SetSortOrderDsc(isDsc bool) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if isDsc {
		obj.SetProperty(OptionQuerySortOrder, "1")
	} else {
		obj.SetProperty(OptionQuerySortOrder, "0")
	}
	return nil
}

// GetSortResultLimit gets sort result limit
func (obj *TGQueryOptionImpl) GetSortResultLimit() int {
	if DoesPropertyExist(obj.optionProperties, OptionQuerySortLimit) {
		cn := NewConfigName(OptionQuerySortLimit, OptionQuerySortLimit, strconv.Itoa(DefaultOptionSortLimit))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultOptionSortLimit)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultOptionSortLimit
}

// SetSortResultLimit sets sort result limit
func (obj *TGQueryOptionImpl) SetSortResultLimit(limit int) tgdb.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if limit <= 0 {
		errMsg := fmt.Sprintf("Invalid sort limit")
		return GetErrorByType(TGErrorGeneralException, INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.SetProperty(OptionQuerySortLimit, strconv.Itoa(limit))
	return nil
}

func (obj *TGQueryOptionImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("TGQueryOptionImpl:{")
	buffer.WriteString(fmt.Sprintf("optionProperties: %+v", obj.optionProperties))
	buffer.WriteString(fmt.Sprintf(", mutable: %+v", obj.mutable))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGProperties
/////////////////////////////////////////////////////////////////

// AddProperty checks whether a property already exists, else adds a new property in the form of Name=value pair
func (obj *TGQueryOptionImpl) AddProperty(name, value string) {
	//logger.Log(fmt.Sprintf("TGQueryOptionImpl::AddProperty received N-V Pair as '%+v':'%+v'", Name, value))
	obj.optionProperties.AddProperty(name, value)
}

// GetProperty gets the property either with value or default value
func (obj *TGQueryOptionImpl) GetProperty(cn tgdb.TGConfigName, value string) string {
	return obj.optionProperties.GetProperty(cn, value)
}

// SetProperty sets existing property value in the form of Name=value pair
func (obj *TGQueryOptionImpl) SetProperty(name, value string) {
	obj.optionProperties.SetProperty(name, value)
}

// GetPropertyAsInt gets Property as int value
func (obj *TGQueryOptionImpl) GetPropertyAsInt(cn tgdb.TGConfigName) int {
	return obj.optionProperties.GetPropertyAsInt(cn)
}

// GetPropertyAsLong gets Property as long value
func (obj *TGQueryOptionImpl) GetPropertyAsLong(cn tgdb.TGConfigName) int64 {
	return obj.optionProperties.GetPropertyAsLong(cn)
}

// GetPropertyAsBoolean gets Property as bool value
func (obj *TGQueryOptionImpl) GetPropertyAsBoolean(cn tgdb.TGConfigName) bool {
	return obj.optionProperties.GetPropertyAsBoolean(cn)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryOptionImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.optionProperties, obj.mutable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryOptionImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *TGQueryOptionImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.optionProperties, &obj.mutable)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning TGQueryOptionImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}


type ResultDataDescriptor struct {
	DataType int
	Annot string
	IsItMap bool
	IsItArray bool
	HasType bool
	ContainedSize int
	ScalarType int
	SysObject tgdb.TGSystemObject
	KeyDesc tgdb.TGResultDataDescriptor
	ValueDesc tgdb.TGResultDataDescriptor
	ContainedDesc []tgdb.TGResultDataDescriptor
}

func (obj* ResultDataDescriptor) GetDataType() int {
	return obj.DataType
}

func (obj* ResultDataDescriptor) GetContainedDataSize() int {
	return obj.ContainedSize
}

func (obj* ResultDataDescriptor) IsMap() bool {
	return obj.IsItMap
}

func (obj* ResultDataDescriptor) IsArray() bool {
	return obj.IsItArray
}

func (obj* ResultDataDescriptor) HasConcreteType() bool {
	return obj.HasType
}

func (obj* ResultDataDescriptor) GetScalarType() int {
	return obj.ScalarType
}

func (obj* ResultDataDescriptor) GetSystemObject() tgdb.TGSystemObject {
	return obj.SysObject
}

func (obj* ResultDataDescriptor) GetKeyDescriptor() tgdb.TGResultDataDescriptor {
	return obj.KeyDesc
}

func (obj* ResultDataDescriptor) GetValueDescriptor() tgdb.TGResultDataDescriptor {
	return obj.ValueDesc
}

func (obj* ResultDataDescriptor) GetContainedDescriptors() []tgdb.TGResultDataDescriptor {
	return obj.ContainedDesc
}

func (obj* ResultDataDescriptor) GetContainedDescriptor(position int) tgdb.TGResultDataDescriptor {
	return obj.ContainedDesc[position]
}


////////////// setters for ResultDataDescriptor
func (obj* ResultDataDescriptor) SetDataType(dataType int) {
	obj.DataType = dataType
}

func (obj* ResultDataDescriptor) SetContainedDataSize(dataSize int) {
	obj.ContainedSize = dataSize
}

func (obj* ResultDataDescriptor) SetIsMap(isMap bool) {
	obj.IsItMap = isMap
}

func (obj* ResultDataDescriptor) SetIsArray(isArray bool) {
	obj.IsItArray = isArray
}

func (obj* ResultDataDescriptor) SetHasConcreteType(hasType bool) {
	obj.HasType = hasType
}

func (obj* ResultDataDescriptor) SetScalarType(scalarType int) {
	obj.ScalarType = scalarType
}

func (obj* ResultDataDescriptor) SetSystemObject(sysObject tgdb.TGSystemObject) {
	obj.SysObject = sysObject
}

func (obj* ResultDataDescriptor) SetKeyDescriptor(keyDesc tgdb.TGResultDataDescriptor) {
	obj.KeyDesc = keyDesc
}

func (obj* ResultDataDescriptor) SetValueDescriptor(valueDesc tgdb.TGResultDataDescriptor) {
	obj.ValueDesc = valueDesc
}

func (obj* ResultDataDescriptor) SetContainedDescriptors(containedDesc []tgdb.TGResultDataDescriptor) {
	obj.ContainedDesc = containedDesc
}

////////////// setters for ResultDataDescriptor

func DefaultResultDataDescriptor () *ResultDataDescriptor {
	newResults := ResultDataDescriptor{
		DataType:      tgdb.TYPE_NODE,
		Annot:         "",
		IsItMap:         false,
		IsItArray:       false,
		HasType:       false,
		ContainedSize: 0,
		ScalarType:    AttributeTypeInvalid,
		SysObject:     nil,
		KeyDesc:       nil,
		ValueDesc:     nil,
		ContainedDesc: nil,
	}
	return &newResults

}

func NewResultDataDescriptor (typeOfDataDescriptor int) *ResultDataDescriptor {
	newResults := DefaultResultDataDescriptor()
	newResults.DataType = typeOfDataDescriptor
	return newResults
}

func NewResultDataDescriptorWithScalarAnnot (typeOfDataDescriptor int, scalarAnnot string) *ResultDataDescriptor {
	newResults := NewResultDataDescriptor(typeOfDataDescriptor)
	newResults.HasType = true

	switch (scalarAnnot) {
	case "?":
		newResults.ScalarType = AttributeTypeBoolean;
		break;
	case "b":
		newResults.ScalarType = AttributeTypeByte;
		break;
	case "c":
		newResults.ScalarType = AttributeTypeChar;
		break;
	case "h":
		newResults.ScalarType = AttributeTypeShort;
		break;
	case "i":
		newResults.ScalarType = AttributeTypeInteger;
		break;
	case "l":
		newResults.ScalarType = AttributeTypeLong;
		break;
	case "f":
		newResults.ScalarType = AttributeTypeFloat;
		break;
	case "d":
		newResults.ScalarType = AttributeTypeDouble;
		break;
	case "s":
		newResults.ScalarType = AttributeTypeString;
		break;
	case "D":
		newResults.ScalarType = AttributeTypeDate;
		break;
	case "e":
		newResults.ScalarType = AttributeTypeTime;
		break;
	case "a":
		newResults.ScalarType = AttributeTypeTimeStamp;
		break;
	case "n":
		newResults.ScalarType = AttributeTypeNumber;
		break;
	default:
		newResults.HasType = false;
		break;
	}
	return newResults
}

type ResultSetMetadata struct {
	ResultDataDescriptor* tgdb.TGResultDataDescriptor
	ResultType int
	Annot string
}

/////////// Constructor Initialization starts
func DefaultResultSetMetadata () *ResultSetMetadata {
	newResults := ResultSetMetadata{
		ResultDataDescriptor: nil,
		ResultType:           tgdb.TYPE_UNKNOWN,
		Annot:                "",
	}
	return &newResults
}

func NewResultSetMetadataWithAnnot (annot string) *ResultSetMetadata {
	newResults := DefaultResultSetMetadata()
	newResults.SetAnnot(annot)
	return newResults
}

/////////// Constructor Initialization end

func (obj *ResultSetMetadata) SetAnnot (annot string) {
	obj.Annot = annot
}

func (obj *ResultSetMetadata) GetAnnot () string {
	return obj.Annot
}

func (obj *ResultSetMetadata) GetResultDataDescriptor () *tgdb.TGResultDataDescriptor {
	return obj.ResultDataDescriptor
}

func (obj *ResultSetMetadata) GetResultType() int {
	return obj.ResultType
}

func (obj *ResultSet) SetResultTypeAnnotation(annot string) tgdb.TGError {
	if  len(annot) == 0 {
		return nil
	}

	obj.MetaData = NewResultSetMetadataWithAnnot(annot)
	gmd, err := obj.conn.GetGraphMetadata(false)


	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Failed to initialize result set metadata"))
		return err
	}

	obj.MetaData.(*ResultSetMetadata).Initialize(gmd)
	return nil
}

func (obj *ResultSetMetadata) Initialize(gmd tgdb.TGGraphMetadata)  {
	/*ddesc*/obj.ResultDataDescriptor = obj.ConstructDataDescriptor(gmd, obj.GetAnnot())
	obj.ResultType = (*obj.ResultDataDescriptor).GetDataType()
	//obj.resultType = (*ddesc).GetDataType()
}

func (obj *ResultSetMetadata) ConstructDataDescriptor(gmd tgdb.TGGraphMetadata, annot string) *tgdb.TGResultDataDescriptor {
	if annot == "" {
		return nil
	}

	var desc tgdb.TGResultDataDescriptor;
	desc = nil

	startsWith := string(annot[0])
	switch startsWith {
	case "V":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_NODE)
			break
		}
	case "E":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_EDGE)
			break;
		}
	case "P":
		{
			desc = *obj.ConstructPathDataDescriptor(gmd, annot);
			break
		}
	case "S":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_SCALAR);
			break
		}
	case "[":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_LIST)
			nDesc := obj.ConstructDataDescriptor(gmd, annot[1:])
			aDesc := make([]tgdb.TGResultDataDescriptor, 1)
			aDesc[0] = *nDesc
			desc.(*ResultDataDescriptor).SetContainedDescriptors(aDesc)
			desc.(*ResultDataDescriptor).SetIsArray(true);
			break;
		}
	case "{":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_MAP)
			keyDesc := NewResultDataDescriptorWithScalarAnnot(tgdb.TYPE_SCALAR, "s")
			valueDesc := NewResultDataDescriptor(tgdb.TYPE_SCALAR)
			desc.(*ResultDataDescriptor).SetKeyDescriptor(keyDesc)
			desc.(*ResultDataDescriptor).SetValueDescriptor(valueDesc)
			desc.(*ResultDataDescriptor).SetIsMap(true)
			break;
		}
	case "(":
		{
			desc = NewResultDataDescriptor(tgdb.TYPE_TUPLE);
			break;
		}
	case "?":
	case "b":
	case "c":
	case "h":
	case "i":
	case "l":
	case "f":
	case "d":
	case "s":
	case "D":
	case "e":
	case "a":
	case "n":
		{
			desc = NewResultDataDescriptorWithScalarAnnot(tgdb.TYPE_SCALAR, string(annot[0]))
			break
		}
	default:
		//gLogger.log(TGLogger.TGLevel.Info, "Failed to initialize data descriptor with type annotation : %c", chars[0]);
		break;
	}
	return &desc
}

func (obj *ResultSetMetadata) ConstructPathDataDescriptor (gmd tgdb.TGGraphMetadata, annot string) *tgdb.TGResultDataDescriptor {
	var desc tgdb.TGResultDataDescriptor;
	desc = NewResultDataDescriptor(tgdb.TYPE_PATH)
	data := strings.Split(annot[2:len(annot)-1], ",")
	edesc := make([]tgdb.TGResultDataDescriptor, len(data))
	for i := 0; i < len(data); i++ {
		edesc[i] = *obj.ConstructDataDescriptor(gmd, data[i])
	}
	desc.(*ResultDataDescriptor).SetContainedDescriptors(edesc)
	return &desc
}
