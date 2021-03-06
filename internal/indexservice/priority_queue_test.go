// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

package indexservice

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

const QueueLen = 10

func newPriorityQueue() *PriorityQueue {
	ret := &PriorityQueue{}
	for i := 0; i < QueueLen; i++ {
		item := &PQItem{
			value:    nil,
			key:      UniqueID(i),
			priority: i,
			index:    i,
		}
		ret.items = append(ret.items, item)
	}
	heap.Init(ret)
	return ret
}

func TestPriorityQueue_Len(t *testing.T) {
	pq := newPriorityQueue()

	assert.Equal(t, QueueLen, pq.Len())
	pq = nil
}

func TestPriorityQueue_Push(t *testing.T) {
	pq := newPriorityQueue()
	for i := 1; i <= QueueLen; i++ {
		item := &PQItem{
			key:      UniqueID(i),
			priority: i,
			index:    i,
		}
		pq.Push(item)
		assert.Equal(t, i+QueueLen, pq.Len())
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := newPriorityQueue()
	cnt := 0
	for i := 0; i < QueueLen; i++ {
		if i%2 == 0 {
			continue
		}
		pq.Remove(UniqueID(i))
		cnt++
	}
	assert.Equal(t, QueueLen-cnt, pq.Len())
}

func TestPriorityQueue_UpdatePriority(t *testing.T) {
	pq := newPriorityQueue()
	key := UniqueID(pq.Len() / 2)
	pq.UpdatePriority(key, -pq.Len())
	item := pq.Peek()
	assert.Equal(t, key, item.(*PQItem).key)
}

func TestPriorityQueue_IncPriority(t *testing.T) {
	pq := newPriorityQueue()
	key := UniqueID(pq.Len() / 2)
	pq.IncPriority(key, -pq.Len())
	item := pq.Peek()
	assert.Equal(t, key, item.(*PQItem).key)
}
