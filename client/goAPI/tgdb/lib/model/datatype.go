package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DataType int

const (
	TypeString DataType = iota
	TypeInteger
	TypeLong
	TypeDouble
	TypeBoolean
	TypeDate
)

var types = [...]string{
	"String",
	"Integer",
	"Long",
	"Double",
	"Boolean",
	"Date",
}

func (dt DataType) String() string {
	return types[dt]
}

// ToTypeEnum get the data type that corresponds to the specified name
func ToTypeEnum(typeStr string) (DataType, bool) {

	fmt.Println("(ToTypeEnum) typeStr - ", typeStr)

	switch strings.ToLower(typeStr) {
	case "string":
		return TypeString, true
	case "integer", "int":
		return TypeInteger, true
	case "long":
		return TypeLong, true
	case "double", "number":
		return TypeDouble, true
	case "boolean", "bool":
		return TypeBoolean, true
	case "date":
		return TypeDate, true
	default:
		return TypeString, false
	}
}

// GetType get the Type of the supplied value
func GetDataType(val interface{}) (DataType, error) {

	switch t := val.(type) {
	case string:
		return TypeString, nil
	case int, int32:
		return TypeInteger, nil
	case int64:
		return TypeLong, nil
	case float64:
		return TypeDouble, nil
	case json.Number:
		if strings.Contains(t.String(), ".") {
			return TypeDouble, nil
		} else {
			return TypeLong, nil
		}
	case bool:
		return TypeBoolean, nil
	case time.Time:
		return TypeDate, nil
	default:
		return TypeString, fmt.Errorf("unable to determine type of %#v", t)
	}
}

func IsSimpleType(val interface{}) bool {

	switch val.(type) {
	case string, int, int32, float32, float64, json.Number, bool:
		return true
	default:
		return false
	}
}
