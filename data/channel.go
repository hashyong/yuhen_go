package data

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

/*
通道

鼓励使用csp通道，以通信来代替内存共享，实现并发安全
通道 channel 行为类似消息队列。不限制收发人数，不可以重复消费

Don't communicate by sharing memory, share memory by communicate
CSP: Communicating Sequential Process


同步

没有数据缓冲区，必须双方到场直接交换数据
- 阻塞，直到另一方准备妥当或者通道关闭
- 可用过 cap == 0 判断为无缓冲通道
*/

func tSync() {
	quit := make(chan struct{})
	data := make(chan int)

	go func() {
		data <- 11
	}()

	go func() {
		defer close(quit)

		println(<-data)
		println(<-data)
	}()

	data <- 22
	<-quit

	/*
		22
		11
	*/
}

// 异步，通道自带固定大小的缓冲区，有数据或者空位的时候，不会阻塞
// 没有数据或者没有空位的时候，会阻塞
// 用cap和len 获取当前缓存区大小和当前缓冲数
func tAsync() {
	quit := make(chan struct{})
	data := make(chan int, 3)

	data <- 11
	data <- 22
	data <- 33
	//data <- 45 // 果然，当前没有空位了，还往channel里边写，自然是会阻塞住

	println(cap(data), len(data), cap(quit), len(quit))

	go func() {
		defer close(quit)

		println(<-data)
		println(<-data)
		println(<-data)
		println(<-data) // 果然如此，队列的缓冲区为3，当前队列已经空了。还在继续出队，没有数据出队就会导致阻塞
	}()

	println("begin")
	data <- 44
	println("end")
	<-quit
}

// 缓冲区大小仅仅是内部属性，不属于类型组成部分
// 通道变量本身就是指针，可以判断是否为同一对象或者nil
func tAsyncCmp() {
	var a, b chan int = make(chan int, 3), make(chan int)
	var c chan bool

	println(a == b)
	println(c == nil)
	println(a, unsafe.Sizeof(a))

	/*
		false //指针本身就不相等
		true  // 确实为空
		0xc00012c100 8 // 就是一个指针大小
	*/
}

/*
关闭
对于 closed 或者 nil 通道，规则如下
- 无论收发，nil 通道都会阻塞. nil通道是指 没有make的变量. 鉴于channel关闭后，所有基于其阻塞都会被解除，可用作通知
- 不能关闭 nil通道, 会panic
- 重复关闭通道，会panic
- 向已关闭通道发送数据，会panic
- 从已关闭通道接收数据，返回缓存数据或者零值

- 没有判断通道是否已被关闭的直接方法，只能透过收发模式获知
-
*/
func tCloseNil() {
	// var a chan int
	// close(a)
	/*
		panic: close of nil channel
	*/

	// println(<-a) //使用dlv attach之后可以看到，目前协程处于 chan receive (nil chan) 状态，已经阻塞, 所以这里会卡住
	/*
		(dlv) goroutines
		  Goroutine 1 - User: /mnt/c/Users/GolandProjects/yuhen/data/channel.go:113 yuhen/data.Channel (0x4f812f) [chan receive (nil chan)]
		  Goroutine 2 - User: /usr/local/go/src/runtime/proc.go:382 runtime.gopark (0x439256) [force gc (idle)]
		  Goroutine 3 - User: /usr/local/go/src/runtime/proc.go:382 runtime.gopark (0x439256) [GC sweep wait]
		  Goroutine 4 - User: /usr/local/go/src/runtime/proc.go:382 runtime.gopark (0x439256) [GC scavenge wait]
		  Goroutine 18 - User: /usr/local/go/src/runtime/proc.go:382 runtime.gopark (0x439256) [finalizer wait]
	*/
}

