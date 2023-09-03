package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const serverPort = 9075

const name = "BIOMETRIC"

func main() {
	handler := http.NewServeMux()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: handler,
	}

	LoadBanData()
	SetupShutdownHandler()

	done := make(chan struct{})

	go func() {
		handler.HandleFunc("/api/ban/", func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodPost {
				http.Error(writer, "Error method", http.StatusMethodNotAllowed)
				return
			}

			var banData Ban
			decoder := json.NewDecoder(request.Body)
			if err := decoder.Decode(&banData); err != nil {
				http.Error(writer, "Error decoding json", http.StatusBadRequest)
			}

			cache = append(cache, banData)

			writer.WriteHeader(http.StatusOK)
		})

		handler.HandleFunc("/api/bans", func(writer http.ResponseWriter, request *http.Request) {
			if request.Method != http.MethodGet {
				http.Error(writer, "Error method", http.StatusMethodNotAllowed)
				return
			}

			fmt.Printf("%v\n", cache)

			if len(cache) == 0 {
				writer.Write([]byte("null"))
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(writer).Encode(cache); err != nil {
				http.Error(writer, "Error encoding json", http.StatusInternalServerError)
				return
			}
		})

		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Printf("(%v) Server has been shut down: %s\n", name, err)
				close(done)
			} else {
				log.Fatalf("(%v) Error running server: %s\n", name, err)
			}
		}
	}()

	<-done
}
