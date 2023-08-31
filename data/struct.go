package data

import (
	"fmt"
	"os"
	"reflect"
	"time"
	"unsafe"
)

/*结构体将多个字段打包为一个复合类型

* 直属字段名必须唯一
* 支持自身指针类型成员
* 可用_忽略字段名

* 支持匿名接口和空结构
* 支持匿名嵌入和其他类型
* 支持为字段添加标签

* 仅所有字段全部支持 才可以做相等操作
* 可用指针选择字段 单不支持多级指针

* 字段的名称 类型 标签 以及排列的顺序属类型组成部分
* 除对其外 编译器不会优化或者调整内存布局
 */

// Struct 主函数
func Struct() {
	inits()
	equals()
	empty()
	tag()
	anonymous()
	mem()
}

/*
	内存

不管结构体包含多少字段，其内存总是一次性分配，各个字段 包括匿名字段成员 在相邻的地址空间按照定义顺序 含对齐 排列
对于引用类型 字符串和指针，结构体内存中只包含其基本数据

借助于unsafe相关函数 输出所有字段的偏移量和长度
*/
func mem() {
	type Point struct {
		x, y int
	}

	type Value struct {
		id   int
		name string
		data []byte
		next *Value
		Point
	}

	v := Value{
		id:    1,
		name:  "test",
		data:  []byte{1, 2, 3, 4},
		Point: Point{x: 100, y: 200},
	}
	fmt.Printf("%p ~ %p, size: %d, align: %d\n",
		&v, unsafe.Add(unsafe.Pointer(&v), unsafe.Sizeof(v)), unsafe.Sizeof(v), unsafe.Alignof(v))

	s := "%p, %d, %d\n"
	fmt.Printf(s, &v.id, unsafe.Offsetof(v.id), unsafe.Sizeof(v.id))
	fmt.Printf(s, &v.name, unsafe.Offsetof(v.name), unsafe.Sizeof(v.name))
	fmt.Printf(s, &v.data, unsafe.Offsetof(v.data), unsafe.Sizeof(v.data))
	fmt.Printf(s, &v.next, unsafe.Offsetof(v.next), unsafe.Sizeof(v.next))
	fmt.Printf(s, &v.x, unsafe.Offsetof(v.x), unsafe.Sizeof(v.x))
	fmt.Printf(s, &v.y, unsafe.Offsetof(v.y), unsafe.Sizeof(v.y))

	/*
		0xc00020a000 ~ 0xc00020a048, size: 72, align: 8
		0xc00020a000, 0, 8     	int
		0xc00020a008, 8, 16		string{ptr,len}
		0xc00020a018, 24, 24	slice{ptr,len,cap}
		0xc00020a030, 48, 8		pointer.go
		0xc00020a038, 56, 8		int
		0xc00020a040, 64, 8		int
	*/

	// 对齐以所有字段中最长的基础类型宽度为准, 划重点 为基础类型
	// 编译器目的：为了最大限度减少读写所需要的指令，也因为某些架构平台自身的要求
	v1 := struct {
		a byte
		b byte
		c int32 // * 根据它的长度来对齐 4byte来对齐
	}{}
	/*
		+-----------+-----------+
		|a |b |? |? |    c      |
		+-----------------------+
		| <--+4+--> | <--+4+--> |

	*/

	v2 := struct {
		a byte
		b byte // * 1byte 来对齐
	}{}
	/*
		+---+---+
		| a | b |
		+---+---+
		0   1   2

	*/

	v3 := struct {
		a byte
		b []int // 8字节来对齐 b.ptr b.len b.cap
		c byte
	}{}
	/*
		+-+-----+------+------+------+-+----+
		|a| ... |b.ptr |b.len |b.cap |c|... |
		+-+-----+------+------+------+-+----+
		|   8   |  8   |  8   |  8   |   8  |
	*/

	fmt.Printf("v1:%d, %d\n", unsafe.Alignof(v1), unsafe.Sizeof(v1))
	fmt.Printf("v2:%d, %d\n", unsafe.Alignof(v2), unsafe.Sizeof(v2))
	fmt.Printf("v3:%d, %d\n", unsafe.Alignof(v3), unsafe.Sizeof(v3))
	/*
		v1:4, 8
		v2:1, 2
		v3:8, 40
	*/

	d := struct {
		a [3]byte
		x int32
	}{
		a: [3]byte{1, 2, 3},
		x: 100,
	}
	/*
	  +--+--+--+--+-----------+
	  |1 |2 |3 |? |    100    |
	  +--+--+--+--+-----------+
	  |    4      |     4     |
	*/
	fmt.Printf("d:%d, %d\n", unsafe.Alignof(d), unsafe.Sizeof(d)) // d:4, 8

	// 空结构
	// 如果空结构是最后一个字段，那么将其当作长度为1的类型 避免越界
	// 其他零长度对象 [0]int 类似
	v10 := struct {
		a struct{}
		b int
		c struct{}
	}{}
	fmt.Printf("%p ~ %p, size: %d, align: %d\n",
		&v10, unsafe.Add(unsafe.Pointer(&v10), unsafe.Sizeof(v10)), unsafe.Sizeof(v10), unsafe.Alignof(v10))
	s10 := "%p, %d, %d\n"
	fmt.Printf(s10, &v10.a, unsafe.Offsetof(v10.a), unsafe.Sizeof(v10.a))
	fmt.Printf(s10, &v10.b, unsafe.Offsetof(v10.b), unsafe.Sizeof(v10.b))
	fmt.Printf(s10, &v10.c, unsafe.Offsetof(v10.c), unsafe.Sizeof(v10.c))
	/*
		0xc0000b2510 ~ 0xc0000b2520, size: 16, align: 8
		0xc0000b2510, 0, 0
		0xc0000b2510, 0, 8
		0xc0000b2518, 8, 0
	*/

	// 如果仅仅有一个空结构字段，那么按照1对齐，只不过其长度为0，且指向 zerobase变量
	// zerobase 是一个全局变量uintptr，所有空结构体的地址都为它，省内存
	v20 := struct {
		a struct{}
	}{}
	fmt.Printf("%p, size:%d, align:%d\n", &v20, unsafe.Sizeof(v20), unsafe.Alignof(v20))
	//0x6f38b0, size:0, align:1
}

