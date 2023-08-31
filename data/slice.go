package data

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

/*slice以指针引用底层数组片段，限定读写区域，类似胖指针，而非动态数组或者数组指针


+----------+        +-----------------------+
|  array   +------> |0 |1 |2 |3 |4 |...  |99|
+----------+        +-----------------------+
|  len     |
+----------+         array
|  cap    |
+----------+

* 引用类型，未init值未nil
* 基于数组 初始化值 或者make函数创建

* 函数len 返回元素数量 cap返回容量
* 仅仅能判断是否为nil，不支持其他 == != 操作

* 以索引访问元素，可获取底层数据元素指针
* 底层数组可能在堆上分配


+-----------------------------+
|0 |1 |2 |3 |4 |5 |6 |7 |8 |9 |   array
+-----------------------------+
	  |           |     |
	  |           |     |        slice: [low:high:max]
	  | <+s.len+> |     |
	  |                 |        len=high-low
	  | <---+s.cap+---> |        cap=max -low

*/

func Slice() {
	a := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s := a[2:6:8]
	fmt.Println(a, s, len(s), cap(s))
	// [0 1 2 3 4 5 6 7 8 9] [2 3 4 5] 4 6

	// 引用原数组
	fmt.Printf("a: %p ~ %p\n", &a[0], &a[len(a)-1])
	fmt.Printf("s: %p ~ %p\n", &s[0], &a[len(s)-1])
	/*
		a: 0xc0000f8000 ~ 0xc0000f8048
		s: 0xc0000f8010 ~ 0xc0000f8018
	*/

	// slice.len 限定索引或迭代读取的数据范围
	// slice.cap 重新切片 reslice的允许范围
	//fmt.Println(s[5]) //panic: runtime error: index out of range [5] with length 4

	for i, x := range s {
		fmt.Printf("s[%d]: %d\n", i, x)
	}
	/**
			s[0]: 2
			s[1]: 3
			s[2]: 4
			s[3]: 5

			+-----------------------------+
			|0 |1 |2 |3 |4 |5 |6 |7 |8 |9 |   a: [10]int
			+-----------------------------+

		          +-----------------+
		          |2 |3 |4 |5 |6 |7 |         s: a[2:6:8]
		          +-----------------+
	**/

	// 转换
	// slice 可以直接转回数组或者数组指针
	var a1 [4]int = [...]int{0, 1, 2, 3}

	// array -》slice ：指向原数组
	var s1 []int = a1[:]
	println(&a1[0] == &s1[0]) //true

	// slice -> array 复制底层数组
	a2 := [4]int(s1)
	fmt.Println(&a2[0] == &a1[0]) // false

	// slice-> *array: 返回底层数组的指针
	p2 := (*[4]int)(s1)
	fmt.Println(p2 == &a1) // true

	// 指针相关，可以获取元素指针，但是不能以切片指针访问元素
	// 指针相关
	// * &array 指向数组的指针
	// &array[0] 指向某个元素的指针
	// [...]*int{&a, &b} 元素是指针的数组

	// 复制，数组传递会整个复制数组，可以使用切片或者指针代替

	// 这个从array 向slice转换的过程是slice是指向原数组的
	a3 := [...]int{0, 1, 2, 3}
	s2 := a3[:]

	p := &s2     // 切片指针, 指向header的内存
	e1 := &s2[1] // 元素指针

	// 数组指针直接指向元素所在内存
	// 切片指针指向header的内存

	//_ = p[1]
	//Invalid operation: 'p[1]' (type '*[]int' does not support indexing)

	_ = (*p)[1]

	// 元素指针指向数组
	*e1 += 100

	fmt.Println(e1 == &a3[1])
	fmt.Println(a3)

	// 基于数组指针创建切片
	var p3 *[4]int = &a3
	var s3 []int = p3[:]
	fmt.Println(&s3[2] == &a3[2])

	// 基于非数组指针创建切片, 最新版本编译不过
	//p3 := (*byte)(unsafe.Pointer(&a3[2])) // 元素指针
	//var s4 []byte = unsafe.Slice(p3, 8)
	//fmt.Println(s4)

	var a4 [3]byte = [...]byte{'a', 'b', 'c'}
	var s4 []byte = a4[:]
	// 返回切片底层数组首个元素指针
	fmt.Println(unsafe.SliceData(s4) == &a4[0]) // true

	// 构建字符串，返回底层数组指针
	var str string = unsafe.String(&a4[0], len(s4))
	fmt.Println(unsafe.StringData(str) == &a4[0]) // true

	// slice
	// 本身只是3个整数字段的小对象，可以直接值传递, 确实看汇编代码全是走的寄存器
	// 另外编译器尽可能将底层数组分配在栈上，以提升性能
	// 这明显全分配在栈上了，效率还是很高的， 开始疑惑了，开始补课，为什么栈上的效率要高于堆上
	s5 := []int{1, 2, 3}
	//MOVQ $0x1, 0x50(SP) ;array
	//MOVQ $0x2, 0x58(SP) ;
	//MOVQ $0x3, 0x60(SP) ;
	fmt.Println(sum(s5)) // header
	//LEAQ  0x50(SP), AX  ; .array = AX
	//MOVL  0x3, BX		  ; .len   = BX
	//MOVQ  BX, CX		  ; .cap   = CX
	//call yuhen/data.sum(SB)

	// 看看汇编是如何做的
	// go build
	// lensm -watch xxxx

	// jagged array 如果元素也是切片 即可实现交错数组的功能
	// 有点意思
	s6 := [][]int{
		{1, 2},
		{10, 20, 30},
		{100},
	}
	s6[1][2] += 1000
	fmt.Println(s6) // [[1 2] [10 20 1030] [100]]

	// 追加 append
	/*
		数据追加到s[len]
		返回新slice对象，通常复用原内存
		超出s.cap限制，即便底层数组尚有空间，也会重新分配内存，并复制数据
		新分配内存大小，通常是s.cap*2 比较大的切片会减少倍数，避免浪费
	*/

	a5 := [...]int{0, 1, 2, 3, 99: 0}
	fmt.Printf("a5:%p ~ %p\n", &a5[0], &a5[len(a5)-1]) // 打出内存地址大致的分布

	s7 := a5[:4:8]
	pslice("s7", &s7)

	// 开始尝试往里边追加看看内存的变化
	// 未超出s.cap的限制, 因为cap还有4空余，所以添加3个元素不会触发内存的重新分配
	s7 = append(s7, []int{4, 5, 6}...)
	pslice("s7", &s7)

	// 超出s.cap, 需要重新分配内存，复制数据
	s7 = append(s7, []int{11, 22}...)
	pslice("s7", &s7)

	fmt.Println("a5:", a5[:len(s7)])
	fmt.Println("s7:", s7)
	/*
		a5:0xc0000f4380 ~ 0xc0000f4698
		s7: reflect.SliceHeader{Data:0xc0000f4380, Len:4, Cap:8}
		s7: reflect.SliceHeader{Data:0xc0000f4380, Len:7, Cap:8}
		s7: reflect.SliceHeader{Data:0xc0000de100, Len:9, Cap:16}
		a5: [0 1 2 3 4 5 6 0 0]
		s7: [0 1 2 3 4 5 6 11 22]
	*/

	// 为slice预留足够的cap，可以有效减少内存分配和复制
	// s := make([]int, 0, 10000)

	// copy 在两个slice之间复制数据
	/*
		允许指向同一个数组
		允许目标区间重叠
		复制长度以较短 len slice为准
	*/
	s8 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// 在同一底层数组的不同区间复制
	src := s8[5:8]
	dst := s8[4:]

	n := copy(dst, src)
	fmt.Println(n, s8)
	// 3 [0 1 2 3 5 6 7 7 8 9]. 可以看到复制了3个元素 以短的为主

	//可以在不同数组之间复制
	dst = make([]int, 6)
	n = copy(dst, src)
	fmt.Println(n, dst)
	//3 [6 7 7 0 0 0]

	// 可以从字符串中复制数据到字节slice
	b := make([]byte, 3)
	n = copy(b, "abcde")
	fmt.Println(n, b) //3 [97 98 99]

	//slicePerfomance()

	mainStack()
	mainQueue()
	mainLoopQueue()
}

