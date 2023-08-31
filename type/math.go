package _type

import "fmt"

func add(a, b int) int // 汇编函数声明
func sub(a, b int) int // 汇编函数声明
func mul(a, b int) int // 汇编函数声明

func Main_1() {
	fmt.Println(add(10, 11))
	fmt.Println(sub(99, 15))
	fmt.Println(mul(11, 12))
}
