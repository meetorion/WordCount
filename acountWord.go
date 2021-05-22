package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type WordCount map[string]int

//实现三个方法以支持排序：Swqp(i,j int),Len() int, Less(i, j int) bool
type Pair struct {
	Key string
	Value int
}

type PairList []Pair

func (p PairList) Swap(i, j int) { //Swap方法必须用指针类型接收器吗,不必
	(p)[i], (p)[j] = (p)[j], (p)[i]
}
func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Key < p[j].Key
	}
	return p[i].Value > p[j].Value
}

func (w WordCount) Count(files []string) { //统计文件列表files的词频
	result := make(chan Pair, len(files)) //指定缓冲区大小
	done := make(chan bool, len(files))
	for _, filename := range files {
		go func(done chan bool, result chan Pair, filename string) {
			tmp := make(WordCount)
			tmp.Add(filename)
			for key, value := range tmp {
				pair := Pair{key, value}
				result <- pair
			}
			done <- true
		}(done, result, filename)
	}


	//监听通道
	for sum := len(files); sum > 0; {
		if sum == 0 {
			break
		}
		select { //因为监听规则是如果多个条件都满足，则随机选择其中一个满足的分支，如果选择了<-done,那么result通道中还有一些数据没有统计完
		case x := <-result:
			w[x.Key] += x.Value
		case <-done:
			sum--

		}

	}

	//处理result通道有可能剩下的数据
	for {
		flag := false
		select {
		case x := <-result:
			w[x.Key] += x.Value
		default:
			flag = true
			break //这个break只是跳出了select
		}
		if flag {
			break
		}
	}

	close(result)
	close(done)
}

func (w WordCount) Add(filename string) { //统计filename文件的词频,不区分大小写
	//读取文件内容
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	//统计词频
	now := ""
	for _,val := range str {
		if val >= 'a' && val <= 'z' || val >= 'A' && val <= 'Z' || val == '\''{
			now += string(val)
		} else {
			if now != "" {
				now = strings.ToLower(now) //不区分大小写
				w[now] += 1

			}
			now = ""
		}
	}
}

//WordCount内数据按照词频大小降序排列，并返回一个PairList类型数据（即结构体切片）
func (w WordCount) Sort() PairList {
	//将WordCount数据转换为PairList数据
	pairs := PairList{}
	for key, value := range w {
		pairs = append(pairs, Pair{key, value})
	}

	//排序
	sort.Sort(pairs)
	return pairs
	//转换回去，map是无序的，所以转换回去没有意义
}

func main() {
	w := make(WordCount)
	files := []string{"a.txt", "c.txt"} //要统计的文件名切片
	w.Count(files)
	p := w.Sort()
	fmt.Println(p)
}
