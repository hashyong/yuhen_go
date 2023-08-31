package data

import (
	"fmt"
	"reflect"
	"strconv"
)

/* 接口 interface 是多个方法声明集合，代表一种调用契约

只要目标类型方法集合包含接口声明的全部方法，就视为实现该接口，无需显示声明。
单个目标类型可以实现多个接口，在设计上，接口解除了显示类型依赖，使用dip既 依赖倒置
提供面向对象的多态性。
应该定义小型 灵活以及组合接口（ISP）接口隔离。减少可视方法，屏蔽对象内部构造和实现细节

不能有字段
只能声明方法，不能实现
可以嵌入其他接口

通用以er作为后缀名称
空接口 interface{} any 没有任何方法声明

接口实现的依据是方法集 method set，所以要区分 T.set 和*T.set

*/

type Yer interface {
	test()
	toString() string
}

type Z2 struct{}

func (*Z2) test()           {}
func (Z2) toString() string { return "" }

func yerTest() {
	var z Z2

	// var t1 Yer = z
	//  Z2 does not implement Yer (method test has pointer receiver), T.set 不包含 *T.set
	var t Yer = &z
	t.toString()
	t.test()
}

// ---------------

// Z3 匿名接口可以直接用于变量定义，或者作为结构字段类型
type Z3 struct{}

func (Z3) test() {}

type node struct {
	// value 为字段名，具体的类型为interface，没有显示的if名字，是匿名接口
	value interface {
		test()
	}
}

func z3Test() {
	// 可以直接作为变量定义
	var t interface {
		test()
	} = Z3{}

	n := node{value: t}
	n.value.test()
}

//----

// 空接口可被赋值任何对象
type any1 = interface{}

func anyTest() {
	var i any1 = 123
	fmt.Println(i)

	i = "abc"
	fmt.Println(i)

	i = Z3{}
	fmt.Println(i)
}

// 接口会复制目标对象，通常以指针来代替原始值

type Y1er interface {
	toString() string
}

type Y1 struct {
	x int
}

func (y Y1) toString() string {
	return strconv.Itoa(y.x)
}

func y1erTest() {
	n := Y1{100}
	var y Y1er = n //copy
	n.x = 200
	fmt.Println(n.toString(), y.toString()) //200,100
}

/*
	匿名嵌入

类似匿名字段, 嵌入其他接口，目标类型 方法集中，必须包含嵌入接口在内的全部方法实现

- 相当于 include 声明，而非继承
- 不能嵌入自身或者循环嵌入
- 不允许声明重载overload
- 允许签名相同的声明（并集去重）

- 鼓励小接口的嵌入和组合
- 签名包括方法名 参数列表 数量，类型，排列顺序和返回值，不包括参数名
*/

type Aer interface {
	Test()
}

type Ber interface {
	ToString(string) string
}

type X1er interface {
	Aer
	Ber // 嵌入接口有相同声明

	ToString(s string) string // 签名相同（并集去重）

	// ToString() string // 不允许重载（签名不同）
	// duplicate method ToString

	Print()
}

type N4 struct{}

func (N4) ToString(string) string { return "" }

func (N4) Test()  {}
func (N4) Print() {}

func anonymousEntry() {
	i := X1er(N4{})
	t := reflect.TypeOf(i)

	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i))
	}
	/*
		{Print  func(data.N4) <func(data.N4) Value> 0}
		{Test  func(data.N4) <func(data.N4) Value> 1}
		{ToString  func(data.N4, string) string <func(data.N4, string) string Value> 2}
	*/
}

/*类型转换

超集接口 即使是非嵌入 可以隐式转换为子集，反之不行
- 超集包含子集的全部声明 含嵌入
- 和声明顺序无关
*/

type A1er interface {
	toString(string) string
}

type X2er interface {
	// A1er
	test()
	toString(string) string
}

type N5 struct{}

func (*N5) test() {}
func (N5) toString(s string) string {
	return ""
}

func interfaceConvert() {
	var x X2er = &N5{} // super
	var a A1er = x     // sub

	a.toString("xxx")

	//var x1 X2er = N5{} //  N5 does not implement X2er (method test has pointer receiver)
	//x1 := X2er(a) // A1er does not implement X2er (missing method test)

	// 类型推断将接口还原为原始类型，或者判断是否实现了某个更具体的接口类型
	// 原始类型
	n, ok := a.(*N5)
	fmt.Println(n, ok)

	// 接口
	x2, ok := a.(X2er)
	fmt.Println(x2, ok)

	/*
		&{} true
		&{} true
	*/

	// 可以使用 switch语句 在多种类型之间做出推断匹配，这样空接口有更多空间
	// 未使用变量会视为错误
	// 不支持 fallthrough
	var i any = &N5{}
	switch v := i.(type) {
	case nil:
	case *int:
	case func() string:
	case *N5:
		fmt.Println("N5", v)
	case X2er:
		fmt.Println("X2er", v)
	default:
	}

	//switch v1 := i.(type) { // v1 declared and not used
	//case *N5:
	//	fallthrough // cannot fallthrough in type switch
	//default:
	//}
	//
}

/*
接口比较
如果实现接口的动态类型支持，可以做相等或者不等于运算

允许方式:
- 接口类型相同
- 接口类型不同，但是声明完全相同
- 接口类型是超集或者子集
*/

type Qer interface {
	test()
}

type Wer interface {
	test()
}

type Eer interface {
	toString() string
}

type Rer interface {
	test()
	toString() string
}

type U struct {
	x int
}

func (U) test() {}
func (U) toString() string {
	return ""
}

