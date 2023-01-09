package main

import (
	"fmt"

	"github.com/EnricoPicci/go-class/src/concurrency-patterns/recursive-examples/hilberthotel"
)

func RoomKeysClerk(upTo int, keysCh chan<- int) {
	for i := 0; i < upTo; i++ {
		keysCh <- i + 1
	}
	close(keysCh)
}

func BusClerk(busNumber int, roomKeysCh <-chan int, welcomeKitsCh chan<- []hilberthotel.WelcomeKit, parallelism int) {
	var count = 0
	var passengerNumber = 1
	var nextClerkCh chan int

	welcomeKits := []hilberthotel.WelcomeKit{}

	for roomKey := range roomKeysCh {
		count++
		if nextClerkCh == nil {
			nextClerkCh = make(chan int, parallelism)
			go BusClerk(busNumber+1, nextClerkCh, welcomeKitsCh, parallelism)
		}
		if count == passengerNumber {
			kit := hilberthotel.NewWelcomeKit(busNumber, passengerNumber, roomKey)
			welcomeKits = append(welcomeKits, kit)
			passengerNumber++
			count = 0
			continue
		}
		nextClerkCh <- roomKey
	}

	if nextClerkCh != nil {
		welcomeKitsCh <- welcomeKits
		close(nextClerkCh)
	} else {
		close(welcomeKitsCh)
	}
}

func GoHilbert(upTo int, parallelism int) []hilberthotel.WelcomeKit {
	if parallelism < 0 {
		parallelism = 0
	}
	keysCh := make(chan int, parallelism)
	go RoomKeysClerk(upTo, keysCh)

	hilbertCh := make(chan []hilberthotel.WelcomeKit, parallelism)
	go BusClerk(1, keysCh, hilbertCh, parallelism)

	kits := []hilberthotel.WelcomeKit{}
	for busKits := range hilbertCh {
		kits = append(kits, busKits...)
	}

	fmt.Println()
	fmt.Printf("%v guests have been given a room by Hilber at his Hotel\n", len(kits))

	return kits
}
