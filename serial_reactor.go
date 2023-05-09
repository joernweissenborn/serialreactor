package serialreactor

import (
	"bufio"

	"github.com/joernweissenborn/eventual2go"
	"github.com/tarm/serial"
)

type ReadEvent struct{}

type SerialReactor struct {
	*eventual2go.Reactor
	s             *serial.Port
	maxPacketSize int
}

func New(port string, baudrate int, maxPacketSize int) (sr *SerialReactor, err error) {
	c := &serial.Config{Name: port, Baud: baudrate}
	s, err := serial.OpenPort(c)
	if err != nil {
		return
	}

	sr = &SerialReactor{
		Reactor:       eventual2go.NewReactor(),
		s:             s,
		maxPacketSize: maxPacketSize,
	}
	return
}

func (sr *SerialReactor) Listen() {
	go sr.listen()
}

func (sr *SerialReactor) listen() {
	scanner := bufio.NewScanner(sr.s)
	for scanner.Scan() {
		sr.Reactor.Fire(ReadEvent{}, scanner.Text())
	}
}

func (sr *SerialReactor) OnRead(handler eventual2go.Subscriber) {
	sr.Reactor.React(ReadEvent{}, handler)
}
