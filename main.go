package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	wg             sync.WaitGroup
	linkList       LinkList
	innerGroupList []sync.WaitGroup
	//m sync.Mutex
)

func main() {
	linkList.InitList()
	// 从文件读入数组
	f, _ := os.Open("numbers_list.txt")
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

	Monitor(1, &linkList, 2, 33)
}

/*
Monitor: 控制函数
*/
func Monitor(id int, linkList *LinkList, step int, pivot int) {
	fmt.Println("id: ", id, "is working.")
	var innerGroup sync.WaitGroup
	innerGroupList = append(innerGroupList, innerGroup)
	if linkList.Length <= 1 {
		innerGroup.Done()
		return
	}
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

	// 寻找当前 pivot 位置
	item := linkList.Head
	s := 0
	for item.Data != pivot {
		item = item.Suc
		s++
	}
	// 查看有几个 pivot 相同的值
	e := 1
	for item.Suc.Data == item.Data {
		item = item.Suc
		e++
	}
	// 切割原链表为 3 段
	aCopy, _ := linkList.Clone(0, s)
	bCopy, _ := linkList.Clone(s, e)
	cCopy, _ := linkList.Clone(e, linkList.Length)
	aCopy.Display(false)
	cCopy.Display(false)

	innerGroup.Add(3)
	go Monitor(id+1, &aCopy, step, aCopy.Head.Data)
	go Monitor(id+3, &cCopy, step, cCopy.Head.Data)
	innerGroup.Done()
	res := ConnectLinkList(aCopy, bCopy)
	res = ConnectLinkList(res, cCopy)
	res.Display(false)
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

		//m.Lock()
		//fmt.Println("ID: ", id)
		//l.Display(false)
		//m.Unlock()
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
