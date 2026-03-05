package main

import (
    "fmt"
    "math/rand"
    "time"
    "sort"
)

// Клиент
type Client struct {
	id       int
	name     string        //Имя
    arrivalTime time.Time // Время прихода
	playTime time.Duration // Сколько времени будет играть
	waitTime time.Duration // Сколько может подождать
	leave    bool          // Ушел ли клиент
}


/*
    Метод для создания клиента
     config - конфигурация
     resultClient - поток в котором будут обрабатываться юзеры
     return: Юзеры согласно их времени прихода 
*/
func (t *Client) CreateClient(config Config,resultClient chan<- Client) {
   
    var clients []Client

    // Создаём клиентов, кол-во согласно конфигу
    for clientId := 1;clientId <= config.quantityClient; clientId ++ {

        // Генерим массив минут 
        randomMinutes := rand.Intn(config.maxClientInterval + 1) 
        arrivalTime := time.Now().Add(time.Duration(randomMinutes) * time.Minute)

        // Генераируем клиента
        clients = append(clients, Client{
            id:			clientId,
            name:       fmt.Sprintf("Client_%d", clientId),
            arrivalTime:  arrivalTime,
            playTime:  	config.gameTimeMin + time.Duration(rand.Int63n(int64(config.gameTimeMax - config.gameTimeMin))),
            waitTime: 	time.Duration(rand.Int63n(int64(config.maxWaitTime))), 
            leave:		false,
        })
    }
        //Сортировочка
        sort.Slice(clients, func(a,b int) bool {
            return clients[a].arrivalTime.Before(clients[b].arrivalTime)
        })  

        //Отправка в канал согласно времени прибытия
        for _, result := range clients {
            time.Sleep(time.Until(result.arrivalTime))

            resultClient <- result
        }
} 
