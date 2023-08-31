package fun

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// 标志性错误：一种明确的状态，比如 io.EOF
// 提示性错误: 返回可读信息，提示错误原因

// 必须检查所有返回错误，这会导致代码不太美观
// 确实看起来有些丑
func errorTest() {
	file, err := os.Open("xxx")
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	c, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(c))
}

// EOF 标志性错误，通常使用 全局变量（指针或者接口）方式定义，提示性错误，一般会直接返回临时对象
// io
var EOF = errors.New("EOF")

/*package errors

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}*/

// 匹配分位几种类型,
// 是否有错 err != nil
// 具体错误 err == ErrVar
// 错误匹配  case ErrVal
// 类型转换 e, ok := err.(T)

// 自定义承载更多上下文的错误类型

type TestError struct {
	x int
}

func (e *TestError) Error() string {
	return fmt.Sprintf("test: %d", e.x)
}

var ErrZero = &TestError{0} // 指针，用来判断是否为同一个对象

func MyContextError() {
	var e error = ErrZero
	fmt.Println(e == ErrZero) // true

	if t, ok := e.(*TestError); ok {
		fmt.Println(t.x) // 0
	}

	err := testError()
	fmt.Println(err == nil)

	fmt.Println(cache())

	mainErr()
	mainPanic()
	mainPanic1()
	mainPanic2()
	mainPanic3()
}

// 函数返回的时候，要注意 error是否真等于 nil，以免错不成错
func testError() error {
	var err *TestError

	// 结构体指针，目前是空的
	fmt.Println(err == nil) // true

	// 接口只有类型和数据都为nil的是偶才等于nil
	// err显然是由类型信息的, 所以返回接口 err必然是不等于nil的
	// 假如需要返回空，那直接返回nil

	// 转换为接口
	return err
}

// 标准库
// errors.New: 创建包含文本信息的错误对象
// errors.Join: 打包一到多个错误对象，构建树状结构
// fmt.Errorf:  可以用w%包装一到多个错误对象

// 包装对象
// errors.Unwrap: 返回被包装错误对象 或者 列表
// errors.Is: 递归查找是否有指定的错误对象
// errors.As: 递归查找并获取 类型匹配的错误对象

func database() error {
	return errors.New("data")
}

func cache() error {
	if err := database(); err != nil {
		fmt.Println(errors.Join(database(), EOF))
		return fmt.Errorf("cache miss:%w", err)
	}

	return nil
}

// TestError1 从错误树获取信息
type TestError1 struct{}

func (t *TestError1) Error() string {
	return ""
}

func mainErr() {
	// 打包
	a := errors.New("a")
	b := fmt.Errorf("b, %w", a)
	c := fmt.Errorf("c, %w", b)

	fmt.Println(c)                     // c, b, a
	fmt.Println(errors.Unwrap(c) == b) // true

	// 递归检查
	fmt.Println(errors.Is(c, a)) // true

	x := &TestError1{}
	y := fmt.Errorf("y, %w", x)
	z := fmt.Errorf("z, %w", y)

	var x2 *TestError1
	if errors.As(z, &x2) {
		fmt.Println(x == x2) // true
	}
}

// panic/recover 的使用方式，和结构化异常很类似，panic类似 throw，recover 类似于 try catch
// panic 立即中断当前流程，执行 defer
// 在defer中，recover 捕获并返回panic数据

/*
- 会沿着调用堆栈向外传递，直至被捕获，或者进程崩溃
- 连续引发panic，仅最后一次可以被捕获

- 先捕获 panic，恢复执行，然后才返回数据
- recover之后，可以再次出现 panic，可以再次被捕获

- 无论recover与否，defer总会执行
- defer中引发panic，不影响后续延迟调用执行
*/

func mainPanic() {
	defer func() {
		// 拦截panic
		// 无法回到panic后的位置继续执行
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	func() {
		panic("p1") // 终止当前函数，直接执行defer
		fmt.Println("p1 after")
	}()

	println("exit")
}

func mainPanic1() {
	defer func() {
		//只有最后一次的panic会被捕获, 最终会输出P2
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	func() {
		defer func() {
			panic("p2") // 第二次
		}()

		panic("p1") // 第一次
	}()

	println("exit")
}

func mainPanic2() {
	fmt.Println("main2")
	defer func() {
		fmt.Println(recover()) // 捕获第二次 p2
	}()
	func() {
		// 看起来defer的执行顺序是个队列 LIFO
		defer func() {
			panic("p2")
		}()

		defer func() {
			fmt.Println(recover()) // 捕获第一次p1
		}()
		panic("p1")
	}()

	fmt.Println("exit")
}

// 由此可见上边的猜测没有问题，一个函数内部有多个defer的话是 串行执行的，执行顺序为 LIFO，后入先出. 本质就是个栈而已
func mainPanic3() {
	fmt.Println("main3")
	defer func() {
		fmt.Println(recover()) // 捕获第二次 p2
	}()

	func() {
		defer func() {
			fmt.Println(recover()) // 捕获第一次p1
		}()

		// 看起来defer的执行顺序是个队列，fifo
		defer func() {
			panic("p2")
		}()

		panic("p1")
	}()

	fmt.Println("exit")
}

// recover 只能在 defer函数内才可以正确执行，否则无法catch panic
func catch() {
	recover()
}

func mainPanic4() {
	defer catch()                // 生效的
	defer log.Println(recover()) // 无效，会作为参数立即执行
	defer recover()              // 无效
}

// 适用场景
// 在调用堆栈任意位置终端，跳转到合适位置
// runtime.Goexit: 不能用于 main goroutine
// os.Exit: 进程立即中止，defer不会被执行
// 无法恢复的现场，输出调用堆栈的现场
// 文件系统损坏
// 数据库无法连接等等
