package channel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/pdu"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sync/atomic"
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
 * File name: ChannelReader.go
 * Created on: Dec 22, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

var gReaders int64

type ChannelReader struct {
	channel   types.TGChannel
	isRunning bool
	name      string
	readerNum int64
}

func DefaultChannelReader() *ChannelReader {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelReader{})

	newChannelReader := ChannelReader{
		isRunning: false,
		readerNum: atomic.AddInt64(&gReaders, 1),
	}

	return &newChannelReader
}

func NewChannelReader(rChannel types.TGChannel) *ChannelReader {
	newChannelReader := DefaultChannelReader()
	newChannelReader.channel = rChannel
	newChannelReader.name = fmt.Sprintf("TGLinkReader@[%s-%d]", rChannel.GetClientId(), newChannelReader.readerNum)
	return newChannelReader
}

/////////////////////////////////////////////////////////////////
// Helper functions for ChannelReader
/////////////////////////////////////////////////////////////////

// readAndProcessLoop reads a message from the network and processes it
func (obj *ChannelReader) readAndProcessLoop() {
	if obj == nil {
		return
	}
	//logger.Log(fmt.Sprintf("Entering ChannelReader:readAndProcessLoop w/ Reader object: '%+v'", obj.String()))
	for {
		// Terminating Conditions for this Infinite Loop are:
		// 	(a) Break if the channel reader is NOT RUNNING
		// 	(b) Break if the channel reader / GO Routine is INTERRUPTED
		// 	(c) Break if the channel is CLOSED
		// 	(d) Break if the request on the wire is to DISCONNECT from the SERVER
		// 	(e) Break - in case of ERROR - if the exceptionResult is NOT RetryOperation
		// Looping Conditions for this Infinite Loop are:
		// 	(a) Continue if the message on the wire is EMPTY or cannot be READ
		// 	(b) Continue if the message on the wire is PING (HeartBeat) message w/o Processing the message
		// 	(c) Continue if the message on the wire is ANYTHING else after Processing the message
		// 	(d) Continue - in case of ERROR - if the exceptionResult is ANYTHING OTHER THAN RetryOperation after setting the reply on channelResponse

		if !obj.isRunning {
			logger.Log(fmt.Sprintf("Breaking ChannelReader:readAndProcessLoop loop since reader is not running '%+v'", obj.isRunning))
			break
		}

		if obj.channel.IsClosed() {
			logger.Log(fmt.Sprintf("Breaking ChannelReader:readAndProcessLoop loop since channel is closed"))
			break
		}

		// Execute Derived Channel's method
		msg, err := obj.channel.ReadWireMsg()
		if err != nil {
			logger.Error(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop Error in obj.channel.ReadWireMsg() w/ '%+v'", err.Error()))
			if !obj.isRunning {
				logger.Log(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop reader is not running (2) '%+v'", obj.isRunning))
				break
			}
			exceptionResult := channelHandleException(obj.channel, err, true)
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop Error in reading message - exceptionResult '%+v'", exceptionResult))
			for _, resp := range obj.channel.GetResponses() {
				resp.SetReply(pdu.NewExceptionMessageWithType(int(exceptionResult.ExceptionType), exceptionResult.ExceptionMessage))
			}
			if exceptionResult.ExceptionType != RetryOperation {
				// AbstractChannel.gLogger.logException("Returning channel Reader thread", e);
				logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop since Reader thread returned w/o Retrying due to error - exceptionResult '%+v'", exceptionResult))
				break
			}
			logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop - Read Wire Message resulted in error: '%+v'", err))
			//break
		}

		if msg == nil {
			logger.Log(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop - Read Message Again since MSG is NIL"))
			continue
		}

		logger.Log(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop - Read Message of type '%+v'", msg.GetVerbId()))
		if msg.GetVerbId() == pdu.VerbPingMessage {
			logger.Log(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop Trying to Read Message Again since MSG is PingMessage"))
			continue
		}

		// Server Requested to disconnect
		if msg.GetVerbId() == pdu.VerbSessionForcefullyTerminated {
			logger.Log(fmt.Sprintf("Breaking ChannelReader:readAndProcessLoop loop w/ Forceful Termination Message is '%+v'", msg.String()))
			channelTerminated(obj.channel, msg.(*pdu.SessionForcefullyTerminatedMessage).GetKillString())
			obj.isRunning = false
			break
		}

		logger.Log(fmt.Sprintf("Inside ChannelReader:readAndProcessLoop Processing Message of type '%+v'", msg.GetVerbId()))
		err = channelProcessMessage(obj.channel, msg)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop Error in channelProcessMessage() w/ '%+v'", err.Error()))
			if !obj.isRunning {
				logger.Log(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop reader is not running (3) '%+v'", obj.isRunning))
				break
			}
			exceptionResult := channelHandleException(obj.channel, err, true)
			logger.Error(fmt.Sprintf("ERROR: Inside ChannelReader:readAndProcessLoop Error in reading message - exceptionResult (2) '%+v'", exceptionResult))
			for _, resp := range obj.channel.GetResponses() {
				resp.SetReply(pdu.NewExceptionMessageWithType(int(exceptionResult.ExceptionType), exceptionResult.ExceptionMessage))
			}
			if exceptionResult.ExceptionType != RetryOperation {
				// AbstractChannel.gLogger.logException("Returning channel Reader thread", e);
				logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop since Reader thread returned w/o Retrying due to error (2) - exceptionResult '%+v'", exceptionResult))
				break
			}
			logger.Error(fmt.Sprintf("ERROR: Breaking ChannelReader:readAndProcessLoop loop - ProcessMessage resulted in error: '%+v'", err))
			//break
		}

		if obj.channel.IsClosed() {
			logger.Log(fmt.Sprintf("Breaking ChannelReader:readAndProcessLoop loop since channel is closed"))
			break
		}
	} // End of Infinite Loop
	obj.isRunning = false
	logger.Log(fmt.Sprintf("Returning ChannelReader:readAndProcessLoop w/ Reader object: '%+v'", obj.String()))
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> ChannelReader
/////////////////////////////////////////////////////////////////

// Start starts the channel reader
func (obj *ChannelReader) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelReader:Start ..."))
	if !obj.isRunning {
		obj.isRunning = true
		// Start reading and processing messages from the wire
		obj.readAndProcessLoop()
		//go obj.readAndProcessLoop()
	}
	//logger.Log(fmt.Sprint("Returning ChannelReader:Start ..."))
}

// Stop stops the channel reader
func (obj *ChannelReader) Stop() {
	//logger.Log(fmt.Sprint("Entering ChannelReader:Stop ..."))
	if obj.isRunning {
		// Finish / Flush any remaining processing
		//obj.readAndProcessLoop()
		go obj.readAndProcessLoop()
		obj.isRunning = false
	}
	//logger.Log(fmt.Sprint("Returning ChannelReader:Stop ..."))
}

func (obj *ChannelReader) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelReader:{")
	buffer.WriteString(fmt.Sprintf("Name: %+v", obj.name))
	buffer.WriteString(fmt.Sprintf(", IsRunning: %+v", obj.isRunning))
	buffer.WriteString(fmt.Sprintf(", ReaderNum: %d", obj.readerNum))
	buffer.WriteString(fmt.Sprintf(", Channel: %s", obj.channel.String()))
	buffer.WriteString("}")
	return buffer.String()
}
