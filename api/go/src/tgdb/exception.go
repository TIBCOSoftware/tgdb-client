/*
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * File name: exception.go
 * Created on: 11/13/2019
 * Created by: nimish
 *
 * SVN Id: $Id: exception.go 4048 2020-06-01 18:53:50Z nimish $
 */

package tgdb

// ======= Various Error Types =======
type TGExceptionType int

type TGError interface {
	error
	// Get the detail error message
	GetErrorCode() string
	GetErrorType() int
	GetErrorMsg() string
	GetErrorDetails() string
	GetServerErrorCode() int
}