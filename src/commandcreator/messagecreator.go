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

package commandcreator

import (
	"fmt"

	"github.com/piot/briskd-go/src/commands"
	"github.com/piot/briskd-go/src/message"
	"github.com/piot/briskd-go/src/packet"
	"github.com/piot/brook-go/src/instream"
)

func createMessageFromStream(oobType packet.PacketCmd) message.Message {
	switch oobType {
	case packet.OobPacketTypeChallenge:
		return &commands.ChallengeMessage{}
	case packet.OobPacketTypeChallengeResponse:
		return &commands.ChallengeResponseMessage{}
	case packet.OobPacketTypeTimeSyncRequest:
		return &commands.TimeSyncRequest{}
	}

	return nil
}

func CreateMessage(stream *instream.InStream) (message.Message, error) {
	packetValue, packetValueErr := stream.ReadUint8()
	if packetValueErr != nil {
		return nil, packetValueErr
	}
	oobType := packet.PacketCmd(packetValue)

	msg := createMessageFromStream(oobType)
	if msg == nil {
		return nil, fmt.Errorf("illegal message type:%02X", packetValue)
	}
	err := msg.Deserialize(stream)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
