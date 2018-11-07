package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"
)

func parseJSON(file string) ([]interface{}, error) {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err = json.Unmarshal(d, &m); err != nil {
		return nil, err
	}

	stn, ok := m["stations"]
	if !ok {
		return nil, fmt.Errorf("'stations' key not found")
	}

	stns, ok := stn.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'stations' is of wrong type, expected an array")
	}

	//_, ok := stns[0].(string)
	//if !ok {
	////fmt.Printf("%T", strStns)
	//}

	return stns, nil
}

func main() {
	var config string
	var up, down bool
	var trains int

	flag.StringVar(&config, "d", "data.json", "data filename")
	flag.IntVar(&trains, "n", 1, "number of trains to run")
	flag.BoolVar(&up, "up", false, "run trains only in UP direction")
	flag.BoolVar(&down, "down", false, "run trains only in DOWN direction")
	flag.Parse()

	stnNames, err := parseJSON(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	line := GreenLine
	stns := make([]*Station, 0)

	for _, n := range stnNames {
		name := n.(string)
		s := NewStation(name)
		s.AddLine(line)
		stns = append(stns, s)
	}

	Connect(line, stns...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	var wg sync.WaitGroup

	for _, stn := range stns {
		go monitor(stn)
		stn.Run(ctx, &wg)
	}

	// TODO: add station monitor

	rand.Seed(time.Now().UnixNano())

	if !down {
		fmt.Printf("running service: %s --> %s\n", stns[0].String(), stns[len(stns)-1])
		stns[0].StartService(line, trains)
	}

	if !up {
		fmt.Printf("running service: %s --> %s\n", stns[len(stns)-1], stns[0])
		stns[len(stns)-1].StartService(line, trains)
	}

	wg.Wait()
	cancel()
}

func monitor(stn *Station) {
	stn.AddProbe(EventTrainArrived)
	stn.AddProbe(EventTrainDeparted)

	for ev := range stn.Events {
		//log.Printf("%s train %s arrived at %s with passengers: %d -> %d\n", k.dir, tr.Name, s.Name, arrived, tr.Passengers)
		switch ev.Type {
		case EventTrainArrived:
			tr := ev.Payload.(*Train)
			log.Printf("train %s arrived at %s\n", tr.Name, stn.Name)
		case EventTrainDeparted:
			tr := ev.Payload.(*Train)
			log.Printf("train %s departed %s\n", tr.Name, stn.Name)
		}
	}
}
