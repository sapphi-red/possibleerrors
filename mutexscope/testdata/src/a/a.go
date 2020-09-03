package a

import "sync"

var mu sync.Mutex

func f1() {
	mu.Lock() // want "Should Unlock inside function."
}

func f2() {
	mu.Lock()
	mu.Unlock()
}

func f3() {
	mu.Lock()
	defer mu.Unlock()
}

func f4() {
	mu.Lock() // want "Should Unlock inside function."
	if (true) {
		mu.Unlock()
	}
}

func f5() {
	mu.Lock() // want "Should Unlock inside function."
	for i := 0; i < 10; i++ {
		switch i % 3 {
		case 0: mu.Unlock()
		case 1: println(1)
		case 2: println(2)
		}
	}
}

func f6() {
	mu.Lock()
	for i := 0; i < 10; i++ {
		switch i % 3 {
		case 0: mu.Unlock()
		case 1: println(1)
		case 2: println(2)
		}
	}
	mu.Unlock()
}

func f7() {
	ok := true
	mu.Lock() // want "Should Unlock inside function."
	if ok {
		return
	}
	mu.Unlock()
}
