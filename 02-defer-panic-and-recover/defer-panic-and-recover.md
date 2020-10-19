# defer-panic-and-recover

## 一、defer

``defer``会即将函数存储到一个列表中，这个列表里的函数会在外层函数返回的时候执行（``不知道是在外层函数返回值之前还是之后执行``）。``defer``通常用来执行清理动作，参看下面的代码。
```Golang
func CopyFile(destName, srcName string) (written int64, err error) {
  src, err := os.OpenFile(destName)
  if err != nil {
    return 0, err
  }
  defer src.Close()

  dest, err := os.OpenFile(srcName)
  if err != nil {
    return 0, err
  }
  defer dest.Close()

  return io.Copy(dest, src)
}
```

### ``defer``的3个规则

1. 程序运行到defer表达式的时候，defer函数的参数会被计算
```Golang
func a() int {
  i := 1
  defer fmt.Println(i)  // 打印出1
  i++
  return i
}
```

2. defer函数的执行顺序，是先进后出的顺序
```Golang
func b() {
  for i := 0; i < 4; i++ {
    defer fmt.Println(i)
  }
}
// 打印的顺序是 3, 2, 1, 0
```

3. ``defer``函数可以读取并分配返回值的值
```Golang
func c() (int i) {
  defer func() {
    i++
  }()
  return 1
  // 返回值是2
}
```
这个机制的好处：可以很方便的修改返回值中的err(就目前而言，没get到这个好处)

### ``defer``函数和return执行的顺序

从``defer``的第三条规则可以看出，``defer``函数会在return之前执行。

## ``panic``和``recover``

``panic``是内置函数，会组织函数继续向下执行，同时会开始``panicking``。函数调用``panic``，已经加载的``defer``函数还会正常执行，同时函数会返回给调用者（不是返回指定的返回值）。对于函数的调用者，其行为也像执行了``panic``一样，这样一直向上，直到程序崩溃。

``recover``是一个内置的函数，会获取``panicking goroutine``的控制权（能继续执行程序，不让程序崩溃）。在程序正常运行时，``recover``返回nil并且不会有任何其他的影响。如果当前的``goroutine``是``panicking``，``recover``会捕获``panic``并且会恢复正常运行。

```Golang
package main

import "fmt"

func main() {
    f()
    fmt.Println("Returned normally from f.")
}

func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    fmt.Println("Calling g.")
    g(0)
    fmt.Println("Returned normally from g.")
}

func g(i int) {
    if i > 3 {
        fmt.Println("Panicking!")
        panic(fmt.Sprintf("%v", i))
    }
    defer fmt.Println("Defer in g", i)
    fmt.Println("Printing in g", i)
    g(i + 1)
}
```
输出如下
```Text
Calling g.
Printing in g 0
Printing in g 1
Printing in g 2
Printing in g 3
Panicking!
Defer in g 3
Defer in g 2
Defer in g 1
Defer in g 0
Recovered in f 4
Returned normally from f.
```






