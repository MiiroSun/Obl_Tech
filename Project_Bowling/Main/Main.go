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
	quantityTracks:		3,						// Сколько дорожек
	quantityClient:		8, 						// Сколько клиентов будет
	maxClientInterval:  1,						// Максимальное время в рамках которого будут "приходить" клиенты (кол-во минут) 
	maxWaitTime: 		20 * time.Second,		// Сколько будут ждать (сек)
	gameTimeMin:		10 * time.Second,		// Мин. время игры (сек)
	gameTimeMax:		30 * time.Second,		// Макс. время игры (сек)
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

	// Запускаем систему
	dispatcher.StartSystem(config)
}