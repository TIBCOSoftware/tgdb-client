package utils

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/logging"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sort"
	"strconv"
	"strings"
	"sync"
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
 * File name: TGProperties.go (Always Sorted by key)
 * Created on: Oct 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/******* Sample code tested in GO Playground *******
package main

import (
"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sort"
	"strings"
	"sync"
)

type PropSet struct {
	Pairs []*NVPair
}

type NVPair struct {
	key string
	Val string
}

func ByKey(c1, c2 *NVPair) bool {
	return strings.ToLower(c1.key) < strings.ToLower(c2.key)
}

func ByVal(c1, c2 *NVPair) bool {
	return strings.ToLower(c1.Val) < strings.ToLower(c2.Val)
}

func GetVal(obj *PropSet, key string) string {
	for _, kvp := range obj.Pairs {
		if kvp.key == key {
			return kvp.Val
		}
	}
	return "key Not Found"
}

func SetVal(obj *PropSet, key string, val string) {
	for _, kvp := range obj.Pairs {
		if kvp.key == key {
			kvp.Val = val
		}
	}
}

func main() {
	fmt.Println("Hello, playground")
	a := &NVPair{"Bob", "SportsFan"}
	b := &NVPair{"Jim", "BooksFan"}
	if ByKey(a, b) {
		fmt.Println("key Comparison: Bob < Jim")
	} else {
		fmt.Println("key Comparison: Bob > Jim")
	}
	if ByVal(a, b) {
		fmt.Println("Val Comparison: Bob < Jim")
	} else {
		fmt.Println("Val Comparison: Bob > Jim")
	}

	nvPairs := []*NVPair{a, b}
	ps := &PropSet{Pairs: nvPairs}
	fmt.Println(ps)
	SetVal(ps, "Bob", "Sports Fanatic")
	newVal := GetVal(ps, "Bob")
	fmt.Printf("NewTGDecimal value set is '%s'\n", newVal)
}
*/

var logger = logging.DefaultTGLogManager().GetLogger()

type KvPair struct {
	KeyName  string
	KeyValue string
}

type sortFunc func(p1, p2 *KvPair) bool

type SortedProperties struct {
	properties   []*KvPair
	sortHandlers []sortFunc // Intentionally kept Private
	mutex        sync.Mutex // rw-lock for synchronizing read-n-update of env configuration
}

// Define Sort Handler functions
// Sort by key name
var kvKey = func(c1, c2 *KvPair) bool {
	return strings.ToLower(c1.KeyName) < strings.ToLower(c2.KeyName)
}

//var kvKey = func(c1, c2 string) bool {
//	return strings.ToLower(c1) < strings.ToLower(c2)
//}

// Make sure that the ConfigName implements the TGConfigName interface
var _ types.TGProperties = (*SortedProperties)(nil)

func defaultSortedProperties() *SortedProperties {
	return &SortedProperties{
		properties:   make([]*KvPair, 0),
		sortHandlers: make([]sortFunc, 0),
	}
}

func NewSortedProperties() *SortedProperties {
	newPropertySet := defaultSortedProperties()
	// Always return a sorted array of properties
	//go orderedBy(newPropertySet, kvKey).Sort(newPropertySet.Properties)
	return newPropertySet
}

/////////////////////////////////////////////////////////////////
// Private functions for SortedProperties
/////////////////////////////////////////////////////////////////

// Usage: orderedBy(kvKey).Sort([]KvPair)
// Usage: orderedBy(kvVal).Sort([]KvPair)
func orderedBy(obj *SortedProperties, sorters ...sortFunc) *SortedProperties {
	obj.sortHandlers = sorters
	return obj
}

func addProperty(obj *SortedProperties, name, value string) {
	//logger.Log(fmt.Sprintf("Entering SortedProperties:addProperty received name '%+v' and value as '%+v'", name, value))
	// Check if the property already exists, if it does, just update the value, else insert it into the set
	if DoesPropertyExist(obj, name) {
		setProperty(obj, name, value)
	}
	// NewTGDecimal Property
	newKVPair := KvPair{KeyName: name, KeyValue: value}
	obj.properties = append(obj.properties, &newKVPair)
	//logger.Log(fmt.Sprintf("Returning SortedProperties:addProperty has properties as '%+v'", obj.Properties))
}

func getProperty(obj *SortedProperties, conf types.TGConfigName, value string) string {
	//logger.Log(fmt.Sprintf("Entering SortedProperties:getProperty received '%+v' and substitute value as '%+v'", conf, value))
	var propVal string
	cn := conf.(*ConfigName)
	if cn == nil || cn.GetName() == "" || cn.GetAlias() == "" {
		return ""
	}
	//logger.Log(fmt.Sprintf("Inside SortedProperties:getProperty obj has properties as '%+v'\n", obj.Properties))
	// Search whether incoming configName has an associated value in existing NV pairs or not
	for _, kvp := range obj.properties {
		//logger.Log(fmt.Sprintf("Inside SortedProperties:getProperty kvp as '%+v'\n", kvp)
		if kvp.KeyName == cn.GetName() || kvp.KeyName == cn.aliasName {
			logger.Log(fmt.Sprintf("Inside SortedProperties:getProperty FOUND a config MATCH w/ kvp as '%+v'", kvp))
			propVal = kvp.KeyValue
			break
		}
	}
	if propVal == "" {
		logger.Log(fmt.Sprintf("Returning SortedProperties:getProperty DID NOT FIND a config MATCH - hence returning substitute value as '%+v'", value))
		return value
	}
	//logger.Log(fmt.Sprintf("Returning SortedProperties:getProperty w/ Config '%+v' and value as '%+v'", conf, propVal))
	return propVal
}

