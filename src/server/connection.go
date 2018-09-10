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
	"os"

	"github.com/piot/briskd-go/src/commands"
	"github.com/piot/briskd-go/src/communication"
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/briskd-go/src/endpoint"
	"github.com/piot/briskd-go/src/sequence"
	brisktime "github.com/piot/briskd-go/src/time"
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
	tend "github.com/piot/tend-go/src"
)

type Connection struct {
	endpoint             *endpoint.Endpoint
	nonce                uint32
	id                   connection.ID
	userConnection       communication.Connection
	server               *Server
	NextOutSequenceID    sequence.ID
	LastReceivedPacketAt int64
	runningStats         connection.RunningStats
	stats                connection.Stats
	debugDumpFile        *os.File
	tendOut              *tend.OutgoingLogic
	tendIn               *tend.IncomingLogic
}

func NewConnection(server *Server, id connection.ID, endpoint *endpoint.Endpoint, nonce uint32) *Connection {
	nextOutSequenceID, _ := sequence.NewID(sequence.MaxIDValue)
	c := &Connection{server: server, id: id, endpoint: endpoint, nonce: nonce, NextOutSequenceID: nextOutSequenceID, LastReceivedPacketAt: brisktime.MonotonicMilliseconds(), tendOut: tend.NewOutgoingLogic(), tendIn: tend.NewIncomingLogic()}
	if server.debugEnabled {
		c.debugDumpFile, _ = os.Create(fmt.Sprintf("connection_%d.ibd", id))
	}
	return c
}

func (c *Connection) SetUserConnection(userConnection communication.Connection) {
	c.userConnection = userConnection
}

func (c *Connection) SentPacket(octetCount uint) {
	c.runningStats.Sent.AddPackets(1, octetCount)
}

func (c *Connection) ReceivedPacket(octetCount uint) {
	c.runningStats.Received.AddPackets(1, octetCount)
}

func (c *Connection) CalculateStats(millisecondsDuration uint) {
	c.stats.SetFromRunningStats(&c.runningStats, millisecondsDuration)
	fmt.Printf("%v %v\n", c.id, &c.stats)
}

func (c *Connection) Addr() *endpoint.Endpoint {
	return c.endpoint
}

func (c *Connection) ID() connection.ID {
	return c.id
}

func (c *Connection) Send(stream *outstream.OutStream) (bool, error) {
	return c.userConnection.SendStream(c.tendOut.OutgoingSequenceID(), stream)
}

func (c *Connection) CanSendReliable() bool {
	return c.tendOut.CanIncrementOutgoingSequence()
}

func (c *Connection) IncreaseOutgoingSequenceID() commands.TendInfo {
	c.tendOut.IncreaseOutgoingSequenceID()
	return commands.TendInfo{
		PacketSequenceID:   c.tendOut.OutgoingSequenceID().Value(),
		ReceivedSequenceID: c.tendIn.ReceivedHeader().SequenceID.Value(),
		ReceivedMask:       c.tendIn.ReceivedHeader().Mask.Bits()}
}

func (c *Connection) writePacket(cmd uint8, monotonicTimeMs int64, b []byte) {
	o := c.debugDumpFile
	s := outstream.New()
	s.WriteUint8(cmd)
	s.WriteUint16(uint16(len(b)))
	s.WriteUint64(uint64(monotonicTimeMs))
	s.WriteOctets(b)
	o.Write(s.Octets())
	o.Sync()
}

func (c *Connection) DebugIncomingPacket(b []byte, monotonicTimeMs int64) {
	if c.debugDumpFile != nil {
		c.writePacket(0x01, monotonicTimeMs, b)
	}
}

func (c *Connection) DebugOutgoingPacket(b []byte, monotonicTimeMs int64) {
	if c.debugDumpFile != nil {
		c.writePacket(0x81, monotonicTimeMs, b)
	}
}

func (c *Connection) handleTend() error {
	for c.tendOut.QueueCount() > 0 {
		status, statusErr := c.tendOut.Dequeue()
		if statusErr != nil {
			return statusErr
		}
		if status.WasDelivered {
			c.userConnection.ReceivedByRemote(status.SequenceID)
		} else {
			c.userConnection.Dropped(status.SequenceID)
		}
	}

	return nil
}

func (c *Connection) handleStream(stream *instream.InStream, octetCount uint) error {
	//fmt.Printf("<< %v %v\n", c, stream)
	tendInfo, tendErr := commands.TendDeserialize(stream)
	if tendErr != nil {
		return tendErr
	}
	c.tendIn.ReceivedToUs(tend.NewSequenceID(tendInfo.PacketSequenceID))
	c.tendOut.ReceivedByRemote(tend.Header{SequenceID: tend.NewSequenceID(tendInfo.ReceivedSequenceID), Mask: tend.NewReceiveMask(tendInfo.ReceivedMask)})
	handleTendErr := c.handleTend()
	if handleTendErr != nil {
		return handleTendErr
	}
	userErr := c.userConnection.HandleStream(stream, octetCount)
	if userErr != nil {
		fmt.Printf("error:%v\n", userErr)
	}
	return nil
}

func (c *Connection) Lost() {
	fmt.Printf("Connection Lost %v\n", c)
	c.userConnection.Lost()
	if c.debugDumpFile != nil {
		c.debugDumpFile.Close()
		c.debugDumpFile = nil
	}
}

func (c *Connection) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[connection ")
	buffer.WriteString(fmt.Sprintf("id:%d", c.id))
	buffer.WriteString("]")
	return buffer.String()
}
