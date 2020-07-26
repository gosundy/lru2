package main

import "sync"

type Lru struct {
	capacity  int
	len       int
	head      *Link
	tail      *Link
	hashTable map[string]*Link
	rwLock    sync.RWMutex
}
type LruNode interface {
	Encode() string
}
type Link struct {
	data LruNode
	pre  *Link
	next *Link
}

type LruNotError struct {
}

func NewLru(capacity int) *Lru {
	return &Lru{capacity: capacity, hashTable: make(map[string]*Link), head: &Link{}}
}
func (lru *Lru) Add(data LruNode) {
	lru.rwLock.Lock()
	defer lru.rwLock.Unlock()

	//if link already in link list, just adjust current link as first link
	if link, ok := lru.hashTable[data.Encode()]; ok {
		//if link is first link, just return
		if link.pre == lru.head {
			return
		} else {
			//if link is latest, adjust tail as link's pre link
			if lru.tail == link {
				lru.tail = link.pre
				lru.tail.next = nil
			} else {
				//1.remove link
				link.pre.next = link.next
				link.next.pre = link.pre
			}

			//2.adjust as first link
			link.next = lru.head.next
			lru.head.next = link
			link.next.pre = link
			link.pre = lru.head
		}
	} else {
		//淘汰最后一个link
		if lru.len >= lru.capacity {
			delete(lru.hashTable, lru.tail.data.Encode())
			lru.tail = lru.tail.pre
			lru.tail.next = nil
			lru.len = lru.len - 1

		}

		link := &Link{data: data}
		link.next = lru.head.next
		lru.head.next = link
		link.pre = lru.head
		if lru.len != 0 {
			link.next.pre = link
		}
		if lru.tail == nil || lru.tail == lru.head {
			lru.tail = link
		}
		lru.len = lru.len + 1
		lru.hashTable[link.data.Encode()] = link
	}
	return
}
func (lru *Lru) Get(node LruNode) (LruNode, error) {
	if link, ok := lru.hashTable[node.Encode()]; ok {
		//if link is first, just return
		if link == lru.head.next {
			return link.data, nil
		}
		//if link is tail
		if link == lru.tail {
			lru.tail = link.pre
			lru.tail.next = nil
		} else {
			link.pre.next = link.next
			link.next.pre = link.pre
		}
		link.next = lru.head.next
		lru.head.next = link
		link.pre = lru.head
		link.next.pre = link
		return link.data, nil

	} else {
		return nil, LruNotError{}
	}
}

func (err LruNotError) Error() string {
	return "lru is empty"
}
func (err LruNotError) Is(target error) bool {
	_, ok := target.(LruNotError)
	return ok
}
