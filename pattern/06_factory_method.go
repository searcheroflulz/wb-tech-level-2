package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/* Применимость:
Данный паттерн применяется в тех ситуациях, когда создание объекта требует от нас выполнения сложной логики или подготовительных действий.
И также когда нам нужно дать возможность подклассам выбирать тип создаваемого объекта.

Плюсы:
Разделение отвественности (наш код для создания объектов в отдельном от всего остального кода месте);
Гибкое добавление новых типов объектов без изменения существующего кода;
Возможность использования множества различных фабрик;

Минусы:
Усложнение кода (множество дополнительных структур и интерфейсов для разных создаваемых объектов);

Примеры:
Игровые движки (создание экземпляров игровых объектов (персонажи, предметы, оружие);
Библиотеки для работы с базами данных (создание DAO);
*/

import "fmt"

// интерфейс транспорта
type ITransport interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

type Transport struct {
	name  string
	power int
}

func (c *Transport) setName(name string) {
	c.name = name
}

func (c *Transport) setPower(power int) {
	c.power = power
}

func (c *Transport) getName() string {
	return c.name
}

func (c *Transport) getPower() int {
	return c.power
}

type car struct {
	Transport
}

func newCar() ITransport {
	return &car{
		Transport: Transport{
			name:  "car",
			power: 150,
		},
	}
}

type motorcycle struct {
	Transport
}

func newMotorcycle() ITransport {
	return &motorcycle{
		Transport: Transport{
			name:  "motorcycel",
			power: 80,
		},
	}
}

func getTransport(transportType string) (ITransport, error) {
	if transportType == "car" {
		return newCar(), nil
	}
	if transportType == "motorcycle" {
		return newMotorcycle(), nil
	}
	return nil, fmt.Errorf("wrong type")
}

func main() {
	audi, _ := getTransport("car")
	suzuki, _ := getTransport("motorcycle")

	fmt.Println(audi.getName())
	fmt.Println(suzuki.getName())
}
