package data

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
并发 concurrency ： 逻辑具备同时处理多个任务的能力
并行 parallesim: 物理上在同一时刻执行多个并发任务

需要程序以并发模型设计 执行的时候依据环境 单core或者多core不同，有不同的运行方式和效率
多核处理器真正同时执行多个任务，而单核只能以间隔切换方式运行
所以，并发是并行的必要条件，并行是并发的立项状态
并行需要多进程或者多线程支持，而并发可以在单线程上以协程实现

携程通常是在 单线程上，通过协作切换执行多个任务的并发设计，比如，将io等待时间，用来执行其他任务
且单线程 无竞态条件，可减少或者避免用锁，某些时间，用户空间下，较少的上下文切换的协程比多线程有更高的执行效率
*/

func Concurrency() {
	task()
}

/*
简单将goroutine 归为 coroutine 并不合适，它类似于 多线程和协程的综合体，最大限度的发挥多核处理能力，提升执行效率

- 按需扩张的极小初始栈，支持海量任务
- 高效无锁内存分配和复用，提升并发性能
- 调度器平衡任务队列。充分利用多核处理器
- 线程自动休眠和唤醒，减少内核开销
- 基于信号 signal实现抢占式任务调度

关键字go 将目标函数和参数打包（非执行）成并发的执行单元，放入待运行队列
无法确定执行时间 执行次序 以及执行线程，由调度器负责处理
*/

func testCon1() {
	// 打包并发任务 函数+参数
	// 并非立即执行
	go println("sdfsdf")
	go func(x int) { println(x) }(123)

	// 上述的并发任务 会被其他ji饿线程取走
	// 执行时间未知，下面这行可能会先输出

	println("main")
	time.Sleep(time.Second)
}

// 参数立即计算并复制
var c int

func inc() int {
	c++
	return c
}

func testCon2() {
	a := 100

	// 立即计算出参数(100,1)并复制
	// 内部sleep 母的是让main先输出

	go func(x, y int) {
		time.Sleep(time.Second)
		println("go:", x, y) // go: 100 1
	}(a, inc())

	a += 100
	println("main:", a, inc()) // main: 200 2

	time.Sleep(time.Second * 1)
}

// 所有用户代码都已 goroutine 执行，包括main.main函数
// 进程结束，不会等待其他正在执行或者尚未执行的任务
func testEnd() {

	go func() {
		defer println("g fin done")
		time.Sleep(time.Second)
	}()

	defer println("main.done")
}

/*
等待任务结束, 可以做的比time.sleep 更优雅一些

- channel: 信号通知
- WaitGroup: 等待多个任务结束
- Context: 上下文通知
- Mutex：锁阻塞

如果只是一次性通知行为，可以使用空机构，只要关闭通道，等待阻塞接口解除
*/
func testWait() {
	q := make(chan struct{})

	go func() {
		defer close(q)
		println("done")
	}()

	<-q
}

// 添加计数 WaitGroup.Add  应该在创建任务和等待之前，否则会导致 等待提前解除
// 可以有多处等待，实现群体性通知
func testWaitGroup() {
	var wg sync.WaitGroup
	total := 3
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(id int) {
			defer wg.Done()
			println(id, "done.")
		}(i)
	}

	wg.Wait()
}

func testWaitGroupAll() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		// 可以写在此处，应对未知的循环数
		// 但不可以放在下边 goroutine 函数内
		// 因为它肯恶搞未执行，下边wait先结束了
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			println(id, "s done.")
		}(i)
	}

	wg.Wait()
}

// 通过context来控制
func testContextNotify() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		println("ctx done.")
	}()

	<-ctx.Done()
}

// 利用锁实现同步，在其他语言很常见，但是go更倾向于以通信代替
func testMutexNotify() {
	var lock sync.Mutex

	lock.Lock()
	go func() {
		defer lock.Unlock()
		println("lock goroutine done.")
	}()

	lock.Lock()
	lock.Unlock()
	println("exit.")
}

/*
终止
主动结束任务，有以下集中方式
- 调用 runtime.Goexit 来终止任务
- 调用 os.Exit 结束进程

在任务调用 堆栈 call stac 的任何位置调用 runtime.Goexit 都能立即终止任务
结束前，defer 被执行，其他任务不受影响
*/
func testExit() {
	q := make(chan struct{})

	go func() {
		defer close(q)
		defer println("done")

		// 如果将f里边的Goexit 换为return
		// 那么只是结束了f(), 而非整个堆栈
		e()
		f()
		g()
	}()

	<-q
}

func g() {
	println("g what")
}

func f() {
	println("f")
	runtime.Goexit()
}

func e() {
	println("e")
}

func testExitMain() {
	defer println("defer exit main")
	println("test exit main begin")
	q := make(chan struct{})

	go func() {
		defer close(q)
		defer println("exit main g done.")
		time.Sleep(time.Second)
	}()

	<-q
	// 这里会卡住，因为还有其他goroutine在执行
	// runtime.Goexit()
	/*
		fatal error: no goroutines (main called runtime.Goexit) - deadlock!
	*/
	println("end")
}

