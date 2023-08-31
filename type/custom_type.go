package _type

import "fmt"

type A int   // 定义新类型
type B = int // 别名

func CustomType() {
	var a A
	var a1 int
	// fmt.Println(a == a1)
	//Invalid operation: a == a1 (mismatched types A and int)
	//即使底层类型相同，也并非同一类型
	//除了运算符外，不继承任何信息
	//	不能隐式转换。不能直接比较

	var a2 B
	fmt.Println(a1 == a2)
	//别名并没有改变变量的类型

	fmt.Println(a)

	// bool int等等都为有明确标识的类型
	// array slice map channel 类型和元素类型或者长度属性相关 称为未命名类型
	// 可以用type 提供具体名称，变为命名类型

	// []int // unnamed type
	type A []int // named type

	// [2]int, [3]int // 未命名类型，长度不同 不是同一类型
	// []int, []byte // 未命名类型，元素类型不同

	/*具有相同声明的未命名类型被看作同一类型
	- 相同基类型的指针 pointer.go
	- 相同的元素类型和长度的数组 array
	- 相同元素类型的切片 slice
	- 相同key类型的字典 map
	- 相同数据类型以及操作方向的通道 channel
	- 相同字段序列的结构体 struct
	- 相同签名的函数 function
	- 相同方法集合的接口 interface
	*/

	var f1, f2 interface {
		test()
	}

	fmt.Println(f1 == f2)

	/*转换规则
	- 所属类型相同
	- 基础类型相同，其中一个是未命名类型
	- 数据类型相同，双向通道赋值给单向通道，其中一个为未命名类型
	- 默认值nil赋值给切片 字典 通道 指针 函数 接口
	*/
}
