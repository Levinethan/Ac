package main

import (
	"fmt"
	"strings"
)

const (
	CAPS int=65536 //关键字为基准  26个字母 +10个数字  大小写62  支持中文65536
	MaxNodeNumber int = 5000000  //最大节点数量
)

type Node struct {
	strNo int  //
	fail  *Node   //失败节点
	next  [CAPS] *Node  //每一个元素都是一个指针
	id int  //编号
}


//树
type ACA struct {
	root *Node
	tail int  //记录尾部的数量
	nodeCount int //节点数量
	stringCount int //字符串长度
	stringList []string  //字符串列表
	linBreak bool  //是否还行
	nodeList  [MaxNodeNumber]*Node  //内部节点数量
}

func NewNode(id int) *Node  {
	n := new(Node)
	n.id = id
	return n
}

func NewACA() *ACA  {
	ac := new(ACA)

	ac.root = NewNode(ac.nodeCount)

	ac.nodeCount ++

	return ac
}

//抓取索引
func Getindex(char byte)int  {
	base := []byte("a")
	return int(char-base[0])
}

//截取字符串
func getString(index int)string  {
	base := []byte("a")
	target := base[0]+byte(index) //锁定目标
	var str[]byte  //字节
	str = append(str,target)  //叠加字符串

	return string(str)
}

func (ac *ACA)Insert(ortString string)  {
	str := strings.ToLower(ortString) //忽略大小写
	p := ac.root
	ac.stringCount ++ //数量叠加
	for i:=0 ; i <len(str); i++{
		index := Getindex(str[i])
		if p.next[index] == nil{
			p.next[index] = NewNode(ac.nodeCount)
			ac.nodeCount ++ //字母不存在 插入节点
		}
		p = p.next[index]
	}
	p.strNo = ac.stringCount
	ac.stringList = append(ac.stringList,ortString)  //叠加字符串
}

func (ac *ACA)PrintTree()  {
	r := ac.root
	fmt.Print("R(0)->")
	ac.RPrintTree(r,0)
}

func (ac *ACA)RPrintTree (cur *Node , depth int)  {
	for i := 0; i < CAPS ; i++{
		if cur.next[i]!=nil{
			//循环遍历每一个节点
			if ac.linBreak == true{
				for j:=0;j<depth;j++{
					//显示深度
					fmt.Printf("       ")
				}
				ac.linBreak = false
			}

			var failNodeid int  //失败节点编号
			if  cur.fail !=nil &&cur.next[i].fail != nil {
				failNodeid = cur.fail.id  //编号循环叠加
			}
			fmt.Printf("%s(%3d %3d [%3d])->",getString(i),cur.next[i].id,failNodeid,cur.next[i].strNo)
			temp := cur //备份当前节点
			cur = cur.next[i] //循环过程
			depth ++
			ac.RPrintTree(cur,depth) //递归下一级调用

			cur = temp

			ac.linBreak = true
		}

	}
	//最后这里处理意外的情况 没有字符串
	if cur.strNo>0{
		fmt.Printf("null \n")
	}
}

func (ac *ACA)BuildAC()  {
	head := 0    //头部

	r := ac.root  //备份根节点
	r.fail = nil  //处理失败节点

	ac.nodeList[head] = r //设置头部节点

	head++

	for {
		if head == ac.tail{
			break //进入失败节点 。终止循环
		}

		temp := ac.nodeList[ac.tail]
		ac.tail++
		var p*Node //节点
		for i:=0;i<CAPS;i++{
			if temp.next[i]!=nil{
				//每个节点不可以等于nil
				if temp == ac.root{
					temp.next[i].fail =ac.root
				}else {
					p = temp.fail  //继续处理失败节点
					for{
						if p == nil{
							break
						}

						if p.next[i]!=nil{
							temp.next[i].fail=p.next[i]
							break
						}
						p = p.fail //处理失败节点

					}
					if p == nil{
						temp.next[i].fail = r //跳到根目录
					}
				}
				ac.nodeList[head] = temp.next[i]
				head ++



			}
		}


	}

}
//AC自动机在内存里玩
func (ac *ACA)Query (findstr string) []string  {
	str := strings.ToLower(findstr) //忽略大小写

	n := len(str)

	var index int  //索引
	var ret []string //结果

	p := ac.root //备份根节点
	for i:=0; i<n;i++{


		index = Getindex(str[i])
		//获取索引
		for {
			//死循环
			if p.next[index]== nil && p!=ac.root{
				p = p.fail  //标注失败节点位置


			}else {
				break //跳出循环后  此时找到非失败节点
			}
		}
		p = p.next[index] //循环下一个级别
		if p == nil{
			p=ac.root  //如果p=nil  跳回根节点 继续找
		}
		temp := p //备份当前节点
		for {
			if temp == ac.root || temp.strNo==0{
				break
			}
			if temp.strNo > 0{
				//如果长度大于零 继续
				stringindex := temp.strNo -1 //获取索引
				ret = append(ret,ac.stringList[stringindex])  //叠加字符
			}

			temp = temp.fail
		}

	}




	return ret
}

func main()  {
	ac := NewACA()

	ac.Insert("sex")
	ac.Insert("she")
	ac.Insert("fuck")
	ac.Insert("coming")
	ac.Insert("法轮功")

	ac.BuildAC()
	ret := ac.Query("ooosex法轮功oohohohsexesessexfuckingcomingssdfjowieru")
	fmt.Println(ret)
	ac.PrintTree()
}

//词少的  可以每个线程挂一个AC自动机
//词量大的  可以加锁  保证线程安全
//ac自动机里用得最多的是   树套树 结构
