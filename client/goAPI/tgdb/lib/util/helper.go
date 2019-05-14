package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func CastGenMap(mapData interface{}) map[string]interface{} {
	if nil != mapData {
		return mapData.(map[string]interface{})
	}
	return make(map[string]interface{})
}

func CastGenArray(arrayData interface{}) []interface{} {
	if nil != arrayData {
		return arrayData.([]interface{})
	}
	return make([]interface{}, 0)
}

func CastString(stringData interface{}) string {
	if nil != stringData {
		return stringData.(string)
	}
	return ""
}

func StringToTypes(data string, dataType string, dateTimeSample string) (interface{}, error) {
	if "null" == data {
		return nil, nil
	}
	switch dataType {
	case "String":
		return data, nil
	case "Integer":
		return strconv.ParseInt(data, 10, 32)
	case "Long":
		return strconv.ParseInt(data, 10, 64)
	case "Boolean":
		return strconv.ParseBool(data)
	case "Double":
		return strconv.ParseFloat(data, 64)
	case "Date":
		return time.Parse(dateTimeSample, data)
	}
	return data, nil
}

func ToString(data interface{}, dataType string, dateTimeSample string) (string, error) {
	golangTypeString := reflect.TypeOf(data).String()

	fmt.Println("data = ", data, ", dataType = ", dataType, ", golangTypeString = ", golangTypeString)

	switch dataType {
	case "String":
		if "string" == golangTypeString {
			stringData := data.(string)
			return stringData, nil
		} else {
			return "", fmt.Errorf("Not a string type data")
		}
	case "Integer":
		if "int32" == golangTypeString {
			intData := int(data.(int32))
			return strconv.Itoa(intData), nil
		} else {
			return "", fmt.Errorf("Not a Integer(int32) type data")
		}
	case "Long":
		if "int64" == golangTypeString {
			longData := data.(int64)
			return strconv.FormatInt(longData, 10), nil
		} else {
			return "", fmt.Errorf("Not a Long(int64) type data")
		}
	case "Boolean":
		if "bool" == golangTypeString {
			booeanlData := data.(bool)
			return strconv.FormatBool(booeanlData), nil
		} else {
			return "", fmt.Errorf("Not a Boolean(bool) type data")
		}
	case "Double":
		if "float64" == golangTypeString {
			doubleData := data.(float64)
			return strconv.FormatFloat(doubleData, 'f', -1, 64), nil
		} else {
			return "", fmt.Errorf("Not a Double(float64) type data")
		}
	case "Date":
		if "time.Time" == golangTypeString {
			dateTimeData := data.(time.Time)
			return dateTimeData.Format(dateTimeSample), nil
		} else {
			return "", fmt.Errorf("Not a Date(time.Time) type data")
		}
	}
	return data.(string), nil
}

func ReplaceCharacter(str string, targetRegex string, replacement string, doReplace bool) string {
	if doReplace {
		var re = regexp.MustCompile(targetRegex)
		str = re.ReplaceAllString(str, replacement)
	}

	return str
}

func SliceContains(slice []string, targetElement string) bool {
	for _, element := range slice {
		if element == targetElement {
			return true
		}
	}
	return false
}

func IsInteger(data string) bool {
	_, err := strconv.ParseInt(data, 10, 64)
	return err == nil
}
