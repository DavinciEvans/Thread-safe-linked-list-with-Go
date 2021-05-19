package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

// 生成随机数列
func TestGenerate(t *testing.T) {
	count := math.Pow(2, 5)
	var numbersList string
	rand.Seed(time.Now().Unix())
	fmt.Println("Generate a number list now...")
	for i := .0; i <= count-1; i++ {
		numbersList = numbersList + strconv.Itoa(rand.Intn(int(count*2))) + "\n"
	}

	err := ioutil.WriteFile("./numbers_list.txt", []byte(numbersList), 0644)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// 测试列表的各项功能
func TestLinkList(t *testing.T) {
	var linkList LinkList
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
	fmt.Println("----- Basic List -----")
	linkList.Display(false)

	testNode1 := Node{Data: 100}
	testNode2 := Node{Data: 101}
	testNode3 := Node{Data: 102}

	// 测试 Insert
	fmt.Println("----- Insert -----")
	_ = linkList.Insert(linkList.Length, &testNode1)
	_ = linkList.Insert(20, &testNode2)
	_ = linkList.Insert(0, &testNode3)
	linkList.Display(false)

	// 测试 Find
	fmt.Println("----- Find -----")
	firstRes := linkList.Find(&testNode1)
	secondRes := linkList.Find(&testNode2)
	thirdRes := linkList.Find(&testNode3)
	fmt.Println("First: ", firstRes)
	fmt.Println("Second: ", secondRes)
	fmt.Println("Third: ", thirdRes)

	// 测试 Erase
	fmt.Println("----- Erase -----")
	_ = linkList.Erase(linkList.Length - 1)
	_ = linkList.Erase(20)
	_ = linkList.Erase(0)
	linkList.Display(false)

	// 测试 Get
	fmt.Println("----- Get -----")
	var item *Node
	item, _ = linkList.Get(0)
	fmt.Println("result: ", *item)

	// 测试 Clone
	fmt.Println("----- Clone -----")
	newList, _ := linkList.Clone(0, 10)
	newList.Display(false)
	linkList.Display(false)

	// 测试 ConnectLinkList
	fmt.Println("----- Connect Link List -----")
	var a, b LinkList
	connect := ConnectLinkList(newList, linkList)
	connect.Display(false)
	newList.Display(false)
	linkList.Display(false)
	connect = ConnectLinkList(a, connect)
	connect = ConnectLinkList(connect, b)
}
