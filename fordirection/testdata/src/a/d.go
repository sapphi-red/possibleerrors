package a

import "fmt"

func d() {
	for i := 0; i < i; i++ { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i <= i; i++ { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i > i; i++ { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i >= i; i++ { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}

	for i := 0; i < i; i-- { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i <= i; i-- { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i > i; i-- { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
	for i := 0; i >= i; i-- { // want "loop direction seems to be wrong."
		fmt.Println(i)
	}
}
