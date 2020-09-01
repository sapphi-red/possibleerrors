package a

import "fmt"

func a2() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	for i := 0; i <= 10; i++ {
		fmt.Println(i)
	}
	for i := 0; 10 > i; i++ {
		fmt.Println(i)
	}
	for i := 0; 10 >= i; i++ {
		fmt.Println(i)
	}

	for i := 0; i > 10; i-- {
		fmt.Println(i)
	}
	for i := 0; i >= 10; i-- {
		fmt.Println(i)
	}
	for i := 0; 10 < i; i-- {
		fmt.Println(i)
	}
	for i := 0; 10 <= i; i-- {
		fmt.Println(i)
	}
}
