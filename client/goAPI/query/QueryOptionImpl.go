package query

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"strconv"
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
 * File name: TGQueryOption.go
 * Created on: Oct 13, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

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
	optionProperties *utils.SortedProperties
	mutable          bool
}

func DefaultQueryOption() *TGQueryOptionImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(TGQueryOptionImpl{})

	newQryOption := TGQueryOptionImpl{
		optionProperties: utils.NewSortedProperties(),
		mutable:          false,
	}
	return &newQryOption
}

// Make sure that the TGQueryOptionImpl implements the TGQueryOption interface
var _ types.TGQueryOption = (*TGQueryOptionImpl)(nil)

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

func (obj *TGQueryOptionImpl) GetQueryOptionProperties() *utils.SortedProperties {
	return obj.optionProperties
}

func (obj *TGQueryOptionImpl) GetIsMutable() bool {
	return obj.mutable
}

func (obj *TGQueryOptionImpl) SetUserAndPassword(user, pwd string) types.TGError {
	return utils.SetUserAndPassword(obj.optionProperties, user, pwd)
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGQueryOption
/////////////////////////////////////////////////////////////////

// GetBatchSize gets the current value of the batch size
func (obj TGQueryOptionImpl) GetBatchSize() int {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQueryBatchSize) {
		cn := utils.NewConfigName(OptionQueryBatchSize, OptionQueryBatchSize, strconv.Itoa(DefaultBatchSize))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultBatchSize)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultPrefetchSize
}

// SetBatchSize sets a limit on the batch. Default is 50
func (obj TGQueryOptionImpl) SetBatchSize(size int) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if size == 0 {
		size = -1
	}
	obj.SetProperty(OptionQueryBatchSize, strconv.Itoa(size))
	return nil
}

// GetPreFetchSize gets the current value of the pre-fetch size
func (obj *TGQueryOptionImpl) GetPreFetchSize() int {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQueryFetchSize) {
		cn := utils.NewConfigName(OptionQueryFetchSize, OptionQueryFetchSize, strconv.Itoa(DefaultPrefetchSize))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultPrefetchSize)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultPrefetchSize
}

// SetPreFetchSize sets a limit on the number of entities(nodes and edges) return in a query. Default is 1000
func (obj *TGQueryOptionImpl) SetPreFetchSize(size int) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if size == 0 {
		size = -1
	}
	obj.SetProperty(OptionQueryFetchSize, strconv.Itoa(size))
	return nil
}

// GetTraversalDepth gets the current value of traversal depth
func (obj *TGQueryOptionImpl) GetTraversalDepth() int {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQueryTraversalDepth) {
		cn := utils.NewConfigName(OptionQueryTraversalDepth, OptionQueryTraversalDepth, strconv.Itoa(DefaultTraversalDepth))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultTraversalDepth)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultTraversalDepth
}

// SetTraversalDepth sets the additional level of traversal from the query result set. Default is 3.
func (obj *TGQueryOptionImpl) SetTraversalDepth(depth int) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if depth == 0 {
		depth = -1
	}
	obj.SetProperty(OptionQueryTraversalDepth, strconv.Itoa(depth))
	return nil
}

// GetEdgeLimit gets the current value of edge limit
func (obj *TGQueryOptionImpl) GetEdgeLimit() int {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQueryEdgeLimit) {
		cn := utils.NewConfigName(OptionQueryEdgeLimit, OptionQueryEdgeLimit, strconv.Itoa(DefaultEdgeLimit))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultEdgeLimit)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultEdgeLimit
}

// SetEdgeLimit sets the number of edges per node to be returned in a query.  Default is 0 which means unlimited.
func (obj *TGQueryOptionImpl) SetEdgeLimit(limit int) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if limit == 0 {
		limit = -1
	}
	obj.SetProperty(OptionQueryEdgeLimit, strconv.Itoa(limit))
	return nil
}

// GetSortAttrName gets sort attribute name
func (obj *TGQueryOptionImpl) GetSortAttrName() string {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQuerySortAttr) {
		cn := utils.NewConfigName(OptionQuerySortAttr, OptionQuerySortAttr, "")
		if s := obj.GetProperty(cn, ""); s != "" {
			return s
		}
	}
	return ""
}

// SetSortAttrName sets sort attribute name
func (obj *TGQueryOptionImpl) SetSortAttrName(name string) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if name == "" || len(name) == 0 {
		errMsg := fmt.Sprintf("AttributeTypeInvalid attribute name")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	obj.SetProperty(OptionQuerySortAttr, name)
	return nil
}

// IsSortOrderDsc gets sort order desc
func (obj *TGQueryOptionImpl) IsSortOrderDsc() bool {
	if utils.DoesPropertyExist(obj.optionProperties, OptionQuerySortOrder) {
		cn := utils.NewConfigName(OptionQuerySortOrder, OptionQuerySortOrder, "false")
		if s := obj.GetProperty(cn, "false"); s != "" {
			v, _ := strconv.ParseBool(s)
			return v
		}
	}
	return false
}

// SetSortOrderDsc sets sort order desc
func (obj *TGQueryOptionImpl) SetSortOrderDsc(isDsc bool) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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
	if utils.DoesPropertyExist(obj.optionProperties, OptionQuerySortLimit) {
		cn := utils.NewConfigName(OptionQuerySortLimit, OptionQuerySortLimit, strconv.Itoa(DefaultOptionSortLimit))
		if s := obj.GetProperty(cn, strconv.Itoa(DefaultOptionSortLimit)); s != "" {
			v, _ := strconv.ParseInt(s, 10, 32)
			return int(v)
		}
	}
	return DefaultOptionSortLimit
}

// SetSortResultLimit sets sort result limit
func (obj *TGQueryOptionImpl) SetSortResultLimit(limit int) types.TGError {
	if !obj.mutable {
		errMsg := fmt.Sprintf("Can't modify a immutable Option")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	if limit <= 0 {
		errMsg := fmt.Sprintf("Invalid sort limit")
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
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

// AddProperty checks whether a property already exists, else adds a new property in the form of name=value pair
func (obj *TGQueryOptionImpl) AddProperty(name, value string) {
	//logger.Log(fmt.Sprintf("TGQueryOptionImpl::AddProperty received N-V Pair as '%+v':'%+v'", name, value))
	obj.optionProperties.AddProperty(name, value)
}

// GetProperty gets the property either with value or default value
func (obj *TGQueryOptionImpl) GetProperty(cn types.TGConfigName, value string) string {
	return obj.optionProperties.GetProperty(cn, value)
}

// SetProperty sets existing property value in the form of name=value pair
func (obj *TGQueryOptionImpl) SetProperty(name, value string) {
	obj.optionProperties.SetProperty(name, value)
}

// GetPropertyAsInt gets Property as int value
func (obj *TGQueryOptionImpl) GetPropertyAsInt(cn types.TGConfigName) int {
	return obj.optionProperties.GetPropertyAsInt(cn)
}

// GetPropertyAsLong gets Property as long value
func (obj *TGQueryOptionImpl) GetPropertyAsLong(cn types.TGConfigName) int64 {
	return obj.optionProperties.GetPropertyAsLong(cn)
}

// GetPropertyAsBoolean gets Property as bool value
func (obj *TGQueryOptionImpl) GetPropertyAsBoolean(cn types.TGConfigName) bool {
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
