package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 测试
	rand.Seed(time.Now().UnixNano())
	sw := NewWindow(100)
	var r int
	for i := 0; i < 1000; i++ {
		r = rand.Intn(3)
		if r == 1 {
			sw.AddSuccess()
		} else {
			sw.AddFail()
		}
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
	}
	fmt.Println(sw.Len())
	for _, item := range sw.Data(3) {
		fmt.Println(item.success, item.fail)
	}
	fmt.Println("==============")
	for _, item := range sw.Data(5) {
		fmt.Println(item.success, item.fail)
	}
}
