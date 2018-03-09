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

package endpoint

import (
	"bytes"
	"fmt"
	"net"
)

type Endpoint struct {
	addr *net.UDPAddr
}

// New : Creates an endpoint
func New(addr *net.UDPAddr) Endpoint {
	point := Endpoint{addr: addr}
	return point
}

func (self *Endpoint) UDPAddr() *net.UDPAddr {
	return self.addr
}

func (addr *Endpoint) Equal(addr2 *Endpoint) bool {
	udp := addr.addr
	udp2 := addr2.addr
	return udp.IP.Equal(udp2.IP) && udp.Port == udp2.Port
}

func (e *Endpoint) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[endpoint ")
	buffer.WriteString(fmt.Sprintf("addr:%s port:%d", e.addr.IP, e.addr.Port))
	buffer.WriteString("]")
	return buffer.String()
}
