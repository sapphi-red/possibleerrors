package a

import "fmt"

func b() {
	i := 0
	for ; i < 10; i -= 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	i = 0
	for ; i <= 10; i -= 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	i = 0
	for ; 10 > i; i -= 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	i = 0
	for ; 10 >= i; i -= 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}

	i = 0
	for ; i > 10; i += 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for ; i >= 10; i += 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for ; 10 < i; i += 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for ; 10 <= i; i += 2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
}
