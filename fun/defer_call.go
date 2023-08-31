package fun

import (
	"fmt"
	"io"
	"log"
	"os"
)

func deferCall() {
	f, err := os.Open("./main.go")
	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		fmt.Println("haha")
		err := f.Close()
		if err != nil {

		}
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
}

func deferSort() {
	defer println(1)
	defer println(2)
	defer println(3)

	fmt.Println("main")
}

func deferPanic() {
	var f = func() {
		fmt.Println("defer")
	}

	var f1 func()

	defer f()
	// nil 会失败, 空指针
	defer f1()
}

func deferCopy() {
	x := 100
	defer fmt.Println("defer", x)

	x++
	fmt.Println("normal", x)
}

func bibao1() (z int) {
	defer func() {
		z += 200
	}()

	return 100 // return = 1. ret_val_z = 100 2.call defer 3. ret
}

func bibao2() int {
	z := 0

	defer func() {
		z += 200 // 本地变量，和返回值无关
	}()

	z = 100
	return z
}

func deferRetBibao() {
	fmt.Println(bibao1())
	fmt.Println(bibao2())
}

// 切记，defer是在函数结束的时候才被执行，不合理的使用方式会浪费更多资源和逻辑错误

func deferError() {
	for i := 0; i < 1000; i++ {
		f, err := os.Open("./main.go")
		if err != nil {
			continue
		}

		// 实际上在循环内部不会执行，只有到函数结束的时候才会执行，无端延长了1000个f的生命周期，平白消耗资源
		// 需要把for循环内部打开文件的动作封装为一个函数
		defer f.Close()
	}
}

func MainDefer() {
	deferSort()
	deferCall()
	// deferPanic()
	// defer 本质上是在汇编 ret前才调用，所以
	deferCopy()
	deferRetBibao()
}
