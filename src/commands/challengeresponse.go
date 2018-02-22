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
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/briskd-go/src/packet"
)

type ChallengeResponseMessage struct {
	ClientDeterminedNonce uint32
	ConnectionID          connection.ID
}

func NewChallengeResponseMessage(clientDeterminedNonce uint32, connectionID connection.ID) *ChallengeResponseMessage {
	return &ChallengeResponseMessage{ClientDeterminedNonce: clientDeterminedNonce, ConnectionID: connectionID}
}

func (self *ChallengeResponseMessage) Serialize(stream *outstream.OutStream) error {
	stream.WriteUint32(self.ClientDeterminedNonce)
	stream.WriteUint16(uint16(self.ConnectionID))

	return nil
}

func (self *ChallengeResponseMessage) Deserialize(stream *instream.InStream) error {
	var err1 error
	self.ClientDeterminedNonce, err1 = stream.ReadUint32()
	if err1 != nil {
		return err1
	}
	id, err1 := stream.ReadUint16()
	if err1 != nil {
		return err1
	}

	self.ConnectionID = connection.ID(id)
	if err1 != nil {
		return err1
	}

	return nil
}

func (self *ChallengeResponseMessage) Command() packet.PacketCmd {
	return packet.OobPacketTypeChallengeResponse
}

func (self *ChallengeResponseMessage) String() string {
	return fmt.Sprintf("[ChallengeResponse Nonce %X Conn %d]", self.ClientDeterminedNonce, self.ConnectionID)
}
