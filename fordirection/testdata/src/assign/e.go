package a

import "fmt"

func e() {
	for i := 0; i > 10; i -= -2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}

	const skip = 2
	for i := 0; i < 10; i -= skip { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}

	const skip2 = -2
	for i := 0; i > 10; i -= skip2 { // want "Loop direction seems to be wrong."
		fmt.Println(i)
	}
}
