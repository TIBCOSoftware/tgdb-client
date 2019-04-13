package admin

import (
	"bytes"
	"encoding/gob"
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
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: UserInfoImpl.go
 * Created on: Apr 06, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type UserInfoImpl struct {
	userId   int
	userType byte
	userName string
}

// Make sure that the UserInfoImpl implements the TGUserInfo interface
var _ TGUserInfo = (*UserInfoImpl)(nil)

func DefaultUserInfoImpl() *UserInfoImpl {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(UserInfoImpl{})

	return &UserInfoImpl{}
}

func NewUserInfoImpl(_userId int, _userName string, _userType byte) *UserInfoImpl {
	newConnectionInfo := DefaultUserInfoImpl()
	newConnectionInfo.userId = _userId
	newConnectionInfo.userType = _userType
	newConnectionInfo.userName = _userName
	return newConnectionInfo
}

/////////////////////////////////////////////////////////////////
// Helper functions from Interface ==> TGUserInfoImpl
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("UserInfoImpl:{")
	buffer.WriteString(fmt.Sprintf("UserId: '%d'", obj.userId))
	buffer.WriteString(fmt.Sprintf(", UserType: '%+v'", obj.userType))
	buffer.WriteString(fmt.Sprintf(", UserName: '%s'", obj.userName))
	buffer.WriteString("}")
	return buffer.String()
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGUserInfo
/////////////////////////////////////////////////////////////////

// GetName returns the user name
func (obj *UserInfoImpl) GetName() string {
	return obj.userName
}

// GetSystemId returns the system ID for this user
func (obj *UserInfoImpl) GetSystemId() int {
	return obj.userId
}

// GetType returns the user type
func (obj *UserInfoImpl) GetType() byte {
	return obj.userType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, obj.userId, obj.userType, obj.userName)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning UserInfoImpl:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

func (obj *UserInfoImpl) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &obj.userId, &obj.userType, &obj.userName)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning UserInfoImpl:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
