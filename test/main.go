package main

import (
	"fmt"
)

type Room struct {
	Number     int
	TypeRoom   int
	PriceNight float32
	Status     bool
}

type Rezerf struct {
	NameCustomer string
	NumberRoom   int
	Start        string
	NumberNight  int
}

type Hotel struct {
	Rooms   []Room
	Rezerfs []Rezerf
}

func (h *Hotel) AddRoom(number, typeRoom int, priceNight float32, status bool) error {
	for _, room := range h.Rooms {
		if room.Number == number {
			return fmt.Errorf("found Room Error")
		}
		h.Rooms = append(h.Rooms, Room{Number: number, TypeRoom: typeRoom, PriceNight: priceNight, Status: status})
	}
	return nil
}

func (h *Hotel) ListStatusRoom() {
	for _, room := range h.Rooms {
		if !room.Status {
			fmt.Println(room)
		}
	}
}

func (h *Hotel) RezervRoom(nameCustomer string, numberRoom int, start string, numberNight int) (bool, error) {
	for _, room := range h.Rooms {
		if room.Number == numberNight && !room.Status {
			room.Status = true

			h.Rezerfs = append(h.Rezerfs, Rezerf{NameCustomer: nameCustomer, NumberRoom: numberRoom, Start: start, NumberNight: numberNight})
			return true, nil
		}
	}
	return false, fmt.Errorf("not found error")
}