// 如果切片引用大数组，需要考虑新建并复制，及时释放大数组，避免内存浪费
// go build && GODEBUG=gctrace=1 ./xxx
func slicePerfomance() {
	d := [...]byte{100 << 20: 10}
	//--1--------
	// s := d[:2]
	/*
		gc 1 @0.034s 0%: 0.012+5.3+0.036 ms clock, 0.14+0/2.8/0+0.43 ms cpu, 100->100->100 MB, 100 MB goal, 0 MB stacks, 0 MB globals, 12 P
		gc 2 @3.608s 0%: 0.045+0.11+0.002 ms clock, 0.54+0/0.11/0+0.028 ms cpu, 100->100->100 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 3 @4.609s 0%: 0.047+0.30+0.002 ms clock, 0.57+0/0.17/0+0.032 ms cpu, 100->100->100 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 4 @5.609s 0%: 0.11+0.16+0.003 ms clock, 1.3+0/0.15/0+0.043 ms cpu, 100->100->100 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 5 @6.610s 0%: 0.048+0.085+0.002 ms clock, 0.58+0/0.11/0+0.031 ms cpu, 100->100->100 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 6 @7.610s 0%: 0.083+0.17+0.003 ms clock, 1.0+0/0.14/0+0.037 ms cpu, 100->100->100 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
	*/

	//--2------
	s := make([]byte, 2)
	copy(s, d[:])
	/*
		gc 1 @0.035s 0%: 0.012+6.5+0.034 ms clock, 0.15+0/3.4/0+0.41 ms cpu, 100->100->100 MB, 100 MB goal, 0 MB stacks, 0 MB globals, 12 P
		gc 2 @3.332s 0%: 0.057+0.26+0.002 ms clock, 0.68+0/0.11/0+0.032 ms cpu, 100->100->0 MB, 200 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 3 @4.333s 0%: 0.038+0.13+0.002 ms clock, 0.46+0/0.16/0+0.030 ms cpu, 0->0->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 4 @5.334s 0%: 0.031+0.16+0.002 ms clock, 0.38+0/0.16/0+0.031 ms cpu, 0->0->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 5 @6.334s 0%: 0.096+0.27+0.019 ms clock, 1.1+0/0.27/0+0.23 ms cpu, 0->0->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
		gc 6 @7.335s 0%: 0.16+0.41+0.004 ms clock, 2.0+0/0.34/0+0.056 ms cpu, 0->0->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P (forced)
	*/

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		runtime.GC()
	}

	runtime.KeepAlive(&s)
}

