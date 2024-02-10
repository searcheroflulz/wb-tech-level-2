package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/* Применимость:
Данный паттерн может быть применим в тех ситуациях, когда нам необходимо создать сложный объект пошагово.
Также с помощью данного паттерна мы сможем создавать различные представления одного и того же объекта.
Строитель помогает разграничить процесс сборки и от представления объекта.

Плюсы:
Пошаговое конструирование сложного объекта;
Возможность переиспользования строителя для создания большего количество экземпляров;
Изолирование ненужного для конечного пользователя кода сборки конечного объекта;

Минусы:
Сложность написания кода для данного паттерна (большое количество дополнительных классов);
Большое количество дополнительных классов для строительства и получения нужного результата;

Примеры:
ORM (TypeORM, bun)
GUI (fyne)
*/

import "fmt"

// наш конечный объект
type CarProduct struct {
	Engine    string
	Wheels    int
	Seats     int
	Bluetooth bool
}

// интерфейс строителя
type CarBuilder interface {
	BuildEngine() CarBuilder
	BuildWheels() CarBuilder
	BuildSeats() CarBuilder
	BuildBluetooth() CarBuilder
	GetResult() CarProduct
}

// структура строителя
type ConcreteCarBuilder struct {
	car CarProduct
}

func (b *ConcreteCarBuilder) BuildEngine() CarBuilder {
	b.car.Engine = "V8"
	return b
}

func (b *ConcreteCarBuilder) BuildWheels() CarBuilder {
	b.car.Wheels = 4
	return b
}

func (b *ConcreteCarBuilder) BuildSeats() CarBuilder {
	b.car.Seats = 5
	return b
}

func (b *ConcreteCarBuilder) BuildBluetooth() CarBuilder {
	b.car.Bluetooth = true
	return b
}

func (b *ConcreteCarBuilder) GetResult() CarProduct {
	return b.car
}

// структура директора, которая управляет строительством нашего объекта
type Director struct {
	builder CarBuilder
}

func NewDirector(builder CarBuilder) *Director {
	return &Director{builder: builder}
}

// директор выдает сразу готовый результат с помощью строителя
func (d *Director) Construct() CarProduct {
	return d.builder.BuildEngine().BuildWheels().BuildSeats().BuildBluetooth().GetResult()
}

func main() {
	//создание строителя
	builder := &ConcreteCarBuilder{}
	//создание директора
	director := NewDirector(builder)

	//поручение директору строительства объекта
	car := director.Construct()
	fmt.Printf("Built Car: %+v\n", car)
}
