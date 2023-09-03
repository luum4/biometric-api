package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cache []Ban

var fileName = "ban_data.json"

type Ban struct {
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	XUID   string `json:"xuid"`
	Reason string `json:"reason"`
}

func LoadBanData() {
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("(%v) No %v found. A new file will be created", name, fileName)
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("(%v) Error creating %v: %v", name, fileName, err)
			return
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("(%v) Error creating %v: %v", name, fileName, err)
				return
			}
		}(file)

		jsonData := []byte("[]")
		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatalf("(%v) Error writing json data into the file: %v", name, err)
			return
		}
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("(%v) Error reading file: %v", name, err)
		return
	}

	if len(content) == 0 || len(content) == 1 {
		log.Fatalf("(%v) Error the file is empty: %v", name, err)
		return
	}

	if string(content) == "[]" {
		return
	}

	if err := json.Unmarshal(content, &cache); err != nil {
		fmt.Printf("(%v) Error loading ban data: %v\n", name, err)
		return
	}
}

func SaveBanData() {
	if len(cache) == 0 {
		log.Printf("(%v) Cache is empty. No data to save.\n", name)
		return
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		log.Printf("(%v) Error saving ban data: %v\n", name, err)
		return
	}

	if err := ioutil.WriteFile("ban_data.json", data, 0644); err != nil {
		log.Printf("(%v) Error saving ban data into the file: %v\n", name, err)
	}
}

func SetupShutdownHandler() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		SaveBanData()
		log.Fatalf("(%v) Server stopping", name)
	}()
}
