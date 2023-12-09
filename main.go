package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Time     string `json:"time"`
	Greeting string `json:"greeting"`
}

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Println(buf.String())

		// Create response
		responce := Response{
			Time:     time.Now().UTC().Format(time.RFC3339),
			Greeting: fmt.Sprintf("Hello, %s!", buf.String()),
		}
		jsonResponce, err := json.Marshal(&responce)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponce)
	})

	port := 80
	if p, found := os.LookupEnv("SERVER_PORT"); found {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			logger.Panic(err.Error())
		}
	}
	logger.Printf("server started at port %d", port)
	err := http.ListenAndServe(":"+fmt.Sprint(port), nil)
	if err != nil {
		logger.Panic(err.Error())
	}
}
