# fordirection

fordirection finds for-loops which likely has a wrong direction.
```go
for i := 0; i < 10; i-- { // this should be `i++`
  fmt.Println(i)
}
```