// 测试下非nil的channel
func tClose() {
	c := make(chan int, 2)
	c <- 1
	close(c)

	// channel已经被关闭，所以可以验证验证
	//close(c) // panic: close of closed channel
	//c <- 1 // panic: send on closed channel
	println(<-c)
	println(<-c)
	println(<-c)
	/*
		  // 不会阻塞
		1 // 获取到缓存的结果
		0 // 0
		0 // 0
	*/
}

// 为了避免重复关闭，可以包装close函数
// 也可以用类似方法封装send和recv等等操作
func closeChan[T any](c chan T) {
	defer func() {
		fmt.Println(recover())
	}()

	close(c)
}

func tCloseWrapper() {
	c := make(chan int)

	closeChan(c)
	closeChan(c)
}

// 保留关闭状态，为了并发安全，关闭和获取关闭状态应该保持同步
// 可以考虑使用 sync.RWMutex sync.Once 来优化设计

type tQueue[T any] struct {
	sync.Mutex

	ch     chan T
	cap    int
	closed bool
}

func NewtQueue[T any](cap int) *tQueue[T] {
	return &tQueue[T]{
		ch: make(chan T, cap),
	}
}

func (q *tQueue[T]) Close() {
	q.Lock()
	defer q.Unlock()

	if !q.closed {
		close(q.ch)
		q.closed = true
	}
}

func (q *tQueue[T]) IsClosed() bool {
	q.Lock()
	defer q.Unlock()

	return q.closed
}

func tConcurrencyClose() {
	var wg sync.WaitGroup
	q := NewtQueue[int](3)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer q.Close()
			println(q.IsClosed())
		}()
	}

	wg.Wait()
	// false
	// true
	// ...
	// true
}

// 利用nil通道阻塞的特性，可以阻止退出
func tNil() {
	//<-(chan struct{})(nil) // 确实会卡住，还是注释掉吧
}

// 收发
// 正常操作是 ok-idiom 或者 range模式
// ok == false 通道被关闭
// for ...range : 循环，直到通道关闭
func tReceive() {
	var wg sync.WaitGroup
	wg.Add(2)

	c := make(chan int)

	go func() {
		defer wg.Done()

		for {
			x, ok := <-c
			if !ok {
				return
			}
			println(x, ok)
			/*
				-1 true
				-2 true
				-3 true
			*/
		}
	}()

	go func() {
		defer wg.Done()
		defer close(c)

		c <- -1
		c <- -2
		c <- -3
	}()

	wg.Wait()
}

func tReceiveRange() {
	var wg sync.WaitGroup
	wg.Add(2)

	c := make(chan int)

	go func() {
		defer wg.Done()

		for x := range c {
			println(x)
			/*
				-1
				-2
				-3
			*/
		}
	}()

	go func() {
		defer wg.Done()
		defer close(c)

		c <- -1
		c <- -2
		c <- -3
	}()

	wg.Wait()
}

// 及时关闭channel，否则可能会导致死锁

// 单向，通道默认是双向的，不区分发送和接收端，可以限制方向获得更严谨的操作逻辑
func tOneway() {
	var wg sync.WaitGroup
	wg.Add(2)

	c := make(chan int)
	var send chan<- int = c
	var recv <-chan int = c

	// recv
	go func() {
		defer wg.Done()
		for x := range recv {
			println(x)
		}
	}()

	// send
	go func() {
		defer wg.Done()
		defer close(c)

		for i := 0; i < 3; i++ {
			send <- i
		}
	}()
}

// 不能在单向通道上做逆向操作
func tOnewayReverse() {
	c := make(chan int, 2)

	var recv <-chan int = c // 从队列也就是chan里边出数据，也就是recv
	var send chan<- int = c // 向队列也就是chan里边发送数据
	println(recv, send)

	// <-send
	// recv <- 1
	// close(recv)

	/*
		    	data/channel.go:332:4: invalid operation: cannot receive from send-only channel send (variable of type chan<- int)
				data/channel.go:333:2: invalid operation: cannot send to receive-only channel recv (variable of type <-chan int)
				data/channel.go:334:8: invalid operation: cannot close receive-only channel recv (variable of type <-chan int)
	*/

	// 无法将单向通道转换回去

	// 可以用make创建单向通道，没啥意义，因为 chan 必须有 收发才可以
	// q := make(<-chan struct{})
	//close(q) // data/channel.go:346:8: invalid operation: cannot close receive-only channel q (variable of type <-chan struct{})
}

