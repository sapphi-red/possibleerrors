package a

import "fmt"

func c() {
	for ;; {
		fmt.Println("")
		break
	}

	for i, j := 0, 0; j < 10; i++ {
		j = i
		fmt.Println(j)
	}
}
