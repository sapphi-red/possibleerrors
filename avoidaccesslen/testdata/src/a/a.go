package a

import "fmt"

func f() {
	slice := []string{}
	fmt.Println(slice[len(slice)]) // want "Will occur index out of range"

	slice2 := slice
	fmt.Println(slice2[len(slice2)]) // want "Will occur index out of range"

	fmt.Println(slice[len(slice)-1])
	fmt.Println(slice[len(slice)+1]) // want "Will occur index out of range"
	fmt.Println(slice[len(slice)+2]) // want "Will occur index out of range"

	fmt.Println(slice[-1+len(slice)])
	fmt.Println(slice[1+len(slice)])   // want "Will occur index out of range"
	fmt.Println(slice[3-2+len(slice)]) // want "Will occur index out of range"
	fmt.Println(slice[3+4+len(slice)]) // want "Will occur index out of range"

	sliceA := []string{}
	fmt.Println(sliceA[len(slice)])
	fmt.Println(sliceA[len(slice)+1])
}
