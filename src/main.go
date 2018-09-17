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

package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/piot/brisk-protocol-go/src/communication"
	"github.com/piot/brisk-protocol-go/src/connection"
	"github.com/piot/briskd-go/src/server"
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
	"github.com/piot/tend-go/src"
)

type FakeConnection struct {
}

func (FakeConnection) HandleStream(stream *instream.InStream, octetCount uint) error {
	return nil
}
func (FakeConnection) SendStream(sequenceID tend.SequenceID, stream *outstream.OutStream) (bool, error) {
	stream.WriteUint8(0x80)
	return false, nil
}
func (FakeConnection) Lost() error {
	return nil
}

func (FakeConnection) Dropped(sequenceID tend.SequenceID) {

}

func (FakeConnection) ReceivedByRemote(sequenceID tend.SequenceID) {

}

type FakeServer struct {
}

func (FakeServer) CreateConnection(id connection.ID) communication.Connection {
	return FakeConnection{}
}

func (FakeServer) Tick() {

}

func main() {
	color.Cyan("briskd example server 0.2\n")
	fakeServer := FakeServer{}
	instance, instanceErr := server.New(32002, fakeServer, true)
	if instanceErr != nil {
		fmt.Printf("Err:%v", instanceErr)
	}
	instance.Forever()
}
