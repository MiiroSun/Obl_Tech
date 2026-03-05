package main

import (
	"fmt"
	"time"
)

//Dispatcher
type Dispatcher struct{
	 tracks      []*Track      // Список дорожек
     queue       []Client      // Очередь ожидания
   
}

// Метод который запускается из main 
func(t *Dispatcher) StartSystem(config Config){

	resultClient := make(chan Client, config.quantityClient)

	// Вызываем метод для создания клиентов 
	client := &Client{}
	go client.CreateClient(config, resultClient)
	
	// Закидываем клиента в очередь 
		// Если есть свободная дорожка - играть 
		// Если клиент ждёт, то начинается его таймер
			// По завершению таймера - удаляем клиента 	

	// Смотрим текущие игры, если счёт изменился - выводим инфу об этом

}

// Методы с очередью
func(t *Dispatcher) addToQueue(client Client) {
	t.queue = append(t.queue, client)
	fmt.Printf("Клиент %s добавлен в очередь. Клиентов в очереди: %d\n",
			client.name, len(t.queue))

	// Запускаем таймер ожидания клиента
	go t.startWaitTimer(client)
}

// Таймер ожидания клиента
func(t *Dispatcher) startWaitTimer(client Client) {
	timer := time.NewTimer(client.waitTime)
	defer timer.Stop()

	<-timer.C
	t.removeFromQueue(client)
}

// Удаление из очереди
func(t *Dispatcher) removeFromQueue(client Client) {
	for i, c := range t.queue {
		if c.id == client.id {
			t.queue = append(t.queue[:i], t.queue[i+1:]... )
			fmt.Printf("Клиент %s не стал ждать и ушел. Клиентов в очереди: %d\n", 
				client.name, len(t.queue))
			break
		}
	}
}

// Метод с Score 
