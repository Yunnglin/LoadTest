package queue

import (
	"container/list"
	"fmt"
	"sync"
)

var lock sync.Mutex

type queue interface {
	Pop() interface{}
	Push(interface{})
	Dump()
}

type Queue struct {
	data *list.List
}

func (q Queue) Pop() interface{} {
	defer lock.Unlock()
	lock.Lock()
	iter := q.data.Back()
	if iter == nil {
		return nil
	}
	v := iter.Value
	q.data.Remove(iter)
	return v
}

func (q *Queue) Push(i interface{}) {
	defer lock.Unlock()
	lock.Lock()
	q.data.PushFront(i)
}

func (q *Queue) Dump() {
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item: ", iter.Value)
	}
}

func NewQueue() *Queue {
	return &Queue{list.New()}
}
