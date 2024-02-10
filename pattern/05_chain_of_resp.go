package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/* Применимость:
Данный паттерн будет особенно полезен, когда у нас есть несколько объектов, которые могут обработать запрос, и наш порядок обработки не фиксирован.
Также он будет полезен и если нам нужен строгий порядок обработки запроса обработчиками или если порядок должен задаваться динамически.

Плюсы:
Возможность легко добавить новые обработчики или изменить текущий порядок без изменений клиентского кода;
Уменьшение количества зависимостей в работе обработчиков, так как они не знают о существовании друг друга;
Отделение отправителя запроса от получателя;

Минусы:
Нет гарантий, что запрос будет обработан хотя бы одним из обработчиков (запрос просто может быть потярян);

Примеры:
Веб-фреймворки (различные middleware);
Веб-серверы (обработчики HTTP-запросов);

*/

import (
	"fmt"
)

// интерфейс обработчика
type ErrorHandler interface {
	HandleError(errorCode int) bool
}

// структура конкретного обработчика
type ConcreteHandler struct {
	nextHandler ErrorHandler
	errorCodes  map[int]bool
}

// создание обработчика
func NewConcreteHandler(errorCodes []int) *ConcreteHandler {
	handler := &ConcreteHandler{
		errorCodes: make(map[int]bool),
	}

	for _, code := range errorCodes {
		handler.errorCodes[code] = true
	}

	return handler
}

// добавление следующего обработчика в нашу цепочку
func (ch *ConcreteHandler) SetNextHandler(nextHandler ErrorHandler) {
	ch.nextHandler = nextHandler
}

// обработка запроса
func (ch *ConcreteHandler) HandleError(errorCode int) bool {
	if ch.errorCodes[errorCode] {
		fmt.Printf("Handled error with code %d\n", errorCode)
		return true
	}

	if ch.nextHandler != nil {
		return ch.nextHandler.HandleError(errorCode)
	}

	fmt.Printf("Error with code %d is unhandled\n", errorCode)
	return false
}

func main() {
	handler1 := NewConcreteHandler([]int{404, 500})
	handler2 := NewConcreteHandler([]int{401, 403})
	handler3 := NewConcreteHandler([]int{400})

	handler1.SetNextHandler(handler2)
	handler2.SetNextHandler(handler3)

	//пример обработки нескольких ошибок
	errorCodes := []int{200, 404, 401, 500, 400, 403}
	for _, code := range errorCodes {
		handler1.HandleError(code)
	}
}
