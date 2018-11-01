package main

import (
	"fmt"
	"math/rand"
)

var trainCount int

type Train struct {
	Name       string
	Passengers int
}

func NewTrain(name string, passengers int) *Train {
	return &Train{name, passengers}
}

func genTrain() *Train {
	//num := rand.Intn(9)
	trainCount++
	p := rand.Intn(100)
	return &Train{
		Name:       fmt.Sprintf("TR%d", trainCount),
		Passengers: p,
	}
}
