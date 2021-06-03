package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	HitCount   int
	VisitCount int
	countLock  sync.Mutex
)

type CacheStruct struct {
	Cache     LinkList
	MaxLength int
}

/*
LruCacheScheduling: LRU 调度函数
id：线程 id
linkList：缓冲区使用的链表
cap：缓冲区大小
*/
func LruCacheScheduling(id int, cacheStruct *CacheStruct, data int) {
	node := cacheStruct.Cache.Search(data)
	// 如果存在缓冲区内，则将其放到缓冲区顶部
	if node != nil {
		pos := cacheStruct.Cache.Find(node)
		_ = cacheStruct.Cache.Erase(pos)
		_ = cacheStruct.Cache.Insert(0, node)
		countLock.Lock()
		HitCount ++
		VisitCount++
	} else { // 不在缓冲区内，则将其放入缓冲区
		node := &Node{Pre: nil, Suc: nil, Data: data}
		if cacheStruct.Cache.Length == 0 {
			_ = cacheStruct.Cache.Add(node)
		} else {
			_ = cacheStruct.Cache.Insert(0, node)
		}
		// 超出缓冲区，删除最末尾的节点
		if cacheStruct.Cache.Length > cacheStruct.MaxLength {
			_ = cacheStruct.Cache.Erase(cacheStruct.Cache.Length - 1)
		}
		countLock.Lock()
		VisitCount++
	}
	countLock.Unlock()
}

/*
实例2：LRU最近最少使用调度
*/
func InstanceLRU() {
	// 从文件读入数组
	f, _ := os.Open("numbers_list2.txt")
	defer f.Close()
	r := bufio.NewReader(f)
	var numbersList []int
	for {
		s, err := r.ReadString('\n')
		// 结尾不再是回车时，读取结束
		if err != nil {
			break
		}
		if len(s) != 0 {
			num, err := strconv.Atoi(s[:len(s)-1])
			// 文件以空行表示结束
			if err != nil {
				break
			}
			numbersList = append(numbersList, num)
		}
	}
	// 创建基本链表
	linkList.InitList()

	// 初始化
	count := 5  // 线程数量
	wg.Add(count)  // 声明 WaitGroup
	cache := CacheStruct{linkList, 10}

	// 设置协程运行
	for i := 0; i < count; i++ {
		go func(id int) {
			for _, j := range numbersList {
				LruCacheScheduling(id, &cache, j)
				// 设置显示屏为临界区
				//m.Lock()
				// 打印
				var str string
				rw.RLock()
				rw.RUnlock()
				node := cache.Cache.Head
				for node != nil {
					str += strconv.Itoa(node.Data)
					str += "\t"
					//fmt.Printf("%d, ", node.Data)
					node = node.Suc
				}
				fmt.Printf("Thread ID: %d, Search Data: %d,\t %s\n", id, j, str)
				//cache.Cache.Display(false)
				//m.Unlock()
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	// 打印命中率
	fmt.Println("______ Result ______")
	fmt.Printf("Hit Rate: %.3f\n", float64(HitCount) / float64(VisitCount))
}