func pslice(name string, p *[]int) {
	fmt.Printf("%s: %#v\n",
		name, *(*reflect.SliceHeader)(unsafe.Pointer(p)))
}

//go:noinline
func sum(s []int) (n int) {
	for _, v := range s {
		n += v
	}

	return n
}

// 栈的一个实现
// 先进后出，栈

type Stack []int

func NewStack() *Stack {
	s := make(Stack, 0, 10)
	return &s
}

func (s *Stack) Push(v int) {
	*s = append(*s, v)
}

func (s *Stack) Pop() (int, bool) {
	if len(*s) == 0 {
		return 0, false
	}

	x, n := *s, len(*s)

	v := x[n-1]
	*s = x[:n-1]

	return v, true
}

func mainStack() {
	s := NewStack()

	// push
	for i := 0; i < 5; i++ {
		s.Push(i + 111)
	}

	// pop
	for i := 0; i < 7; i++ {
		fmt.Println(s.Pop())
	}
}

// Queue 队列，先进先出
type Queue []int

func NewQueue() *Queue {
	q := make(Queue, 0, 10)
	return &q
}

func (q *Queue) Put(v int) {
	*q = append(*q, v)
}

func (q *Queue) Get() (int, bool) {
	if len(*q) == 0 {
		return 0, false
	}

	x := *q
	v := x[0]

	// copy (x, x[1:])
	// *q = x[:len(x)-1]
	*q = append(x[:0], x[1:]...)
	// 目前来看，queue的实现还是相对比较简单的，从数组的头部出 尾部进，使用尾插的方式，每次出队列直接将数组元素向前挪一个位置

	return v, true
}

func mainQueue() {
	q := NewQueue()

	// put
	for i := 0; i < 5; i++ {
		q.Put(i + 111)
	}

	// ge/**/t
	for i := 0; i < 7; i++ {
		fmt.Println(q.Get())
	}
}

// LoopQueue 容量固定, 且无需移动数据的环形队列
// 这是个啥
type LoopQueue struct {
	sync.Mutex // 如果没有写实际的变量名，就相当于继承
	data       []int
	head       int
	tail       int
}

// NewLoopQueue 目前看这个方法还是比较简单的
func NewLoopQueue(cap int) *LoopQueue {
	return &LoopQueue{data: make([]int, cap)}
}

// Put 因为数组已经提前分配好了，所以这里只是找到最后的位置直接插入即可
func (l *LoopQueue) Put(v int) bool {
	// 加锁是因为避免多线程问题，假如两个goroutine同时执行，会导致tail异常和结果异常，产生冲突
	l.Lock()
	defer l.Unlock()

	// 表示当前队列已经满了， 无法再向里边插入
	if l.tail-l.head == len(l.data) {
		return false
	}

	// tail始终指向当前最后一个元素的下一个位置
	l.data[l.tail%len(l.data)] = v
	l.tail++

	return true
}

func (l *LoopQueue) Get() (int, bool) {
	l.Lock()
	defer l.Unlock()

	// 表示当前队列是空的, 没有可以出队的元素
	// 或者已经写满了， 同时也全部弹出了
	if l.tail-l.head == 0 {
		return 0, false
	}

	// 返回的是头部的元素 也还ok
	v := l.data[l.head%len(l.data)]
	l.head++

	return v, true
}

func mainLoopQueue() {
	fmt.Println("lp begin")
	rand.Seed(time.Now().UnixNano())

	const max = 10000
	src := rand.Perm(max) // 随机测试数据
	dst := make([]int, 0, max)

	q := NewLoopQueue(6)

	var wg sync.WaitGroup
	wg.Add(2)

	//put
	go func() {
		defer wg.Done()
		// 直到
		for _, v := range src {
			for {
				// 假如失败，会持续往里边写，直到成功
				if ok := q.Put(v); !ok {
					continue
				}
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		// 一直弹出 直到向dst写满为止, 之后才会退出循环
		for len(dst) < max {
			if v, ok := q.Get(); ok {
				dst = append(dst, v)
				continue
			}
		}
	}()

	wg.Wait()

	if *(*[max]int)(src) != *(*[max]int)(dst) {
		log.Fatalln("xxx")
	}
	fmt.Printf("%+v\n", q)
}
