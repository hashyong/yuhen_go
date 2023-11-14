package data

/*泛型 是一种代码复用技术，也称作 模板 template
允许在强类型语言代码中，使用实例化时才指定的类型参数type parameter

- 函数和类型（包括接口）支持类型参数，方法暂不支持
- 支持推导，可以省略类型实参 type argument

- 通常以单个大写字母命名类型参数
- 类型参数必须有约束 constraints

*/

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func testGen1() {
	fmt.Println(Max[int](1, 2)) // 实例化，指定了类型实参
	fmt.Println(Max(1.1, 1.2))  // 类型推导，可以省略
	d := Data[int]{x: 1}
	d.test()
}

type Data[T any] struct {
	x T
}

func (d Data[T]) test() {
	fmt.Println(d)
}

//func (d Data[T]) test1[X any](x X) {} // method must have no type parameters
// 目前类型参数只有 函数和类型支持，方法不支持, 也不知道啥时候能支持

// 可以使用接口对类型参数进行约束. 指示它可以做什么， 如果为any，表示什么都可以做
// [T any] 还不如写为 [T], 约束直接写在[]里边会导致函数签名很丑, 不如有个where更优雅，看看后边会不会优化

type NGen struct{}

func (NGen) String() string { return "N" }

func testGen2[T fmt.Stringer](v T) { fmt.Println(v) }

// 类型集合
// 接口有两种定义
// 普通接口，为方法集合 methods sets, 指示可以做什么
// 类型约束，类型集合 type sets， 指示谁来做

// 相比普通接口 被动和隐式的实现，类型约束显示指定实现接口的类型集合
// |： 竖线，类型集合，匹配其中任一类型即可
// 波浪线：底层类型 时该类型的所有类型

/*
// Signed is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
type Ordered interface {
	Integer | Float | ~string
}

*/

type GenA int
type GenB string

func (a GenA) Test() {
	fmt.Println("A:", a)
}

func (b GenB) Test() {
	fmt.Println("B:", b)
}

type Tester interface {
	GenA | GenB

	Test() // 集合里边的类型需要实现
}

type TesterPointer interface {
	*GenA | *GenB

	Test() // 集合里边的类型需要实现
}

//go:noinline
func testTester[T Tester](x T) {
	x.Test()
}

//go:noinline
func testTesterPoniter[T TesterPointer](x T) {
	x.Test()
}

func testT1[T int | float32](x T) {
	fmt.Println(x)
}

func testMakeSlice[T int | float64](x T) []T {
	s := make([]T, 10)
	for i := 0; i < cap(s); i++ {
		s[i] = x
	}

	return s
}

func testT2[T int | float64, E ~[]T](x E) {
	for v := range x {
		fmt.Println(v)
	}
}

func testT3[T any](x T) {
	//switch x.(type) {
	//}
	//  cannot use type switch on type parameter value x (variable of type T constrained by any)

	// 要转换为普通接口
	switch any(x).(type) {
	case int:
		fmt.Println("int:", x)
	}
}

func testGen3() {
	// fmt.Println(Max(struct{}{}, struct{}{}))
	// struct{} does not satisfy constraints.Ordered (struct{} missing in ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string)

	// 普通接口可以用做约束，但类型约束 不能当作普通接口使用
	//var x constraints.Integer = 1
	//  cannot use type constraints.Integer outside a type constraint: interface contains type constraints

	// 类型约束除了类型集合外，还可以有方法声明
	testTester[GenA](1)     // GenA(1)
	testTester[GenB]("xxx") // GenB("xxx")

	// var c Tester = GenA(1)
	// 还是那个问题，类型约束不能当作 interface来使用
	// cannot use type Tester outside a type constraint: interface contains type constraints

	// 类型约束
	// 除了包含类型集合的接口类型意外，也可以直接写入参数列表

	/*
		[T any]
		[T int]
		[T ~int] // int或者任意以int作为底层类型的数组
		[T int | string] // int or string
		[T io.reader]
	*/
	testT1(1)
	testT1[float32](1.1)
	//testT1("xxx") // string does not satisfy int | float32 (string missing in int | float32)

	_ = testMakeSlice(1)
	testT2([]int{})

	// 泛型函数内部不自持定义新类型
	// 如果类型约束不是接口，则无法调用其成员

	// 类型转换
	// 不支持 switch 类型推断，要转换为普通接口
	testT3(1)
}

