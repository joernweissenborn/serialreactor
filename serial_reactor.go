package serialreactor

import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/tarm/serial"
)

type ReadEvent struct{}

type SerialReactor struct {
	*eventual2go.Reactor
	s *serial.Port
}

func New() (sr *SerialReactor, err error) {
	c := &serial.Config{Name: "/dev/ttyUSB1", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		return
	}

	sr = &SerialReactor{
		Reactor: eventual2go.NewReactor(),
		s:       s,
	}
	go sr.listen()
	return
}

func (sr *SerialReactor) listen() {
	buf := make([]byte, 1024)
	for {
		n, _ := sr.s.Read(buf)
		sr.Reactor.Fire(ReadEvent{}, buf[:n])
	}
}

func (sr *SerialReactor) OnRead(handler eventual2go.Subscriber) {
	sr.Reactor.React(ReadEvent{}, handler)
}