// 选择
// 用select 语句处理多个通道，随机选择可用通道做收发操作
// 将失效通道设置为nil （阻塞 不可用）可用作结束判断
func tSelect() {
	var wg sync.WaitGroup
	wg.Add(4)

	cha, chb, chc := make(chan int), make(chan int), make(chan int)

	// recv
	go func() {
		defer wg.Done()

		for {
			x := 0
			ok := false

			// random
			select {
			case x, ok = <-cha:
				if !ok {
					cha = nil
				}
			case x, ok = <-chb:
				if !ok {
					chb = nil
				}
			}

			if (cha == nil) && (chb == nil) {
				return
			}

			println(x)
		}
	}()

	// send
	go func() {
		defer wg.Done()
		defer close(cha)
		defer close(chb)

		for i := 0; i < 10; i++ {
			// random
			select {
			case cha <- i:
			case chb <- i * 10:
			}
		}
	}()

	// 即使是单个通道，也会随机选择一个
	// recv
	go func() {
		defer wg.Done()

		for {
			x := 0
			ok := false

			// random
			select {
			case x, ok = <-chc:
				println("c1:", x)
			case x, ok = <-chc:
				println("c2:", x)
			}
			/*
				c2: 0
				c1: 10
				c2: 2
				c2: 3
				c1: 40
				c2: 5
				c2: 6
				c2: 7
				c1: 8
				c2: 9
				c1: 0
			*/

			if !ok {
				return
			}
		}
	}()

	// send
	go func() {
		defer wg.Done()
		defer close(chc)

		for i := 0; i < 10; i++ {
			// random
			select {
			case chc <- i:
			case chc <- i * 10:
			}
		}
	}()

	wg.Wait()

	// 所有通道都不可用，需要执行 default分支，避免阻塞
	// 空 select预计，一直阻塞或者死锁
}

// 缺省，利用default实现判断逻辑
func tDefault() {
	done := make(chan struct{})

	// 多chan 构成的数据区
	data := []chan int{
		make(chan int, 3),
	}

	go func() {
		defer close(done)

		for i := 0; i < 10; i++ {
			// 向随后一个slot添加数据，添加失败（满了），则创建新slot
			select {
			case data[len(data)-1] <- i:
			default:
				data = append(data, make(chan int, 3))
				data[len(data)-1] <- i
			}
		}
	}()

	<-done

	for _, c := range data {
		close(c)
		for x := range c {
			println("tDefault:", x)
		}
	}
}