/*
匿名字段 anonymous filed, 是指没有名字，仅有类型的字段，也被称作嵌入字段，嵌入类型
- 隐式以类型名作为字段名
- 嵌入类型与其指针类型隐式字段名相同
- 可以像直属字段那样直接访问嵌入类型成员
*/
func anonymous() {
	type Attr struct {
		perm int
	}

	type File struct {
		Attr // Attr:Attr
		name string
	}

	f := File{
		name: "test",
		// 必须显示的init
		//perm: 0544, //unknown field perm in struct literal of type File

		Attr: Attr{
			perm: 0544,
		},
	}

	// 像直属字段访问嵌入字段
	fmt.Println(f.name, f.perm) //test 356

	// 嵌入其他包里的类型，隐式字段不包含包名
	type data struct {
		os.File // File:os.File
	}

	d := data{File: os.File{}}
	fmt.Println(d) //{{<nil>}}

	// 除了接口指针和多级指针以外的任何命名类型都可作为匿名字段
	type data1 struct {
		*int // int:*int(星号不是名字的组成部分)
		//int  // data/struct.go:78:3: int redeclared
		//        data/struct.go:76:3: other declaration of int

		fmt.Stringer
		//*fmt.Stringer // embedded field type cannot be a pointer.go to an interface
	}
	_ = data1{
		int:      nil,
		Stringer: nil,
	}

	// 重名
	// 直属字段与嵌入类型成员 存在重名问题
	// 编译器优先选择 直属命名字段 或者按照嵌入层次逐级查找匿名类型成员
	// 如果匿名类型成员被外层同名遮蔽，那么必须用显示的字段名
	type Files struct {
		name []byte
	}

	type Data struct {
		Files
		name string // 和File.name重名
	}

	d1 := Data{
		name:  "data",
		Files: Files{[]byte("file")},
	}

	d1.name = "data22"                  // 优先选择直属命名字段
	d1.Files.name = []byte("sdf")       // 显示的字段名
	fmt.Println(d1.name, d1.Files.name) //data22 [115 100 102]

	// 如果多个相同层次的匿名类型成员重名 只能用显示的字段名 编译器不发确定用那个
	type Gift struct {
		name string
	}

	type Gift1 struct {
		name string
	}

	type Gift2 struct {
		Gift
		Gift1
	}

	g := Gift2{}
	//g.name = "name" // ambiguous selector g.name
	g.Gift1.name = "xx"
	g.Gift.name = "x"

	// go并非传统意义的oop，包括封装 继承 多态，仅仅实现了最小机制
	// 匿名嵌入是组合 而非继承
	// 结构体 虽然承载了class的功能，但无法多态，只能以方法集配合接口实现
}

