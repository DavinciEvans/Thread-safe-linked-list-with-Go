package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Node struct {
	Pre  *Node
	Suc  *Node
	Data int
}

type LinkList struct {
	Head   *Node
	Tail   *Node
	Length int
}

var rw sync.RWMutex

// InitList: 初始化链表
func (l *LinkList) InitList() {
	(*l).Head = nil
	(*l).Tail = nil
	(*l).Length = 0
}

/*
Add 在链表尾部添加
node: 需要添加的节点
*/
func (l *LinkList) Add(node *Node) error {
	// 如果 node 为空
	if node == nil {
		return errors.New("the parameter node can't be nil")
	}

	node.Pre = nil
	node.Suc = nil
	// 将新元素放入链表当中
	rw.Lock()
	defer rw.Unlock()
	if (*l).Length == 0 {
		(*l).Head = node
	} else {
		oldTail := (*l).Tail
		(*oldTail).Suc = node
		node.Pre = oldTail
	}

	// 调整尾部与链表数量
	(*l).Tail = node
	(*l).Length++
	return nil
}

/*
Insert: 在某一位置插入具体数值
position: 链表当中的位置
node: 需要插入的节点
*/
func (l *LinkList) Insert(position int, node *Node) error {
	// 节点不存在或插入点超出范围
	if node == nil {
		return errors.New("the parameter node can not be nil")
	} else if position > (*l).Length || position < 0 {
		return errors.New("the parameter position is too big")
	}

	node.Pre = nil
	node.Suc = nil
	// 插入数据
	rw.Lock()
	defer rw.Unlock()
	if (*l).Length == 0 { // 如果此时链表为0
		(*l).Head = node
		(*l).Tail = node
	} else if position == (*l).Length { // 在最末尾插入
		(*l).Tail.Suc = node
		node.Pre = (*l).Tail
		(*l).Tail = node
	} else if position == 0 { // 在第一个位置插入
		(*l).Head.Pre = node
		node.Suc = (*l).Head
		(*l).Head = node
	}  else { // 中间插入
		item := (*l).Head
		for i := 0; i < position-1; i++ {
			item = item.Suc
		}
		preItem := item.Pre
		node.Pre = preItem
		preItem.Suc = node
		node.Suc = item
		item.Pre = node
	}

	l.Length++
	return nil
}

/*
Erase: 删除具体位置的元素
position: 位置
*/
func (l *LinkList) Erase(position int) error {
	// 异常判断
	if position > (*l).Length || (*l).Length == 0 {
		return errors.New("the parameter position is too big")
	}

	rw.Lock()
	defer rw.Unlock()
	if position == 0 { // 删除头部
		(*l).Head = (*l).Head.Suc
		if (*l).Head != nil {
			(*l).Head.Pre = nil
		}
	} else if position == (*l).Length-1 { // 删除末尾
		(*l).Tail = (*l).Tail.Pre
		if (*l).Tail != nil {
			(*l).Tail.Suc = nil
		}
	} else {
		node := (*l).Head
		for i := 0; i < position; i++ {
			node = node.Suc
		}
		node.Pre.Suc = node.Suc
		node.Suc.Pre = node.Pre
		node.Pre = nil
		node.Suc = nil
	}

	(*l).Length--
	return nil
}

/*
Find: 查找
data: 所需要查找的节点地址
返回的节点地址和返回查找到的位置，如果未找到，返回 -1
*/
func (l *LinkList) Find(data *Node) int {
	rw.RLock()
	defer rw.RUnlock()
	node := (*l).Head
	p := 0
	for node != nil {
		if node == data {
			return p
		}
		node = node.Suc
		p++
	}
	return -1
}

/*
Get: 得到具体位置的节点
position: 查找的位置
返回节点地址
*/
func (l *LinkList) Get(position int) (*Node, error) {
	// 异常判断
	if position > (*l).Length || (*l).Length == 0 {
		return nil, errors.New("the parameter position is too big")
	}
	rw.RLock()
	defer rw.RUnlock()

	node := (*l).Head
	for i := 0; i < position; i++ {
		node = node.Suc
	}
	return node, nil
}

/*
Search: 找到数据匹配的节点
data: 所要找的数据
返回节点地址，未找到返回 -1
 */
func (l *LinkList) Search(data int) *Node {
	rw.RLock()
	defer rw.RUnlock()

	node := (*l).Head
	for i := 0; i < (*l).Length; i++ {
		if node.Data == data {
			return node
		}
		node = node.Suc
	}
	return nil
}

/*
Clone: 链表克隆
*/
func (l *LinkList) Clone(start int, end int) (LinkList, error) {
	ll := LinkList{}
	if end < start || start < 0 && end > l.Length {
		return ll, errors.New("parameter is abnormal")
	}

	for i := start; i < end-start; i++ {
		node, _ := l.Get(i)
		item := *node
		_ = ll.Add(&item)
	}

	return ll, nil
}

/*
ConnectLinkList: 链表连接
*/
func ConnectLinkList(a LinkList, b LinkList) LinkList {
	var res LinkList
	aCopy, _ := a.Clone(0, a.Length)
	bCopy, _ := b.Clone(0, b.Length)

	if aCopy.Tail == nil {
		return bCopy
	} else if bCopy.Head == nil {
		return aCopy
	} else if aCopy.Tail == nil && bCopy.Head == nil {
		return res
	}

	aCopy.Tail.Suc = bCopy.Head
	bCopy.Head.Pre = aCopy.Tail
	res = aCopy
	return res
}

// Display: 展示链表
func (l *LinkList) Display(reverse bool) {
	//fmt.Printf("Link List: ")
	var str string
	rw.RLock()
	rw.RUnlock()
	if !reverse {
		node := (*l).Head
		for node != nil {
			str += strconv.Itoa(node.Data)
			str += " "
			//fmt.Printf("%d, ", node.Data)
			node = node.Suc
		}

	} else {
		node := (*l).Tail
		for node != nil {
			str += strconv.Itoa(node.Data)
			str += " "
			//fmt.Printf("%d, ", node.Data)
			node = node.Pre
		}
	}
	fmt.Println(str)
	//fmt.Printf("\n")
}
