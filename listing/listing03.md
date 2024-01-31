Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}

```

Ответ:
```
Программа выведет nil и затем false. error, который возвращает функция не будет равна nil,
так как нил равно только значение переменной err, а ее тип равен *os.PathError.
То есть мы сравниваем типизированный nil с nil без типа. Поэтому они не равны.
```
