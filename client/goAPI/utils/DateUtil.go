package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
 * File name: DateUtil.go
 * Created on: Dec 08, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	OUTPUTDESIREDTSLAYOUT = time.RFC3339              // Format is same as TROPOS Log TS format i.e. yyyy-mm-ddTHH:MM:SSZ<TZ Diff>
	SINCEEPOCH            = "19700101T12:00:00Z00.00" // Alternative form for UNITY Queries
	UNITYDESIREDLAYOUT    = "2006-Jan-02 15:04:05"
	UNITYDATETIMEFORMAT   = "yyyy-mmm-dd HH:MM:SS" // Should be in sync with UNITYDESIREDLAYOUT
	INPUTDESIREDTSLAYOUT  = "2006-01-02T15:04"
	DATETIMEFORMAT        = "yyyy-mm-ddTHH:MM" // Should be in sync with INPUTDESIREDTSLAYOUT
	INPUTDATEONLYLAYOUT   = "2006-01-02"
	DATEONLYFORMAT        = "yyyy-mm-dd" // Should be in sync with INPUTDATEONLYLAYOUT
	INPUTTIMEONLYLAYOUT   = "15:04"
	TIMEONLYFORMAT        = "HH:MM"                // Should be in sync with INPUTTIMEONLYLAYOUT
	WEBSERVER_LAYOUT      = "2006-Jan-02 15:04:05" // the input/output log format from web server log API
)

// Convienient predefined commonly-used time intervals
var PresetIntervals = struct {
	One_Min,
	Five_Mins,
	Thirty_Mins,
	One_Hr,
	Twelve_Hrs,
	One_Day,
	Seven_Days,
	Thirty_Days,
	One_Yr time.Duration
}{
	One_Min:     time.Duration(1 * time.Minute),
	Five_Mins:   time.Duration(5 * time.Minute),
	Thirty_Mins: time.Duration(30 * time.Minute),
	One_Hr:      time.Duration(1 * time.Hour),
	Twelve_Hrs:  time.Duration(12 * time.Hour),
	One_Day:     time.Duration(24 * time.Hour),
	Seven_Days:  time.Duration(7 * 24 * time.Hour),
	Thirty_Days: time.Duration(30 * 24 * time.Hour),
	One_Yr:      time.Duration(365 * 24 * time.Hour),
}

// FormatStringAsDateTime identifies the input format in which users have entered date or time or date/time
// and appropriately convert it into standard RFC3339 format, if user has not specified the format
func FormatStringAsDateTime(input string, outputDateFormat string) (string, error) {
	var output string = "InvalidFormat"
	if outputDateFormat == "" {
		outputDateFormat = OUTPUTDESIREDTSLAYOUT
	}
	if strings.Contains(input, "T") {
		// implies both date and time specified on the command line
		from, err := time.ParseInLocation(INPUTDESIREDTSLAYOUT, input, time.Local)
		if err != nil {
			return output, err
		}
		if outputDateFormat == SINCEEPOCH {
			output = strconv.FormatInt(from.Unix()*1000, 10)
		} else {
			output = from.Local().Format(outputDateFormat)
		}
	} else if strings.Contains(input, ":") {
		// implies only time specified on the command line - date should be assumed TODAY's date
		var currentTime = time.Now().Local()
		var hr, min int
		yr, mth, day := currentTime.Date()
		if len(input) == 4 {
			hr, _ = strconv.Atoi(FindSubstring(input, 0, 1))
			min, _ = strconv.Atoi(FindSubstring(input, 2, 2))
		} else if len(input) == 5 {
			hr, _ = strconv.Atoi(FindSubstring(input, 0, 2))
			min, _ = strconv.Atoi(FindSubstring(input, 3, 2))
		}
		if hr < 0 || hr > 23 || min < 0 || min > 59 {
			err := errors.New("Please enter valid values for Hour [00-23] and/or Minutes [00-59].")
			return output, err
		}
		from := time.Date(yr, mth, day, hr, min, 0, 0, time.Local)
		if outputDateFormat == SINCEEPOCH {
			output = strconv.FormatInt(from.Unix()*1000, 10)
		} else {
			output = from.Local().Format(outputDateFormat)
		}
	} else {
		// implies only date specified on the command line - time should be assumed NOW's time
		from, err := time.ParseInLocation(INPUTDATEONLYLAYOUT, input, time.Local)
		if err != nil {
			return output, err
		}
		if outputDateFormat == SINCEEPOCH {
			output = strconv.FormatInt(from.Unix()*1000, 10)
		} else {
			output = from.Local().Format(outputDateFormat)
		}
	}
	//	fmt.Printf("=====> input: %s \toutputDateFormat: %s \toutput: %s\n", input, outputDateFormat, output)
	return output, nil
}

