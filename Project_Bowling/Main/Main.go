package main

import (
	"fmt"
	"math/rand"
	"sync"
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

