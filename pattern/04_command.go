package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/* Применимость:
Данный паттерн может быть применим, когда нам нужно удобное управление различными операциями в нашей программе.
Нам нужно ставить выполняемые операции в очередь, иметь поддержку их отмены.

Плюсы:
Возможность хранения истории операций;
Поддержка отмены операций;
Возможность комбинирования нескольких простых комманд в одну более сложную;
Отделение отправителя команды от получателя;

Минусы:
Увеличение числа классов;
Сложность дальнейшего дебаггинга;

Примеры:
Графические редакторы (множество различных встроенных программ для построения фигур, изменения цвета...);
Транзакции для баз данных (возможность построение сложной операции из простых и их отмена);
Управление умными устройствами (множество команд и операций над различными устройствами);
*/

import "fmt"

// интерфейс команды с методом execute
type Command interface {
	Execute()
}

// объект получатель
type Light struct {
	isOn bool
}

// методы самого объекта получателя
func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Println("Light is ON")
}

func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Println("Light is OFF")
}

// отдельная команда для работы с объектом получателем
type TurnOnCommand struct {
	light *Light
}

// выполнение действия над объектом получателем
func (c *TurnOnCommand) Execute() {
	c.light.TurnOn()
}

// отдельная команда
type TurnOffCommand struct {
	light *Light
}

func (c *TurnOffCommand) Execute() {
	c.light.TurnOff()
}

// объект отправитель команд
type RemoteControl struct {
	command Command
}

// отправка команды
func (rc *RemoteControl) PressButton() {
	rc.command.Execute()
}

func main() {
	light := &Light{}
	turnOnCommand := &TurnOnCommand{light: light}
	turnOffCommand := &TurnOffCommand{light: light}

	remoteControl := &RemoteControl{}

	remoteControl.command = turnOnCommand
	remoteControl.PressButton()

	remoteControl.command = turnOffCommand
	remoteControl.PressButton()
}
