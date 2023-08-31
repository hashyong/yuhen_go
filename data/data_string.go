package data

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"
)

/*
+------------+         +---+---+---+---+---+
|            |         |   |   |   |   |   |
|   pointer.go +--------> | h | e | l | l | o |
|            |         |   |   |   |   |   |
+------------+         +---+---+---+---+---+
|            |
|   len=5    |
|            |         [...]byte, UTF-8
+------------+
    header
*/

/*
runtime/string.go

type stringStruct struct {
	str unsafe.Pointer
	len int
}

- 编码支持utf-8，无null结尾，默认值 ""

- 使用 raw string，定义原始字符串
- 支持各种运算符

- 索引访问字节数组（非字符），不能获取元素地址
- 切片返回子串，依旧指向原数组
- 内置函数len 返回字节数组的长度
*/

/*
utf8.RuneCount 和 utf8.RuneCountInString 都是 Go 语言中 utf8 包提供的用于计算字符串中 Unicode 字符数量的函数。

Unicode 是一种字符编码标准，它定义了世界上几乎所有的符号，包括字母、数字、标点符号和表情符号等。在 Go 语言中，Unicode 字符串是以 rune 类型表示的。

utf8.RuneCount 函数接收一个字节数组参数，返回该字节数组中包含的 Unicode 字符数量

utf8.RuneCountInString 函数接收一个字符串参数，返回该字符串中包含的 Unicode 字符数量
*/

func MainString() {
	s := "YB哈哈哈\x61\142\u0012"

	byteString := []byte(s)
	rs := []rune(s)

	fmt.Printf("%X, %d\n", s, len(s))
	fmt.Printf("%X, %d\n", byteString, utf8.RuneCount(byteString))
	fmt.Printf("%X, %d\n", rs, utf8.RuneCountInString(s))

	s1 := "abcdef"
	fmt.Printf("%X %d\n", s1, len(s1)) // 616263646566,6字节, 十六进制下一个数字是4bit，2个是8bit为1byte 合理
	fmt.Printf("%X\n", s1[1])          // 62
	fmt.Println(&s1)
	// fmt.Println((&s1), (&s1)[1])       // Invalid operation: '(&s1)[1]' (type '*string' does not support indexing)
	//索引访问字节数组（非字符），不能获取元素地址, 原来还有这个限制在里边

	var s2 string
	fmt.Println(s2 == "") //true
	//fmt.Println(s2 == nil) //Cannot convert 'nil' to type 'string'

	mainString1()
	loop()
	opt()
	convert()
	performance()
}

func mainString1() {
	// raw string 内的转义 换行 前置空格 注释等 都视作内容 不予处理
	s := `sdf
sdfsf sdf
  sfsdfsd
    sd
fs
f
sd`
	fmt.Println(s)

	s1 := s + " aa"
	fmt.Println(s1)

	// 以加号链接字面量的时候 注意操作符的位置
	s1 = s1 +
		"sdfsfsdfsf"
	fmt.Println(s1)

	// 字串内部指针指向原来的字节数组
	s2 := "hello, world"
	s3 := s2[:4]

	p1 := (*reflect.StringHeader)(unsafe.Pointer(&s2))
	p2 := (*reflect.StringHeader)(unsafe.Pointer(&s3))

	fmt.Printf("%#v\n%#v\n", p1, p2)
	/*&reflect.StringHeader{Data:0x571af0, Len:12}
	&reflect.StringHeader{Data:0x571af0, Len:4}*/
}

func loop() {
	s := "王者荣耀"

	//byte 按照字节来遍历
	for i := 0; i < len(s); i++ {
		fmt.Printf("%d:%v\n", i, s[i])
	}
	/*
		0:231
		1:142
		2:139
		3:232
		4:128
		5:133
		6:232
		7:141
		8:163
		9:232
		10:128
		11:128
	*/
	for i, c := range s {
		fmt.Printf("%d:%c\n", i, c)
	}
	/*
		0:王
		3:者
		6:荣
		9:耀
	*/
}

