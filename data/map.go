package data

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/* 字典 hashtable 存储kv对，使用频率极高
 *
- 引用类型 无序的键值对的集合
- 以初始化表达式胡总和 make 函数创建

- 主键要支持 ==  != 的类型
- 可以判断字典是否为nil 不支持比较操作

- 函数len反馈键值对的数量
- 访问不存在的主键 返回零值
- 迭代以随机次序返回

- 值不可寻址，需要整体赋值 没理解，这啥意思
- 按需扩张 但不会收缩内存(shrink)
*/

func p(m map[string]int) {
	fmt.Printf("%t, %x\n", m == nil, *(*uintptr)(unsafe.Pointer(&m)))
}

// 可变参数列表还是很方便的呀, 点个赞
func plist(ms ...map[string]int) {
	for _, v := range ms {
		p(v)
	}
}

// 字典 变量是指针类型
// 未初始化 未空指针
func map1() {
	var m1 map[string]int
	m2 := map[string]int{}
	m3 := make(map[string]int, 0)
	plist(m1, m2, m3)
	/*
		true, 0
		false, c000125c48
		false, c000125c78
	*/
}

func Map() {
	map1()
	map2()
	map3()
	map4()
	map5()
	map6()
	map7()
	//mainTest()
}

/*
安全
迭代期间，新增和删除操作都是安全的，但是无法控制次序
运行时候会对地点并发操作做出检测
- 启用竞争检测 data race 查找此类问题
- 使用 sync.map 来代替
*/
func map7() {
	m := make(map[int]int)
	for i := 0; i < 10; i++ {
		m[i] = i + 10
	}

	for k := range m {
		if k == 5 {
			m[100] = 1000
		}

		delete(m, k)
		fmt.Println(k, m)
	}
	/*
		9 map[0:10 1:11 2:12 3:13 4:14 5:15 6:16 7:17 8:18]
		0 map[1:11 2:12 3:13 4:14 5:15 6:16 7:17 8:18]
		1 map[2:12 3:13 4:14 5:15 6:16 7:17 8:18]
		2 map[3:13 4:14 5:15 6:16 7:17 8:18]
		5 map[3:13 4:14 6:16 7:17 8:18 100:1000]
		7 map[3:13 4:14 6:16 8:18 100:1000]
		3 map[4:14 6:16 8:18 100:1000]
		4 map[6:16 8:18 100:1000]
		6 map[8:18 100:1000]
		8 map[100:1000]
	*/

	/*	m1 := make(map[string]int)
		// write
		go func() {
			for {
				m1["a"] += 1
			}
		}()
		// read
		go func() {
			for {
				_ = m1["b"]
			}
		}()*/

	// select {}
	//	fatal error: concurrent map read and map write

	/*
		WARNING: DATA RACE
		Read at 0x00c000092480 by goroutine 10:
		  runtime.mapaccess1_faststr()
		      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
		  yuhen/data.map7.func2()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:103 +0x44

		Previous write at 0x00c000092480 by goroutine 9:
		  runtime.mapassign_faststr()
		      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
		  yuhen/data.map7.func1()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:97 +0x7b

		Goroutine 10 (running) created at:
		  yuhen/data.map7()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:101 +0x2c4
		  yuhen/data.Map()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:56 +0x29e
		  main.main()
		      /mnt/c/Users/ /GolandProjects/yuhen/main.go:37 +0x49d
		  yuhen/fun.mainPanic2()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:224 +0xa8
		  yuhen/fun.MyContextError()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:87 +0x218
		  yuhen/fun.mainPanic1()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:204 +0x44
		  yuhen/fun.MyContextError()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:86 +0x213
		  main.main()
		      /mnt/c/Users/ /GolandProjects/yuhen/main.go:33 +0x489

		Goroutine 9 (running) created at:
		  yuhen/data.map7()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:95 +0x256
		  yuhen/data.Map()
		      /mnt/c/Users/ /GolandProjects/yuhen/data/map.go:56 +0x29e
		  main.main()
		      /mnt/c/Users/ /GolandProjects/yuhen/main.go:37 +0x49d
		  yuhen/fun.mainPanic2()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:224 +0xa8
		  yuhen/fun.MyContextError()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:87 +0x218
		  yuhen/fun.mainPanic1()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:204 +0x44
		  yuhen/fun.MyContextError()
		      /mnt/c/Users/ /GolandProjects/yuhen/fun/error_opt.go:86 +0x213
		  main.main()
		      /mnt/c/Users/ /GolandProjects/yuhen/main.go:33 +0x489
	*/

}

