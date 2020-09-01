package a

import "fmt"

func b2() {
	i := 0
	for ; i < 10; i++ {
		fmt.Println(i)
	}
	i = 0
	for ; i <= 10; i++ {
		fmt.Println(i)
	}
	i = 0
	for ; 10 > i; i++ {
		fmt.Println(i)
	}
	i = 0
	for ; 10 >= i; i++ {
		fmt.Println(i)
	}

	i = 0
	for ; i > 10; i-- {
		fmt.Println(i)
	}
	for ; i >= 10; i-- {
		fmt.Println(i)
	}
	for ; 10 < i; i-- {
		fmt.Println(i)
	}
	for ; 10 <= i; i-- {
		fmt.Println(i)
	}
}
