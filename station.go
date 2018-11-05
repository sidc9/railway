package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Station struct {
	Name      string
	Platforms map[lineID]*Platform
}

func (s *Station) String() string {
	return fmt.Sprintf(s.Name)
}

func NewStation(name string) *Station {
	return &Station{name, make(map[lineID]*Platform)}
}

func (s *Station) AddLine(id lineID) {
	st := NewPlatform()
	s.Platforms[id] = st
}

func (s *Station) StartService(line lineID, trains int) {
	sendTrain := func(trk *Track) {
		if trk.Out != nil {
			log.Printf("svc started from stn %s in %s direction \n", s.Name, trk.dir)
			for i := 0; i < trains; i++ {
				tr := genTrain()
				trk.Out <- tr
				time.Sleep(time.Second * 1)
			}
		}
	}

	trk := s.Platforms[line].Up
	go sendTrain(trk)

	trk = s.Platforms[line].Down
	go sendTrain(trk)
}

type Platform struct {
	Up   *Track
	Down *Track
}

func NewPlatform() *Platform {
	return &Platform{
		Up:   NewTrack(Up),
		Down: NewTrack(Down),
	}
}

func (s *Station) Run() {
	for _, pf := range s.Platforms {
		processTrains := func(k *Track) {
			for tr := range k.In {
				arrived := tr.Passengers
				tr.Passengers = rand.Intn(100)

				log.Printf("%s train %s arrived at %s with passengers: %d -> %d\n", k.dir, tr.Name, s.Name, arrived, tr.Passengers)

				time.Sleep(time.Second * 2)

				if k.Out != nil {
					k.Out <- tr
				}
			}
		}

		u := pf.Up
		d := pf.Down
		go processTrains(u)
		go processTrains(d)

	}
}

func Connect(line lineID, stations ...*Station) {
	for i := 0; i < len(stations)-1; i++ {
		stn := stations[i]
		nextStn := stations[i+1]

		pf := stn.Platforms[line]
		nextPf := nextStn.Platforms[line]
		pf.Up.Connect(nextPf.Up)
		nextPf.Down.Connect(pf.Down)
	}
}
