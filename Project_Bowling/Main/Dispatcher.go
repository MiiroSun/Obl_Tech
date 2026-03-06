package main

import (
	"fmt"
	"time"
	"sync"
)

//Dispatcher
type Dispatcher struct{
	 tracks      []*Track			// Список дорожек
     queue       []Client			// Очередь ожидания
	 queueMu	 sync.Mutex
	 doneTrack   chan *Track		// Поток для свободных дорожек
	 resultGame  chan GameResult	// Поток для результата игр
	 results     []GameResult		// Результаты всех игр
	 resultsMu   sync.Mutex
	 activeGame  int				// Кол-во завершённых игр
	 countClient int				// Кол-во ушедших клиентов
	 totalClient int				// Общее кол-во клиентов которые сегодня придут 
}

// Метод который запускается из main 
func(t *Dispatcher) StartSystem(config Config){

	resultClient := make(chan Client, config.quantityClient) // Созданные клиенты
	t.resultGame = make(chan GameResult, config.quantityClient) // результаты игры
	t.doneTrack = make(chan *Track, config.quantityTracks) // Свободный трек

	t.totalClient = config.quantityClient

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
			t.displayGameStatus()
		
		case freeTrack := <- t.doneTrack:
			fmt.Printf("Дорожка %d освободилась\n", freeTrack.id)
			t.handleFreeTrack(freeTrack)
		
		case GameResult := <- t.resultGame:
			t.handleGameResult(GameResult)
			//Смотрим текущие игры, если счётчик меняется - информируем
			t.displayGameStatus()
			
		if t.checkCompletion() {
			return
		}
		}
	}
}

//region Первичная обработка

// Обрабатываем нового клиента
func(t *Dispatcher) handleNewClient(client Client) {
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
	t.activeGame++  

	go track.Start(&client, t.doneTrack, t.resultGame)
}

func(t *Dispatcher) handleFreeTrack(track *Track) {
	t.queueMu.Lock()
	// ПРоверяем есть ли кто в очереди
	if len(t.queue) > 0 {
		nextClient := t.queue[0]
		t.queue = t.queue[1:]
		t.queueMu.Unlock()
		
		fmt.Printf("Клиенту %s назначена дорожка %d\n", nextClient.name, track.id)
		t.activeGame++  

		// Запускаем игру
		go track.Start(&nextClient, t.doneTrack, t.resultGame)
	} else {
		t.queueMu.Unlock()
	}
}

//endregion

//region Очередь
func(t *Dispatcher) addToQueue(client Client) {
	t.queueMu.Lock()
	t.queue = append(t.queue, client)
	t.queueMu.Unlock()

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
	t.removeFromQueue(client.id)
}

// Удаление из очереди
func(t *Dispatcher) removeFromQueue(clientId int) {
	t.queueMu.Lock()
	defer t.queueMu.Unlock()
	
	for i, c := range t.queue {
		if c.id == clientId {
			t.queue = append(t.queue[:i], t.queue[i+1:]... )
			fmt.Printf("Клиент %s не стал ждать и ушел. Клиентов в очереди: %d\n", 
				c.name, len(t.queue))
			t.countClient++
			break
		}
	}
}

//endregion

//region Результаты

//Обрабатываем результаты игры
func(t *Dispatcher) handleGameResult(result GameResult) {
	t.activeGame--
	t.countClient++

	t.resultsMu.Lock()
	t.results = append(t.results, result)
	t.resultsMu.Unlock()
	
	fmt.Printf("Игра завершена: клиент %d на дорожке %d, счёт: %d, время игры: %v\n",
		result.clientId, result.trackId, result.score, result.timeGameEnd.Sub(result.timeGameStart))
}

// Информация
func (t *Dispatcher) displayGameStatus() {

	t.queueMu.Lock()
	defer t.queueMu.Unlock()

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

//Проверка нужно ли "закрываться"
func (t *Dispatcher) checkCompletion() bool{
	t.queueMu.Lock()
	queueLen := len(t.queue)
	t.queueMu.Unlock()

	if queueLen == 0 && t.activeGame == 0 && t.countClient == t.totalClient {
			return t.printFinalStats()
	}	
	return false 
}

//Выводим итоговую статистику и завершаем приложение
func(t *Dispatcher) printFinalStats() bool{
	fmt.Println("ИТОГИ")
	t.resultsMu.Lock()
	defer t.resultsMu.Unlock()

		fmt.Printf("%-10s %-15s %-15s %-12s %-10s\n", 
		"Клиент", "Время прихода", "Длительность", "Дорожка", "Счёт")
		fmt.Println("------------------------------------------------")
	
	for _, r := range t.results {

		fmt.Printf("%-10d %-15s %-15v %-12d %-10d\n",
			r.clientId,
			r.timeGameStart.Format("15:04:05"),
			r.timeGameEnd.Sub(r.timeGameStart).Round(time.Second),
			r.trackId,
			r.score)
	}
	return true
}

//endregion