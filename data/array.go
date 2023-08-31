package data

import (
	"fmt"
	"unsafe"
)

/*
+--------------------+
|0 |1 |2 |3 |.....|99|  [100]int
+--------------------+

数组是单一内存块，并无其他附加结构

* 数据长度必须是 非负整形常量 表达式
* 长度是数组类型组成部分

* 初始化多维数据，仅第一维度允许用...
* 初始化符合类型，可省略元素类型标签

* 内置函数和len和cap都返回一维长度
* 元素类型支持 == != 那么数组也支持

* 支持数组指针直接访问
* 直接获取元素指针
* 值传递，会复制整个数组


*/

func Array() {
	// 编译期确定的大小
	var d1 [unsafe.Sizeof(0)]int
	var d2 [2]int
	//d2 = d1 // data/array.go:31:7: cannot use d1 (variable of type [8]int) as [2]int value in assignment
	fmt.Println(d1, d2)

	// 更灵活的初始化方式
	var a [4]int            // 元素会自动初始化为0
	b := [4]int{2, 5}       // 自动初始化未0
	c := [4]int{2, 3: 5}    // 支持指定位置init
	d := [...]int{2, 5}     // 自动推导元素个数
	e := [...]int{2, 11: 5} // 自动推导元素个数

	fmt.Println(a, b, c, d, e) // [0 0 0 0] [2 5 0 0] [2 0 0 5] [2 5] [2 0 0 0 0 0 0 0 0 0 0 5]

	type user struct {
		id   int
		name string
	}
	_ = [...]user{ //不能省略
		{1, "ss"}, //元素类型标签可以省略
		{2, "ps"},
	}

	//var x [2]int = {2, 5} // Syntax error: unexpected end of the statement, expecting ':=','=', or ','
	//x := [2]int{2,5} right

	//var y [...]int = [2]int{2, 5} //data/array.go:59:8: invalid use of [...] array (outside a composite literal)

	// 多维数组依旧是连续内存存储
	// 二维数组
	x := [...][2]int{
		{1, 2},
		{3, 4},
	}

	// 三维数组
	y := [...][3][2]int{
		{
			{1, 2},
			{1, 2},
			{1, 2},
		},
		{
			{3, 4},
			{3, 4},
			{3, 4},
		},
	}

	fmt.Println(x, len(x), cap(x))
	fmt.Println(y, len(y), cap(y))

	// 比较需要元素本身支持才可以

}
