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
 * File name: GremlinResult.go
 * Created on: Apr 27, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

package query

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

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

func FillCollection(entityStream types.TGInputStream, gof types.TGGraphObjectFactory, col []interface{}) types.TGError {
	//logger.Log(fmt.Sprint("Entering GremlinResult:FillCollection"))
	eleType, err := entityStream.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:FillCollection - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:FillCollection - unable to read element type in the response stream"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted eleType: '%+v'", eleType))

	if ElementType(eleType) == ElementTypeList {
		return ConstructList(entityStream, gof, col)
	} else {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:FillCollection - Invalid gremlin response collection type : %+v", ElementType(eleType)))
		errMsg := fmt.Sprintf("GremlinResult:FillCollection - Invalid gremlin response collection type : %+v", ElementType(eleType))
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, "")
	}
	return nil
}

func ConstructList(entityStream types.TGInputStream, gof types.TGGraphObjectFactory, col []interface{}) types.TGError {
	logger.Log(fmt.Sprint("Entering GremlinResult:ConstructList"))
	size, err := entityStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to read size in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructList unable to read size in the response stream"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted size: '%+v'", size))

	eleType, err := entityStream.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructList unable to read element type in the response stream"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:ConstructList extracted eleType: '%+v'", eleType))

	var dummyNode types.TGNode
	if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
		node1, err := gof.CreateNode()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
			errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
			return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		dummyNode = node1
	}

	for i:=0; i<size; i++ {
		if ElementType(eleType) == ElementTypeEntity {
			entityType, err := entityStream.(*iostream.ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "GremlinResult:ConstructList - unable to read entity type in the entity stream"
				return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := types.TGEntityKind(entityType)
			switch kindId {
			case types.EntityKindNode:
				node, nErr := gof.CreateNode()
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = node.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to node.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructList - unable to node.ReadExternal in the entity stream"
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				col = append(col, node)
			case types.EntityKindEdge:
				edge, nErr := gof.CreateEntity(types.EntityKindEdge)
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = edge.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList  - unable to edge.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructList - unable to edge.ReadExternal in the entity stream"
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				col = append(col, edge)
			case types.EntityKindGraph:
				fallthrough
			case types.EntityKindInvalid:
				fallthrough
			default:
				logger.Error(fmt.Sprint("ERROR: Returning GremlinResult:ConstructList - Invalid entity kind from gremlin response stream"))
				break
			}
		} else if ElementType(eleType) == ElementTypeList {
			colElem := make([]interface{}, 0)
			_ = ConstructList(entityStream, gof, colElem)
			col = append(col, colElem)
		} else if ElementType(eleType) == ElementTypeMap {
			mapElem := make(map[string]interface{}, 0)
			_ = ConstructMap(entityStream, gof, mapElem)
			col = append(col, mapElem)
		} else if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			attr, err := model.ReadExternalForEntity(dummyNode, entityStream)
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning AbstractEntity:AbstractEntityReadExternal - unable to read attr w/ Error: '%+v'", err.Error()))
				return err
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

	logger.Log(fmt.Sprint("Returning GremlinResult:ConstructList"))
	return nil
}

func ConstructMap(entityStream types.TGInputStream, gof types.TGGraphObjectFactory, colMap map[string]interface{}) types.TGError {
	logger.Log(fmt.Sprint("Entering GremlinResult:ConstructMap"))
	size, err := entityStream.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to read size in the response stream w/ error: '%s'", err.Error()))
		errMsg := "GremlinResult:ConstructMap unable to read size in the response stream"
		return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
	}
	logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted size: '%+v'", size))

	for i:=0; i<size; i++ {
		key, err := entityStream.(*iostream.ProtocolDataInputStream).ReadUTF()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to read key in the response stream w/ error: '%s'", err.Error()))
			errMsg := "GremlinResult:ConstructMap - unable to read key in the entity stream"
			return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:FillCollection extracted key: '%+v'", key))

		eleType, err := entityStream.(*iostream.ProtocolDataInputStream).ReadByte()
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to read eleType in the response stream w/ error: '%s'", err.Error()))
			errMsg := "GremlinResult:ConstructMap unable to read element type in the response stream"
			return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
		}
		logger.Debug(fmt.Sprintf("Inside Returning GremlinResult:ConstructMap extracted eleType: '%+v'", eleType))

		var dummyNode types.TGNode
		if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			node1, err := gof.CreateNode()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructList - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
				errMsg := fmt.Sprintf("GremlinResult:ConstructList unable to create node for element type: '%+v'", ElementType(eleType))
				return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			dummyNode = node1
		}

		if ElementType(eleType) == ElementTypeEntity {
			entityType, err := entityStream.(*iostream.ProtocolDataInputStream).ReadByte()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to read entityType in the response stream w/ error: '%s'", err.Error()))
				errMsg := "GremlinResult:ConstructMap - unable to read entity type in the entity stream"
				return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
			}
			kindId := types.TGEntityKind(entityType)

			switch kindId {
			case types.EntityKindNode:
				node, nErr := gof.CreateNode()
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructMap unable to create node for element type: '%+v'", ElementType(eleType))
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = node.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to node.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructMap - unable to node.ReadExternal in the entity stream"
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				colMap[key] = node
			case types.EntityKindEdge:
				edge, nErr := gof.CreateEntity(types.EntityKindEdge)
				if nErr != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap - unable to gof.CreateNode() for element type: '%+v'", ElementType(eleType)))
					errMsg := fmt.Sprintf("GremlinResult:ConstructMap unable to create node for element type: '%+v'", ElementType(eleType))
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, nErr.Error())
				}
				err = edge.ReadExternal(entityStream)
				if err != nil {
					logger.Error(fmt.Sprintf("ERROR: Returning GremlinResult:ConstructMap  - unable to edge.ReadExternal w/ error: '%s'", err.Error()))
					errMsg := "GremlinResult:ConstructMap - unable to edge.ReadExternal in the entity stream"
					return exception.GetErrorByType(types.TGErrorGeneralException, "", errMsg, err.Error())
				}
				colMap[key] = edge
			case types.EntityKindGraph:
				fallthrough
			case types.EntityKindInvalid:
				fallthrough
			default:
				logger.Error(fmt.Sprint("ERROR: Returning GremlinResult:ConstructList - Invalid entity kind from gremlin response stream"))
				break
			}
		} else if ElementType(eleType) == ElementTypeAttr || ElementType(eleType) == ElementTypeAttrValue || ElementType(eleType) == ElementTypeAttrValueTransient {
			attr, err := model.ReadExternalForEntity(dummyNode, entityStream)
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
			_ = ConstructList(entityStream, gof, colElem)
			colMap[key] = colElem
		} else if ElementType(eleType) == ElementTypeMap {
			mapElem := make(map[string]interface{}, 0)
			_ = ConstructMap(entityStream, gof, mapElem)
			colMap[key] = mapElem
		}
	}	// End of for loop

	logger.Log(fmt.Sprint("Returning GremlinResult:ConstructMap"))
	return nil
}

