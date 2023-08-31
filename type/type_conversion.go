package _type

import "fmt"

func convertType() {
	a := 10
	b := byte(a)
	c := a + int(b) // 类型统一才可以
	fmt.Println(c)

	x := 100
	d := make(chan int)
	// var b bool = x
	// 'b' redeclared in this block. Unused variable 'b'.
	// Cannot use 'x' (type int) as the type bool

	// if x {
	// }
	// The non-bool value 'x' (type int) used as a condition
	fmt.Println(x)

	// p := *int(&x)
	fmt.Println((*int)(&x)) // 如果没有括号, *(int(&x))
	fmt.Println((<-chan int)(d))
	//Invalid indirect of 'int(&x)' (type 'int'). Cannot convert an expression of the type '*int' to the type 'int'.
}