func setProperty(obj *SortedProperties, name, value string) {
	//logger.Log(fmt.Sprintf("Entering SortedProperties:setProperty received obj '%+v' name '%+v' and value as '%+v'", obj, name, value))
	// Check if the property already exists, if it does, just update the value, else insert it into the set
	if DoesPropertyExist(obj, name) {
		// Existing property - so set the new value
		for _, kvp := range obj.properties {
			if strings.ToLower(kvp.KeyName) != strings.ToLower(name) {
				continue
			}
			kvp.KeyValue = value
		}
	}
	//logger.Log(fmt.Sprintf("Returning SortedProperties:setProperty returning property set as '%+v'", obj.Properties))
}

/////////////////////////////////////////////////////////////////
// Helper Public functions for SortedProperties
/////////////////////////////////////////////////////////////////

func DoesPropertyExist(obj *SortedProperties, name string) bool {
	//flogger.Log(fmt.Sprintf("Entering SortedProperties:DoesPropertyExist searching '%+v' in properties as '%+v'", name, obj.Properties))
	for _, kvp := range obj.properties {
		if strings.ToLower(kvp.KeyName) == strings.ToLower(name) {
			//logger.Log(fmt.Sprintf("Returning SortedProperties:DoesPropertyExist as Property '%+v' Exists in properties as '%+v'", name, obj.Properties))
			return true
		}
	}
	//logger.Log(fmt.Sprintf("Returning SortedProperties:DoesPropertyExist as Property '%+v' does not Exist in properties as '%+v'", name, obj.Properties))
	return false
}

func SetUserAndPassword(obj *SortedProperties, user, pwd string) types.TGError {
	err := obj.SetUser(user)
	if err != nil {
		return err
	}
	err = obj.SetPassword(pwd)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SortedProperties) GetAllProperties() []*KvPair {
	return obj.properties
}

func (obj *SortedProperties) SetUser(user string) types.TGError {
	userConfig := GetConfigFromKey(ChannelUserID)
	if len(user) < 1 {
		u := obj.GetProperty(userConfig, NewTGEnvironment().GetChannelDefaultUser())
		if u == "" {
			return types.NewTGDBError("", types.TGErrorBadAuthentication, "Username not specified", "")
		}
		user = u
	}
	// AddProperty either sets the property or adds it
	obj.AddProperty(userConfig.GetName(), user)
	return nil
}

func (obj *SortedProperties) SetPassword(pwd string) types.TGError {
	pwdConfig := GetConfigFromKey(ChannelPassword)
	if len(pwd) < 1 {
		p := obj.GetProperty(pwdConfig, "")
		pwd = p
	}
	// AddProperty either sets the property or adds it
	obj.AddProperty(pwdConfig.GetName(), pwd)
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> sort.Interface
/////////////////////////////////////////////////////////////////

//func customSort(obj *SortedProperties) {
//	obj.mutex.Lock()
//	defer obj.mutex.Unlock()
//	var ss []KvPair
//	for k, v := range obj.Prop {
//		ss = append(ss, KvPair{k, v})
//	}
//
//	sort.Slice(ss, func(i, j int) bool {
//		return ss[i].KeyName < ss[j].KeyName
//	})
//}

func (obj *SortedProperties) Len() int {
	return len(obj.properties)
}

func (obj *SortedProperties) Less(i, j int) bool {
	p, q := obj.properties[i], obj.properties[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(obj.sortHandlers)-1; k++ {
		less := obj.sortHandlers[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return obj.sortHandlers[k](p, q)
}

func (obj *SortedProperties) Sort(props []*KvPair) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.properties = props
	sort.Sort(obj)
}

func (obj *SortedProperties) Swap(i, j int) {
	obj.properties[i], obj.properties[j] = obj.properties[j], obj.properties[i]
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGProperties
/////////////////////////////////////////////////////////////////

// AddProperty checks whether a property already exists, else adds a new property in the form of name=value pair
func (obj *SortedProperties) AddProperty(name, value string) {
	addProperty(obj, name, value)
	// Always return a sorted array of properties
	orderedBy(obj, kvKey).Sort(obj.properties)
	//logger.Log(fmt.Sprintf("Returning SortedProperties:AddProperty has properties as '%+v'", obj.Properties))
}

// GetProperty gets the property either with value or default value
func (obj *SortedProperties) GetProperty(conf types.TGConfigName, value string) string {
	propVal := getProperty(obj, conf, value)
	return propVal
}

// GetPropertyAsInt gets Property as int value
func (obj *SortedProperties) GetPropertyAsBoolean(conf types.TGConfigName) bool {
	value := obj.GetProperty(conf, "")
	if value != "" {
		//return value.(bool)
		v, _ := strconv.ParseBool(value)
		return v
	}
	return false
}

// GetPropertyAsInt gets Property as int value
func (obj *SortedProperties) GetPropertyAsInt(conf types.TGConfigName) int {
	value := obj.GetProperty(conf, "")
	if value != "" {
		//return value.(int)
		v, _ := strconv.Atoi(value)
		return v
	}
	return 0
}

// GetPropertyAsLong gets Property as long value
func (obj *SortedProperties) GetPropertyAsLong(conf types.TGConfigName) int64 {
	value := obj.GetProperty(conf, "")
	if value != "" {
		//return value.(int64)
		v, _ := strconv.ParseInt(value, 10, 64)
		return v
	}
	return 0
}

// GetPropertyAsBoolean gets Property as bool value
func (obj *SortedProperties) SetProperty(name, value string) {
	setProperty(obj, name, value)
	// Always return a sorted array of properties
	orderedBy(obj, kvKey).Sort(obj.properties)
}
