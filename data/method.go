package data

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
)

// Method
/*
方法 method 是与对象实例 instance 相绑定的特殊函数

方法是面向对象编程的基本概念，用于维护和展示对象的自身状态。对象是内敛的。每个实例都有各自不同的独立特征。以属性和方法来对外暴露
普通函数专注于算法流程，接受参数完成漏极运算，返回结果并清理现场
也就是说。方法有持续性的状态，而函数通常是没有的

前置接收参数 receiver , 代表方法的所属类型
可为当前包内除接口和指针意外的任何类型定义方法
不支持静态方法或者关联函数
不支持重载
*/
func Method() {
	var a N2 = 25
	fmt.Println(a.toString())

	n3Test()
	anonymousEmb()
}

// -----
/*匿名嵌入
像访问匿名类型成员一样来使用，由编译器负责查找
*/
type data struct {
	sync.Mutex
	buf [1024]byte
}

// 同名遮蔽。利用其可实现类似 覆盖 override的操作
// 按照最小深度优先原则，如果两个同名方法深度相同，编译器无法选出来，需要显示指定

type E struct{}
type T struct{ E }

func (E) toString() string {
	return "E"
}

func (t T) toString() string {

	return "T." + t.E.toString()
}

// 这里不小心写出来死循环，引以为戒
//func (t T) toString() string {
//
//	return "T." + t.toString()
//}

// 同名，但是函数签名不同, 没法override

type T1 struct{ E }

func (T1) toString(s string) string {
	return "T1：" + s
}

func anonymousEmb() {
	d := data{}
	d.Lock() // sync.(*Mutex).lock()
	defer d.Unlock()

	var t T
	fmt.Println(t.toString(), t.E.toString()) // T.E E

	var t1 T1
	//fmt.Println(t1.toString())
	/*
		data/method.go:75:14: not enough arguments in call to t1.toString
		        have ()
		        want (string)
	*/
	// 选择深度最小的方法
	fmt.Println(t1.toString("sdfs"))
	// 明确目标
	fmt.Println(t1.E.toString())

	// 匿名类型的方法只能访问自己的字段，对外层一无所知

	set()
	setValue()
}

/*
方法值
和函数一样，方法除了直接调用外，还可以复制给变量，或者作为参数传递
按照引用方式不同，分为表达式 expression和值 value两种
*/

type V int

func (v V) copy() {
	println(v)
}
func (v *V) ref() {
	if v == nil {
		return
	}
	println(*v)
}

// 表达式 expr很好理解，将方法还原为普通的函数，显示传递接收参数 receiver
// 而方法值 value 似乎打包了接收参数和方法，导致签名不一样

func testF(f func()) {
	f()
}

func setValue() {
	var n V = 100
	// 表达式
	var e func(*V) = (*V).ref
	e(&n)

	// value
	var v func() = n.ref
	v()

	var n1 V = 200
	var v1 func() = n1.copy

	n1++
	n1.copy() // 201
	v1()      // 200
	testF(v1) // 200

	// 空指针要注意安全
	var p *V
	p.ref()         // value
	(*V)(nil).ref() // value
	(*V).ref(nil)   // exp
	// p.copy(), V.copy(*p)
	// runtime error: invalid memory address or nil pointer dereference

}

/*
类型有与之相关的方法集合 method set，这决定其是否实现了某个接口
根据接收参数 receiver的不同 分位 T和*T两种视角

	 T.set = T
	*T.set = T + *T
*/

type T2 int

func (T2) A()  {} //导出成员，否则反射无法获取
func (T2) B()  {} //导出成员，否则反射无法获取
func (*T2) C() {} //导出成员，否则反射无法获取
func (*T2) D() {} //导出成员，否则反射无法获取
func show(i interface{}) {
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i), t.Method(i).Name)
	}
}

// 直接方法调用，不涉及方法集，编译器自动转换所需参数 receiver
// 而转换 赋值 接口 interface 时，要检查方法及是否完全实现接口声明

type Xer interface {
	B()
	C()
}

// 别名扩展
// 通过类型别名，对方法集合进行分类，便于维护，或者新增别名 为类型添加扩展方法

