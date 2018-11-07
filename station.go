package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type EventType int

const (
	_ EventType = iota
	EventServiceStart
	EventTrainArrived
	EventTrainDeparted
)

type Event struct {
	Type    EventType
	Msg     string
	Payload interface{}
}

type Station struct {
	Name      string
	Platforms map[lineID]*Platform
	Events    chan Event

	probes map[EventType]bool
}

func (s *Station) String() string {
	return fmt.Sprintf(s.Name)
}

func NewStation(name string) *Station {
	return &Station{
		Name:      name,
		Platforms: make(map[lineID]*Platform),
		Events:    make(chan Event, 1),
		probes:    make(map[EventType]bool),
	}
}

func (s *Station) AddProbe(event EventType) {
	s.probes[event] = true
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

func (s *Station) Run(ctx context.Context, wg *sync.WaitGroup) {
	for _, pf := range s.Platforms {
		processTrains := func(k *Track) {
			wg.Add(1)
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					//fmt.Printf("station %s is closing\n", s.Name)
					return
				case tr, ok := <-k.In:
					if !ok {
						return
					}
					//arrived := tr.Passengers
					tr.Passengers = rand.Intn(100)

					//log.Printf("%s train %s arrived at %s with passengers: %d -> %d\n", k.dir, tr.Name, s.Name, arrived, tr.Passengers)
					if _, ok := s.probes[EventTrainArrived]; ok {
						s.Events <- Event{Type: EventTrainArrived, Payload: tr}
					}

					time.Sleep(time.Second * 1)

					if k.Out != nil {
						s.depart(k, tr)
					}
				}
			}
		}

		u := pf.Up
		d := pf.Down
		go processTrains(u)
		go processTrains(d)
	}
}

// depart emits the train departed event and simulates travel time
// before writing the train to the IN channel of the next station.
func (s *Station) depart(k *Track, tr *Train) {
	if _, ok := s.probes[EventTrainDeparted]; ok {
		s.Events <- Event{Type: EventTrainDeparted, Payload: tr}
	}
	time.Sleep(time.Second * 2) // travel time
	k.Out <- tr
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
