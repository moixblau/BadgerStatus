package transport

import (
	"errors"
	"testing"
	"time"

	"go.bug.st/serial"
)

type fakePort struct {
	written  []byte
	closed   bool
	writeErr error
	closeErr error
}

func (f *fakePort) SetMode(mode *serial.Mode) error { return nil }
func (f *fakePort) Read(p []byte) (int, error) {
	return 0, nil
}
func (f *fakePort) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	f.written = append(f.written, p...)
	return len(p), nil
}
func (f *fakePort) Drain() error { return nil }
func (f *fakePort) ResetInputBuffer() error {
	return nil
}
func (f *fakePort) ResetOutputBuffer() error {
	return nil
}
func (f *fakePort) SetDTR(dtr bool) error {
	return nil
}
func (f *fakePort) SetRTS(rts bool) error {
	return nil
}
func (f *fakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return &serial.ModemStatusBits{}, nil
}
func (f *fakePort) SetReadTimeout(t time.Duration) error {
	return nil
}
func (f *fakePort) Close() error {
	f.closed = true
	return f.closeErr
}
func (f *fakePort) Break(duration time.Duration) error {
	return nil
}

func TestNewBadgerTransport_OpenError(t *testing.T) {
	orig := serialOpen
	defer func() { serialOpen = orig }()

	serialOpen = func(name string, mode *serial.Mode) (serial.Port, error) {
		return nil, errors.New("open error")
	}

	_, err := NewBadgerTransport("/dev/ttyX", 9600)
	if err == nil {
		t.Fatalf("expected error from serialOpen, got nil")
	}
}

func TestBadgerTransport_WriteClose(t *testing.T) {
	orig := serialOpen
	defer func() { serialOpen = orig }()

	fp := &fakePort{}
	serialOpen = func(name string, mode *serial.Mode) (serial.Port, error) {
		return fp, nil
	}

	tp, err := NewBadgerTransport("/dev/ttyX", 9600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	n, err := tp.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5 bytes written got %d", n)
	}
	if string(fp.written) != "hello" {
		t.Fatalf("expected port to receive 'hello' got %s", string(fp.written))
	}

	if err := tp.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !fp.closed {
		t.Fatalf("expected port to be closed")
	}
}
