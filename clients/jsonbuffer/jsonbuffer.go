// Package jsonbuffer implements a stream protocol using JSON blobs as segments
// in the buffer, with a type and envelope.
//
// Each message is sent as JSON (see the Message struct) and contains a type
// and payload. The type determines how the payload is handled; but arguably
// what happens as a result of that belongs to the client.
//
// There is also a javascript client for this in the ci-ui repository.
package jsonbuffer

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"errors"
)

const (
	// TypeMessage is a message and the payload is the message we're trying to send.
	TypeMessage = "message"
	// TypeError is an error and the payload is the error text.
	TypeError = "error"
)

// Message is the websocket message transmitted.
type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// ReadWrapper is a read-only wrapper for jsonbuffer.
type ReadWrapper struct {
	io.Reader
	dec *json.Decoder

	leftoverReadMutex sync.Mutex
	leftoverRead      []byte
}

// WriteWrapper is a write-only wrapper for jsonbuffer.
type WriteWrapper struct {
	io.Writer
	enc *json.Encoder
}

// Wrapper is a wrapper for io.Reader that lets us power it via
// protocol through our client.
type Wrapper struct {
	*ReadWrapper
	*WriteWrapper
}

// NewWrapper creates a Wrapper from a pre-established connection.
func NewWrapper(rw io.ReadWriter) *Wrapper {
	return &Wrapper{NewReadWrapper(rw), NewWriteWrapper(rw)}
}

// NewReadWrapper is a read-only wrapper for jsonbuffer.
func NewReadWrapper(r io.Reader) *ReadWrapper {
	return &ReadWrapper{dec: json.NewDecoder(r)}
}

// NewWriteWrapper is a write-only wrapper for jsonbuffer.
func NewWriteWrapper(w io.Writer) *WriteWrapper {
	return &WriteWrapper{enc: json.NewEncoder(w)}
}

// Send sends a message by writing through the Wrapper. The outgoing message will be of
// type websocket.TypeMessage.
func (w *WriteWrapper) Send(message string) error {
	if err := w.enc.Encode(Message{Type: TypeMessage, Payload: message}); err != nil {
		return err
	}

	return nil
}

// SendError is like Send, but for errors.
func (w *WriteWrapper) SendError(err error) error {
	return w.enc.Encode(Message{Type: TypeError, Payload: err.Error()})
}

// Recv reads a single message from the reader and returns it. If the type is
// message, it returns the payload, otherwise if it hits EOF it will return
// ErrEOF; on any error return it returns the error.
func (w *ReadWrapper) Recv() (string, error) {
	var msg Message
	var eof bool
	if err := w.dec.Decode(&msg); err != nil && err != io.EOF {
		return msg.Payload, err
	} else if err == io.EOF && msg.Type != "" {
		eof = true
	} else if err == io.EOF {
		return msg.Payload, err
	}

	switch msg.Type {
	case TypeMessage:
		var err error
		if eof {
			err = io.EOF
		}

		return msg.Payload, err
	case TypeError:
		return "", errors.New(msg.Payload)
	default:
		return "", fmt.Errorf("invalid type %v", msg.Type)
	}
}

// Write is Send that conforms to the io.Writer spec. Each buffer will be sent as a single Message.
func (w *WriteWrapper) Write(buf []byte) (int, error) {
	err := w.Send(string(buf))
	if err != nil {
		return len(buf), err
	}

	return len(buf), nil
}

// Read conforms to the io.Reader interface. If it cannot fill buf with the
// payload, it will keep the buffer for the next read call.
func (w *ReadWrapper) Read(buf []byte) (int, error) {
	w.leftoverReadMutex.Lock()
	if len(w.leftoverRead) > 0 {
		l := copy(buf, w.leftoverRead)
		if l < len(w.leftoverRead) {
			w.leftoverRead = w.leftoverRead[l:]
		} else {
			w.leftoverRead = nil
		}

		w.leftoverReadMutex.Unlock()
		return l, nil
	}
	w.leftoverReadMutex.Unlock()

	rcv, err := w.Recv()
	if err != nil && err != io.EOF {
		return 0, err
	}

	l := copy(buf, rcv)

	if l < len(rcv) {
		w.leftoverReadMutex.Lock()
		w.leftoverRead = []byte(rcv[l:])
		w.leftoverReadMutex.Unlock()
	}

	if err == io.EOF {
		return l, err
	}

	return l, nil
}
