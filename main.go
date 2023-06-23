package main

import (
	"fmt"
	"sync"
)




func main(){
/*
	t2()
	t3()
	t4()
 */
	t1()

}

// У нас отсутствует  выход из 2го for
// select у нас получается блокирующий:
		// можно обеспечть разблокировку закрытием канала
		// или чтением из закрывающего канала

		// решив проблемы мы должны подумать когда закрывать канал или когда писать в закрывающий канал - назовём это механизмом остановки
		// итак нам нужен механизм остановки чтения (из 2  самых очевидных вариантов) - держим в уме


// т. к нам абсолютно важен запуск всех горутин цикла wg.add правильнее разместить до самого цикла
// соответственно необходимо контролировать запуск всех горутин пишущих в канал, количество чтений из канала
// для этих целей заводим константу times
// обращаем внимание что wg в примере котролирует только запись в канал но не чтение из него
// т.к чтение из буферизированного канала у нас блокирует запись
//	wg.Wait моожно размесить только после цикла чтения что будет уже не логично т.к "зачем контролировать запись если мы уже прочитали?"
// пототму гораздо идеоматичнее в данном примере wg контролировать именно запуск горутин, соответственно убираем defer
// раз убрали defer - переносим wait до цикла чтения

// итак полседнее - нам нужен механизм остановки чтения функцией чтения с контролем количества прочитанных сообщщений
//	из 2 предложенных ранее вариантов, выберу с закрытием канала
//	т. к функция чтения у нас "однопоточная" достаточно простого счётчика - причитали условленное количество раз - запускаем мезанизм остановки чтения
// заводим переменную счётчик
// проверяем счётчик на каждом чтении..

func t1() {
	wg := sync.WaitGroup{}
	ch := make(chan string)
	mu := sync.Mutex{}
	const times int = 5  // именно столько запусков горутин и чтений наш код должен контролировать

	wg.Add(times)
	for i := 0; i < times; i++ {

		go func(i int, group *sync.WaitGroup) {
			group.Done()

			mu.Lock()
			ch <- fmt.Sprintf("Goroutine %d", i)
			mu.Unlock()
		}(i, &wg)
	}
	wg.Wait() // - проконтролировали запуск всех горутин
	i := 0
	for {
		select {
		case s, ok := <-ch:
			i ++
			if i == times{
				close(ch)
			}
			if ! ok {  // решение проблемы выхода из цикла
				return
			}
			fmt.Println(s)
		}
	}

	//wg.Wait() намеренно его здесь оставил показать где он был
}


// Ещё варианты..

func t4() {
	wg := sync.WaitGroup{}
	ch := make(chan string)
	mu := sync.Mutex{}
	times := 5

	wg.Add(times)
	for i := 0; i < times; i++ {

		go func(i int, wg *sync.WaitGroup) {
			wg.Done()

			mu.Lock()
			ch <- fmt.Sprintf("Goroutine %d", i)
			mu.Unlock()
		}(i, &wg)
	}
	wg.Wait()

	for {
		if times < 1{
			return
		}
		select {
		case s := <-ch:
			fmt.Println(s)
			times --
		}
	}

}




func t3() {
	wg := sync.WaitGroup{}
	ch := make(chan string)
	mu := sync.Mutex{}
	times := 5

	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(i int, group *sync.WaitGroup) {
			defer wg.Done()

			mu.Lock()
			ch <- fmt.Sprintf("Goroutine %d", i)
			mu.Unlock()
		}(i, &wg)
	}
	for {
		select {
		case s, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println(s)
			times --
			if times < 1{
				close(ch)
			}
		}
	}

	wg.Wait()
}





func t2() {
	wg := sync.WaitGroup{}
	ch := make(chan string)
	mu := sync.Mutex{}

	wg.Add(5)
	for i := 0; i < 5; i++ {

		go func(group *sync.WaitGroup, i int) {

			mu.Lock()
			ch <- fmt.Sprintf("Goroutine %d", i)
			mu.Unlock()
		}(&wg, i)
	}

	go func(ch chan string, wg *sync.WaitGroup){
		for {
			select {
			case s := <-ch:
				wg.Done()
				fmt.Println(s)
			}
		}

	}(ch, &wg)
	wg.Wait()
	close(ch)

}