// FormatDateAsDateTime is a convenience function that converts the input date into the date format
// the user has specified, if no format is specified, it uses the default time.RFC3339 format
func FormatDateAsDateTime(input time.Time, outputDateFormat string) string {
	var output string = "InvalidFormat"
	if outputDateFormat == "" {
		outputDateFormat = OUTPUTDESIREDTSLAYOUT
	}
	//	output = input.In(time.Local).Format(outputDateFormat)
	if outputDateFormat == SINCEEPOCH {
		output = strconv.FormatInt(input.Unix()*1000, 10)
	} else {
		output = input.Format(outputDateFormat)
	}
	//	fmt.Printf("=====> input: %s \toutputDateFormat: %s \toutput: %s\n", input, outputDateFormat, output)
	return output
}

// Generate a slice of times between 2 time points(exclusive), with the oldest time point first
//
// start - start time point
// end - end time point
// n - number of time points needed
// randomized - specify if generated timer points should be randomized
//
// If the n is such that the durations between time points fell below a nanosecond resolution, an empty
// slice is returned instead.
func Range(start, end time.Time, n uint, randomized bool) ([]time.Time, error) {
	start = start.Local()
	end = end.Local()
	if start.Equal(end) {
		return nil, fmt.Errorf("Start time '%+v' specified is same as end time '%+v'", start, end)
	}
	if start.After(end) {
		return nil, fmt.Errorf("Start time '%+v' specified is later than end time '%+v'", start, end)
	}
	d := end.Sub(start)
	var interval int64
	if n > 0 {
		interval = d.Nanoseconds() / int64(n)
	}
	if n == 0 || d.Nanoseconds() == 0 || interval == 0 {
		return []time.Time{}, nil
	}
	var sequences []time.Time
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for i := 1; i <= int(n); i++ {
		var t time.Time
		if !randomized {
			t = start.Add(time.Duration(int64(i) * interval))
		} else {
			t = start.Add(time.Duration(int64(i-1) * interval)).Add(time.Duration(interval / 2)).Add(time.Duration(r.Int63n(7 * interval / 8)))
		}
		sequences = append(sequences, t)
	}

	return sequences, nil
}

// Generate a slice of times from now into the past of specified interval, with the oldest time point first
//
// interval - past time interval from time.Now()
// n - number of time points needed
// randomized - specify if generated timer points should be randomized
//
// If the n is such that the durations between time points fell below nanosecond resolution, an empty
// slice is returned instead.
func RangePast(interval time.Duration, n uint, randomized bool) ([]time.Time, error) {
	end := time.Now().Local()
	start := end.Add(-1 * interval)
	return Range(start, end, n, randomized)
}

// Generate a slice of times from now into the future of specified interval, with the oldest time point first
//
// interval - future time interval from time.Now()
// n - number of time points needed
// randomized - specify if generated timer points should be randomized
//
// If the n is such that the durations between time points fell below nanosecond resolution, an empty
// slice is returned instead
func RangeFuture(interval time.Duration, n uint, randomized bool) ([]time.Time, error) {
	start := time.Now().Local()
	end := start.Add(interval)
	return Range(start, end, n, randomized)
}

// Rounded up time.Duration display string eliminating the precision
// Example: 500ms, 4s, 1h4m20s
//
func RoundDuration(d time.Duration) string {
	d_str := d.String()
	unit_suffix := extractDurationUnitSuffix(d)
	split_d_strs := strings.SplitN(d_str, ".", 2)
	if len(split_d_strs) == 1 {
		return d_str
	} else {
		return split_d_strs[0] + string(unit_suffix)
	}

}

// find substring within another string
func FindSubstring(s string, pos int, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func extractDurationUnitSuffix(d time.Duration) string {
	var unit string
	d_str := d.String()
	for i := len(d_str) - 1; i > 0; i-- {
		s := string(d_str[i])
		_, err := strconv.Atoi(s)
		if err != nil {
			unit += s
		} else {
			break
		}
	}
	runes := []rune(unit)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
