package _type

import "fmt"

func output(int) (int, int, int)

func Main_2() {
	a, b, c := output(987654321)
	fmt.Println(a, b, c)
}
