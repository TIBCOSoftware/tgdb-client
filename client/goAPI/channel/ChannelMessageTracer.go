package channel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
	"os"
	"path/filepath"
	"syscall"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: ChannelMessageTracer.go
 * Created on: Apr 13, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type ChannelMessageTracer struct {
	currentSuffix int
	traceFile     *os.File
	isRunning     bool
	msgQueue      *utils.SimpleQueue
	traceFileName string
}

const MaxFileSize int64 = 1 << 20

func DefaultChannelMessageTracer() *ChannelMessageTracer {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelMessageTracer{})

	newChannelMessageTracer := ChannelMessageTracer{
		currentSuffix: 0,
		isRunning:     false,
		msgQueue:      utils.NewSimpleQueue(),
		traceFileName: "",
	}

	return &newChannelMessageTracer
}

func NewChannelMessageTracer(queue *utils.SimpleQueue, client, traceDir string) *ChannelMessageTracer {
	newChannelMessageTracer := DefaultChannelMessageTracer()
	newChannelMessageTracer.msgQueue = queue
	newChannelMessageTracer.traceFileName = filepath.FromSlash(fmt.Sprint(traceDir, "/", client, ".trace"))
	newChannelMessageTracer.createTraceFile(0, newChannelMessageTracer.currentSuffix)
	return newChannelMessageTracer
}

/////////////////////////////////////////////////////////////////
// Private functions for ChannelMessageTracer
/////////////////////////////////////////////////////////////////

// exists returns whether the given file or directory exists, and whether it is a directory or not
func exists(path string) (bool, bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, false, err
	}
	if os.IsNotExist(err) {
		return false, false, nil
	}
	if info.IsDir() {
		return true, true, nil
	}
	return true, false, err
}

// isFileReadyForRollover checks if the file needs to be rolled over with incremented suffix
func (obj *ChannelMessageTracer) isFileReadyForRollover(suffix, msgBufLen int) bool {
	if obj == nil {
		return false
	}

	fileWithSuffix := fmt.Sprintf("%s%d", obj.traceFileName, suffix)
	fileInfo, err := os.Stat(fileWithSuffix)
	if err != nil {
		return false
	}

	if fileInfo.Size()+int64(msgBufLen) < MaxFileSize {
		return false
	}
	return true
}

// createTraceFile creates a new trace file with new suffix
func (obj *ChannelMessageTracer) createTraceFile(oldSuffix, newSuffix int) bool {
	if obj == nil {
		return false
	}

	// First check the existence of file with old suffix and close it for writing
	traceFileWithOldSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, oldSuffix)
	oFlag, _, err := exists(traceFileWithOldSuffix)
	if err != nil {
		return false
	}

	if !oFlag {
		return false
	} else {
		_ = obj.traceFile.Sync()  // Flush
		_ = obj.traceFile.Close() // Close FD
		obj.traceFile = nil
	}

	// Next check the existence of file with new suffix and open it for writing
	traceFileWithNewSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, newSuffix)
	nFlag, _, err := exists(traceFileWithNewSuffix)
	if err != nil {
		return false
	}

	if nFlag {
		return false
	}

	fp, err := os.OpenFile(traceFileWithNewSuffix, syscall.O_RDWR, 0644)
	if err != nil {
		return false
	}
	obj.traceFile = fp
	obj.currentSuffix = newSuffix
	return true
}

// extractAndTraceMessage reads a message from the message queue and traces it in the trace file
func (obj *ChannelMessageTracer) extractAndTraceMessage() {
	if obj == nil {
		return
	}

	//logger.Log(fmt.Sprintf("Entering ChannelMessageTracer:extractAndTraceMessage w/ Message Tracer object: '%+v'", obj.String()))
	for {
		if !obj.isRunning {
			logger.Log(fmt.Sprintf("Breaking ChannelMessageTracer:extractAndTraceMessage loop since message tracer is not running '%+v'", obj.isRunning))
			break
		}

		// At this point, the trace file with suffix is expected to be ready for writing contents in it
		msg := obj.msgQueue.Dequeue()
		if msg != nil {
			msgBuf, msgLen, err := msg.(types.TGMessage).ToBytes()
			if err != nil {
				logger.Error(fmt.Sprintf("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in msg.ToBytes() w/ '%+v'", err.Error()))
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			// Check if the file is about to exceed the limit after appending current message buffer
			if obj.isFileReadyForRollover(obj.currentSuffix, msgLen) {
				// Create a new file with incremented suffix
				cFlag := obj.createTraceFile(obj.currentSuffix, obj.currentSuffix+1)
				if !cFlag {
					logger.Error(fmt.Sprint("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in obj.createTraceFile()"))
					time.Sleep(1000 * time.Millisecond)
					continue
				}
			}

			_, err1 := obj.traceFile.Write(msgBuf)
			if err1 != nil {
				logger.Error(fmt.Sprintf("ERROR: Inside ChannelMessageTracer:extractAndTraceMessage Error in obj.traceFile.Write() w/ '%+v'", err1.Error()))
				time.Sleep(1000 * time.Millisecond)
				continue
			}
		} else {
			logger.Debug(fmt.Sprint("Inside ChannelMessageTracer:extractAndTraceMessage - No pending messages in the queue"))
			time.Sleep(1000 * time.Millisecond)
			continue
		}
	} // End of Infinite Loop
	obj.isRunning = false
	//logger.Log(fmt.Sprintf("Returning ChannelMessageTracer:extractAndTraceMessage w/ Message Tracer object: '%+v'", obj.String()))
}

func (obj *ChannelMessageTracer) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelMessageTracer:{")
	buffer.WriteString(fmt.Sprintf("TraceFileName: %+v", obj.traceFileName))
	buffer.WriteString(fmt.Sprintf(", MsgQueue: %d", obj.msgQueue))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for TGTracer
/////////////////////////////////////////////////////////////////

// Start starts the channel message tracer
func (obj *ChannelMessageTracer) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelMessageTracer:Start ..."))
	if !obj.isRunning {
		obj.isRunning = true
		// Start reading and processing messages from the wire
		obj.extractAndTraceMessage()
	}
	//logger.Log(fmt.Sprint("Returning ChannelMessageTracer:Start ..."))
}

// Stop stops the channel message tracer
func (obj *ChannelMessageTracer) Stop() {
	//logger.Log(fmt.Sprint("Entering ChannelMessageTracer:Stop ..."))
	if obj.isRunning {
		// Finish / Flush any remaining processing
		traceFileWithSuffix := fmt.Sprintf("%s.%d", obj.traceFileName, obj.currentSuffix)
		flag, _, _ := exists(traceFileWithSuffix)
		if flag {
			_ = obj.traceFile.Sync()  // Flush
			_ = obj.traceFile.Close() // Close FD
			obj.traceFile = nil
		}
		obj.isRunning = false
	}
	//logger.Log(fmt.Sprint("Returning ChannelMessageTracer:Stop ..."))
}
