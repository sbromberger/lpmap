# lpmap: linear probing hash tables in Go

This is a generic implementation of linear probing hash tables in Go.

Example:

```go
type MyKey uint64

func (k MyKey) hash uint64 {
    return k
}
    
func main() {
    h := lpmap.New[MyKey, string](10, 0.5)  // capacity of 10, with a max load factor of 0.5
    h.Set(9, "nine")
    
    val, found := h.Get(MyKey(9))     // "nine", true
    val, found = h.Get(MyKey(1))      // nil, false
    h.Set(MyKey(9)) = "four"  // replace value
```