/*
因 内联存储，键值迁移 以及 内存安全 需要，字典被设计为不可寻址
不能直接修改值成员 结构体或者数据 应整体赋值 或者以指针代替
*/
func map6() {
	type user struct {
		name string
		age  byte
	}

	// 指针
	m := map[int]*user{
		1: &user{"u1", 19},
	}

	// 返回指针来修改外部对象
	m[1].age = 20

	m1 := map[string]int{
		"a": 1,
	}

	m1["a"]++ // m[a]=m[a]+1
}

/*
迭代map 以随机次序返回键值
- 以struct可忽略值类型
- 随机效果依赖键值对数量和迭代次数
- 可用作简易版本随机挑选算法
*/
func map5() {
	m := make(map[int]struct{})

	for i := 0; i < 10; i++ {
		m[i] = struct{}{}
	}

	for i := 0; i < 10; i++ {
		for k := range m {
			fmt.Print(k, ",")
		}
		fmt.Println()
	}
	/* 基本都是随机返回的
	2,5,0,1,6,7,8,9,3,4,
	7,8,9,3,4,6,5,0,1,2,
	0,1,2,5,3,4,6,7,8,9,
	4,6,7,8,9,3,1,2,5,0,
	0,1,2,5,3,4,6,7,8,9,
	1,2,5,0,4,6,7,8,9,3,
	0,1,2,5,9,3,4,6,7,8,
	3,4,6,7,8,9,0,1,2,5,
	8,9,3,4,6,7,0,1,2,5,
	9,3,4,6,7,8,0,1,2,5,
	*/

}

// 打印的小技巧
func pmap(ms ...interface{}) {
	var size int
	for _, v := range ms {
		for {
			res, ok := v.(map[string]int)
			if ok {
				size = len(res)
				break
			}

			res1, ok := v.(map[string]struct {
				id   int
				name string
			})
			if ok {
				size = len(res1)
				break
			}

			break
		}
		fmt.Println(v, size)
	}
}

// 初始化的各种方式
func map2() {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
	}

	// 省略复合类型标签
	// 确实还是蛮省事的
	m2 := map[string]struct {
		id   int
		name string
	}{
		"a": {1, "u1"},
		"b": {2, "u2"},
	}

	// 分配足够的容量，可以减少后续扩张
	m3 := make(map[string]int, 10)
	m3["a"] = 1

	pmap(m1, m2, m3)
}

// map 之间不能比较，只能和nil比较
func map3() {
	m1 := map[int]int{}
	m2 := map[int]int{}
	fmt.Println(m1 == nil, m2 == nil)

	//fmt.Println(m1 == m2)
	// invalid operation: m1 == m2 (map can only be compared to nil)
}

/*
建议以 ok-idiom 访问主键 来确认是否存在
删除不存在的键值 不会引发错误
可对nil字段读和删除 但写会引发panic
*/
func map4() {
	m := map[string]int{}

	// a:0 是否存在
	x := m["a"]
	fmt.Println(x) // 因为可以对不存在的key来读 0

	x, ok := m["a"]
	fmt.Println(x, ok) // 0, false

	m["a"] = 0
	x, ok = m["a"]
	fmt.Println(x, ok) // 0, true

	var m1 map[string]int // nil
	_ = m1["a"]           // 读没有问题
	delete(m1, "a")       // 删除没有问题

	//m1["a"] = 0
	//写会有问题，assignment to entry in nil map
}