/*
泛型的实现方式，通常有俩
- 通过模板来实现：为每次调用生成代码实例。即便类型参数相同
- 通过字典来实现：单份代码实例，以字段传递类型参数信息

模板方式性能比较好，但是编译时间比较长，而且生成的文件比较大
字典方式代码比较少了但时复杂度高，而且性能比较差
*/

//go:noinline
func testGenImp[T any](x T) {
	println(x)

	// 操作步骤
	// 1. go build -gcflags "-l"
	// 2. go tool objdump -s "data\.testGenImp" ./yuhen
	/*	俩函数并没有使用字典
		TEXT yuhen/data.testGenImp[go.shape.string](SB) /mnt/c/Users/GolandProjects/yuhen/data/generic.go
		  generic.go:211        0x4fc180                493b6610                CMPQ 0x10(R14), SP
		  generic.go:211        0x4fc184                7642                    JBE 0x4fc1c8
		  generic.go:211        0x4fc186                4883ec18                SUBQ $0x18, SP
		  generic.go:211        0x4fc18a                48896c2410              MOVQ BP, 0x10(SP)
		  generic.go:211        0x4fc18f                488d6c2410              LEAQ 0x10(SP), BP
		  generic.go:211        0x4fc194                48894c2430              MOVQ CX, 0x30(SP)
		  generic.go:211        0x4fc199                48895c2428              MOVQ BX, 0x28(SP)
		  generic.go:212        0x4fc19e                6690                    NOPW
		  generic.go:212        0x4fc1a0                e85bbbf3ff              CALL runtime.printlock(SB)
		  generic.go:212        0x4fc1a5                488b442428              MOVQ 0x28(SP), AX
		  generic.go:212        0x4fc1aa                488b5c2430              MOVQ 0x30(SP), BX
		  generic.go:212        0x4fc1af                e84cc4f3ff              CALL runtime.printstring(SB)
		  generic.go:212        0x4fc1b4                e8a7bdf3ff              CALL runtime.printnl(SB)
		  generic.go:212        0x4fc1b9                e8c2bbf3ff              CALL runtime.printunlock(SB)
		  generic.go:213        0x4fc1be                488b6c2410              MOVQ 0x10(SP), BP
		  generic.go:213        0x4fc1c3                4883c418                ADDQ $0x18, SP
		  generic.go:213        0x4fc1c7                c3                      RET
		  generic.go:211        0x4fc1c8                4889442408              MOVQ AX, 0x8(SP)
		  generic.go:211        0x4fc1cd                48895c2410              MOVQ BX, 0x10(SP)
		  generic.go:211        0x4fc1d2                48894c2418              MOVQ CX, 0x18(SP)
		  generic.go:211        0x4fc1d7                e88492f6ff              CALL runtime.morestack_noctxt.abi0(SB)
		  generic.go:211        0x4fc1dc                488b442408              MOVQ 0x8(SP), AX
		  generic.go:211        0x4fc1e1                488b5c2410              MOVQ 0x10(SP), BX
		  generic.go:211        0x4fc1e6                488b4c2418              MOVQ 0x18(SP), CX
		  generic.go:211        0x4fc1eb                eb93                    JMP yuhen/data.testGenImp[go.shape.string](SB)

		TEXT yuhen/data.testGenImp[go.shape.int](SB) /mnt/c/Users/GolandProjects/yuhen/data/generic.go
		  generic.go:211        0x4fc200                493b6610                CMPQ 0x10(R14), SP
		  generic.go:211        0x4fc204                7636                    JBE 0x4fc23c
		  generic.go:211        0x4fc206                4883ec18                SUBQ $0x18, SP
		  generic.go:211        0x4fc20a                48896c2410              MOVQ BP, 0x10(SP)
		  generic.go:211        0x4fc20f                488d6c2410              LEAQ 0x10(SP), BP
		  generic.go:211        0x4fc214                48895c2408              MOVQ BX, 0x8(SP)
		  generic.go:212        0x4fc219                e8e2baf3ff              CALL runtime.printlock(SB)
		  generic.go:212        0x4fc21e                488b442408              MOVQ 0x8(SP), AX
		  generic.go:212        0x4fc223                e8d8c1f3ff              CALL runtime.printint(SB)
		  generic.go:212        0x4fc228                e833bdf3ff              CALL runtime.printnl(SB)
		  generic.go:212        0x4fc22d                e84ebbf3ff              CALL runtime.printunlock(SB)
		  generic.go:213        0x4fc232                488b6c2410              MOVQ 0x10(SP), BP
		  generic.go:213        0x4fc237                4883c418                ADDQ $0x18, SP
		  generic.go:213        0x4fc23b                c3                      RET
		  generic.go:211        0x4fc23c                4889442408              MOVQ AX, 0x8(SP)
		  generic.go:211        0x4fc241                48895c2410              MOVQ BX, 0x10(SP)
		  generic.go:211        0x4fc246                e81592f6ff              CALL runtime.morestack_noctxt.abi0(SB)
		  generic.go:211        0x4fc24b                488b442408              MOVQ 0x8(SP), AX
		  generic.go:211        0x4fc250                488b5c2410              MOVQ 0x10(SP), BX
		  generic.go:211        0x4fc255                eba9                    JMP yuhen/data.testGenImp[go.shape.int](SB)
	*/

	/*
		TEXT yuhen/data.testGenImp[go.shape.*uint8](SB) /mnt/c/Users/GolandProjects/yuhen/data/generic.go
		TEXT yuhen/data.testGenImp[go.shape.string](SB) /mnt/c/Users/GolandProjects/yuhen/data/generic.go
		TEXT yuhen/data.testGenImp[go.shape.int](SB) /mnt/c/Users/GolandProjects/yuhen/data/generic.go
	*/
}

