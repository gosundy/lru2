package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

type Integer int

func init() {
	rand.Seed(123456)
}
func TestLruAddAndGet0(t *testing.T) {
	lru := NewLru(1)
	datas := []Integer{1, 2, 3, 4, 5}
	for _, data := range datas {
		lru.Add(data)
	}
	lruNode, err := lru.Get(Integer(5))
	if err != nil {
		t.Fatalf("expect: value 5, actual: with error %+v", err)
	}
	if lruNode.(Integer) != 5 {
		t.Fatalf("expect: value 5, actual: %v", lruNode)
	}
}
func TestLruAddAndGet1(t *testing.T) {
	lru := NewLru(5)
	datas := []Integer{1, 2, 3, 4, 5}
	for _, data := range datas {
		lru.Add(data)
	}
	lruNode, err := lru.Get(Integer(1))
	if err != nil {
		t.Fatalf("expect: value 1, actual: with error %+v", err)
	}
	if lruNode.(Integer) != 1 {
		t.Fatalf("expect: value 1, actual: %v", lruNode)
	}
}

func TestLruAddGet2(t *testing.T) {
	lru := NewLru(5)
	datas := rand.Perm(100000)
	datas2 := datas[100000-5:]
	for _, data := range datas {
		lru.Add(Integer(data))
	}
	for _, data := range datas2 {
		_, err := lru.Get(Integer(data))
		if errors.Is(err, LruNotError{}) {
			t.Fatalf("expect: %d, actual: not found", data)
		}
	}

}
func (i Integer) Encode() string {
	return fmt.Sprintf("%d", i)
}
