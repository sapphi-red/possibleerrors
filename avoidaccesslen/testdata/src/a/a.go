package a

import "fmt"

func f() {
	arr := []string{}
	fmt.Println(arr[len(arr)]) // want "Likely cause index out of range"

	arr2 := arr
	fmt.Println(arr2[len(arr2)]) // want "Likely cause index out of range"

	arr3 := []string{}
	fmt.Println(arr3[len(arr)])
}
