package fun

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
因为闭包就是匿名函数+一堆环境变量的指针组成的结构体
所以 拿到手的闭包除非执行，否则是只有两个指针的结构体
这会导致延迟求值，这里就是指闭包执行的时候的环境变量的值是目前这个变量的值，不一定是我们预期的值
*/
func test() (s []func()) {
	for i := 0; i < 2; i++ {
		// 循环使用的i是始终为同一个变量
		// 闭包 存储i的指针 被添加到切片内
		s = append(s, func() {
			println(&i, i)
		})

	}
	return
}

func test1() (s1 []func()) {
	for i := 0; i < 2; i++ {
		x := i
		s1 = append(s1, func() {
			println(&x, x)
		})
	}
	return
}

// 静态局部变量，有点意思，类似于线程的私有存储一样，相当于函数的私有存储, 这不就相当于可以夹带私货了
// 局部变量定义在函数内部，将操作这个变量的方案以函数的方式提供出去, 可以 够骚
func call() func() {
	n := 0

	return func() {
		n++
		println("call", n)
	}
}

type Database struct{}

func (d Database) Get() string {
	return ""
}

// 模拟方法，绑定状态
func handle(db Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := db.Get()
		io.WriteString(w, url)
	}

}

// 包装函数，改变其签名或者增加欸外功能
// 修改签名
func partial(f func(int), x int) func() {
	return func() {
		f(x)
		fmt.Println("hahaha")
	}
}

// 增加功能
func proxy(f func()) func() {
	return func() {
		n := time.Now()
		defer func() {
			fmt.Println(time.Now().Sub(n))
		}()

		f()
	}
}

func Anonymous() {
	for _, f := range test() {
		// 因为test先执行，所以此时 i已经为2了
		// 此时拿到的f函数执行的时候，i已经是2了，啊这 有点意思，确实很容易写出bug
		f()
	}

	for _, f := range test1() {
		// 因为没有依赖i 这个环境变量，所以问题不大
		f()
	}

	f := call()
	f()
	f()
	f()

	db := Database{}
	http.HandleFunc("/url", handle(db))

	test := func(x int) {
		fmt.Println(x)
	}

	var f1 func() = partial(test, 100)
	f1()

	proxy(f1)()
}
