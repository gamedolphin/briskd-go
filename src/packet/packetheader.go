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
package packet

import (
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/briskd-go/src/sequence"
)

type PacketHeader struct {
	Mode         Mode
	Sequence     sequence.ID
	ConnectionID connection.ID
}

func ReadHeader(stream *instream.InStream) (PacketHeader, error) {
	modeValue, err := stream.ReadUint8()
	if err != nil {
		return PacketHeader{}, err
	}
	s, err := stream.ReadUint8()
	if err != nil {
		return PacketHeader{}, err
	}
	connectionID, err := stream.ReadUint16()
	if err != nil {
		return PacketHeader{}, err
	}
	sequenceID, _ := sequence.NewID(sequence.IDType(s))
	return PacketHeader{Mode: Mode(modeValue), Sequence: sequenceID, ConnectionID: connection.ID(connectionID)}, nil
}

func WriteHeader(stream *outstream.OutStream, header *PacketHeader) {
	stream.WriteUint8(uint8(header.Mode))
	stream.WriteUint8(uint8(header.Sequence.Raw()))
	stream.WriteUint16(uint16(header.ConnectionID))
}