// 反射，如果运行期才能确认通道数量，可以利用 反射 reflect实现
// 说实话，没太看明白. 已经看明白了，
// 可以支持更多的io类型的监听时间
func tReflect() {
	// 一个chanel。用作退出的信号
	exit := make(chan struct{})

	// 运行的时候动态创建
	// 看起来就是一个 chan 的 slice
	chans := make([]chan int, 0)
	// 添加了2个 chan进来
	chans = append(chans, make(chan int))
	chans = append(chans, make(chan int))

	// 起了一个goroutine
	go func() {
		// 日常操作，在主函数肯定在等这个<-chan, 一般就是用来阻塞下主函数
		defer close(exit)

		// 反射构建，select操作
		// 这是个啥新奇的方法，问下gpt
		/*
				reflect.SelectCase 是 Go 语言中反射包（reflect）中的一个类型，它表示一个可以在 select 语句中使用的 case 语句。

				在 Go 语言的 select 语句中，可以使用多个 case 语句来监听多个 channel，
				一旦其中一个 channel 中有数据可以读取，相应的 case 语句就会被执行。
				reflect.SelectCase 类型提供了一种在 select 语句中监听 channel 的灵活方式，它可以监听多个 channel 以及 default case，
				还可以指定相应的操作，例如读取或写入 channel。

				reflect.SelectCase 类型定义如下：

				```
				type SelectCase struct {
				    Dir  SelectDir
				    Chan Value
				    Send Value
				    Recv Value
				}
				```

				其中，Dir 表示操作方向，可以是 SelectSend（写入）或 SelectRecv（读取）；
				Chan 表示要监听的 channel；Send 表示要写入 channel 的值，如果该 case 是读取操作，则 Send 为零值；
				Recv 表示要接收的值的类型，如果该 case 是写入操作，则 Recv 为零值。

				reflect.SelectCase 和 select 语句在功能上有一些相似之处，它们都可以用于监听多个 channel，一旦其中一个 channel 中有数据可以读取或写入，相应的操作就会被执行。

			然而，reflect.SelectCase 和 select 语句还是有一些区别的：

			1. 语法不同：select 语句是一种特殊的语法结构，可以直接在 Go 语言的代码中使用，而 reflect.SelectCase 是一个类型，需要通过反射机制来使用。

			2. 灵活性不同：使用 select 语句可以方便地在代码中监听多个 channel，并且可以使用 default case 来处理超时等特殊情况。
				但是，它只能监听 channel，而 reflect.SelectCase 则可以监听任意类型的数据源，包括 channel、文件、网络连接等等。

			3. 性能不同：使用 select 语句可以在编译时进行优化，因此在性能上可能比 reflect.SelectCase 更好。
				而 reflect.SelectCase 的灵活性则会带来一些运行时的性能开销，例如反射操作会比直接读写 channel 更慢
		*/
		// 整了一个case的list，把之前的chans
		cases := make([]reflect.SelectCase, len(chans))
		for i, cs := range chans {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(cs),
			}
		}

		for {

			index, value, ok := reflect.Select(cases)

			// 检查并退出
			// 这里根据注释猜测下，假如返回false 表示可能当前channel已经被close了
			if !ok {
				// 置为空值
				chans[index] = nil

				// 当前为空的chanel的数量
				n := 0
				for _, c := range chans {
					if c == nil {
						n++
					}

					// 假如当前的chans都为nil，那可以返回了
					if n == len(chans) {
						println("所有的channel都结束了, 该返回了")
						return
					}
				}

				continue
			}

			println(index, value.Int(), ok)
		}
	}()

	chans[1] <- 101
	chans[0] <- 100

	for _, c := range chans {
		close(c)
	}

	<-exit
}

// 通常以工厂方法将 goroutine和channel想绑定
func newRecv[T any](cap int) (data chan T, done chan struct{}) {
	// 返回值有类型称之为 命名返回值，一般来说不推荐这种写法
	data = make(chan T, cap)
	done = make(chan struct{})

	go func() {
		defer close(done)

		// 一直循环消费，直到通道被close才会退出
		for v := range data {
			println(v)
		}
	}()

	return
}

func tPattern() {
	data, done := newRecv[int](3)
	for i := 0; i < 10; i++ {
		data <- i
	}

	close(data)
	<-done
}

// 如果channel阻塞且没有关闭，那么可能导致goroutine泄露 leak
// 解决办法使用 select default 或者 time.After设置超时
func tTimeout() {
	quit := make(chan struct{})
	c := make(chan int)

	go func() {
		defer close(quit)

		// 还是蛮巧妙的
		select {
		case x, ok := <-c:
			println(x, ok)
		case <-time.After(time.Second):
			println("channel timeout")
			return
		}
	}()

	<-quit
}

// 用通道实现信号量，semaphore，在同一时刻仅指定数量的goroutine参与工作
type sema struct {
	c chan struct{}
}

func newSema(n int) *sema {
	return &sema{
		c: make(chan struct{}, n),
	}
}

func (s *sema) acquire() {
	s.c <- struct{}{}
}

