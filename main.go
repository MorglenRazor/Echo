package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var messageChan chan string
var message string
var interval time.Duration = 1
var serverAdress  = "localhost:3000"

//Функция бесконечного вывода сообщения
func echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Client connection is Open")

		//Заголовки...
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Открытие канала
		messageChan = make(chan string)

		// Выход из канала
		defer func() {
			close(messageChan)
			messageChan = nil
			log.Printf("Client connection is Closed")
		}()
		message = r.URL.Query().Get("w")
		// Объявление Flusher
		flusher, _ := w.(http.Flusher)
		for {

			select {

			// Вывод сообщения раз в секунду
			case <-time.After(interval * time.Second):
				fmt.Fprintf(w, "%s\n", message)
				flusher.Flush()

			// закрытие конекта к каналу и переход в defer
			case <-r.Context().Done():
				return

			}
		}
	}
}
//Функция изменения выводимого сообщения
func say() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if messageChan != nil {
			log.Printf("Change message...")
			//обявление нового сообщения
			message = r.URL.Query().Get("w")
		}

	}
}

func main() {

	http.HandleFunc("/echo", echo())

	http.HandleFunc("/say", say())

	log.Fatal("HTTP server error: ", http.ListenAndServe(serverAdress, nil))
}