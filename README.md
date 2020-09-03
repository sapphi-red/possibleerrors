# possibleerrors

Golang linters for finding code which is likely a logic error.

- [fordirection](./fordirection): `for i := 0; i < 10; i--`
- [avoidaccesslen](./avoidaccesslen): `slice[len(slice)]`
- [mutexscope](./mutexscope): forgotten `.Unlock()`.
