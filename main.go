package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// data variable
var messageChan chan string

func writeData(w http.ResponseWriter) (int, error) {
	// set data into response writer
	return fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
}

func sseStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// initialize messageChan
		messageChan = make(chan string)

		// calling anonymous function that
		// closing messageChan channel and
		// set it to nil
		defer func() {
			close(messageChan)
			messageChan = nil
		}()

		// handle server sent event header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// create http http.Flusher that allows
		// http handler to flush buffered data to
		// client until closed
		flusher, _ := w.(http.Flusher)
		for {
			write, err := writeData(w)
			if err != nil {
				log.Println(err)
			}
			log.Println(write)
			flusher.Flush()
		}
	}
}

// sendMessage used to write data into messageChan and flushed to client through sseStream
func sseMessage() http.HandlerFunc {

	type DataRequest struct {
		Tweet string `json:"tweet"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var data DataRequest
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if messageChan != nil {
			messageChan <- data.Tweet
		}
	}
}

func client() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// Define the data to be displayed in the template
		data := struct {
			Title string
		}{
			Title: "Welcome to My Website",
		}

		// Parse the HTML template
		htmlTemplate, err := template.ParseFiles("streamer.html")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the HTML template with the data
		err = htmlTemplate.Execute(writer, data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func post() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// Define the data to be displayed in the template
		data := struct {
			Title string
			Body  string
		}{
			Title: "Client that make triger to server",
		}

		// Parse the HTML template
		htmlTemplate, err := template.ParseFiles("post.html")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the HTML template with the data
		err = htmlTemplate.Execute(writer, data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func main() {
	http.HandleFunc("/", client())
	http.HandleFunc("/post", post())
	http.HandleFunc("/stream", sseStream())
	http.HandleFunc("/send", sseMessage())

	log.Println("Server running on port :8080")
	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}
