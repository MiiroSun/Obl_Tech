package main

import (
	"fmt"
	"math/rand"
	"time"
)


// Config - конфигурация системы
type Config struct {
	quantityTracks  int  // Сколько дорожек
	quantityClient	int  // Сколько будет клиентов 
	maxClientInterval int  // Максимальное время в рамках которого будут "приходить" клиенты (кол-во минут) 
	maxWaitTime   time.Duration // Максимальное время ожидания у клиентов
	gameTimeMin   time.Duration // Минимальное время игры
	gameTimeMax   time.Duration // Максимальное время игры 
}


// Запускаем
func main() {
	rand.Seed(time.Now().UnixNano())

	//Конфиг
	config := Config{
		quantityTracks: 3,
		quantityClient: 8, 
		maxClientInterval: 1 ,  // Максимальное время в рамках которого будут "приходить" клиенты (кол-во минут) 
		maxWaitTime: 30 * time.Minute,      // ждут 30 секунд
	gameTimeMin:       10 * time.Second,      // играют 10-30 секунд
	gameTimeMax:       30 * time.Second,
	}

	fmt.Println("=== СИСТЕМА БОУЛИНГ-КЛУБА ЗАПУЩЕНА ===")
	fmt.Printf("Конфигурация: %d дорожек, %d клиентов\n", config.quantityTracks, config.quantityClient)
	fmt.Printf("Клиенты приходят в течение %d минут\n", config.maxClientInterval)
	fmt.Printf("Время игры: от %v до %v\n", config.gameTimeMin, config.gameTimeMax)
	fmt.Printf("Максимальное время ожидания: %v\n", config.maxWaitTime)
	fmt.Println("=====================================\n")

	// Создаем диспетчер
	dispatcher := &Dispatcher{
		tracks:     []*Track{},
		queue:      []Client{},
		doneTrack:  make(chan *Track, config.quantityTracks),
		resultGame: make(chan GameResult, config.quantityClient),
	}

	// Запускаем систему (блокирующий вызов)
	dispatcher.StartSystem(config)
}

/* Проблемы:
1. Не показывает если клиент ушёл
2. Не показывает набранные очки клиента
3. Не корректно отображается время прибытия клиента
4. Не отображается информация о том, что игра завершена 

*/