type Z int

func (*Z) A() {
	fmt.Println("Z.A")
}

type Z1 = Z // 别名
func (*Z1) B() {
	fmt.Println("Z1.B")
}

// 注意，不同包的类型可以定义别名，但不能定义方法

type Y = bytes.Buffer

// func (*Y) test() {}
//  cannot define new methods on non-local type bytes.Buffer

func set() {
	var n T2 = 1
	var p *T2 = &n
	show(n) // A B
	show(p) // A B C D

	// ---
	// 方法调用，不涉及方法集合
	n.B()
	n.C()

	// 接口，检查方法集合
	// var x Xer = n
	// cannot use n (variable of type T2) as Xer value in variable declaration: T2 does not implement Xer (method C has pointer receiver)
	var x Xer = p
	x.B()
	x.C()
	/*
		    - 接口会复制对象，且复制品不能寻址
			- 如果T实现接口，通过接口调用时候，receiver可以被复制，不能获取指针
			-   *T实现接口，目标对象在接口以外，无论是及取值还是复制指针都没有问题
			- 这就是方法集与接口有关且 T = T, *T = T + *T


			除了直属方法外，列表里边还包括匿名类型 （E）的方法
			- T{E}  = T + E
			- T{*E} = T + E +*E
			- *T{E|*E} = T + *T + E + *E
	*/

	var z1 Z1
	z1.A()
	z1.B()

	var z Z
	z.A()
	z.B()
	show(&z)
	/*
		{A  func(*data.Z) <func(*data.Z) Value> 0} A
		{B  func(*data.Z) <func(*data.Z) Value> 1} B
	*/
}

//func (int) test() {}
//cannot define new methods on non-local type int

type N1 *int

//func (N1) test() {}
//invalid receiver type N1 (pointer or interface type)
// 类型为指针，不能整

type M int

func (M) test() {}

//func (*M) test() {}
//method M.test already declared at data/method.go:31:10

// 对接受参数命名无限制，按照管理选用简单有意义的名称
// 如方法内部不应用实例，可省略接收参数名，仅留类型

type N2 int

func (n N2) toString() string {
	return fmt.Sprintf("%#x", n)
}

func (N2) test() { // 省略
	fmt.Println("test")
}

// 接收参数可以是指针类型，调用时 根据此决定是否要复制 pass by value
/*
不能为指针和接口定义方法，是说类型 N本身不能是接口和指针
这与作为参数列表成员的 receiver *N意思完全不同

方法本质上就是特殊函数，接收参数无非就是其第一参数，只不过，在某些语言里边它是隐式的this
*/

type N3 int

func (n N3) copy() {
	fmt.Printf("%p, %v\n", &n, n)
}

func (n *N3) ref() {
	fmt.Printf("%p, %v\n", n, *n)
}

func n3Test() {
	var a N3 = 25
	fmt.Printf("%p\n", &a) // 0xc0000a6070
	a.copy()               // 0xc0000a6078, 25
	N3.copy(a)             // 0xc0000a6080, 25

	a++

	a.ref()       // 0xc0000a6070, 26
	(*N3).ref(&a) // 0xc0000a6070, 26

	var p *N3 = &a
	a.copy()
	a.ref() // (*N3).ref(&a)

	p.copy() //N3.copy(*p)
	p.ref()

	p2 := &p
	//p2.copy()
	//p2.ref()
	/*data/method.go:91:5: p2.copy undefined (type **N3 has no field or method copy)
	data/method.go:92:5: p2.ref undefined (type **N3 has no field or method ref)*/
	(*p2).copy()
	(*p2).ref()
}

/* 如何确定接收参数的类型

修改实例状态，用*T
不修改状态的小对象用固定值，用T
大对象用*T，减少复制成本
引用类型。字符串。函数等指针包装对象，用T
含mutex等同步手段， 用*T，避免因复制造成锁无效
其他无法确定的 都用*T

原则
1. 看是否要修改状态，默认使用*T
2. 假如为本就为指针类型。用T
3. 其次假如大小较小，用T
4. 其他都用*T
*/
