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
 * File name: ChannelTracer.go
 * Created on: Apr 13, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

package channel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
)

type ChannelTracer struct {
	msgQueue  *utils.SimpleQueue
	msgTracer *ChannelMessageTracer
	clientId  string
	isRunning bool
}

func DefaultChannelTracer() *ChannelTracer {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(ChannelTracer{})

	newChannelTracer := ChannelTracer{
		msgQueue: utils.NewSimpleQueue(),
		clientId: "",
		isRunning: false,
	}

	return &newChannelTracer
}

func NewChannelTracer(client, traceDir string) *ChannelTracer {
	newChannelTracer := DefaultChannelTracer()
	newChannelTracer.clientId = client
	msgTracer := NewChannelMessageTracer(newChannelTracer.msgQueue, client, traceDir)
	newChannelTracer.msgTracer = msgTracer
	return newChannelTracer
}

/////////////////////////////////////////////////////////////////
// Private functions for ChannelTracer
/////////////////////////////////////////////////////////////////

func (obj *ChannelTracer) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("ChannelTracer:{")
	buffer.WriteString(fmt.Sprintf("ClientId: %+v", obj.clientId))
	buffer.WriteString(fmt.Sprintf(", MsgQueue: %d", obj.msgQueue))
	buffer.WriteString(fmt.Sprintf(", MsgTracer: %d", obj.msgTracer))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for TGTracer
/////////////////////////////////////////////////////////////////

// Start starts the channel tracer
func (obj *ChannelTracer) Start() {
	//logger.Log(fmt.Sprint("Entering ChannelTracer:Start ..."))
	if !obj.isRunning {
		obj.msgTracer.Start()
		obj.isRunning = true
	}
	//logger.Log(fmt.Sprint("Returning ChannelTracer:Start ..."))
}

// Stop stops the channel tracer
func (obj *ChannelTracer) Stop() {
	logger.Log(fmt.Sprint("Entering ChannelTracer:Stop ..."))
	if obj.isRunning {
		// Finish / Flush any remaining processing
		obj.msgTracer.Stop()
		obj.isRunning = false
	}
	logger.Log(fmt.Sprint("Returning ChannelTracer:Stop ..."))
}

// Trace traces the path the message has taken
func (obj *ChannelTracer) Trace(msg types.TGMessage) {
	//logger.Log(fmt.Sprint("Entering ChannelTracer:Trace"))
	obj.msgQueue.Enqueue(msg)
	//logger.Log(fmt.Sprintf("Returning ChannelTracer:Trace ..."))
}
