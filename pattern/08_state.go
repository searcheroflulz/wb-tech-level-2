package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/* Применимость:
Данный паттерн удобно использовать, когда нам нужно, чтобы объект менял свое поведение в зависимости от внутреннего состояния.
Таким паттерном мы можем заменить большое количество условных операторов, которые выбириали бы поведение объекта по текущим значениям полей структуры.

Плюсы:
Избавление от большого количества условных конструкций;
Возможность удобного увеличения количества состояний объекта;
Инкапсулирование состояний в отдельной структуре;

Минусы:
Усложнение структуры программы из-за дополнительных переходов от одного состояния к другому;
Увеличение количества объектов-состояний, переходов к ним и вызовов их методов;

Примеры:
Игровые движки (поведение игровых персонажей - бег, атака, защита);
Конечный автомат (искусственный интеллект в играх);
*/

import "fmt"

// состояние объекта
type OrderState interface {
	HandleOrder(order *Order)
}

// контекст нашего заказа
type Order struct {
	state OrderState
}

func NewOrder() *Order {
	return &Order{state: &CreatedState{}}
}

// изменение состояния объекта
func (o *Order) SetState(state OrderState) {
	o.state = state
}

// базовый метод заказа для его продвижения дальше до состояния "доставлен"
func (o *Order) HandleOrder() {
	o.state.HandleOrder(o)
}

// первое состояние заказа - создан
type CreatedState struct{}

func (cs *CreatedState) HandleOrder(order *Order) {
	fmt.Println("Order is created. Waiting for payment.")
	order.SetState(&PaidState{})
}

// второе состояние заказа - оплачен
type PaidState struct{}

func (ps *PaidState) HandleOrder(order *Order) {
	fmt.Println("Order is paid. Preparing for shipping.")
	order.SetState(&ShippedState{})
}

// третье состояние заказа - отправлен
type ShippedState struct{}

func (ss *ShippedState) HandleOrder(order *Order) {
	fmt.Println("Order is shipped. Waiting for delivery.")
	order.SetState(&DeliveredState{})
}

// последнее состояние заказа - доставлен
type DeliveredState struct{}

func (ds *DeliveredState) HandleOrder(order *Order) {
	fmt.Println("Order is delivered. Thank you for your purchase!")
}

func main() {
	order := NewOrder()

	order.HandleOrder()
	order.HandleOrder()
	order.HandleOrder()
	order.HandleOrder()
}
