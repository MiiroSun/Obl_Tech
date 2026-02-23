package main

import (
	"sync"
	"time"
	"math/rand"
)

// Дорожка
type Track struct {
	id     int
	client *Client    // Клиент, который играет на дорожке
	use    bool       // Используется ли дорожка
	mu     sync.Mutex // Мьютекс для синхронизации доступа к дорожке
}

/*
    Метод для запуска дорожки 
     client - получаем структуру клиента
     doneChan - Канал для отправки, в него отправляем структуру "Track"
     resultChan - Канал для отправки, в него отправляем структуру "GameResult"
     return: В потоках возвращается объект дорожки и статистики по этой дорожке
*/
func (t *Track) Start(client *Client, doneChan chan<- *Track, resultChan chan<- GameResult) {
    //Инициализация дорожки
    // Сделал так, чтобы иницилазиция была конкретной дорожки 
    t.mu.Lock()
    t.use = true 
    t.client = client 
    startTime := time.Now()
    t.mu.Unlock()

    // Запускаем игру в отдельной горутине
    go t.playGame(client, startTime, doneChan, resultChan)
}

/*
    Метод для запуска игры
     client - Структура клиента
     startTime - Время начала TODO: Мб всё таки внутри запуска игры ?
     doneChan - Поток с Track 
     resultChan - Поток со статистикой
     return: В каналах возвращается пустая дорожка (игра заканчивается) и заполненная статистика
*/
func (t *Track) playGame(client *Client, startTime time.Time, doneChan chan<- *Track, resultChan chan<- GameResult) {
    defer func() {
        t.mu.Lock()
        t.use = false
        t.client = nil
        t.mu.Unlock()

        doneChan <- t // дорожка освободилась. Возвращаю структуру, т.к. дорожек то у меня n и создавать новые нет необходимости, работаю с одним объектом
    }()

    endTime := startTime.Add(client.playTime)

    var score int
    ticker := time.NewTicker(800 * time.Millisecond)
    defer ticker.Stop()

    for now := range ticker.C {
        if now.After(endTime) {
            break
        }

        if rand.Float64() < 0.35 {
            points := rand.Intn(15) + 5
            score += points
            // опционально: логировать или отправлять обновление счёта в реальном времени
        }
    }

    // По логике, тут храним статистику по только что завершённой игре 
    result := GameResult{
        clientId:      client.id,
        trackId:       t.id,
        timeGameStart: startTime,
        timeGameEnd:   time.Now(), 
        score:         score,
    }

    resultChan <- result
}
