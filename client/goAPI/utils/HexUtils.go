package utils

import (
	"bytes"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: HexUtils.go
 * Created on: Dec 14, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

const (
	NullString string = "0000"
	Space      string = " "
	NewLine    string = "\r\n"
)

func FormatHex(byteArray []byte) (string, error) {
	if byteArray == nil {
		return NullString, nil
	}
	var buffer bytes.Buffer
	_, err := FormatHexToWriter(byteArray, buffer, 0)
	if err != nil {
		return NullString, nil
	}
	return buffer.String(), nil
}

func FormatHexForLength(byteArray []byte, actualLength int) (string, error) {
	if byteArray == nil {
		return NullString, nil
	}
	var buffer bytes.Buffer
	_, err := FormatHexToWriter(byteArray, buffer, actualLength)
	if err != nil {
		return NullString, nil
	}
	return buffer.String(), nil
}

func FormatHexToWriter(buf []byte, writer bytes.Buffer, actualLength int) (int, error) {
	return FormatHexToWriterInChunks(buf, writer, 48, actualLength)
}

func FormatHexToWriterInChunks(buf []byte, writer bytes.Buffer, lineLength int, actualLength int) (int, error) {
	bLen := len(buf)
	bNewLine := false
	lineNo := 1

	writer.WriteString("Formatted Byte Array:")
	writer.WriteString(NewLine)
	writer.WriteString(fmt.Sprintf("%08x", 0))
	writer.WriteString(Space)

	if actualLength > 0 {
		bLen = actualLength
	}
	for i := 0; i < bLen; i++ {
		if bNewLine {
			bNewLine = false
			writer.WriteString(NewLine)
			writer.WriteString(fmt.Sprintf("%08x", lineNo*lineLength))
			writer.WriteString(Space)
		}

		writer.WriteString(fmt.Sprintf("%02x", buf[i]))
		if (i+1)%2 == 0 {
			writer.WriteString(Space)
		}

		if (i+1)%lineLength == 0 {
			bNewLine = true
			lineNo += 1
			//writer.flush()
		}
	} // End of for loop
	return lineNo, nil
}