const max = 10000

//go:noinline
func test(cap int) map[int]int {
	m := make(map[int]int, cap)

	for i := 0; i < max; i++ {
		m[i] = i
	}

	return m
}

// 字典内敛存储有长度限制，如果超出 则重新分配内存并复制
const max1 = 1000000

//go:noinline
func test1[T any](v T) map[int]T {
	m := make(map[int]T, max1)
	for i := 0; i < max1; i++ {
		m[i] = v
	}
	return m
}

// 扩展的内存不会因键值删除而收缩, 必要的时候要新建字典
func mainTest() {
	//m := test1([128]byte{1, 2, 3})
	m := test1([128 + 1]byte{1, 2, 3})
	runtime.GC()
	for k := range m {
		delete(m, k)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		runtime.GC()
	}

	runtime.KeepAlive(&m)
}

// 内联版本 字段内存无法收缩 时间较短
// 说实话，每台看明白，等后边看了gc之后再回过来看看todo:123
/*
gc 1 @0.056s 0%: 0.010+0.38+0.040 ms clock, 0.12+0.058/0.095/0+0.48 ms cpu, 294->294->293 MB, 294 MB goal, 0 MB stacks, 0 MB globals, 12 P
gc 2 @0.394s 0%: 0.014+0.27+0.003 ms clock, 0.17+0/0.20/0+0.043 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 3 @1.405s 0%: 0.038+0.18+0.002 ms clock, 0.46+0/0.20/0+0.033 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 4 @2.406s 0%: 0.083+0.22+0.002 ms clock, 1.0+0/0.19/0+0.033 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 5 @3.407s 0%: 0.19+0.28+0.003 ms clock, 2.3+0/0.28/0+0.036 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 6 @4.408s 0%: 0.059+0.16+0.002 ms clock, 0.71+0/0.12/0+0.031 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 7 @5.408s 0%: 0.058+0.19+0.002 ms clock, 0.70+0/0.15/0+0.032 ms cpu, 293->293->293 MB, 587 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)


外联版本，外部对象被回收，时间较长
gc 1 @0.051s 1%: 0.010+2.5+0.048 ms clock, 0.13+0.10/6.7/6.3+0.58 ms cpu, 39->39->38 MB, 39 MB goal, 0 MB stacks, 0 MB globals, 12 P
gc 2 @0.120s 1%: 0.011+7.7+0.003 ms clock, 0.14+4.5/17/2.1+0.042 ms cpu, 75->77->76 MB, 77 MB goal, 0 MB stacks, 0 MB globals, 12 P
gc 3 @0.220s 2%: 0.020+18+0.058 ms clock, 0.24+12/31/10+0.70 ms cpu, 149->153->152 MB, 153 MB goal, 0 MB stacks, 0 MB globals, 12 P
gc 4 @0.275s 3%: 0.018+56+0.004 ms clock, 0.22+0/56/0+0.050 ms cpu, 176->176->175 MB, 304 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 5 @1.338s 0%: 0.13+3.7+0.002 ms clock, 1.6+0/4.2/3.2+0.033 ms cpu, 175->175->38 MB, 351 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 6 @2.352s 0%: 0.066+2.5+0.003 ms clock, 0.79+0/5.8/5.2+0.038 ms cpu, 38->38->38 MB, 77 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 7 @3.355s 0%: 0.098+2.6+0.003 ms clock, 1.1+0/7.3/4.3+0.046 ms cpu, 38->38->38 MB, 77 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 8 @4.358s 0%: 0.29+9.5+0.002 ms clock, 3.5+0/9.7/0+0.027 ms cpu, 38->38->38 MB, 77 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
gc 9 @5.368s 0%: 0.093+11+0.005 ms clock, 1.1+0/11/0+0.066 ms cpu, 38->38->38 MB, 77 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
*/