// os.Exit 可以在任意位置结束进程，而不用等待其他任务，也不执行defer
func testOsExit() {
	go func() {
		defer println("g done")
		time.Sleep(time.Second)
	}()

	defer println("main done")

	//os.Exit(-1)
}

// 运行的时候可能会创建很多线程，但任何时候都只有几个参与并发任务执行，其他处于休眠状态
// 默认与逻辑处理器（logic core）数量相等，或使用GOMAXPROCS  修改

//go:noinline
func sumCon() (n int) {
	//for i := 0; i < math.MaxUint32; i++ {
	for i := 0; i < 11; i++ {
		n += i
	}

	return
}

// 可以使用 time GOMAXPROCS=1 ./test
// 可以使用 time GOMAXPROCS=2 ./test
func testLimit() {
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			sumCon()
		}()
	}

	wg.Wait()
}

/*
	调度

除了运行的时候自动调度外，某些时候需要手动控制任务运行

挂起
暂时挂起任务，释放线程去执行其他任务
当前任务被放回任务队列。等待下次被某个线程重新获取后 继续执行

也就是说，一个任务不一定由同一线程执行，实际上，除了主动协作调度之外, 还要考虑运行的时候抢占调度这些因素
长时间运行的任务会被暂停. 让其他等待任务由机会执行，以确保公平
*/
func testSchedule() {
	// 限制任务并发数
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	a, b := make(chan struct{}), make(chan struct{})

	go func() {
		defer wg.Done()

		<-a // 等待 管道a的通知
		for i := 0; i < 5; i++ {
			fmt.Println("a-------", i)
		}
	}()

	go func() {
		defer wg.Done()

		<-b // 等待b通知
		for i := 0; i < 5; i++ {
			fmt.Println("b-------", i)
			if i == 2 {
				runtime.Gosched()
			}
		}
	}()

	// 安排执行次序
	close(b)
	close(a)

	/*
		b------- 0
		b------- 1
		b------- 2 // Gosched, pause
		a------- 0
		a------- 1
		a------- 2
		a------- 3
		a------- 4
		b------- 3 // CONT
		b------- 4
	*/
}

// 发令，暂停一批任务，直到某个信号发出
// 也可以反向使用 sync.WaitGroup, 让多个 goroutine wait，然后 main Done，类似还可以使用 信号量，控制启动任务数量
func testSchedNotice() {
	var wg sync.WaitGroup
	r := make(chan struct{})

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			<-r // 阻塞，等待信号
			println("--------", id)
		}(i)
	}

	close(r)
	wg.Wait()

	/*
		-------- 9
		-------- 0
		-------- 1
		-------- 2
		-------- 3
		-------- 4
		-------- 5
		-------- 6
		-------- 7
		-------- 8
	*/
}

/*
次序
多个任务按特定次序执行
*/
func testSchedSort() {
	const CNT = 5

	var wg sync.WaitGroup
	wg.Add(CNT)

	var chans [CNT]chan struct{}
	for i := 0; i < CNT; i++ {
		chans[i] = make(chan struct{})

		go func(id int) {
			defer wg.Done()
			<-chans[id]
			println("*********************", id)
		}(i)
	}

	// 次序（延时，给调度器时间处理）
	for _, x := range []int{2, 3, 0, 1, 4} {
		close(chans[x])
		//time.Sleep(time.Millisecond * 10)
		// 假如将这行注释掉，输出的顺序就不一致了。这是为毛，待以后再仔细看看
		// TODO:研究研究
		/*
		********************* 4
		********************* 0
		********************* 1
		********************* 2
		********************* 3
		 */
	}

	wg.Wait()

	/*
	********************* 2
	********************* 3
	********************* 0
	********************* 1
	********************* 4
	 */
}

// 存储，默认情形下，任务和线程并非绑定关系，所以不能用TLS之类的本地存储
// 亿参数方式传入外部容器，但要避免和其他并发任务竞争 data race
// runtime.LockOSThread 可锁定调用线程，常用于syscall 和 CGO
func testSchedStorage() {
	var gls [2]struct {
		id  int
		ret int
	}

	var wg sync.WaitGroup
	wg.Add(len(gls))

	for i := 0; i < len(gls); i++ {
		go func(id int) {
			defer wg.Done()

			gls[id].id = id
			gls[id].ret = id * 100
		}(i)
	}

	wg.Wait()
	fmt.Println(gls)
	/*
		[{0 0} {1 100}]
	*/
}

func task() {
	testCon1()
	testCon2()

	testEnd()
	testWait()
	testWaitGroup()
	testWaitGroupAll()
	testContextNotify()
	testMutexNotify()

	testExit()
	testExitMain()
	testOsExit()

	testLimit()
	testSchedule()
	testSchedNotice()
	testSchedSort()
	testSchedStorage()
}
