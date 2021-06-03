package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
Monitor: 快排控制函数
*/
func Monitor(id int, linkList *LinkList, step int, pivot int) {
	wg.Add(step)
	for i := 0; i < step; i++ {
		var rangeNode [][]*Node
		start, _ := linkList.Get(i * linkList.Length / step)
		end, _ := linkList.Get((i+1)*linkList.Length/step - 1)
		rangeNode = append(rangeNode, []*Node{start, end})
		go linkList.MultiQuickSort(i, pivot, start, end)
	}
	wg.Wait()
	linkList.Display(false)
}

/*
MultiQuickSort: 链表多线程快排
*/
func (l *LinkList) MultiQuickSort(id int, pivot int, start *Node, end *Node) {
	item := start
	for item != end {
		temp := item.Suc
		if item.Data > pivot { // 比标杆大，移到后面
			pos := l.Find(item)
			_ = l.Erase(pos)
			//endPos := linkList.Find(end)
			//_ = linkList.Insert(endPos, item)
			_ = l.Add(item)
		} else if item.Data < pivot { // 比标杆大，移到前面
			pos := l.Find(item)
			_ = l.Erase(pos)
			_ = l.Insert(0, item)
		}
		item = temp

		m.Lock()
		fmt.Println("ID: ", id)
		l.Display(false)
		m.Unlock()
	}

	// 对于临界节点，还要再判断一次
	if item.Data > pivot { // 比标杆大，移到后面
		pos := l.Find(item)
		_ = l.Erase(pos)
		_ = l.Add(item)
	} else if item.Data < pivot { // 比标杆大，移到前面
		pos := l.Find(item)
		_ = l.Erase(pos)
		_ = l.Insert(0, item)
	}

	wg.Done()
}

/*
实例1：快排
*/
func InstanceQuickSort() {
	linkList.InitList()
	// 从文件读入数组
	f, _ := os.Open("numbers_list1.txt")
	defer f.Close()
	r := bufio.NewReader(f)

	// 创建基本链表
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
			node := Node{
				Pre:  nil,
				Suc:  nil,
				Data: num,
			}
			err = linkList.Add(&node)
			if err != nil {
				panic(err)
			}
		}
	}
	linkList.Display(false)
	fmt.Println("----- Start Sort -----")

	Monitor(1, &linkList, 4, 33)
}
