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
    client - получаем структуру клиента
     doneChan - Канал для отправки, в него отправляем структуру "Track"
    resultChan - Канал для отправки, в него отправляем структуру со Статистикой
     Start запускает игру асинхронно и сразу возвращает управление
*/
func (t *Track) Start(client *Client, doneChan chan<- *Track, resultChan chan<- GameResult) {
    t.mu.Lock()
    t.use = true 
    t.client = client 
    startTime := time.Now()
    t.mu.Unlock()

    // Запускаем игру в отдельной горутине
    go t.playGame(client, startTime, doneChan, resultChan)
}

// Запуск игры
func (t *Track) playGame(client *Client, startTime time.Time, doneChan chan<- *Track, resultChan chan<- GameResult) {
    defer func() {
        t.mu.Lock()
        t.use = false
        t.client = nil
        t.mu.Unlock()

        doneChan <- t // дорожка освободилась
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

    result := GameResult{
        clientId:      client.id,
        trackId:       t.id,
        timeGameStart: startTime,
        timeGameEnd:   time.Now(), // реальное время окончания
        score:         score,
    }

    resultChan <- result
}
