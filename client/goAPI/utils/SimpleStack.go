package utils

import (
	"errors"
	"sync"
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
 * File name: SimpleStack.go
 * Created on: Apr 14, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**
package main

import (
    "errors"
    "fmt"
    "sync"
)

// SimpleStack is a basic LIFO stack that re-sizes as needed.
type SimpleStack struct {
	items []interface{}
	sLock *sync.Mutex
}

// NewSimpleStack returns a new stack.
func NewSimpleStack() *SimpleStack {
	return &SimpleStack{
		items: make([]interface{}, 0),
		sLock:  &sync.Mutex{},
	}
}

// Push adds a node to the stack.
func (s *SimpleStack) Push(e interface{}) {
	s.sLock.Lock()
	defer s.sLock.Unlock()

	s.items = append(s.items, e)
}

// Pop removes and returns a node from the stack in last to first order.
func (s *SimpleStack) Pop() (interface{}, error) {
	s.sLock.Lock()
	defer s.sLock.Unlock()

	len := len(s.items)
	if len == 0 {
		return 0, errors.New("Empty Stack")
	}

	entry := s.items[len-1]
	s.items = s.items[:len-1]
	return entry, nil
}

func (s *SimpleStack) Items() []interface{} {
	return s.items
}

func (s *SimpleStack) Size() int {
	return len(s.items)
}

func main() {
	fmt.Println("Hello, playground")

	s := NewSimpleStack()
    	s.Push(1)
    	s.Push(2)
    	s.Push(3)
    	fmt.Println(s.Pop())
    	fmt.Println(s.Pop())
    	fmt.Println(s.Pop())

	fmt.Println("done")
}
 */

// SimpleStack is a basic LIFO stack that re-sizes as needed.
type SimpleStack struct {
	items []interface{}
	sLock *sync.Mutex
}

// NewSimpleStack returns a new stack.
func NewSimpleStack() *SimpleStack {
	return &SimpleStack{
		items: make([]interface{}, 0),
		sLock:  &sync.Mutex{},
	}
}

// Push adds a node to the stack.
func (s *SimpleStack) Push(e interface{}) {
	s.sLock.Lock()
	defer s.sLock.Unlock()

	s.items = append(s.items, e)
}

// Pop removes and returns a node from the stack in last to first order.
func (s *SimpleStack) Pop() (interface{}, error) {
	s.sLock.Lock()
	defer s.sLock.Unlock()

	len := len(s.items)
	if len == 0 {
		return 0, errors.New("Empty Stack")
	}

	entry := s.items[len-1]
	s.items = s.items[:len-1]
	return entry, nil
}

func (s *SimpleStack) Items() []interface{} {
	return s.items
}

func (s *SimpleStack) Size() int {
	return len(s.items)
}