func (s *sema) release() {
	<-s.c
}

// 蛮好的，有点意思，控制并发还能这么玩
func tSemaphore() {
	var wg sync.WaitGroup

	// runtime:4
	// sema:2
	runtime.GOMAXPROCS(4)
	sem := newSema(3)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 一个结构体入到队列里边. 假如当前队列满了，这里就会卡住，入队会阻塞
			// 当前缓冲区的长度为2，所以前两个执行不会被阻塞，此处会正常执行
			// 当第三个协程开始执行的时候，因为当前缓存区已满，所以会阻塞住
			// 这样基本就能保障同一时刻只有俩协程在执行
			sem.acquire()
			// defer 里边会出队
			// 假如当前缓存区里边有数据，出队不会被阻塞
			// 假如当前缓冲区没有数据了，出队会阻塞
			// 当前逻辑为 会先入队再出队，所以不会出现
			defer sem.release()

			for i := 0; i < 3; i++ {
				// 等了2s
				time.Sleep(time.Second * 2)
				fmt.Println(id, time.Now())
				/*
					1 2023-06-06 10:59:15.129424414 +0800 CST m=+7.073192299
					4 2023-06-06 10:59:15.129541816 +0800 CST m=+7.073309701
					0 2023-06-06 10:59:15.129580416 +0800 CST m=+7.073348301

					0 2023-06-06 10:59:17.131826496 +0800 CST m=+9.075594381
					1 2023-06-06 10:59:17.131860096 +0800 CST m=+9.075627981
					4 2023-06-06 10:59:17.131864596 +0800 CST m=+9.075632

				*/
			}
		}(i)
	}

	wg.Wait()
}

// 鉴于通道本身就是一个并发安全的队列，可以用作 ID generator pool等用途
type pool[T any] chan T

func newPool[T any](cap int) pool[T] {
	return make(chan T, cap)
}

func (p pool[T]) get() (v T, ok bool) {
	select {
	case v = <-p:
		ok = true
	default:
	}
	return
}

func (p pool[T]) put(v T) bool {
	select {
	case p <- v:
		return true
	default:
	}
	return false
}

func tPool() {
	p := newPool[int](2)

	println(p.put(1))
	println(p.put(2))
	println(p.put(3))

	for {
		v, ok := p.get()
		if !ok {
			break
		}
		println(v)
	}

	/*
		true
		true
		false
		1
		2
	*/
}

// 捕获INT TREM 信号，顺便实现一个简易的 atexit函数
var exs = &struct {
	sync.RWMutex
	funcs   []func()
	signals chan os.Signal
}{}

func atexit(f func()) {
	exs.Lock()
	defer exs.Unlock()
	exs.funcs = append(exs.funcs, f)
}

func wait(code int) {
	// 信号注册
	if exs.signals == nil {
		exs.signals = make(chan os.Signal)
		signal.Notify(exs.signals, syscall.SIGINT, syscall.SIGTERM)
	}

	// 独立函数，确保 atexit 函数可以按照 filo顺序来执行
	// 不受 os.Exit影响
	// 果然，这里加上go之后 因为主程序挂掉，所以这里就GG了，也不会执行到
	func() {
		exs.RLock()
		for _, f := range exs.funcs {
			defer f()
		}
		exs.RUnlock()
		<-exs.signals
	}()

	//终止进程
	os.Exit(code)
}

func tExit() {
	atexit(func() {
		time.Sleep(time.Second * 3)
		println("atexit 1.......")
	})

	atexit(func() {
		println("atexit 2.......")
	})

	println("press ctrl+c to exit")
	wait(1)
}