/*
操作步骤
go build -gcflags "-l"
go tool objdump -S -s "data\.impGen" ./yuhen
*/

//go:noinline
func impGen() {
	testGenImp(1)
	/*
	  0x4f5cf4              488d05edf21d00          LEAQ yuhen/data..dict.testGenImp[int](SB), AX
	  0x4f5cfb              bb01000000              MOVL $0x1, BX
	  0x4f5d00              e8fb640000              CALL yuhen/data.testGenImp[go.shape.int](SB)
	*/

	// same underlying type
	type X int
	testGenImp(X(2))
	/*
	  0x4f5d05              488d05ecf21d00          LEAQ yuhen/data..dict.testGenImp[yuhen/data.X.2](SB), AX 	; 字典不同
	  0x4f5d0c              bb02000000              MOVL $0x2, BX
	  0x4f5d11              e8ea640000              CALL yuhen/data.testGenImp[go.shape.int](SB)				; 函数相同
	*/

	testGenImp("abc")
	/*
	  0x4f5d16              488d05d3f21d00          LEAQ yuhen/data..dict.testGenImp[string](SB), AX	; 字典不同
	  0x4f5d1d              488d1d5e491700          LEAQ 0x17495e(IP), BX
	  0x4f5d24              b903000000              MOVL $0x3, CX
	  0x4f5d29              e852640000              CALL yuhen/data.testGenImp[go.shape.string](SB)		; 函数不同
	*/

	//所有指针（任意类型）同步，与指针目标类型不同组
	a := 1
	/*
		0x4f5d32              48c744242001000000      MOVQ $0x1, 0x20(SP)
	*/
	testGenImp(a)
	/*
	  0x4f5d3b              488d05b6f21d00          LEAQ yuhen/data..dict.testGenImp[int](SB), AX
	  0x4f5d42              bb01000000              MOVL $0x1, BX
	  0x4f5d47              e874650000              CALL yuhen/data.testGenImp[go.shape.int](SB)
	*/
	testGenImp(&a) // 与&int &float同组，与int float不同组
	/*
	  0x4f5d4c              488d059df21d00          LEAQ yuhen/data..dict.testGenImp[*int](SB), AX
	  0x4f5d53              488d5c2420              LEAQ 0x20(SP), BX
	  0x4f5d58              e883640000              CALL yuhen/data.testGenImp[go.shape.*uint8](SB); 和testGenImp(a)不同
	*/

	b := 1.2
	/*
	  0x4f5d5d              f20f100583f01d00        MOVSD_XMM $f64.3ff3333333333333(SB), X0
	  0x4f5d65              f20f11442418            MOVSD_XMM X0, 0x18(SP)
	*/
	testGenImp(&b)
	/*
	  0x4f5d6b              488d0576f21d00          LEAQ yuhen/data..dict.testGenImp[*float64](SB), AX
	  0x4f5d72              488d5c2418              LEAQ 0x18(SP), BX
	  0x4f5d77              e864640000              CALL yuhen/data.testGenImp[go.shape.*uint8](SB);和testGenImp(&a)相同
	*/

}