// 可以作为append 和 copy相关的参数
func opt() {
	s := "王者"

	bs := make([]byte, 0)
	// 语法糖，假如是append slice 需要再后边跟着...
	bs = append(bs, "abc"...)
	bs = append(bs, s...)

	// 一个中文要占3字节
	buf := make([]byte, 9)
	copy(buf, "abc")
	copy(buf[3:], s)

	fmt.Printf("%s\n", bs)
	fmt.Printf("%s\n", buf)
	// 标准库相关的操作
	/*
		    - bytes: 字节切片
			- fmt: 格式化
			- strconv: 转换
			- strings: 函数
			- text: 模板
			- unicode: 码点
	*/
}

// 可以在rune byte string 之间转换
// 单引号字符字面量是 rune类型，代码unicode字符
func convert() {
	var r rune = '我'

	var s string = string(r)
	var b byte = byte(r)
	var s2 string = string(b)
	var r2 rune = rune(b)

	fmt.Printf("%c, %U\n", r, r)
	fmt.Printf("%s, %X, %X, %X\n", s, b, s2, r2)

	// 要修改字符串，必须转换为可变类型，[]rune or []byte, 修改完之后再转回来
	// 但是不管如何转换, 都需要重新分配内存, 并复制数据, 字面量在二进制文件里边被分配为 RODATA，只读，其内存不能修改

	s3 := strings.Repeat("a", 1<<10)
	// 分配内存和复制
	bs := []byte(s3)
	bs[1] = 'B'

	// 分配内存和复制
	s4 := string(bs)

	hs3 := (*reflect.StringHeader)(unsafe.Pointer(&s3))
	hbs := (*reflect.StringHeader)(unsafe.Pointer(&bs))
	hs4 := (*reflect.StringHeader)(unsafe.Pointer(&s4))

	fmt.Printf("%#v\n%#v\n%#v\n", hs3, hbs, hs4)
	/*&reflect.StringHeader{Data:0xc0000e0400, Len:1024}
	&reflect.StringHeader{Data:0xc0000e0800, Len:1024}
	&reflect.StringHeader{Data:0xc0000e0c00, Len:1024}*/

	// 验证下编码是否正确
	s5 := "雨痕"
	s6 := string(s5[0:2] + s5[4:]) //非法拼接

	fmt.Printf("%X, %X\n", s5, s6)
	// go语言果然还是NB的 可以判断出来是否是utf编码
	fmt.Println(utf8.ValidString(s6))
}

// 核心是使用不安全的方法进行转换，改善性能
// 不动底层数组，直接构建 string 或者 slice头
// 如果修改要注意安全

// slice{data, len, cap} -> string{data, len}
// slice{string.data, string,len, string.len}

var S = strings.Repeat("a", 100)

func normalConv() bool {
	// 底层转换会调用 runtime.stringtoslicebyte runtime.slicebytetostring 会引发mallocgc，memmove等等操作, 比较花费时间和占用内存资源
	b := []byte(S)
	s2 := string(b)
	return s2 == S
}

func unsafeConv() bool {
	// []byte(S)
	b := unsafe.Slice(unsafe.StringData(S), len(S))

	// string(b)
	s2 := unsafe.String(unsafe.SliceData(b), len(b))

	return s2 == S
}

func performance() {
	fmt.Println(normalConv(), unsafeConv())
}

// 性能差异非常大，详细见 test benchmark的结果

//---
// 字符串的构造也很容易造成性能问题
// 加法操作符拼接字符串，每次都需要重新分配内存和复制数据
// 该井方法是预先分配内存，然后一次性返回

const N = 1000
const C = "a"

var S1 = strings.Repeat(C, N)

// 编译器对于字面+号拼接会做优化
// 还可以用fmt.sprintf text/template
// 字符串某些看似简单的操作，都可能引发内存分配和复制
func concat() bool {
	var s2 string

	for i := 0; i < N; i++ {
		s2 += C
	}

	return s2 == S1
}

func join() bool {
	b := make([]string, N)
	for i := 0; i < N; i++ {
		b[i] = C
	}

	return strings.Join(b, "") == S1
}

func buffer() bool {
	var b bytes.Buffer
	b.Grow(N)

	for i := 0; i < N; i++ {
		b.WriteString(C)
	}

	return b.String() == S1
}