func interfaceCmp() {
	fmt.Println(any(nil) == any(nil)) // true
	fmt.Println(any(100) == any(100)) // true

	a, b := U{100}, U{100}
	fmt.Println(a == b) // true

	// 相同类型的接口
	fmt.Println(Qer(a) == Qer(b))   // true
	fmt.Println(Qer(&a) == Qer(&b)) // false

	// 不同接口类型，但声明相同（或者超集）
	fmt.Println(Qer(a) == Wer(b)) // true
	fmt.Println(Qer(a) == Rer(b)) // true

	// 不同接口类型，声明不同
	//fmt.Println(Qer(a) == Eer(b)) // Qer(a) == Eer(b) (mismatched types Qer and Eer)

	// 对象和它实现的接口类型 b impl Aer
	fmt.Println(Qer(a) == b) // true

	// 接口内部由两个字段组成：itab和data
	// 只有两个字段都为nil，接口才等于nil，可以利用反射完善判断结果
}

/*
实现
接口本身是静态类型，内部使用itab结构存储 运行期所需要的类型信息

		type iface struct {
			tab  *itab   // 类型和方法表
			data unsafe.Pointer // 目标对象的指针
		}

		type itab struct {
		inter *interfacetype // 接口类型
		_type *_type			// 目标类型
		hash  uint32 // copy of _type.hash. Used for type switches.
		_     [4]byte
		fun   [1]uintptr // variable sized. 方法表
	}
*/
func ifaceImp() {
	var n L
	var r Ner = &n
	r.A()
}

// 编译器将接口和目标类型组合，生成实例 go.itab,*"".L, "".Ner
// 曲中有接口和目标类型引用，并在方法表中静态填入具体方法地址
// 每太看懂这段在说啥.
//TODO:需要研究研究

// 前文提及，接口有个重要特征. 赋值给接口时，会赋值目标对象

// 输出编译结果, 查看接口结构存储的具体内容
func ifaceCopy() {
	var l L = 100

	var t1 Mer = l // copy l
	t1.A()

	var t2 Ner = &l //copy ptr
	t2.B(1)
}

type Mer interface {
	A()
}

type Ner interface {
	Mer

	B(int)
	C(string) string
}

type L int

/*尝试展开下方法表的地址
fun
Mer.A() +0
Ner.Mer.A()+8
Ner.B()+16
Ner.C()+24
.D()   +32
*/

func (L) A()               {}
func (*L) B(int)           {}
func (*L) C(string) string { return "" }
func (*L) D()              {}

// 动态调用
// 以接口调用方法，需要通过方法表动态完成

//go:noinline
func testNew(n Ner) {
	n.B(9)
}

// 编译器尝试用内联这些方式优化，消除或者减少 因为动态调用和内存逃逸导致的性能问题
func testNew1(n Ner) {
	n.B(9)
}

func main() {
	var n L = 100
	var i Ner = &n
	testNew1(i)
}

/*
	go build -gcflags "-m" ./data 2>&1|grep interface|grep inlin

	data/interface.go:371:6: can inline testNew1
	data/interface.go:375:6: can inline main
	data/interface.go:379:10: inlining call to testNew1
*/

/*
go tool objdump -S  -s "data.Interface" ./yuhen

	var n L = 100
  0x4f664d		488d050cb41300		LEAQ 0x13b40c(IP), AX
  0x4f6654		e8278ef1ff		CALL runtime.newobject(SB);  heap alloc
  0x4f6659		48c70064000000		MOVQ $0x64, 0(AX)
	testNew(i)
  0x4f6660		4889c3			MOVQ AX, BX
  0x4f6663		488d05760e1e00		LEAQ go:itab.*yuhen/data.L,yuhen/data.Ner(SB), AX
  0x4f666a		e811ffffff		CALL yuhen/data.testNew(SB)



*/

/*
func testNew(n Ner) {
  0x4f6580              493b6610                CMPQ 0x10(R14), SP
  0x4f6584              7630                    JBE 0x4f65b6
  0x4f6586              4883ec18                SUBQ $0x18, SP
  0x4f658a              48896c2410              MOVQ BP, 0x10(SP)
  0x4f658f              488d6c2410              LEAQ 0x10(SP), BP
  0x4f6594              4889442420              MOVQ AX, 0x20(SP)
  0x4f6599              48895c2428              MOVQ BX, 0x28(SP)
        n.B(9)
  0x4f659e              488b4820                MOVQ 0x20(AX), CX  	; .tab+32->B
  0x4f65a2              4889d8                  MOVQ BX, AX        	; .data
  0x4f65a5              bb09000000              MOVL $0x9, BX		; arg
  0x4f65aa              ffd1                    CALL CX				; B.(.data, 0x9)
}
*/

// go build -gcflags "-m" ./data 2>&1|grep interface|grep heap
// data/interface.go:362:6: moved to heap: n
// 相比动态调用，内存逃逸才是接口导致的最大性能问题

type X int

type FuncString func() string

// 当前这个函数FuncString 实现了 String 方法
func (f FuncString) String() string {
	return f()
}

func skill() {
	// 让编译器检查，确保类型实现了指定接口
	//var _ fmt.Stringer = X(0) // X does not implement fmt.Stringer (missing method String)

	// 定义函数类型，包装函数，使其实现特定接口
	f := func() string {
		return "hello,world"
	}

	// f 实现了 string这个方法. 所以其可以转换为 fmt.Stringer
	// 感觉还是比较有用的，后边看看在哪儿可以应用
	var t fmt.Stringer = FuncString(f)
	fmt.Println(t)
}

func Interface() {
	yerTest()
	z3Test()
	anyTest()
	y1erTest()
	anonymousEntry()
	interfaceConvert()
	interfaceCmp()
	ifaceImp()
	ifaceCopy()
	main()
	skill()
}
