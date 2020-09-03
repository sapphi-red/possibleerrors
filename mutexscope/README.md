# mutexscope

mutexscope finds sync.Mutex which likely forgotten `.Unlock()`.

```go
var mu sync.Mutex

func f() {
  mu.Lock()
  if exists {
    return // forgot .Unlock()
  }
  mu.Unlock()
}
```
