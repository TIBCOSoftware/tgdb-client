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
 * File name: TGConnectionPool.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type TGConnectionPool interface {
	// Connect establishes connection from this pool of available/configured connections to the TGDB server
	// Exception could be BadAuthentication or BadUrl
	Connect() TGError
	// Disconnect breaks the connection from the TGDB server and returns the connection back to this connection pool for reuse
	Disconnect() TGError
	// Get returns a free connection from the connection pool
	// The property ConnectionReserveTimeoutSeconds or tgdb.connectionpool.ConnectionReserveTimeoutSeconds specifies the time
	// to wait in seconds. It has the following meaning
	// 0 :      Indefinite
	// -1 :     Immediate
	// &gt; :   That many seconds
	Get() (TGConnection, TGError)
	// GetPoolSize gets pool size
	GetPoolSize() int
	// ReleaseConnection frees the connection and sends back to the pool
	ReleaseConnection(conn TGConnection) (TGConnectionPool, TGError)
	// SetExceptionListener sets exception listener
	SetExceptionListener(lsnr TGConnectionExceptionListener)
}
