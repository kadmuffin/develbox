// Copyright 2022 Kevin Ledesma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Wrapper around net package to make it convinient to use for
// UNIX domain sockets.
package socket

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/kpango/glg"
)

type Socket struct {
	Path       string
	Listener   net.Listener
	Connection net.Conn
}

// Returns a new socket instance.
// To create the socker use the Create method.
func New(path string) *Socket {
	glg.Debug("Creating new socket instance")
	return &Socket{Path: path}
}

// Checks if the socket exists.
func (s *Socket) Exists() bool {
	_, err := os.Stat(s.Path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// Creates a socket.
func (s *Socket) Create() error {
	glg.Debug("Creating socket...")
	// Create socket directory if it doesn't exist
	dir := filepath.Dir(s.Path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		glg.Debugf("Creating directory %s", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	// Create socket
	l, err := net.Listen("unix", s.Path)
	if err != nil {
		return err
	}
	s.Listener = l
	return nil
}

// Connects to the socket.
// Returns a net.Conn.
func (s *Socket) Connect() error {
	glg.Debug("Connecting to socket...")
	conn, err := net.Dial("unix", s.Path)
	if err != nil {
		return err
	}
	s.Connection = conn
	return nil
}

// Closes the socket.
func (s *Socket) Close() error {
	return s.Listener.Close()
}

// Only closes the connection.
func (s *Socket) CloseConnection() error {
	return s.Connection.Close()
}

// Accepts a new connection.
// Returns a net.Conn.
func (s *Socket) Accept() (net.Conn, error) {
	return s.Listener.Accept()
}

// Waits for a connection.
// When a connection is made, it will run the callback function
// with the connection as argument.
func (s *Socket) Listen(cb func()) {
	glg.Debug("Listening for connections...")
	for {
		conn, err := s.Accept()
		if err != nil {
			glg.Fatal(err)
		}
		s.Connection = conn
		go cb()
	}
}

// Sends a message to the socket.
// Returns the amount of bytes sent.
func (s *Socket) Send(msg string) (int, error) {
	return s.Connection.Write([]byte(msg))
}

// Receives a message from the socket.
// Returns the message received.
func (s *Socket) Receive() ([]byte, error) {
	var buf [1024]byte
	n, err := s.Connection.Read(buf[:])
	if err != nil {
		return []byte{}, err
	}
	return buf[:n], nil
}

// Receives JSON from the socket.
// Returns an unmarshaled JSON object using the provided interface.
func (s *Socket) ReceiveJSON(v any) error {
	dec := json.NewDecoder(s.Connection)
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

// Sends JSON to the socket.
// Returns the amount of bytes sent.
func (s *Socket) SendJSON(v any) (int, error) {
	enc := json.NewEncoder(s.Connection)
	if err := enc.Encode(v); err != nil {
		return 0, err
	}
	return 0, nil
}

// Returns a reader for the socket.
func (s *Socket) Reader() (io.Reader, error) {
	return s.Connection, nil
}

// Returns a writer for the socket.
func (s *Socket) Writer() (io.Writer, error) {
	return s.Connection, nil
}
