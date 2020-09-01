package a

import "fmt"

func f() {
	arr := []string{}
	fmt.Println(arr[len(arr)]) // want "Likely cause index out of range"
}
