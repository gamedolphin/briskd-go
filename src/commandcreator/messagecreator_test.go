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
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/fluxd-go/src/commands"
	"testing"
)

func TestChallengeResponseMessage(t *testing.T) {
	octets := []byte{0x02, 0xca, 0xfe, 0xde, 0xad, 0xc0, 0xde, 0xff}
	stream := instream.New(octets)

	msg, err := CreateMessage(&stream)
	if err != nil {
		t.Error(err)
	}

	responseMsg := msg.(*commands.ChallengeResponseMessage)

	if responseMsg.ClientDeterminedNonce != 0xcafedead {
		t.Errorf("Not correct nonce:%X", responseMsg.ClientDeterminedNonce)
	}
}
