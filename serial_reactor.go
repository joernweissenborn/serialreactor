package serialreactor

import (
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
	buf := make([]byte, sr.maxPacketSize)
	for {
		n, _ := sr.s.Read(buf)
		sr.Reactor.Fire(ReadEvent{}, buf[:n])
	}
}

func (sr *SerialReactor) OnRead(handler eventual2go.Subscriber) {
	sr.Reactor.React(ReadEvent{}, handler)
}
