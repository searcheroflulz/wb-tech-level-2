package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/* Применимость:
Данный паттерн будет очень полезен, если нам нужно какую-то операцию для множества объектов одного интерфейса,
но при этом мы не можем изменять внутреннюю логику этих объектов. Также с помощью посетителя мы можем выполнять внутренние методы классов
в одном методе посетителя для большего удобства.

Плюсы:
Позволяет добавлять новые операции без изменения внутреннего кода класса;
Возможность использовать внутренние методы класса в одной функции;

Минусы:
Слжность понимания, если операций много, а структура элементов объемна;


Примеры:
Обработка синтаксического дерева, где каждый узел может быть представлен различными конструкциями языка;
Вычисление статистики посещаемости сайтов, где могут быть различные типы страниц;
Итераторы для обхода структуры и обработки ее данных;
*/

// интерфейс посетителя
type Visitor interface {
	VisitCircle(circle *Circle)
	VisitRectangle(rectangle *Rectangle)
}

// интерфейс фигуры
type Shape interface {
	Accept(visitor Visitor)
}

// реализация круга
type Circle struct {
	Radius float64
}

// метод круга, который принимает посетителя
func (c *Circle) Accept(visitor Visitor) {
	visitor.VisitCircle(c)
}

// реализация прямоугольника
type Rectangle struct {
	Width  float64
	Height float64
}

// метод прямоугольника, который принимает посетителя
func (r *Rectangle) Accept(visitor Visitor) {
	visitor.VisitRectangle(r)
}

// реализация интерфейса посетитель для вычисления площади
type AreaVisitor struct {
	TotalArea float64
}

// посещаем круг, чтобы вычислеть его площадь
func (a *AreaVisitor) VisitCircle(circle *Circle) {
	area := 3.14 * circle.Radius * circle.Radius
	a.TotalArea += area
}

// посещаем прямоугольник, чтобы вычислить его площадь
func (a *AreaVisitor) VisitRectangle(rectangle *Rectangle) {
	area := rectangle.Width * rectangle.Height
	a.TotalArea += area
}

// реализация интерфейса посетитель для вычисления периметра
type PerimeterVisitor struct {
	TotalPerimeter float64
}

// посещаем круг для вычисление периметра
func (p *PerimeterVisitor) VisitCircle(circle *Circle) {
	perimeter := 2 * 3.14 * circle.Radius
	p.TotalPerimeter += perimeter
}

// посещаем прямоугольник, чтобы вычислить периметр
func (p *PerimeterVisitor) VisitRectangle(rectangle *Rectangle) {
	perimeter := 2 * (rectangle.Width + rectangle.Height)
	p.TotalPerimeter += perimeter
}

func main() {
	circle := &Circle{Radius: 5}
	rectangle := &Rectangle{Width: 4, Height: 6}

	areaVisitor := &AreaVisitor{}
	perimeterVisitor := &PerimeterVisitor{}

	circle.Accept(areaVisitor)
	rectangle.Accept(areaVisitor)

	circle.Accept(perimeterVisitor)
	rectangle.Accept(perimeterVisitor)

	fmt.Printf("Total Area: %f\n", areaVisitor.TotalArea)
	fmt.Printf("Total Perimeter: %f\n", perimeterVisitor.TotalPerimeter)
}
