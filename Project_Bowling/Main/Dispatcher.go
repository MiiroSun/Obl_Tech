package main

import (
	"fmt"
	"time"
)

//Dispatcher
type Dispatcher struct{
	 tracks      []*Track        // Список дорожек
     queue       []Client        // Очередь ожидания
	 doneTrack   chan *Track     // Поток для свободных дорожек
	 resultGame  chan GameResult // Поток для результата игр
   
}

// Метод который запускается из main 
func(t *Dispatcher) StartSystem(config Config){

	resultClient := make(chan Client, config.quantityClient)
	resultGame := make(chan GameResult, config.quantityClient)
	doneTrack := make(chan *Track, config.quantityTracks)

	// Вызываем метод для создания клиентов 
	client := &Client{}
	go client.CreateClient(config, resultClient)
	
	// Инициализируем дорожки
	for i := 0; i < config.quantityTracks; i++ {
		t.tracks = append(t.tracks, &Track{id: i+1, use:false})
	}

	//Основной цикл
	for {
		select {
		case newClient := <- resultClient:
			fmt.Printf("Клиент %s пришёл в %s\n", newClient.name, newClient.arrivalTime.Format("00:00:00"))
			t.handleNewClient(newClient)
		
		case freeTrack := <- doneTrack:
			fmt.Printf("Дорожка %d освободилась\n", freeTrack.id)
			t.handleFreeTrack(freeTrack)
		
		case GameResult := <- resultGame:
			t.handleGameResult(GameResult)
		}

		//Смотрим текущие игры, если счётчик меняется - информируем
		t.displayGameStatus()
	}
}

//region Первичная обработка

// Обрабатываем нового клиента
func(t *Dispatcher) handleNewClient(client Client){
	//Смотрим, есть ли свободная дорожка
	for _, track := range t.tracks {
		track.mu.Lock()
		if !track.use {
			//Если дорожка свободна - играем
			track.mu.Unlock()
			t.startGameOnTrack(track, client)
			return
		}
		track.mu.Unlock()
	}

	//Если не нашли свободных - в очередь
	t.addToQueue(client)
}

//Запускаем игры на дорожках
func(t *Dispatcher) startGameOnTrack(track *Track, client Client) {
	fmt.Printf("Клиент %s начал игру на дорожке %d, будет играть %v\n",
		client.name, track.id, client.playTime)	

	go track.Start(&client, t.doneTrack, t.resultGame)
}

func(t *Dispatcher) handleFreeTrack(track *Track) {
	// ПРоверяем есть ли кто в очереди
	if len(t.queue) > 0 {
		nextClient := t.queue[0]
		t.queue = t.queue[1:]

		fmt.Printf("Клиенту %s назначена дорожка %d\n", nextClient.name, track.id)

		// Запускаем игру
		go track.Start(&nextClient, t.doneTrack, t.resultGame)
	}
}

//endregion

//region Очередь
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

//endregion

//region Результаты

//Обрабатываем результаты игры
func(t *Dispatcher) handleGameResult(result GameResult) {
	fmt.Printf("Игра завершена: клиент %d на дорожке %d, счёт: %d, время игры: %v\n",
		result.clientId, result.trackId, result.score, result.timeGameEnd.Sub(result.timeGameStart))
}

// Информация
func (t *Dispatcher) displayGameStatus() {
	//Информация о текущих играх
	fmt.Println("\n ТЕКУЩИЕ ИГРЫ")

	//Информация о дорожках
	for _, track := range t.tracks {
		track.mu.Lock()
		if track.use && track.client != nil {
			fmt.Printf("Дорожка %d: занята клиентом %s\n", track.id, track.client.name)
		} else {
			fmt.Printf("Дорожка %d: свободна \n", track.id)
		}
		track.mu.Unlock()
	}

	//Информация про оччередь
	fmt.Printf("В очереди: %d клиентов\n", len(t.queue))
	if len(t.queue) > 0 {
		for i, client := range t.queue {
			fmt.Printf(" %d. %s (ждёт %v)\n", i+1, client.name, client.waitTime)
		}
	}
	fmt.Println()
}

//endregion