// 通道本身就是队列。需要关心的是如何优雅的关闭通道
func tPatternQueue() {
	max := int64(100) // 最大发送计数
	m := 3            // 接收者数量
	n := 3            // 发送者数量

	var wg sync.WaitGroup
	wg.Add(m + n)

	data := make(chan int)      // 数据通道
	done := make(chan struct{}) // 结束通知

	// m recv
	for i := 0; i < m; i++ {
		go func(id int) {
			defer wg.Done()

			for {
				select {
				// 假如接收到这个信号，则返回, 这个语法糖真别扭，就不能搞一个 done.in() done.out()
				case <-done:
					return
				case v := <-data:
					println("gid:", id, "|val:", v)
				}
			}

		}(i)
	}

	// n send
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			defer func() { recover() }()

			for {
				select {
				case <-done:
					return
				case data <- id:
					if atomic.AddInt64(&max, -1) <= 0 {
						close(done)
						return
					}
				default:
				}
			}
		}(i)
	}

	wg.Wait()
}

// 将发往通道的数据打包，减少传输次数，可以有效的提升性能
// 通道内部实现有锁和数据复制的操作，单次发更多的数据（批处理）, 可以改善性能
const (
	cmax   = 50000000 // 数据统计上限
	cblock = 500      // 数据块大小
	ccap   = 100      // 缓冲区大小
)

//go:noinline
func normal() {
	done := make(chan struct{})
	c := make(chan int, ccap)

	go func() {
		defer close(done)

		count := 0
		for x := range c {
			count += x
		}
	}()

	for i := 0; i < cmax; i++ {
		c <- i
	}

	close(c)
	<-done
}

//go:noinline
func block() {
	done := make(chan struct{})
	c := make(chan [cblock]int, ccap)

	go func() {
		defer close(done)

		count := 0
		for a := range c {
			for _, x := range a {
				count += x
			}
		}
	}()

	for i := 0; i < cmax; i += cblock {
		// 使用数组对数据打包
		var b [cblock]int
		for n := 0; n < cblock; n++ {
			b[n] = i + n
			if i+n == cmax-1 {
				break
			}
		}
		c <- b
	}

	close(c)
	<-done
}

func tPerformance() {
	/*
		BenchmarkNoraml-12                     1        2679765546 ns/op            1528 B/op          5 allocs/op
		BenchmarkBlock-12                     14          77240504 ns/op          401535 B/op          3 allocs/op
	*/
}

// 如果channel 一直处于阻塞状态，那么会导致 goroutine无法结束和回收，形成资源泄露

func leak() chan byte {
	c := make(chan byte)

	go func() {
		buf := make([]byte, 0, 10<<20) // 10MB
		for {
			d, ok := <-c
			if !ok {
				return
			}
			buf = append(buf, d)
		}
	}()

	return c
}

