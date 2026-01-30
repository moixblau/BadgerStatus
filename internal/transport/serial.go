package transport

import (
	"io"

	"go.bug.st/serial"
)

// package-level seam to allow tests to replace serial.Open
var serialOpen = serial.Open

type BadgerTransport struct {
	port serial.Port
}

func NewBadgerTransport(portName string, baudRate int) (*BadgerTransport, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serialOpen(portName, mode)
	if err != nil {
		return nil, err
	}
	return &BadgerTransport{port: port}, nil
}

func (t *BadgerTransport) Write(p []byte) (n int, err error) {
	return t.port.Write(p)
}

func (t *BadgerTransport) Close() error {
	return t.port.Close()
}

var _ io.WriteCloser = (*BadgerTransport)(nil)
