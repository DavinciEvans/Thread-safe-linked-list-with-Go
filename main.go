package main

import (
	"sync"
)

var (
	wg       sync.WaitGroup
	linkList LinkList
	m        sync.Mutex
)

func main() {
	//InstanceQuickSort()
	InstanceLRU()
}
