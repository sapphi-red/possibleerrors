package a

import "fmt"

func a() {
	for i := 0; i < 10; i-- { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i <= 10; i-- { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; 10 > i; i-- { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; 10 >= i; i-- { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}

	for i := 0; i > 10; i++ { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i >= 10; i++ { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; 10 < i; i++ { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; 10 <= i; i++ { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
}
