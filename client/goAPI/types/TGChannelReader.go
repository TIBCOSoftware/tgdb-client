package types

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
 * File name: TGChannelReader.go
 * Created on: Oct 27, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

// Channel Reader is an independent thread that starts and stops with the channel, and continuously monitors
// network communication socket to read and process any message that is sent by the TGDB server
type TGChannelReader interface {
	// Start starts the channel reader
	Start()
	// Stop stops the channel reader
	Stop()
	// Additional Method to help debugging
	String() string
}