/*
泛型实现有性能问题，争论源于 所有类型的指针都属于统一 gcshape，共享同一份代码实例

编译器会插入一条指令，以字段装载类型信息，通过寄存器传递给目标代码
必要时以接口方式处理，可能引发 动态调用，无法内联，内存逃逸 等性能问题

基于前文内容，我们了解go的泛型时 gcshape，兼有静态模板和动态字典特点，这就导致其某些行为看起来是接口套壳
选用泛型的时候，本就该了解并接收与接口版本莱斯的性能损耗，而非将其当作泛型的固有缺陷

主要还是因为目前是半成品，也没有编译优化，要看看后续版本
*/

//go:noinline
func performanceGen() {
	var a GenA = 1
	var b GenB = "2"
	/*
	        var a GenA = 1
	  0x4f5d74              488d0565cd1300          LEAQ 0x13cd65(IP), AX
	  0x4f5d7b              0f1f440000              NOPL 0(AX)(AX*1)
	  0x4f5d80              e8fb96f1ff              CALL runtime.newobject(SB) ;跑到堆上去了
	  0x4f5d85              4889442418              MOVQ AX, 0x18(SP)
	  0x4f5d8a              48c70001000000          MOVQ $0x1, 0(AX)
	        var b GenB = "2"
	  0x4f5d91              488d05a8cd1300          LEAQ 0x13cda8(IP), AX
	  0x4f5d98              e8e396f1ff              CALL runtime.newobject(SB)
	  0x4f5d9d              4889442410              MOVQ AX, 0x10(SP)
	  0x4f5da2              48c7400801000000        MOVQ $0x1, 0x8(AX)
	  0x4f5daa              488d0df7ed1d00          LEAQ 0x1dedf7(IP), CX
	  0x4f5db1              488908                  MOVQ CX, 0(AX)
	*/

	// testTester(a)
	// testTester(b)

	testTesterPoniter(&a)
	testTesterPoniter(&b) // 确实，指针版本存在内存逃逸
	/*
			func performanceGen() {
			  0x4f5da0              493b6610                CMPQ 0x10(R14), SP
			  0x4f5da4              7641                    JBE 0x4f5de7
			  0x4f5da6              4883ec20                SUBQ $0x20, SP
			  0x4f5daa              48896c2418              MOVQ BP, 0x18(SP)
			  0x4f5daf              488d6c2418              LEAQ 0x18(SP), BP
			        testTester(a)
			  0x4f5db4              488d05e5111e00          LEAQ yuhen/data..dict.testTester[yuhen/data.GenA](SB), AX
			  0x4f5dbb              bb01000000              MOVL $0x1, BX
			  0x4f5dc0              e8db640000              CALL yuhen/data.testTester[go.shape.int](SB)
			        testTester(b)
			  0x4f5dc5              488d05e4111e00          LEAQ yuhen/data..dict.testTester[yuhen/data.GenB](SB), AX
			  0x4f5dcc              488d1dcdec1d00          LEAQ 0x1deccd(IP), BX
			  0x4f5dd3              b901000000              MOVL $0x1, CX
			  0x4f5dd8              e863640000              CALL yuhen/data.testTester[go.shape.string](SB)
			}
			  0x4f5ddd              488b6c2418              MOVQ 0x18(SP), BP
			  0x4f5de2              4883c420                ADDQ $0x20, SP
			  0x4f5de6              c3                      RET
			func performanceGen() {
			  0x4f5de7              e874f6f6ff              CALL runtime.morestack_noctxt.abi0(SB)
			  0x4f5dec              ebb2                    JMP yuhen/data.performanceGen(SB)


			  testTesterPoniter(&a)
		  0x4f5db4              488b5c2418              MOVQ 0x18(SP), BX
		  0x4f5db9              488d05e0121e00          LEAQ yuhen/data..dict.testTesterPoniter[*yuhen/data.GenA](SB), AX
		  0x4f5dc0              e8db640000              CALL yuhen/data.testTesterPoniter[go.shape.*yuhen/data.GenA](SB)
		        testTesterPoniter(&b)
		  0x4f5dc5              488d05e4121e00          LEAQ yuhen/data..dict.testTesterPoniter[*yuhen/data.GenB](SB), AX
		  0x4f5dcc              488b5c2410              MOVQ 0x10(SP), BX
		  0x4f5dd1              e86a640000              CALL yuhen/data.testTesterPoniter[go.shape.*yuhen/data.GenB](SB)
	*/

}

func Generic() {
	testGen1()
	//testGen2(1) // int does not satisfy fmt.Stringer (missing method String)
	// 成功约束，int方法 没有string方法
	testGen2(NGen{})
	testGen3()
	impGen()
	performanceGen()
}
