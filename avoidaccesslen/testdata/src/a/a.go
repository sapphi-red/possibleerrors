package a

import "fmt"

func f() {
	arr := []string{}
	fmt.Println(arr[len(arr)]) // want "Will occur index out of range"

	arr2 := arr
	fmt.Println(arr2[len(arr2)]) // want "Will occur index out of range"

	arr3 := []string{}
	fmt.Println(arr3[len(arr)])
}
