package utils

import "sync"

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
 * File name: SimpleQueue.go
 * Created on: Apr 14, 2019
 * Created by: achavan
 * SVN id: $id: $
 *
 */

/**
package main

import (
    "fmt"
    "sync"
)

// SimpleQueue is a basic FIFO queue based on a circular list that re-sizes as needed.
type SimpleQueue struct {
	qLock  *sync.Mutex
	values []interface{}
}

// NewSimpleQueue returns a new queue with the given initial size.
func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		qLock:  &sync.Mutex{},
		values: make([]interface{}, 0),
	}
}

// Enqueue adds an entry to the end of the queue.
func (q *SimpleQueue) Enqueue(x interface{}) {
	for {
		q.qLock.Lock()
		q.values = append(q.values, x)
		q.qLock.Unlock()
		return
	}
}

// Dequeue removes and returns an entry from the queue in first to last order.
func (q *SimpleQueue) Dequeue() interface{} {
	for {
		if len(q.values) > 0 {
			q.qLock.Lock()
			x := q.values[0]
			q.values = q.values[1:]
			q.qLock.Unlock()
			return x
		}
		break
		//return nil
	}
	return nil
}

func (q *SimpleQueue) Len() int {
	return len(q.values)
}

func main() {
	fmt.Println("Hello, playground")
	queue := NewSimpleQueue()
	queue.Enqueue(3)
	queue.Enqueue(4)
	queue.Enqueue(5)
	queue.Enqueue(6)
	queue.Enqueue(9)

	pick := queue.Dequeue()
	fmt.Println("Pick ", pick)
	fmt.Println("Remain ", queue.Items())

	pick = queue.Dequeue()
	fmt.Println("Pick ", pick)
	fmt.Println("Remain ", queue.Items())

	pick = queue.Dequeue()
	fmt.Println("Pick ", pick)
	fmt.Println("Remain ", queue.Items())

	pick = queue.Dequeue()
	fmt.Println("Pick ", pick)
	fmt.Println("Remain ", queue.Items())

	pick = queue.Dequeue()
	fmt.Println("Pick ", pick)
	fmt.Println("Remain ", queue.Items())

	fmt.Println("---In case nothing else left---")
	fmt.Println("Pick ", queue.Dequeue())

	fmt.Println("done")
}
*/

// SimpleQueue is a basic FIFO queue based on a circular list that re-sizes as needed.
type SimpleQueue struct {
	qLock  *sync.Mutex
	values []interface{}
}

// NewSimpleQueue returns a new queue with the given initial size.
func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		qLock:  &sync.Mutex{},
		values: make([]interface{}, 0),
	}
}

// Enqueue adds an entry to the end of the queue.
func (q *SimpleQueue) Enqueue(x interface{}) {
	for {
		q.qLock.Lock()
		q.values = append(q.values, x)
		q.qLock.Unlock()
		return
	}
}

// Dequeue removes and returns an entry from the queue in first to last order.
func (q *SimpleQueue) Dequeue() interface{} {
	for {
		if len(q.values) > 0 {
			q.qLock.Lock()
			x := q.values[0]
			q.values = q.values[1:]
			q.qLock.Unlock()
			return x
		}
		break
		//return nil
	}
	return nil
}

func (q *SimpleQueue) Items() []interface{} {
	return q.values
}

func (q *SimpleQueue) Len() int {
	return len(q.values)
}