/*
标签 tag不是注释，是对字段类型进行描述的元数据
- 不是数据成员 是类型的组成部分
- 内存布局相同 允许进行显示转换
*/
func tag() {
	type user struct {
		id int `id`
	}

	type user2 struct {
		id int `uid`
	}

	u1 := user{1}
	u2 := user2{2}
	// 类型不同 因为tag不同
	//_ = u1 == u2 // invalid operation: u1 == u2 (mismatched types user and user2)

	// 内存布局相同 支持转换
	u1 = user(u2)
	fmt.Println(u1) //{2}

	// 运行期 可用反射获取标签信息 经常被用作格式校验 数据库关系映射
	// 感觉还是挺方便的
	type User struct {
		id   int    `field:"uid" type:"integer"`
		name string `field:"name" type:"text"`
	}

	t := reflect.TypeOf(User{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Name, f.Type.Name(), f.Tag.Get("field"), f.Tag.Get("type"))
	}
	/*
		id int uid integer
		name string name text
	*/
}

/*
空结构 struct{}是指没有字段的结构类型，常用于值可被忽略的场合
无论是其自身，还是作为元素类型，其长度都为0，但这不影响其作为实体存在
*/
func empty() {
	var a struct{}
	var b [10000]struct{}

	s := b[:]
	fmt.Println(unsafe.Sizeof(a), unsafe.Sizeof(b))
	fmt.Println(len(s), len(b))
	/*
		0 0
		10000 10000
	*/

	// 结束通知
	done := make(chan struct{})
	go func() {
		time.Sleep(time.Second)
		close(done)
	}()

	<-done

	// 随机选择 不需要值
	nul := struct{}{}
	users := make(map[int]struct{})
	for i := 0; i < 100; i++ {
		users[i] = nul
	}
	for k := range users {
		fmt.Println(k)
		break
	}
}

// 相等操作限制
// - 全部字段支持
// - 内存布局相同 但是类型不同 不能比较
// - 布局相同（字段名和标签）的匿名结构视为同一类型
func equals() {
	// map不支持比较 所以包含其的结构体data不支持比较
	type data struct {
		x int
		y map[string]int
	}

	d1 := data{x: 100}
	d2 := data{x: 100}
	//_ = d1 == d2 //data/struct.go:40:6: invalid operation: d1 == d2 (struct containing map[string]int cannot be compared)
	fmt.Println(d1, d2)

	// 类型不同 不能比较
	type data1 struct {
		x int
	}
	type data2 struct {
		x int
	}
	d3 := data1{x: 100}
	d4 := data2{x: 100}
	//_ = d3 == d4 // invalid operation: d3 == d4 (mismatched types data1 and data2)
	fmt.Println(d3, d4)

	// 匿名结构类型相同的前提是 字段名 字段类型 标签和排列顺序相等
	d5 := struct {
		x int
		s string
	}{100, "sdf"}

	var d6 struct {
		x int
		s string
	}
	d6.x = 100
	d6.s = "sdf"
	fmt.Println(d5 == d6) // true

	// 不能以多级指针字段访问成员，应该是做了啥保护在里边，以后看了源码再分析分析
	type user struct {
		name string
		age  int
	}
	p := &user{
		name: "sdfsdf",
		age:  21,
	}
	p.age++
	p.name = "sdfsd"

	//二级指针
	p2 := &p
	(*p2).age++ //这样可以
	//p2.age++ //p2.age undefined (type **user has no field or method age)
}

// 顺序init 必须包含全部字段值
// 命名init，不用全部字段，也无关顺序
// 建议使用命名方式，后边新增字段或者调整排列顺序也不受影响
func inits() {
	type node struct {
		id    int
		value string
		next  *node
	}

	// 顺序提供所有字段值或者全部忽略
	//n := node{1, "a"} // too few values in struct literal of type node

	// 按字段名init
	n := node{
		id:    2,
		value: "123",
	}
	fmt.Println(n)

	// 匿名结构
	user := struct {
		id   int
		name string
	}{
		id:   1,
		name: "123",
	}
	fmt.Printf("%+v\n", user)

	var color struct {
		r int
		g int
		b int
	}
	color.g = 1
	color.r = 2
	color.b = 3
	fmt.Printf("%+v\n", color)
	/*
		{id:1 name:123}
		{r:2 g:1 b:3}
	*/

	type file struct {
		name string
		attr struct {
			owner int
			perm  int
		}
	}

	f := file{
		name: "xxx", // 可以init成功
		// 因为缺少类型标签，无法直接init
		//attr: {
		//	owner: 1,
		//	perm:  0755,
		//},
		// data/struct.go:81:9: missing type in composite literal
	}
	// 这样可以成功
	f.attr.owner = 1
	f.attr.perm = 0755
	fmt.Println(f)
}
