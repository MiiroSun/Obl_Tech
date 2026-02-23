package main

import "time"

// Клиент
type Client struct {
	id       int
	name     string        //Имя
	playTime time.Duration // Сколько времени будет играть
	waitTime time.Duration // Сколько может подождать
	leave    bool          // Ушел ли клиент
}
