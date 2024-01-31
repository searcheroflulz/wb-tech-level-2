Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}

```

Ответ:
```
Программа выведет сначала 2 затем 1. Defer вызывается после вызова return. 
В функции test нам заранее известно возвращаемое значение (не произойдет копирования значения),
поэтому defer изменит именно его.
В функции anotherTest return скопирует x и вернет уже скопированное значение x.
```
