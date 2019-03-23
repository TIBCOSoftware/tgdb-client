package channel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"sync"
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
 * File name: BlockingChannelResponse.go
 * Created on: Dec 22, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type BlockingChannelResponse struct {
	status    types.ChannelResponseStatus
	requestId int64
	timeout   int64
	reply     types.TGMessage
	lock      sync.Mutex // reentrant-lock for synchronizing sending/receiving messages over the wire
	cond      *sync.Cond // Condition for lock
}

func DefaultBlockingChannelResponse(reqId int64) *BlockingChannelResponse {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(BlockingChannelResponse{})

	newBlockingChannelResponse := BlockingChannelResponse{
		requestId: reqId,
		status:    types.Waiting,
		timeout:   -1,
	}
	newBlockingChannelResponse.cond = sync.NewCond(&newBlockingChannelResponse.lock) // Condition for lock

	return &newBlockingChannelResponse
}

func NewBlockingChannelResponse(reqId, rTimeout int64) *BlockingChannelResponse {
	newBlockingChannelResponse := DefaultBlockingChannelResponse(reqId)
	newBlockingChannelResponse.timeout = rTimeout
	return newBlockingChannelResponse
}

/////////////////////////////////////////////////////////////////
// Private functions for BlockingChannelResponse
/////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////
// Implement functions for TGChannelResponse
/////////////////////////////////////////////////////////////////

// Await waits (loops) till the channel response receives reply message from the server
func (obj *BlockingChannelResponse) Await(tester types.StatusTester) {
	//logger.Log(fmt.Sprintf("Entering BlockingChannelResponse:Await w/ contents as '%+v'", obj.String()))
	logger.Log(fmt.Sprint("Entering BlockingChannelResponse:Await"))
	//obj.lock.Lock()
	//defer obj.lock.Unlock()

	//go func() {
	count := 0
	for {
		// Terminating Condition for this Infinite Loop is:
		// 	(a) Break if the channel response object status is NOT WAITING - Status is set via SetReply()/Signal() execution
		if !tester.Test(obj.status) {
			logger.Log(fmt.Sprintf("Breaking out from BlockingChannelResponse:Await w/ contents as '%+v'", obj.String()))
			break
		}
		//obj.cond.Wait()
		time.Sleep(time.Duration(obj.timeout) * time.Millisecond)
		//obj.cond.Signal()
		// TODO: Remove this block once testing is over
		count++
		if (count%10000) == 0 {
			logger.Log(fmt.Sprintf("Inside BlockingChannelResponse:Await(%d) ... BlockingChannelResponse:status = Waiting", count))
		}
	}
	//}()

	logger.Log(fmt.Sprintf("Returning BlockingChannelResponse:Await ..."))
}

// GetCallback gets a Callback object
func (obj *BlockingChannelResponse) GetCallback() types.Callback {
	// Not applicable / available for BlockingChannelResponse
	return nil
}

// GetReply gets Reply object
func (obj *BlockingChannelResponse) GetReply() types.TGMessage {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	return obj.reply
}

// GetRequestId gets Request id
func (obj *BlockingChannelResponse) GetRequestId() int64 {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	return obj.requestId
}

// GetStatus gets Status
func (obj *BlockingChannelResponse) GetStatus() types.ChannelResponseStatus {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	return obj.status
}

// IsBlocking checks whether this channel response is blocking or not
func (obj *BlockingChannelResponse) IsBlocking() bool {
	return true
}

// Reset resets the state of channel response and initializes everything
func (obj *BlockingChannelResponse) Reset() {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:Reset ..."))
	obj.status = types.Waiting
	obj.reply = nil
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:Reset ..."))
}

// SetReply sets the reply message received from the server
func (obj *BlockingChannelResponse) SetReply(msg types.TGMessage) {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:SetReply ..."))
	obj.reply = msg
	obj.status = types.Ok
	obj.cond.Broadcast()
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:SetReply ..."))
}

// SetRequestId sets Request id
func (obj *BlockingChannelResponse) SetRequestId(reqId int64) {
	//obj.lock.Lock()
	//defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:SetRequestId ..."))
	obj.requestId = reqId
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:SetRequestId ..."))
}

// Signal lets other listeners of channel response know the status of this channel response
func (obj *BlockingChannelResponse) Signal(cStatus types.ChannelResponseStatus) {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	//logger.Log(fmt.Sprint("Entering BlockingChannelResponse:Signal ..."))
	obj.status = cStatus
	obj.cond.Broadcast()
	//logger.Log(fmt.Sprint("Returning BlockingChannelResponse:Signal ..."))
}

func (obj *BlockingChannelResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("BlockingChannelResponse:{")
	buffer.WriteString(fmt.Sprintf("Status: %+v", obj.status))
	buffer.WriteString(fmt.Sprintf(", RequestId: %d", obj.requestId))
	buffer.WriteString(fmt.Sprintf(", Timeout: %d", obj.timeout))
	buffer.WriteString(fmt.Sprintf(", Cond: %+v", obj.cond))
	buffer.WriteString(fmt.Sprintf(", Reply: %+v", obj.reply))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions for StatusTester
/////////////////////////////////////////////////////////////////

// Test checks whether the channel response is in WAIT mode or not
func (obj *BlockingChannelResponse) Test(status types.ChannelResponseStatus) bool {
	obj.lock.Lock()
	defer obj.lock.Unlock()
	if obj.status == types.Waiting {
		return true
	}
	return false
}
