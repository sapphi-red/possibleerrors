package b

import "sync"

type MapMutex struct {
	sync.Mutex
	mp map[interface{}]interface{}
}

func (m *MapMutex) Get(key interface{}) (interface{}, bool) {
	m.Lock()
	v, ok := m.mp[key]
	m.Unlock()
	return v, ok
}
func (m *MapMutex) Set(key interface{}, value interface{}) {
	m.Mutex.Lock()
	m.mp[key] = value
	m.Unlock()
}

func f() {
	mm := MapMutex{}
	mm.Lock() // want "Should Unlock inside function."
}