func tLeak() {
	for i := 0; i < 5; i++ {
		leak()
	}

	for {
		time.Sleep(time.Second)
		runtime.GC()

		/*
				go build && GODEBUG=gctrace=1 ./yuhen

				gc 39 @55.163s 0%: 0.12+0.20+0.006 ms clock, 0.50+0/0.14/0.052+0.027 ms cpu, 51->51->51 MB, 102 MB goal, 0 MB stacks, 0 MB globals, 4 P (forced)
				gc 40 @56.164s 0%: 0.11+0.24+0.007 ms clock, 0.46+0/0.13/0.054+0.031 ms cpu, 51->51->51 MB, 102 MB goal, 0 MB stacks, 0 MB globals, 4 P (forced)
				gc 41 @57.165s 0%: 0.072+0.16+0.002 ms clock, 0.29+0/0.14/0.030+0.009 ms cpu, 51->51->51 MB, 102 MB goal, 0 MB stacks, 0 MB globals, 4 P (forced)
				gc 42 @58.165s 0%: 0.024+0.17+0.005 ms clock, 0.098+0/0.15/0.065+0.022 ms cpu, 51->51->51 MB, 102 MB goal, 0 MB stacks, 0 MB globals, 4 P (forced)
				gc 43 @59.166s 0%: 0.12+0.25+0.007 ms clock, 0.49+0/0.13/0.070+0.028 ms cpu, 51->51->51 MB, 102 MB goal, 0 MB stacks, 0 MB globals, 4 P (forced)

				GODEBUG="schedtrace=1000,scheddetail=1" ./yuhen

				SCHED 24131ms: gomaxprocs=4 idleprocs=4 threads=8 spinningthreads=0 needspinning=0 idlethreads=4 runqueue=0 gcwaiting=false nmidlelocked=0 stopwait=0 sysmonwait=false
			  P0: status=0 schedtick=142 syscalltick=124 m=nil runqsize=0 gfreecnt=4 timerslen=0
			  P1: status=0 schedtick=76 syscalltick=281 m=nil runqsize=0 gfreecnt=0 timerslen=1
			  P2: status=0 schedtick=33 syscalltick=7 m=nil runqsize=0 gfreecnt=5 timerslen=0
			  P3: status=0 schedtick=98 syscalltick=1 m=nil runqsize=0 gfreecnt=0 timerslen=0
			  M7: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=nil
			  M6: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=nil
			  M5: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=nil
			  M4: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=nil
			  M3: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=false lockedg=nil
			  M2: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=2 dying=0 spinning=false blocked=false lockedg=nil
			  M1: p=nil curg=17 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=false lockedg=17
			  M0: p=nil curg=nil mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=nil
			  G1: status=4(sleep) m=nil lockedm=nil
			  G17: status=6() m=1 lockedm=1
			  G2: status=4(force gc (idle)) m=nil lockedm=nil
			  G3: status=4(GC sweep wait) m=nil lockedm=nil
			  G4: status=4(GC scavenge wait) m=nil lockedm=nil
			  G18: status=4(finalizer wait) m=nil lockedm=nil
			  G19: status=4(trace reader (blocked)) m=nil lockedm=nil
			  G110: status=6() m=nil lockedm=nil
			  G87: status=6() m=nil lockedm=nil
			  G107: status=6() m=nil lockedm=nil
			  G115: status=4(chan receive) m=nil lockedm=nil
			  G103: status=6() m=nil lockedm=nil
			  G114: status=4(chan receive) m=nil lockedm=nil
			  G106: status=6() m=nil lockedm=nil
			  G102: status=6() m=nil lockedm=nil
			  G79: status=6() m=nil lockedm=nil
			  G78: status=6() m=nil lockedm=nil
			  G130: status=4(GC worker (idle)) m=nil lockedm=nil
			  G113: status=4(GC worker (idle)) m=nil lockedm=nil
			  G67: status=6() m=nil lockedm=nil
			  G116: status=4(chan receive) m=nil lockedm=nil
			  G117: status=4(chan receive) m=nil lockedm=nil
			  G118: status=4(chan receive) m=nil lockedm=nil
			  G119: status=4(GC worker (idle)) m=nil lockedm=nil
			  G24: status=4(GC worker (idle)) m=nil lockedm=nil
		*/
	}
}

func Channel() {
	tSync()
	tAsync()
	tAsyncCmp()

	tCloseNil()
	tClose()
	tCloseWrapper()
	tConcurrencyClose()
	tNil()

	tReceive()
	tReceiveRange()
	tOneway()
	tOnewayReverse()
	tSelect()
	tDefault()
	tReflect()
	tPattern()
	tTimeout()
	tSemaphore()
	tPool()
	//tExit()
	tPatternQueue()
	tPerformance()
	tLeak()
}

func Println(a ...any) {
	fmt.Println(len(a))
	if len(a) == 0 {
		fmt.Println(printFuncName())
	} else {
		fmt.Println(printFuncName(), a)
	}
}

func printFuncName() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Failed to get Caller info")
		return ""
	}

	funcName := runtime.FuncForPC(pc).Name()
	functionName := filepath.Ext(funcName)
	return fmt.Sprintf("%s|%s[%s]:%d|", printStartTime(), file, functionName[1:], line)
}

func printStartTime() string {
	start := time.Now()
	return fmt.Sprintf("%s", start.Format("2006-01-02 15:04:05.000"))
}
