package main

type lineID int

const (
	GreenLine lineID = 1
	RedLine   lineID = 2
)

type Line struct {
	LineID   lineID
	Stations []*Station
}

func NewLine(id lineID) *Line {
	return &Line{LineID: id}
}

func (l *Line) AddStation(station *Station) {
	l.Stations = append(l.Stations, station)
}
