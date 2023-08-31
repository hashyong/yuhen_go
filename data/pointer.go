package data

import (
	"fmt"
	"unsafe"
)

// Pointer
/*
不能将内存地址与指针混为一谈

地址是内存中每个字节单元的唯一编号，而指针则是实体，指针需要分配内存空间，为专门用来保存地址的整形变量

        x:int             p:*int
-------+-----------------+----------+---------
 ...   |100              |0xc000000 |  ...      memory
-------+-----------------+----------+---------
        0xc000000         0xc000008             address
*/
func Pointer() {
	begin()
	operator()
}

// 常见操作
/*
取址运算符 & 用于获取目标地址
指针运算符 * 间接引用目标对象 解引用
二级指针 **T 含报名则携程 **package.T

指针默认值 nil，支持 == != 操作
不支持指针运算，可以借助unsafe变相实现
*/
func operator() {
	//	空指针也会分配内存
	var p *int
	fmt.Println(unsafe.Sizeof(p)) // 8
	/*
		 p: *int
		+---------+
		| 0       |
		+---------+
		 0xc000008
	*/

	var x int
	/*
		 p: *int
		+---------+
		| 0       |
		+---------+
		 0xc000008


		 x: int
		+---------+
		| 0       |
		+---------+
		 0xc000000
	*/

	// 二级指针，比较好整 这里就不画图了
	var pp **int = &p
	*pp = &x

	*p = 100
	**pp += 1
	fmt.Println(**pp, *p, x) // 101 101 101

	// 并非所有对象都可以进行取址操作
	m := map[string]int{
		"a": 1,
	}
	fmt.Println(m)

	//_ = &m["a"]
	// invalid operation: cannot take address of m["a"] (map index expression of type int)

	// 支持相等运算符，但不能做加减法运算，不能做类型转换
	// 如果两个地址指向同一地址 或都为nil，则相等
	var b byte
	var p1 *byte = &b
	fmt.Println(b, p1)
	//p1++            // invalid operation: p1++ (non-numeric type *byte)
	//n := (*int)(p1) // cannot convert p1 (variable of type *byte) to type *int

	var p2, p3 *int
	fmt.Println(p2 == p3) // true

	var x1 int
	p2, p3 = &x1, &x1
	fmt.Println(p2 == p3) // true

	var y int
	p3 = &y
	fmt.Println(p2 == p3) // false

	// 指针没有专门指向成员的 ->运算符，统一使用 . 表达式
	a := struct {
		x int
	}{100}

	p4 := &a
	p4.x += 100
	fmt.Println(p4.x)

	// 0 byte 对象指针是否相等，与版本以及编译优化有关，不等于nil

	var a1, a2 struct{}
	var a3 [0]int

	pa1, pa2, pa3 := &a1, &a2, &a3

	println(pa1, pa2, pa3)
	fmt.Println(pa1 == nil || pa2 == nil)
	fmt.Println(pa1 == pa2)
	/*
		0xc0000c9d18 0xc0000c9d18 0xc0000c9d18
		false
		false

		0xc000141d18 0xc000141d18 0xc000141d18
		false
		true
	*/

	// 借助unsafe 实现指针转换和运算，要保证内存安全
	//	普通指针：*T 包含类型信息
	//	通用指针：pointer 只有地址 没有类型
	//	指针整数，uintptr，足以存储地址的整数

	// 普通指针和通用指针 都能构成引用，影响垃圾回收，unintptr只是整数，不构成引用关系，无法组织gc
	d := [...]int{1, 2, 3}
	p5 := &d

	// *[3]int --》 *int
	p6 := (*int)(unsafe.Pointer(p5))

	// p6++
	p6 = (*int)(unsafe.Add(unsafe.Pointer(p6), unsafe.Sizeof(p5[0])))
	*p6 += 100
	fmt.Println(d) // [1 102 3], 这尼玛就离谱
}

// 简单开始下
func begin() {
	var x int
	var p *int = &x
	*p = 100
	fmt.Println(p, *p) // 0xc000192050 100
}
