package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

func or(channels ...<-chan interface{}) <-chan interface{} {
	//создаем waitGroup
	var wg sync.WaitGroup
	//наш объединяющий канал
	result := make(chan interface{})

	//функция, в которой считываются все значения из канала и отправляются в наш результирующий канал
	closeChan := func(ch <-chan interface{}) {
		for val := range ch {
			result <- val
		}
		//как только прочитали все значения, мы уменьшаем счетчик waitGroup
		wg.Done()
	}
	//количество прослушиваемых каналов
	wg.Add(len(channels))

	//слушаем каналы в цикле
	for _, channel := range channels {
		go closeChan(channel)
	}

	//как только все каналы отдали свои значения, мы закрываем наш главный канал
	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("Done after %v\n", time.Since(start))
}
