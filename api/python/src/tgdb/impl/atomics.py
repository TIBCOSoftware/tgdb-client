"""
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
 *  File name :attr.py
 *  Created on: 5/15/2019
 *  Created by: suresh
 *
 *		SVN Id: $Id: atomics.py 3256 2019-06-10 03:31:30Z ssubrama $
 *
 *  This file encapsulates basic Atomic functions needed
 """

from multiprocessing import *

"""
* typecode is one of the following
* 'c' = char
* 'u' = wchar 
* 'b' = byte
* 'B' = ubyte 
* 'h' = short     
* 'H' = ushort 
* 'i' = int       
* 'I' = uint 
* 'l' = long 
* 'L' = ulong 
* 'q' = longlong  
* 'Q' = ulonglong 
* 'f' = float     
* 'd' = double
"""


class AtomicReference(object):
    ref: Value = None

    def __init__(self, typecode, initial):
        self.ref = Value(typecode, initial, lock=True)

    def increment(self):
        with self.ref.get_lock():
            self.ref.value += 1
        return self.ref.value

    def decrement(self):
        with self.ref.get_lock():
            self.ref.value -= 1
        return self.ref.value

    def get(self):
        with self.ref.get_lock():
            return self.ref.value

    def set(self, v):
        with self.ref.get_lock():
            oldv = self.ref.value
            self.ref.value = v
            return oldv

    # @property
    # def value(self):
    #     with self.ref.get_lock():
    #         return self.ref
    #
    # @value.setter
    # def value(self, value):
    #     with self.ref.get_lock():
    #         self.ref = value





