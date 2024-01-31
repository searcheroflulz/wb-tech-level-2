Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

```

Ответ:
```
Программа выведет error. Это произойдет по той причине, что мы имеем interface,
в котором значение равно nil, но также имеется тип, который не равен nil а равен *customError. 
Поэтому при сравнении с обычным nil, у которого нет типа, мы получаем неравенство.
```
