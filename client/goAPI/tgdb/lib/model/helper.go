/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/tgdb/lib/util"
)

//-====================-//
//    Globle Function
//-====================-//

func NewNodeId(nodeType string, nodeKey []interface{}) NodeId {
	var nodeId NodeId
	nodeId._keyHash = Hash(nodeKey)
	nodeId._type = nodeType
	return nodeId
}

func Hash(key []interface{}) string {
	keyBytes := []byte{}
	for _, element := range key {
		elementBytes, _ := json.Marshal(element)
		keyBytes = append(keyBytes, elementBytes...)
	}
	hasher := md5.New()
	hasher.Write(keyBytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func ToString(attribute *Attribute) string {
	if nil == attribute {
		return ""
	}

	return strings.TrimSpace(util.CastString(attribute.GetValue()))
}
