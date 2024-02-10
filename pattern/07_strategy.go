package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/* Применимость:
Данный паттерн программирования можно применить, когда нам нужно выбрать какой-то алгоритм из множества на лету, то есть у нас есть понимание,
когда какой-то из алгоритмов будет правильнее использовать. Этот паттерн поможет реализовать выбор и смену алгоритмов внутри одного объекта.
Это также помогает нам держать все наши алгоритмы закрытыми для других функций и объектов.

Плюсы:
Выбор нужного нам алгоритма на лету;
Закрытие кода алгоритмов от других сущностей;
Изолирование алгоритмов от клиента;

Минусы:
Усложнение программы из-за дополнительного числа интерфейсов и зависимостей;
Клиенту необходимо понимание того, когда ему стоит менять алгоритмы;

Примеры:
Различные системы поиска (двоичный, полный перебор);
Сжатие данных (разные стратегии сжатия файлов);
Игровые движки (разная физика, отрисовка объектов);
*/

import (
	"fmt"
)

// интерфейс семейства алгоритмов сортировки
type SortStrategy interface {
	Sort([]int) []int
}

// алгоритм сортировки пузырьком, реализующий наш основной интерфейс
type BubbleSortStrategy struct{}

func (bs *BubbleSortStrategy) Sort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// алгоритм сортировки слиянием, реализующий наш основной интерфейс
type MergeSortStrategy struct{}

func (ms *MergeSortStrategy) Sort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := ms.Sort(arr[:mid])
	right := ms.Sort(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	for len(left) > 0 || len(right) > 0 {
		if len(left) == 0 {
			return append(result, right...)
		}
		if len(right) == 0 {
			return append(result, left...)
		}
		if left[0] <= right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}
	return result
}

// наша структура, в которой мы сможем менять алгоритмы на лету
type Context struct {
	strategy SortStrategy
}

func NewContext(strategy SortStrategy) *Context {
	return &Context{strategy: strategy}
}

func (c *Context) SetStrategy(strategy SortStrategy) {
	c.strategy = strategy
}

func (c *Context) Sort(arr []int) []int {
	return c.strategy.Sort(arr)
}

func main() {
	arr := []int{5, 3, 8, 2, 1, 4, 6, 7}
	context := NewContext(&BubbleSortStrategy{})
	sorted := context.Sort(arr)
	fmt.Println("Sorted using BubbleSortStrategy:", sorted)

	context.SetStrategy(&MergeSortStrategy{})
	sorted = context.Sort(arr)
	fmt.Println("Sorted using MergeSortStrategy:", sorted)
}
