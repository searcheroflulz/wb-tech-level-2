package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/* Применимость:
Данный паттерн может быть отличным решением, когда нам нужно сделать программу со сложной подсистемой максимально простой в использовании для конечного пользователя.
То есть для пользователя есть лишь необходимые для него способы взаимодействия, а вся сложная подсистема со множеством различных вызовов и зависимостей будет спрятана.

Плюсы:
Упрощение взаимодействия с пользователем;
Снижение количества зависимостей;
Упрощение поддержки основного кода, скрытого за фасадом;

Минусы:
Скрытие некоторые важных для пользователя деталей и методов;
Отсутствие гибкости для разных типов пользователей;

Примеры:
Различные фреймворки для разработки веб-приложений (React, Angular, Vue.js);
Различные библиотеки для работы со внешними сервисами (pgx для Postgres, go-redis для Redis);
*/

// система инвентаря
type Inventory struct {
}

func (i *Inventory) checkInventory(productId string) bool {
	fmt.Printf("Checking inventory for product %s\n", productId)
	//проверка наличия продукта в инвенторе
	return true
}

// система платежей
type PaymentProcessor struct {
}

func (p *PaymentProcessor) makePayment(amount float64) {
	fmt.Printf("Processing payment of $%f\n", amount)
	//обработка платежей
}

// система заказов
type OrderFulfillment struct {
}

func (o *OrderFulfillment) fulfillOrder(productId string, quantity int) {
	fmt.Printf("Fulfilling order for product %s, quantity %d\n", productId, quantity)
	//обработка выполнения заказа
}

// фасад
type OrderFacade struct {
	Inventory        *Inventory
	PaymentProcessor *PaymentProcessor
	OrderFulfillment *OrderFulfillment
}

func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		Inventory:        &Inventory{},
		PaymentProcessor: &PaymentProcessor{},
		OrderFulfillment: &OrderFulfillment{},
	}
}

// создание нового заказа через фасад
func (of *OrderFacade) PlaceOrder(productId string, quantity int, amount float64) {
	if of.Inventory.checkInventory(productId) {
		of.PaymentProcessor.makePayment(amount)
		of.OrderFulfillment.fulfillOrder(productId, quantity)
		fmt.Println("Order placed successfully!")
	} else {
		fmt.Println("Failed to place order due to insufficient inventory.")
	}
}

func main() {
	//размещение заказа через один метод фасада
	orderFacade := NewOrderFacade()
	orderFacade.PlaceOrder("ABC123", 2, 50.0)
}
