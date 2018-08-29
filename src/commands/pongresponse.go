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

package commands

import (
	"fmt"

	"github.com/piot/briskd-go/src/packet"
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
)

type PongResponse struct {
	EchoedTime uint64
	LocalTime  uint64
	Info       TendInfo
}

func NewPongResponse(echoedTime uint64, localTime uint64, info TendInfo) *PongResponse {
	return &PongResponse{EchoedTime: echoedTime, LocalTime: localTime, Info: info}
}

func (c *PongResponse) Serialize(stream *outstream.OutStream) error {
	stream.WriteUint64(c.EchoedTime)
	stream.WriteUint64(c.LocalTime)
	c.Info.Serialize(stream)
	return nil
}

func (c *PongResponse) Deserialize(stream *instream.InStream) error {
	stream.ReadUint64()
	stream.ReadUint64()
	c.Info, _ = TendDeserialize(stream)
	return nil
}

func (c *PongResponse) Command() packet.PacketCmd {
	return packet.OobPacketTypePongResponse
}

func (c *PongResponse) String() string {
	return fmt.Sprintf("[PongResponse echo %v local %v tend: %v", c.EchoedTime, c.LocalTime, c.Info)
}
