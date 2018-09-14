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
	"fmt"

	"net"
	"time"

	"github.com/piot/brisk-protocol-go/src/connection"
	"github.com/piot/brisk-protocol-go/src/endpoint"
)

type Server struct {
	connection *net.UDPConn

	//userServer communication.Server

	lastTimeStatsCalculatedAt int64
	debugEnabled              bool
	incomingHandler           *connection.IncomingHandler
}

// New : Creates a new server
func New(userServer connection.UserServer, enableDebug bool) *Server {
	s := &Server{debugEnabled: enableDebug}
	incomingHandler := connection.NewIncomingHandler(s, userServer)
	s.incomingHandler = incomingHandler
	return s
}

func (s *Server) WriteToUDP(addr *endpoint.Endpoint, octets []byte) {
	s.connection.WriteToUDP(octets, addr.UDPAddr())
}

func (s *Server) handleIncomingUDP() {
	for {
		buf := make([]byte, 1800)
		n, addr, err := s.connection.ReadFromUDP(buf)
		packet := buf[0:n]

		addrEndpoint := endpoint.New(addr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		packetErr := s.incomingHandler.HandlePacket(packet, &addrEndpoint)
		if packetErr != nil {
			fmt.Printf("Problem with packet:%s\n", packetErr)
		}
	}
}

func (s *Server) tick() error {
	s.incomingHandler.Update()
	//s.userServer.Tick()
	/*
		if resultErr != nil {
			fmt.Printf("Error: %s\n", resultErr)
		}
		return resultErr
	*/
	return nil
}

func (s *Server) start(ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			err := s.tick()
			if err != nil {
				fmt.Printf("Start err %s \n", err)
			}
		}
	}()
}

func (s *Server) Forever(listenPort int) error {
	portString := fmt.Sprintf(":%d", listenPort)
	serverAddress, err := net.ResolveUDPAddr("udp", portString)
	if err != nil {
		return fmt.Errorf("Error:%v ", err)
	}
	serverConnection, err := net.ListenUDP("udp", serverAddress)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Listening to %s\n", portString)

	go s.handleIncomingUDP()
	//defer serverConnection.Close()
	ticker := time.NewTicker(time.Millisecond * 33)

	s.connection = serverConnection
	s.start(ticker)

	select {}

	return nil
}
