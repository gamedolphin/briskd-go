/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package server

import (
	"bytes"
	"fmt"

	"github.com/piot/briskd-go/src/communication"
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/briskd-go/src/endpoint"
	"github.com/piot/briskd-go/src/sequence"
	brisktime "github.com/piot/briskd-go/src/time"
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
)

type Connection struct {
	endpoint             *endpoint.Endpoint
	nonce                uint32
	id                   connection.ID
	userConnection       communication.Connection
	server               *Server
	NextOutSequenceID    sequence.ID
	LastReceivedPacketAt int64
}

func NewConnection(server *Server, id connection.ID, endpoint *endpoint.Endpoint, nonce uint32) *Connection {
	nextOutSequenceID, _ := sequence.NewID(sequence.MaxIDValue)
	c := &Connection{server: server, id: id, endpoint: endpoint, nonce: nonce, NextOutSequenceID: nextOutSequenceID, LastReceivedPacketAt: brisktime.MonotonicMilliseconds()}
	return c
}

func (c *Connection) SetUserConnection(userConnection communication.Connection) {
	c.userConnection = userConnection
}

func (c *Connection) Addr() *endpoint.Endpoint {
	return c.endpoint
}

func (c *Connection) ID() connection.ID {
	return c.id
}

func (c *Connection) Send(stream *outstream.OutStream) error {
	return c.userConnection.SendStream(stream)
}

func (c *Connection) handleStream(stream *instream.InStream) error {
	//fmt.Printf("<< %v %v\n", c, stream)
	userErr := c.userConnection.HandleStream(stream)
	if userErr != nil {
		fmt.Printf("error:%v\n", userErr)
	}
	return nil
}

func (c *Connection) Lost() {
	fmt.Printf("Connection Lost %v\n", c)
	c.userConnection.Lost()
}

func (c *Connection) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[connection ")
	buffer.WriteString(fmt.Sprintf("id:%d", c.id))
	buffer.WriteString("]")
	return buffer.String()
}
