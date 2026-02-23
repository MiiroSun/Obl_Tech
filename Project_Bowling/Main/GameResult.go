package main

import (
	"time"
)

// Статистика
type GameResult struct {
	clientId     int       // id клиента
	trackId       int       // id дорожки
	timeGameStart time.Time // время начала игры
	timeGameEnd   time.Time // время окончания игры
	score         int       // счет
}