package util

import (
	"container/heap"
)

type QueueItem struct {
	Value    string
	Priority int64
	Index    int
	Data     interface{}
}

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Pop gives us the highest priority first
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*QueueItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) PushItem(item *QueueItem) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue) PopItem() *QueueItem {
	return heap.Pop(pq).(*QueueItem)
}

func (pq PriorityQueue) IsEmpty() bool {
	return pq.Len() <= 0
}

func InitQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}
