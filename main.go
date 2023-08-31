package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime/trace"
	"unsafe"
	"yuhen/data"
	"yuhen/fun"
	_type "yuhen/type"
)

var a = 898

//go:noinline
func stringParam(s string) {}

func main() {

	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	s := new([]int)
	m := *new(map[string]int)
	c := *new(chan int)

	fmt.Println(s, reflect.TypeOf(s), unsafe.Sizeof(s), unsafe.Sizeof(s))
	fmt.Println(m, reflect.TypeOf(m), unsafe.Sizeof(m), unsafe.Sizeof(m))
	fmt.Println(c, reflect.TypeOf(c), unsafe.Sizeof(c), unsafe.Sizeof(c))

	newTest()
	_type.Main_1()
	_type.Main_2()
	_type.Main_output()
	sliceTest()
	fun.Anonymous()
	fun.MainDefer()
	fun.MyContextError()
	data.MainString()
	data.Array()
	data.Slice()
	data.Map()
	data.Struct()
	data.Pointer()
	data.Method()
	data.Interface()
	data.Generic()
	data.Concurrency()
	data.Channel()
}

func sliceTest() {
	var x = "abcc"
	stringParam(x)
}

func newTest() {
	p := new(int)
	*p = 100
}
