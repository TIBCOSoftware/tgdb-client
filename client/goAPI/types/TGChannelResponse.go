package types

import "bytes"

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
 * File name: TGChannelResponse.go
 * Created on: Oct 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// ======= Channel Response Status =======
type ChannelResponseStatus int

const (
	Waiting ChannelResponseStatus = iota
	Ok
	Pushed
	Resend
	Disconnected
	Closed
)

func (responseStatus ChannelResponseStatus) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if responseStatus&Waiting == Waiting {
		buffer.WriteString("Channel Response is Waiting")
	} else if responseStatus&Ok == Ok {
		buffer.WriteString("Channel Response is Ok")
	} else if responseStatus&Pushed == Pushed {
		buffer.WriteString("Channel Response is Pushed")
	} else if responseStatus&Resend == Resend {
		buffer.WriteString("Channel Response is Resend")
	} else if responseStatus&Disconnected == Disconnected {
		buffer.WriteString("Channel Response is Disconnected")
	} else if responseStatus&Closed == Closed {
		buffer.WriteString("Channel Response is Closed")
	}
	return buffer.String()
}

// Channel Response is an independent thread that starts and stops with the channel, and continuously monitors
// whether server has replied with a message event or not
type TGChannelResponse interface {
	// Await waits (loops) till the channel response receives reply message from the server
	Await(tester StatusTester)
	// GetCallback gets a Callback object
	GetCallback() Callback
	// GetReply gets Reply object
	GetReply() TGMessage
	// GetRequestId gets Request id
	GetRequestId() int64
	// GetStatus gets Status
	GetStatus() ChannelResponseStatus
	// IsBlocking checks whether this channel response is blocking or not
	IsBlocking() bool
	// Reset resets the state of channel response and initializes everything
	Reset()
	// SetReply sets the reply message received from the server
	SetReply(msg TGMessage)
	// SetRequestId sets Request id
	SetRequestId(requestId int64)
	// Signal lets other listeners of channel response know the status of this channel response
	Signal(status ChannelResponseStatus)
}

type Callback interface {
	OnResponse(msg TGMessage)
}

type StatusTester interface {
	// Test checks whether the channel response is in WAIT mode or not
	Test(Status ChannelResponseStatus) bool
}
