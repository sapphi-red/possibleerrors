package a

import "fmt"

func f() {
	slice := []string{}
	fmt.Println(slice[len(slice)]) // want "Will occur index out of range"

	slice2 := slice
	fmt.Println(slice2[len(slice2)]) // want "Will occur index out of range"

	slice3 := []string{}
	fmt.Println(slice3[len(slice)])
}
