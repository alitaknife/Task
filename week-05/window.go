package main

import (
	"container/list"
	"log"
	"sync"
	"time"
)

const (
	SUCCEED int = 0
	FAILED  int = 1
)

type metrics struct {
	success int64
	fail    int64
}

type Window struct {
	sync.RWMutex
	buckets    int
	key        int64
	statistics map[int64]*metrics
	data       *list.List
}

func NewWindow(bucket int) *Window {
	sw := &Window{}
	sw.buckets = bucket
	sw.data = list.New()
	return sw
}

func (sw *Window) AddSuccess() {
	sw.increase(SUCCEED)
}

func (sw *Window) AddFail() {
	sw.increase(FAILED)
}

func (sw *Window) increase(t int) {
	sw.Lock()
	defer sw.Unlock()
	nowTime := time.Now().Unix()
	if _, ok := sw.statistics[nowTime]; !ok {
		sw.statistics = make(map[int64]*metrics)
		sw.statistics[nowTime] = &metrics{}
	}
	if sw.key == 0 {
		sw.key = nowTime
	}
	// 一秒一个 bucket
	if sw.key != nowTime {
		sw.data.PushBack(sw.statistics[nowTime])
		delete(sw.statistics, sw.key)
		sw.key = nowTime
		if sw.data.Len() > sw.buckets {
			for i := 0; i <= sw.data.Len()-sw.buckets; i++ {
				sw.data.Remove(sw.data.Front())
			}
		}
	}

	switch t {
	case SUCCEED:
		sw.statistics[nowTime].success++
	case FAILED:
		sw.statistics[nowTime].fail++
	default:
		log.Fatal("error")
	}
}

// Len 获取数据长度
func (sw *Window) Len() int {
	return sw.data.Len()
}

// Data 获取数据
func (sw *Window) Data(space int) []*metrics {
	sw.RLock()
	defer sw.RUnlock()
	var data []*metrics
	var num = 0
	var m = &metrics{}
	for i := sw.data.Front(); i != nil; i = i.Next() {
		one := i.Value.(*metrics)
		m.success += one.success
		m.fail += one.fail
		if num%space == 0 {
			data = append(data, m)
			m = &metrics{}
		}
		num++
	}
	return data
}
