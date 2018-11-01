package main

type direction int

const (
	Up   direction = 1
	Down direction = 2
)

func (d direction) String() string {
	switch d {
	case Up:
		return "UP"
	case Down:
		return "DOWN"
	}
	return "unknown"
}

type Track struct {
	In, Out chan *Train
	dir     direction
}

func NewTrack(dir direction) *Track {
	return &Track{dir: dir}
}

func (k *Track) Connect(trackB *Track) {
	k.Out = make(chan *Train, 1)
	trackB.In = k.Out
}
